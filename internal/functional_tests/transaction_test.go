//go:build functional
// +build functional

package functionaltests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fbveronez/go-test-case-api/internal/handlers"
	"github.com/fbveronez/go-test-case-api/internal/model"
	"github.com/fbveronez/go-test-case-api/internal/repository"
	"github.com/fbveronez/go-test-case-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func transactionRouter(t *testing.T) *gin.Engine {
	db := ConnectTestDB()
	db.Exec("DELETE FROM transactions")

	t.Cleanup(func() {
		db.Exec("DELETE FROM transactions")
	})
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	router := gin.New()

	router.POST("/transactions", transactionHandler.CreateTransaction)
	router.GET("/transactions/:id", transactionHandler.GetTransactionsByAccountID)
	router.PUT("/transactions/:id", transactionHandler.UpdateTransaction)

	return router
}

func TestCreateTransactionSuccess(t *testing.T) {
	router := transactionRouter(t)

	payload := model.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          22.3,
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.Transaction
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, payload.Amount, resp.Amount)
	assert.NotZero(t, resp.AccountID)
}

func TestCreateAccountMissingRequiredField(t *testing.T) {
	router := transactionRouter(t)

	payload := model.CreateTransactionRequest{
		AccountID: 1,
		Amount:    33.2,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Contains(t, resp["error"], "Missing required parameter")
}

func TestGetAllTransactionsByAccountByID(t *testing.T) {
	router := transactionRouter(t)

	payload := model.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 2,
		Amount:          22.2,
	}
	body, _ := json.Marshal(payload)

	//create first record
	firstReq, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	firstReq.Header.Set("Content-Type", "application/json")
	wCreate1 := httptest.NewRecorder()
	router.ServeHTTP(wCreate1, firstReq)
	assert.Equal(t, http.StatusCreated, wCreate1.Code)

	//create second record
	secondReq, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	secondReq.Header.Set("Content-Type", "application/json")
	wCreate2 := httptest.NewRecorder()

	router.ServeHTTP(wCreate2, secondReq)
	assert.Equal(t, http.StatusCreated, wCreate2.Code)

	reqGet, _ := http.NewRequest(http.MethodGet, "/transactions/1", nil)
	wGet := httptest.NewRecorder()
	router.ServeHTTP(wGet, reqGet)

	assert.Equal(t, http.StatusOK, wGet.Code)

	var transactions []model.Transaction
	assert.NoError(t, json.Unmarshal(wGet.Body.Bytes(), &transactions))
	assert.Len(t, transactions, 2)
	assert.Equal(t, uint64(1), transactions[0].AccountID)
	assert.Equal(t, uint64(2), transactions[1].OperationTypeID)
}

func TestUpdateTransactionsByID(t *testing.T) {
	router := transactionRouter(t)

	payload := model.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 2,
		Amount:          22.2,
	}
	body, _ := json.Marshal(payload)

	//create first record
	firstReq, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	firstReq.Header.Set("Content-Type", "application/json")
	wCreate1 := httptest.NewRecorder()
	router.ServeHTTP(wCreate1, firstReq)
	assert.Equal(t, http.StatusCreated, wCreate1.Code)

	var transaction model.Transaction
	assert.NoError(t, json.Unmarshal(wCreate1.Body.Bytes(), &transaction))
	//update record
	updateReq := model.UpdateTransactionRequest{
		Amount: 300.32,
	}
	body2, _ := json.Marshal(updateReq)

	secondReq, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/transactions/%d", transaction.TransactionID), bytes.NewBuffer(body2))
	secondReq.Header.Set("Content-Type", "application/json")
	wUpdate := httptest.NewRecorder()

	router.ServeHTTP(wUpdate, secondReq)
	assert.Equal(t, http.StatusOK, wUpdate.Code)

}

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fbveronez/go-test-case-api/internal/model"
	"github.com/fbveronez/go-test-case-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	payload := model.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 3,
		Amount:          22.2,
	}

	mockService.On("CreateTransaction", mock.AnythingOfType("*model.Transaction")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*model.Transaction)
		arg.AccountID = 1
	})

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateTransaction(c)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.Transaction
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, resp.AccountID)
	assert.Equal(t, payload.Amount, resp.Amount)

	mockService.AssertExpectations(t)

}

func TestCreateTransactionBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(service.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer([]byte(`{"amount":`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateTransaction(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTransactionsByAccountIDSuccess(t *testing.T) {

	gin.SetMode(gin.TestMode)

	mockService := new(service.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	transactions := []model.Transaction{
		{
			AccountID:       1,
			TransactionID:   1,
			OperationTypeID: 2,
			Amount:          22.45,
		},
	}

	mockService.
		On("GetTransactionsByAccountID", uint64(1)).
		Return(transactions, nil)

	router := gin.New()
	router.GET("/transactions/:id", handler.GetTransactionsByAccountID)

	req, _ := http.NewRequest(http.MethodGet, "/transactions/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []model.Transaction
	err := json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(t, err)
	assert.Len(t, resp, 1)

	assert.EqualValues(t, 1, resp[0].AccountID)
	assert.EqualValues(t, 2, resp[0].OperationTypeID)
	assert.Equal(t, 22.45, resp[0].Amount)

	mockService.AssertExpectations(t)
}

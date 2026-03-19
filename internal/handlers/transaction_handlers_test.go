package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fbveronez/go-test-case-api/internal/mocks"
	"github.com/fbveronez/go-test-case-api/internal/model"
	"github.com/fbveronez/go-test-case-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateTransactionSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockTransactionService)
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

	mockService := new(mocks.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer([]byte(`{"amount":`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateTransaction(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateTransactionMissingRequiredParameter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	payload := model.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 3,
	}

	mockService.On("CreateTransaction", mock.Anything).
		Return(nil)

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateTransaction(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error": "Missing required parameter"}`, w.Body.String())
}

func TestCreateTransaction_AccountNotFound(t *testing.T) {
	mockService := new(mocks.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	body := `{
		"account_id": 1,
		"operation_type_id": 1,
		"amount": 100
	}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mockService.On("CreateTransaction", mock.Anything).
		Return(service.ErrAccountNotFound)

	handler.CreateTransaction(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "account")

}

func TestGetTransactionsByAccountIDSuccess(t *testing.T) {

	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockTransactionService)
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

func TestGetTransactionInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	req, _ := http.NewRequest(http.MethodPut, "/transactions/abc", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetTransactionsByAccountID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"Invalid account ID"}`, w.Body.String())
}

func TestGetTransactionNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	mockService.On("GetTransactionsByAccountID", uint64(1), mock.Anything).
		Return(nil, gorm.ErrRecordNotFound)

	req, _ := http.NewRequest(http.MethodPut, "/transactions/1", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	handler.GetTransactionsByAccountID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error":"no transactions found"}`, w.Body.String())
}

func TestUpdateTransactionInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	req, _ := http.NewRequest(http.MethodPut, "/transactions/abc", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UpdateTransaction(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"invalid transaction ID"}`, w.Body.String())
}

func TestUpdateTransactionInvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	req, _ := http.NewRequest(http.MethodPut, "/transactions/1", strings.NewReader("invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	handler.UpdateTransaction(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"invalid payload"}`, w.Body.String())
}

func TestUpdateTransactionNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	mockService.On("UpdateTransaction", uint64(1), mock.Anything).
		Return(nil, gorm.ErrRecordNotFound)

	body := `{"amount":100}`

	req, _ := http.NewRequest(http.MethodPut, "/transactions/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	handler.UpdateTransaction(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error":"transaction not found"}`, w.Body.String())
}

func TestUpdateTransactionInternal(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	mockService.On("UpdateTransaction", uint64(1), mock.Anything).
		Return(nil, errors.New("db error"))

	body := `{"amount":100}`

	req, _ := http.NewRequest(http.MethodPut, "/transactions/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	handler.UpdateTransaction(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":"internal server error"}`, w.Body.String())
}

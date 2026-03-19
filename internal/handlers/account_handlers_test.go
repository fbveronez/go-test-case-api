package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fbveronez/go-test-case-api/internal/mocks"
	"github.com/fbveronez/go-test-case-api/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateAccountSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockAccountService)
	handler := NewAccountHandler(mockService)

	payload := model.CreateAccountRequest{
		DocumentNumber: "12345678900",
	}

	mockService.On("CreateAccount", mock.AnythingOfType("*model.Account")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*model.Account)
		arg.AccountID = 1
	})

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateAccount(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.Account
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, resp.AccountID)
	assert.Equal(t, payload.DocumentNumber, resp.DocumentNumber)

	mockService.AssertExpectations(t)
}

func TestCreateAccountBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockAccountService)
	handler := NewAccountHandler(mockService)

	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer([]byte(`{"document_number":`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateAccount(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAccountByIDSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockAccountService)
	handler := NewAccountHandler(mockService)

	account := &model.Account{
		AccountID:      1,
		DocumentNumber: "12345678900",
	}

	mockService.On("GetAccountByID", uint(1)).Return(account, nil)

	req, _ := http.NewRequest(http.MethodGet, "/accounts/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.GetAccountByID(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.Account
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, resp.AccountID)
	assert.Equal(t, "12345678900", resp.DocumentNumber)

	mockService.AssertExpectations(t)
}

func TestGetAccountByIDNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockAccountService)
	handler := NewAccountHandler(mockService)

	mockService.On("GetAccountByID", uint(2)).Return((*model.Account)(nil), gorm.ErrRecordNotFound)

	req, _ := http.NewRequest(http.MethodGet, "/accounts/2", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "2"}}

	handler.GetAccountByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAccountByIDBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockAccountService)
	handler := NewAccountHandler(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/accounts/abc", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler.GetAccountByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAccountByIDInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockAccountService)
	handler := NewAccountHandler(mockService)

	mockService.On("GetAccountByID", uint(3)).Return((*model.Account)(nil), assert.AnError)

	req, _ := http.NewRequest(http.MethodGet, "/accounts/3", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "3"}}

	handler.GetAccountByID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteByIDSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockAccountService)
	handler := NewAccountHandler(mockService)

	mockService.On("DeleteAccountByID", uint64(3)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/account/3", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "3"}}

	handler.DeleteAccount(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteAccountNotFound(t *testing.T) {
	mockService := new(mocks.MockAccountService)
	handler := NewAccountHandler(mockService)

	mockService.On("DeleteAccountByID", uint64(1)).
		Return(gorm.ErrRecordNotFound)

	req, _ := http.NewRequest(http.MethodDelete, "/accounts/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.DeleteAccount(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteAccountInternalError(t *testing.T) {
	mockService := new(mocks.MockAccountService)
	handler := NewAccountHandler(mockService)

	mockService.On("DeleteAccountByID", uint64(1)).
		Return(errors.New("db error"))

	req, _ := http.NewRequest(http.MethodDelete, "/accounts/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.DeleteAccount(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

//go:build functional
// +build functional

package functionaltests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/fbveronez/go-test-case-api/internal/handlers"
	"github.com/fbveronez/go-test-case-api/internal/model"
	"github.com/fbveronez/go-test-case-api/internal/repository"
	"github.com/fbveronez/go-test-case-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func accountRouter() *gin.Engine {
	db := ConnectTestDB()
	db.Exec("DELETE FROM accounts")

	accountRepo := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepo)
	accountHandler := handlers.NewAccountHandler(accountService)

	router := gin.New()

	router.POST("/accounts", accountHandler.CreateAccount)
	router.GET("/accounts/:id", accountHandler.GetAccountByID)
	router.DELETE("/accounts/:id", accountHandler.DeleteAccount)

	return router
}

func TestCreateAccountSuccess(t *testing.T) {
	router := accountRouter()

	payload := model.CreateAccountRequest{DocumentNumber: "12345678901"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.Account
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, payload.DocumentNumber, resp.DocumentNumber)
	assert.NotZero(t, resp.AccountID)
}

func TestCreateAccountDuplicate(t *testing.T) {
	router := accountRouter()

	payload := model.CreateAccountRequest{DocumentNumber: "98765432100"}
	body, _ := json.Marshal(payload)

	req1, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	var resp1 model.Account
	assert.NoError(t, json.Unmarshal(w1.Body.Bytes(), &resp1))
	assert.Equal(t, payload.DocumentNumber, resp1.DocumentNumber)
	assert.NotZero(t, resp1.AccountID)

	req2, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Dependendo de como você implementou duplicidade, pode ser 500 ou 409
	assert.Equal(t, http.StatusInternalServerError, w2.Code)

	// Optional: verificar mensagem de erro
	expectedMsg := "document number already exists"
	assert.Contains(t, w2.Body.String(), expectedMsg)
}

func TestCreateAccountDocumentNumberEmpty(t *testing.T) {
	router := accountRouter()

	payload := model.CreateAccountRequest{DocumentNumber: ""}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Contains(t, resp["error"], "document number is required")
}

func TestCreateAccountEmptyPayload(t *testing.T) {
	router := accountRouter()

	payload := map[string]string{}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Contains(t, resp["error"], "document number is required")
}

func TestGetAccountByID(t *testing.T) {
	router := accountRouter()

	payload := model.CreateAccountRequest{DocumentNumber: "55555555555"}
	body, _ := json.Marshal(payload)

	reqCreate, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	router.ServeHTTP(wCreate, reqCreate)

	assert.Equal(t, http.StatusCreated, wCreate.Code)

	var createdAccount model.Account
	assert.NoError(t, json.Unmarshal(wCreate.Body.Bytes(), &createdAccount))
	assert.NotZero(t, createdAccount.AccountID)

	reqGet, _ := http.NewRequest(http.MethodGet, "/accounts/"+strconv.FormatUint(createdAccount.AccountID, 10), nil)
	wGet := httptest.NewRecorder()
	router.ServeHTTP(wGet, reqGet)

	assert.Equal(t, http.StatusOK, wGet.Code)

	var fetchedAccount model.Account
	assert.NoError(t, json.Unmarshal(wGet.Body.Bytes(), &fetchedAccount))
	assert.Equal(t, createdAccount.AccountID, fetchedAccount.AccountID)
	assert.Equal(t, createdAccount.DocumentNumber, fetchedAccount.DocumentNumber)

	reqGetNotFound, _ := http.NewRequest(http.MethodGet, "/accounts/999999", nil)
	wGetNotFound := httptest.NewRecorder()
	router.ServeHTTP(wGetNotFound, reqGetNotFound)

	assert.Equal(t, http.StatusNotFound, wGetNotFound.Code)

	var respNotFound map[string]string
	assert.NoError(t, json.Unmarshal(wGetNotFound.Body.Bytes(), &respNotFound))
	assert.Contains(t, respNotFound["error"], "not found")
}

func TestDeleteAccountByID(t *testing.T) {
	router := accountRouter()

	payload := model.CreateAccountRequest{DocumentNumber: "55555555555"}
	body, _ := json.Marshal(payload)

	reqCreate, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	router.ServeHTTP(wCreate, reqCreate)

	assert.Equal(t, http.StatusCreated, wCreate.Code)

	var createdAccount model.Account
	assert.NoError(t, json.Unmarshal(wCreate.Body.Bytes(), &createdAccount))

	reqDel, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/accounts/%d", createdAccount.AccountID), nil)
	wDel := httptest.NewRecorder()
	router.ServeHTTP(wDel, reqDel)

	assert.Equal(t, http.StatusNoContent, wDel.Code)
}

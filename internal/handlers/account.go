// Package handlers
package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fbveronez/go-test-case-api/internal/model"
	"github.com/fbveronez/go-test-case-api/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountHandler struct {
	Service service.AccountService
}

func NewAccountHandler(s service.AccountService) *AccountHandler {
	return &AccountHandler{Service: s}
}

// CreateAccount godoc
// @Summary Create account
// @Description Create a new account
// @Tags Accounts
// @Accept json
// @Produce json
// @Param account body model.CreateAccountRequest true "Account payload"
// @Success 201 {object} model.Account
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req model.CreateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document number is required"})
		return
	}

	account := model.Account{
		DocumentNumber: req.DocumentNumber,
	}

	if err := h.Service.CreateAccount(&account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccountByID godoc
// @Summary Get account by ID
// @Description Retrieve an account by its ID
// @Tags Accounts
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} model.Account
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	idParam := c.Param("id")

	accountID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := h.Service.GetAccountByID(uint(accountID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// DeleteAccount godoc
// @Summary Delete an account
// @Description Delete an account by ID
// @Tags Accounts
// @Produce json
// @Param id path int true "Account ID"
// @Success 204 {string} string "No Content"
// @Router /accounts/{id} [delete]
func (h *AccountHandler) DeleteAccount(c *gin.Context) {

	idParam := c.Param("id")

	accountID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid account ID",
		})
		return
	}

	err = h.Service.DeleteAccountByID(accountID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "account not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

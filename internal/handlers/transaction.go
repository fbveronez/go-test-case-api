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

type TransactionHandler struct {
	Service service.TransactionService
}

func NewTransactionHandler(s service.TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: s}
}

// CreateTransaction godoc
// @Summary Create a transaction
// @Description Creates a new financial transaction for an account
// @Tags Transactions
// @Accept json
// @Produce json
// @Param transaction body model.CreateTransactionRequest true "Transaction payload"
// @Success 201 {object} model.Transaction
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {

	var req model.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameter"})
		return
	}

	transaction := model.Transaction{
		AccountID:       req.AccountID,
		OperationTypeID: req.OperationTypeID,
		Amount:          req.Amount,
	}

	if err := h.Service.CreateTransaction(&transaction); err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// GetTransactionsByAccountID godoc
// @Summary Get transactions by ID
// @Description Retrieve all transactions by accountID
// @Tags Transactions
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} []model.Transaction
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransactionsByAccountID(c *gin.Context) {
	idParam := c.Param("id")

	accountID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	transactions, err := h.Service.GetTransactionsByAccountID(uint64(accountID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no transactions found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// UpdateTransaction godoc
// @Summary Update a transaction
// @Description Update an existing transaction by transactionID
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param transaction body model.UpdateTransactionRequest true "Transaction update payload"
// @Success 200 {object} model.Transaction
// @Router /transactions/{id} [put]
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {

	idParam := c.Param("id")

	transactionID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	var req model.UpdateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	update := model.Transaction{
		Amount: req.Amount,
	}

	transaction, err := h.Service.UpdateTransaction(transactionID, &update)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

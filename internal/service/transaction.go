// Package service.
package service

import (
	"errors"
	"strings"

	"github.com/fbveronez/go-test-case-api/internal/model"
	"github.com/fbveronez/go-test-case-api/internal/repository"
)

var ErrOperationTypeNotValid = errors.New("operation type not valid")

// TransactionService defines the business operations related to transactions.
type TransactionService interface {
	CreateTransaction(transaction *model.Transaction) error
	GetTransactionsByAccountID(accountID uint64) ([]model.Transaction, error)
	UpdateTransaction(transactionID uint64, update *model.Transaction) (*model.Transaction, error)
}

// transactionService is the concrete implementation of TransactionService.
type transactionService struct {
	repo repository.TransactionRepository
}

// NewTransactionService creates a new TransactionService with the provided repository.
func NewTransactionService(r repository.TransactionRepository) TransactionService {
	return &transactionService{repo: r}
}

// CreateTransaction inserts a new transaction into the database.
func (s *transactionService) CreateTransaction(t *model.Transaction) error {
	err := s.repo.Create(t)

	if err != nil {
		if strings.Contains(err.Error(), "fk_account") {
			return ErrAccountNotFound
		}

		if strings.Contains(err.Error(), "fk_operation_type") {
			return ErrOperationTypeNotValid
		}

		return err
	}
	return nil
}

// GetTransactionsByAccountID retrieves all transactions for the specified account.
func (s *transactionService) GetTransactionsByAccountID(accountID uint64) ([]model.Transaction, error) {
	return s.repo.GetAllByAccountID(accountID)
}

// UpdateTransaction updates an existing transaction by its ID.
func (s *transactionService) UpdateTransaction(transactionID uint64, update *model.Transaction) (*model.Transaction, error) {
	return s.repo.UpdateByTransactionID(transactionID, update)
}

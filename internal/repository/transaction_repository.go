// Package repository provides persistence layer implementations for the application.
package repository

import (
	"github.com/fbveronez/go-test-case-api/internal/model"
	"gorm.io/gorm"
)

// TransactionRepository defines the interface for managing transactions in the database.
type TransactionRepository interface {
	Create(transaction *model.Transaction) error
	GetAllByAccountID(accountID uint64) ([]model.Transaction, error)
	UpdateByTransactionID(transactionID uint64, update *model.Transaction) (*model.Transaction, error)
}

// transactionRepository is the concrete implementation of TransactionRepository
// using GORM as the underlying ORM.
type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new TransactionRepository using the provided GORM database connection.
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// Create inserts a new transaction record into the database.
func (r *transactionRepository) Create(t *model.Transaction) error {
	return r.db.Create(t).Error
}

// GetAllByAccountID fetches all transactions for the given account ID.
func (r *transactionRepository) GetAllByAccountID(accountID uint64) ([]model.Transaction, error) {
	var transactions []model.Transaction

	result := r.db.Where("account_id = ?", accountID).Find(&transactions)

	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

// UpdateByTransactionID updates the transaction record identified by transactionID
// with the provided update values and returns the updated transaction.
func (r *transactionRepository) UpdateByTransactionID(transactionID uint64, update *model.Transaction) (*model.Transaction, error) {

	var transaction model.Transaction

	if err := r.db.First(&transaction, transactionID).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&transaction).Updates(update).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

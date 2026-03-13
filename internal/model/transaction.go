// Package model contains the data models used in the application.
package model

import (
	"time"
)

// UpdateTransactionRequest defines the payload for updating a transaction.
// Only the amount can be updated via this request.
type UpdateTransactionRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

// CreateTransactionRequest defines the payload for creating a new transaction.
type CreateTransactionRequest struct {
	AccountID       uint64  `json:"account_id" binding:"required" example:"1"`
	OperationTypeID uint64  `json:"operation_type_id" binding:"required" example:"1"`
	Amount          float64 `json:"amount" binding:"required" example:"22.2"`
}

// Transaction represents the database.
type Transaction struct {
	TransactionID   uint64    `gorm:"primaryKey;column:transaction_id" json:"transaction_id"`
	AccountID       uint64    `gorm:"not null;column:account_id" json:"account_id"`
	OperationTypeID uint64    `gorm:"not null;column:operation_type_id" json:"operation_type_id"`
	Amount          float64   `gorm:"type:numeric(12,2);not null" json:"amount"`
	EventDate       time.Time `gorm:"not null;default:now()" json:"event_date"`
}

// TableName overrides the default GORM table name for Transaction.
// Ensures the struct maps to the "transactions" table in the database.
func (Transaction) TableName() string {
	return "transactions"
}

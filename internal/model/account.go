// Package model contains the data models used in the application.
package model

import (
	"time"
)

// CreateAccountRequest represents the payload required to create a new account.
type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" binding:"required" example:"12345678900"`
}

// Account represents the database.
type Account struct {
	AccountID      uint64    `gorm:"primaryKey;column:account_id" json:"account_id"`
	DocumentNumber string    `gorm:"size:20;not null;unique;column:document_number" json:"document_number"`
	CreatedAt      time.Time `gorm:"not null;default:now()" json:"created_at"`
}

// TableName overrides the default table name for GORM.
// This ensures the Account struct maps to the "accounts" table in the database.
func (Account) TableName() string {
	return "accounts"
}

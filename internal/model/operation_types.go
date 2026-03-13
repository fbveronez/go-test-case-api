// Package model contains the data models used in the application.
package model

import (
	"time"
)

// OperationType represents the type of a financial transaction operation.
// Each transaction must be associated with one operation type (e.g., purchase, payment, withdrawal).
type OperationType struct {
	OperationTypeID uint64    `gorm:"primaryKey;column:operation_type_id" json:"operation_type_id"`
	Description     string    `gorm:"size:100;not null;unique;column:description" json:"description"`
	CreatedAt       time.Time `gorm:"not null;default:now()" json:"created_at"`
}

// TableName overrides the default table name for GORM.
// This ensures the OperationType struct maps to the "operation_types" table in the database.
func (OperationType) TableName() string {
	return "operation_types"
}

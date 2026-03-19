// Package repository
package repository

import (
	"github.com/fbveronez/go-test-case-api/internal/model"
	"gorm.io/gorm"
)

// AccountRepository defines the persistence layer for account operations.
type AccountRepository interface {
	Create(account *model.Account) error
	FindByID(id uint) (*model.Account, error)
	FindByDocumentNumber(document string) (*model.Account, error)
	DeleteByID(id uint64) error
	UpdateCredit(account *model.Account, amount float64) error
}

// accountRepository is the concrete implementation of AccountRepository
// using GORM as the underlying database ORM.
type accountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new AccountRepository instance
// with the provided GORM database connection.
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

// Create persists a new account in the database.
func (r *accountRepository) Create(t *model.Account) error {
	return r.db.Create(t).Error
}

// FindByID retrieves an account by its unique ID.
// Returns gorm.ErrRecordNotFound if the account does not exist.
func (r *accountRepository) FindByID(id uint) (*model.Account, error) {
	var account model.Account
	result := r.db.First(&account, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &account, nil
}

// FindByDocumentNumber retrieves an account using its document number.
// Returns gorm.ErrRecordNotFound if no account is found.
func (r *accountRepository) FindByDocumentNumber(document string) (*model.Account, error) {
	var account model.Account
	err := r.db.Where("document_number = ?", document).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// DeleteByID removes an account from the database by its ID.
// Returns gorm.ErrRecordNotFound if the account does not exist.
func (r *accountRepository) DeleteByID(id uint64) error {

	result := r.db.Delete(&model.Account{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *accountRepository) UpdateCredit(account *model.Account, amount float64) error {

	value := account.AvailableCreditLimit + amount

	if value <= 0 {
		return gorm.ErrInvalidTransaction
	}
	update := model.Account{
		AvailableCreditLimit: value,
	}
	if err := r.db.Model(&account).Updates(update).Error; err != nil {
		return err
	}
	return nil
}

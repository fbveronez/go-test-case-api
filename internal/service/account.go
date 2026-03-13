// Package service
package service

import (
	"errors"

	"github.com/fbveronez/go-test-case-api/internal/model"
	"github.com/fbveronez/go-test-case-api/internal/repository"
	"gorm.io/gorm"
)

// Predefined service errors.
var (
	ErrAccountNotFound     = errors.New("account not found")
	ErrDocumentAlreadyUsed = errors.New("document number already exists")
)

// AccountService defines the business operations related to accounts.
type AccountService interface {
	CreateAccount(account *model.Account) error
	GetAccountByID(id uint) (*model.Account, error)
	DeleteAccountByID(id uint64) error
}

// accountService is the concrete implementation of AccountService.
type accountService struct {
	repo repository.AccountRepository
}

// NewAccountService creates a new AccountService with the given repository.
func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{
		repo: repo,
	}
}

// CreateAccount creates a new account in the database.
func (s *accountService) CreateAccount(account *model.Account) error {

	if account.DocumentNumber == "" {
		return errors.New("document number is required")
	}

	existingAccount, err := s.repo.FindByDocumentNumber(account.DocumentNumber)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			existingAccount = nil
		} else {
			return err
		}
	}

	if existingAccount != nil {
		return ErrDocumentAlreadyUsed
	}

	return s.repo.Create(account)
}

// GetAccountByID retrieves an account by its ID.
func (s *accountService) GetAccountByID(id uint) (*model.Account, error) {

	account, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return account, nil
}

// DeleteAccountByID deletes an account by its ID.
func (s *accountService) DeleteAccountByID(id uint64) error {
	err := s.repo.DeleteByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	return nil
}

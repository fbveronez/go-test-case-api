// Package service
package service

import (
	"errors"
	"testing"

	"github.com/fbveronez/go-test-case-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockAccountRepository struct {
	account *model.Account
	err     error
	mock.Mock
}

func (m *mockAccountRepository) Create(account *model.Account) error {
	m.account = account
	return m.err
}

func (m *mockAccountRepository) FindByID(id uint) (*model.Account, error) {
	return m.account, m.err
}

func (m *mockAccountRepository) DeleteByID(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockAccountRepository) FindByDocumentNumber(document string) (*model.Account, error) {
	return m.account, m.err
}

func TestCreateAccountSuccess(t *testing.T) {

	repo := &mockAccountRepository{}

	service := NewAccountService(repo)

	account := &model.Account{
		DocumentNumber: "12345678900",
	}

	err := service.CreateAccount(account)

	assert.NoError(t, err)
}

func TestCreateAccountDuplicateDocument(t *testing.T) {

	repo := &mockAccountRepository{
		account: &model.Account{
			DocumentNumber: "12345678900",
		},
	}

	service := NewAccountService(repo)

	account := &model.Account{
		DocumentNumber: "12345678900",
	}

	err := service.CreateAccount(account)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrDocumentAlreadyUsed))
}

func TestCreateAccountWithoutDocument(t *testing.T) {

	repo := &mockAccountRepository{}

	service := NewAccountService(repo)

	account := &model.Account{}

	err := service.CreateAccount(account)

	assert.Error(t, err)
	assert.Equal(t, "document number is required", err.Error())
}

func TestGetAccountByIDSuccess(t *testing.T) {

	repo := &mockAccountRepository{
		account: &model.Account{
			AccountID:      1,
			DocumentNumber: "12345678900",
		},
	}

	service := NewAccountService(repo)

	account, err := service.GetAccountByID(1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), uint(account.AccountID))
}

func TestDeleteByIDSuccess(t *testing.T) {

	repo := new(mockAccountRepository)

	repo.On("DeleteByID", uint64(1)).Return(nil)

	service := NewAccountService(repo)

	err := service.DeleteAccountByID(1)

	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestDeleteByIDNotFound(t *testing.T) {

	repo := new(mockAccountRepository)

	repo.On("DeleteByID", uint64(1)).Return(gorm.ErrRecordNotFound)

	service := NewAccountService(repo)

	err := service.DeleteAccountByID(1)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	repo.AssertExpectations(t)
}

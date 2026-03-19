// Package mocks
package mocks

import (
	"github.com/fbveronez/go-test-case-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) CreateAccount(account *model.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountService) GetAccountByID(id uint) (*model.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountService) DeleteAccountByID(accountID uint64) error {
	args := m.Called(accountID)
	return args.Error(0)
}

func (m *MockAccountService) UpdateCreditLimit(id uint64, amount float64) error {
	args := m.Called(id)
	return args.Error(0)
}

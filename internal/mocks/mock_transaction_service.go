package mocks

import (
	"github.com/fbveronez/go-test-case-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) CreateTransaction(t *model.Transaction) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTransactionService) GetTransactionsByAccountID(accountID uint64) ([]model.Transaction, error) {
	args := m.Called(accountID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *MockTransactionService) UpdateTransaction(transactionID uint64, update *model.Transaction) (*model.Transaction, error) {
	args := m.Called(transactionID, update)

	transaction, _ := args.Get(0).(*model.Transaction)

	return transaction, args.Error(1)
}

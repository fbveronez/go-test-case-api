// Package service
package service

import (
	"errors"
	"testing"
	"time"

	"github.com/fbveronez/go-test-case-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTransactionRepository struct {
	fail bool
	mock.Mock
}

func (m *mockTransactionRepository) Create(t *model.Transaction) error {

	args := m.Called(t)

	return args.Error(0)
}

func (m *mockTransactionRepository) GetAllByAccountID(id uint64) ([]model.Transaction, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *mockTransactionRepository) UpdateByTransactionID(transactionID uint64, update *model.Transaction) (*model.Transaction, error) {
	args := m.Called(transactionID, update)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.Transaction), args.Error(1)
}

func TestCreateTransactionSuccess(t *testing.T) {

	repo := &mockTransactionRepository{}

	repo.On("Create", mock.AnythingOfType("*model.Transaction")).Return(nil)

	service := NewTransactionService(repo)

	transaction := &model.Transaction{
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          100,
	}

	err := service.CreateTransaction(transaction)

	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestCreateTransactionError(t *testing.T) {

	repo := new(mockTransactionRepository)

	repo.On("Create", mock.Anything).
		Return(errors.New("failed to create transaction"))

	service := NewTransactionService(repo)

	transaction := &model.Transaction{
		AccountID: 1,
		Amount:    50.0,
		EventDate: time.Now(),
	}

	err := service.CreateTransaction(transaction)

	assert.Error(t, err)
	assert.Equal(t, "failed to create transaction", err.Error())

	repo.AssertExpectations(t)
}

func TestCreateTransactionAccountNotFound(t *testing.T) {

	repo := new(mockTransactionRepository)

	service := NewTransactionService(repo)

	transaction := &model.Transaction{
		AccountID:       999,
		OperationTypeID: 1,
		Amount:          100,
	}

	dbError := errors.New(
		"ERROR: insert or update on table \"transactions\" violates foreign key constraint \"fk_account\" (SQLSTATE 23503)",
	)

	repo.On("Create", mock.AnythingOfType("*model.Transaction")).Return(dbError)

	err := service.CreateTransaction(transaction)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrAccountNotFound))

	repo.AssertExpectations(t)
}

func TestCreateTransactionOperationTypeNotValid(t *testing.T) {

	repo := new(mockTransactionRepository)

	service := NewTransactionService(repo)

	transaction := &model.Transaction{
		AccountID:       999,
		OperationTypeID: 8,
		Amount:          100,
	}

	dbError := errors.New(
		"ERROR: insert or update on table \"transactions\" violates foreign key constraint \"fk_operation_types\" (SQLSTATE 23503)",
	)

	repo.On("Create", mock.AnythingOfType("*model.Transaction")).Return(dbError)

	err := service.CreateTransaction(transaction)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrOperationTypeNotValid))

	repo.AssertExpectations(t)
}

func TestGetTransactionsByAccountID(t *testing.T) {

	repo := &mockTransactionRepository{}

	service := NewTransactionService(repo)

	accountID := uint64(1)

	expectedTransactions := []model.Transaction{
		{
			TransactionID: 1,
			AccountID:     1,
			Amount:        100,
		},
		{
			TransactionID: 2,
			AccountID:     1,
			Amount:        50,
		},
	}

	repo.On("GetAllByAccountID", uint64(accountID)).Return(expectedTransactions, nil)

	transactions, err := service.GetTransactionsByAccountID(accountID)

	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Len(t, transactions, 2)
	assert.Equal(t, expectedTransactions, transactions)

	repo.AssertExpectations(t)
}

func TestUpdateTransaction_Success(t *testing.T) {

	repo := &mockTransactionRepository{}

	service := NewTransactionService(repo)

	transactionID := uint64(1)

	updatePayload := &model.Transaction{
		Amount: 100,
	}

	expectedTransaction := &model.Transaction{
		TransactionID: 1,
		AccountID:     1,
		Amount:        100,
	}

	repo.On("UpdateByTransactionID", transactionID, updatePayload).Return(expectedTransaction, nil)

	transaction, err := service.UpdateTransaction(transactionID, updatePayload)

	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, expectedTransaction.Amount, transaction.Amount)

	repo.AssertExpectations(t)
}

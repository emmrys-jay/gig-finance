package service

import (
	"github.com/emmrys-jay/gigmile/internal/models"
	"github.com/emmrys-jay/gigmile/internal/repository"
)

type TransactionService interface {
	GetTransactionsByCustomer(customerID int64) ([]*models.Transaction, error)
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *transactionService) GetTransactionsByCustomer(customerID int64) ([]*models.Transaction, error) {
	return s.transactionRepo.GetByCustomerID(customerID)
}

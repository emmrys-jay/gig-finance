package service

import (
	"github.com/emmrys-jay/gigmile/internal/models"
	"github.com/emmrys-jay/gigmile/internal/repository"
)

type AccountService interface {
	GetAccountByCustomer(customerID int64) (*models.Account, error)
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{
		accountRepo: accountRepo,
	}
}

func (s *accountService) GetAccountByCustomer(customerID int64) (*models.Account, error) {
	return s.accountRepo.GetByCustomerID(customerID)
}


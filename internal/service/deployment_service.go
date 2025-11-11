package service

import (
	"fmt"

	"github.com/emmrys-jay/gigmile/internal/models"
	"github.com/emmrys-jay/gigmile/internal/repository"
	"github.com/emmrys-jay/gigmile/internal/utils"
)

type DeploymentService interface {
	RecordDeployment(req *models.CreateDeploymentRequest) error
}

type deploymentService struct {
	customerRepo    repository.CustomerRepository
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func NewDeploymentService(
	customerRepo repository.CustomerRepository,
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
) DeploymentService {
	return &deploymentService{
		customerRepo:    customerRepo,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

const DeploymentAmount = 1000000.00 // 1 million

func (s *deploymentService) RecordDeployment(req *models.CreateDeploymentRequest) error {
	// Parse customer ID (remove GIG prefix if present)
	customerID, err := utils.ParseCustomerID(req.CustomerID)
	if err != nil {
		return fmt.Errorf("invalid customer_id: %w", err)
	}

	// Get customer
	_, err = s.customerRepo.GetByID(customerID)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}

	// Get customer's account
	account, err := s.accountRepo.GetByCustomerID(customerID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Create transaction
	if req.Description == "" {
		req.Description = "Deployment"
	}
	createTransactionReq := &models.CreateTransactionRequest{
		CustomerID:  customerID,
		AccountID:   account.ID,
		Reference:   req.Reference,
		Amount:      DeploymentAmount,
		Status:      models.PaymentStatusPending,
		Description: req.Description,
	}

	transaction, err := s.transactionRepo.Create(createTransactionReq)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Debit the account
	err = s.accountRepo.Debit(account.ID, transaction.ID, DeploymentAmount)
	if err != nil {
		return fmt.Errorf("failed to debit account: %w", err)
	}

	return nil
}

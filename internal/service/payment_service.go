package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/emmrys-jay/gigmile/internal/models"
	"github.com/emmrys-jay/gigmile/internal/repository"
	"github.com/emmrys-jay/gigmile/internal/utils"
)

type PaymentService interface {
	ProcessPaymentNotification(req *models.PaymentNotificationRequest) error
}

type paymentService struct {
	customerRepo    repository.CustomerRepository
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func NewPaymentService(
	customerRepo repository.CustomerRepository,
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
) PaymentService {
	return &paymentService{
		customerRepo:    customerRepo,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *paymentService) ProcessPaymentNotification(req *models.PaymentNotificationRequest) error {
	// Validate payment status
	if req.PaymentStatus != "COMPLETE" {
		return fmt.Errorf("only COMPLETE payment status is currently supported")
	}

	// Parse customer ID (remove GIG prefix if present)
	customerID, err := utils.ParseCustomerID(req.CustomerID)
	if err != nil {
		return fmt.Errorf("invalid customer_id: %w", err)
	}

	// Parse transaction amount
	amount, err := strconv.ParseFloat(req.TransactionAmount, 64)
	if err != nil {
		return fmt.Errorf("invalid transaction_amount: %w", err)
	}

	// Get or create the customer's account (one account per customer)
	account, err := s.accountRepo.GetByCustomerID(customerID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Parse transaction date
	// Expected format: "2025-11-07 14:54:16"
	transactionDate, err := time.Parse("2006-01-02 15:04:05", req.TransactionDate)
	if err != nil {
		return fmt.Errorf("invalid transaction_date format: %w", err)
	}

	// Create transaction
	createTransactionReq := &models.CreateTransactionRequest{
		CustomerID:      customerID,
		AccountID:       account.ID,
		Reference:       req.TransactionReference,
		Amount:          amount,
		Status:          models.PaymentStatus(strings.ToUpper(req.PaymentStatus)),
		Description:     req.TransactionDate,
		TransactionDate: &transactionDate,
	}

	transaction, err := s.transactionRepo.Create(createTransactionReq)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Credit the account (only if status is COMPLETE)
	if createTransactionReq.Status == models.PaymentStatusComplete {
		go func() {
			err := s.accountRepo.Credit(account.ID, transaction.ID, amount)
			if err != nil {
				log.Printf("failed to credit account: %v", err)
			}
		}()
	}

	return nil
}

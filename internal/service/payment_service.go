package service

import (
	"context"
	json "encoding/json/v2"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/emmrys-jay/gigmile/internal/cache"
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
	cache           cache.Cache
}

func NewPaymentService(
	customerRepo repository.CustomerRepository,
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
	cache cache.Cache,
) PaymentService {
	return &paymentService{
		customerRepo:    customerRepo,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		cache:           cache,
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

	// Get the customer's account from cache or database
	ctx := context.Background()
	cacheKey := fmt.Sprintf("account:customer:%d", customerID)

	var account *models.Account
	cachedData, err := s.cache.Get(ctx, cacheKey)

	if err == nil && cachedData != nil {

		if err := json.Unmarshal(cachedData, &account); err != nil {
			return fmt.Errorf("failed to unmarshal cached account: %w", err)
		}

	} else {

		account, err = s.accountRepo.GetByCustomerID(customerID)
		if err != nil {
			return fmt.Errorf("account not found: %w", err)
		}

		accountData, _ := json.Marshal(struct {
			ID         int64     `json:"id"`
			CustomerID int64     `json:"customer_id"`
			Balance    float64   `json:"balance"`
			CreatedAt  time.Time `json:"created_at"`
			UpdatedAt  time.Time `json:"updated_at"`
		}{
			ID:         account.ID,
			CustomerID: account.CustomerID,
			Balance:    account.Balance,
			CreatedAt:  account.CreatedAt,
			UpdatedAt:  account.UpdatedAt,
		})

		if err := s.cache.Set(ctx, cacheKey, accountData, 5*time.Minute); err != nil {
			log.Printf("failed to set account in cache: %v", err)
		}

	}

	go func() {
		// Parse transaction date
		// Expected format: "2025-11-07 14:54:16"
		transactionDate, err := time.Parse("2006-01-02 15:04:05", req.TransactionDate)
		if err != nil {
			log.Printf("invalid transaction_date format: %v", err)
			return
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
			log.Printf("failed to create transaction: %v", err)
			return
		}

		// Credit the account (only if status is COMPLETE)
		if createTransactionReq.Status == models.PaymentStatusComplete {
			err := s.accountRepo.Credit(account.ID, transaction.ID, amount)
			if err != nil {
				log.Printf("failed to credit account: %v", err)
				return
			}

		}
	}()

	return nil
}

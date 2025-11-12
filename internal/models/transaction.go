package models

import (
	"encoding/json"
	"time"

	"github.com/emmrys-jay/gigmile/internal/utils"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusComplete  PaymentStatus = "COMPLETE"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
)

type Transaction struct {
	ID              int64         `json:"-"`
	CustomerID      int64         `json:"-"`
	AccountID       int64         `json:"-"`
	Reference       string        `json:"reference"`
	Amount          float64       `json:"amount"`
	Status          PaymentStatus `json:"status"`
	Description     *string       `json:"description,omitempty"`
	TransactionDate time.Time     `json:"transaction_date"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// MarshalJSON customizes JSON marshaling to include formatted transaction_id, customer_id, and account_id
func (t *Transaction) MarshalJSON() ([]byte, error) {
	type Alias Transaction

	return json.Marshal(struct {
		ID         string `json:"id"`
		CustomerID string `json:"customer_id"`
		AccountID  string `json:"account_id"`
		Alias
	}{
		ID:         utils.FormatTransactionID(t.ID),
		CustomerID: utils.FormatCustomerID(t.CustomerID),
		AccountID:  utils.FormatAccountID(t.AccountID),
		Alias:      (Alias)(*t),
	})
}

type CreateTransactionRequest struct {
	CustomerID      int64         `json:"customer_id" validate:"required"`
	AccountID       int64         `json:"account_id" validate:"required"`
	Reference       string        `json:"reference"`
	Amount          float64       `json:"amount" validate:"required"`
	Status          PaymentStatus `json:"status" validate:"required"`
	Description     string        `json:"description"`
	TransactionDate *time.Time    `json:"transaction_date,omitempty"` // If nil, will use NOW() in database
}

type UpdateTransactionRequest struct {
	Amount      *float64       `json:"amount,omitempty"`
	Status      *PaymentStatus `json:"status,omitempty"`
	Description *string        `json:"description,omitempty"`
}

type PaymentNotificationRequest struct {
	CustomerID           string `json:"customer_id" validate:"required"`
	PaymentStatus        string `json:"payment_status" validate:"required"`
	TransactionAmount    string `json:"transaction_amount" validate:"required"`
	TransactionDate      string `json:"transaction_date" validate:"required"`
	TransactionReference string `json:"transaction_reference" validate:"required"`
}

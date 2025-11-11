package models

import "time"

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusComplete  PaymentStatus = "COMPLETE"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
)

type Transaction struct {
	ID              int64         `json:"id"`
	CustomerID      int64         `json:"customer_id"`
	AccountID       int64         `json:"account_id"`
	Reference       string        `json:"reference"`
	Amount          float64       `json:"amount"`
	Status          PaymentStatus `json:"status"`
	Description     *string       `json:"description,omitempty"`
	TransactionDate time.Time     `json:"transaction_date"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
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

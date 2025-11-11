package models

import "time"

type TransactionType string

const (
	TransactionTypeDebit  TransactionType = "debit"
	TransactionTypeCredit TransactionType = "credit"
)

type AccountEvent struct {
	ID             int64          `json:"id"`
	TransactionID  int64          `json:"transaction_id"`
	AccountID      int64          `json:"account_id"`
	Type           TransactionType `json:"type"`
	PreviousBalance float64        `json:"previous_balance"`
	NewBalance     float64        `json:"new_balance"`
	CreatedAt      time.Time      `json:"created_at"`
}

type CreateAccountEventRequest struct {
	TransactionID  int64          `json:"transaction_id"`
	AccountID      int64          `json:"account_id"`
	Type           TransactionType `json:"type"`
	PreviousBalance float64        `json:"previous_balance"`
	NewBalance     float64        `json:"new_balance"`
}


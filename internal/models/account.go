package models

import "time"

type Account struct {
	ID        int64     `json:"id"`
	CustomerID int64     `json:"customer_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAccountRequest struct {
	CustomerID  int64   `json:"customer_id"`
	Balance float64 `json:"balance"`
}

type UpdateAccountRequest struct {
	Balance *float64 `json:"balance,omitempty"`
}


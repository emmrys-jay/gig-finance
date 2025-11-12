package models

import (
	"encoding/json"
	"time"

	"github.com/emmrys-jay/gigmile/internal/utils"
)

type Account struct {
	ID         int64     `json:"-"`
	CustomerID int64     `json:"-"`
	Balance    float64   `json:"balance"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MarshalJSON customizes JSON marshaling to include formatted account_id and customer_id
func (a *Account) MarshalJSON() ([]byte, error) {
	type Alias Account

	return json.Marshal(struct {
		ID         string `json:"id"`
		CustomerID string `json:"customer_id"`
		Alias
	}{
		ID:         utils.FormatAccountID(a.ID),
		CustomerID: utils.FormatCustomerID(a.CustomerID),
		Alias:      (Alias)(*a),
	})
}

type CreateAccountRequest struct {
	CustomerID int64   `json:"customer_id"`
	Balance    float64 `json:"balance"`
}

type UpdateAccountRequest struct {
	Balance *float64 `json:"balance,omitempty"`
}

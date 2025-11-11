package models

import (
	"encoding/json"
	"time"

	"github.com/emmrys-jay/gigmile/internal/utils"
)

type Customer struct {
	ID        int64      `json:"-"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// MarshalJSON customizes JSON marshaling to include formatted customer_id
func (c *Customer) MarshalJSON() ([]byte, error) {
	type Alias Customer

	return json.Marshal(struct {
		ID string `json:"id"`
		Alias
	}{
		ID:    utils.FormatCustomerID(c.ID),
		Alias: (Alias)(*c),
	})
}

type CreateCustomerRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type UpdateCustomerRequest struct {
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty"`
}

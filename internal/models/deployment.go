package models

type CreateDeploymentRequest struct {
	CustomerID  string `json:"customer_id" validate:"required"`
	Reference   string `json:"reference" validate:"required"`
	Description string `json:"description"`
}

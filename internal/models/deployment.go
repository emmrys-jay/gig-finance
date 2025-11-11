package models

type CreateDeploymentRequest struct {
	CustomerID  string `json:"customer_id" validate:"required"`
	Reference   string `json:"reference" validate:"required"`
	Description string `json:"description"`
}

type TransactionWithAccountEvent struct {
	Transaction  *Transaction  `json:"transaction"`
	AccountEvent *AccountEvent `json:"account_event"`
}

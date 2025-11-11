package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/emmrys-jay/gigmile/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	Create(transaction *models.CreateTransactionRequest) (*models.Transaction, error)
	GetByID(id int64) (*models.Transaction, error)
	GetByReference(reference string) (*models.Transaction, error)
	GetByCustomerID(customerID int64) ([]*models.Transaction, error)
	GetByCustomerIDWithAccountEvents(customerID int64) ([]*models.TransactionWithAccountEvent, error)
	GetByAccountID(accountID int64) ([]*models.Transaction, error)
	GetByCustomerAndAccountID(customerID, accountID int64) ([]*models.Transaction, error)
	GetAll() ([]*models.Transaction, error)
	Update(id int64, transaction *models.UpdateTransactionRequest) (*models.Transaction, error)
	Delete(id int64) error
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transactionReq *models.CreateTransactionRequest) (*models.Transaction, error) {
	ctx := context.Background()

	var query string
	var args []interface{}

	// Handle description - convert empty string to nil
	var description *string
	if transactionReq.Description != "" {
		description = &transactionReq.Description
	}

	// If transaction_date is provided, include it; otherwise use database default (NOW())
	if transactionReq.TransactionDate != nil {
		query = `
			INSERT INTO transactions (customer_id, account_id, reference, amount, status, description, transaction_date, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
			RETURNING id, customer_id, account_id, reference, amount, status, description, transaction_date, created_at, updated_at
		`
		args = []interface{}{
			transactionReq.CustomerID,
			transactionReq.AccountID,
			transactionReq.Reference,
			transactionReq.Amount,
			transactionReq.Status,
			description,
			*transactionReq.TransactionDate,
		}
	} else {
		query = `
			INSERT INTO transactions (customer_id, account_id, reference, amount, status, description, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
			RETURNING id, customer_id, account_id, reference, amount, status, description, transaction_date, created_at, updated_at
		`
		args = []interface{}{
			transactionReq.CustomerID,
			transactionReq.AccountID,
			transactionReq.Reference,
			transactionReq.Amount,
			transactionReq.Status,
			description,
		}
	}

	transaction := &models.Transaction{}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&transaction.ID,
		&transaction.CustomerID,
		&transaction.AccountID,
		&transaction.Reference,
		&transaction.Amount,
		&transaction.Status,
		&transaction.Description,
		&transaction.TransactionDate,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

func (r *transactionRepository) GetByID(id int64) (*models.Transaction, error) {
	ctx := context.Background()
	query := `
		SELECT id, customer_id, account_id, reference, amount, status, description, transaction_date, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	transaction := &models.Transaction{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.CustomerID,
		&transaction.AccountID,
		&transaction.Reference,
		&transaction.Amount,
		&transaction.Status,
		&transaction.Description,
		&transaction.TransactionDate,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("transaction with id %d not found", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

func (r *transactionRepository) GetByReference(reference string) (*models.Transaction, error) {
	ctx := context.Background()
	query := `
		SELECT id, customer_id, account_id, reference, amount, status, description, transaction_date, created_at, updated_at
		FROM transactions
		WHERE reference = $1
	`

	transaction := &models.Transaction{}
	err := r.db.QueryRow(ctx, query, reference).Scan(
		&transaction.ID,
		&transaction.CustomerID,
		&transaction.AccountID,
		&transaction.Reference,
		&transaction.Amount,
		&transaction.Status,
		&transaction.Description,
		&transaction.TransactionDate,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("transaction with reference %s not found", reference)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

func (r *transactionRepository) GetByCustomerID(customerID int64) ([]*models.Transaction, error) {
	ctx := context.Background()
	query := `
		SELECT id, customer_id, account_id, reference, amount, status, description, transaction_date, created_at, updated_at
		FROM transactions
		WHERE customer_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	transactions := []*models.Transaction{}
	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.CustomerID,
			&transaction.AccountID,
			&transaction.Reference,
			&transaction.Amount,
			&transaction.Status,
			&transaction.Description,
			&transaction.TransactionDate,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

func (r *transactionRepository) GetByCustomerIDWithAccountEvents(customerID int64) ([]*models.TransactionWithAccountEvent, error) {
	ctx := context.Background()
	query := `
		SELECT 
			t.id, t.customer_id, t.account_id, t.reference, t.amount, t.status, t.description, t.transaction_date, t.created_at, t.updated_at,
			ae.id, ae.transaction_id, ae.account_id, ae.type, ae.previous_balance, ae.new_balance, ae.created_at
		FROM transactions t
		LEFT JOIN account_events ae ON t.id = ae.transaction_id
		WHERE t.customer_id = $1
		ORDER BY t.transaction_date DESC
	`

	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	transactions := []*models.TransactionWithAccountEvent{}
	for rows.Next() {
		transaction := &models.Transaction{}
		accountEvent := &models.AccountEvent{}

		var accountEventID *int64
		var accountEventTransactionID *int64
		var accountEventAccountID *int64
		var accountEventType *models.TransactionType
		var accountEventPreviousBalance *float64
		var accountEventNewBalance *float64
		var accountEventCreatedAt *time.Time

		err := rows.Scan(
			&transaction.ID,
			&transaction.CustomerID,
			&transaction.AccountID,
			&transaction.Reference,
			&transaction.Amount,
			&transaction.Status,
			&transaction.Description,
			&transaction.TransactionDate,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&accountEventID,
			&accountEventTransactionID,
			&accountEventAccountID,
			&accountEventType,
			&accountEventPreviousBalance,
			&accountEventNewBalance,
			&accountEventCreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Only include account event if it exists
		if accountEventID != nil {
			accountEvent.ID = *accountEventID
			accountEvent.TransactionID = *accountEventTransactionID
			accountEvent.AccountID = *accountEventAccountID
			accountEvent.Type = *accountEventType
			accountEvent.PreviousBalance = *accountEventPreviousBalance
			accountEvent.NewBalance = *accountEventNewBalance
			accountEvent.CreatedAt = *accountEventCreatedAt
		} else {
			accountEvent = nil
		}

		transactions = append(transactions, &models.TransactionWithAccountEvent{
			Transaction:  transaction,
			AccountEvent: accountEvent,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

func (r *transactionRepository) GetByAccountID(accountID int64) ([]*models.Transaction, error) {
	ctx := context.Background()
	query := `
		SELECT id, customer_id, account_id, reference, amount, status, description, transaction_date, created_at, updated_at
		FROM transactions
		WHERE account_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	transactions := []*models.Transaction{}
	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.CustomerID,
			&transaction.AccountID,
			&transaction.Reference,
			&transaction.Amount,
			&transaction.Status,
			&transaction.Description,
			&transaction.TransactionDate,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

func (r *transactionRepository) GetByCustomerAndAccountID(customerID, accountID int64) ([]*models.Transaction, error) {
	ctx := context.Background()
	query := `
		SELECT id, customer_id, account_id, reference, amount, status, description, transaction_date, created_at, updated_at
		FROM transactions
		WHERE customer_id = $1 AND account_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, customerID, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	transactions := []*models.Transaction{}
	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.CustomerID,
			&transaction.AccountID,
			&transaction.Reference,
			&transaction.Amount,
			&transaction.Status,
			&transaction.Description,
			&transaction.TransactionDate,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

func (r *transactionRepository) GetAll() ([]*models.Transaction, error) {
	ctx := context.Background()
	query := `
		SELECT id, customer_id, account_id, reference, amount, status, description, transaction_date, created_at, updated_at
		FROM transactions
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	transactions := []*models.Transaction{}
	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.CustomerID,
			&transaction.AccountID,
			&transaction.Reference,
			&transaction.Amount,
			&transaction.Status,
			&transaction.Description,
			&transaction.TransactionDate,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

func (r *transactionRepository) Update(id int64, transactionReq *models.UpdateTransactionRequest) (*models.Transaction, error) {
	ctx := context.Background()
	// Build dynamic update query
	query := "UPDATE transactions SET updated_at = NOW()"
	args := []interface{}{}
	argPos := 1

	if transactionReq.Amount != nil {
		query += fmt.Sprintf(", amount = $%d", argPos)
		args = append(args, *transactionReq.Amount)
		argPos++
	}

	if transactionReq.Status != nil {
		query += fmt.Sprintf(", status = $%d", argPos)
		args = append(args, *transactionReq.Status)
		argPos++
	}

	if transactionReq.Description != nil {
		query += fmt.Sprintf(", description = $%d", argPos)
		args = append(args, *transactionReq.Description)
		argPos++
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, customer_id, account_id, reference, amount, status, description, transaction_date, created_at, updated_at", argPos)
	args = append(args, id)

	transaction := &models.Transaction{}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&transaction.ID,
		&transaction.CustomerID,
		&transaction.AccountID,
		&transaction.Reference,
		&transaction.Amount,
		&transaction.Status,
		&transaction.Description,
		&transaction.TransactionDate,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("transaction with id %d not found", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return transaction, nil
}

func (r *transactionRepository) Delete(id int64) error {
	ctx := context.Background()
	query := "DELETE FROM transactions WHERE id = $1"

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("transaction with id %d not found", id)
	}

	return nil
}

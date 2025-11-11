package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/emmrys-jay/gigmile/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository interface {
	Create(account *models.CreateAccountRequest) (*models.Account, error)
	GetByID(id int64) (*models.Account, error)
	GetByCustomerID(customerID int64) (*models.Account, error)
	GetAll() ([]*models.Account, error)
	Update(id int64, account *models.UpdateAccountRequest) (*models.Account, error)
	Delete(id int64) error
	Debit(accountID int64, transactionID int64, amount float64) error
	Credit(accountID int64, transactionID int64, amount float64) error
}

type accountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(accountReq *models.CreateAccountRequest) (*models.Account, error) {
	ctx := context.Background()
	query := `
		INSERT INTO accounts (customer_id, balance, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, customer_id, balance, created_at, updated_at
	`

	account := &models.Account{}
	err := r.db.QueryRow(
		ctx,
		query,
		accountReq.CustomerID,
		accountReq.Balance,
	).Scan(
		&account.ID,
		&account.CustomerID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (r *accountRepository) GetByID(id int64) (*models.Account, error) {
	ctx := context.Background()
	query := `
		SELECT id, customer_id, balance, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`

	account := &models.Account{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.CustomerID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("account with id %d not found", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

func (r *accountRepository) GetByCustomerID(customerID int64) (*models.Account, error) {
	ctx := context.Background()
	query := `
		SELECT id, customer_id, balance, created_at, updated_at
		FROM accounts
		WHERE customer_id = $1
		LIMIT 1
	`

	account := &models.Account{}
	err := r.db.QueryRow(ctx, query, customerID).Scan(
		&account.ID,
		&account.CustomerID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("account not found for customer_id %d", customerID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

func (r *accountRepository) GetAll() ([]*models.Account, error) {
	ctx := context.Background()
	query := `
		SELECT id, customer_id, balance, created_at, updated_at
		FROM accounts
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}
	defer rows.Close()

	accounts := []*models.Account{}
	for rows.Next() {
		account := &models.Account{}
		err := rows.Scan(
			&account.ID,
			&account.CustomerID,
			&account.Balance,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating accounts: %w", err)
	}

	return accounts, nil
}

func (r *accountRepository) Update(id int64, accountReq *models.UpdateAccountRequest) (*models.Account, error) {
	ctx := context.Background()
	// Build dynamic update query
	query := "UPDATE accounts SET updated_at = NOW()"
	args := []interface{}{}
	argPos := 1

	if accountReq.Balance != nil {
		query += fmt.Sprintf(", balance = $%d", argPos)
		args = append(args, *accountReq.Balance)
		argPos++
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, customer_id, balance, created_at, updated_at", argPos)
	args = append(args, id)

	account := &models.Account{}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&account.ID,
		&account.CustomerID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("account with id %d not found", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return account, nil
}

func (r *accountRepository) Delete(id int64) error {
	ctx := context.Background()
	query := "DELETE FROM accounts WHERE id = $1"

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("account with id %d not found", id)
	}

	return nil
}

func (r *accountRepository) Debit(accountID int64, transactionID int64, amount float64) error {
	ctx := context.Background()

	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Lock the account row for update to prevent race conditions
	var previousBalance float64
	var customerID int64
	lockQuery := `
		SELECT balance, customer_id
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`
	err = tx.QueryRow(ctx, lockQuery, accountID).Scan(&previousBalance, &customerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("account with id %d not found", accountID)
	}
	if err != nil {
		return fmt.Errorf("failed to lock account: %w", err)
	}

	// Calculate new balance (debit decreases balance)
	newBalance := previousBalance - amount

	// Update account balance
	updateQuery := `
		UPDATE accounts
		SET balance = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.Exec(ctx, updateQuery, newBalance, accountID)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	// Create account event
	eventQuery := `
		INSERT INTO account_events (transaction_id, account_id, type, previous_balance, new_balance, created_at)
		VALUES ($1, $2, 'debit', $3, $4, NOW())
	`
	_, err = tx.Exec(ctx, eventQuery, transactionID, accountID, previousBalance, newBalance)
	if err != nil {
		return fmt.Errorf("failed to create account event: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *accountRepository) Credit(accountID int64, transactionID int64, amount float64) error {
	ctx := context.Background()

	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Lock the account row for update to prevent race conditions
	var previousBalance float64
	var customerID int64
	lockQuery := `
		SELECT balance, customer_id
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`
	err = tx.QueryRow(ctx, lockQuery, accountID).Scan(&previousBalance, &customerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("account with id %d not found", accountID)
	}
	if err != nil {
		return fmt.Errorf("failed to lock account: %w", err)
	}

	// Calculate new balance (credit increases balance)
	newBalance := previousBalance + amount

	// Update account balance
	updateQuery := `
		UPDATE accounts
		SET balance = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.Exec(ctx, updateQuery, newBalance, accountID)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	// Create account event
	eventQuery := `
		INSERT INTO account_events (transaction_id, account_id, type, previous_balance, new_balance, created_at)
		VALUES ($1, $2, 'credit', $3, $4, NOW())
	`
	_, err = tx.Exec(ctx, eventQuery, transactionID, accountID, previousBalance, newBalance)
	if err != nil {
		return fmt.Errorf("failed to create account event: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

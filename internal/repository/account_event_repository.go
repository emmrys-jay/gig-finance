package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/emmrys-jay/gigmile/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountEventRepository interface {
	Create(event *models.CreateAccountEventRequest) (*models.AccountEvent, error)
	GetByID(id int64) (*models.AccountEvent, error)
	GetByTransactionID(transactionID int64) (*models.AccountEvent, error)
	GetByAccountID(accountID int64) ([]*models.AccountEvent, error)
	GetAll() ([]*models.AccountEvent, error)
}

type accountEventRepository struct {
	db *pgxpool.Pool
}

func NewAccountEventRepository(db *pgxpool.Pool) AccountEventRepository {
	return &accountEventRepository{db: db}
}

func (r *accountEventRepository) Create(eventReq *models.CreateAccountEventRequest) (*models.AccountEvent, error) {
	ctx := context.Background()
	query := `
		INSERT INTO account_events (transaction_id, account_id, type, previous_balance, new_balance, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, transaction_id, account_id, type, previous_balance, new_balance, created_at
	`

	event := &models.AccountEvent{}
	err := r.db.QueryRow(
		ctx,
		query,
		eventReq.TransactionID,
		eventReq.AccountID,
		eventReq.Type,
		eventReq.PreviousBalance,
		eventReq.NewBalance,
	).Scan(
		&event.ID,
		&event.TransactionID,
		&event.AccountID,
		&event.Type,
		&event.PreviousBalance,
		&event.NewBalance,
		&event.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create account event: %w", err)
	}

	return event, nil
}

func (r *accountEventRepository) GetByID(id int64) (*models.AccountEvent, error) {
	ctx := context.Background()
	query := `
		SELECT id, transaction_id, account_id, type, previous_balance, new_balance, created_at
		FROM account_events
		WHERE id = $1
	`

	event := &models.AccountEvent{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.TransactionID,
		&event.AccountID,
		&event.Type,
		&event.PreviousBalance,
		&event.NewBalance,
		&event.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("account event with id %d not found", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get account event: %w", err)
	}

	return event, nil
}

func (r *accountEventRepository) GetByTransactionID(transactionID int64) (*models.AccountEvent, error) {
	ctx := context.Background()
	query := `
		SELECT id, transaction_id, account_id, type, previous_balance, new_balance, created_at
		FROM account_events
		WHERE transaction_id = $1
	`

	event := &models.AccountEvent{}
	err := r.db.QueryRow(ctx, query, transactionID).Scan(
		&event.ID,
		&event.TransactionID,
		&event.AccountID,
		&event.Type,
		&event.PreviousBalance,
		&event.NewBalance,
		&event.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("account event with transaction_id %d not found", transactionID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get account event: %w", err)
	}

	return event, nil
}

func (r *accountEventRepository) GetByAccountID(accountID int64) ([]*models.AccountEvent, error) {
	ctx := context.Background()
	query := `
		SELECT id, transaction_id, account_id, type, previous_balance, new_balance, created_at
		FROM account_events
		WHERE account_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account events: %w", err)
	}
	defer rows.Close()

	events := []*models.AccountEvent{}
	for rows.Next() {
		event := &models.AccountEvent{}
		err := rows.Scan(
			&event.ID,
			&event.TransactionID,
			&event.AccountID,
			&event.Type,
			&event.PreviousBalance,
			&event.NewBalance,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating account events: %w", err)
	}

	return events, nil
}

func (r *accountEventRepository) GetAll() ([]*models.AccountEvent, error) {
	ctx := context.Background()
	query := `
		SELECT id, transaction_id, account_id, type, previous_balance, new_balance, created_at
		FROM account_events
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get account events: %w", err)
	}
	defer rows.Close()

	events := []*models.AccountEvent{}
	for rows.Next() {
		event := &models.AccountEvent{}
		err := rows.Scan(
			&event.ID,
			&event.TransactionID,
			&event.AccountID,
			&event.Type,
			&event.PreviousBalance,
			&event.NewBalance,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating account events: %w", err)
	}

	return events, nil
}

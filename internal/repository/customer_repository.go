package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/emmrys-jay/gigmile/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerRepository interface {
	Create(customer *models.CreateCustomerRequest) (*models.Customer, error)
	GetByID(id int64) (*models.Customer, error)
	GetAll() ([]*models.Customer, error)
	Update(id int64, customer *models.UpdateCustomerRequest) (*models.Customer, error)
	Delete(id int64) error
}

type customerRepository struct {
	db *pgxpool.Pool
}

func NewCustomerRepository(db *pgxpool.Pool) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(customerReq *models.CreateCustomerRequest) (*models.Customer, error) {
	ctx := context.Background()
	query := `
		INSERT INTO customers (email, first_name, last_name, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, email, first_name, last_name, created_at, updated_at, deleted_at
	`

	customer := &models.Customer{}
	err := r.db.QueryRow(
		ctx,
		query,
		customerReq.Email,
		customerReq.FirstName,
		customerReq.LastName,
	).Scan(
		&customer.ID,
		&customer.Email,
		&customer.FirstName,
		&customer.LastName,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&customer.DeletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

func (r *customerRepository) GetByID(id int64) (*models.Customer, error) {
	ctx := context.Background()
	query := `
		SELECT id, email, first_name, last_name, created_at, updated_at, deleted_at
		FROM customers
		WHERE id = $1 AND deleted_at IS NULL
	`

	customer := &models.Customer{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&customer.ID,
		&customer.Email,
		&customer.FirstName,
		&customer.LastName,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&customer.DeletedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("customer with id %d not found", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

func (r *customerRepository) GetAll() ([]*models.Customer, error) {
	ctx := context.Background()
	query := `
		SELECT id, email, first_name, last_name, created_at, updated_at, deleted_at
		FROM customers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}
	defer rows.Close()

	customers := []*models.Customer{}
	for rows.Next() {
		customer := &models.Customer{}
		err := rows.Scan(
			&customer.ID,
			&customer.Email,
			&customer.FirstName,
			&customer.LastName,
			&customer.CreatedAt,
			&customer.UpdatedAt,
			&customer.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating customers: %w", err)
	}

	return customers, nil
}

func (r *customerRepository) Update(id int64, customerReq *models.UpdateCustomerRequest) (*models.Customer, error) {
	ctx := context.Background()
	// Build dynamic update query
	query := "UPDATE customers SET updated_at = NOW()"
	args := []interface{}{}
	argPos := 1

	if customerReq.Email != nil {
		query += fmt.Sprintf(", email = $%d", argPos)
		args = append(args, *customerReq.Email)
		argPos++
	}

	if customerReq.FirstName != nil {
		query += fmt.Sprintf(", first_name = $%d", argPos)
		args = append(args, *customerReq.FirstName)
		argPos++
	}

	if customerReq.LastName != nil {
		query += fmt.Sprintf(", last_name = $%d", argPos)
		args = append(args, *customerReq.LastName)
		argPos++
	}

	query += fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL RETURNING id, email, first_name, last_name, created_at, updated_at, deleted_at", argPos)
	args = append(args, id)

	customer := &models.Customer{}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&customer.ID,
		&customer.Email,
		&customer.FirstName,
		&customer.LastName,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&customer.DeletedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("customer with id %d not found", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return customer, nil
}

func (r *customerRepository) Delete(id int64) error {
	ctx := context.Background()
	query := "UPDATE customers SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL"

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("customer with id %d not found", id)
	}

	return nil
}


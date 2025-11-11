package database

import (
	"context"
	"fmt"

	"github.com/emmrys-jay/gigmile/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func NewDB(cfg *config.Config) (*DB, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(cfg.GetDBURL())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure connection pool settings
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 0 // No limit
	config.MaxConnIdleTime = 0 // No limit

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

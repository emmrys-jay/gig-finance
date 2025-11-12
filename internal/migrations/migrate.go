package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/emmrys-jay/gigmile/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func RunMigrations(cfg *config.Config) error {
	dsn := cfg.GetDBURL()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Use a single connection for migrations to avoid prepared statement conflicts
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Use the newer provider API instead of global SetDialect
	provider, err := goose.NewProvider(goose.DialectPostgres, db, os.DirFS("migrations"))
	if err != nil {
		return fmt.Errorf("failed to create goose provider: %w", err)
	}

	// Run migrations
	ctx := context.Background()
	_, err = provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func RunMigrationsDown(cfg *config.Config) error {
	dsn := cfg.GetDBURL()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Use a single connection for migrations to avoid prepared statement conflicts
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Use the newer provider API instead of global SetDialect
	provider, err := goose.NewProvider(goose.DialectPostgres, db, os.DirFS("migrations"))
	if err != nil {
		return fmt.Errorf("failed to create goose provider: %w", err)
	}

	// Run migrations down
	ctx := context.Background()
	_, err = provider.Down(ctx)
	if err != nil {
		return fmt.Errorf("failed to run migrations down: %w", err)
	}

	return nil
}

func GetMigrationStatus(cfg *config.Config) error {
	dsn := cfg.GetDBURL()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Use a single connection for migrations to avoid prepared statement conflicts
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Use the newer provider API instead of global SetDialect
	provider, err := goose.NewProvider(goose.DialectPostgres, db, os.DirFS("migrations"))
	if err != nil {
		return fmt.Errorf("failed to create goose provider: %w", err)
	}

	// Get migration status
	ctx := context.Background()
	_, err = provider.Status(ctx)
	if err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	return nil
}

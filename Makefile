.PHONY: help migrate-up migrate-down migrate-status migrate-create run build test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

migrate-up: ## Run database migrations up
	@go run cmd/migrate/main.go -command=up

migrate-down: ## Rollback the last migration
	@go run cmd/migrate/main.go -command=down

migrate-status: ## Check migration status
	@go run cmd/migrate/main.go -command=status

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@goose -dir migrations create $(NAME) sql

run: ## Run the application
	@go run cmd/main.go

build: ## Build the application
	@go build -o bin/gigmile cmd/main.go

test: ## Run tests
	@go test ./...



# Gigmile API

A well-structured, decoupled Go API with PostgreSQL database support. This project demonstrates clean architecture principles with clear separation of concerns.

## Architecture

The project follows a layered architecture pattern:

- **Handler Layer**: HTTP request/response handling
- **Service Layer**: Business logic and validation
- **Repository Layer**: Data access and database operations
- **Model Layer**: Data structures and DTOs
- **Database Layer**: Database connection management
- **Config Layer**: Configuration management

## Project Structure

```
gigmile/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   └── config.go               # Configuration management
├── internal/
│   ├── database/
│   │   └── database.go         # Database connection
│   ├── handler/
│   │   └── user_handler.go     # HTTP handlers
│   ├── models/
│   │   └── user.go             # Data models
│   ├── repository/
│   │   └── user_repository.go  # Data access layer
│   ├── router/
│   │   └── router.go           # Route definitions
│   └── service/
│       └── user_service.go     # Business logic
├── migrations/
│   ├── 20240101000000_create_users_table.up.sql   # Migration up
│   └── 20240101000000_create_users_table.down.sql # Migration down
├── internal/
│   └── migrations/
│       └── migrate.go               # Migration helper functions
├── .env.example                # Environment variables template
├── .gitignore
├── go.mod                      # Go dependencies
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Make (optional, for convenience commands)

## Setup

1. **Clone the repository** (if applicable)

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**:
   ```bash
   createdb gigmile
   ```

4. **Run database migrations**:
   ```bash
   go run cmd/migrate/main.go -command=up
   ```
   
   Or using the goose CLI directly:
   ```bash
   goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=gigmile sslmode=disable" up
   ```

5. **Configure environment variables**:
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` with your database credentials:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=gigmile
   DB_SSLMODE=disable
   SERVER_PORT=8080
   ```

6. **Run the application**:
   ```bash
   go run cmd/main.go
   ```

   Or build and run:
   ```bash
   go build -o bin/gigmile cmd/main.go
   ./bin/gigmile
   ```

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Users

- `POST /api/v1/users` - Create a new user
  ```json
  {
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe"
  }
  ```

- `GET /api/v1/users` - Get all users

- `GET /api/v1/users/{id}` - Get user by ID

- `PUT /api/v1/users/{id}` - Update user
  ```json
  {
    "email": "newemail@example.com",
    "first_name": "Jane",
    "last_name": "Smith"
  }
  ```
  Note: All fields are optional in the update request.

- `DELETE /api/v1/users/{id}` - Delete user

## Example Usage

### Create a user
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Get all users
```bash
curl http://localhost:8080/api/v1/users
```

### Get user by ID
```bash
curl http://localhost:8080/api/v1/users/1
```

### Update user
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane"
  }'
```

### Delete user
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Database Migrations

This project uses [goose](https://github.com/pressly/goose) for database migrations.

### Running Migrations

**Using the migration command:**
```bash
# Run all pending migrations
go run cmd/migrate/main.go -command=up

# Rollback the last migration
go run cmd/migrate/main.go -command=down

# Check migration status
go run cmd/migrate/main.go -command=status
```

**Using goose CLI directly:**
```bash
# Install goose CLI (if not already installed)
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations up
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=gigmile sslmode=disable" up

# Rollback last migration
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=gigmile sslmode=disable" down

# Check migration status
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=gigmile sslmode=disable" status
```

### Creating New Migrations

To create a new migration:
```bash
goose -dir migrations create migration_name sql
```

This will create two files:
- `YYYYMMDDHHMMSS_migration_name.up.sql` - Migration up
- `YYYYMMDDHHMMSS_migration_name.down.sql` - Migration down

## Development

### Running Tests
```bash
go test ./...
```

### Code Formatting
```bash
go fmt ./...
```

### Linting
```bash
golangci-lint run
```

## Dependencies

- [gorilla/mux](https://github.com/gorilla/mux) - HTTP router
- [jackc/pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit with connection pooling
- [joho/godotenv](https://github.com/joho/godotenv) - Environment variable management
- [pressly/goose](https://github.com/pressly/goose) - Database migration tool

## License

MIT


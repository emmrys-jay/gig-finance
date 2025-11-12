# Gigmile API

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gigmile
DB_SSLMODE=disable
SERVER_PORT=8080
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

You can copy the example file:
```bash
cp env.example .env
```

## Database Migrations

Run migrations to set up the database schema:

```bash
GOEXPERIMENT=jsonv2 go run cmd/migrate/main.go -command=up
```

Or using Make:
```bash
make migrate-up
```

## Start Server

Run the application:

```bash
GOEXPERIMENT=jsonv2 go run cmd/main.go
```

Or using Make:
```bash
make run
```

The server will start on the port specified in `SERVER_PORT` (default: 8080).

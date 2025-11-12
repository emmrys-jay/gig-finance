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

## Key Endpoints

### 1. Create Customer

Creates a new customer and automatically creates an associated account with a balance of 0.00.

**Endpoint:** `POST /api/v1/customers`

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response (201 Created):**
```json
{
  "status": true,
  "data": {
    "id": "GIG00001",
    "email": "john.doe@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
  },
  "error": "",
  "message": ""
}
```

**Note:** Customer IDs are returned in the format `GIGXXXXX` where `XXXXX` is a zero-padded numeric ID (e.g., `GIG00001` for the first customer).

---

### 2. Notify Payment

Processes a payment notification and credits the customer's account when the payment status is `COMPLETE`. This endpoint also records the transaction.

**Endpoint:** `POST /api/v1/payments/notify`

**Request Body:**
```json
{
  "customer_id": "GIG00001",
  "payment_status": "COMPLETE",
  "transaction_amount": "10000",
  "transaction_date": "2025-11-07 14:54:16",
  "transaction_reference": "VPAY25110713542114478761522000"
}
```

**Response (200 OK):**
```json
{
  "status": true,
  "data": {},
  "error": "",
  "message": ""
}
```

**Notes:**
- The `customer_id` can be provided with or without the `GIG` prefix (e.g., `GIG00001` or `00001`).
- Only `COMPLETE` payment status is currently supported.
- When payment status is `COMPLETE`, the customer's account balance is automatically credited with the transaction amount.
- The transaction is recorded with the provided transaction date and reference.

---

### 3. Record Deployment

Records a deployment and debits 1,000,000 from the customer's wallet. This determines the customer's position that requires settlement.

**Endpoint:** `POST /api/v1/deployments`

**Request Body:**
```json
{
  "customer_id": "GIG00001",
  "reference": "DEPLOY-2025-01-15-001",
  "description": "Production deployment for customer portal"
}
```

**Response (200 OK):**
```json
{
  "status": true,
  "data": null,
  "error": "",
  "message": ""
}
```

**Notes:**
- Each deployment costs exactly 1,000,000.
- The customer's account balance is debited immediately upon recording the deployment.
- A transaction is created with `PENDING` status.
- The `customer_id` can be provided with or without the `GIG` prefix.

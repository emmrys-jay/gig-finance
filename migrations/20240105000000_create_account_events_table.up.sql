-- +goose Up
-- +goose StatementBegin
CREATE TYPE transaction_type AS ENUM ('debit', 'credit');

CREATE TABLE IF NOT EXISTS account_events (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    type transaction_type NOT NULL,
    previous_balance DECIMAL(15, 2) NOT NULL,
    new_balance DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_account_events_transaction_id ON account_events(transaction_id);
CREATE INDEX IF NOT EXISTS idx_account_events_account_id ON account_events(account_id);
CREATE INDEX IF NOT EXISTS idx_account_events_type ON account_events(type);
CREATE INDEX IF NOT EXISTS idx_account_events_created_at ON account_events(created_at);
-- +goose StatementEnd


CREATE TABLE IF NOT EXISTS accounts (
    id          SERIAL PRIMARY KEY,
    user_id     INT         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    currency    VARCHAR(26)  NOT NULL DEFAULT 'TJS',
    balance     NUMERIC(20,2) NOT NULL DEFAULT 0,
    blocked     BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ
);
-- CONSTRAINT chk_accounts_balance_non_negative CHECK (balance >= 0)

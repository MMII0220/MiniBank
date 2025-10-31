CREATE TABLE IF NOT EXISTS transactions (
    id          SERIAL PRIMARY KEY,
    account_id  INT          NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    amount      NUMERIC(20,2) NOT NULL,
    currency    VARCHAR(26)   NOT NULL DEFAULT 'TJS',
    type        VARCHAR(255)  NOT NULL CHECK (type IN ('deposit','withdraw','transfer')),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ
);

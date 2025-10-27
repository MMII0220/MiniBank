CREATE TABLE IF NOT EXISTS cards (
    id               SERIAL PRIMARY KEY,
    account_id       INT         NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    card_number      VARCHAR(255) NOT NULL UNIQUE,      -- keep formatted digits only if you want; else TEXT
    card_holder_name VARCHAR(255) NULL,
    expiry_date      DATE        NOT NULL,      -- last day of month recommended
    cvv              CHAR(255)     NOT NULL,      -- consider NOT storing CVV at all
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ
);

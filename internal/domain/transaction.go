package domain

import "time"

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
	Transfer   TransactionType = "transfer"
)

type Transaction struct {
	ID        int             `json:"id" db:"id"`
	AccountID int             `json:"account_id" db:"account_id"`
	Amount    float64         `json:"amount,omitempty" db:"amount"`
	Currency  string          `json:"currency" db:"currency"`
	Blocked   bool            `json:"blocked" db:"blocked"`
	Type      TransactionType `json:"type" db:"type"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at,omitempty" db:"updated_at"`
}

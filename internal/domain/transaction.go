package domain

import "time"

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdraw"
	Transfer   TransactionType = "transfer"
)

// Чистая доменная модель транзакции
type Transaction struct {
	ID        int
	AccountID int
	Amount    float64
	Currency  string
	Blocked   bool
	Type      TransactionType
	CreatedAt time.Time
	UpdatedAt time.Time
}

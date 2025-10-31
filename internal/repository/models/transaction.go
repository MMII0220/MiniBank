package models

import (
	"time"

	"github.com/MMII0220/MiniBank/internal/domain"
)

// TransactionData для работы с БД в запросах лимитов
type TransactionData struct {
	Amount   float64 `db:"amount"`
	Currency string  `db:"currency"`
}

// TransactionModel для полной работы с транзакциями в БД
type TransactionModel struct {
	ID        int       `db:"id"`
	AccountID int       `db:"account_id"`
	Amount    float64   `db:"amount"`
	Currency  string    `db:"currency"`
	Blocked   bool      `db:"blocked"`
	Type      string    `db:"type"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (tm *TransactionModel) ToDomain() domain.Transaction {
	return domain.Transaction{
		ID:        tm.ID,
		AccountID: tm.AccountID,
		Amount:    tm.Amount,
		Currency:  tm.Currency,
		Blocked:   tm.Blocked,
		Type:      domain.TransactionType(tm.Type),
		CreatedAt: tm.CreatedAt,
		UpdatedAt: tm.UpdatedAt,
	}
}

func TransactionFromDomain(t domain.Transaction) TransactionModel {
	return TransactionModel{
		ID:        t.ID,
		AccountID: t.AccountID,
		Amount:    t.Amount,
		Currency:  t.Currency,
		Blocked:   t.Blocked,
		Type:      string(t.Type),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

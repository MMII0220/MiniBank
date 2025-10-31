package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MMII0220/MiniBank/internal/domain"
)

// AccountModel для работы с аккаунтами в БД
type AccountModel struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Balance   float64   `db:"balance"`
	Currency  string    `db:"currency"`
	Blocked   bool      `db:"blocked"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (am *AccountModel) ToDomain() domain.Account {
	return domain.Account{
		ID:        am.ID,
		UserID:    am.UserID,
		Balance:   fmt.Sprintf("%.2f", am.Balance),
		Currency:  am.Currency,
		Blocked:   am.Blocked,
		CreatedAt: am.CreatedAt,
		UpdatedAt: am.UpdatedAt,
	}
}

func AccountFromDomain(a domain.Account) AccountModel {
	balance, _ := strconv.ParseFloat(a.Balance, 64)

	return AccountModel{
		ID:        a.ID,
		UserID:    a.UserID,
		Balance:   balance,
		Currency:  a.Currency,
		Blocked:   a.Blocked,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

package domain

import "time"

// Чистая доменная модель аккаунта
type Account struct {
	ID        int
	UserID    int
	Currency  string
	Balance   string
	Blocked   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ddd стурктура, прочитать

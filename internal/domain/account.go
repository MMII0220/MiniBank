package domain

import "time"

type Account struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Currency  string    `json:"currency,omitempty" db:"currency"`
	Balance   float64   `json:"balance" db:"balance"`
	Blocked   bool      `json:"blocked,omitempty" db:"blocked"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

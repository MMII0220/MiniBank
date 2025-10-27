package domain

import "time"

type Card struct {
	ID             int       `json:"id" db:"id"`
	AccountID      int       `json:"account_id" db:"account_id"`
	CardNumber     string    `json:"card_number" db:"card_number"`
	CardHolderName string    `json:"card_holder_name" db:"card_holder_name"`
	ExpiryDate     time.Time `json:"expiry_date" db:"expiry_date"`
	CVV            string    `json:"cvv" db:"cvv"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

package domain

import "time"

// Чистая доменная модель карты
type Card struct {
	ID             int
	AccountID      int
	CardNumber     string
	CardHolderName string
	ExpiryDate     time.Time
	CVV            string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

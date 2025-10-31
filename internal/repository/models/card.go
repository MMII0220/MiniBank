package models

import (
	"time"

	"github.com/MMII0220/MiniBank/internal/domain"
)

// CardModel для работы с картами в БД
type CardModel struct {
	ID             int       `db:"id"`
	AccountID      int       `db:"account_id"`
	CardNumber     string    `db:"card_number"`
	CardHolderName string    `db:"cardholder_name"`
	ExpiryDate     time.Time `db:"expiry_date"`
	CVV            string    `db:"cvv"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (cm *CardModel) ToDomain() domain.Card {
	return domain.Card{
		ID:             cm.ID,
		AccountID:      cm.AccountID,
		CardNumber:     cm.CardNumber,
		CardHolderName: cm.CardHolderName,
		ExpiryDate:     cm.ExpiryDate,
		CVV:            cm.CVV,
		CreatedAt:      cm.CreatedAt,
		UpdatedAt:      cm.UpdatedAt,
	}
}

func CardFromDomain(c domain.Card) CardModel {
	return CardModel{
		ID:             c.ID,
		AccountID:      c.AccountID,
		CardNumber:     c.CardNumber,
		CardHolderName: c.CardHolderName,
		ExpiryDate:     c.ExpiryDate,
		CVV:            c.CVV,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

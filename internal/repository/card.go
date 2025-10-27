package repository

import (
	"database/sql"
	"errors"
	"github.com/MMII0220/MiniBank/config"
	"github.com/MMII0220/MiniBank/internal/domain"
)

// Предполагается, что в БД есть unique index на cards.card_number
// CREATE UNIQUE INDEX IF NOT EXISTS idx_cards_number ON cards(card_number);
func CreateCard(card *domain.Card) error {
	query := `
		INSERT INTO cards (account_id, card_number, card_holder_name, expiry_date, cvv, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id
	`

	db := config.GetDBConfig()
	err := db.QueryRow(query,
		card.AccountID,
		card.CardNumber,
		card.CardHolderName,
		card.ExpiryDate,
		card.CVV,
	).Scan(&card.ID)

	// если уникальность нарушена — вернуть понятную ошибку вызывающему
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("failed to create card")
		}
		return err
	}
	return nil
}

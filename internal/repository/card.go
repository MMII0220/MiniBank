package repository

import (
	"database/sql"
	"errors"

	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/repository/models"
)

// Предполагается, что в БД есть unique index на cards.card_number
// CREATE UNIQUE INDEX IF NOT EXISTS idx_cards_number ON cards(card_number);
func (r *Repository) CreateCard(card *domain.Card) error {
	cardModel := models.CardFromDomain(*card)
	query := `
		INSERT INTO cards (account_id, card_number, card_holder_name, expiry_date, cvv, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id
	`

	err := r.db.QueryRow(query,
		cardModel.AccountID,
		cardModel.CardNumber,
		cardModel.CardHolderName,
		cardModel.ExpiryDate,
		cardModel.CVV,
	).Scan(&cardModel.ID)

	// если уникальность нарушена — вернуть понятную ошибку вызывающему
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("failed to create card")
		}
		return r.translateError(err)
	}

	// Обновляем ID в доменном объекте
	card.ID = cardModel.ID
	return nil
}

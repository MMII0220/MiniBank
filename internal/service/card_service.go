// internal/service/card_service.go - ПРОСТАЯ ВЕРСИЯ
package service

import (
	"errors"
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/utils"
	"time"
)

// CreateCardForAccount - простое создание карты
func (s *Service) CreateCardForAccount(accountID int, holderName string) (*domain.Card, error) {
	// Генерируем данные
	cardNumber, _ := utils.GenerateCardNumber()
	cvv, _ := utils.GenerateCVV()
	expiryDate := utils.GenerateExpiry(4) // 4 года

	card := &domain.Card{
		AccountID:      accountID,
		CardNumber:     cardNumber,
		CardHolderName: holderName,
		ExpiryDate:     expiryDate,
		CVV:            cvv,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Сохраняем в БД
	err := s.repo.CreateCard(card)
	if err != nil {
		return nil, errors.New("не удалось создать карту")
	}

	return card, nil
}

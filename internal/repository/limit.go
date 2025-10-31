package repository

import (
	"github.com/MMII0220/MiniBank/config"
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/repository/models"
)

func GetDailyLimitByUserID(userID int) (domain.Limit, error) {
	var limitModel models.LimitModel
	query := `SELECT id, user_id, daily_amount, last_reset FROM limits WHERE user_id = $1`
	err := config.GetDBConfig().Get(&limitModel, query, userID)
	if err != nil {
		return domain.Limit{}, err
	}
	return limitModel.ToDomain(), nil
}

func GetTodayUsageInTJS(userID int) (float64, error) {
	// Курсы валют для конвертации в TJS (должны совпадать с service)
	currencyRates := map[string]float64{
		"TJS": 1.0,
		"USD": 9.21,
		"EUR": 10.72,
	}

	query := `
		SELECT t.amount, t.currency 
		FROM transactions t
		JOIN accounts a ON a.id = t.account_id
		WHERE a.user_id = $1 
		AND t.type IN ('withdraw', 'transfer')
		AND DATE(t.created_at) = CURRENT_DATE
	`

	var transactions []models.TransactionData
	err := config.GetDBConfig().Select(&transactions, query, userID)
	if err != nil {
		return 0, err
	}

	// Суммируем все операции, конвертируя в TJS
	var totalInTJS float64
	for _, tx := range transactions {
		rate, exists := currencyRates[tx.Currency]
		if !exists {
			rate = 1.0 // по умолчанию как TJS
		}

		// Конвертируем в TJS
		amountInTJS := tx.Amount * rate
		totalInTJS += amountInTJS
	}

	return totalInTJS, nil
}

// CreateDailyLimitForUser создает стандартный лимит для нового пользователя
func CreateDailyLimitForUser(userID int, dailyAmount float64) error {
	query := `INSERT INTO limits (user_id, daily_amount, last_reset) VALUES ($1, $2, NOW())`
	_, err := config.GetDBConfig().Exec(query, userID, dailyAmount)
	return err
}

// ResetDailyLimit сбрасывает дневной лимит (обновляет last_reset на сегодня)
func ResetDailyLimit(userID int) error {
	query := `UPDATE limits SET last_reset = NOW() WHERE user_id = $1`
	_, err := config.GetDBConfig().Exec(query, userID)
	return err
}

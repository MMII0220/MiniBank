package service

import (
	"errors"
	"time"
)

// Курсы валют к TJS (базовая валюта)
var currencyRates = map[string]float64{
	"TJS": 1.0,   // базовая валюта
	"USD": 9.21,  // 1 USD = 9.21 TJS
	"EUR": 10.72, // 1 EUR = 10.72 TJS (примерный курс)
}

// Конвертируем любую валюту в TJS для сравнения с лимитом
func (s *Service) ConvertToBaseCurrency(amount float64, currency string) (float64, error) {
	rate, exists := currencyRates[currency]
	if !exists {
		return 0, errors.New("unsupported currency")
	}

	// Конвертируем в TJS
	return amount * rate, nil
}

// Проверяем лимит и рассчитываем комиссию
func (s *Service) CheckLimitAndCalculateFee(userID int, amount float64, currency string) (float64, error) {
	// Конвертируем сумму операции в TJS
	amountInTJS, err := s.ConvertToBaseCurrency(amount, currency)
	if err != nil {
		return 0, s.translateError(err)
	}

	// Получаем лимит пользователя (в TJS)
	limit, err := s.repo.GetDailyLimitByUserID(userID)
	if err != nil {
		// Если лимита нет - создаем стандартный лимит 1000 TJS
		if err.Error() == "sql: no rows in result set" {
			// Возвращаем 0 комиссии если лимитов нет
			return 0, nil
		}
		return 0, s.translateError(err)
	}

	// Проверяем нужно ли сбросить лимит (если прошел день)
	var usedTodayInTJS float64
	if s.IsNewDay(limit.LastReset) {
		err = s.repo.ResetDailyLimit(userID)
		if err != nil {
			return 0, s.translateError(err)
		}
		// После сброса лимита - потрачено 0, комиссии не будет
		usedTodayInTJS = 0
	} else {
		// Получаем уже потраченную сумму сегодня (в TJS)
		usedTodayInTJS, err = s.repo.GetTodayUsageInTJS(userID)
		if err != nil {
			return 0, s.translateError(err)
		}
	}

	// Проверяем превышение лимита
	totalUsageInTJS := usedTodayInTJS + amountInTJS
	if totalUsageInTJS > limit.DailyAmount {
		// Рассчитываем комиссию только с превышающей части
		overlimitAmountInTJS := totalUsageInTJS - limit.DailyAmount

		// Конвертируем превышающую сумму обратно в валюту операции
		rate, exists := currencyRates[currency]
		if !exists {
			return 0, errors.New("unsupported currency")
		}
		overlimitAmountInCurrency := overlimitAmountInTJS / rate

		// Рассчитываем комиссию только с превышающей части
		fee := s.CalculateOverlimitFee(overlimitAmountInCurrency /*, currency*/)
		return fee, nil
	}

	// В пределах лимита - комиссия 0
	return 0, nil
}

// Рассчитываем размер комиссии за превышение лимита
func (s *Service) CalculateOverlimitFee(amount float64 /*, currency string*/) float64 {
	// Комиссия 2% за превышение лимита
	feePercent := 0.02

	// тут в данный момент чуть не равные комиссии настроены типа 100 сомон 1 % это 1 сомн а в долларах 1% это 10 сомон получиться надо будет пзменить потом! Вот таким образом.
	// // Можно настроить разные комиссии для разных валют
	// switch currency {
	// case "TJS":
	// 	feePercent = 0.015 // 1.5% для сомони
	// case "USD", "EUR":
	// 	feePercent = 0.025 // 2.5% для валюты
	// }

	return amount * feePercent
}

// IsNewDay проверяет, прошел ли день с последнего сброса лимита
func (s *Service) IsNewDay(lastReset time.Time) bool {
	now := time.Now()
	return now.Day() != lastReset.Day() || now.Month() != lastReset.Month() || now.Year() != lastReset.Year()
}

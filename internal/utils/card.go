// internal/utils/card.go - ПРОСТАЯ ВЕРСИЯ
package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

// generateRandomDigits - простая генерация N случайных цифр
func generateRandomDigits(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)

	result := ""
	for i := 0; i < n; i++ {
		result += fmt.Sprintf("%d", bytes[i]%10)
	}
	return result
}

// GenerateCardNumber - простой номер карты 16 цифр
func GenerateCardNumber() (string, error) {
	// Префикс 4000 (тестовые карты Visa)
	prefix := "4000"
	// Генерируем 12 случайных цифр
	rest := generateRandomDigits(12)

	return prefix + rest, nil
}

// GenerateCVV - простой CVV 3 цифры
func GenerateCVV() (string, error) {
	return generateRandomDigits(3), nil
}

// GenerateExpiry - простая дата истечения
func GenerateExpiry(years int) time.Time {
	return time.Now().AddDate(years, 0, 0)
}

package service

import (
	"errors"
	"fmt"
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/repository"
)

// GetUserAccounts получает все счета пользователя
func GetUserAccounts(userID int) ([]domain.Account, error) {
	return repository.GetAccountsByUserID(userID)
}

// GetAccountByCard получает счёт по номеру карты с проверками
func GetAccountByCard(cardNumber string) (*domain.Account, error) {
	if cardNumber == "" {
		return nil, errors.New("card number is required")
	}
	
	account, err := repository.GetAccountByCardNumber(cardNumber)
	if err != nil {
		return nil, errors.New("account not found")
	}
	
	return account, nil
}

// GetAccountByPhone получает счёт по телефону с проверками  
func GetAccountByPhone(phoneNumber string) (*domain.Account, error) {
	if phoneNumber == "" {
		return nil, errors.New("phone number is required")
	}
	
	account, err := repository.GetAccountByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, errors.New("account not found")
	}
	
	return account, nil
}

// CreateAccountsForUser создаёт дефолтные счета при регистрации
func CreateAccountsForUser(userID int, currencies []string) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}
	
	for _, currency := range currencies {
		account := domain.Account{
			UserID:   userID,
			Currency: currency,
			Balance:  0,
			Blocked:  false,
		}
		
		err := repository.CreateAccount(&account)
		if err != nil {
			return fmt.Errorf("failed to create %s account: %w", currency, err)
		}
		
		// Создаём карту для каждого счёта
		_, err = CreateCardForAccount(account.ID, "") // имя заполним позже
		if err != nil {
			return fmt.Errorf("failed to create card for %s account: %w", currency, err)
		}
	}
	
	return nil
}

// ValidateAccountForOperation проверяет, можно ли проводить операции со счётом
func ValidateAccountForOperation(account *domain.Account, amount float64, operation string) error {
	if account == nil {
		return errors.New("account not found")
	}
	
	if account.Blocked {
		return errors.New("account is blocked")
	}
	
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	
	// Для операций списания проверяем баланс
	if operation == "withdraw" || operation == "transfer" {
		if account.Balance < amount {
			return errors.New("insufficient funds")
		}
	}
	
	return nil
}

// BlockUserAccount блокирует счёт (для админов)
func BlockUserAccount(accountID int, blocked bool) error {
	return repository.BlockAccount(accountID, blocked)
}
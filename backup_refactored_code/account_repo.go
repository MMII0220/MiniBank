package repository

import (
	"github.com/MMII0220/MiniBank/config"
	"github.com/MMII0220/MiniBank/internal/domain"
)

// CreateAccount создаёт новый счёт
func CreateAccount(account *domain.Account) error {
	query := `
        INSERT INTO accounts (user_id, currency, balance, blocked, created_at)
        VALUES ($1, $2, $3, $4, NOW())
        RETURNING id
    `
	return config.GetDBConfig().QueryRow(query, 
		account.UserID, account.Currency, account.Balance, account.Blocked).Scan(&account.ID)
}

// GetAccountByID получает счёт по ID
func GetAccountByID(accountID int) (*domain.Account, error) {
	var account domain.Account
	query := `SELECT id, user_id, currency, balance, blocked, created_at, updated_at 
			  FROM accounts WHERE id = $1`
	
	err := config.GetDBConfig().Get(&account, query, accountID)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccountsByUserID получает все счета пользователя
func GetAccountsByUserID(userID int) ([]domain.Account, error) {
	var accounts []domain.Account
	query := `SELECT id, user_id, currency, balance, blocked, created_at, updated_at 
			  FROM accounts WHERE user_id = $1 ORDER BY currency`
	
	err := config.GetDBConfig().Select(&accounts, query, userID)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAccountByCardNumber получает счёт по номеру карты
func GetAccountByCardNumber(cardNumber string) (*domain.Account, error) {
	var account domain.Account
	query := `
		SELECT a.id, a.user_id, a.currency, a.balance, a.blocked, a.created_at, a.updated_at
		FROM accounts a
		JOIN cards c ON c.account_id = a.id
		WHERE c.card_number = $1
	`
	
	err := config.GetDBConfig().Get(&account, query, cardNumber)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccountByPhoneNumber получает счёт по телефону (первый TJS счёт)
func GetAccountByPhoneNumber(phoneNumber string) (*domain.Account, error) {
	var account domain.Account
	query := `
		SELECT a.id, a.user_id, a.currency, a.balance, a.blocked, a.created_at, a.updated_at
		FROM accounts a
		JOIN users u ON u.id = a.user_id
		WHERE u.phone = $1 AND a.currency = 'TJS'
		LIMIT 1
	`
	
	err := config.GetDBConfig().Get(&account, query, phoneNumber)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// UpdateAccountBalance обновляет баланс счёта (атомарно)
func UpdateAccountBalance(accountID int, newBalance float64) error {
	query := `UPDATE accounts SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err := config.GetDBConfig().Exec(query, newBalance, accountID)
	return err
}

// BlockAccount блокирует/разблокирует счёт
func BlockAccount(accountID int, blocked bool) error {
	query := `UPDATE accounts SET blocked = $1, updated_at = NOW() WHERE id = $2`
	_, err := config.GetDBConfig().Exec(query, blocked, accountID)
	return err
}
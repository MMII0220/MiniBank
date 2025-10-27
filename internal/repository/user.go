package repository

import (
	"github.com/MMII0220/MiniBank/config"
	"github.com/MMII0220/MiniBank/internal/domain"
)

func CreateUser(user *domain.User) error {
	query := `
		INSERT INTO users (full_name, phone, email, password, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return config.GetDBConfig().Get(&user.ID, query,
		user.FullName, user.Phone, user.Email, user.Password, user.Role)
}

func GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, email, password, role FROM users WHERE email = $1`
	err := config.GetDBConfig().Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateAccount(account *domain.Account) error {
	query := `
        INSERT INTO accounts (user_id, currency, balance, blocked, created_at)
        VALUES ($1, $2, $3, $4, NOW())
        RETURNING id
    `
	return config.GetDBConfig().QueryRow(query, account.UserID, account.Currency, account.Balance, account.Blocked).Scan(&account.ID)
}

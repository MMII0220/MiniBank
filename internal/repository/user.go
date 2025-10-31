package repository

import (
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/repository/models"
)

func (r *Repository) CreateUser(user *domain.User) error {
	userModel := models.UserFromDomain(*user)
	query := `
		INSERT INTO users (full_name, phone, email, password, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.Get(&userModel.ID, query,
		userModel.FullName, userModel.Phone, userModel.Email, userModel.Password, userModel.Role)
	if err == nil {
		user.ID = userModel.ID
	}
	return err
}

func (r *Repository) GetUserByEmail(email string) (*domain.User, error) {
	var userModel models.UserModel
	query := `SELECT id, email, password, role FROM users WHERE email = $1`
	err := r.db.Get(&userModel, query, email)
	if err != nil {
		return nil, err
	}
	domainUser := userModel.ToDomain()
	return &domainUser, nil
}

func (r *Repository) CreateAccount(account *domain.Account) error {
	accountModel := models.AccountFromDomain(*account)
	query := `
        INSERT INTO accounts (user_id, currency, balance, blocked, created_at)
        VALUES ($1, $2, $3, $4, NOW())
        RETURNING id
    `
	err := r.db.QueryRow(query, accountModel.UserID, accountModel.Currency, accountModel.Balance, accountModel.Blocked).Scan(&accountModel.ID)
	if err == nil {
		account.ID = accountModel.ID
	}
	return err
}

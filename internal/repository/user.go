package repository

import (
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/errs"
	"github.com/MMII0220/MiniBank/internal/logger"
	"github.com/MMII0220/MiniBank/internal/redis"
	"github.com/MMII0220/MiniBank/internal/repository/models"
)

func (r *Repository) CreateUser(user *domain.User) error {
	log := logger.GetLogger()
	log.Info().Str("email", user.Email).Msg("Creating new user")

	userModel := models.UserFromDomain(*user)
	query := `
		INSERT INTO users (full_name, phone, email, password, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.Get(&userModel.ID, query,
		userModel.FullName, userModel.Phone, userModel.Email, userModel.Password, userModel.Role)
	if err != nil {
		log.Error().Err(err).Str("email", user.Email).Msg("Failed to create user")
		return r.translateError(err)
	}

	user.ID = userModel.ID
	log.Info().Int("user_id", user.ID).Str("email", user.Email).Msg("User created successfully")
	return nil
}

func (r *Repository) GetUserByEmail(email string) (*domain.User, error) {
	log := logger.GetLogger()
	log.Debug().Str("email", email).Msg("Searching user by email")

	var userModel models.UserModel
	query := `SELECT id, email, password, role FROM users WHERE email = $1`
	err := r.db.Get(&userModel, query, email)
	if err != nil {
		// Если пользователь не найден, возвращаем специфичную ошибку
		if err.Error() == "sql: no rows in result set" {
			log.Debug().Str("email", email).Msg("User not found")
			return nil, errs.ErrUserNotFound
		}
		log.Error().Err(err).Str("email", email).Msg("Failed to get user by email")
		return nil, r.translateError(err)
	}

	domainUser := userModel.ToDomain()
	log.Debug().Int("user_id", domainUser.ID).Str("email", email).Msg("User found successfully")
	return &domainUser, nil
}

func (r *Repository) CreateAccount(account *domain.Account) error {
	log := logger.GetLogger()
	log.Info().
		Int("user_id", account.UserID).
		Str("currency", account.Currency).
		Str("initial_balance", account.Balance).
		Msg("Creating new account")

	accountModel := models.AccountFromDomain(*account)
	query := `
        INSERT INTO accounts (user_id, currency, balance, blocked, created_at)
        VALUES ($1, $2, $3, $4, NOW())
        RETURNING id
    `
	err := r.db.QueryRow(query, accountModel.UserID, accountModel.Currency, accountModel.Balance, accountModel.Blocked).Scan(&accountModel.ID)
	if err != nil {
		log.Error().Err(err).Int("user_id", account.UserID).Msg("Failed to create account")
		return r.translateError(err)
	}

	account.ID = accountModel.ID

	// Удаляем кеш аккаунтов пользователя после создания нового аккаунта
	if cacheErr := redis.DeleteAccountsCache(account.UserID); cacheErr != nil {
		log.Warn().Err(cacheErr).Int("user_id", account.UserID).Msg("Failed to delete account cache after creating account")
	}

	log.Info().
		Int("account_id", account.ID).
		Int("user_id", account.UserID).
		Str("currency", account.Currency).
		Msg("Account created successfully")
	return nil
}

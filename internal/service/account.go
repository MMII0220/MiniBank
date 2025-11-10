package service

import (
	"errors"
	"time"

	"github.com/MMII0220/MiniBank/internal/domain"

	"github.com/MMII0220/MiniBank/internal/logger"
	"github.com/MMII0220/MiniBank/internal/redis"
	redis_client "github.com/redis/go-redis/v9"
)

func (s *Service) BlockUnblockAccount(accountID int, block bool, adminID int, reason string) error {

	if reason == "" {
		return errors.New("reason is required for account blocking/unblocking")
	}

	auditLog := domain.AdminAuditLog{
		AccountID: accountID,
		AdminID:   adminID,
		Action:    map[bool]string{true: "block", false: "unblock"}[block],
		Reason:    reason,
		CreatedAt: time.Now(),
	}

	return s.translateError(s.repo.SetAccountBlock(accountID, block, auditLog))
}

func (s *Service) AuditLogs() ([]domain.AdminAuditLog, error) {
	logs, err := s.repo.GetAuditLogs()
	if err != nil {
		return nil, s.translateError(err)
	}
	return logs, nil
}

func (s *Service) GetAllAccounts(userID int) ([]domain.Account, error) {
	log := logger.GetLogger()
	log.Info().Int("user_id", userID).Msg("Getting all accounts for user")

	// Сначала пробуем кеш
	var cached []domain.Account
	if err := redis.GetAccountsCache(userID, &cached); err == nil {
		log.Info().Int("user_id", userID).Msg("Accounts retrieved from cache")
		return cached, nil
	} else if err != redis_client.Nil {
		// логируем другие ошибки кеша и идем в БД
		log.Debug().Err(err).Int("user_id", userID).Msg("Cache miss/error, fetching from DB")
	}

	// Получаем из базы данных
	accounts, err := s.repo.GetAllAccountsByUserID(userID)
	if err != nil {
		log.Error().Err(err).Int("user_id", userID).Msg("Failed to retrieve accounts from database")
		return nil, s.translateError(err)
	}

	// Пишем в кеш, ошибки игнорируем
	if cacheErr := redis.SetAccountsCache(userID, accounts); cacheErr != nil {
		log.Debug().Err(cacheErr).Int("user_id", userID).Msg("Failed to cache accounts (ignored)")
	}

	log.Info().Int("user_id", userID).Int("accounts_count", len(accounts)).Msg("Accounts retrieved from database")
	return accounts, nil
}

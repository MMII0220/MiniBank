package repository

import (
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/errs"
	"github.com/MMII0220/MiniBank/internal/logger"
	"github.com/MMII0220/MiniBank/internal/redis"
	"github.com/MMII0220/MiniBank/internal/repository/models"
)

func (r *Repository) SetAccountBlock(accountID int, block bool, reqLogs domain.AdminAuditLog) error {
	log := logger.GetLogger()
	log.Info().
		Int("account_id", accountID).
		Bool("block", block).
		Int("admin_id", reqLogs.AdminID).
		Str("reason", reqLogs.Reason).
		Msg("Setting account block status")

	// Начинаем транзакцию
	tx, err := r.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("account_id", accountID).Msg("Failed to begin transaction")
		return err
	}
	defer tx.Rollback() // откатываем если что-то пойдет не так

	// Обновляем аккаунт
	result, err := tx.Exec(`UPDATE accounts SET blocked = $1 WHERE id = $2`, block, accountID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrAccountNotFound
	}

	// Записываем лог
	logModel := models.AdminAuditLogFromDomain(reqLogs)
	_, err = tx.Exec(`INSERT INTO account_audit (account_id, admin_id, action, reason) VALUES ($1, $2, $3, $4)`,
		logModel.AccountID, logModel.AdminID, logModel.Action, logModel.Reason)
	if err != nil {
		return r.translateError(err)
	}

	// Коммитим все изменения
	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Int("account_id", accountID).Msg("Failed to commit transaction")
		return r.translateError(err)
	}

	// Удаляем кеш после успешного изменения статуса блокировки
	if cacheErr := redis.DeleteAccountCacheByAccountID(accountID); cacheErr != nil {
		log.Warn().Err(cacheErr).Int("account_id", accountID).Msg("Failed to delete account cache after block/unblock")
	}

	log.Info().
		Int("account_id", accountID).
		Bool("block", block).
		Msg("Account block status updated successfully")
	return nil
}

func (r *Repository) GetAuditLogs() ([]domain.AdminAuditLog, error) {
	log := logger.GetLogger()
	log.Debug().Msg("Retrieving audit logs")

	var logModels []models.AdminAuditLogModel
	query := `SELECT id, account_id, admin_id, action, reason, created_at FROM account_audit ORDER BY created_at DESC`
	err := r.db.Select(&logModels, query)
	if err != nil {
		return nil, r.translateError(err)
	}

	// Конвертируем в доменные объекты
	logs := make([]domain.AdminAuditLog, len(logModels))
	for i, logModel := range logModels {
		logs[i] = logModel.ToDomain()
	}
	return logs, nil
}

func (r *Repository) GetAllAccountsByUserID(userID int) ([]domain.Account, error) {
	log := logger.GetLogger()
	log.Debug().Int("user_id", userID).Msg("Retrieving all accounts for user")

	var accountModels []models.AccountModel
	query := `SELECT id, user_id, balance, currency, blocked, created_at, updated_at 
			  FROM accounts WHERE user_id = $1 ORDER BY created_at DESC`

	err := r.db.Select(&accountModels, query, userID)
	if err != nil {
		log.Error().Err(err).Int("user_id", userID).Msg("Failed to retrieve accounts")
		return nil, r.translateError(err)
	}

	// Конвертируем в доменные объекты
	accounts := make([]domain.Account, len(accountModels))
	for i, accountModel := range accountModels {
		accounts[i] = accountModel.ToDomain()
	}

	log.Info().Int("user_id", userID).Int("accounts_count", len(accounts)).Msg("Accounts retrieved successfully")
	return accounts, nil
}

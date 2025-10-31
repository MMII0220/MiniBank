package repository

import (
	"errors"

	"github.com/MMII0220/MiniBank/config"
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/repository/models"
)

func SetAccountBlock(accountID int, block bool, reqLogs domain.AdminAuditLog) error {
	// Начинаем транзакцию
	tx, err := config.GetDBConfig().Beginx()
	if err != nil {
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
		return errors.New("account not found")
	}

	// Записываем лог
	logModel := models.AdminAuditLogFromDomain(reqLogs)
	_, err = tx.Exec(`INSERT INTO account_audit (account_id, admin_id, action, reason) VALUES ($1, $2, $3, $4)`,
		logModel.AccountID, logModel.AdminID, logModel.Action, logModel.Reason)
	if err != nil {
		return err
	}

	// Коммитим все изменения
	return tx.Commit()
}

func GetAuditLogs() ([]domain.AdminAuditLog, error) {
	var logModels []models.AdminAuditLogModel
	query := `SELECT id, account_id, admin_id, action, reason, created_at FROM account_audit ORDER BY created_at DESC`
	err := config.GetDBConfig().Select(&logModels, query)
	if err != nil {
		return nil, err
	}

	// Конвертируем в доменные объекты
	logs := make([]domain.AdminAuditLog, len(logModels))
	for i, logModel := range logModels {
		logs[i] = logModel.ToDomain()
	}
	return logs, nil
}

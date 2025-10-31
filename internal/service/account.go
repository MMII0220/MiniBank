package service

import (
	"errors"
	"time"
	
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/repository"
)

func BlockUnblockAccount(accountID int, block bool, adminID int, reason string) error {
	// Валидация бизнес-правил
	if reason == "" {
		return errors.New("reason is required for account blocking/unblocking")
	}
	
	// Создание audit log (это бизнес-логика сервиса)
	auditLog := domain.AdminAuditLog{
		AccountID: accountID,
		AdminID:   adminID,
		Action:    map[bool]string{true: "block", false: "unblock"}[block],
		Reason:    reason,
		CreatedAt: time.Now(),
	}

	return repository.SetAccountBlock(accountID, block, auditLog)
}

func AuditLogs() ([]domain.AdminAuditLog, error) {
	return repository.GetAuditLogs()
}

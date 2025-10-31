package models

import (
	"time"

	"github.com/MMII0220/MiniBank/internal/domain"
)

// AdminAuditLogModel для работы с админскими логами в БД
type AdminAuditLogModel struct {
	ID        int       `db:"id"`
	AccountID int       `db:"account_id"`
	AdminID   int       `db:"admin_id"`
	Action    string    `db:"action"`
	Reason    string    `db:"reason"`
	CreatedAt time.Time `db:"created_at"`
}

func (aal *AdminAuditLogModel) ToDomain() domain.AdminAuditLog {
	return domain.AdminAuditLog{
		ID:        aal.ID,
		AccountID: aal.AccountID,
		AdminID:   aal.AdminID,
		Action:    aal.Action,
		Reason:    aal.Reason,
		CreatedAt: aal.CreatedAt,
	}
}

func AdminAuditLogFromDomain(a domain.AdminAuditLog) AdminAuditLogModel {
	return AdminAuditLogModel{
		ID:        a.ID,
		AccountID: a.AccountID,
		AdminID:   a.AdminID,
		Action:    a.Action,
		Reason:    a.Reason,
		CreatedAt: a.CreatedAt,
	}
}

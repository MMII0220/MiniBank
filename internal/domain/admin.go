package domain

import "time"

// Чистая доменная модель для аудита админских действий
type AdminAuditLog struct {
	ID        int
	AccountID int
	AdminID   int
	Action    string
	Reason    string
	CreatedAt time.Time
}

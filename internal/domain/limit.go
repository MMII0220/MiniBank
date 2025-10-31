package domain

import "time"

// Чистая доменная модель лимитов
type Limit struct {
	ID          int
	UserID      int
	DailyAmount float64
	LastReset   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

package models

import (
	"time"

	"github.com/MMII0220/MiniBank/internal/domain"
)

// LimitModel для работы с лимитами в БД
type LimitModel struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	DailyAmount float64   `db:"daily_amount"`
	LastReset   time.Time `db:"last_reset"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (lm *LimitModel) ToDomain() domain.Limit {
	return domain.Limit{
		ID:          lm.ID,
		UserID:      lm.UserID,
		DailyAmount: lm.DailyAmount,
		LastReset:   lm.LastReset,
		CreatedAt:   lm.CreatedAt,
		UpdatedAt:   lm.UpdatedAt,
	}
}

func LimitFromDomain(l domain.Limit) LimitModel {
	return LimitModel{
		ID:          l.ID,
		UserID:      l.UserID,
		DailyAmount: l.DailyAmount,
		LastReset:   l.LastReset,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}

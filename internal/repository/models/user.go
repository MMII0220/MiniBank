package models

import (
	"time"

	"github.com/MMII0220/MiniBank/internal/domain"
)

// UserModel для работы с пользователями в БД
type UserModel struct {
	ID        int       `db:"id"`
	FullName  string    `db:"full_name"`
	Phone     string    `db:"phone"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (um *UserModel) ToDomain() domain.User {
	return domain.User{
		ID:        um.ID,
		FullName:  um.FullName,
		Phone:     um.Phone,
		Email:     um.Email,
		Password:  um.Password,
		Role:      domain.Role(um.Role),
		CreatedAt: um.CreatedAt.Format(time.RFC3339),
		UpdatedAt: um.UpdatedAt.Format(time.RFC3339),
	}
}

func UserFromDomain(u domain.User) UserModel {
	createdAt, _ := time.Parse(time.RFC3339, u.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, u.UpdatedAt)

	return UserModel{
		ID:        u.ID,
		FullName:  u.FullName,
		Phone:     u.Phone,
		Email:     u.Email,
		Password:  u.Password,
		Role:      string(u.Role),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

package service

import (
	"errors"

	"github.com/MMII0220/MiniBank/internal/domain/contracts"
	"github.com/MMII0220/MiniBank/internal/errs"
)

type Service struct {
	repo contracts.RepositoryI
}

func NewService(repo contracts.RepositoryI) *Service {
	return &Service{
		repo: repo,
	}
}

// translateError - функция service слоя для перевода repository ошибок в business логику
func (s *Service) translateError(err error) error {
	if err == nil {
		return nil
	}

	// Преобразуем repository ошибки в бизнес-логику
	switch {
	case errors.Is(err, errs.ErrUserNotFound):
		return errs.ErrInvalidCredentials // Не раскрываем что пользователь не найден
	case errors.Is(err, errs.ErrAccountNotFound):
		return errs.ErrAccessDenied
	case errors.Is(err, errs.ErrUserAlreadyExists):
		return errs.ErrUserAlreadyRegistered
	case errors.Is(err, errs.ErrCardAlreadyExists):
		return errs.ErrInvalidOperation
	case errors.Is(err, errs.ErrDatabaseError):
		// Не маскируем ошибку БД как 400, пусть контроллер вернёт 500
		return errs.ErrDatabaseError
	default:
		return err // Возвращаем как есть, если не знаем как перевести
	}
}

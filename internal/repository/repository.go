package repository

import (
	"fmt"
	"strings"

	"github.com/MMII0220/MiniBank/internal/errs"
	"github.com/MMII0220/MiniBank/internal/logger"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// PostgreSQL error codes (константы для repository слоя)
const (
	PgUniqueViolation     = "23505" // duplicate key value violates unique constraint
	PgForeignKeyViolation = "23503" // foreign key constraint violation
	PgNotNullViolation    = "23502" // not null constraint violation
	PgCheckViolation      = "23514" // check constraint violation
)

// translateError - функция repository слоя для перевода ошибок БД в доменные ошибки
func (r *Repository) translateError(err error) error {
	if err == nil {
		return nil
	}

	log := logger.GetLogger()

	// Проверяем специальные случаи SQL
	errMsg := err.Error()
	if errMsg == "sql: no rows in result set" {
		log.Debug().Msg("No rows found in database query")
		return errs.ErrUserNotFound // контекстно зависимая ошибка
	}

	// Парсим PostgreSQL ошибки
	if pqErr, ok := err.(*pq.Error); ok {
		log.Error().
			Str("pg_code", string(pqErr.Code)).
			Str("pg_message", pqErr.Message).
			Str("pg_detail", pqErr.Detail).
			Msg("PostgreSQL error occurred")

		switch pqErr.Code {
		case PgUniqueViolation:
			// Определяем какое именно поле нарушает уникальность
			detail := strings.ToLower(pqErr.Detail)
			switch {
			case strings.Contains(detail, "email"):
				log.Warn().Str("field", "email").Msg("Unique constraint violation")
				return fmt.Errorf("user with this email already exists")
			case strings.Contains(detail, "phone"):
				log.Warn().Str("field", "phone").Msg("Unique constraint violation")
				return fmt.Errorf("user with this phone already exists")
			case strings.Contains(detail, "card_number"):
				log.Warn().Str("field", "card_number").Msg("Unique constraint violation")
				return errs.ErrCardAlreadyExists
			default:
				log.Warn().Msg("Unknown unique constraint violation")
				return errs.ErrUserAlreadyExists
			}
		case PgForeignKeyViolation:
			log.Error().Str("constraint_type", "foreign_key").Msg("Foreign key constraint violation")
			return fmt.Errorf("related record does not exist")
		case PgNotNullViolation:
			log.Error().Str("column", pqErr.Column).Msg("Not null constraint violation")
			return fmt.Errorf("required field is missing: %s", pqErr.Column)
		case PgCheckViolation:
			log.Error().Str("constraint_type", "check").Msg("Check constraint violation")
			return fmt.Errorf("invalid data: %s", pqErr.Message)
		default:
			log.Error().Str("unknown_pg_code", string(pqErr.Code)).Msg("Unknown PostgreSQL error")
			return fmt.Errorf("database error: %s", pqErr.Message)
		}
	}

	// Возвращаем общую ошибку БД
	log.Error().Err(err).Msg("Unhandled database error")
	return fmt.Errorf("%w: %v", errs.ErrDatabaseError, err)
}

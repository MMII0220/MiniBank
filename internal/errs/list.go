package errs

import (
	"errors"
	// "fmt"
	// "strings"
	// "github.com/lib/pq"
)

// Все ошибки приложения в одном месте
var (
	// Repository/Database errors
	ErrUserNotFound        = errors.New("user not found")
	ErrAccountNotFound     = errors.New("account not found")
	ErrCardNotFound        = errors.New("card not found")
	ErrLimitNotFound       = errors.New("daily limit not found")
	ErrTransactionNotFound = errors.New("transaction not found")

	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrCardAlreadyExists    = errors.New("card already exists")

	ErrDatabaseError    = errors.New("database operation failed")
	ErrConnectionFailed = errors.New("database connection failed")
	ErrInvalidData      = errors.New("invalid data provided")

	// Authentication errors
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrUserAlreadyRegistered = errors.New("user already registered")
	ErrWeakPassword          = errors.New("password is too weak")
	ErrInvalidEmail          = errors.New("invalid email format")
	ErrInvalidPhone          = errors.New("invalid phone number format")

	// JWT Token errors
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenMalformed     = errors.New("token is malformed")
	ErrInvalidTokenType   = errors.New("invalid token type")
	ErrInvalidTokenClaims = errors.New("invalid token claims")

	// Refresh token errors
	ErrRefreshTokenExpired = errors.New("refresh token has expired")
	ErrRefreshTokenRevoked = errors.New("refresh token has been revoked")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	// Authorization errors
	ErrUnauthorized            = errors.New("unauthorized access")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrSessionExpired          = errors.New("session has expired")
	ErrAccessDenied            = errors.New("access denied")

	// Banking domain errors
	ErrAccountBlocked     = errors.New("account is blocked")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrInvalidAmount      = errors.New("amount must be greater than zero")
	ErrInvalidCurrency    = errors.New("unsupported currency")
	ErrDailyLimitExceeded = errors.New("daily limit exceeded")
	ErrInvalidLimit       = errors.New("invalid limit amount")

	// Card errors
	ErrInvalidCardNumber = errors.New("invalid card number")
	ErrCardExpired       = errors.New("card has expired")
	ErrCardBlocked       = errors.New("card is blocked")

	// Transaction errors
	ErrTransactionFailed    = errors.New("transaction failed")
	ErrDuplicateTransaction = errors.New("duplicate transaction detected")
	ErrSameAccount          = errors.New("cannot transfer to the same account")
	ErrInvalidRecipient     = errors.New("invalid recipient")
	ErrTransferNotAllowed   = errors.New("transfer not allowed")

	// Security errors
	ErrTooManyAttempts    = errors.New("too many failed attempts")
	ErrAccountLocked      = errors.New("account is temporarily locked")
	ErrSuspiciousActivity = errors.New("suspicious activity detected")

	// Operation errors
	ErrOperationNotAllowed = errors.New("operation not allowed")
	ErrInvalidOperation    = errors.New("invalid operation")
)

// PostgreSQL error codes
const (
	PgUniqueViolation     = "23505" // duplicate key value violates unique constraint
	PgForeignKeyViolation = "23503" // foreign key constraint violation
	PgNotNullViolation    = "23502" // not null constraint violation
	PgCheckViolation      = "23514" // check constraint violation
)

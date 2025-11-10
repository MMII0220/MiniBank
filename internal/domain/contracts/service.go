package contracts

import (
	"time"

	"github.com/MMII0220/MiniBank/internal/domain"
)

type ServiceI interface {
	BlockUnblockAccount(accountID int, block bool, adminID int, reason string) error
	AuditLogs() ([]domain.AdminAuditLog, error)

	Register(req domain.ReqRegister, role domain.Role) (domain.User, error)
	Login(req domain.ReqLogin) (domain.TokenResponse, error)
	RefreshToken(req domain.ReqRefreshToken) (domain.TokenResponse, error)
	ParseToken(tokenStr string) (domain.User, error)

	CreateCardForAccount(accountID int, holderName string) (*domain.Card, error)

	ConvertToBaseCurrency(amount float64, currency string) (float64, error)
	CheckLimitAndCalculateFee(userID int, amount float64, currency string) (float64, error)
	CalculateOverlimitFee(amount float64 /*, currency string*/) float64
	IsNewDay(lastReset time.Time) bool

	Deposit(currentUserID int, req domain.ReqTransaction) error
	Withdraw(currentUserID int, req domain.ReqTransaction) error
	Transfer(currentUserID int, req domain.ReqTransfer) error
	HistoryLogs(idUser int) ([]domain.Transaction, error)
	GetAllAccounts(userID int) ([]domain.Account, error)
}

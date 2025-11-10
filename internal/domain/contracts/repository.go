package contracts

import (
	"github.com/MMII0220/MiniBank/internal/domain"
)

type RepositoryI interface {
	SetAccountBlock(accountID int, block bool, reqLogs domain.AdminAuditLog) error
	GetAuditLogs() ([]domain.AdminAuditLog, error)

	CreateCard(card *domain.Card) error

	GetDailyLimitByUserID(userID int) (domain.Limit, error)
	GetTodayUsageInTJS(userID int) (float64, error)
	CreateDailyLimitForUser(userID int, dailyAmount float64) error
	ResetDailyLimit(userID int) error

	DepositToAccount(accountID int, amount float64) error
	WithdrawFromAccount(accountID int, amount float64, currency string) error
	TransferFunds(fromAccountID, toAccountID int, amount float64) error
	GetAccountByCardNumber(account *domain.Account, cardNumber string, currency string) error
	GetAccountByPhoneNumber(account *domain.Account, phoneNumber string, currency string) error
	GetTransactionHistory(idUser int) ([]domain.Transaction, error)

	CreateUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
	CreateAccount(account *domain.Account) error
	GetAllAccountsByUserID(userID int) ([]domain.Account, error)
}

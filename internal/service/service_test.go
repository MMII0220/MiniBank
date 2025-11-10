package service

import (
	"errors"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/errs"
	"golang.org/x/crypto/bcrypt"
)

type mockRepo struct {
	setAccountBlockFn         func(accountID int, block bool, reqLogs domain.AdminAuditLog) error
	getAuditLogsFn            func() ([]domain.AdminAuditLog, error)
	getAllAccountsByUserIDFn  func(userID int) ([]domain.Account, error)
	getTransactionHistoryFn   func(idUser int) ([]domain.Transaction, error)
	createUserFn              func(user *domain.User) error
	createAccountFn           func(account *domain.Account) error
	createCardFn              func(card *domain.Card) error
	createDailyLimitFn        func(userID int, dailyAmount float64) error
	getUserByEmailFn          func(email string) (*domain.User, error)
	getAccountByCardNumberFn  func(account *domain.Account, cardNumber string, currency string) error
	getAccountByPhoneNumberFn func(account *domain.Account, phoneNumber string, currency string) error
	depositToAccountFn        func(accountID int, amount float64) error
	withdrawFromAccountFn     func(accountID int, amount float64, currency string) error
	transferFundsFn           func(fromAccountID, toAccountID int, amount float64) error
	getDailyLimitByUserIDFn   func(userID int) (domain.Limit, error)
	getTodayUsageInTJSFn      func(userID int) (float64, error)
	resetDailyLimitFn         func(userID int) error
}

func (m *mockRepo) SetAccountBlock(accountID int, block bool, reqLogs domain.AdminAuditLog) error {
	if m.setAccountBlockFn != nil {
		return m.setAccountBlockFn(accountID, block, reqLogs)
	}
	return nil
}
func (m *mockRepo) GetAuditLogs() ([]domain.AdminAuditLog, error) {
	if m.getAuditLogsFn != nil {
		return m.getAuditLogsFn()
	}
	return []domain.AdminAuditLog{}, nil
}
func (m *mockRepo) GetDailyLimitByUserID(userID int) (domain.Limit, error) {
	if m.getDailyLimitByUserIDFn != nil {
		return m.getDailyLimitByUserIDFn(userID)
	}
	return domain.Limit{}, nil
}
func (m *mockRepo) GetTodayUsageInTJS(userID int) (float64, error) {
	if m.getTodayUsageInTJSFn != nil {
		return m.getTodayUsageInTJSFn(userID)
	}
	return 0, nil
}
func (m *mockRepo) CreateDailyLimitForUser(userID int, dailyAmount float64) error {
	if m.createDailyLimitFn != nil {
		return m.createDailyLimitFn(userID, dailyAmount)
	}
	return nil
}
func (m *mockRepo) ResetDailyLimit(userID int) error {
	if m.resetDailyLimitFn != nil {
		return m.resetDailyLimitFn(userID)
	}
	return nil
}
func (m *mockRepo) DepositToAccount(accountID int, amount float64) error {
	if m.depositToAccountFn != nil {
		return m.depositToAccountFn(accountID, amount)
	}
	return nil
}
func (m *mockRepo) WithdrawFromAccount(accountID int, amount float64, currency string) error {
	if m.withdrawFromAccountFn != nil {
		return m.withdrawFromAccountFn(accountID, amount, currency)
	}
	return nil
}
func (m *mockRepo) TransferFunds(fromAccountID, toAccountID int, amount float64) error {
	if m.transferFundsFn != nil {
		return m.transferFundsFn(fromAccountID, toAccountID, amount)
	}
	return nil
}
func (m *mockRepo) GetAccountByCardNumber(account *domain.Account, cardNumber string, currency string) error {
	if m.getAccountByCardNumberFn != nil {
		return m.getAccountByCardNumberFn(account, cardNumber, currency)
	}
	return nil
}
func (m *mockRepo) GetAccountByPhoneNumber(account *domain.Account, phoneNumber string, currency string) error {
	if m.getAccountByPhoneNumberFn != nil {
		return m.getAccountByPhoneNumberFn(account, phoneNumber, currency)
	}
	return nil
}
func (m *mockRepo) GetTransactionHistory(idUser int) ([]domain.Transaction, error) {
	if m.getTransactionHistoryFn != nil {
		return m.getTransactionHistoryFn(idUser)
	}
	return nil, nil
}
func (m *mockRepo) CreateUser(user *domain.User) error {
	if m.createUserFn != nil {
		return m.createUserFn(user)
	}
	return nil
}
func (m *mockRepo) GetUserByEmail(email string) (*domain.User, error) {
	if m.getUserByEmailFn != nil {
		return m.getUserByEmailFn(email)
	}
	return nil, nil
}
func (m *mockRepo) CreateAccount(account *domain.Account) error {
	if m.createAccountFn != nil {
		return m.createAccountFn(account)
	}
	return nil
}
func (m *mockRepo) CreateCard(card *domain.Card) error {
	if m.createCardFn != nil {
		return m.createCardFn(card)
	}
	return nil
}
func (m *mockRepo) GetAllAccountsByUserID(userID int) ([]domain.Account, error) {
	if m.getAllAccountsByUserIDFn != nil {
		return m.getAllAccountsByUserIDFn(userID)
	}
	return []domain.Account{}, nil
}

func TestService_BlockUnblockAccount_RequiresReason(t *testing.T) {
	s := NewService(&mockRepo{})
	if err := s.BlockUnblockAccount(1, true, 99, ""); err == nil {
		t.Fatalf("expected error for empty reason")
	}
}

func TestService_BlockUnblockAccount_Success(t *testing.T) {
	called := false
	s := NewService(&mockRepo{setAccountBlockFn: func(accountID int, block bool, reqLogs domain.AdminAuditLog) error {
		called = true
		if accountID != 10 || !block || reqLogs.AdminID != 7 {
			t.Fatalf("wrong args: %+v", reqLogs)
		}
		return nil
	}})
	if err := s.BlockUnblockAccount(10, true, 7, "fraud"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected repo.SetAccountBlock called")
	}
}

func TestService_AuditLogs_DBError(t *testing.T) {
	s := NewService(&mockRepo{getAuditLogsFn: func() ([]domain.AdminAuditLog, error) {
		return nil, errs.ErrDatabaseError
	}})
	_, err := s.AuditLogs()
	if !errors.Is(err, errs.ErrDatabaseError) {
		t.Fatalf("expected ErrDatabaseError, got %v", err)
	}
}

func TestService_GetAllAccounts_DBSuccess(t *testing.T) {
	accounts := []domain.Account{{ID: 1, UserID: 5}, {ID: 2, UserID: 5}}
	s := NewService(&mockRepo{getAllAccountsByUserIDFn: func(userID int) ([]domain.Account, error) {
		if userID != 5 {
			t.Fatalf("expected userID 5, got %d", userID)
		}
		return accounts, nil
	}})
	got, err := s.GetAllAccounts(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, accounts) {
		t.Fatalf("accounts mismatch: %+v", got)
	}
}

func TestService_GetAllAccounts_DBError(t *testing.T) {
	s := NewService(&mockRepo{getAllAccountsByUserIDFn: func(userID int) ([]domain.Account, error) {
		return nil, errs.ErrDatabaseError
	}})
	_, err := s.GetAllAccounts(5)
	if !errors.Is(err, errs.ErrDatabaseError) {
		t.Fatalf("expected ErrDatabaseError, got %v", err)
	}
}

func TestService_translateError_Mapping(t *testing.T) {
	s := NewService(&mockRepo{})
	cases := []struct {
		in  error
		out error
	}{
		{errs.ErrUserNotFound, errs.ErrInvalidCredentials},
		{errs.ErrAccountNotFound, errs.ErrAccessDenied},
		{errs.ErrUserAlreadyExists, errs.ErrUserAlreadyRegistered},
		{errs.ErrCardAlreadyExists, errs.ErrInvalidOperation},
		{errs.ErrDatabaseError, errs.ErrDatabaseError},
	}
	for _, tc := range cases {
		got := s.translateError(tc.in)
		if !errors.Is(got, tc.out) {
			t.Fatalf("expected %v, got %v", tc.out, got)
		}
	}
}

func TestService_LimitHelpers(t *testing.T) {
	s := NewService(&mockRepo{})
	// ConvertToBaseCurrency success
	v, err := s.ConvertToBaseCurrency(2, "USD")
	if err != nil || v <= 0 {
		t.Fatalf("unexpected convert err=%v v=%f", err, v)
	}
	// unsupported currency
	if _, err := s.ConvertToBaseCurrency(1, "ABC"); err == nil {
		t.Fatalf("expected error for unsupported currency")
	}
	// CalculateOverlimitFee
	fee := s.CalculateOverlimitFee(100)
	if math.Abs(fee-2) > 1e-9 {
		t.Fatalf("expected fee 2 got %f", fee)
	}
	// IsNewDay
	if !s.IsNewDay(time.Now().Add(-24 * time.Hour)) {
		t.Fatalf("expected new day")
	}
}

func TestService_HistoryLogs(t *testing.T) {
	s := NewService(&mockRepo{getTransactionHistoryFn: func(idUser int) ([]domain.Transaction, error) {
		if idUser != 5 {
			t.Fatalf("want 5")
		}
		return []domain.Transaction{{ID: 1}}, nil
	}})
	got, err := s.HistoryLogs(5)
	if err != nil || len(got) != 1 {
		t.Fatalf("unexpected: %v %+v", err, got)
	}
}

// extend mockRepo to allow history injection

func TestService_Register_Success(t *testing.T) {
	s := NewService(&mockRepo{createUserFn: func(user *domain.User) error {
		user.ID = 99
		return nil
	}})
	u, err := s.Register(domain.ReqRegister{FullName: "John Doe", Phone: "123", Email: "a@b.c", Password: "password123"}, domain.RoleUser)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if u.ID == 0 {
		t.Fatalf("expected non-zero user id")
	}
}

func TestService_Login_SuccessAndInvalid(t *testing.T) {
	// create hashed password
	pw := "password123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	s := NewService(&mockRepo{getUserByEmailFn: func(email string) (*domain.User, error) {
		return &domain.User{ID: 5, Email: email, Password: string(hashed), Role: domain.RoleUser}, nil
	}})
	// success
	tok, err := s.Login(domain.ReqLogin{Email: "a@b.c", Password: pw})
	if err != nil || tok.AccessToken == "" {
		t.Fatalf("unexpected: %v tok=%+v", err, tok)
	}
	// invalid password
	if _, err := s.Login(domain.ReqLogin{Email: "a@b.c", Password: "wrong"}); err == nil {
		t.Fatalf("expected error on bad password")
	}
}

func TestService_Tokens_ParseAndRefresh(t *testing.T) {
	s := NewService(&mockRepo{})
	// create access token and parse
	at, err := s.createAccessToken(7, domain.RoleAdmin)
	if err != nil {
		t.Fatalf("createAccessToken: %v", err)
	}
	user, err := s.ParseToken(at)
	if err != nil || user.ID != 7 || user.Role != domain.RoleAdmin {
		t.Fatalf("parse fail: %v user=%+v", err, user)
	}
	// create refresh and refresh
	rt, err := s.createRefreshToken(7, domain.RoleAdmin)
	if err != nil {
		t.Fatalf("createRefreshToken: %v", err)
	}
	tr, err := s.RefreshToken(domain.ReqRefreshToken{RefreshToken: rt})
	if err != nil || tr.AccessToken == "" {
		t.Fatalf("refresh failed: %v %+v", err, tr)
	}
}

func TestService_Deposit_And_Withdraw_Success(t *testing.T) {
	s := NewService(&mockRepo{
		getAccountByCardNumberFn: func(acc *domain.Account, card string, currency string) error {
			*acc = domain.Account{ID: 1, UserID: 5, Balance: "100.00", Currency: currency}
			return nil
		},
		depositToAccountFn: func(accountID int, amount float64) error {
			if accountID != 1 || amount != 10 {
				t.Fatalf("bad deposit args")
			}
			return nil
		},
		resetDailyLimitFn: func(userID int) error { return nil },
		withdrawFromAccountFn: func(accountID int, amount float64, currency string) error {
			if accountID != 1 || amount <= 0 {
				t.Fatalf("bad withdraw args")
			}
			return nil
		},
	})
	// deposit
	if err := s.Deposit(5, domain.ReqTransaction{CardNumber: "4000", Amount: 10, Currency: "TJS"}); err != nil {
		t.Fatalf("deposit err: %v", err)
	}
	// withdraw
	if err := s.Withdraw(5, domain.ReqTransaction{CardNumber: "4000", Amount: 10, Currency: "TJS"}); err != nil {
		t.Fatalf("withdraw err: %v", err)
	}
}

func TestService_CreateCardForAccount_Success(t *testing.T) {
	called := false
	s := NewService(&mockRepo{createCardFn: func(card *domain.Card) error {
		called = true
		card.ID = 77
		return nil
	}})
	card, err := s.CreateCardForAccount(10, "John Doe")
	if err != nil || card == nil || card.ID != 77 {
		t.Fatalf("unexpected: %v %+v", err, card)
	}
	if !called {
		t.Fatalf("expected repo.CreateCard called")
	}
}

func TestService_Withdraw_InsufficientIncludingFee(t *testing.T) {
	s := NewService(&mockRepo{
		getAccountByPhoneNumberFn: func(acc *domain.Account, phone string, currency string) error {
			*acc = domain.Account{ID: 1, UserID: 5, Balance: "10.00", Currency: currency}
			return nil
		},
		getDailyLimitByUserIDFn: func(userID int) (domain.Limit, error) {
			return domain.Limit{UserID: userID, DailyAmount: 5.0, LastReset: time.Now()}, nil
		},
		getTodayUsageInTJSFn: func(userID int) (float64, error) { return 5.0, nil },
	})
	// amount 10 > balance 10 after fee > 0, expect error
	err := s.Withdraw(5, domain.ReqTransaction{PhoneNumber: "992", Amount: 10, Currency: "TJS"})
	if err == nil || err.Error() != "insufficient funds including overlimit fee" {
		t.Fatalf("expected overlimit insufficient, got %v", err)
	}
}

func TestService_Transfer_Success(t *testing.T) {
	called := false
	s := NewService(&mockRepo{
		getAccountByCardNumberFn: func(acc *domain.Account, card string, currency string) error {
			if card == "4000" {
				*acc = domain.Account{ID: 1, UserID: 5, Balance: "100.00", Currency: currency}
			}
			if card == "5000" {
				*acc = domain.Account{ID: 2, UserID: 6, Balance: "20.00", Currency: currency}
			}
			return nil
		},
		transferFundsFn: func(fromID, toID int, amount float64) error {
			called = true
			if fromID != 1 || toID != 2 || amount <= 10 {
				t.Fatalf("bad transfer args")
			}
			return nil
		},
	})
	if err := s.Transfer(5, domain.ReqTransfer{FromCardNumber: "4000", ToCardNumber: "5000", Amount: 10, Currency: "TJS"}); err != nil {
		t.Fatalf("transfer err: %v", err)
	}
	if !called {
		t.Fatalf("expected transfer called")
	}
}

package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/errs"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func newMockRepo(t *testing.T) (*Repository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := NewRepository(sqlxDB)
	cleanup := func() { db.Close() }
	return r, mock, cleanup
}

func TestTranslateError_NoRows(t *testing.T) {
	r := &Repository{}
	err := r.translateError(errors.New("sql: no rows in result set"))
	if !errors.Is(err, errs.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestTranslateError_PgUnique_Email(t *testing.T) {
	r := &Repository{}
	pgErr := &pq.Error{Code: pq.ErrorCode(PgUniqueViolation), Detail: "Key (email)=(a) already exists", Message: "duplicate"}
	err := r.translateError(pgErr)
	if err == nil || !regexp.MustCompile(`email`).MatchString(err.Error()) {
		t.Fatalf("expected email unique error, got %v", err)
	}
}

func TestSetAccountBlock_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET blocked = $1 WHERE id = $2")).
		WithArgs(true, 10).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO account_audit (account_id, admin_id, action, reason) VALUES ($1, $2, $3, $4)")).
		WithArgs(10, 99, sqlmock.AnyArg(), "reason").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	audit := domain.AdminAuditLog{AccountID: 10, AdminID: 99, Action: "block", Reason: "reason", CreatedAt: time.Now()}
	if err := r.SetAccountBlock(10, true, audit); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetAuditLogs_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "account_id", "admin_id", "action", "reason", "created_at"}).
		AddRow(1, 10, 99, "block", "r", time.Now())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, admin_id, action, reason, created_at FROM account_audit ORDER BY created_at DESC")).
		WillReturnRows(rows)

	logs, err := r.GetAuditLogs()
	if err != nil || len(logs) != 1 {
		t.Fatalf("expected 1 log, got %v, err=%v", len(logs), err)
	}
}

func TestGetAllAccountsByUserID_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "blocked", "created_at", "updated_at"}).
		AddRow(1, 7, 100.50, "TJS", false, time.Now(), sql.NullTime{})
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, balance, currency, blocked, created_at, updated_at \n\t\t\t\t  FROM accounts WHERE user_id = $1 ORDER BY created_at DESC")).
		WithArgs(7).
		WillReturnRows(rows)

	accounts, err := r.GetAllAccountsByUserID(7)
	if err != nil || len(accounts) != 1 {
		t.Fatalf("expected 1 account, got %v, err=%v", len(accounts), err)
	}
}

func TestDepositToAccount_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET balance = balance + $1 WHERE id = $2")).
		WithArgs(25.0, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT currency FROM accounts WHERE id = $1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"currency"}).AddRow("TJS"))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (account_id, amount, currency, type) VALUES ($1, $2, $3, 'deposit')")).
		WithArgs(1, 25.0, "TJS").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := r.DepositToAccount(1, 25.0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWithdrawFromAccount_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET balance = CAST(balance AS NUMERIC) - $1 WHERE id = $2 AND currency = $3")).
		WithArgs(10.0, 2, "USD").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (account_id, amount, currency, type) VALUES ($1, $2, $3, 'withdraw')")).
		WithArgs(2, 10.0, "USD").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := r.WithdrawFromAccount(2, 10.0, "USD"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWithdrawFromAccount_NoRow(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET balance = CAST(balance AS NUMERIC) - $1 WHERE id = $2 AND currency = $3")).
		WithArgs(10.0, 2, "USD").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := r.WithdrawFromAccount(2, 10.0, "USD")
	if !errors.Is(err, errs.ErrAccountNotFound) {
		t.Fatalf("expected ErrAccountNotFound, got %v", err)
	}
}

func TestTransferFunds_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1")).
		WithArgs(5.0, 3).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET balance = balance + $1 WHERE id = $2")).
		WithArgs(5.0, 4).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO transactions (account_id, amount, type) VALUES ($1, $2, 'transfer')")).
		WithArgs(3, 5.0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := r.TransferFunds(3, 4, 5.0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTransferFunds_Insufficient(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1")).
		WithArgs(5.0, 3).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := r.TransferFunds(3, 4, 5.0)
	if !errors.Is(err, errs.ErrInsufficientFunds) {
		t.Fatalf("expected ErrInsufficientFunds, got %v", err)
	}
}

func TestGetAccountByCardNumber_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "user_id", "currency", "balance", "blocked"}).
		AddRow(10, 20, "TJS", 100.0, false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.id, a.user_id, a.currency, a.balance, a.blocked
		FROM accounts a
		JOIN cards c ON c.account_id = a.id
		WHERE c.card_number = $1 AND a.currency = $2`)).
		WithArgs("4000", "TJS").
		WillReturnRows(rows)

	var acc domain.Account
	if err := r.GetAccountByCardNumber(&acc, "4000", "TJS"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acc.ID != 10 || acc.UserID != 20 {
		t.Fatalf("unexpected account: %+v", acc)
	}
}

func TestGetAccountByPhoneNumber_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "user_id", "currency", "balance", "blocked"}).
		AddRow(11, 21, "USD", 55.0, false)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.id, a.user_id, a.currency, a.balance, a.blocked
        FROM accounts a
        JOIN users u ON u.id = a.user_id
        WHERE u.phone = $1 AND a.currency = $2`)).
		WithArgs("+123", "USD").
		WillReturnRows(rows)

	var acc domain.Account
	if err := r.GetAccountByPhoneNumber(&acc, "+123", "USD"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if acc.ID != 11 || acc.UserID != 21 {
		t.Fatalf("unexpected account: %+v", acc)
	}
}

func TestGetTransactionHistory_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "amount", "currency", "type", "created_at"}).
		AddRow(1, 10.5, "TJS", "deposit", time.Now()).
		AddRow(2, 3.0, "TJS", "withdraw", time.Now())
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.id, t.amount, t.currency, t.type, t.created_at
		FROM transactions t
		JOIN accounts a ON a.id = t.account_id
		WHERE a.user_id = $1
		ORDER BY t.created_at DESC;`)).
		WithArgs(77).
		WillReturnRows(rows)

	trs, err := r.GetTransactionHistory(77)
	if err != nil || len(trs) != 2 {
		t.Fatalf("expected 2 transactions, got %v, err=%v", len(trs), err)
	}
}

func TestCreateUser_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(123)
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (full_name, phone, email, password, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`)).
		WithArgs("John Doe", "+1", "john@example.com", sqlmock.AnyArg(), domain.RoleUser).
		WillReturnRows(rows)

	u := &domain.User{FullName: "John Doe", Phone: "+1", Email: "john@example.com", Password: "hash", Role: domain.RoleUser}
	if err := r.CreateUser(u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID != 123 {
		t.Fatalf("expected ID 123, got %d", u.ID)
	}
}

func TestGetUserByEmail_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "email", "password", "role"}).AddRow(5, "a@b.c", "hash", string(domain.RoleUser))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password, role FROM users WHERE email = $1")).
		WithArgs("a@b.c").
		WillReturnRows(rows)

	u, err := r.GetUserByEmail("a@b.c")
	if err != nil || u == nil || u.ID != 5 {
		t.Fatalf("unexpected result user=%v err=%v", u, err)
	}
}

func TestTranslateError_PgForeignKey(t *testing.T) {
	r := &Repository{}
	pgErr := &pq.Error{Code: pq.ErrorCode(PgForeignKeyViolation), Message: "fk", Detail: "fk"}
	err := r.translateError(pgErr)
	if err == nil || !regexp.MustCompile(`related record`).MatchString(err.Error()) {
		t.Fatalf("expected related record error, got %v", err)
	}
}

func TestTranslateError_PgNotNull(t *testing.T) {
	r := &Repository{}
	pgErr := &pq.Error{Code: pq.ErrorCode(PgNotNullViolation), Column: "email", Message: "not null"}
	err := r.translateError(pgErr)
	if err == nil || !regexp.MustCompile(`required field`).MatchString(err.Error()) {
		t.Fatalf("expected required field missing error, got %v", err)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password, role FROM users WHERE email = $1")).
		WithArgs("none@example.com").
		WillReturnError(sql.ErrNoRows)

	u, err := r.GetUserByEmail("none@example.com")
	if !errors.Is(err, errs.ErrUserNotFound) || u != nil {
		t.Fatalf("expected ErrUserNotFound, got %v, user=%v", err, u)
	}
}

func TestCreateAccount_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(777)
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO accounts (user_id, currency, balance, blocked, created_at)\n        VALUES ($1, $2, $3, $4, NOW())\n        RETURNING id")).
		WithArgs(55, "TJS", 0.0, false).
		WillReturnRows(rows)

	a := &domain.Account{UserID: 55, Currency: "TJS", Balance: "0.00", Blocked: false, CreatedAt: time.Now()}
	if err := r.CreateAccount(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.ID != 777 {
		t.Fatalf("expected account id 777, got %d", a.ID)
	}
}

func TestCreateCard_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(42)
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO cards (account_id, card_number, card_holder_name, expiry_date, cvv, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id`)).
		WithArgs(7, "4000123412341234", "JOHN DOE", sqlmock.AnyArg(), "123").
		WillReturnRows(rows)

	c := &domain.Card{AccountID: 7, CardNumber: "4000123412341234", CardHolderName: "JOHN DOE", ExpiryDate: time.Now(), CVV: "123"}
	if err := r.CreateCard(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ID != 42 {
		t.Fatalf("expected id 42 got %d", c.ID)
	}
}

func TestGetDailyLimitByUserID_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "user_id", "daily_amount", "last_reset"}).AddRow(1, 9, 1000.0, time.Now())
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, daily_amount, last_reset FROM limits WHERE user_id = $1")).
		WithArgs(9).
		WillReturnRows(rows)

	lim, err := r.GetDailyLimitByUserID(9)
	if err != nil || lim.UserID != 9 || lim.DailyAmount != 1000.0 {
		t.Fatalf("unexpected limit: %+v err=%v", lim, err)
	}
}

func TestGetTodayUsageInTJS_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"amount", "currency"}).
		AddRow(10.0, "TJS").
		AddRow(2.0, "USD")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.amount, t.currency 
		FROM transactions t
		JOIN accounts a ON a.id = t.account_id
		WHERE a.user_id = $1 
		AND t.type IN ('withdraw', 'transfer')
		AND DATE(t.created_at) = CURRENT_DATE`)).
		WithArgs(5).
		WillReturnRows(rows)

	total, err := r.GetTodayUsageInTJS(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total <= 10.0 { // should be > 10 due to USD conversion
		t.Fatalf("expected converted total > 10, got %f", total)
	}
}

func TestCreateDailyLimitForUser_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO limits (user_id, daily_amount, last_reset) VALUES ($1, $2, NOW())")).
		WithArgs(8, 500.0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := r.CreateDailyLimitForUser(8, 500.0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResetDailyLimit_Success(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta("UPDATE limits SET last_reset = NOW() WHERE user_id = $1")).
		WithArgs(8).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := r.ResetDailyLimit(8); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetAuditLogs_Error(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, admin_id, action, reason, created_at FROM account_audit ORDER BY created_at DESC")).
		WillReturnError(errors.New("db down"))

	_, err := r.GetAuditLogs()
	if !errors.Is(err, errs.ErrDatabaseError) {
		t.Fatalf("expected ErrDatabaseError, got %v", err)
	}
}

func TestGetAllAccountsByUserID_Error(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, balance, currency, blocked, created_at, updated_at \n\t\t\t\t  FROM accounts WHERE user_id = $1 ORDER BY created_at DESC")).
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	_, err := r.GetAllAccountsByUserID(1)
	if !errors.Is(err, errs.ErrDatabaseError) {
		t.Fatalf("expected ErrDatabaseError, got %v", err)
	}
}

func TestDepositToAccount_NotFound(t *testing.T) {
	r, mock, cleanup := newMockRepo(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET balance = balance + $1 WHERE id = $2")).
		WithArgs(1.0, 999).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := r.DepositToAccount(999, 1.0)
	if !errors.Is(err, errs.ErrAccountNotFound) {
		t.Fatalf("expected ErrAccountNotFound, got %v", err)
	}
}

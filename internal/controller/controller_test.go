package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/errs"
	"github.com/gin-gonic/gin"
)

type mockService struct {
	blockUnblockFn   func(accountID int, block bool, adminID int, reason string) error
	parseTokenFn     func(tokenStr string) (domain.User, error)
	getAllAccountsFn func(userID int) ([]domain.Account, error)
	depositFn        func(currentUserID int, req domain.ReqTransaction) error
	withdrawFn       func(currentUserID int, req domain.ReqTransaction) error
	transferFn       func(currentUserID int, req domain.ReqTransfer) error
	historyFn        func(idUser int) ([]domain.Transaction, error)
	registerFn       func(req domain.ReqRegister, role domain.Role) (domain.User, error)
	loginFn          func(req domain.ReqLogin) (domain.TokenResponse, error)
	refreshFn        func(req domain.ReqRefreshToken) (domain.TokenResponse, error)
	// other methods not used in these tests
}

func (m *mockService) BlockUnblockAccount(accountID int, block bool, adminID int, reason string) error {
	if m.blockUnblockFn != nil {
		return m.blockUnblockFn(accountID, block, adminID, reason)
	}
	return nil
}
func (m *mockService) AuditLogs() ([]domain.AdminAuditLog, error) { return nil, nil }
func (m *mockService) Register(req domain.ReqRegister, role domain.Role) (domain.User, error) {
	if m.registerFn != nil {
		return m.registerFn(req, role)
	}
	return domain.User{}, nil
}
func (m *mockService) Login(req domain.ReqLogin) (domain.TokenResponse, error) {
	if m.loginFn != nil {
		return m.loginFn(req)
	}
	return domain.TokenResponse{}, nil
}
func (m *mockService) RefreshToken(req domain.ReqRefreshToken) (domain.TokenResponse, error) {
	if m.refreshFn != nil {
		return m.refreshFn(req)
	}
	return domain.TokenResponse{}, nil
}
func (m *mockService) ParseToken(tokenStr string) (domain.User, error) {
	if m.parseTokenFn != nil {
		return m.parseTokenFn(tokenStr)
	}
	return domain.User{}, nil
}
func (m *mockService) CreateCardForAccount(accountID int, holderName string) (*domain.Card, error) {
	return nil, nil
}
func (m *mockService) ConvertToBaseCurrency(amount float64, currency string) (float64, error) {
	return 0, nil
}
func (m *mockService) CheckLimitAndCalculateFee(userID int, amount float64, currency string) (float64, error) {
	return 0, nil
}
func (m *mockService) CalculateOverlimitFee(amount float64) float64 { return 0 }
func (m *mockService) IsNewDay(lastReset time.Time) bool            { return false }
func (m *mockService) Deposit(currentUserID int, req domain.ReqTransaction) error {
	if m.depositFn != nil {
		return m.depositFn(currentUserID, req)
	}
	return nil
}
func (m *mockService) Withdraw(currentUserID int, req domain.ReqTransaction) error {
	if m.withdrawFn != nil {
		return m.withdrawFn(currentUserID, req)
	}
	return nil
}
func (m *mockService) Transfer(currentUserID int, req domain.ReqTransfer) error {
	if m.transferFn != nil {
		return m.transferFn(currentUserID, req)
	}
	return nil
}
func (m *mockService) HistoryLogs(idUser int) ([]domain.Transaction, error) {
	if m.historyFn != nil {
		return m.historyFn(idUser)
	}
	return nil, nil
}
func (m *mockService) GetAllAccounts(userID int) ([]domain.Account, error) {
	if m.getAllAccountsFn != nil {
		return m.getAllAccountsFn(userID)
	}
	return []domain.Account{}, nil
}

func TestGetAllAccountsHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockService{getAllAccountsFn: func(userID int) ([]domain.Account, error) {
		return []domain.Account{{ID: 1, UserID: userID}}, nil
	}}
	ctr := NewController(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/accounts", nil)
	c.Set("currentUser", domain.User{ID: 5, Role: domain.RoleUser})

	ctr.getAllAccountsHandler(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "\"total_count\":1") {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestDepositHandler_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/deposit", strings.NewReader(`{"amount":1,"currency":"TJS"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("currentUser", domain.User{ID: 5, Role: domain.RoleUser})

	ctr.depositHandler(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "must be provided") {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestDepositHandler_DBErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{depositFn: func(currentUserID int, req domain.ReqTransaction) error {
		return errs.ErrDatabaseError
	}})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/deposit", strings.NewReader(`{"card_number":"4000","amount":1,"currency":"TJS"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("currentUser", domain.User{ID: 5, Role: domain.RoleUser})

	ctr.depositHandler(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Database error") {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestTransferHandler_Validation_FromTo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// from is set, to is missing -> should fail
	c.Request = httptest.NewRequest(http.MethodPost, "/api/transfer", strings.NewReader(`{"from_card_number":"4000","amount":1,"currency":"TJS"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("currentUser", domain.User{ID: 5, Role: domain.RoleUser})

	ctr.transferHandler(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestWithdrawHandler_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/withdraw", strings.NewReader(`{"amount":2,"currency":"TJS"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("currentUser", domain.User{ID: 5, Role: domain.RoleUser})

	ctr.withdrawHandler(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDepositWithdrawTransfer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{
		depositFn:  func(id int, req domain.ReqTransaction) error { return nil },
		withdrawFn: func(id int, req domain.ReqTransaction) error { return nil },
		transferFn: func(id int, req domain.ReqTransfer) error { return nil },
	})

	// deposit
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/deposit", strings.NewReader(`{"card_number":"4000","amount":1,"currency":"TJS"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("currentUser", domain.User{ID: 5, Role: domain.RoleUser})
	ctr.depositHandler(c)
	if w.Code != http.StatusOK {
		t.Fatalf("deposit expected 200 got %d", w.Code)
	}

	// withdraw
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodPost, "/api/withdraw", strings.NewReader(`{"card_number":"4000","amount":1,"currency":"TJS"}`))
	c2.Request.Header.Set("Content-Type", "application/json")
	c2.Set("currentUser", domain.User{ID: 5, Role: domain.RoleUser})
	ctr.withdrawHandler(c2)
	if w2.Code != http.StatusOK {
		t.Fatalf("withdraw expected 200 got %d", w2.Code)
	}

	// transfer
	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)
	c3.Request = httptest.NewRequest(http.MethodPost, "/api/transfer", strings.NewReader(`{"from_card_number":"4000","to_card_number":"5000","amount":1,"currency":"TJS"}`))
	c3.Request.Header.Set("Content-Type", "application/json")
	c3.Set("currentUser", domain.User{ID: 5, Role: domain.RoleUser})
	ctr.transferHandler(c3)
	if w3.Code != http.StatusOK {
		t.Fatalf("transfer expected 200 got %d", w3.Code)
	}
}

func TestGetAllAccountsHandler_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{getAllAccountsFn: func(userID int) ([]domain.Account, error) { return nil, errs.ErrDatabaseError }})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/accounts", nil)
	c.Set("currentUser", domain.User{ID: 5, Role: domain.RoleUser})
	ctr.getAllAccountsHandler(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", w.Code)
	}
}

func TestHistoryLogs_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{historyFn: func(id int) ([]domain.Transaction, error) { return []domain.Transaction{{ID: 1}}, nil }})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/history", nil)
	c.Set("currentUser", domain.User{ID: 5, Role: domain.RoleUser})
	ctr.historyLogs(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}

func TestAuthMiddleware_HeaderIssues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{})
	r := gin.New()
	r.GET("/p", ctr.AuthMiddleware(domain.RoleUser), func(c *gin.Context) { c.String(200, "ok") })
	// Missing header
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/p", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
	// Bad format
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/p", nil)
	req2.Header.Set("Authorization", "Token abc")
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w2.Code)
	}
}

func TestBindErrorHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{})
	// Register bind error
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{`))
	c.Request.Header.Set("Content-Type", "application/json")
	ctr.registerHandler(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
	// Login bind error
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{`))
	c2.Request.Header.Set("Content-Type", "application/json")
	ctr.loginHandler(c2)
	if w2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w2.Code)
	}
	// Refresh bind error
	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)
	c3.Request = httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(`{`))
	c3.Request.Header.Set("Content-Type", "application/json")
	ctr.refreshTokenHandler(c3)
	if w3.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w3.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{parseTokenFn: func(tokenStr string) (domain.User, error) { return domain.User{}, errs.ErrInvalidToken }})

	r := gin.New()
	r.GET("/protected", ctr.AuthMiddleware(domain.RoleUser), func(c *gin.Context) { c.String(200, "ok") })

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer badtoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_RoleMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{parseTokenFn: func(tokenStr string) (domain.User, error) { return domain.User{ID: 1, Role: domain.RoleUser}, nil }})

	r := gin.New()
	r.GET("/admin", ctr.AuthMiddleware(domain.RoleAdmin), func(c *gin.Context) { c.String(200, "ok") })

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer good")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestAuthMiddleware_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{parseTokenFn: func(tokenStr string) (domain.User, error) { return domain.User{ID: 1, Role: domain.RoleAdmin}, nil }})

	r := gin.New()
	r.GET("/admin", ctr.AuthMiddleware(domain.RoleAdmin), func(c *gin.Context) { c.String(200, "ok") })

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer good")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHealthHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{})

	// /ping
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ping", nil)
	ctr.healthCheck(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// /health/redis - likely returns 503 when redis not initialized
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodGet, "/health/redis", nil)
	ctr.redisHealth(c2)
	if w2.Code != http.StatusServiceUnavailable && w2.Code != http.StatusOK {
		t.Fatalf("expected 503 or 200, got %d", w2.Code)
	}
}

func TestTranslateError_Mapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{})

	cases := []struct {
		err      error
		status   int
		contains string
	}{
		{errs.ErrDatabaseError, http.StatusInternalServerError, "Database error"},
		{errs.ErrInvalidCredentials, http.StatusUnauthorized, "Invalid credentials"},
		{errs.ErrUserAlreadyRegistered, http.StatusConflict, "User already exists"},
		{errs.ErrAccessDenied, http.StatusForbidden, "Access denied"},
		{errs.ErrAccountBlocked, http.StatusForbidden, "Account is blocked"},
		{errs.ErrInsufficientFunds, http.StatusBadRequest, "Insufficient funds"},
		{errs.ErrInvalidAmount, http.StatusBadRequest, "Invalid amount"},
		{errs.ErrDailyLimitExceeded, http.StatusBadRequest, "Daily limit exceeded"},
		{errs.ErrInvalidToken, http.StatusUnauthorized, "Invalid token"},
		{errs.ErrTokenExpired, http.StatusUnauthorized, "Token expired"},
		{errs.ErrRefreshTokenExpired, http.StatusUnauthorized, "Refresh token expired"},
		{errs.ErrOperationNotAllowed, http.StatusBadRequest, "Operation not allowed"},
		{errors.New("unknown"), http.StatusInternalServerError, "Internal server error"},
	}

	for _, tc := range cases {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		ctr.translateError(c, tc.err)
		if w.Code != tc.status {
			t.Fatalf("expected %d got %d for %v", tc.status, w.Code, tc.err)
		}
		if !strings.Contains(w.Body.String(), tc.contains) {
			t.Fatalf("body missing %q: %s", tc.contains, w.Body.String())
		}
	}
}

func TestRegisterHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockService{registerFn: func(req domain.ReqRegister, role domain.Role) (domain.User, error) {
		return domain.User{ID: 123}, nil
	}}
	ctr := NewController(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"full_name":"John Doe","phone":"123","email":"a@b.c","password":"password123","role":"user"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	ctr.registerHandler(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "\"user_id\":123") {
		t.Fatalf("unexpected: %s", w.Body.String())
	}
}

func TestLoginHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockService{loginFn: func(req domain.ReqLogin) (domain.TokenResponse, error) {
		return domain.TokenResponse{AccessToken: "a", RefreshToken: "r"}, nil
	}}
	ctr := NewController(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"email":"a@b.c","password":"pw"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	ctr.loginHandler(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "access_token") {
		t.Fatalf("unexpected: %s", w.Body.String())
	}
}

func TestRefreshTokenHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &mockService{refreshFn: func(req domain.ReqRefreshToken) (domain.TokenResponse, error) {
		return domain.TokenResponse{AccessToken: "a2", RefreshToken: "r2"}, nil
	}}
	ctr := NewController(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"refresh_token":"r1"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	ctr.refreshTokenHandler(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "access_token") {
		t.Fatalf("unexpected: %s", w.Body.String())
	}
}

func TestRegisterHandler_DefaultRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	called := false
	svc := &mockService{registerFn: func(req domain.ReqRegister, role domain.Role) (domain.User, error) {
		called = true
		if role != domain.RoleUser {
			t.Fatalf("expected default role user, got %s", role)
		}
		return domain.User{ID: 1}, nil
	}}
	ctr := NewController(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"full_name":"John","phone":"1","email":"a@b.c","password":"password123","role":"unknown"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	ctr.registerHandler(c)
	if w.Code != http.StatusOK || !called {
		t.Fatalf("expected 200 and called, got %d called=%v", w.Code, called)
	}
}

func TestBlockUnblockAccountHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	called := false
	ctr := NewController(&mockService{blockUnblockFn: func(accountID int, block bool, adminID int, reason string) error {
		called = true
		if accountID != 10 || !block || adminID != 1 || reason != "abuse" {
			t.Fatalf("wrong args: id=%d block=%v admin=%d reason=%s", accountID, block, adminID, reason)
		}
		return nil
	}})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "10"})
	c.Request = httptest.NewRequest(http.MethodPost, "/admin/blockUnblock/10", strings.NewReader(`{"block":true,"reason":"abuse"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("currentUser", domain.User{ID: 1, Role: domain.RoleAdmin})

	ctr.blockUnblockAccountHandler(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	if !called {
		t.Fatalf("expected service to be called")
	}
}

func TestBlockUnblockAccountHandler_NotAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "10"})
	c.Request = httptest.NewRequest(http.MethodPost, "/admin/blockUnblock/10", strings.NewReader(`{"block":true,"reason":"abuse"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("currentUser", domain.User{ID: 2, Role: domain.RoleUser})
	ctr.blockUnblockAccountHandler(c)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", w.Code)
	}
}

func TestGetAuditLogsHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctr := NewController(&mockService{historyFn: func(id int) ([]domain.Transaction, error) { return nil, nil }})
	// But we need AuditLogs method; our mock returns default nil,nil so ok
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/getAuditLogs", nil)
	ctr.getAuditLogsHandler(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}

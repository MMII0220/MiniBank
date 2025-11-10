package domain

// Business DTOs - описывают бизнес-операции и контракты

// Operation for transaction requests
type ReqTransaction struct {
	Amount      float64
	CardNumber  string
	PhoneNumber string
	Currency    string
}

type ReqTransfer struct {
	ToCardNumber    string
	FromCardNumber  string
	ToPhoneNumber   string
	FromPhoneNumber string
	Amount          float64
	Currency        string
}

// For user registration/login requests
type ReqRegister struct {
	FullName string
	Phone    string
	Email    string
	Password string
	Role     string // optional, default "user"
}

type ReqLogin struct {
	Email    string
	Password string
}

// JWT response with access and refresh tokens
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds until access token expires
	TokenType    string `json:"token_type"` // "Bearer"
}

// Refresh token request
type ReqRefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

// For admin account block/unblock requests
type ReqAdminAccountAction struct {
	Block  bool
	Reason string
}

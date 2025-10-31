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

// For admin account block/unblock requests
type ReqAdminAccountAction struct {
	Block  bool
	Reason string
}

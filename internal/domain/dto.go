package domain

// Operation for transaction requests
type ReqTransaction struct {
	Amount      float64 `json:"amount"`
	CardNumber  *string `json:"card_number"`
	PhoneNumber *string `json:"phone_number"`
	Currency    string  `json:"currency"`
}

type ReqTransfer struct {
	ToCardNumber    *string `json:"to_card_number"`
	FromCardNumber  *string `json:"from_card_number"`
	ToPhoneNumber   *string `json:"to_phone_number"`
	FromPhoneNumber *string `json:"from_phone_number"`
	Amount          float64 `json:"amount"`
}

// For user registration/login requests
type ReqRegister struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReqLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

package dto

import "github.com/MMII0220/MiniBank/internal/domain"

// HTTP Request DTOs с JSON тегами

type ReqTransactionHTTP struct {
	Amount      float64 `json:"amount" binding:"required"`
	CardNumber  string  `json:"card_number,omitempty"`
	PhoneNumber string  `json:"phone_number,omitempty"`
	Currency    string  `json:"currency,omitempty"`
}

func (r *ReqTransactionHTTP) ToDomain() domain.ReqTransaction {
	return domain.ReqTransaction{
		Amount:      r.Amount,
		CardNumber:  r.CardNumber,
		PhoneNumber: r.PhoneNumber,
		Currency:    r.Currency,
	}
}

type ReqTransferHTTP struct {
	ToCardNumber    string  `json:"to_card_number,omitempty"`
	FromCardNumber  string  `json:"from_card_number,omitempty"`
	ToPhoneNumber   string  `json:"to_phone_number,omitempty"`
	FromPhoneNumber string  `json:"from_phone_number,omitempty"`
	Amount          float64 `json:"amount" binding:"required"`
	Currency        string  `json:"currency,omitempty"`
}

func (r *ReqTransferHTTP) ToDomain() domain.ReqTransfer {
	return domain.ReqTransfer{
		ToCardNumber:    r.ToCardNumber,
		FromCardNumber:  r.FromCardNumber,
		ToPhoneNumber:   r.ToPhoneNumber,
		FromPhoneNumber: r.FromPhoneNumber,
		Amount:          r.Amount,
		Currency:        r.Currency,
	}
}

type ReqRegisterHTTP struct {
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role,omitempty"`
}

func (r *ReqRegisterHTTP) ToDomain() domain.ReqRegister {
	return domain.ReqRegister{
		FullName: r.FullName,
		Phone:    r.Phone,
		Email:    r.Email,
		Password: r.Password,
		Role:     r.Role,
	}
}

type ReqLoginHTTP struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (r *ReqLoginHTTP) ToDomain() domain.ReqLogin {
	return domain.ReqLogin{
		Email:    r.Email,
		Password: r.Password,
	}
}

type ReqRefreshTokenHTTP struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (r *ReqRefreshTokenHTTP) ToDomain() domain.ReqRefreshToken {
	return domain.ReqRefreshToken{
		RefreshToken: r.RefreshToken,
	}
}

type ReqAdminAccountActionHTTP struct {
	Block  bool   `json:"block"`
	Reason string `json:"reason" binding:"required"`
}

func (r *ReqAdminAccountActionHTTP) ToDomain() domain.ReqAdminAccountAction {
	return domain.ReqAdminAccountAction{
		Block:  r.Block,
		Reason: r.Reason,
	}
}

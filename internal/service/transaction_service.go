package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/MMII0220/MiniBank/internal/domain"
)

func (s *Service) Deposit(currentUserID int, req domain.ReqTransaction) error {
	var account domain.Account
	var err error

	if req.Currency == "" {
		req.Currency = "TJS"
	}

	if req.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if req.CardNumber != "" {
		err = s.repo.GetAccountByCardNumber(&account, req.CardNumber, req.Currency)
	} else if req.PhoneNumber != "" {
		err = s.repo.GetAccountByPhoneNumber(&account, req.PhoneNumber, req.Currency)
	} else {
		return errors.New("either card_number or phone_number must be provided")
	}

	if err != nil {
		return err
	}

	if account.Blocked {
		return errors.New("account is blocked")
	}

	// account.UserID = currentUserID
	if account.UserID != currentUserID {
		return errors.New("access denied")
	}

	return s.repo.DepositToAccount(account.ID, req.Amount)
}

func (s *Service) Withdraw(currentUserID int, req domain.ReqTransaction) error {
	var account domain.Account
	var err error

	if req.Currency == "" {
		req.Currency = "TJS"
	}

	if req.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if req.CardNumber != "" {
		err = s.repo.GetAccountByCardNumber(&account, req.CardNumber, req.Currency)
	} else if req.PhoneNumber != "" {
		err = s.repo.GetAccountByPhoneNumber(&account, req.PhoneNumber, req.Currency)
	} else {
		return errors.New("either card_number or phone_number must be provided")
	}

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New("account not found for the provided phone number and currency")
		}
		return err
	}

	if account.Blocked {
		return errors.New("account is blocked")
	}

	// account.UserID = currentUserID
	if account.UserID != currentUserID {
		return errors.New("access denied")
	}
	balance, err := strconv.ParseFloat(account.Balance, 64)
	if err != nil {
		return errors.New("invalid balance format")
	}

	if req.Amount > balance {
		return errors.New("insufficient funds")
	}

	// Проверяем лимит и получаем комиссию (НЕ перезаписываем req.Amount!)
	fee, err := s.CheckLimitAndCalculateFee(account.UserID, req.Amount, req.Currency)
	if err != nil {
		return err
	}

	// Если есть комиссия - добавляем к основной сумме
	totalAmount := req.Amount + fee
	if totalAmount > balance {
		return errors.New("insufficient funds including overlimit fee")
	}

	// Обновляем req.Amount для списания основной суммы + комиссии
	req.Amount = totalAmount

	return s.repo.WithdrawFromAccount(account.ID, req.Amount, req.Currency)
}

func (s *Service) Transfer(currentUserID int, req domain.ReqTransfer) error {
	var fromAccount, toAccount domain.Account
	var err error

	// Используем TJS по умолчанию для переводов
	if req.Currency == "" {
		req.Currency = "TJS"
	}

	if req.FromCardNumber != "" {
		err = s.repo.GetAccountByCardNumber(&fromAccount, req.FromCardNumber, req.Currency)
	} else if req.FromPhoneNumber != "" {
		err = s.repo.GetAccountByPhoneNumber(&fromAccount, req.FromPhoneNumber, req.Currency)
	}

	if req.ToCardNumber != "" {
		err = s.repo.GetAccountByCardNumber(&toAccount, req.ToCardNumber, req.Currency)
	} else if req.ToPhoneNumber != "" {
		err = s.repo.GetAccountByPhoneNumber(&toAccount, req.ToPhoneNumber, req.Currency)
	}

	if err != nil {
		return err
	}

	if fromAccount.Blocked || toAccount.Blocked {
		return errors.New("one of the accounts is blocked")
	}

	if req.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	balance, err := strconv.ParseFloat(fromAccount.Balance, 64)
	if err != nil {
		return errors.New("invalid balance format")
	}
	if req.Amount > balance {
		return errors.New("insufficient funds")
	}

	// Проверяем лимит и получаем комиссию для переводов (НЕ перезаписываем req.Amount!)
	fee, err := s.CheckLimitAndCalculateFee(fromAccount.UserID, req.Amount, req.Currency)
	if err != nil {
		return err
	}

	fmt.Printf("DEBUG: Transfer fee = %f\n", fee)

	// Если есть комиссия - добавляем к основной сумме
	totalAmount := req.Amount + fee
	if totalAmount > balance {
		return errors.New("insufficient funds including overlimit fee")
	}

	// Обновляем req.Amount для списания основной суммы + комиссии
	req.Amount = totalAmount

	fromAccount.UserID = currentUserID

	// Атомарная операция через репозиторий
	return s.repo.TransferFunds(fromAccount.ID, toAccount.ID, req.Amount)
}

// HistoryLogs возвращает историю операций пользователя
func (s *Service) HistoryLogs(idUser int) ([]domain.Transaction, error) {
	return s.repo.GetTransactionHistory(idUser)
}

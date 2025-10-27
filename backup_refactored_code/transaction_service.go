package service

import (
	"errors"
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/repository"
)

func Deposit(currentUserID int, req domain.ReqTransaction) error {
	var account *domain.Account
	var err error

	if req.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if req.CardNumber != nil {
		account, err = GetAccountByCard(*req.CardNumber)
	} else if req.PhoneNumber != nil {
		account, err = GetAccountByPhone(*req.PhoneNumber)
	} else {
		return errors.New("either card_number or phone_number must be provided")
	}

	if err != nil {
		return err
	}

	// Проверяем, что аккаунт принадлежит текущему пользователю
	if account.UserID != currentUserID {
		return errors.New("access denied: account does not belong to current user")
	}

	err = ValidateAccountForOperation(account, req.Amount, "deposit")
	if err != nil {
		return err
	}

	return repository.DepositToAccount(account.ID, req.Amount)
}

func Withdraw(currentUserID int, req domain.ReqTransaction) error {
	var account *domain.Account
	var err error

	if req.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if req.CardNumber != nil {
		account, err = GetAccountByCard(*req.CardNumber)
	} else if req.PhoneNumber != nil {
		account, err = GetAccountByPhone(*req.PhoneNumber)
	} else {
		return errors.New("either card_number or phone_number must be provided")
	}

	if err != nil {
		return err
	}

	// Проверяем, что аккаунт принадлежит текущему пользователю
	if account.UserID != currentUserID {
		return errors.New("access denied: account does not belong to current user")
	}

	err = ValidateAccountForOperation(account, req.Amount, "withdraw")
	if err != nil {
		return err
	}

	return repository.WithdrawFromAccount(account.ID, req.Amount)
}

func Transfer(currentUserID int, req domain.ReqTransfer) error {
	var fromAccount, toAccount *domain.Account
	var err error

	// Получаем счёт отправителя
	if req.FromCardNumber != nil {
		fromAccount, err = GetAccountByCard(*req.FromCardNumber)
	} else if req.FromPhoneNumber != nil {
		fromAccount, err = GetAccountByPhone(*req.FromPhoneNumber)
	} else {
		return errors.New("from_card_number or from_phone_number is required")
	}

	if err != nil {
		return errors.New("sender account not found")
	}

	// Проверяем, что аккаунт отправителя принадлежит текущему пользователю
	if fromAccount.UserID != currentUserID {
		return errors.New("access denied: sender account does not belong to current user")
	}

	// Получаем счёт получателя
	if req.ToCardNumber != nil {
		toAccount, err = GetAccountByCard(*req.ToCardNumber)
	} else if req.ToPhoneNumber != nil {
		toAccount, err = GetAccountByPhone(*req.ToPhoneNumber)
	} else {
		return errors.New("to_card_number or to_phone_number is required")
	}

	if err != nil {
		return errors.New("recipient account not found")
	}

	// Валидация
	err = ValidateAccountForOperation(fromAccount, req.Amount, "transfer")
	if err != nil {
		return err
	}

	if toAccount.Blocked {
		return errors.New("recipient account is blocked")
	}

	if req.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	// Атомарная операция через репозиторий
	return repository.TransferFunds(fromAccount.ID, toAccount.ID, req.Amount)
}

// HistoryLogs возвращает историю операций пользователя
func HistoryLogs(idUser int) ([]domain.Transaction, error) {
	return repository.GetTransactionHistory(idUser)
}

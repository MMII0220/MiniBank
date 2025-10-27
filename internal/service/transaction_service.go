package service

import (
	"errors"
	// "github.com/MMII0220/MiniBank/internal/controller"
	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/repository"
)

func Deposit(currentUserID int, req domain.ReqTransaction) error {
	var account domain.Account
	var err error

	if req.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if req.CardNumber != nil {
		err = repository.GetAccountByCardNumber(&account, *req.CardNumber)
	} else if req.PhoneNumber != nil {
		err = repository.GetAccountByPhoneNumber(&account, *req.PhoneNumber)
	}

	if err != nil {
		return err
	}

	if account.Blocked {
		return errors.New("account is blocked")
	}

	account.ID = currentUserID

	return repository.DepositToAccount(account.ID, req.Amount)
}

func Withdraw(currentUserID int, req domain.ReqTransaction) error {
	var account domain.Account
	var err error

	if req.Amount > account.Balance {
		return errors.New("balance cannot be negative")
	}

	if req.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if req.CardNumber != nil {
		err = repository.GetAccountByCardNumber(&account, *req.CardNumber)
	}

	if req.PhoneNumber != nil {
		err = repository.GetAccountByPhoneNumber(&account, *req.PhoneNumber)
	}

	if err != nil {
		return err
	}

	if account.Blocked {
		return errors.New("account is blocked")
	}

	account.ID = currentUserID

	return repository.WithdrawFromAccount(account.ID, req.Amount)
}

func Transfer(currentUserID int, req domain.ReqTransfer) error {
	var fromAccount, toAccount domain.Account
	var err error

	if req.FromCardNumber != nil {
		err = repository.GetAccountByCardNumber(&fromAccount, *req.FromCardNumber)
	} else if req.FromPhoneNumber != nil {
		err = repository.GetAccountByPhoneNumber(&fromAccount, *req.FromPhoneNumber)
	}

	if req.ToCardNumber != nil {
		err = repository.GetAccountByCardNumber(&toAccount, *req.ToCardNumber)
	} else if req.ToPhoneNumber != nil {
		err = repository.GetAccountByPhoneNumber(&toAccount, *req.ToPhoneNumber)
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

	if req.Amount > fromAccount.Balance {
		return errors.New("insufficient funds")
	}

	fromAccount.ID = currentUserID

	// Атомарная операция через репозиторий
	return repository.TransferFunds(fromAccount.ID, toAccount.ID, req.Amount)
}

// HistoryLogs возвращает историю операций пользователя
func HistoryLogs(idUser int) ([]domain.Transaction, error) {
	return repository.GetTransactionHistory(idUser)
}

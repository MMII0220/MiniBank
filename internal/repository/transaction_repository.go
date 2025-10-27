package repository

import (
	"fmt"
	"github.com/MMII0220/MiniBank/config"
	"github.com/MMII0220/MiniBank/internal/domain"
)

func DepositToAccount(accountID int, amount float64) error {
	tx, err := config.GetDBConfig().Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, accountID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO transactions (account_id, amount, type) VALUES ($1, $2, 'deposit')`, accountID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func WithdrawFromAccount(accountID int, amount float64) error {
	tx, err := config.GetDBConfig().Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1`, amount, accountID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO transactions (user_id, amount, type) VALUES ($1, $2, 'withdraw')`, accountID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func TransferFunds(fromAccountID, toAccountID int, amount float64) error {
	tx, err := config.GetDBConfig().Beginx()
	if err != nil {
		return err
	}

	// Гарантированный откат, если дальше что-то сломается
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Списываем деньги, но только если хватает баланса
	res, err := tx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1`, amount, fromAccountID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("недостаточно средств на счёте")
	}

	// Зачисляем получателю
	_, err = tx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, toAccountID)
	if err != nil {
		return err
	}

	// Логируем операцию
	_, err = tx.Exec(`INSERT INTO transactions (account_id, amount, type) VALUES ($1, $2, 'transfer')`, fromAccountID, amount)
	if err != nil {
		return err
	}

	// Завершаем успешно
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func GetAccountByCardNumber(account *domain.Account, cardNumber string) error {
	query := `
		SELECT a.id, a.user_id, a.currency, a.balance, a.blocked
		FROM accounts a
		JOIN cards c ON c.account_id = a.id
		WHERE c.card_number = $1
	`

	err := config.GetDBConfig().Get(account, query, cardNumber)
	if err != nil {
		return err
	}

	return nil
}

func GetAccountByPhoneNumber(account *domain.Account, phoneNumber string) error {
	query := `
		SELECT a.id, a.user_id, a.currency, a.balance, a.blocked
		FROM accounts a
		JOIN users u ON u.id = a.user_id
		WHERE u.phone = $1
	`

	err := config.GetDBConfig().Get(account, query, phoneNumber)
	if err != nil {
		return err
	}

	return nil
}

// GetTransactionHistory возвращает список транзакций пользователя
func GetTransactionHistory(idUser int) ([]domain.Transaction, error) {
	query := `
		SELECT t.id, t.amount, t.currency, t.type, t.created_at
		FROM transactions t
		JOIN accounts a ON a.id = t.account_id
		WHERE a.user_id = $1
		ORDER BY t.created_at DESC;
	`

	var transactions []domain.Transaction
	err := config.GetDBConfig().Select(&transactions, query, idUser)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

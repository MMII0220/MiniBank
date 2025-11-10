package repository

import (
	"log"

	"github.com/MMII0220/MiniBank/internal/domain"
	"github.com/MMII0220/MiniBank/internal/errs"
	"github.com/MMII0220/MiniBank/internal/logger"
	"github.com/MMII0220/MiniBank/internal/redis"
	"github.com/MMII0220/MiniBank/internal/repository/models"
)

func (r *Repository) DepositToAccount(accountID int, amount float64) error {
	log := logger.GetLogger()
	log.Info().
		Int("account_id", accountID).
		Float64("amount", amount).
		Msg("Starting deposit transaction")

	tx, err := r.db.Beginx()
	if err != nil {
		log.Error().Err(err).Int("account_id", accountID).Msg("Failed to begin deposit transaction")
		return r.translateError(err)
	}
	defer tx.Rollback()

	// Обновляем баланс счета
	result, err := tx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, accountID)
	if err != nil {
		return r.translateError(err)
	}

	// Проверяем, что строка действительно обновилась
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return r.translateError(err)
	}
	if rowsAffected == 0 {
		return errs.ErrAccountNotFound
	}

	// Получаем валюту счета для транзакции
	var currency string
	err = tx.Get(&currency, `SELECT currency FROM accounts WHERE id = $1`, accountID)
	if err != nil {
		return r.translateError(err)
	}

	// Создаем запись транзакции
	_, err = tx.Exec(`INSERT INTO transactions (account_id, amount, currency, type) VALUES ($1, $2, $3, 'deposit')`, accountID, amount, currency)
	if err != nil {
		return r.translateError(err)
	}

	// Коммитим транзакцию
	err = tx.Commit()
	if err != nil {
		return r.translateError(err)
	}

	// Удаляем кеш после успешного депозита
	if cacheErr := redis.DeleteAccountCacheByAccountID(accountID); cacheErr != nil {
		log.Warn().Err(cacheErr).Int("account_id", accountID).Msg("Failed to delete account cache after deposit")
	}

	log.Info().Int("account_id", accountID).Float64("amount", amount).Msg("Deposit completed successfully")
	return nil
}

func (r *Repository) WithdrawFromAccount(accountID int, amount float64, currency string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return r.translateError(err)
	}
	defer tx.Rollback()

	// Обновляем баланс (убираем проверку balance >= $1 из UPDATE)
	result, err := tx.Exec(`UPDATE accounts SET balance = CAST(balance AS NUMERIC) - $1 WHERE id = $2 AND currency = $3`, amount, accountID, currency)
	if err != nil {
		log.Printf("ERROR: Failed to update balance: %v", err)
		return r.translateError(err)
	}

	// Проверяем что строка обновилась
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return r.translateError(err)
	}
	if rowsAffected == 0 {
		return errs.ErrAccountNotFound
	}

	_, err = tx.Exec(`INSERT INTO transactions (account_id, amount, currency, type) VALUES ($1, $2, $3, 'withdraw')`, accountID, amount, currency)
	if err != nil {
		log.Printf("ERROR: Failed to insert transaction: %v", err)
		return r.translateError(err)
	}

	// Коммитим транзакцию
	err = tx.Commit()
	if err != nil {
		return r.translateError(err)
	}

	// Удаляем кеш после успешного снятия
	if cacheErr := redis.DeleteAccountCacheByAccountID(accountID); cacheErr != nil {
		log.Printf("WARNING: Failed to delete account cache after withdrawal: %v", cacheErr)
	}

	return nil
}

func (r *Repository) TransferFunds(fromAccountID, toAccountID int, amount float64) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return r.translateError(err)
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
		return r.translateError(err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrInsufficientFunds
	}

	// Зачисляем получателю
	_, err = tx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, toAccountID)
	if err != nil {
		return r.translateError(err)
	}

	// Логируем операцию
	_, err = tx.Exec(`INSERT INTO transactions (account_id, amount, type) VALUES ($1, $2, 'transfer')`, fromAccountID, amount)
	if err != nil {
		return r.translateError(err)
	}

	// Завершаем успешно
	err = tx.Commit()
	if err != nil {
		return r.translateError(err)
	}

	// Удаляем кеш для обоих аккаунтов после успешного перевода
	if cacheErr := redis.DeleteAccountCacheByAccountID(fromAccountID); cacheErr != nil {
		log.Printf("WARNING: Failed to delete cache for sender account: %v", cacheErr)
	}
	if cacheErr := redis.DeleteAccountCacheByAccountID(toAccountID); cacheErr != nil {
		log.Printf("WARNING: Failed to delete cache for recipient account: %v", cacheErr)
	}

	return nil
}

func (r *Repository) GetAccountByCardNumber(account *domain.Account, cardNumber string, currency string) error {
	var accountModel models.AccountModel
	query := `
		SELECT a.id, a.user_id, a.currency, a.balance, a.blocked
		FROM accounts a
		JOIN cards c ON c.account_id = a.id
		WHERE c.card_number = $1 AND a.currency = $2
	`

	err := r.db.Get(&accountModel, query, cardNumber, currency)
	if err != nil {
		return r.translateError(err)
	}

	*account = accountModel.ToDomain()
	return nil
}

func (r *Repository) GetAccountByPhoneNumber(account *domain.Account, phoneNumber string, currency string) error {
	var accountModel models.AccountModel
	query := `
        SELECT a.id, a.user_id, a.currency, a.balance, a.blocked
        FROM accounts a
        JOIN users u ON u.id = a.user_id
        WHERE u.phone = $1 AND a.currency = $2
    `

	err := r.db.Get(&accountModel, query, phoneNumber, currency)
	if err != nil {
		log.Printf("ERROR: Аккаунт не найден: %v", err)
		return r.translateError(err)
	}

	*account = accountModel.ToDomain()
	return nil
}

// GetTransactionHistory возвращает список транзакций пользователя
func (r *Repository) GetTransactionHistory(idUser int) ([]domain.Transaction, error) {
	query := `
		SELECT t.id, t.amount, t.currency, t.type, t.created_at
		FROM transactions t
		JOIN accounts a ON a.id = t.account_id
		WHERE a.user_id = $1
		ORDER BY t.created_at DESC;
	`

	var transactionModels []models.TransactionModel
	err := r.db.Select(&transactionModels, query, idUser)
	if err != nil {
		return nil, r.translateError(err)
	}

	// Конвертируем в доменные объекты
	transactions := make([]domain.Transaction, len(transactionModels))
	for i, tm := range transactionModels {
		transactions[i] = tm.ToDomain()
	}

	return transactions, nil
}

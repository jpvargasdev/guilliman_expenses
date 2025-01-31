package models

import (
	"database/sql"
	"fmt"
	"guilliman/internal/utils"
	"guilliman/internal/utils/timeutils"
	"log"
	"strings"
	"time"
)

const (
	TransactionTypeIncome   = "Income"
	TransactionTypeExpense  = "Expense"
	TransactionTypeSavings  = "Savings"
	TransactionTypeTransfer = "Transfer"
)

const (
	MainCategoryNeeds    = "Needs"
	MainCategoryWants    = "Wants"
	MainCategorySavings  = "Savings"
	MainCategoryTransfer = "Transfer"
)

type Transaction struct {
	ID                   int     `json:"id"`
	Description          string  `json:"description"`
	Amount               float64 `json:"amount"`
	Currency             string  `json:"currency"`
	AmountInBaseCurrency float64 `json:"amount_in_base_currency"`
	ExchangeRate         float64 `json:"exchange_rate"`
	Date                 int64   `json:"date"`
	MainCategory         string  `json:"main_category"`
	Subcategory          string  `json:"subcategory"`
	CategoryID           int     `json:"category_id"`
	AccountID            int     `json:"account_id"`
	RelatedAccountID     int     `json:"related_account_id"`
	TransactionType      string  `json:"transaction_type"`
	Fees                 int     `json:"fees"`
	UserID               string  `json:"user_id"`
}

func GetTransactionsByMainCategory(mainCategory string, startDay string, endDay string, uid string) ([]Transaction, error) {
	var start, end int64

	startDate, endDate := timeutils.GetSalaryMonthRange(startDay, endDay)
	start = startDate.Unix()
	end = endDate.Unix()

	query := `SELECT 
	  id,	
	  description,	
	  amount,	
	  currency,	
	  amount_in_base_currency,	
	  exchange_rate,	
	  main_category,	
	  subcategory,	
	  date,	
	  category_id,
	  account_id,
	  related_account_id,
	  transaction_type
	FROM transactions`

	var conditions []string
	var args []interface{}

	conditions = append(conditions, "user_id = ?")
	args = append(args, uid)

	if mainCategory != "" {
		conditions = append(conditions, "main_category = ?")
		args = append(args, mainCategory)
	}

	conditions = append(conditions, "date BETWEEN ? AND ?")
	args = append(args, start, end)

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.AmountInBaseCurrency,
			&transaction.ExchangeRate,
			&transaction.MainCategory,
			&transaction.Subcategory,
			&transaction.Date,
			&transaction.CategoryID,
			&transaction.AccountID,
			&transaction.RelatedAccountID,
			&transaction.TransactionType,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func GetTransactionByID(transactionID int, userID string) (Transaction, error) {
	var transaction Transaction

	query := `SELECT 
		id, description, amount, currency, amount_in_base_currency, exchange_rate, 
		date, main_category, subcategory, category_id, account_id, 
		related_account_id, transaction_type
	FROM transactions 
	WHERE id = ? AND user_id = ?`

	err := db.QueryRow(query, transactionID, userID).Scan(
		&transaction.ID,
		&transaction.Description,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.AmountInBaseCurrency,
		&transaction.ExchangeRate,
		&transaction.Date,
		&transaction.MainCategory,
		&transaction.Subcategory,
		&transaction.CategoryID,
		&transaction.AccountID,
		&transaction.RelatedAccountID,
		&transaction.TransactionType,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Transaction{}, fmt.Errorf("transaction not found")
		}
		return Transaction{}, fmt.Errorf("failed to retrieve transaction: %v", err)
	}

	return transaction, nil
}

func GetTransactions(transactionType string, accountId int, uid string) ([]Transaction, error) {
	query := `SELECT 
	  id,	
	  description,	
	  amount,	
	  currency,	
	  amount_in_base_currency,	
	  exchange_rate,	
	  main_category,	
	  subcategory,	
	  date,	
	  category_id,
	  account_id,
	  related_account_id,
	  transaction_type
	FROM transactions`

	var conditions []string
	var args []interface{}

	conditions = append(conditions, "user_id = ?")
	args = append(args, uid)

	if transactionType != "" {
		conditions = append(conditions, "transaction_type = ?")
		args = append(args, transactionType)
	}

	if accountId > 0 {
		conditions = append(conditions, "account_id = ?")
		args = append(args, accountId)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.AmountInBaseCurrency,
			&transaction.ExchangeRate,
			&transaction.MainCategory,
			&transaction.Subcategory,
			&transaction.Date,
			&transaction.CategoryID,
			&transaction.AccountID,
			&transaction.RelatedAccountID,
			&transaction.TransactionType,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func GetTransactionsForPeriod(start int64, end int64, transactionType string, accountId int, uid string) ([]Transaction, error) {
	query := `SELECT 
	  id,	
	  description,	
	  amount,	
	  currency,	
	  amount_in_base_currency,	
	  exchange_rate,	
	  main_category,	
	  subcategory,	
	  date,	
	  category_id,
	  account_id,
	  related_account_id,
	  transaction_type
	FROM transactions`

	var conditions []string
	var args []interface{}

	conditions = append(conditions, "user_id = ?")
	args = append(args, uid)

	if transactionType != "" {
		conditions = append(conditions, "transaction_type = ?")
		args = append(args, transactionType)
	}

	if accountId > 0 {
		conditions = append(conditions, "account_id = ?")
		args = append(args, accountId)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.AmountInBaseCurrency,
			&transaction.ExchangeRate,
			&transaction.MainCategory,
			&transaction.Subcategory,
			&transaction.Date,
			&transaction.CategoryID,
			&transaction.AccountID,
			&transaction.RelatedAccountID,
			&transaction.TransactionType,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

/**
* Add a new transaction to the database
* Can add TransactionType = "expense", "income"
 */
func AddTransaction(transaction Transaction) (Transaction, error) {
	sourceAccount, err := GetAccountByID(transaction.AccountID, transaction.UserID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid account: %v", err)
	}

	if transaction.TransactionType == TransactionTypeExpense {
		if sourceAccount.Balance < transaction.Amount {
			return Transaction{}, fmt.Errorf("insufficient balance in account: %v", err)
		}
	}

	// Determine the main category based on the subcategory
	mainCategory, err := GetMainCategory(transaction.CategoryID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid subcategory: %v", err)
	}
	subcategory, err := GetSubCategory(transaction.CategoryID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid subcategory: %v", err)
	}
	transaction.MainCategory = mainCategory
	transaction.Subcategory = subcategory

	if transaction.Date == 0 {
		transaction.Date = time.Now().Unix()
	}

	var exchangeRate float64
	var amountInBaseCurrency float64

	rate, err := utils.GetExchangeRate(transaction.Currency)
	if err != nil {
		// Log the error but proceed without exchange rate
		log.Printf("Warning: Exchange rate not found for currency '%s'. Transaction will be saved without conversion.", transaction.Currency)
		exchangeRate = 0
		amountInBaseCurrency = 0
	} else {
		exchangeRate = rate
		// Convert the transaction amount to the base currency
		amountInBaseCurrency = transaction.Amount * exchangeRate
	}

	transaction.ExchangeRate = exchangeRate
	transaction.AmountInBaseCurrency = amountInBaseCurrency

	// Start a database transaction
	tx, err := db.Begin()
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to start database transaction: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// Insert the transaction into the database
	_, err = tx.Exec(
		`INSERT INTO transactions (
		  description,
		  amount,
		  currency,
		  amount_in_base_currency,
		  exchange_rate,
		  date,
		  main_category,
		  subcategory,
		  category_id,
		  account_id,
		  related_account_id,
		  transaction_type,
      user_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		transaction.Description,
		transaction.Amount,
		transaction.Currency,
		transaction.AmountInBaseCurrency,
		transaction.ExchangeRate,
		transaction.Date,
		transaction.MainCategory,
		transaction.Subcategory,
		transaction.CategoryID,
		transaction.AccountID,
		transaction.RelatedAccountID,
		transaction.TransactionType,
		transaction.UserID,
	)
	if err != nil {
		tx.Rollback()
		return Transaction{}, fmt.Errorf("failed to insert transaction: %v", err)
	}

	// Update the account balance for the source account
	_, err = tx.Exec(
		`UPDATE accounts SET balance = balance + ? WHERE id = ?`,
		transaction.Amount, transaction.AccountID,
	)

	if err != nil {
		tx.Rollback()
		return Transaction{}, fmt.Errorf("failed to update source account balance: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to commit database transaction: %v", err)
	}

	return transaction, nil
}

func UpdateTransaction(transactionID int, updatedTransaction Transaction) (Transaction, error) {
	// Retrieve the existing transaction
	existingTransaction, err := GetTransactionByID(transactionID, updatedTransaction.UserID)
	if err != nil {
		return Transaction{}, fmt.Errorf("transaction not found: %v", err)
	}

	// Ensure the category exists
	mainCategory, err := GetMainCategory(updatedTransaction.CategoryID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid subcategory: %v", err)
	}
	subcategory, err := GetSubCategory(updatedTransaction.CategoryID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid subcategory: %v", err)
	}
	updatedTransaction.MainCategory = mainCategory
	updatedTransaction.Subcategory = subcategory

	// Handle currency conversion if the currency changed
	if updatedTransaction.Currency != existingTransaction.Currency {
		rate, err := utils.GetExchangeRate(updatedTransaction.Currency)
		if err != nil {
			log.Printf("Warning: Exchange rate not found for currency '%s'. Transaction will be saved without conversion.", updatedTransaction.Currency)
			updatedTransaction.ExchangeRate = 0
			updatedTransaction.AmountInBaseCurrency = 0
		} else {
			updatedTransaction.ExchangeRate = rate
			updatedTransaction.AmountInBaseCurrency = updatedTransaction.Amount * rate
		}
	} else {
		// Retain the previous exchange rate if the currency hasn't changed
		updatedTransaction.ExchangeRate = existingTransaction.ExchangeRate
		updatedTransaction.AmountInBaseCurrency = updatedTransaction.Amount * existingTransaction.ExchangeRate
	}

	// Start a database transaction
	tx, err := db.Begin()
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to start database transaction: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// Adjust the account balance: First revert the old transaction
	_, err = tx.Exec(
		`UPDATE accounts SET balance = balance - ? WHERE id = ?`,
		existingTransaction.Amount, existingTransaction.AccountID,
	)
	if err != nil {
		tx.Rollback()
		return Transaction{}, fmt.Errorf("failed to revert old transaction amount: %v", err)
	}

	// Then apply the updated transaction amount
	_, err = tx.Exec(
		`UPDATE accounts SET balance = balance + ? WHERE id = ?`,
		updatedTransaction.Amount, updatedTransaction.AccountID,
	)
	if err != nil {
		tx.Rollback()
		return Transaction{}, fmt.Errorf("failed to update account balance: %v", err)
	}

	// Update the transaction in the database
	_, err = tx.Exec(
		`UPDATE transactions SET
		  description = ?, amount = ?, currency = ?, amount_in_base_currency = ?, exchange_rate = ?, 
		  date = ?, main_category = ?, subcategory = ?, category_id = ?, account_id = ?, 
		  related_account_id = ?, transaction_type = ?
		WHERE id = ?`,
		updatedTransaction.Description,
		updatedTransaction.Amount,
		updatedTransaction.Currency,
		updatedTransaction.AmountInBaseCurrency,
		updatedTransaction.ExchangeRate,
		updatedTransaction.Date,
		updatedTransaction.MainCategory,
		updatedTransaction.Subcategory,
		updatedTransaction.CategoryID,
		updatedTransaction.AccountID,
		updatedTransaction.RelatedAccountID,
		updatedTransaction.TransactionType,
		transactionID,
	)
	if err != nil {
		tx.Rollback()
		return Transaction{}, fmt.Errorf("failed to update transaction: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to commit database transaction: %v", err)
	}

	return updatedTransaction, nil
}

/**
* Add a new transfer to the database
* Can add TransactionType = "transfer" "savings"
 */
func AddTransfer(transaction Transaction) (Transaction, error) {
	mainCategory, err := GetMainCategory(transaction.CategoryID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid subcategory: %v", err)
	}
	subcategory, err := GetSubCategory(transaction.CategoryID)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid subcategory: %v", err)
	}
	transaction.MainCategory = mainCategory
	transaction.Subcategory = subcategory

	if transaction.Date == 0 {
		transaction.Date = time.Now().Unix()
	}

	var exchangeRate float64
	var amountInBaseCurrency float64

	rate, err := utils.GetExchangeRate(transaction.Currency)
	if err != nil {
		// Log the error but proceed without exchange rate
		log.Printf("Warning: Exchange rate not found for currency '%s'. Transaction will be saved without conversion.", transaction.Currency)
		exchangeRate = 0
		amountInBaseCurrency = 0
	} else {
		exchangeRate = rate
		// Convert the transaction amount to the base currency
		amountInBaseCurrency = transaction.Amount * exchangeRate
	}

	transaction.ExchangeRate = exchangeRate
	transaction.AmountInBaseCurrency = amountInBaseCurrency

	tx, err := db.Begin()
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to start database transaction: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// Insert the transaction into the database
	_, err = tx.Exec(
		`INSERT INTO transactions (
		  description,
		  amount,
		  currency,
		  amount_in_base_currency,
		  exchange_rate,
		  date,
		  main_category,
		  subcategory,
		  category_id,
		  account_id,
		  related_account_id,
		  transaction_type,
		  fees,
      	  user_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		transaction.Description,
		transaction.Amount,
		transaction.Currency,
		transaction.AmountInBaseCurrency,
		transaction.ExchangeRate,
		transaction.Date,
		transaction.MainCategory,
		transaction.Subcategory,
		transaction.CategoryID,
		transaction.AccountID,
		transaction.RelatedAccountID,
		transaction.TransactionType,
		transaction.Fees,
		transaction.UserID,
	)
	if err != nil {
		tx.Rollback()
		return Transaction{}, fmt.Errorf("failed to insert transaction: %v", err)
	}

	// Update the account balance for the source account
	_, err = tx.Exec(
		`UPDATE accounts SET balance = balance - (? + ?) WHERE id = ?`,
		transaction.Amount, transaction.Fees, transaction.AccountID,
	)

	if err != nil {
		tx.Rollback()
		return Transaction{}, fmt.Errorf("failed to update source account balance: %v", err)
	}

	// Update the account balance for the destination account
	_, err = tx.Exec(
		`UPDATE accounts SET balance = balance + (? + ?) WHERE id = ?`,
		transaction.Amount, transaction.Fees, transaction.RelatedAccountID,
	)

	if err != nil {
		tx.Rollback()
		return Transaction{}, fmt.Errorf("failed to update destination account balance: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to commit database transaction: %v", err)
	}

	return transaction, nil
}

func DeleteTransaction(id int) error {
	// Start a database transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start database transaction: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// Fetch the transaction details to retrieve its amount and account ID
	var transaction struct {
		Amount           float64
		AccountID        int
		RelatedAccountID int // For transfers
		TransactionType  string
		Fees             int
	}

	err = tx.QueryRow(
		`SELECT amount, account_id, related_account_id, transaction_type, fees
		 FROM transactions 
		 WHERE id = ?`, id,
	).Scan(
		&transaction.Amount,
		&transaction.AccountID,
		&transaction.RelatedAccountID,
		&transaction.TransactionType,
		&transaction.Fees,
	)

	if err == sql.ErrNoRows {
		tx.Rollback()
		return fmt.Errorf("transaction with ID %d not found", id)
	} else if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to retrieve transaction: %v", err)
	}

	// Reverse the balance change for the source account
	_, err = tx.Exec(
		`UPDATE accounts SET balance = balance + (? + ?) WHERE id = ?`,
		transaction.Amount, transaction.Fees, transaction.AccountID,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update source account balance: %v", err)
	}

	// If the transaction is a transfer, update the related account balance as well
	if transaction.TransactionType == TransactionTypeTransfer ||
		transaction.TransactionType == TransactionTypeSavings &&
			transaction.RelatedAccountID > 0 {
		_, err = tx.Exec(
			`UPDATE accounts SET balance = balance - (? + ?) WHERE id = ?`,
			transaction.Amount, transaction.Fees, transaction.RelatedAccountID,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update destination account balance: %v", err)
		}
	}

	// Delete the transaction
	_, err = tx.Exec(`DELETE FROM transactions WHERE id = ?`, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete transaction: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit database transaction: %v", err)
	}

	return nil
}

func GetTransactionsByAccount(accountID string, uid string) ([]Transaction, error) {
	// Prepare the SQL query
	query := `
		SELECT 
			id, 
			description, 
			amount, 
			currency, 
			amount_in_base_currency, 
			exchange_rate, 
			date, 
			main_category, 
			subcategory, 
			category_id, 
			account_id, 
			related_account_id, 
			transaction_type, 
			fees 
		FROM transactions 
		WHERE account_id = ? AND user_id = ?
	`

	// Initialize a slice to hold the transactions
	var transactions []Transaction

	// Execute the query
	rows, err := db.Query(query, accountID, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	defer rows.Close() // Ensure rows are closed

	// Iterate through the result set and scan into the slice
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.AmountInBaseCurrency,
			&transaction.ExchangeRate,
			&transaction.Date,
			&transaction.MainCategory,
			&transaction.Subcategory,
			&transaction.CategoryID,
			&transaction.AccountID,
			&transaction.RelatedAccountID,
			&transaction.TransactionType,
			&transaction.Fees,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
	}

	return transactions, nil
}

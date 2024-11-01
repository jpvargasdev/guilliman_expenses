package models

import (
	"time"
  "fmt"
  "log"

  "guilliman/internal/utils"
)

type Expense struct {
  ID                    int     `json:"id"`
  Description           string  `json:"description"`
  TransactionAmount     float64 `json:"transaction_amount"`        // Amount in transaction currency
  TransactionCurrency   string  `json:"transaction_currency"`      // Currency code of the transaction
  AmountInBaseCurrency  float64 `json:"amount_in_base_currency"`   // Amount converted to base currency
  ExchangeRate          float64 `json:"exchange_rate"`             // Exchange rate used for conversion
  MainCategory          string  `json:"main_category"`             // Needs, Wants, Savings
  Subcategory           string  `json:"subcategory"`               // Specific subcategory
  Date                  int64   `json:"date"`                      // Unix timestamp
  CategoryID            int     `json:"category_id"`
}


func GetExpenses() ([]Expense, error) {
  rows, err := db.Query(`
    SELECT 
      id, 
      description, 
      transaction_amount, 
      transaction_currency, 
      amount_in_base_currency, 
      exchange_rate, 
      main_category, 
      subcategory, 
      date, 
      category_id
    FROM expenses
  `)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var expenses []Expense
  for rows.Next() {
    var expense Expense
    if err := rows.Scan(
      &expense.ID,
      &expense.Description,
      &expense.TransactionAmount,
      &expense.TransactionCurrency,
      &expense.AmountInBaseCurrency,
      &expense.ExchangeRate,
      &expense.MainCategory,
      &expense.Subcategory,
      &expense.Date,
      &expense.CategoryID,
    ); err != nil {
      return nil, err
    }
    expenses = append(expenses, expense)
  }
  return expenses, nil 
}

func GetExpensesForPeriod(start, end int64) ([]Expense, error) {
  rows, err := db.Query(`
    SELECT 
      id, 
      description, 
      transaction_amount, 
      transaction_currency, 
      amount_in_base_currency, 
      exchange_rate, 
      main_category, 
      subcategory, 
      date, 
      category_id
    FROM expenses
    WHERE e.date BETWEEN ? AND ?`, start, end)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

 var expenses []Expense
  for rows.Next() {
      var expense Expense
      if err := rows.Scan(
        &expense.ID, 
        &expense.Description, 
        &expense.TransactionAmount,
        &expense.TransactionCurrency,
        &expense.AmountInBaseCurrency,
        &expense.ExchangeRate,
        &expense.MainCategory,
        &expense.Subcategory,
        &expense.Date, 
        &expense.CategoryID, 
      ); err != nil {
        return nil, err
      }
      expenses = append(expenses, expense)
  }
  return expenses, nil
}

func AddExpense(expense Expense) error {
  // Determine the main category based on the subcategory
  mainCategory, err := GetMainCategory(expense.CategoryID)
  subcategory, err := GetSubCategory(expense.CategoryID)
  if err != nil {
    return fmt.Errorf("invalid subcategory: %v", err)
  }
  expense.MainCategory = mainCategory
  expense.Subcategory = subcategory

  if expense.Date == 0 {
    expense.Date = time.Now().Unix()
  }

  var exchangeRate float64
  var amountInBaseCurrency float64

  rate, err := utils.GetExchangeRate(expense.TransactionCurrency)
  if err != nil {
    // Log the error but proceed without exchange rate
    log.Printf("Warning: Exchange rate not found for currency '%s'. Expense will be saved without conversion.", expense.TransactionCurrency)
    exchangeRate = 0
    amountInBaseCurrency = 0
  } else {
    exchangeRate = rate
    // Convert the transaction amount to the base currency
    amountInBaseCurrency = expense.TransactionAmount * exchangeRate
  }

  expense.ExchangeRate = exchangeRate
  expense.AmountInBaseCurrency = amountInBaseCurrency

  // Insert the expense into the database
  _, err = db.Exec(
    `INSERT INTO expenses (
      description,
      transaction_amount,
      transaction_currency,
      amount_in_base_currency,
      exchange_rate,
      main_category,
      subcategory,
      date,
      category_id
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
    expense.Description,
    expense.TransactionAmount,
    expense.TransactionCurrency,
    expense.AmountInBaseCurrency, // May be zero or null
    expense.ExchangeRate,         // May be zero or null
    expense.MainCategory,
    expense.Subcategory,
    expense.Date,
    expense.CategoryID,
  )
  if err != nil {
    return err
  }

  return nil
}


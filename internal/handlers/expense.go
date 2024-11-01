package handlers

import (
  "time"
  "net/http"
  "strconv"

  "guilliman/internal/utils/timeutils"
  "guilliman/internal/models"

  "github.com/gin-gonic/gin"
)

func AddExpenseHandler(c *gin.Context) {
  var expense models.Expense
  if err := c.ShouldBindJSON(&expense); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
  }
  models.AddExpense(expense)
  c.JSON(http.StatusCreated, expense)
}

func GetExpensesHandler(c *gin.Context) {
  expenses, err := models.GetExpenses() 
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, expenses)
}

func GetExpensesForPeriodHandler(c *gin.Context) {
  dateParam := c.Query("date")
  var date time.Time
  if dateParam == "" {
      date = time.Now()  
  } else {
    timestamp, err := strconv.ParseInt(dateParam, 10, 64)
    if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use a Unix timestamp."})
      return
    }
    date = time.Unix(timestamp, 0)
  }

  start, end := timeutils.CalculatePeriodBoundaries(date)

  expenses, err := models.GetExpensesForPeriod(start, end)
  if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }

  c.JSON(http.StatusOK, expenses)
}

package controller

import (
	"net/http"

	"guilliman/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Controller) GetSavingsController(c *gin.Context) {
	accountParam := c.Query("account")

	expenses, err := models.GetTransactions(models.TransactionTypeSavings, accountParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}
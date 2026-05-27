package handlers

import (
	"encoding/json"
	"net/http"

	"xendit-va/database"
	"xendit-va/models"
	"xendit-va/services"

	"github.com/gin-gonic/gin"
)

type XenditWebhook struct {
	ExternalID   string  `json:"external_id"`
	Amount       float64 `json:"amount"`
	Status       string  `json:"status"`
	BankCode     string  `json:"bank_code"`
	AccountNumber string `json:"account_number"`
}

func XenditWebhookHandler(c *gin.Context) {

	var payload XenditWebhook

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var transaction models.Transaction

	err := database.DB.
		Where("external_id = ?", payload.ExternalID).
		First(&transaction).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Transaction not found",
		})
		return
	}

	transaction.Status = "PAID"
	transaction.PaymentReceived = true

	database.DB.Save(&transaction)

	event, _ := json.Marshal(transaction)

	services.PublishMessage(string(event))

	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook processed",
	})
}
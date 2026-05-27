package handlers

import (
	"net/http"

	"xendit-va/database"
	"xendit-va/models"
	"xendit-va/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateVAInput struct {
	BankCode string  `json:"bank_code"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
}

func CreateVA(c *gin.Context) {

	var input CreateVAInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	externalID := uuid.New().String()

	vaResponse, err := services.CreateVirtualAccount(
		services.CreateVARequest{
			ExternalID: externalID,
			BankCode:   input.BankCode,
			Name:       input.Name,
			Amount:     input.Amount,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	transaction := models.Transaction{
		ExternalID:      externalID,
		BankCode:        vaResponse.BankCode,
		AccountNumber:   vaResponse.AccountNumber,
		Name:            vaResponse.Name,
		Amount:          input.Amount,
		Status:          "PENDING",
		XenditVAID:      vaResponse.ID,
		PaymentReceived: false,
	}

	database.DB.Create(&transaction)

	c.JSON(http.StatusOK, gin.H{
		"message": "VA Created",
		"data":    transaction,
	})
}
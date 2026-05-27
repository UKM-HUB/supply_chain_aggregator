package controllers

import (
	"net/http"
	"project/config"
	"project/models"

	"github.com/gin-gonic/gin"
)

func CreateTransaction(c *gin.Context) {
	var trx models.Transaction

	if err := c.ShouldBindJSON(&trx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	config.DB.Create(&trx)

	c.JSON(http.StatusOK, trx)
}

func GetTransactions(c *gin.Context) {
	var trx []models.Transaction

	config.DB.Find(&trx)

	c.JSON(http.StatusOK, trx)
}
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReports(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Report data",
	})
}
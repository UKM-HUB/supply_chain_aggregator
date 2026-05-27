package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Communication(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Communication API",
	})
}
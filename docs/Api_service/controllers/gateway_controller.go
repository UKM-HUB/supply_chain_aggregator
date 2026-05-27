package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Gateway(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Gateway API",
	})
}
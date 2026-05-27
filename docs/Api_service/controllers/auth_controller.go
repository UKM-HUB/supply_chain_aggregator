package controllers

import (
	"net/http"
	"project/config"
	"project/models"
	"project/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(hash)

	config.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "Register success",
	})
}

func Login(c *gin.Context) {
	var input models.User
	var user models.User

	c.ShouldBindJSON(&input)

	config.DB.Where("email = ?", input.Email).First(&user)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials",
		})
		return
	}

	token, _ := utils.GenerateJWT(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
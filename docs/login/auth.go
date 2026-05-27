package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("SECRET_KEY")

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GenerateToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"id":   user.ID,
		"name": user.Name,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
		Role:     "ADMIN",
	}

	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email already exists",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Register success",
	})
}

func Login(c *gin.Context) {
	var input LoginInput
	var user User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	if err := DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email/password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email/password",
		})
		return
	}

	token, err := GenerateToken(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed generate token",
		})
		return
	}

	refreshToken := uuid.New().String()

	user.RefreshToken = refreshToken
	DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":   user.ID,
			"name": user.Name,
			"role": user.Role,
		},
	})
}

func Refresh(c *gin.Context) {
	type RefreshInput struct {
		RefreshToken string `json:"refresh_token"`
	}

	var input RefreshInput
	var user User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	if err := DB.Where("refresh_token = ?", input.RefreshToken).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid refresh token",
		})
		return
	}

	token, err := GenerateToken(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func Logout(c *gin.Context) {
	type LogoutInput struct {
		RefreshToken string `json:"refresh_token"`
	}

	var input LogoutInput
	var user User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	if err := DB.Where("refresh_token = ?", input.RefreshToken).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		return
	}

	user.RefreshToken = ""
	DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout success",
	})
}
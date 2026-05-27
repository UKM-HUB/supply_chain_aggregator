package handlers

import (
	"fmt"
	"math/rand"
	"net/http"

	"register-api/config"
	"register-api/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type VerifyRequest struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

// Generate OTP 6 digit
func generateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// POST /api/v1/register
func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	// cek user
	var existing models.User

	config.DB.Where("phone = ?", req.Phone).First(&existing)

	if existing.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Phone already registered",
		})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed hash password",
		})
		return
	}

	otp := generateOTP()

	user := models.User{
		Name:       req.Name,
		Phone:      req.Phone,
		Password:   string(hashedPassword),
		OTP:        otp,
		IsVerified: false,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed register",
		})
		return
	}

	// simulasi kirim OTP
	fmt.Println("OTP:", otp)

	c.JSON(http.StatusOK, gin.H{
		"message": "Register success, verify OTP",
		"otp": otp,
	})
}

// POST /api/v1/register/verify
func VerifyOTP(c *gin.Context) {
	var req VerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	var user models.User

	if err := config.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	if user.OTP != req.OTP {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid OTP",
		})
		return
	}

	user.IsVerified = true
	user.OTP = ""

	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verified, user activated",
	})
}
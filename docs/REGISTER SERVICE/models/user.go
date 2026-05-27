package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string `json:"name"`
	Phone     string `json:"phone" gorm:"unique"`
	Password  string `json:"-"`
	OTP       string `json:"otp"`
	IsVerified bool  `json:"is_verified"`
}
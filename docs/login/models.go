package main

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"unique"`
	Password     string `json:"-"`
	Role         string `json:"role"`
	RefreshToken string `json:"-"`
}
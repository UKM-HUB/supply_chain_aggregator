package entity

import "time"

type User struct {
	ID           string
	Name         string
	Email        string
	Phone        string
	PasswordHash string
	Role         string
	Latitude     float64
	Longitude    float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

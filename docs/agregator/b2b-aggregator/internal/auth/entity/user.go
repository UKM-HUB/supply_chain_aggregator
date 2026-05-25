package entity

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Role      string    `gorm:"not null"`
	Latitude  float64
	Longitude float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

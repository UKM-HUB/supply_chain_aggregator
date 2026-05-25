package entity

import "time"

type Product struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductCode string    `gorm:"uniqueIndex;not null"`
	Name        string    `gorm:"not null"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

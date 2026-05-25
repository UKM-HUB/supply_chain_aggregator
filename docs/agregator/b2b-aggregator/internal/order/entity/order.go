package entity

import "time"

type Order struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FactoryID   string    `gorm:"not null"`
	UMKMID      string    `gorm:"not null"`
	ProductCode string    `gorm:"not null"`
	Quantity    int32     `gorm:"not null"`
	Status      string    `gorm:"not null"` // PENDING, PAID, CANCELED
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SubOrder struct {
	UMKMID      string
	ProductCode string
	Quantity    int32
}

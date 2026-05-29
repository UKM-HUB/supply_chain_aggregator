package entity

import "time"

const (
	StatusPending   = "pending"
	StatusPaid      = "paid"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"
)

type Transaction struct {
	ID            string
	InvoiceNumber string
	UserID        string
	Amount        float64
	Status        string
	PaymentMethod string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

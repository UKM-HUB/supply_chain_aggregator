package entity

import "time"

const (
	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"
	PaymentStatusFailed  = "failed"
)

type PaymentLog struct {
	ID             string
	InvoiceNumber  string
	Amount         float64
	UserPhone      string
	PaymentURL     string
	XenditInvoiceID string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

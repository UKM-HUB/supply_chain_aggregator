package entity

import "time"

type TransactionRecord struct {
	ID            string
	InvoiceNumber string
	UserID        string
	Amount        float64
	Status        string
	CreatedAt     time.Time
}

type ReportSummary struct {
	TotalTransaction int     `json:"total_transaction"`
	TotalPaid        float64 `json:"total_paid"`
	TotalPending     int     `json:"total_pending"`
}

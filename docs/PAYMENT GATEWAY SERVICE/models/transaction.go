package models

import "time"

type Transaction struct {
	ID              uint      `gorm:"primaryKey"`
	ExternalID      string    `json:"external_id"`
	BankCode        string    `json:"bank_code"`
	AccountNumber   string    `json:"account_number"`
	Name            string    `json:"name"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	XenditVAID      string    `json:"xendit_va_id"`
	PaymentReceived bool      `json:"payment_received"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
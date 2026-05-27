package main

import "time"

type Transaction struct {
    ID            string    `json:"id"`
    InvoiceNumber string    `json:"invoice_number"`
    UserID        string    `json:"user_id"`
    Amount        float64   `json:"amount"`
    Status        string    `json:"status"`
    PaymentMethod string    `json:"payment_method"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
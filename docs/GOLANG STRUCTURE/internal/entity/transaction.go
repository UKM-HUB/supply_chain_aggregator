package entity

import "github.com/google/uuid"

type Transaction struct {
    ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
    InvoiceNumber string
    Amount        float64
    Status        string
    UserID        uuid.UUID
}
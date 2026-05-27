package repository

import (
    "project/internal/entity"

    "gorm.io/gorm"
)

type TransactionRepository struct {
    DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
    return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) Create(tx *entity.Transaction) error {
    return r.DB.Create(tx).Error
}

func (r *TransactionRepository) UpdateStatus(id string, status string) error {
    return r.DB.Model(&entity.Transaction{}).
        Where("id = ?", id).
        Update("status", status).Error
}
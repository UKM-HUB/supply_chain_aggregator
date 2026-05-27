package usecase

import (
    "project/internal/entity"
    "project/internal/helper"
    "project/internal/repository"

    "github.com/google/uuid"
)

type TransactionUsecase struct {
    Repo *repository.TransactionRepository
}

func NewTransactionUsecase(repo *repository.TransactionRepository) *TransactionUsecase {
    return &TransactionUsecase{Repo: repo}
}

func (u *TransactionUsecase) Create(tx *entity.Transaction) error {
    tx.ID = uuid.New()
    tx.Status = "PENDING"

    err := u.Repo.Create(tx)
    if err != nil {
        return err
    }

    return helper.Publish("transaction.created", tx)
}
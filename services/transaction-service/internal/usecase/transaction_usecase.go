package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"supply-chain-aggregator/services/transaction-service/internal/entity"
	"supply-chain-aggregator/services/transaction-service/internal/repository"
)

var (
	ErrNotFound          = errors.New("transaction not found")
	ErrInvalidAmount     = errors.New("amount must be greater than zero")
	ErrInvalidUserID     = errors.New("user_id is required")
	ErrInvalidStatus     = errors.New("invalid status value")
	ErrInvalidPayment    = errors.New("payment_method is required")

	invoiceCounter uint64 = 1000
)

var validStatuses = map[string]bool{
	entity.StatusPending:   true,
	entity.StatusPaid:      true,
	entity.StatusFailed:    true,
	entity.StatusCancelled: true,
}

type CreateTransactionInput struct {
	UserID        string
	Amount        float64
	PaymentMethod string
}

type TransactionUsecase struct {
	repo repository.TransactionRepository
}

func NewTransactionUsecase(repo repository.TransactionRepository) *TransactionUsecase {
	return &TransactionUsecase{repo: repo}
}

func (u *TransactionUsecase) Create(ctx context.Context, input CreateTransactionInput) (entity.Transaction, error) {
	input.UserID = strings.TrimSpace(input.UserID)
	input.PaymentMethod = strings.TrimSpace(input.PaymentMethod)

	if input.UserID == "" {
		return entity.Transaction{}, ErrInvalidUserID
	}
	if input.Amount <= 0 {
		return entity.Transaction{}, ErrInvalidAmount
	}
	if input.PaymentMethod == "" {
		return entity.Transaction{}, ErrInvalidPayment
	}

	now := time.Now()
	tx := entity.Transaction{
		ID:            generateID(),
		InvoiceNumber: generateInvoiceNumber(now),
		UserID:        input.UserID,
		Amount:        input.Amount,
		Status:        entity.StatusPending,
		PaymentMethod: input.PaymentMethod,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := u.repo.Create(ctx, tx); err != nil {
		return entity.Transaction{}, err
	}

	return tx, nil
}

func (u *TransactionUsecase) List(ctx context.Context, userID string) ([]entity.Transaction, error) {
	return u.repo.List(ctx, strings.TrimSpace(userID))
}

func (u *TransactionUsecase) GetByID(ctx context.Context, id string) (entity.Transaction, error) {
	tx, err := u.repo.GetByID(ctx, strings.TrimSpace(id))
	if errors.Is(err, repository.ErrNotFound) {
		return entity.Transaction{}, ErrNotFound
	}
	return tx, err
}

func (u *TransactionUsecase) UpdateStatus(ctx context.Context, id, status string) error {
	status = strings.ToLower(strings.TrimSpace(status))
	if !validStatuses[status] {
		return ErrInvalidStatus
	}
	err := u.repo.UpdateStatus(ctx, strings.TrimSpace(id), status)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func generateInvoiceNumber(t time.Time) string {
	seq := atomic.AddUint64(&invoiceCounter, 1)
	return fmt.Sprintf("INV-%s-%04d", t.Format("20060102"), seq)
}

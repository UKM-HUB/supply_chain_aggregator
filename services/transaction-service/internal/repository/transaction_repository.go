package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"supply-chain-aggregator/services/transaction-service/internal/entity"
)

var ErrNotFound = errors.New("transaction not found")

type TransactionRepository interface {
	Create(ctx context.Context, tx entity.Transaction) error
	List(ctx context.Context, userID string) ([]entity.Transaction, error)
	GetByID(ctx context.Context, id string) (entity.Transaction, error)
	UpdateStatus(ctx context.Context, id, status string) error
}

type InMemoryTransactionRepository struct {
	mu           sync.RWMutex
	transactions []entity.Transaction
}

func NewInMemoryTransactionRepository() *InMemoryTransactionRepository {
	now := time.Now()
	return &InMemoryTransactionRepository{
		transactions: []entity.Transaction{
			{
				ID:            "txn-seed-001",
				InvoiceNumber: "INV-20260526-0001",
				UserID:        "user-corp-001",
				Amount:        5000000,
				Status:        entity.StatusPaid,
				PaymentMethod: "virtual_account",
				CreatedAt:     now.Add(-48 * time.Hour),
				UpdatedAt:     now.Add(-47 * time.Hour),
			},
			{
				ID:            "txn-seed-002",
				InvoiceNumber: "INV-20260526-0002",
				UserID:        "user-corp-001",
				Amount:        2500000,
				Status:        entity.StatusPending,
				PaymentMethod: "virtual_account",
				CreatedAt:     now.Add(-2 * time.Hour),
				UpdatedAt:     now.Add(-2 * time.Hour),
			},
			{
				ID:            "txn-seed-003",
				InvoiceNumber: "INV-20260526-0003",
				UserID:        "user-corp-002",
				Amount:        7800000,
				Status:        entity.StatusPending,
				PaymentMethod: "bank_transfer",
				CreatedAt:     now.Add(-1 * time.Hour),
				UpdatedAt:     now.Add(-1 * time.Hour),
			},
		},
	}
}

func (r *InMemoryTransactionRepository) Create(ctx context.Context, tx entity.Transaction) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.transactions = append(r.transactions, tx)
	return nil
}

func (r *InMemoryTransactionRepository) List(ctx context.Context, userID string) ([]entity.Transaction, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]entity.Transaction, 0)
	for _, tx := range r.transactions {
		if userID != "" && tx.UserID != userID {
			continue
		}
		result = append(result, tx)
	}

	return result, nil
}

func (r *InMemoryTransactionRepository) GetByID(ctx context.Context, id string) (entity.Transaction, error) {
	select {
	case <-ctx.Done():
		return entity.Transaction{}, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, tx := range r.transactions {
		if tx.ID == id {
			return tx, nil
		}
	}

	return entity.Transaction{}, ErrNotFound
}

func (r *InMemoryTransactionRepository) UpdateStatus(ctx context.Context, id, status string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for i, tx := range r.transactions {
		if tx.ID == id {
			r.transactions[i].Status = status
			r.transactions[i].UpdatedAt = time.Now()
			return nil
		}
	}

	return ErrNotFound
}

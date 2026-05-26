package repository

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"sync"

	"supply-chain-aggregator/services/payment-service/internal/entity"
)

var ErrNotFound = errors.New("payment log not found")

type PaymentRepository interface {
	Create(ctx context.Context, log entity.PaymentLog) error
	GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (entity.PaymentLog, error)
	UpdateStatus(ctx context.Context, invoiceNumber, status, xenditInvoiceID string) error
}

type InMemoryPaymentRepository struct {
	mu   sync.RWMutex
	logs []entity.PaymentLog
}

func NewInMemoryPaymentRepository() *InMemoryPaymentRepository {
	return &InMemoryPaymentRepository{
		logs: make([]entity.PaymentLog, 0),
	}
}

func (r *InMemoryPaymentRepository) Create(ctx context.Context, log entity.PaymentLog) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if log.ID == "" {
		log.ID = generateID()
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs = append(r.logs, log)
	return nil
}

func (r *InMemoryPaymentRepository) GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (entity.PaymentLog, error) {
	select {
	case <-ctx.Done():
		return entity.PaymentLog{}, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, log := range r.logs {
		if log.InvoiceNumber == invoiceNumber {
			return log, nil
		}
	}

	return entity.PaymentLog{}, ErrNotFound
}

func (r *InMemoryPaymentRepository) UpdateStatus(ctx context.Context, invoiceNumber, status, xenditInvoiceID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for i, log := range r.logs {
		if log.InvoiceNumber == invoiceNumber {
			r.logs[i].Status = status
			if xenditInvoiceID != "" {
				r.logs[i].XenditInvoiceID = xenditInvoiceID
			}
			return nil
		}
	}

	return ErrNotFound
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

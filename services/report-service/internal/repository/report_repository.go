package repository

import (
	"context"
	"sync"
	"time"

	"supply-chain-aggregator/services/report-service/internal/entity"
)

type DateRange struct {
	From time.Time
	To   time.Time
}

type ReportRepository interface {
	ListByDateRange(ctx context.Context, r DateRange) ([]entity.TransactionRecord, error)
}

type InMemoryReportRepository struct {
	mu      sync.RWMutex
	records []entity.TransactionRecord
}

func NewInMemoryReportRepository() *InMemoryReportRepository {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return &InMemoryReportRepository{
		records: []entity.TransactionRecord{
			{
				ID:            "txn-seed-001",
				InvoiceNumber: "INV-20260524-1001",
				UserID:        "user-corp-001",
				Amount:        5000000,
				Status:        "paid",
				CreatedAt:     today.AddDate(0, 0, -2).Add(10 * time.Hour),
			},
			{
				ID:            "txn-seed-002",
				InvoiceNumber: "INV-20260524-1002",
				UserID:        "user-corp-002",
				Amount:        1200000,
				Status:        "paid",
				CreatedAt:     today.AddDate(0, 0, -2).Add(14 * time.Hour),
			},
			{
				ID:            "txn-seed-003",
				InvoiceNumber: "INV-20260525-1003",
				UserID:        "user-corp-001",
				Amount:        3800000,
				Status:        "pending",
				CreatedAt:     today.AddDate(0, 0, -1).Add(9 * time.Hour),
			},
			{
				ID:            "txn-seed-004",
				InvoiceNumber: "INV-20260525-1004",
				UserID:        "user-corp-003",
				Amount:        2500000,
				Status:        "paid",
				CreatedAt:     today.AddDate(0, 0, -1).Add(16 * time.Hour),
			},
			{
				ID:            "txn-seed-005",
				InvoiceNumber: "INV-20260526-1005",
				UserID:        "user-corp-002",
				Amount:        7800000,
				Status:        "paid",
				CreatedAt:     today.Add(8 * time.Hour),
			},
			{
				ID:            "txn-seed-006",
				InvoiceNumber: "INV-20260526-1006",
				UserID:        "user-corp-003",
				Amount:        4100000,
				Status:        "pending",
				CreatedAt:     today.Add(11 * time.Hour),
			},
			{
				ID:            "txn-seed-007",
				InvoiceNumber: "INV-20260526-1007",
				UserID:        "user-corp-001",
				Amount:        950000,
				Status:        "failed",
				CreatedAt:     today.Add(13 * time.Hour),
			},
		},
	}
}

func (r *InMemoryReportRepository) ListByDateRange(ctx context.Context, dr DateRange) ([]entity.TransactionRecord, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]entity.TransactionRecord, 0)
	for _, rec := range r.records {
		if (rec.CreatedAt.Equal(dr.From) || rec.CreatedAt.After(dr.From)) &&
			rec.CreatedAt.Before(dr.To) {
			result = append(result, rec)
		}
	}

	return result, nil
}

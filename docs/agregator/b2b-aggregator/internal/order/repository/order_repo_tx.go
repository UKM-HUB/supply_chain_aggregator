package repository

import (
	"context"
	"b2b-aggregator/internal/order/entity"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Implementasi Database Transaction (ACID)
func (r *OrderRepository) CreateSplitOrdersWithTx(ctx context.Context, factoryID string, orders []entity.SubOrder) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, sub := range orders {
		newOrder := entity.Order{
			FactoryID:   factoryID,
			UMKMID:      sub.UMKMID,
			ProductCode: sub.ProductCode,
			Quantity:    sub.Quantity,
			Status:      "PENDING",
		}
		if err := tx.Create(&newOrder).Error; err != nil {
			tx.Rollback() // Batalkan semua jika ada 1 yang gagal
			return err
		}
	}
	return tx.Commit().Error
}

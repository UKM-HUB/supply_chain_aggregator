package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"b2b-aggregator/internal/order/entity"
	"b2b-aggregator/internal/order/repository"
	"b2b-aggregator/pkg/rabbitmq"
)

type OrderUsecase struct {
	repo      *repository.OrderRepository
	publisher *rabbitmq.Publisher
}

func NewOrderUsecase(repo *repository.OrderRepository, pub *rabbitmq.Publisher) *OrderUsecase {
	return &OrderUsecase{repo: repo, publisher: pub}
}

func (u *OrderUsecase) ProcessCheckout(ctx context.Context, factoryID, productCode string, qty int32) error {
	if qty <= 0 {
		return errors.New("kuantitas pesanan tidak valid")
	}

	// SIMULASI: Hasil dari Inventory Service via gRPC
	// Di sistem nyata, kita akan memanggil gRPC client di sini
	mockSupplies := []entity.SubOrder{
		{UMKMID: "UMKM-1", ProductCode: productCode, Quantity: qty}, // Anggap UMKM-1 punya stok cukup
	}

	// Eksekusi Transaction
	err := u.repo.CreateSplitOrdersWithTx(ctx, factoryID, mockSupplies)
	if err != nil {
		return err
	}

	// Publish ke RabbitMQ
	eventData, _ := json.Marshal(map[string]interface{}{
		"factory_id": factoryID,
		"status":     "ORDER_CREATED",
	})
	
	// Fire and forget
	go u.publisher.Publish(context.Background(), "order_events", eventData)

	return nil
}

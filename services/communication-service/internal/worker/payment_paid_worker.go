package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"supply-chain-aggregator/services/communication-service/internal/consumer"
	"supply-chain-aggregator/services/communication-service/internal/entity"
	"supply-chain-aggregator/services/communication-service/internal/whatsapp"
)

const queuePaymentPaid = "payment.paid"

type PaymentPaidWorker struct {
	consumer  *consumer.Consumer
	whatsapp  *whatsapp.Client
}

func NewPaymentPaidWorker(c *consumer.Consumer, wa *whatsapp.Client) *PaymentPaidWorker {
	return &PaymentPaidWorker{consumer: c, whatsapp: wa}
}

func (w *PaymentPaidWorker) Start(ctx context.Context) error {
	return w.consumer.Consume(ctx, queuePaymentPaid, w.handle)
}

func (w *PaymentPaidWorker) handle(ctx context.Context, body []byte) error {
	var event entity.PaymentPaidEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to parse payment.paid event: %w", err)
	}

	if event.Phone == "" || event.Invoice == "" {
		log.Printf("[worker] invalid event, missing phone or invoice: %s\n", string(body))
		return nil
	}

	message := whatsapp.FormatPaymentMessage(event.Invoice, event.Amount)

	log.Printf("[worker] sending WhatsApp to phone=%s invoice=%s amount=%.0f\n",
		event.Phone, event.Invoice, event.Amount)

	if err := w.whatsapp.Send(ctx, event.Phone, message); err != nil {
		return fmt.Errorf("failed to send WhatsApp notification: %w", err)
	}

	log.Printf("[worker] notification sent for invoice=%s\n", event.Invoice)
	return nil
}

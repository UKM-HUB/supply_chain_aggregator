package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"supply-chain-aggregator/services/payment-service/internal/entity"
	"supply-chain-aggregator/services/payment-service/internal/rabbitmq"
	"supply-chain-aggregator/services/payment-service/internal/repository"
	"supply-chain-aggregator/services/payment-service/internal/xendit"
)

const queuePaymentPaid = "payment.paid"

var (
	ErrInvalidInvoice     = errors.New("invoice_number is required")
	ErrInvalidAmount      = errors.New("amount must be greater than zero")
	ErrInvalidPhone       = errors.New("user_phone is required")
	ErrUnauthorized       = errors.New("invalid callback token")
	ErrNotFound           = errors.New("payment log not found")
)

type CreateVAInput struct {
	InvoiceNumber string
	Amount        float64
	UserPhone     string
}

type PaymentUsecase struct {
	repo          repository.PaymentRepository
	xenditClient  *xendit.Client
	publisher     *rabbitmq.Publisher
	callbackToken string
}

func NewPaymentUsecase(
	repo repository.PaymentRepository,
	xenditClient *xendit.Client,
	publisher *rabbitmq.Publisher,
	callbackToken string,
) *PaymentUsecase {
	return &PaymentUsecase{
		repo:          repo,
		xenditClient:  xenditClient,
		publisher:     publisher,
		callbackToken: callbackToken,
	}
}

func (u *PaymentUsecase) CreateVirtualAccount(ctx context.Context, input CreateVAInput) (entity.PaymentLog, error) {
	input.InvoiceNumber = strings.TrimSpace(input.InvoiceNumber)
	input.UserPhone = strings.TrimSpace(input.UserPhone)

	if input.InvoiceNumber == "" {
		return entity.PaymentLog{}, ErrInvalidInvoice
	}
	if input.Amount <= 0 {
		return entity.PaymentLog{}, ErrInvalidAmount
	}
	if input.UserPhone == "" {
		return entity.PaymentLog{}, ErrInvalidPhone
	}

	xenditResp, err := u.xenditClient.CreateInvoice(ctx, entity.XenditCreateInvoiceRequest{
		ExternalID:  input.InvoiceNumber,
		Amount:      input.Amount,
		Currency:    "IDR",
		Description: fmt.Sprintf("Payment for %s", input.InvoiceNumber),
	})
	if err != nil {
		return entity.PaymentLog{}, fmt.Errorf("xendit create invoice failed: %w", err)
	}

	now := time.Now()
	log := entity.PaymentLog{
		InvoiceNumber:   input.InvoiceNumber,
		Amount:          input.Amount,
		UserPhone:       input.UserPhone,
		PaymentURL:      xenditResp.InvoiceURL,
		XenditInvoiceID: xenditResp.ID,
		Status:          entity.PaymentStatusPending,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := u.repo.Create(ctx, log); err != nil {
		return entity.PaymentLog{}, err
	}

	return log, nil
}

func (u *PaymentUsecase) HandleWebhook(ctx context.Context, callbackToken string, payload entity.XenditWebhookPayload) error {
	if u.callbackToken != "" && callbackToken != u.callbackToken {
		return ErrUnauthorized
	}

	if strings.ToUpper(payload.Status) != "PAID" {
		return nil
	}

	if err := u.repo.UpdateStatus(ctx, payload.ExternalID, entity.PaymentStatusPaid, payload.ID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	log, err := u.repo.GetByInvoiceNumber(ctx, payload.ExternalID)
	if err != nil {
		return err
	}

	event, _ := json.Marshal(map[string]interface{}{
		"invoice": log.InvoiceNumber,
		"amount":  log.Amount,
		"phone":   log.UserPhone,
	})

	go u.publisher.Publish(context.Background(), queuePaymentPaid, event)

	return nil
}

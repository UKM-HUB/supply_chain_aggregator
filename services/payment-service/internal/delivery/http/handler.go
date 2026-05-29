package http

import (
	"errors"
	"net/http"

	"supply-chain-aggregator/services/payment-service/internal/entity"
	"supply-chain-aggregator/services/payment-service/internal/usecase"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	paymentUsecase *usecase.PaymentUsecase
}

func NewHandler(paymentUsecase *usecase.PaymentUsecase) *Handler {
	return &Handler{paymentUsecase: paymentUsecase}
}

func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "payment-service",
	})
}

type createVARequest struct {
	InvoiceNumber string  `json:"invoice_number"`
	Amount        float64 `json:"amount"`
	UserPhone     string  `json:"user_phone"`
}

func (h *Handler) CreateVirtualAccount(c echo.Context) error {
	var req createVARequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	log, err := h.paymentUsecase.CreateVirtualAccount(c.Request().Context(), usecase.CreateVAInput{
		InvoiceNumber: req.InvoiceNumber,
		Amount:        req.Amount,
		UserPhone:     req.UserPhone,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidInvoice):
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invoice_number is required"})
		case errors.Is(err, usecase.ErrInvalidAmount):
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "amount must be greater than zero"})
		case errors.Is(err, usecase.ErrInvalidPhone):
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "user_phone is required"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to create virtual account"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": map[string]interface{}{
			"invoice_number": log.InvoiceNumber,
			"payment_url":    log.PaymentURL,
			"amount":         log.Amount,
			"status":         log.Status,
		},
	})
}

func (h *Handler) HandleXenditWebhook(c echo.Context) error {
	callbackToken := c.Request().Header.Get("x-callback-token")

	var payload entity.XenditWebhookPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid webhook payload"})
	}

	if err := h.paymentUsecase.HandleWebhook(c.Request().Context(), callbackToken, payload); err != nil {
		switch {
		case errors.Is(err, usecase.ErrUnauthorized):
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid callback token"})
		case errors.Is(err, usecase.ErrNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"message": "payment log not found for invoice"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to process webhook"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "webhook received"})
}

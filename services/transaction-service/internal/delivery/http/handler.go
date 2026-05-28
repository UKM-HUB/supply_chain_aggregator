package http

import (
	"errors"
	"net/http"

	"supply-chain-aggregator/services/transaction-service/internal/entity"
	"supply-chain-aggregator/services/transaction-service/internal/usecase"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	txUsecase *usecase.TransactionUsecase
}

func NewHandler(txUsecase *usecase.TransactionUsecase) *Handler {
	return &Handler{txUsecase: txUsecase}
}

func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "transaction-service",
	})
}

type createTransactionRequest struct {
	UserID        string  `json:"user_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
}

func (h *Handler) CreateTransaction(c echo.Context) error {
	var req createTransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	tx, err := h.txUsecase.Create(c.Request().Context(), usecase.CreateTransactionInput{
		UserID:        req.UserID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidUserID):
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "user_id is required"})
		case errors.Is(err, usecase.ErrInvalidAmount):
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "amount must be greater than zero"})
		case errors.Is(err, usecase.ErrInvalidPayment):
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "payment_method is required"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to create transaction"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": formatTransaction(tx),
	})
}

func (h *Handler) ListTransactions(c echo.Context) error {
	userID := c.QueryParam("user_id")

	transactions, err := h.txUsecase.List(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to list transactions"})
	}

	data := make([]map[string]interface{}, 0, len(transactions))
	for _, tx := range transactions {
		data = append(data, formatTransaction(tx))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  data,
		"total": len(data),
	})
}

func (h *Handler) GetTransaction(c echo.Context) error {
	id := c.Param("id")

	tx, err := h.txUsecase.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "transaction not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to get transaction"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": formatTransaction(tx),
	})
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) UpdateTransactionStatus(c echo.Context) error {
	id := c.Param("id")

	var req updateStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	if err := h.txUsecase.UpdateStatus(c.Request().Context(), id, req.Status); err != nil {
		switch {
		case errors.Is(err, usecase.ErrNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"message": "transaction not found"})
		case errors.Is(err, usecase.ErrInvalidStatus):
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid status; valid values: pending, paid, failed, cancelled"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to update transaction status"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "status updated"})
}

func formatTransaction(tx entity.Transaction) map[string]interface{} {
	return map[string]interface{}{
		"id":             tx.ID,
		"invoice_number": tx.InvoiceNumber,
		"user_id":        tx.UserID,
		"amount":         tx.Amount,
		"status":         tx.Status,
		"payment_method": tx.PaymentMethod,
		"created_at":     tx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"updated_at":     tx.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

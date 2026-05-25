package http

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	// Akan di-inject dengan gRPC Client Order Service
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) Checkout(c echo.Context) error {
	var input struct {
		FactoryID   string `json:"factory_id"`
		ProductCode string `json:"product_code"`
		Quantity    int32  `json:"quantity"`
	}

	// Aplikasi kembali ke respons JSON aman jika input error, tidak akan terminate
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Format data tidak valid. Silakan periksa kembali input Anda.",
		})
	}

	// TODO: Panggil gRPC Client ke Order Service

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Pesanan sedang diproses",
	})
}

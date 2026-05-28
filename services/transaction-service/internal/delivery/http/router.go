package http

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, handler *Handler) {
	e.GET("/health", handler.Health)

	api := e.Group("/api/v1")
	api.POST("/transactions", handler.CreateTransaction)
	api.GET("/transactions", handler.ListTransactions)
	api.GET("/transactions/:id", handler.GetTransaction)
	api.PATCH("/transactions/:id/status", handler.UpdateTransactionStatus)
}

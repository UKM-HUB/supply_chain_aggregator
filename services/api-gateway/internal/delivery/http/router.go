package http

import (
	gatewayMiddleware "supply-chain-aggregator/services/api-gateway/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *Handler, jwtSecret string) {
	e.GET("/health", handler.Health)
	e.GET("/openapi.yaml", handler.OpenAPI)

	api := e.Group("/api/v1")
	api.GET("/health", handler.Health)

	auth := api.Group("/auth")
	auth.POST("/login", handler.Login)
	auth.POST("/register", handler.Register)

	protected := api.Group("")
	protected.Use(gatewayMiddleware.JWT(jwtSecret))

	protected.GET("/users", handler.ListUsers)
	protected.GET("/umkm", handler.ListSMEs)
	protected.GET("/categories", handler.ListCategories)
	protected.GET("/categories/:id", handler.GetCategory)
	protected.GET("/nearby/umkm", handler.FindNearbySMEs)
	protected.POST("/transactions", handler.CreateTransaction)
	protected.GET("/transactions", handler.ListTransactions)
	protected.GET("/transactions/:id", handler.GetTransaction)
	protected.POST("/gateway/create-va", handler.CreateVirtualAccount)
	protected.GET("/reports/daily", handler.DailyReport)
	protected.GET("/reports/monthly", handler.MonthlyReport)
	protected.GET("/reports/export", handler.ExportReport)
}

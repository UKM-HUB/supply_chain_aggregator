package http

import (
	gatewayMiddleware "supply-chain-aggregator/services/api-gateway/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *Handler, jwtSecret string) {
	// ── Health ───────────────────────────────────────────────────────────────
	e.GET("/health", handler.Health)
	e.GET("/api/v1/health", handler.Health)

	// ── Swagger UI ───────────────────────────────────────────────────────────
	// GET /swagger         → HTML halaman Swagger UI
	// GET /openapi/*       → static YAML contract files (misal: /openapi/auth.yaml)
	e.GET("/swagger", handler.SwaggerUI)
	e.Static("/openapi", handler.ContractsPath()) // serve seluruh direktori contracts

	// Fallback untuk /openapi.yaml (backward compat)
	e.GET("/openapi.yaml", handler.OpenAPI)

	// ── Auth (public) ────────────────────────────────────────────────────────
	api := e.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/login", handler.Login)
	auth.POST("/register", handler.Register)

	// ── Protected routes ─────────────────────────────────────────────────────
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

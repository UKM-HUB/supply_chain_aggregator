package http

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, handler *Handler) {
	e.GET("/health", handler.Health)

	api := e.Group("/api/v1")
	api.GET("/reports/daily", handler.DailyReport)
	api.GET("/reports/monthly", handler.MonthlyReport)
	api.GET("/reports/export", handler.ExportReport)
}

package http

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, handler *Handler) {
	e.GET("/health", handler.Health)

	api := e.Group("/api/v1")
	api.GET("/nearby/umkm", handler.SearchNearbySMEs)
}

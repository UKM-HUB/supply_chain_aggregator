package http

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, handler *Handler) {
	e.GET("/health", handler.Health)

	api := e.Group("/api/v1")
	api.POST("/umkm", handler.CreateSME)
	api.GET("/umkm", handler.ListSMEs)
	api.GET("/categories", handler.ListCategories)
}

package http

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, handler *Handler) {
	e.GET("/health", handler.Health)

	api := e.Group("/api/v1/auth")
	api.POST("/register", handler.Register)
	api.POST("/login", handler.Login)
	api.POST("/refresh", handler.Refresh)
	api.POST("/logout", handler.Logout)
}

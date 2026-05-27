package http

import (
    "project/internal/delivery/http/handler"

    "github.com/labstack/echo/v4"
)

func RegisterRoutes(
    e *echo.Echo,
    authHandler *handler.AuthHandler,
    txHandler *handler.TransactionHandler,
) {

    e.POST("/register", authHandler.Register)
    e.POST("/transactions", txHandler.Create)
}
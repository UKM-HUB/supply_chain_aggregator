package main

import (
	deliveryHttp "b2b-aggregator/internal/gateway/delivery/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Mencegah aplikasi crash/terminate jika terjadi panic di layer bawah
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	orderHandler := deliveryHttp.NewOrderHandler()

	api := e.Group("/api/v1")
	api.POST("/checkout", orderHandler.Checkout)

	e.Logger.Fatal(e.Start(":8080"))
}

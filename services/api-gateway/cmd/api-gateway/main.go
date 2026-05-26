package main

import (
	"fmt"

	"supply-chain-aggregator/services/api-gateway/internal/config"
	deliveryHTTP "supply-chain-aggregator/services/api-gateway/internal/delivery/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	handler := deliveryHTTP.NewHandler(cfg.AppName, cfg.Environment, cfg.OpenAPIPath)
	deliveryHTTP.RegisterRoutes(e, handler, cfg.JWTSecret)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)))
}

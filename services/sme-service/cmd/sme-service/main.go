package main

import (
	"fmt"

	"supply-chain-aggregator/services/sme-service/internal/config"
	deliveryHTTP "supply-chain-aggregator/services/sme-service/internal/delivery/http"
	"supply-chain-aggregator/services/sme-service/internal/repository"
	"supply-chain-aggregator/services/sme-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	smeRepo := repository.NewInMemorySMERepository()
	smeUsecase := usecase.NewSMEUsecase(smeRepo)
	handler := deliveryHTTP.NewHandler(smeUsecase)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	deliveryHTTP.RegisterRoutes(e, handler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)))
}

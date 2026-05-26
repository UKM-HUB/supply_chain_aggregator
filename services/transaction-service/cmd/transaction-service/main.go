package main

import (
	"fmt"

	"supply-chain-aggregator/services/transaction-service/internal/config"
	deliveryHTTP "supply-chain-aggregator/services/transaction-service/internal/delivery/http"
	"supply-chain-aggregator/services/transaction-service/internal/repository"
	"supply-chain-aggregator/services/transaction-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	txRepo := repository.NewInMemoryTransactionRepository()
	txUsecase := usecase.NewTransactionUsecase(txRepo)
	handler := deliveryHTTP.NewHandler(txUsecase)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	deliveryHTTP.RegisterRoutes(e, handler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)))
}

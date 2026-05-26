package main

import (
	"fmt"

	"supply-chain-aggregator/services/report-service/internal/config"
	deliveryHTTP "supply-chain-aggregator/services/report-service/internal/delivery/http"
	"supply-chain-aggregator/services/report-service/internal/repository"
	"supply-chain-aggregator/services/report-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	reportRepo := repository.NewInMemoryReportRepository()
	reportUsecase := usecase.NewReportUsecase(reportRepo)
	handler := deliveryHTTP.NewHandler(reportUsecase)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	deliveryHTTP.RegisterRoutes(e, handler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)))
}

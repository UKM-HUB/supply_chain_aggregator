package main

import (
	"fmt"

	"supply-chain-aggregator/services/nearby-service/internal/config"
	deliveryHTTP "supply-chain-aggregator/services/nearby-service/internal/delivery/http"
	"supply-chain-aggregator/services/nearby-service/internal/repository"
	"supply-chain-aggregator/services/nearby-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	locationRepo := repository.NewInMemoryLocationRepository()
	nearbyUsecase := usecase.NewNearbyUsecase(locationRepo)
	handler := deliveryHTTP.NewHandler(nearbyUsecase)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	deliveryHTTP.RegisterRoutes(e, handler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)))
}

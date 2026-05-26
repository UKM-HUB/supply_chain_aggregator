package main

import (
	"fmt"

	"supply-chain-aggregator/services/auth-service/internal/config"
	deliveryHTTP "supply-chain-aggregator/services/auth-service/internal/delivery/http"
	jwtmanager "supply-chain-aggregator/services/auth-service/internal/jwt"
	"supply-chain-aggregator/services/auth-service/internal/repository"
	"supply-chain-aggregator/services/auth-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	userRepo := repository.NewInMemoryUserRepository()
	jwtManager := jwtmanager.NewManager(cfg.JWTSecret)
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtManager)
	handler := deliveryHTTP.NewHandler(authUsecase)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	deliveryHTTP.RegisterRoutes(e, handler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)))
}

package main

import (
	"fmt"
	"log"

	"supply-chain-aggregator/services/payment-service/internal/config"
	deliveryHTTP "supply-chain-aggregator/services/payment-service/internal/delivery/http"
	"supply-chain-aggregator/services/payment-service/internal/rabbitmq"
	"supply-chain-aggregator/services/payment-service/internal/repository"
	"supply-chain-aggregator/services/payment-service/internal/usecase"
	"supply-chain-aggregator/services/payment-service/internal/xendit"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	publisher, err := rabbitmq.NewPublisher(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer publisher.Close()

	xenditClient := xendit.NewClient(cfg.XenditSecretKey)
	paymentRepo := repository.NewInMemoryPaymentRepository()
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo, xenditClient, publisher, cfg.XenditCallbackToken)
	handler := deliveryHTTP.NewHandler(paymentUsecase)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	deliveryHTTP.RegisterRoutes(e, handler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)))
}

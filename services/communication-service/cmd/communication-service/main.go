package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"supply-chain-aggregator/services/communication-service/internal/config"
	"supply-chain-aggregator/services/communication-service/internal/consumer"
	"supply-chain-aggregator/services/communication-service/internal/whatsapp"
	"supply-chain-aggregator/services/communication-service/internal/worker"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	c, err := consumer.NewConsumer(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}
	defer c.Close()

	waClient := whatsapp.NewClient(cfg.WhatsAppAPIURL, cfg.WhatsAppToken)
	w := worker.NewPaymentPaidWorker(c, waClient)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.Println("[main] starting payment.paid worker")
		if err := w.Start(ctx); err != nil {
			log.Printf("[main] worker stopped: %v\n", err)
			cancel()
		}
	}()

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "communication-service",
		})
	})

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)); err != nil && err != http.ErrServerClosed {
			log.Printf("[main] HTTP server stopped: %v\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("[main] shutting down")
	e.Shutdown(context.Background())
}

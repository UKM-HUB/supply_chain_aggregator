package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"supply-chain-aggregator/services/transaction-service/internal/config"
	deliveryGRPC "supply-chain-aggregator/services/transaction-service/internal/delivery/grpc"
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

	// ── gRPC server (background goroutine) ───────────────────────────────────
	go func() {
		if err := deliveryGRPC.Start(cfg.GRPCPort, txUsecase); err != nil {
			log.Fatalf("[transaction-service] gRPC server failed: %v", err)
		}
	}()

	// ── Graceful shutdown ────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("[transaction-service] shutting down...")
		os.Exit(0)
	}()

	log.Fatal(e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)))
}

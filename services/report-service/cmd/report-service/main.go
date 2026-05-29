package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"supply-chain-aggregator/services/report-service/internal/config"
	"supply-chain-aggregator/services/report-service/internal/cron"
	deliveryHTTP "supply-chain-aggregator/services/report-service/internal/delivery/http"
	"supply-chain-aggregator/services/report-service/internal/repository"
	"supply-chain-aggregator/services/report-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	// ── Core dependencies ────────────────────────────────────────────────────
	reportRepo := repository.NewInMemoryReportRepository()
	reportUsecase := usecase.NewReportUsecase(reportRepo)
	handler := deliveryHTTP.NewHandler(reportUsecase)

	// ── Cron Scheduler ───────────────────────────────────────────────────────
	scheduler := cron.NewScheduler(reportUsecase, cfg.ReportOutputDir)
	scheduler.Register()
	scheduler.Start()
	defer scheduler.Stop()

	// ── HTTP server ──────────────────────────────────────────────────────────
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	deliveryHTTP.RegisterRoutes(e, handler)

	// ── Graceful shutdown ────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("[report-service] shutting down...")
		os.Exit(0)
	}()

	log.Fatal(e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)))
}

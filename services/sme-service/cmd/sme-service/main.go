package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	pkges "supply-chain-aggregator/pkg/elasticsearch"
	pkgredis "supply-chain-aggregator/pkg/redis"
	"supply-chain-aggregator/services/sme-service/internal/cache"
	"supply-chain-aggregator/services/sme-service/internal/config"
	deliveryGRPC "supply-chain-aggregator/services/sme-service/internal/delivery/grpc"
	deliveryHTTP "supply-chain-aggregator/services/sme-service/internal/delivery/http"
	smees "supply-chain-aggregator/services/sme-service/internal/elasticsearch"
	"supply-chain-aggregator/services/sme-service/internal/repository"
	"supply-chain-aggregator/services/sme-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	// ── Redis (optional) ─────────────────────────────────────────────────────
	var smeCache *cache.SMECache
	redisClient, err := pkgredis.New(pkgredis.Config{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	if err != nil {
		log.Printf("[sme-service] WARNING: Redis tidak tersedia (%v) — berjalan tanpa cache", err)
	} else {
		smeCache = cache.NewSMECache(redisClient)
		log.Println("[sme-service] Redis terhubung")
	}
	_ = smeCache

	// ── Elasticsearch (optional) ─────────────────────────────────────────────
	var smeIndexer *smees.SMEIndexer
	esClient, err := pkges.New(pkges.Config{
		Addresses: []string{cfg.ElasticsearchURL},
	})
	if err != nil {
		log.Printf("[sme-service] WARNING: Elasticsearch tidak tersedia (%v) — berjalan tanpa search index", err)
	} else {
		smeIndexer, err = smees.NewSMEIndexer(esClient)
		if err != nil {
			log.Printf("[sme-service] WARNING: ES indexer gagal init (%v)", err)
			smeIndexer = nil
		} else {
			log.Println("[sme-service] Elasticsearch terhubung")
		}
	}
	_ = smeIndexer

	// ── Core dependencies ────────────────────────────────────────────────────
	smeRepo := repository.NewInMemorySMERepository()
	smeUsecase := usecase.NewSMEUsecase(smeRepo)
	handler := deliveryHTTP.NewHandler(smeUsecase)

	// ── HTTP server ──────────────────────────────────────────────────────────
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	deliveryHTTP.RegisterRoutes(e, handler)

	// ── gRPC server (background goroutine) ───────────────────────────────────
	go func() {
		if err := deliveryGRPC.Start(cfg.GRPCPort, smeUsecase); err != nil {
			log.Fatalf("[sme-service] gRPC server error: %v", err)
		}
	}()

	// ── Graceful shutdown ────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("[sme-service] shutting down...")
		if redisClient != nil {
			_ = redisClient.Close()
		}
		os.Exit(0)
	}()

	log.Printf("[sme-service] HTTP :%s | gRPC :%s", cfg.HTTPPort, cfg.GRPCPort)
	log.Fatal(e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)))
}

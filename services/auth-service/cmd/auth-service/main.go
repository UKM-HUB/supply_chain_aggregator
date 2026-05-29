package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	pkgredis "supply-chain-aggregator/pkg/redis"
	"supply-chain-aggregator/services/auth-service/internal/cache"
	"supply-chain-aggregator/services/auth-service/internal/config"
	deliveryGRPC "supply-chain-aggregator/services/auth-service/internal/delivery/grpc"
	deliveryHTTP "supply-chain-aggregator/services/auth-service/internal/delivery/http"
	jwtmanager "supply-chain-aggregator/services/auth-service/internal/jwt"
	"supply-chain-aggregator/services/auth-service/internal/repository"
	"supply-chain-aggregator/services/auth-service/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	// ── Redis (optional — service tetap jalan tanpa Redis) ───────────────────
	var tokenCache *cache.TokenCache
	redisClient, err := pkgredis.New(pkgredis.Config{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	if err != nil {
		log.Printf("[auth-service] WARNING: Redis tidak tersedia (%v) — berjalan tanpa cache", err)
	} else {
		tokenCache = cache.NewTokenCache(redisClient)
		log.Println("[auth-service] Redis terhubung")
	}
	_ = tokenCache

	// ── Core dependencies ────────────────────────────────────────────────────
	userRepo := repository.NewInMemoryUserRepository()
	jwtManager := jwtmanager.NewManager(cfg.JWTSecret)
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtManager)
	handler := deliveryHTTP.NewHandler(authUsecase)

	// ── HTTP server ──────────────────────────────────────────────────────────
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	deliveryHTTP.RegisterRoutes(e, handler)

	// ── gRPC server (background goroutine) ───────────────────────────────────
	go func() {
		if err := deliveryGRPC.Start(cfg.GRPCPort, authUsecase); err != nil {
			log.Fatalf("[auth-service] gRPC server error: %v", err)
		}
	}()

	// ── Graceful shutdown ────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("[auth-service] shutting down...")
		if redisClient != nil {
			_ = redisClient.Close()
		}
		os.Exit(0)
	}()

	log.Printf("[auth-service] HTTP :%s | gRPC :%s", cfg.HTTPPort, cfg.GRPCPort)
	log.Fatal(e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)))
}

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	pkgconfig "supply-chain-aggregator/pkg/config"
	pkgredis "supply-chain-aggregator/pkg/redis"
	"supply-chain-aggregator/services/auth-service/internal/cache"
	"supply-chain-aggregator/services/auth-service/internal/config"
	deliveryGRPC "supply-chain-aggregator/services/auth-service/internal/delivery/grpc"
	deliveryHTTP "supply-chain-aggregator/services/auth-service/internal/delivery/http"
	jwtmanager "supply-chain-aggregator/services/auth-service/internal/jwt"
	"supply-chain-aggregator/services/auth-service/internal/repository"
	"supply-chain-aggregator/services/auth-service/internal/usecase"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	// ── Redis (optional — degraded gracefully if unavailable) ────────────────
	var tokenCache *cache.TokenCache
	redisClient, err := pkgredis.New(pkgredis.Config{
		Host:     pkgconfig.GetEnv("REDIS_HOST", "localhost"),
		Port:     pkgconfig.GetEnv("REDIS_PORT", "6379"),
		Password: pkgconfig.GetEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})
	if err != nil {
		log.Printf("[auth-service] WARNING: Redis not available (%v) — running without cache", err)
	} else {
		tokenCache = cache.NewTokenCache(redisClient)
		log.Println("[auth-service] Redis connected")
	}
	_ = tokenCache // injected into usecase when needed

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

	// ── gRPC server (background goroutine) ───────────────────────────────────
	go func() {
		if err := deliveryGRPC.Start(cfg.GRPCPort, authUsecase); err != nil {
			log.Fatalf("[auth-service] gRPC server failed: %v", err)
		}
	}()

	// ── Graceful shutdown ────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("[auth-service] shutting down...")
		if redisClient != nil {
			redisClient.Close()
		}
		os.Exit(0)
	}()

	log.Fatal(e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)))
}

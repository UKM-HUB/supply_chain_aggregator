package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"supply-chain-aggregator/services/api-gateway/internal/config"
	deliveryHTTP "supply-chain-aggregator/services/api-gateway/internal/delivery/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	// Resolve ContractsPath ke absolute path agar tidak bergantung
	// pada working directory saat service dijalankan.
	contractsPath := cfg.ContractsPath
	if !filepath.IsAbs(contractsPath) {
		// Coba resolve relatif ke executable dulu
		exe, err := os.Executable()
		if err == nil {
			abs := filepath.Join(filepath.Dir(exe), contractsPath)
			if _, statErr := os.Stat(abs); statErr == nil {
				contractsPath = abs
			}
		}
		// Jika masih tidak ditemukan, coba relatif ke CWD
		if _, err := os.Stat(contractsPath); err != nil {
			cwd, _ := os.Getwd()
			abs := filepath.Join(cwd, contractsPath)
			if _, statErr := os.Stat(abs); statErr == nil {
				contractsPath = abs
			}
		}
	}
	log.Printf("[api-gateway] contracts path: %s", contractsPath)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	handler := deliveryHTTP.NewHandler(cfg.AppName, cfg.Environment, cfg.OpenAPIPath, contractsPath)
	deliveryHTTP.RegisterRoutes(e, handler, cfg.JWTSecret)

	log.Printf("[api-gateway] HTTP :%s", cfg.HTTPPort)
	log.Printf("[api-gateway] Swagger UI  → http://localhost:%s/swagger", cfg.HTTPPort)
	log.Printf("[api-gateway] Health      → http://localhost:%s/health", cfg.HTTPPort)
	log.Fatal(e.Start(fmt.Sprintf(":%s", cfg.HTTPPort)))
}

package config

import "os"

type Config struct {
	AppName     string
	Environment string
	HTTPPort    string
	JWTSecret   string
	OpenAPIPath string
}

func Load() Config {
	return Config{
		AppName:     getEnv("APP_NAME", "supply-chain-api-gateway"),
		Environment: getEnv("APP_ENV", "development"),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "development-secret"),
		OpenAPIPath: getEnv("OPENAPI_PATH", "../../contracts/openapi/api-gateway.yaml"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

package config

import "os"

type Config struct {
	AppName     string
	Environment string
	HTTPPort    string
	GRPCPort    string
	JWTSecret   string
}

func Load() Config {
	return Config{
		AppName:     getEnv("APP_NAME", "auth-service"),
		Environment: getEnv("APP_ENV", "development"),
		HTTPPort:    getEnv("HTTP_PORT", "8081"),
		GRPCPort:    getEnv("GRPC_PORT", "50051"),
		JWTSecret:   getEnv("JWT_SECRET", "development-secret"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

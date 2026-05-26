package config

import "os"

type Config struct {
	AppName     string
	Environment string
	HTTPPort    string
}

func Load() Config {
	return Config{
		AppName:     getEnv("APP_NAME", "report-service"),
		Environment: getEnv("APP_ENV", "development"),
		HTTPPort:    getEnv("HTTP_PORT", "8087"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

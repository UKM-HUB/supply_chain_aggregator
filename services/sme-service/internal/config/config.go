package config

import "os"

type Config struct {
	AppName          string
	Environment      string
	HTTPPort         string
	GRPCPort         string
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	ElasticsearchURL string
}

func Load() Config {
	return Config{
		AppName:          pkgconfig.GetEnv("APP_NAME", "sme-service"),
		Environment:      pkgconfig.GetEnv("APP_ENV", "development"),
		HTTPPort:         pkgconfig.GetEnv("HTTP_PORT", "8082"),
		GRPCPort:         pkgconfig.GetEnv("GRPC_PORT", "50052"),
		RedisHost:        pkgconfig.GetEnv("REDIS_HOST", "localhost"),
		RedisPort:        pkgconfig.GetEnv("REDIS_PORT", "6379"),
		RedisPassword:    pkgconfig.GetEnv("REDIS_PASSWORD", ""),
		ElasticsearchURL: pkgconfig.GetEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

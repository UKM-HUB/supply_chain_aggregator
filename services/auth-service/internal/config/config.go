package config

import pkgconfig "supply-chain-aggregator/pkg/config"

type Config struct {
	AppName       string
	Environment   string
	HTTPPort      string
	GRPCPort      string
	JWTSecret     string
	RedisHost     string
	RedisPort     string
	RedisPassword string
}

func Load() Config {
	return Config{
		AppName:       pkgconfig.GetEnv("APP_NAME", "auth-service"),
		Environment:   pkgconfig.GetEnv("APP_ENV", "development"),
		HTTPPort:      pkgconfig.GetEnv("HTTP_PORT", "8081"),
		GRPCPort:      pkgconfig.GetEnv("GRPC_PORT", "50051"),
		JWTSecret:     pkgconfig.GetEnv("JWT_SECRET", "development-secret-key"),
		RedisHost:     pkgconfig.GetEnv("REDIS_HOST", "localhost"),
		RedisPort:     pkgconfig.GetEnv("REDIS_PORT", "6379"),
		RedisPassword: pkgconfig.GetEnv("REDIS_PASSWORD", ""),
	}
}

package config

import pkgconfig "supply-chain-aggregator/pkg/config"

type Config struct {
	AppName     string
	Environment string
	HTTPPort    string
	GRPCPort    string
	JWTSecret   string
}

func Load() Config {
	return Config{
		AppName:     pkgconfig.GetEnv("APP_NAME", "auth-service"),
		Environment: pkgconfig.GetEnv("APP_ENV", "development"),
		HTTPPort:    pkgconfig.GetEnv("HTTP_PORT", "8081"),
		GRPCPort:    pkgconfig.GetEnv("GRPC_PORT", "50051"),
		JWTSecret:   pkgconfig.GetEnv("JWT_SECRET", "development-secret"),
	}
}

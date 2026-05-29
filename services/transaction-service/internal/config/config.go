package config

import pkgconfig "supply-chain-aggregator/pkg/config"

type Config struct {
	AppName     string
	Environment string
	HTTPPort    string
	GRPCPort    string
}

func Load() Config {
	return Config{
		AppName:     pkgconfig.GetEnv("APP_NAME", "transaction-service"),
		Environment: pkgconfig.GetEnv("APP_ENV", "development"),
		HTTPPort:    pkgconfig.GetEnv("HTTP_PORT", "8084"),
		GRPCPort:    pkgconfig.GetEnv("GRPC_PORT", "50054"),
	}
}

package config

import pkgconfig "supply-chain-aggregator/pkg/config"

type Config struct {
	AppName     string
	Environment string
	HTTPPort    string
}

func Load() Config {
	return Config{
		AppName:     pkgconfig.GetEnv("APP_NAME", "report-service"),
		Environment: pkgconfig.GetEnv("APP_ENV", "development"),
		HTTPPort:    pkgconfig.GetEnv("HTTP_PORT", "8087"),
	}
}

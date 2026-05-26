package config

import pkgconfig "supply-chain-aggregator/pkg/config"

type Config struct {
	AppName       string
	Environment   string
	HTTPPort      string
	JWTSecret     string
	OpenAPIPath   string
	ContractsPath string
}

func Load() Config {
	return Config{
		AppName:       pkgconfig.GetEnv("APP_NAME", "supply-chain-api-gateway"),
		Environment:   pkgconfig.GetEnv("APP_ENV", "development"),
		HTTPPort:      pkgconfig.GetEnv("HTTP_PORT", "8080"),
		JWTSecret:     pkgconfig.GetEnv("JWT_SECRET", "development-secret"),
		OpenAPIPath:   pkgconfig.GetEnv("OPENAPI_PATH", "../../contracts/openapi/api-gateway.yaml"),
		ContractsPath: pkgconfig.GetEnv("CONTRACTS_PATH", "../../contracts/openapi"),
	}
}

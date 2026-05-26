package config

import pkgconfig "supply-chain-aggregator/pkg/config"

type Config struct {
	AppName             string
	Environment         string
	HTTPPort            string
	XenditSecretKey     string
	XenditCallbackToken string
	RabbitMQURL         string
}

func Load() Config {
	return Config{
		AppName:             pkgconfig.GetEnv("APP_NAME", "payment-service"),
		Environment:         pkgconfig.GetEnv("APP_ENV", "development"),
		HTTPPort:            pkgconfig.GetEnv("HTTP_PORT", "8085"),
		XenditSecretKey:     pkgconfig.GetEnv("XENDIT_SECRET_KEY", ""),
		XenditCallbackToken: pkgconfig.GetEnv("XENDIT_CALLBACK_TOKEN", ""),
		RabbitMQURL:         pkgconfig.GetEnv("RABBITMQ_URL", ""),
	}
}

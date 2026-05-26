package config

import pkgconfig "supply-chain-aggregator/pkg/config"

type Config struct {
	AppName        string
	Environment    string
	HTTPPort       string
	RabbitMQURL    string
	WhatsAppAPIURL string
	WhatsAppToken  string
}

func Load() Config {
	return Config{
		AppName:        pkgconfig.GetEnv("APP_NAME", "communication-service"),
		Environment:    pkgconfig.GetEnv("APP_ENV", "development"),
		HTTPPort:       pkgconfig.GetEnv("HTTP_PORT", "8086"),
		RabbitMQURL:    pkgconfig.GetEnv("RABBITMQ_URL", ""),
		WhatsAppAPIURL: pkgconfig.GetEnv("WHATSAPP_API_URL", ""),
		WhatsAppToken:  pkgconfig.GetEnv("WHATSAPP_TOKEN", ""),
	}
}

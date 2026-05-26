package config

import "os"

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
		AppName:        getEnv("APP_NAME", "communication-service"),
		Environment:    getEnv("APP_ENV", "development"),
		HTTPPort:       getEnv("HTTP_PORT", "8086"),
		RabbitMQURL:    getEnv("RABBITMQ_URL", ""),
		WhatsAppAPIURL: getEnv("WHATSAPP_API_URL", ""),
		WhatsAppToken:  getEnv("WHATSAPP_TOKEN", ""),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

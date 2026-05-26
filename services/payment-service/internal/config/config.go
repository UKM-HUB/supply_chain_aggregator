package config

import "os"

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
		AppName:             getEnv("APP_NAME", "payment-service"),
		Environment:         getEnv("APP_ENV", "development"),
		HTTPPort:            getEnv("HTTP_PORT", "8085"),
		XenditSecretKey:     getEnv("XENDIT_SECRET_KEY", ""),
		XenditCallbackToken: getEnv("XENDIT_CALLBACK_TOKEN", ""),
		RabbitMQURL:         getEnv("RABBITMQ_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

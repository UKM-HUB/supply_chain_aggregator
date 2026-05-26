package config

import "os"

// GetEnv returns the value of the environment variable named by key.
// If the variable is empty or unset, fallback is returned instead.
func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

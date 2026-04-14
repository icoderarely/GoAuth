package config

import (
	"os"
	"time"
)

type Config struct {
	JWTSecret     string
	Port          string
	TokenTTLHours time.Duration
}

func LoadConfig() Config {
	return Config{
		JWTSecret:     getString("JWT_SECRET", "secret-123"),
		Port:          getString("PORT", "8080"),
		TokenTTLHours: 24 * time.Hour,
	}
}

func getString(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

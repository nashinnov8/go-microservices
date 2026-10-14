package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServiceName     string
	Environment     string
	HTTPPort        string
	PostgresDSN     string
	DBMaxConns      int32
	DBMinConns      int32
	ShutdownTimeout time.Duration
}

func Load() *Config {
	return &Config{
		ServiceName:     getEnv("SERVICE_NAME", "order-service"),
		Environment:     getEnv("APP_ENV", "development"),
		HTTPPort:        getEnv("HTTP_PORT", "8082"),
		PostgresDSN:     getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/product_db?sslmode=disable"),
		DBMaxConns:      int32(getEnvInt("DB_MAX_CONNS", 25)),
		DBMinConns:      int32(getEnvInt("DB_MIN_CONNS", 5)),
		ShutdownTimeout: time.Duration(getEnvInt("SHUTDOWN_TIMEOUT_SECONDS", 10)) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}

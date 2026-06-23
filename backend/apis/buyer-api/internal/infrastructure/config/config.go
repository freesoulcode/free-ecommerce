package config

import "os"

const (
	defaultServiceName = "buyer-api"
	defaultEnv         = "local"
	defaultHTTPAddr    = ":8080"
	defaultLogLevel    = "info"
)

type Config struct {
	ServiceName string
	Env         string
	HTTPAddr    string
	LogLevel    string
}

func Load() Config {
	return Config{
		ServiceName: getEnv("BUYER_API_SERVICE_NAME", defaultServiceName),
		Env:         getEnv("BUYER_API_ENV", defaultEnv),
		HTTPAddr:    getEnv("BUYER_API_HTTP_ADDR", defaultHTTPAddr),
		LogLevel:    getEnv("BUYER_API_LOG_LEVEL", defaultLogLevel),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

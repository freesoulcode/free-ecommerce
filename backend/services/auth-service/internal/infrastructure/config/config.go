package config

import "os"

const (
	defaultServiceName = "auth-service"
	defaultEnv         = "local"
	defaultHTTPAddr    = ":8081"
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
		ServiceName: getEnv("AUTH_SERVICE_NAME", defaultServiceName),
		Env:         getEnv("AUTH_SERVICE_ENV", defaultEnv),
		HTTPAddr:    getEnv("AUTH_SERVICE_HTTP_ADDR", defaultHTTPAddr),
		LogLevel:    getEnv("AUTH_SERVICE_LOG_LEVEL", defaultLogLevel),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

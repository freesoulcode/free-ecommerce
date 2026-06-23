package config

import "os"

const (
	defaultServiceName = "user-service"
	defaultEnv         = "local"
	defaultHTTPAddr    = ":8082"
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
		ServiceName: getEnv("USER_SERVICE_NAME", defaultServiceName),
		Env:         getEnv("USER_SERVICE_ENV", defaultEnv),
		HTTPAddr:    getEnv("USER_SERVICE_HTTP_ADDR", defaultHTTPAddr),
		LogLevel:    getEnv("USER_SERVICE_LOG_LEVEL", defaultLogLevel),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

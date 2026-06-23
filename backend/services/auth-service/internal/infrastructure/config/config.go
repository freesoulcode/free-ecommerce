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
	GRPCAddr    string
	LogLevel    string
	MySQL       MySQLConfig
}

type MySQLConfig struct {
	DSN string
}

func Load() Config {
	return Config{
		ServiceName: getEnv("AUTH_SERVICE_NAME", defaultServiceName),
		Env:         getEnv("AUTH_SERVICE_ENV", defaultEnv),
		HTTPAddr:    getEnv("AUTH_SERVICE_HTTP_ADDR", defaultHTTPAddr),
		GRPCAddr:    getEnv("AUTH_SERVICE_GRPC_ADDR", ":9081"),
		LogLevel:    getEnv("AUTH_SERVICE_LOG_LEVEL", defaultLogLevel),
		MySQL: MySQLConfig{
			DSN: getEnv("AUTH_SERVICE_MYSQL_DSN", "root:password@tcp(127.0.0.1:30306)/auth_service?charset=utf8mb4&parseTime=True&loc=Local"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

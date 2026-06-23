package config

import "os"

const (
	defaultServiceName = "product-service"
	defaultEnv         = "local"
	defaultHTTPAddr    = ":8083"
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
		ServiceName: getEnv("PRODUCT_SERVICE_NAME", defaultServiceName),
		Env:         getEnv("PRODUCT_SERVICE_ENV", defaultEnv),
		HTTPAddr:    getEnv("PRODUCT_SERVICE_HTTP_ADDR", defaultHTTPAddr),
		GRPCAddr:    getEnv("PRODUCT_SERVICE_GRPC_ADDR", ":9083"),
		LogLevel:    getEnv("PRODUCT_SERVICE_LOG_LEVEL", defaultLogLevel),
		MySQL: MySQLConfig{
			DSN: getEnv("PRODUCT_SERVICE_MYSQL_DSN", "root:password@tcp(127.0.0.1:30306)/product_service?charset=utf8mb4&parseTime=True&loc=Local"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

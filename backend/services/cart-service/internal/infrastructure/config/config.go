package config

import (
	"os"
	"strconv"
)

const (
	defaultServiceName = "cart-service"
	defaultEnv         = "local"
	defaultHTTPAddr    = ":8084"
	defaultLogLevel    = "info"
)

type Config struct {
	ServiceName    string
	Env            string
	HTTPAddr       string
	GRPCAddr       string
	LogLevel       string
	MySQL          MySQLConfig
	Snowflake      SnowflakeConfig
	ProductService ProductServiceConfig
}

type MySQLConfig struct {
	DSN string
}

type SnowflakeConfig struct {
	Node int64
}

type ProductServiceConfig struct {
	GRPCAddr string
}

func Load() Config {
	return Config{
		ServiceName: getEnv("CART_SERVICE_NAME", defaultServiceName),
		Env:         getEnv("CART_SERVICE_ENV", defaultEnv),
		HTTPAddr:    getEnv("CART_SERVICE_HTTP_ADDR", defaultHTTPAddr),
		GRPCAddr:    getEnv("CART_SERVICE_GRPC_ADDR", ":9084"),
		LogLevel:    getEnv("CART_SERVICE_LOG_LEVEL", defaultLogLevel),
		MySQL: MySQLConfig{
			DSN: getEnv("CART_SERVICE_MYSQL_DSN", "root:password@tcp(127.0.0.1:30306)/cart_service?charset=utf8mb4&parseTime=True&loc=Local"),
		},
		Snowflake: SnowflakeConfig{
			Node: getEnvInt64("CART_SERVICE_SNOWFLAKE_NODE", 3),
		},
		ProductService: ProductServiceConfig{
			GRPCAddr: getEnv("CART_SERVICE_PRODUCT_SERVICE_GRPC_ADDR", "127.0.0.1:9083"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}

	return parsed
}

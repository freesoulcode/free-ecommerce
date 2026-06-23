package config

import (
	"os"
	"strconv"
)

const (
	defaultServiceName = "payment-service"
	defaultEnv         = "local"
	defaultHTTPAddr    = ":8086"
	defaultLogLevel    = "info"
)

type Config struct {
	ServiceName  string
	Env          string
	HTTPAddr     string
	GRPCAddr     string
	LogLevel     string
	MySQL        MySQLConfig
	Snowflake    SnowflakeConfig
	OrderService OrderServiceConfig
}

type MySQLConfig struct{ DSN string }
type SnowflakeConfig struct{ Node int64 }
type OrderServiceConfig struct{ GRPCAddr string }

func Load() Config {
	return Config{
		ServiceName: getEnv("PAYMENT_SERVICE_NAME", defaultServiceName),
		Env:         getEnv("PAYMENT_SERVICE_ENV", defaultEnv),
		HTTPAddr:    getEnv("PAYMENT_SERVICE_HTTP_ADDR", defaultHTTPAddr),
		GRPCAddr:    getEnv("PAYMENT_SERVICE_GRPC_ADDR", ":9086"),
		LogLevel:    getEnv("PAYMENT_SERVICE_LOG_LEVEL", defaultLogLevel),
		MySQL: MySQLConfig{
			DSN: getEnv("PAYMENT_SERVICE_MYSQL_DSN", "root:password@tcp(127.0.0.1:30306)/payment_service?charset=utf8mb4&parseTime=True&loc=Local"),
		},
		Snowflake:    SnowflakeConfig{Node: getEnvInt64("PAYMENT_SERVICE_SNOWFLAKE_NODE", 5)},
		OrderService: OrderServiceConfig{GRPCAddr: getEnv("PAYMENT_SERVICE_ORDER_SERVICE_GRPC_ADDR", "127.0.0.1:9085")},
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

package config

import (
	"os"
	"strconv"
)

const (
	defaultServiceName = "order-service"
	defaultEnv         = "local"
	defaultHTTPAddr    = ":8085"
	defaultLogLevel    = "info"
)

type Config struct {
	ServiceName string
	Env         string
	HTTPAddr    string
	GRPCAddr    string
	LogLevel    string
	MySQL       MySQLConfig
	Snowflake   SnowflakeConfig
	UserService UserServiceConfig
	CartService CartServiceConfig
}

type MySQLConfig struct{ DSN string }
type SnowflakeConfig struct{ Node int64 }
type UserServiceConfig struct{ GRPCAddr string }
type CartServiceConfig struct{ GRPCAddr string }

func Load() Config {
	return Config{
		ServiceName: getEnv("ORDER_SERVICE_NAME", defaultServiceName),
		Env:         getEnv("ORDER_SERVICE_ENV", defaultEnv),
		HTTPAddr:    getEnv("ORDER_SERVICE_HTTP_ADDR", defaultHTTPAddr),
		GRPCAddr:    getEnv("ORDER_SERVICE_GRPC_ADDR", ":9085"),
		LogLevel:    getEnv("ORDER_SERVICE_LOG_LEVEL", defaultLogLevel),
		MySQL: MySQLConfig{
			DSN: getEnv("ORDER_SERVICE_MYSQL_DSN", "root:password@tcp(127.0.0.1:30306)/order_service?charset=utf8mb4&parseTime=True&loc=Local"),
		},
		Snowflake:   SnowflakeConfig{Node: getEnvInt64("ORDER_SERVICE_SNOWFLAKE_NODE", 4)},
		UserService: UserServiceConfig{GRPCAddr: getEnv("ORDER_SERVICE_USER_SERVICE_GRPC_ADDR", "127.0.0.1:9082")},
		CartService: CartServiceConfig{GRPCAddr: getEnv("ORDER_SERVICE_CART_SERVICE_GRPC_ADDR", "127.0.0.1:9084")},
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

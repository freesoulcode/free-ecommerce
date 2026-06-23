package config

import (
	"os"
	"strconv"
)

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
	GRPCAddr    string
	LogLevel    string
	MySQL       MySQLConfig
	Snowflake   SnowflakeConfig
}

type MySQLConfig struct {
	DSN string
}

type SnowflakeConfig struct {
	Node int64
}

func Load() Config {
	return Config{
		ServiceName: getEnv("USER_SERVICE_NAME", defaultServiceName),
		Env:         getEnv("USER_SERVICE_ENV", defaultEnv),
		HTTPAddr:    getEnv("USER_SERVICE_HTTP_ADDR", defaultHTTPAddr),
		GRPCAddr:    getEnv("USER_SERVICE_GRPC_ADDR", ":9082"),
		LogLevel:    getEnv("USER_SERVICE_LOG_LEVEL", defaultLogLevel),
		MySQL: MySQLConfig{
			DSN: getEnv("USER_SERVICE_MYSQL_DSN", "root:password@tcp(127.0.0.1:30306)/user_service?charset=utf8mb4&parseTime=True&loc=Local"),
		},
		Snowflake: SnowflakeConfig{
			Node: getEnvInt64("USER_SERVICE_SNOWFLAKE_NODE", 1),
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

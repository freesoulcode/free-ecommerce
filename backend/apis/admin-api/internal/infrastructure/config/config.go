package config

import "os"

const (
	defaultServiceName = "admin-api"
	defaultEnv         = "local"
	defaultHTTPAddr    = ":8088"
	defaultLogLevel    = "info"
)

type Config struct {
	ServiceName    string
	Env            string
	HTTPAddr       string
	LogLevel       string
	UserService    UserServiceConfig
	OrderService   OrderServiceConfig
	ProductService ProductServiceConfig
	PaymentService PaymentServiceConfig
}

type UserServiceConfig struct {
	GRPCAddr string
}

type OrderServiceConfig struct {
	GRPCAddr string
}

type ProductServiceConfig struct {
	GRPCAddr string
}

type PaymentServiceConfig struct {
	GRPCAddr string
}

func Load() Config {
	return Config{
		ServiceName: getEnv("ADMIN_API_SERVICE_NAME", defaultServiceName),
		Env:         getEnv("ADMIN_API_ENV", defaultEnv),
		HTTPAddr:    getEnv("ADMIN_API_HTTP_ADDR", defaultHTTPAddr),
		LogLevel:    getEnv("ADMIN_API_LOG_LEVEL", defaultLogLevel),
		UserService: UserServiceConfig{
			GRPCAddr: getEnv("ADMIN_API_USER_SERVICE_GRPC_ADDR", "127.0.0.1:9082"),
		},
		OrderService: OrderServiceConfig{
			GRPCAddr: getEnv("ADMIN_API_ORDER_SERVICE_GRPC_ADDR", "127.0.0.1:9085"),
		},
		ProductService: ProductServiceConfig{
			GRPCAddr: getEnv("ADMIN_API_PRODUCT_SERVICE_GRPC_ADDR", "127.0.0.1:9083"),
		},
		PaymentService: PaymentServiceConfig{
			GRPCAddr: getEnv("ADMIN_API_PAYMENT_SERVICE_GRPC_ADDR", "127.0.0.1:9086"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

package config

import "os"

const (
	defaultServiceName = "buyer-api"
	defaultEnv         = "local"
	defaultHTTPAddr    = ":8080"
	defaultLogLevel    = "info"
)

type Config struct {
	ServiceName    string
	Env            string
	HTTPAddr       string
	LogLevel       string
	UserService    UserServiceConfig
	AuthService    AuthServiceConfig
	ProductService ProductServiceConfig
	CartService    CartServiceConfig
	OrderService   OrderServiceConfig
}

type UserServiceConfig struct {
	GRPCAddr string
}

type AuthServiceConfig struct {
	GRPCAddr string
}

type ProductServiceConfig struct {
	GRPCAddr string
}

type CartServiceConfig struct {
	GRPCAddr string
}

type OrderServiceConfig struct {
	GRPCAddr string
}

func Load() Config {
	return Config{
		ServiceName: getEnv("BUYER_API_SERVICE_NAME", defaultServiceName),
		Env:         getEnv("BUYER_API_ENV", defaultEnv),
		HTTPAddr:    getEnv("BUYER_API_HTTP_ADDR", defaultHTTPAddr),
		LogLevel:    getEnv("BUYER_API_LOG_LEVEL", defaultLogLevel),
		UserService: UserServiceConfig{
			GRPCAddr: getEnv("BUYER_API_USER_SERVICE_GRPC_ADDR", "127.0.0.1:9082"),
		},
		AuthService: AuthServiceConfig{
			GRPCAddr: getEnv("BUYER_API_AUTH_SERVICE_GRPC_ADDR", "127.0.0.1:9081"),
		},
		ProductService: ProductServiceConfig{
			GRPCAddr: getEnv("BUYER_API_PRODUCT_SERVICE_GRPC_ADDR", "127.0.0.1:9083"),
		},
		CartService: CartServiceConfig{
			GRPCAddr: getEnv("BUYER_API_CART_SERVICE_GRPC_ADDR", "127.0.0.1:9084"),
		},
		OrderService: OrderServiceConfig{
			GRPCAddr: getEnv("BUYER_API_ORDER_SERVICE_GRPC_ADDR", "127.0.0.1:9085"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

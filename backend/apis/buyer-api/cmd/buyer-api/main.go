package main

import (
	"log"

	applicationbuyer "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/application/buyer"
	servicehttp "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/handler/http"
	serviceauthgrpc "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/infrastructure/authgrpc"
	servicecartgrpc "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/infrastructure/cartgrpc"
	serviceconfig "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/infrastructure/config"
	serviceordergrpc "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/infrastructure/ordergrpc"
	serviceproductgrpc "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/infrastructure/productgrpc"
	serviceusergrpc "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/infrastructure/usergrpc"
	sharedlogger "github.com/freesoulcode/free-ecommerce/backend/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg := serviceconfig.Load()

	logger, err := sharedlogger.New(cfg.ServiceName, cfg.Env, cfg.LogLevel)
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	userServiceClient, err := serviceusergrpc.New(cfg.UserService.GRPCAddr)
	if err != nil {
		logger.Fatal("init user service grpc client", zap.Error(err))
	}
	defer func() {
		_ = userServiceClient.Close()
	}()

	authServiceClient, err := serviceauthgrpc.New(cfg.AuthService.GRPCAddr)
	if err != nil {
		logger.Fatal("init auth service grpc client", zap.Error(err))
	}
	defer func() {
		_ = authServiceClient.Close()
	}()

	productServiceClient, err := serviceproductgrpc.New(cfg.ProductService.GRPCAddr)
	if err != nil {
		logger.Fatal("init product service grpc client", zap.Error(err))
	}
	defer func() {
		_ = productServiceClient.Close()
	}()

	cartServiceClient, err := servicecartgrpc.New(cfg.CartService.GRPCAddr)
	if err != nil {
		logger.Fatal("init cart service grpc client", zap.Error(err))
	}
	defer func() {
		_ = cartServiceClient.Close()
	}()

	orderServiceClient, err := serviceordergrpc.New(cfg.OrderService.GRPCAddr)
	if err != nil {
		logger.Fatal("init order service grpc client", zap.Error(err))
	}
	defer func() {
		_ = orderServiceClient.Close()
	}()

	registerBuyerService := applicationbuyer.NewRegisterBuyerService(userServiceClient, authServiceClient)
	loginBuyerService := applicationbuyer.NewLoginBuyerService(authServiceClient, userServiceClient)
	addressBuyerService := applicationbuyer.NewAddressBuyerService(userServiceClient)
	buyerHandler := servicehttp.NewBuyerHandler(registerBuyerService, loginBuyerService, addressBuyerService)
	productBrowseService := applicationbuyer.NewProductBrowseService(productServiceClient)
	productHandler := servicehttp.NewProductHandler(productBrowseService)
	cartBuyerService := applicationbuyer.NewCartBuyerService(cartServiceClient)
	cartHandler := servicehttp.NewCartHandler(cartBuyerService)
	orderBuyerService := applicationbuyer.NewOrderBuyerService(orderServiceClient)
	orderHandler := servicehttp.NewOrderHandler(orderBuyerService)

	router := servicehttp.NewRouter(servicehttp.RouterParams{
		ServiceName:    cfg.ServiceName,
		BuyerHandler:   buyerHandler,
		ProductHandler: productHandler,
		CartHandler:    cartHandler,
		OrderHandler:   orderHandler,
	})
	logger.Info("starting http server",
		zap.String("addr", cfg.HTTPAddr),
		zap.String("service", cfg.ServiceName),
		zap.String("user_service_grpc_addr", cfg.UserService.GRPCAddr),
		zap.String("auth_service_grpc_addr", cfg.AuthService.GRPCAddr),
		zap.String("product_service_grpc_addr", cfg.ProductService.GRPCAddr),
		zap.String("cart_service_grpc_addr", cfg.CartService.GRPCAddr),
		zap.String("order_service_grpc_addr", cfg.OrderService.GRPCAddr),
	)

	if err := router.Run(cfg.HTTPAddr); err != nil {
		logger.Fatal("run http server", zap.Error(err))
	}
}

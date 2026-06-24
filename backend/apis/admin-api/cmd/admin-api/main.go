package main

import (
	"log"

	applicationadmin "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/application/admin"
	servicehttp "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/handler/http"
	serviceconfig "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/infrastructure/config"
	serviceordergrpc "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/infrastructure/ordergrpc"
	servicepaymentgrpc "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/infrastructure/paymentgrpc"
	serviceproductgrpc "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/infrastructure/productgrpc"
	serviceusergrpc "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/infrastructure/usergrpc"
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

	orderServiceClient, err := serviceordergrpc.New(cfg.OrderService.GRPCAddr)
	if err != nil {
		logger.Fatal("init order service grpc client", zap.Error(err))
	}
	defer func() {
		_ = orderServiceClient.Close()
	}()

	productServiceClient, err := serviceproductgrpc.New(cfg.ProductService.GRPCAddr)
	if err != nil {
		logger.Fatal("init product service grpc client", zap.Error(err))
	}
	defer func() {
		_ = productServiceClient.Close()
	}()

	paymentServiceClient, err := servicepaymentgrpc.New(cfg.PaymentService.GRPCAddr)
	if err != nil {
		logger.Fatal("init payment service grpc client", zap.Error(err))
	}
	defer func() {
		_ = paymentServiceClient.Close()
	}()

	adminUserService := applicationadmin.NewUserAdminService(userServiceClient)
	adminUserHandler := servicehttp.NewAdminUserHandler(adminUserService)
	adminOrderGroupService := applicationadmin.NewOrderGroupAdminService(orderServiceClient)
	adminOrderGroupHandler := servicehttp.NewAdminOrderGroupHandler(adminOrderGroupService)
	adminShopOrderService := applicationadmin.NewShopOrderAdminService(orderServiceClient)
	adminShopOrderHandler := servicehttp.NewAdminShopOrderHandler(adminShopOrderService)
	adminProductService := applicationadmin.NewProductAdminService(productServiceClient)
	adminProductHandler := servicehttp.NewAdminProductHandler(adminProductService)
	adminPaymentService := applicationadmin.NewPaymentAdminService(paymentServiceClient, orderServiceClient)
	adminPaymentHandler := servicehttp.NewAdminPaymentHandler(adminPaymentService)

	router := servicehttp.NewRouter(servicehttp.RouterParams{
		ServiceName:            cfg.ServiceName,
		AdminUserHandler:       adminUserHandler,
		AdminOrderGroupHandler: adminOrderGroupHandler,
		AdminShopOrderHandler:  adminShopOrderHandler,
		AdminProductHandler:    adminProductHandler,
		AdminPaymentHandler:    adminPaymentHandler,
	})

	logger.Info("starting http server",
		zap.String("addr", cfg.HTTPAddr),
		zap.String("service", cfg.ServiceName),
		zap.String("user_service_grpc_addr", cfg.UserService.GRPCAddr),
		zap.String("order_service_grpc_addr", cfg.OrderService.GRPCAddr),
		zap.String("product_service_grpc_addr", cfg.ProductService.GRPCAddr),
		zap.String("payment_service_grpc_addr", cfg.PaymentService.GRPCAddr),
	)

	if err := router.Run(cfg.HTTPAddr); err != nil {
		logger.Fatal("run http server", zap.Error(err))
	}
}

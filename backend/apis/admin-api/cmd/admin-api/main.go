package main

import (
	"log"

	applicationadmin "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/application/admin"
	servicehttp "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/handler/http"
	serviceconfig "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/infrastructure/config"
	serviceordergrpc "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/infrastructure/ordergrpc"
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

	adminUserService := applicationadmin.NewUserAdminService(userServiceClient)
	adminUserHandler := servicehttp.NewAdminUserHandler(adminUserService)
	adminShopOrderService := applicationadmin.NewShopOrderAdminService(orderServiceClient)
	adminShopOrderHandler := servicehttp.NewAdminShopOrderHandler(adminShopOrderService)

	router := servicehttp.NewRouter(servicehttp.RouterParams{
		ServiceName:           cfg.ServiceName,
		AdminUserHandler:      adminUserHandler,
		AdminShopOrderHandler: adminShopOrderHandler,
	})

	logger.Info("starting http server",
		zap.String("addr", cfg.HTTPAddr),
		zap.String("service", cfg.ServiceName),
		zap.String("user_service_grpc_addr", cfg.UserService.GRPCAddr),
		zap.String("order_service_grpc_addr", cfg.OrderService.GRPCAddr),
	)

	if err := router.Run(cfg.HTTPAddr); err != nil {
		logger.Fatal("run http server", zap.Error(err))
	}
}

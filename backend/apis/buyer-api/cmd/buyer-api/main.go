package main

import (
	"log"

	applicationbuyer "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/application/buyer"
	sharedlogger "github.com/freesoulcode/free-ecommerce/backend/pkg/logger"
	servicehttp "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/handler/http"
	serviceconfig "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/infrastructure/config"
	serviceusergrpc "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/infrastructure/usergrpc"
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

	registerBuyerService := applicationbuyer.NewRegisterBuyerService(userServiceClient)
	buyerHandler := servicehttp.NewBuyerHandler(registerBuyerService)

	router := servicehttp.NewRouter(servicehttp.RouterParams{
		ServiceName: cfg.ServiceName,
		BuyerHandler: buyerHandler,
	})
	logger.Info("starting http server",
		zap.String("addr", cfg.HTTPAddr),
		zap.String("service", cfg.ServiceName),
		zap.String("user_service_grpc_addr", cfg.UserService.GRPCAddr),
	)

	if err := router.Run(cfg.HTTPAddr); err != nil {
		logger.Fatal("run http server", zap.Error(err))
	}
}

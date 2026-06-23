package main

import (
	"fmt"
	"log"
	"net"
	"time"

	sharedlogger "github.com/freesoulcode/free-ecommerce/backend/pkg/logger"
	applicationorder "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/application/order"
	servicegrpc "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/handler/grpc"
	servicehttp "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/handler/http"
	servicecartgrpc "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/infrastructure/cartgrpc"
	serviceconfig "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/infrastructure/config"
	serviceid "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/infrastructure/id"
	servicemysql "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/infrastructure/mysql"
	servicepersistence "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/infrastructure/persistence"
	serviceusergrpc "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/infrastructure/usergrpc"
	orderv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/order/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	cfg := serviceconfig.Load()

	logger, err := sharedlogger.New(cfg.ServiceName, cfg.Env, cfg.LogLevel)
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	if err := servicemysql.Migrate(cfg.MySQL.DSN); err != nil {
		logger.Fatal("run mysql migrations", zap.Error(err))
	}

	db, err := servicemysql.Open(cfg.MySQL.DSN)
	if err != nil {
		logger.Fatal("open mysql", zap.Error(err))
	}

	idGenerator, err := serviceid.NewSnowflakeGenerator(cfg.Snowflake.Node)
	if err != nil {
		logger.Fatal("init snowflake generator", zap.Error(err))
	}

	userClient, err := serviceusergrpc.New(cfg.UserService.GRPCAddr)
	if err != nil {
		logger.Fatal("init user service grpc client", zap.Error(err))
	}
	defer func() { _ = userClient.Close() }()

	cartClient, err := servicecartgrpc.New(cfg.CartService.GRPCAddr)
	if err != nil {
		logger.Fatal("init cart service grpc client", zap.Error(err))
	}
	defer func() { _ = cartClient.Close() }()

	repo := servicepersistence.NewOrderRepository(db)
	paymentTTL := time.Duration(cfg.Payment.TimeoutMinutes) * time.Minute
	submitService := applicationorder.NewSubmitOrderService(repo, idGenerator, userClient, cartClient, paymentTTL, nil)
	listService := applicationorder.NewListBuyerOrderGroupsService(repo)
	getService := applicationorder.NewGetBuyerOrderGroupDetailService(repo)
	listMerchantService := applicationorder.NewListMerchantShopOrdersService(repo)
	getMerchantService := applicationorder.NewGetMerchantShopOrderDetailService(repo)
	markProcessingSvc := applicationorder.NewMarkMerchantShopOrderProcessingService(repo, nil)
	markShippedSvc := applicationorder.NewMarkMerchantShopOrderShippedService(repo, nil)
	receiveShopOrderSvc := applicationorder.NewMarkBuyerShopOrderReceivedService(repo, nil)
	getPaymentInfoService := applicationorder.NewGetOrderGroupPaymentInfoService(repo, nil)
	markPaidService := applicationorder.NewMarkOrderGroupPaidService(repo, nil)
	closeTimeoutService := applicationorder.NewCloseOrderGroupByPaymentTimeoutService(repo, nil)
	orderGRPCServer := servicegrpc.NewOrderServiceServer(submitService, listService, getService, listMerchantService, getMerchantService, markProcessingSvc, markShippedSvc, receiveShopOrderSvc, getPaymentInfoService, markPaidService, closeTimeoutService)

	grpcListener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Fatal("listen grpc", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(grpcServer, orderGRPCServer)

	go func() {
		logger.Info("starting grpc server", zap.String("addr", cfg.GRPCAddr), zap.String("service", cfg.ServiceName))
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Fatal("run grpc server", zap.Error(err))
		}
	}()

	router := servicehttp.NewRouter(servicehttp.RouterParams{})
	logger.Info("starting http server",
		zap.String("addr", cfg.HTTPAddr),
		zap.String("grpc_addr", cfg.GRPCAddr),
		zap.String("service", cfg.ServiceName),
		zap.String("mysql_dsn", maskDSN(cfg.MySQL.DSN)),
		zap.Int64("payment_timeout_minutes", cfg.Payment.TimeoutMinutes),
		zap.String("user_service_grpc_addr", cfg.UserService.GRPCAddr),
		zap.String("cart_service_grpc_addr", cfg.CartService.GRPCAddr),
	)
	if err := router.Run(cfg.HTTPAddr); err != nil {
		logger.Fatal("run http server", zap.Error(err))
	}
}

func maskDSN(dsn string) string {
	for i := 0; i < len(dsn); i++ {
		if dsn[i] == ':' {
			return fmt.Sprintf("%s:***", dsn[:i])
		}
	}
	return "***"
}

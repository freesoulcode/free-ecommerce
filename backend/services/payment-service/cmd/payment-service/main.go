package main

import (
	"fmt"
	"log"
	"net"

	sharedlogger "github.com/freesoulcode/free-ecommerce/backend/pkg/logger"
	applicationpayment "github.com/freesoulcode/free-ecommerce/backend/services/payment-service/internal/application/payment"
	servicegrpc "github.com/freesoulcode/free-ecommerce/backend/services/payment-service/internal/handler/grpc"
	servicehttp "github.com/freesoulcode/free-ecommerce/backend/services/payment-service/internal/handler/http"
	serviceconfig "github.com/freesoulcode/free-ecommerce/backend/services/payment-service/internal/infrastructure/config"
	serviceid "github.com/freesoulcode/free-ecommerce/backend/services/payment-service/internal/infrastructure/id"
	servicemysql "github.com/freesoulcode/free-ecommerce/backend/services/payment-service/internal/infrastructure/mysql"
	serviceordergrpc "github.com/freesoulcode/free-ecommerce/backend/services/payment-service/internal/infrastructure/ordergrpc"
	servicepersistence "github.com/freesoulcode/free-ecommerce/backend/services/payment-service/internal/infrastructure/persistence"
	paymentv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/payment/v1"
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
	orderClient, err := serviceordergrpc.New(cfg.OrderService.GRPCAddr)
	if err != nil {
		logger.Fatal("init order service grpc client", zap.Error(err))
	}
	defer func() { _ = orderClient.Close() }()

	repo := servicepersistence.NewPaymentRepository(db)
	createService := applicationpayment.NewCreatePaymentOrderService(repo, idGenerator, orderClient, nil)
	getService := applicationpayment.NewGetPaymentOrderService(repo, idGenerator, orderClient, nil)
	simulatePayService := applicationpayment.NewSimulatePayService(repo, idGenerator, orderClient, nil)
	listAdminOrdersService := applicationpayment.NewListAdminPaymentOrdersService(repo)
	getAdminOrderService := applicationpayment.NewGetAdminPaymentOrderService(repo)
	paymentGRPCServer := servicegrpc.NewPaymentServiceServer(createService, getService, simulatePayService, listAdminOrdersService, getAdminOrderService)

	grpcListener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Fatal("listen grpc", zap.Error(err))
	}
	grpcServer := grpc.NewServer()
	paymentv1.RegisterPaymentServiceServer(grpcServer, paymentGRPCServer)
	go func() {
		logger.Info("starting grpc server", zap.String("addr", cfg.GRPCAddr), zap.String("service", cfg.ServiceName))
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Fatal("run grpc server", zap.Error(err))
		}
	}()

	router := servicehttp.NewRouter(servicehttp.RouterParams{ServiceName: cfg.ServiceName})
	logger.Info("starting http server", zap.String("addr", cfg.HTTPAddr), zap.String("grpc_addr", cfg.GRPCAddr), zap.String("service", cfg.ServiceName), zap.String("mysql_dsn", maskDSN(cfg.MySQL.DSN)), zap.String("order_service_grpc_addr", cfg.OrderService.GRPCAddr))
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

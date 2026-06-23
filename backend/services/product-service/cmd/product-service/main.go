package main

import (
	"log"
	"net"

	sharedlogger "github.com/freesoulcode/free-ecommerce/backend/pkg/logger"
	applicationproduct "github.com/freesoulcode/free-ecommerce/backend/services/product-service/internal/application/product"
	servicegrpc "github.com/freesoulcode/free-ecommerce/backend/services/product-service/internal/handler/grpc"
	servicehttp "github.com/freesoulcode/free-ecommerce/backend/services/product-service/internal/handler/http"
	serviceconfig "github.com/freesoulcode/free-ecommerce/backend/services/product-service/internal/infrastructure/config"
	servicemysql "github.com/freesoulcode/free-ecommerce/backend/services/product-service/internal/infrastructure/mysql"
	servicepersistence "github.com/freesoulcode/free-ecommerce/backend/services/product-service/internal/infrastructure/persistence"
	productv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/product/v1"
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

	repo := servicepersistence.NewProductRepository(db)
	listService := applicationproduct.NewListPublicProductsService(repo)
	getService := applicationproduct.NewGetPublicProductService(repo)
	productGRPCServer := servicegrpc.NewProductServiceServer(listService, getService)

	grpcListener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Fatal("listen grpc", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	productv1.RegisterProductServiceServer(grpcServer, productGRPCServer)

	go func() {
		logger.Info("starting grpc server", zap.String("addr", cfg.GRPCAddr), zap.String("service", cfg.ServiceName))
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Fatal("run grpc server", zap.Error(err))
		}
	}()

	router := servicehttp.NewRouter(servicehttp.RouterParams{})
	logger.Info("starting http server", zap.String("addr", cfg.HTTPAddr), zap.String("grpc_addr", cfg.GRPCAddr), zap.String("service", cfg.ServiceName))
	if err := router.Run(cfg.HTTPAddr); err != nil {
		logger.Fatal("run http server", zap.Error(err))
	}
}

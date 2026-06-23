package main

import (
	"fmt"
	"log"
	"net"

	sharedlogger "github.com/freesoulcode/free-ecommerce/backend/pkg/logger"
	applicationuser "github.com/freesoulcode/free-ecommerce/backend/services/user-service/internal/application/user"
	servicegrpc "github.com/freesoulcode/free-ecommerce/backend/services/user-service/internal/handler/grpc"
	servicehttp "github.com/freesoulcode/free-ecommerce/backend/services/user-service/internal/handler/http"
	serviceconfig "github.com/freesoulcode/free-ecommerce/backend/services/user-service/internal/infrastructure/config"
	serviceid "github.com/freesoulcode/free-ecommerce/backend/services/user-service/internal/infrastructure/id"
	servicemysql "github.com/freesoulcode/free-ecommerce/backend/services/user-service/internal/infrastructure/mysql"
	servicepersistence "github.com/freesoulcode/free-ecommerce/backend/services/user-service/internal/infrastructure/persistence"
	userv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/user/v1"
	"google.golang.org/grpc"
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

	userRepo := servicepersistence.NewUserRepository(db)
	createUserService := applicationuser.NewCreateUserService(userRepo, idGenerator, nil)
	userGRPCServer := servicegrpc.NewUserServiceServer(createUserService)

	grpcListener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Fatal("listen grpc", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, userGRPCServer)

	go func() {
		logger.Info("starting grpc server",
			zap.String("addr", cfg.GRPCAddr),
			zap.String("service", cfg.ServiceName),
		)
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Fatal("run grpc server", zap.Error(err))
		}
	}()

	router := servicehttp.NewRouter(servicehttp.RouterParams{
		ServiceName: cfg.ServiceName,
	})
	logger.Info("starting http server",
		zap.String("addr", cfg.HTTPAddr),
		zap.String("grpc_addr", cfg.GRPCAddr),
		zap.String("service", cfg.ServiceName),
		zap.String("mysql_dsn", maskDSN(cfg.MySQL.DSN)),
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

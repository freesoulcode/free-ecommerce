package main

import (
	"log"
	"net"

	sharedlogger "github.com/freesoulcode/free-ecommerce/backend/pkg/logger"
	applicationcredential "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/application/credential"
	servicegrpc "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/handler/grpc"
	servicehttp "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/handler/http"
	serviceconfig "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/infrastructure/config"
	servicecrypto "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/infrastructure/crypto"
	serviceid "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/infrastructure/id"
	servicemysql "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/infrastructure/mysql"
	servicepersistence "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/infrastructure/persistence"
	servicetoken "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/infrastructure/token"
	authv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/auth/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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

	credentialRepo := servicepersistence.NewPasswordCredentialRepository(db)
	hasher := servicecrypto.NewArgon2idHasher()
	idGenerator, err := serviceid.NewSnowflakeGenerator(cfg.Snowflake.Node)
	if err != nil {
		logger.Fatal("init snowflake generator", zap.Error(err))
	}

	accessSigner, err := servicetoken.NewRS256Signer(cfg.JWT.RSAPrivateKeyPEM)
	if err != nil {
		logger.Fatal("init jwt signer", zap.Error(err))
	}
	refreshTokenGenerator := servicetoken.NewRandomTokenGenerator()
	createPasswordCredentialService := applicationcredential.NewCreatePasswordCredentialService(credentialRepo, hasher, nil)
	loginService := applicationcredential.NewLoginService(
		credentialRepo,
		credentialRepo,
		hasher,
		idGenerator,
		accessSigner,
		refreshTokenGenerator,
		cfg.JWT.Issuer,
		cfg.JWT.Audience,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		nil,
	)
	refreshTokenService := applicationcredential.NewRefreshTokenService(
		credentialRepo,
		credentialRepo,
		idGenerator,
		accessSigner,
		refreshTokenGenerator,
		cfg.JWT.Issuer,
		cfg.JWT.Audience,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		nil,
	)
	logoutService := applicationcredential.NewLogoutService(credentialRepo, nil)
	authGRPCServer := servicegrpc.NewAuthServiceServer(createPasswordCredentialService, loginService, refreshTokenService, logoutService)

	grpcListener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Fatal("listen grpc", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authGRPCServer)

	go func() {
		logger.Info("starting grpc server",
			zap.String("addr", cfg.GRPCAddr),
			zap.String("service", cfg.ServiceName),
		)
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Fatal("run grpc server", zap.Error(err))
		}
	}()

	router := servicehttp.NewRouter(servicehttp.RouterParams{})
	logger.Info("starting http server",
		zap.String("addr", cfg.HTTPAddr),
		zap.String("grpc_addr", cfg.GRPCAddr),
		zap.String("service", cfg.ServiceName),
	)

	if err := router.Run(cfg.HTTPAddr); err != nil {
		logger.Fatal("run http server", zap.Error(err))
	}
}

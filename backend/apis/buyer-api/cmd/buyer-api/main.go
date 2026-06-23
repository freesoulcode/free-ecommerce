package main

import (
	"log"

	sharedlogger "github.com/freesoulcode/free-ecommerce/backend/pkg/logger"
	servicehttp "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/handler/http"
	serviceconfig "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/infrastructure/config"
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

	router := servicehttp.NewRouter(servicehttp.RouterParams{ServiceName: cfg.ServiceName})
	logger.Info("starting http server",
		zap.String("addr", cfg.HTTPAddr),
		zap.String("service", cfg.ServiceName),
	)

	if err := router.Run(cfg.HTTPAddr); err != nil {
		logger.Fatal("run http server", zap.Error(err))
	}
}

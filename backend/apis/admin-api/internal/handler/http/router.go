package http

import (
	"net/http"
	"time"

	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type RouterParams struct {
	ServiceName           string
	AdminUserHandler      *AdminUserHandler
	AdminShopOrderHandler *AdminShopOrderHandler
}

func NewRouter(params RouterParams) *gin.Engine {
	serviceName := params.ServiceName
	if serviceName == "" {
		serviceName = "admin-api"
	}

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		httpx.OK(c, gin.H{
			"status":    "up",
			"service":   serviceName,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	router.GET("/readyz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	if params.AdminUserHandler != nil {
		params.AdminUserHandler.RegisterRoutes(router)
	}
	if params.AdminShopOrderHandler != nil {
		params.AdminShopOrderHandler.RegisterRoutes(router)
	}

	return router
}

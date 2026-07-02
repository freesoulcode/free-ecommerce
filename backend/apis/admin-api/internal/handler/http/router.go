package http

import (
	"net/http"
	"time"

	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterParams struct {
	ServiceName            string
	AdminUserHandler       *AdminUserHandler
	AdminOrderGroupHandler *AdminOrderGroupHandler
	AdminShopOrderHandler  *AdminShopOrderHandler
	AdminProductHandler    *AdminProductHandler
	AdminPaymentHandler    *AdminPaymentHandler
}

func NewRouter(params RouterParams) *gin.Engine {
	serviceName := params.ServiceName
	if serviceName == "" {
		serviceName = "admin-api"
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5174", "http://127.0.0.1:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
	if params.AdminOrderGroupHandler != nil {
		params.AdminOrderGroupHandler.RegisterRoutes(router)
	}
	if params.AdminShopOrderHandler != nil {
		params.AdminShopOrderHandler.RegisterRoutes(router)
	}
	if params.AdminProductHandler != nil {
		params.AdminProductHandler.RegisterRoutes(router)
	}
	if params.AdminPaymentHandler != nil {
		params.AdminPaymentHandler.RegisterRoutes(router)
	}

	return router
}

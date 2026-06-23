package http

import (
	"net/http"
	"time"

	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type RouterParams struct {
	ServiceName    string
	BuyerHandler   *BuyerHandler
	ProductHandler *ProductHandler
	CartHandler    *CartHandler
	OrderHandler   *OrderHandler
	PaymentHandler *PaymentHandler
}

func NewRouter(params RouterParams) *gin.Engine {
	serviceName := params.ServiceName
	if serviceName == "" {
		serviceName = "buyer-api"
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

	if params.BuyerHandler != nil {
		params.BuyerHandler.RegisterRoutes(router)
	}
	if params.ProductHandler != nil {
		params.ProductHandler.RegisterRoutes(router)
	}
	if params.CartHandler != nil {
		params.CartHandler.RegisterRoutes(router)
	}
	if params.OrderHandler != nil {
		params.OrderHandler.RegisterRoutes(router)
	}
	if params.PaymentHandler != nil {
		params.PaymentHandler.RegisterRoutes(router)
	}

	return router
}

package http

import (
	"net/http"

	applicationbuyer "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/application/buyer"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service *applicationbuyer.PaymentBuyerService
}

type createPaymentOrderRequest struct {
	Channel string `json:"channel"`
}

func NewPaymentHandler(service *applicationbuyer.PaymentBuyerService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/api/v1/buyers/:buyerID/orders/:orderGroupID/payment-order", h.create)
	router.GET("/api/v1/buyers/:buyerID/orders/:orderGroupID/payment-order", h.get)
	router.POST("/api/v1/buyers/:buyerID/orders/:orderGroupID/payment-order/simulate-success", h.simulatePay)
}

func (h *PaymentHandler) create(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}
	orderGroupID, ok := parseInt64PathParam(c, "orderGroupID")
	if !ok {
		return
	}
	var req createPaymentOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		httpx.Error(c, appErrors.InvalidArgument("invalid request body"))
		return
	}
	order, err := h.service.Create(c.Request.Context(), applicationbuyer.CreatePaymentOrderInput{UserID: buyerID, OrderGroupID: orderGroupID, Channel: req.Channel})
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, httpx.Response{Code: "OK", Message: "success", Data: gin.H{"payment_order": paymentOrderResponse(order)}, RequestID: c.GetHeader(httpx.HeaderRequestID), TraceID: c.GetHeader(httpx.HeaderTraceID)})
}

func (h *PaymentHandler) get(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}
	orderGroupID, ok := parseInt64PathParam(c, "orderGroupID")
	if !ok {
		return
	}
	order, err := h.service.Get(c.Request.Context(), buyerID, orderGroupID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.OK(c, gin.H{"payment_order": paymentOrderResponse(order)})
}

func (h *PaymentHandler) simulatePay(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}
	orderGroupID, ok := parseInt64PathParam(c, "orderGroupID")
	if !ok {
		return
	}
	order, err := h.service.SimulatePay(c.Request.Context(), buyerID, orderGroupID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.OK(c, gin.H{"payment_order": paymentOrderResponse(order)})
}

func paymentOrderResponse(order *applicationbuyer.PaymentOrder) gin.H {
	if order == nil {
		return nil
	}
	return gin.H{"id": order.ID, "user_id": order.UserID, "order_group_id": order.OrderGroupID, "status": order.Status, "channel": order.Channel, "pay_amount": order.PayAmount, "currency": order.Currency, "expire_at": order.ExpireAt, "paid_at": order.PaidAt, "created_at": order.CreatedAt, "updated_at": order.UpdatedAt}
}

package http

import (
	"strconv"

	applicationadmin "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/application/admin"
	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type AdminPaymentHandler struct {
	service *applicationadmin.PaymentAdminService
}

func NewAdminPaymentHandler(service *applicationadmin.PaymentAdminService) *AdminPaymentHandler {
	return &AdminPaymentHandler{service: service}
}

func (h *AdminPaymentHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/v1/admin/payments", h.list)
	router.GET("/api/v1/admin/payments/:paymentOrderID", h.detail)
	router.GET("/api/v1/admin/orders/:orderGroupID/payment", h.getByOrderGroup)
	router.POST("/api/v1/admin/orders/:orderGroupID/payments/mock-pay", h.mockPay)
}

func (h *AdminPaymentHandler) list(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 32)
	userID, _ := strconv.ParseInt(c.DefaultQuery("user_id", "0"), 10, 64)
	orderGroupID, _ := strconv.ParseInt(c.DefaultQuery("order_group_id", "0"), 10, 64)

	result, err := h.service.List(c.Request.Context(), applicationadmin.ListPaymentOrdersInput{
		Page:         int32(page),
		PageSize:     int32(pageSize),
		Status:       c.Query("status"),
		UserID:       userID,
		OrderGroupID: orderGroupID,
		Channel:      c.Query("channel"),
	})
	if err != nil {
		httpx.Error(c, err)
		return
	}

	items := make([]gin.H, 0, len(result.PaymentOrders))
	for _, item := range result.PaymentOrders {
		items = append(items, paymentOrderResponse(item))
	}

	httpx.OK(c, gin.H{"payment_orders": items, "total": result.Total, "page": result.Page, "page_size": result.PageSize})
}

func (h *AdminPaymentHandler) detail(c *gin.Context) {
	paymentOrderID, ok := parseInt64PathParam(c, "paymentOrderID")
	if !ok {
		return
	}

	order, err := h.service.Get(c.Request.Context(), paymentOrderID)
	if err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{"payment_order": paymentOrderResponse(order)})
}

func (h *AdminPaymentHandler) getByOrderGroup(c *gin.Context) {
	orderGroupID, ok := parseInt64PathParam(c, "orderGroupID")
	if !ok {
		return
	}

	order, err := h.service.GetByOrderGroup(c.Request.Context(), orderGroupID)
	if err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{"payment_order": paymentOrderResponse(order)})
}

func (h *AdminPaymentHandler) mockPay(c *gin.Context) {
	orderGroupID, ok := parseInt64PathParam(c, "orderGroupID")
	if !ok {
		return
	}

	order, err := h.service.MockPay(c.Request.Context(), orderGroupID)
	if err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{"payment_order": paymentOrderResponse(order)})
}

func paymentOrderResponse(order *applicationadmin.PaymentOrder) gin.H {
	if order == nil {
		return nil
	}
	return gin.H{
		"id":             order.ID,
		"user_id":        order.UserID,
		"order_group_id": order.OrderGroupID,
		"status":         order.Status,
		"channel":        order.Channel,
		"pay_amount":     order.PayAmount,
		"currency":       order.Currency,
		"expire_at":      order.ExpireAt,
		"paid_at":        order.PaidAt,
		"created_at":     order.CreatedAt,
		"updated_at":     order.UpdatedAt,
	}
}

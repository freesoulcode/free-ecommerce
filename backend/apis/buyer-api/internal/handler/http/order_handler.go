package http

import (
	"net/http"
	"strconv"

	applicationbuyer "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/application/buyer"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *applicationbuyer.OrderBuyerService
}

type submitOrderRequest struct {
	AddressID   int64   `json:"address_id" binding:"required"`
	CartItemIDs []int64 `json:"cart_item_ids"`
	Source      string  `json:"source"`
}

func NewOrderHandler(service *applicationbuyer.OrderBuyerService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/api/v1/buyers/:buyerID/orders", h.submit)
	router.GET("/api/v1/buyers/:buyerID/orders", h.list)
	router.GET("/api/v1/buyers/:buyerID/orders/:orderGroupID", h.detail)
}

func (h *OrderHandler) submit(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}
	var req submitOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, appErrors.InvalidArgument("invalid request body"))
		return
	}
	group, err := h.service.Submit(c.Request.Context(), applicationbuyer.SubmitOrderInput{UserID: buyerID, AddressID: req.AddressID, CartItemIDs: req.CartItemIDs, Source: req.Source})
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, httpx.Response{Code: "OK", Message: "success", Data: gin.H{"order_group": orderGroupDetailResponse(group)}, RequestID: c.GetHeader(httpx.HeaderRequestID), TraceID: c.GetHeader(httpx.HeaderTraceID)})
}

func (h *OrderHandler) list(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 32)
	result, err := h.service.List(c.Request.Context(), applicationbuyer.ListOrdersInput{UserID: buyerID, Page: int32(page), PageSize: int32(pageSize), Status: c.Query("status")})
	if err != nil {
		httpx.Error(c, err)
		return
	}
	items := make([]gin.H, 0, len(result.OrderGroups))
	for _, group := range result.OrderGroups {
		items = append(items, orderGroupSummaryResponse(group))
	}
	httpx.OK(c, gin.H{"order_groups": items, "total": result.Total, "page": result.Page, "page_size": result.PageSize})
}

func (h *OrderHandler) detail(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}
	orderGroupID, ok := parseInt64PathParam(c, "orderGroupID")
	if !ok {
		return
	}
	group, err := h.service.Detail(c.Request.Context(), buyerID, orderGroupID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.OK(c, gin.H{"order_group": orderGroupDetailResponse(group)})
}

func orderAddressSnapshotResponse(address *applicationbuyer.OrderAddressSnapshot) gin.H {
	if address == nil {
		return nil
	}
	return gin.H{"id": address.ID, "order_group_id": address.OrderGroupID, "user_id": address.UserID, "receiver_name": address.ReceiverName, "receiver_phone": address.ReceiverPhone, "country_code": address.CountryCode, "province": address.Province, "city": address.City, "district": address.District, "address_line1": address.AddressLine1, "address_line2": address.AddressLine2, "postal_code": address.PostalCode, "tag": address.Tag, "created_at": address.CreatedAt, "updated_at": address.UpdatedAt}
}

func orderItemResponse(item *applicationbuyer.OrderItem) gin.H {
	return gin.H{"id": item.ID, "order_group_id": item.OrderGroupID, "shop_order_id": item.ShopOrderID, "user_id": item.UserID, "shop_id": item.ShopID, "product_id": item.ProductID, "sku_id": item.SKUID, "product_title": item.ProductTitle, "product_sub_title": item.ProductSubTitle, "main_image_url": item.MainImageURL, "sku_name": item.SKUName, "price_amount": item.PriceAmount, "currency": item.Currency, "quantity": item.Quantity, "item_amount": item.ItemAmount, "review_status_snapshot": item.ReviewStatusSnapshot, "product_sale_status_snapshot": item.ProductSaleStatusSnapshot, "sku_sale_status_snapshot": item.SKUSaleStatusSnapshot, "created_at": item.CreatedAt, "updated_at": item.UpdatedAt}
}

func shopOrderSummaryResponse(shopOrder *applicationbuyer.ShopOrderSummary) gin.H {
	return gin.H{"id": shopOrder.ID, "order_group_id": shopOrder.OrderGroupID, "user_id": shopOrder.UserID, "shop_id": shopOrder.ShopID, "shop_name": shopOrder.ShopName, "status": shopOrder.Status, "item_amount": shopOrder.ItemAmount, "shipping_amount": shopOrder.ShippingAmount, "pay_amount": shopOrder.PayAmount, "currency": shopOrder.Currency, "item_count": shopOrder.ItemCount, "paid_at": shopOrder.PaidAt, "created_at": shopOrder.CreatedAt, "updated_at": shopOrder.UpdatedAt}
}

func shopOrderResponse(shopOrder *applicationbuyer.ShopOrder) gin.H {
	items := make([]gin.H, 0, len(shopOrder.Items))
	for _, item := range shopOrder.Items {
		items = append(items, orderItemResponse(item))
	}
	return gin.H{"id": shopOrder.ID, "order_group_id": shopOrder.OrderGroupID, "user_id": shopOrder.UserID, "shop_id": shopOrder.ShopID, "shop_name": shopOrder.ShopName, "status": shopOrder.Status, "item_amount": shopOrder.ItemAmount, "shipping_amount": shopOrder.ShippingAmount, "pay_amount": shopOrder.PayAmount, "currency": shopOrder.Currency, "item_count": shopOrder.ItemCount, "paid_at": shopOrder.PaidAt, "items": items, "created_at": shopOrder.CreatedAt, "updated_at": shopOrder.UpdatedAt}
}

func orderGroupSummaryResponse(group *applicationbuyer.OrderGroupSummary) gin.H {
	shopOrders := make([]gin.H, 0, len(group.ShopOrders))
	for _, shopOrder := range group.ShopOrders {
		shopOrders = append(shopOrders, shopOrderSummaryResponse(shopOrder))
	}
	return gin.H{"id": group.ID, "user_id": group.UserID, "status": group.Status, "source": group.Source, "total_item_amount": group.TotalItemAmount, "total_shipping_amount": group.TotalShippingAmount, "total_pay_amount": group.TotalPayAmount, "currency": group.Currency, "shop_order_count": group.ShopOrderCount, "item_count": group.ItemCount, "payment_deadline_at": group.PaymentDeadlineAt, "paid_at": group.PaidAt, "shop_orders": shopOrders, "created_at": group.CreatedAt, "updated_at": group.UpdatedAt}
}

func orderGroupDetailResponse(group *applicationbuyer.OrderGroupDetail) gin.H {
	shopOrders := make([]gin.H, 0, len(group.ShopOrders))
	for _, shopOrder := range group.ShopOrders {
		shopOrders = append(shopOrders, shopOrderResponse(shopOrder))
	}
	return gin.H{"id": group.ID, "user_id": group.UserID, "status": group.Status, "source": group.Source, "total_item_amount": group.TotalItemAmount, "total_shipping_amount": group.TotalShippingAmount, "total_pay_amount": group.TotalPayAmount, "currency": group.Currency, "shop_order_count": group.ShopOrderCount, "item_count": group.ItemCount, "payment_deadline_at": group.PaymentDeadlineAt, "paid_at": group.PaidAt, "address": orderAddressSnapshotResponse(group.Address), "shop_orders": shopOrders, "created_at": group.CreatedAt, "updated_at": group.UpdatedAt}
}

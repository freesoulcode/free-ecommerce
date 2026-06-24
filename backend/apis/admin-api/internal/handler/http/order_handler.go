package http

import (
	"strconv"

	applicationadmin "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/application/admin"
	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type AdminShopOrderHandler struct {
	service *applicationadmin.ShopOrderAdminService
}

type AdminOrderGroupHandler struct {
	service *applicationadmin.OrderGroupAdminService
}

func NewAdminShopOrderHandler(service *applicationadmin.ShopOrderAdminService) *AdminShopOrderHandler {
	return &AdminShopOrderHandler{service: service}
}

func NewAdminOrderGroupHandler(service *applicationadmin.OrderGroupAdminService) *AdminOrderGroupHandler {
	return &AdminOrderGroupHandler{service: service}
}

func (h *AdminOrderGroupHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/v1/admin/orders", h.list)
	router.GET("/api/v1/admin/orders/:orderGroupID", h.detail)
}

func (h *AdminShopOrderHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/v1/admin/shops/:shopID/orders", h.list)
	router.GET("/api/v1/admin/shops/:shopID/orders/:shopOrderID", h.detail)
	router.POST("/api/v1/admin/shops/:shopID/orders/:shopOrderID/processing", h.markProcessing)
	router.POST("/api/v1/admin/shops/:shopID/orders/:shopOrderID/ship", h.markShipped)
}

func (h *AdminOrderGroupHandler) list(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 32)
	userID, _ := strconv.ParseInt(c.DefaultQuery("user_id", "0"), 10, 64)
	shopID, _ := strconv.ParseInt(c.DefaultQuery("shop_id", "0"), 10, 64)

	result, err := h.service.List(c.Request.Context(), applicationadmin.ListAdminOrderGroupsInput{
		UserID:   userID,
		ShopID:   shopID,
		Page:     int32(page),
		PageSize: int32(pageSize),
		Status:   c.Query("status"),
	})
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

func (h *AdminOrderGroupHandler) detail(c *gin.Context) {
	orderGroupID, ok := parseInt64PathParam(c, "orderGroupID")
	if !ok {
		return
	}

	group, err := h.service.Detail(c.Request.Context(), orderGroupID)
	if err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{"order_group": orderGroupDetailResponse(group)})
}

func (h *AdminShopOrderHandler) list(c *gin.Context) {
	shopID, ok := parseInt64PathParam(c, "shopID")
	if !ok {
		return
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 32)

	result, err := h.service.List(c.Request.Context(), applicationadmin.ListShopOrdersInput{
		ShopID:   shopID,
		Page:     int32(page),
		PageSize: int32(pageSize),
		Status:   c.Query("status"),
	})
	if err != nil {
		httpx.Error(c, err)
		return
	}

	items := make([]gin.H, 0, len(result.ShopOrders))
	for _, shopOrder := range result.ShopOrders {
		items = append(items, merchantShopOrderSummaryResponse(shopOrder))
	}

	httpx.OK(c, gin.H{"shop_orders": items, "total": result.Total, "page": result.Page, "page_size": result.PageSize})
}

func (h *AdminShopOrderHandler) detail(c *gin.Context) {
	shopID, ok := parseInt64PathParam(c, "shopID")
	if !ok {
		return
	}
	shopOrderID, ok := parseInt64PathParam(c, "shopOrderID")
	if !ok {
		return
	}

	shopOrder, err := h.service.Detail(c.Request.Context(), shopID, shopOrderID)
	if err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{"shop_order": merchantShopOrderDetailResponse(shopOrder)})
}

func (h *AdminShopOrderHandler) markProcessing(c *gin.Context) {
	shopID, ok := parseInt64PathParam(c, "shopID")
	if !ok {
		return
	}
	shopOrderID, ok := parseInt64PathParam(c, "shopOrderID")
	if !ok {
		return
	}

	shopOrder, err := h.service.MarkProcessing(c.Request.Context(), shopID, shopOrderID)
	if err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{"shop_order": merchantShopOrderDetailResponse(shopOrder)})
}

func (h *AdminShopOrderHandler) markShipped(c *gin.Context) {
	shopID, ok := parseInt64PathParam(c, "shopID")
	if !ok {
		return
	}
	shopOrderID, ok := parseInt64PathParam(c, "shopOrderID")
	if !ok {
		return
	}

	shopOrder, err := h.service.MarkShipped(c.Request.Context(), shopID, shopOrderID)
	if err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{"shop_order": merchantShopOrderDetailResponse(shopOrder)})
}

func orderAddressSnapshotResponse(address *applicationadmin.OrderAddressSnapshot) gin.H {
	if address == nil {
		return nil
	}

	return gin.H{
		"id":             address.ID,
		"order_group_id": address.OrderGroupID,
		"user_id":        address.UserID,
		"receiver_name":  address.ReceiverName,
		"receiver_phone": address.ReceiverPhone,
		"country_code":   address.CountryCode,
		"province":       address.Province,
		"city":           address.City,
		"district":       address.District,
		"address_line1":  address.AddressLine1,
		"address_line2":  address.AddressLine2,
		"postal_code":    address.PostalCode,
		"tag":            address.Tag,
		"created_at":     address.CreatedAt,
		"updated_at":     address.UpdatedAt,
	}
}

func orderItemResponse(item *applicationadmin.OrderItem) gin.H {
	if item == nil {
		return nil
	}

	return gin.H{
		"id":                           item.ID,
		"order_group_id":               item.OrderGroupID,
		"shop_order_id":                item.ShopOrderID,
		"user_id":                      item.UserID,
		"shop_id":                      item.ShopID,
		"product_id":                   item.ProductID,
		"sku_id":                       item.SKUID,
		"product_title":                item.ProductTitle,
		"product_sub_title":            item.ProductSubTitle,
		"main_image_url":               item.MainImageURL,
		"sku_name":                     item.SKUName,
		"price_amount":                 item.PriceAmount,
		"currency":                     item.Currency,
		"quantity":                     item.Quantity,
		"item_amount":                  item.ItemAmount,
		"review_status_snapshot":       item.ReviewStatusSnapshot,
		"product_sale_status_snapshot": item.ProductSaleStatusSnapshot,
		"sku_sale_status_snapshot":     item.SKUSaleStatusSnapshot,
		"created_at":                   item.CreatedAt,
		"updated_at":                   item.UpdatedAt,
	}
}

func shopOrderResponse(shopOrder *applicationadmin.ShopOrder) gin.H {
	if shopOrder == nil {
		return nil
	}

	items := make([]gin.H, 0, len(shopOrder.Items))
	for _, item := range shopOrder.Items {
		items = append(items, orderItemResponse(item))
	}

	return gin.H{
		"id":              shopOrder.ID,
		"order_group_id":  shopOrder.OrderGroupID,
		"user_id":         shopOrder.UserID,
		"shop_id":         shopOrder.ShopID,
		"shop_name":       shopOrder.ShopName,
		"status":          shopOrder.Status,
		"item_amount":     shopOrder.ItemAmount,
		"shipping_amount": shopOrder.ShippingAmount,
		"pay_amount":      shopOrder.PayAmount,
		"currency":        shopOrder.Currency,
		"item_count":      shopOrder.ItemCount,
		"paid_at":         shopOrder.PaidAt,
		"items":           items,
		"created_at":      shopOrder.CreatedAt,
		"updated_at":      shopOrder.UpdatedAt,
	}
}

func merchantShopOrderSummaryResponse(shopOrder *applicationadmin.MerchantShopOrderSummary) gin.H {
	if shopOrder == nil {
		return nil
	}

	return gin.H{
		"id":                  shopOrder.ID,
		"order_group_id":      shopOrder.OrderGroupID,
		"user_id":             shopOrder.UserID,
		"shop_id":             shopOrder.ShopID,
		"shop_name":           shopOrder.ShopName,
		"status":              shopOrder.Status,
		"item_amount":         shopOrder.ItemAmount,
		"shipping_amount":     shopOrder.ShippingAmount,
		"pay_amount":          shopOrder.PayAmount,
		"currency":            shopOrder.Currency,
		"item_count":          shopOrder.ItemCount,
		"paid_at":             shopOrder.PaidAt,
		"created_at":          shopOrder.CreatedAt,
		"updated_at":          shopOrder.UpdatedAt,
		"order_group_status":  shopOrder.OrderGroupStatus,
		"payment_deadline_at": shopOrder.PaymentDeadlineAt,
	}
}

func merchantShopOrderDetailResponse(shopOrder *applicationadmin.MerchantShopOrderDetail) gin.H {
	if shopOrder == nil {
		return nil
	}

	return gin.H{
		"order_group_id":      shopOrder.OrderGroupID,
		"user_id":             shopOrder.UserID,
		"order_group_status":  shopOrder.OrderGroupStatus,
		"source":              shopOrder.Source,
		"payment_deadline_at": shopOrder.PaymentDeadlineAt,
		"paid_at":             shopOrder.PaidAt,
		"address":             orderAddressSnapshotResponse(shopOrder.Address),
		"shop_order":          shopOrderResponse(shopOrder.ShopOrder),
	}
}

func orderGroupSummaryResponse(group *applicationadmin.OrderGroupSummary) gin.H {
	if group == nil {
		return nil
	}

	shopOrders := make([]gin.H, 0, len(group.ShopOrders))
	for _, shopOrder := range group.ShopOrders {
		shopOrders = append(shopOrders, shopOrderResponse(shopOrder))
	}

	return gin.H{
		"id":                    group.ID,
		"user_id":               group.UserID,
		"status":                group.Status,
		"source":                group.Source,
		"total_item_amount":     group.TotalItemAmount,
		"total_shipping_amount": group.TotalShippingAmount,
		"total_pay_amount":      group.TotalPayAmount,
		"currency":              group.Currency,
		"shop_order_count":      group.ShopOrderCount,
		"item_count":            group.ItemCount,
		"payment_deadline_at":   group.PaymentDeadlineAt,
		"paid_at":               group.PaidAt,
		"shop_orders":           shopOrders,
		"created_at":            group.CreatedAt,
		"updated_at":            group.UpdatedAt,
	}
}

func orderGroupDetailResponse(group *applicationadmin.OrderGroupDetail) gin.H {
	if group == nil {
		return nil
	}

	shopOrders := make([]gin.H, 0, len(group.ShopOrders))
	for _, shopOrder := range group.ShopOrders {
		shopOrders = append(shopOrders, shopOrderResponse(shopOrder))
	}

	return gin.H{
		"id":                    group.ID,
		"user_id":               group.UserID,
		"status":                group.Status,
		"source":                group.Source,
		"total_item_amount":     group.TotalItemAmount,
		"total_shipping_amount": group.TotalShippingAmount,
		"total_pay_amount":      group.TotalPayAmount,
		"currency":              group.Currency,
		"shop_order_count":      group.ShopOrderCount,
		"item_count":            group.ItemCount,
		"payment_deadline_at":   group.PaymentDeadlineAt,
		"paid_at":               group.PaidAt,
		"address":               orderAddressSnapshotResponse(group.Address),
		"shop_orders":           shopOrders,
		"created_at":            group.CreatedAt,
		"updated_at":            group.UpdatedAt,
	}
}

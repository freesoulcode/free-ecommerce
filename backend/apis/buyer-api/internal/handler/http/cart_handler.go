package http

import (
	"net/http"

	applicationbuyer "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/application/buyer"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service *applicationbuyer.CartBuyerService
}

type addCartItemRequest struct {
	SKUID    int64 `json:"sku_id" binding:"required"`
	Quantity int32 `json:"quantity" binding:"required"`
}

type updateCartItemRequest struct {
	Quantity int32 `json:"quantity" binding:"required"`
	Selected bool  `json:"selected"`
}

func NewCartHandler(service *applicationbuyer.CartBuyerService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/v1/buyers/:buyerID/cart/items", h.list)
	router.POST("/api/v1/buyers/:buyerID/cart/items", h.add)
	router.PUT("/api/v1/buyers/:buyerID/cart/items/:itemID", h.update)
	router.DELETE("/api/v1/buyers/:buyerID/cart/items/:itemID", h.delete)
}

func (h *CartHandler) list(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}

	items, err := h.service.List(c.Request.Context(), buyerID)
	if err != nil {
		httpx.Error(c, err)
		return
	}

	result := make([]gin.H, 0, len(items))
	for _, item := range items {
		result = append(result, cartItemResponse(item))
	}

	httpx.OK(c, gin.H{"items": result})
}

func (h *CartHandler) add(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}

	var req addCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, appErrors.InvalidArgument("invalid request body"))
		return
	}

	item, err := h.service.Add(c.Request.Context(), applicationbuyer.AddCartItemInput{UserID: buyerID, SKUID: req.SKUID, Quantity: req.Quantity})
	if err != nil {
		httpx.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, httpx.Response{Code: "OK", Message: "success", Data: gin.H{"item": cartItemResponse(item)}, RequestID: c.GetHeader(httpx.HeaderRequestID), TraceID: c.GetHeader(httpx.HeaderTraceID)})
}

func (h *CartHandler) update(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}
	itemID, ok := parseInt64PathParam(c, "itemID")
	if !ok {
		return
	}

	var req updateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, appErrors.InvalidArgument("invalid request body"))
		return
	}

	item, err := h.service.Update(c.Request.Context(), applicationbuyer.UpdateCartItemInput{ID: itemID, UserID: buyerID, Quantity: req.Quantity, Selected: req.Selected})
	if err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{"item": cartItemResponse(item)})
}

func (h *CartHandler) delete(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}
	itemID, ok := parseInt64PathParam(c, "itemID")
	if !ok {
		return
	}

	if err := h.service.Delete(c.Request.Context(), applicationbuyer.DeleteCartItemInput{ID: itemID, UserID: buyerID}); err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{"deleted": true})
}

func cartItemResponse(item *applicationbuyer.CartItem) gin.H {
	return gin.H{
		"id":                  item.ID,
		"user_id":             item.UserID,
		"sku_id":              item.SKUID,
		"product_id":          item.ProductID,
		"shop_id":             item.ShopID,
		"shop_name":           item.ShopName,
		"product_title":       item.ProductTitle,
		"product_sub_title":   item.ProductSubTitle,
		"main_image_url":      item.MainImageURL,
		"sku_name":            item.SKUName,
		"price_amount":        item.PriceAmount,
		"currency":            item.Currency,
		"stock":               item.Stock,
		"quantity":            item.Quantity,
		"selected":            item.Selected,
		"review_status":       item.ReviewStatus,
		"product_sale_status": item.ProductSaleStatus,
		"sku_sale_status":     item.SKUSaleStatus,
		"available":           item.Available,
		"created_at":          item.CreatedAt,
		"updated_at":          item.UpdatedAt,
	}
}

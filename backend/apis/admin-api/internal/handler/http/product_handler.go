package http

import (
	"strconv"

	applicationadmin "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/application/admin"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type AdminProductHandler struct {
	service *applicationadmin.ProductAdminService
}

func NewAdminProductHandler(service *applicationadmin.ProductAdminService) *AdminProductHandler {
	return &AdminProductHandler{service: service}
}

func (h *AdminProductHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/v1/admin/products", h.list)
	router.GET("/api/v1/admin/products/:productID", h.detail)
	router.POST("/api/v1/admin/products/:productID/approve", h.approve)
	router.POST("/api/v1/admin/products/:productID/reject", h.reject)
}

func (h *AdminProductHandler) list(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 32)
	shopID, _ := strconv.ParseInt(c.DefaultQuery("shop_id", "0"), 10, 64)

	result, err := h.service.List(c.Request.Context(), applicationadmin.ListProductsInput{
		Page:         int32(page),
		PageSize:     int32(pageSize),
		Keyword:      c.Query("keyword"),
		ShopID:       shopID,
		ReviewStatus: applicationadmin.NormalizeProductFilterStatus(c.Query("review_status")),
		SaleStatus:   applicationadmin.NormalizeProductFilterStatus(c.Query("sale_status")),
	})
	if err != nil {
		httpx.Error(c, err)
		return
	}
	items := make([]gin.H, 0, len(result.Products))
	for _, product := range result.Products {
		items = append(items, productSummaryResponse(product))
	}
	httpx.OK(c, gin.H{"products": items, "total": result.Total, "page": result.Page, "page_size": result.PageSize})
}

func (h *AdminProductHandler) detail(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("productID"), 10, 64)
	if err != nil || productID <= 0 {
		httpx.Error(c, appErrors.InvalidArgument("invalid product id"))
		return
	}
	product, err := h.service.Get(c.Request.Context(), productID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.OK(c, gin.H{"product": productDetailResponse(product)})
}

func (h *AdminProductHandler) approve(c *gin.Context) {
	h.review(c, true)
}

func (h *AdminProductHandler) reject(c *gin.Context) {
	h.review(c, false)
}

func (h *AdminProductHandler) review(c *gin.Context, approve bool) {
	productID, err := strconv.ParseInt(c.Param("productID"), 10, 64)
	if err != nil || productID <= 0 {
		httpx.Error(c, appErrors.InvalidArgument("invalid product id"))
		return
	}
	var product *applicationadmin.ProductDetail
	if approve {
		product, err = h.service.Approve(c.Request.Context(), productID)
	} else {
		product, err = h.service.Reject(c.Request.Context(), productID)
	}
	if err != nil {
		httpx.Error(c, err)
		return
	}
	httpx.OK(c, gin.H{"product": productDetailResponse(product)})
}

func productSummaryResponse(product *applicationadmin.ProductSummary) gin.H {
	if product == nil {
		return nil
	}
	return gin.H{
		"id":               product.ID,
		"shop_id":          product.ShopID,
		"shop_name":        product.ShopName,
		"title":            product.Title,
		"sub_title":        product.SubTitle,
		"main_image_url":   product.MainImageURL,
		"min_price_amount": product.MinPriceAmount,
		"max_price_amount": product.MaxPriceAmount,
		"currency":         product.Currency,
		"total_stock":      product.TotalStock,
		"review_status":    product.ReviewStatus,
		"sale_status":      product.SaleStatus,
		"created_at":       product.CreatedAt,
		"updated_at":       product.UpdatedAt,
	}
}

func productDetailResponse(product *applicationadmin.ProductDetail) gin.H {
	if product == nil {
		return nil
	}
	skus := make([]gin.H, 0, len(product.SKUs))
	for _, sku := range product.SKUs {
		skus = append(skus, gin.H{"id": sku.ID, "name": sku.Name, "price_amount": sku.PriceAmount, "currency": sku.Currency, "stock": sku.Stock, "sale_status": sku.SaleStatus})
	}
	return gin.H{
		"id":               product.ID,
		"shop_id":          product.ShopID,
		"shop_name":        product.ShopName,
		"title":            product.Title,
		"sub_title":        product.SubTitle,
		"main_image_url":   product.MainImageURL,
		"description":      product.Description,
		"review_status":    product.ReviewStatus,
		"sale_status":      product.SaleStatus,
		"min_price_amount": product.MinPriceAmount,
		"max_price_amount": product.MaxPriceAmount,
		"currency":         product.Currency,
		"total_stock":      product.TotalStock,
		"skus":             skus,
	}
}

package http

import (
	"strconv"

	applicationbuyer "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/application/buyer"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *applicationbuyer.ProductBrowseService
}

func NewProductHandler(service *applicationbuyer.ProductBrowseService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/v1/products", h.list)
	router.GET("/api/v1/products/:productID", h.detail)
}

func (h *ProductHandler) list(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 32)
	shopID, _ := strconv.ParseInt(c.DefaultQuery("shop_id", "0"), 10, 64)

	result, err := h.service.List(c.Request.Context(), applicationbuyer.ListProductsInput{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Keyword:  c.Query("keyword"),
		ShopID:   shopID,
	})
	if err != nil {
		httpx.Error(c, err)
		return
	}

	items := make([]gin.H, 0, len(result.Products))
	for _, product := range result.Products {
		items = append(items, gin.H{
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
		})
	}

	httpx.OK(c, gin.H{"products": items, "total": result.Total, "page": result.Page, "page_size": result.PageSize})
}

func (h *ProductHandler) detail(c *gin.Context) {
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

	skus := make([]gin.H, 0, len(product.SKUs))
	for _, sku := range product.SKUs {
		skus = append(skus, gin.H{"id": sku.ID, "name": sku.Name, "price_amount": sku.PriceAmount, "currency": sku.Currency, "stock": sku.Stock, "sale_status": sku.SaleStatus})
	}

	httpx.OK(c, gin.H{"product": gin.H{
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
	}})
}

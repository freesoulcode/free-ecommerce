package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domainproduct "github.com/freesoulcode/free-ecommerce/backend/services/product-service/internal/domain/product"
	"gorm.io/gorm"
)

type ProductModel struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	ShopID       int64     `gorm:"column:shop_id"`
	ShopName     string    `gorm:"column:shop_name"`
	Title        string    `gorm:"column:title"`
	SubTitle     string    `gorm:"column:sub_title"`
	MainImageURL string    `gorm:"column:main_image_url"`
	Description  string    `gorm:"column:description"`
	ReviewStatus string    `gorm:"column:review_status"`
	SaleStatus   string    `gorm:"column:sale_status"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (ProductModel) TableName() string { return "products" }

type ProductSKUModel struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	ProductID   int64     `gorm:"column:product_id"`
	Name        string    `gorm:"column:name"`
	PriceAmount int64     `gorm:"column:price_amount"`
	Currency    string    `gorm:"column:currency"`
	Stock       int32     `gorm:"column:stock"`
	SaleStatus  string    `gorm:"column:sale_status"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (ProductSKUModel) TableName() string { return "product_skus" }

type listPublicProductRow struct {
	ID             int64  `gorm:"column:id"`
	ShopID         int64  `gorm:"column:shop_id"`
	ShopName       string `gorm:"column:shop_name"`
	Title          string `gorm:"column:title"`
	SubTitle       string `gorm:"column:sub_title"`
	MainImageURL   string `gorm:"column:main_image_url"`
	MinPriceAmount int64  `gorm:"column:min_price_amount"`
	MaxPriceAmount int64  `gorm:"column:max_price_amount"`
	Currency       string `gorm:"column:currency"`
	TotalStock     int32  `gorm:"column:total_stock"`
}

type ProductRepository struct{ db *gorm.DB }

func NewProductRepository(db *gorm.DB) *ProductRepository { return &ProductRepository{db: db} }

func (r *ProductRepository) ListPublicProducts(ctx context.Context, query domainproduct.ListPublicProductsQuery) ([]*domainproduct.Summary, int64, error) {
	whereClause, args := buildPublicProductWhere(query)

	countSQL := `
SELECT COUNT(1)
FROM (
    SELECT p.id
    FROM products p
    JOIN product_skus s ON s.product_id = p.id AND s.sale_status = ?
    ` + whereClause + `
    GROUP BY p.id
) t`
	countArgs := append([]any{domainproduct.SaleStatusOnSale}, args...)
	var total int64
	if err := r.db.WithContext(ctx).Raw(countSQL, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count public products: %w", err)
	}

	offset := (query.Page - 1) * query.PageSize
	listSQL := `
SELECT
    p.id,
    p.shop_id,
    p.shop_name,
    p.title,
    p.sub_title,
    p.main_image_url,
    MIN(s.price_amount) AS min_price_amount,
    MAX(s.price_amount) AS max_price_amount,
    MIN(s.currency) AS currency,
    COALESCE(SUM(s.stock), 0) AS total_stock
FROM products p
JOIN product_skus s ON s.product_id = p.id AND s.sale_status = ?
` + whereClause + `
GROUP BY p.id, p.shop_id, p.shop_name, p.title, p.sub_title, p.main_image_url
ORDER BY p.updated_at DESC, p.id DESC
LIMIT ? OFFSET ?`
	listArgs := append([]any{domainproduct.SaleStatusOnSale}, args...)
	listArgs = append(listArgs, query.PageSize, offset)

	var rows []listPublicProductRow
	if err := r.db.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list public products: %w", err)
	}

	products := make([]*domainproduct.Summary, 0, len(rows))
	for _, row := range rows {
		products = append(products, &domainproduct.Summary{
			ID:             row.ID,
			ShopID:         row.ShopID,
			ShopName:       row.ShopName,
			Title:          row.Title,
			SubTitle:       row.SubTitle,
			MainImageURL:   row.MainImageURL,
			MinPriceAmount: row.MinPriceAmount,
			MaxPriceAmount: row.MaxPriceAmount,
			Currency:       row.Currency,
			TotalStock:     row.TotalStock,
		})
	}

	return products, total, nil
}

func (r *ProductRepository) GetPublicProduct(ctx context.Context, id int64) (*domainproduct.Detail, error) {
	var model ProductModel
	if err := r.db.WithContext(ctx).
		Where("id = ? AND review_status = ? AND sale_status = ?", id, domainproduct.ReviewStatusApproved, domainproduct.SaleStatusOnSale).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("product not found")
		}
		return nil, fmt.Errorf("find public product: %w", err)
	}

	var skuModels []ProductSKUModel
	if err := r.db.WithContext(ctx).
		Where("product_id = ? AND sale_status = ?", id, domainproduct.SaleStatusOnSale).
		Order("price_amount ASC, id ASC").
		Find(&skuModels).Error; err != nil {
		return nil, fmt.Errorf("list product skus: %w", err)
	}
	if len(skuModels) == 0 {
		return nil, appErrors.NotFound("product not found")
	}

	detail := &domainproduct.Detail{
		ID:           model.ID,
		ShopID:       model.ShopID,
		ShopName:     model.ShopName,
		Title:        model.Title,
		SubTitle:     model.SubTitle,
		MainImageURL: model.MainImageURL,
		Description:  model.Description,
		ReviewStatus: model.ReviewStatus,
		SaleStatus:   model.SaleStatus,
		Currency:     skuModels[0].Currency,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}

	for i, sku := range skuModels {
		detail.SKUs = append(detail.SKUs, &domainproduct.SKU{
			ID:          sku.ID,
			Name:        sku.Name,
			PriceAmount: sku.PriceAmount,
			Currency:    sku.Currency,
			Stock:       sku.Stock,
			SaleStatus:  sku.SaleStatus,
		})
		detail.TotalStock += sku.Stock
		if i == 0 || sku.PriceAmount < detail.MinPriceAmount {
			detail.MinPriceAmount = sku.PriceAmount
		}
		if i == 0 || sku.PriceAmount > detail.MaxPriceAmount {
			detail.MaxPriceAmount = sku.PriceAmount
		}
	}

	return detail, nil
}

func buildPublicProductWhere(query domainproduct.ListPublicProductsQuery) (string, []any) {
	clauses := []string{"WHERE p.review_status = ?", "AND p.sale_status = ?"}
	args := []any{domainproduct.ReviewStatusApproved, domainproduct.SaleStatusOnSale}
	if query.ShopID > 0 {
		clauses = append(clauses, "AND p.shop_id = ?")
		args = append(args, query.ShopID)
	}
	if keyword := strings.TrimSpace(query.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		clauses = append(clauses, "AND (p.title LIKE ? OR p.sub_title LIKE ? OR p.shop_name LIKE ?)")
		args = append(args, like, like, like)
	}

	return strings.Join(clauses, " "), args
}

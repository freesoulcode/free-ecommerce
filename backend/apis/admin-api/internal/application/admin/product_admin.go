package admin

import (
	"context"
	"strings"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
)

const (
	ProductReviewStatusApproved = "approved"
	ProductReviewStatusRejected = "review_rejected"
)

type ProductSummary struct {
	ID             int64
	ShopID         int64
	ShopName       string
	Title          string
	SubTitle       string
	MainImageURL   string
	MinPriceAmount int64
	MaxPriceAmount int64
	Currency       string
	TotalStock     int32
	ReviewStatus   string
	SaleStatus     string
	CreatedAt      int64
	UpdatedAt      int64
}

type ProductSKU struct {
	ID          int64
	Name        string
	PriceAmount int64
	Currency    string
	Stock       int32
	SaleStatus  string
}

type ProductDetail struct {
	ID             int64
	ShopID         int64
	ShopName       string
	Title          string
	SubTitle       string
	MainImageURL   string
	Description    string
	ReviewStatus   string
	SaleStatus     string
	MinPriceAmount int64
	MaxPriceAmount int64
	Currency       string
	TotalStock     int32
	SKUs           []*ProductSKU
}

type ListProductsInput struct {
	Page         int32
	PageSize     int32
	Keyword      string
	ShopID       int64
	ReviewStatus string
	SaleStatus   string
}

type ListProductsResult struct {
	Products []*ProductSummary
	Total    int64
	Page     int32
	PageSize int32
}

type ProductAdminClient interface {
	ListAdminProducts(ctx context.Context, input ListProductsInput) (*ListProductsResult, error)
	GetAdminProduct(ctx context.Context, id int64) (*ProductDetail, error)
	ReviewProduct(ctx context.Context, id int64, reviewStatus string) (*ProductDetail, error)
}

type ProductAdminService struct {
	client ProductAdminClient
}

func NewProductAdminService(client ProductAdminClient) *ProductAdminService {
	return &ProductAdminService{client: client}
}

func (s *ProductAdminService) List(ctx context.Context, input ListProductsInput) (*ListProductsResult, error) {
	if s.client == nil {
		return nil, appErrors.Internal("product service is not configured")
	}
	return s.client.ListAdminProducts(ctx, input)
}

func (s *ProductAdminService) Get(ctx context.Context, id int64) (*ProductDetail, error) {
	if s.client == nil {
		return nil, appErrors.Internal("product service is not configured")
	}
	return s.client.GetAdminProduct(ctx, id)
}

func (s *ProductAdminService) Approve(ctx context.Context, id int64) (*ProductDetail, error) {
	if s.client == nil {
		return nil, appErrors.Internal("product service is not configured")
	}
	return s.client.ReviewProduct(ctx, id, ProductReviewStatusApproved)
}

func (s *ProductAdminService) Reject(ctx context.Context, id int64) (*ProductDetail, error) {
	if s.client == nil {
		return nil, appErrors.Internal("product service is not configured")
	}
	return s.client.ReviewProduct(ctx, id, ProductReviewStatusRejected)
}

func NormalizeProductFilterStatus(value string) string {
	return strings.TrimSpace(value)
}

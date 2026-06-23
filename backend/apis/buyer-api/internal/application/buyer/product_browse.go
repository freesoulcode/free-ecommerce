package buyer

import (
	"context"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
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

type ProductBrowseServiceClient interface {
	ListPublicProducts(ctx context.Context, input ListProductsInput) (*ListProductsResult, error)
	GetPublicProduct(ctx context.Context, id int64) (*ProductDetail, error)
}

type ListProductsInput struct {
	Page     int32
	PageSize int32
	Keyword  string
	ShopID   int64
}

type ListProductsResult struct {
	Products []*ProductSummary
	Total    int64
	Page     int32
	PageSize int32
}

type ProductBrowseService struct {
	client ProductBrowseServiceClient
}

func NewProductBrowseService(client ProductBrowseServiceClient) *ProductBrowseService {
	return &ProductBrowseService{client: client}
}

func (s *ProductBrowseService) List(ctx context.Context, input ListProductsInput) (*ListProductsResult, error) {
	if s.client == nil {
		return nil, appErrors.Internal("product service is not configured")
	}

	return s.client.ListPublicProducts(ctx, input)
}

func (s *ProductBrowseService) Get(ctx context.Context, id int64) (*ProductDetail, error) {
	if s.client == nil {
		return nil, appErrors.Internal("product service is not configured")
	}

	return s.client.GetPublicProduct(ctx, id)
}

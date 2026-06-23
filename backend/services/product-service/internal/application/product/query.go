package product

import (
	"context"
	"strings"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domainproduct "github.com/freesoulcode/free-ecommerce/backend/services/product-service/internal/domain/product"
)

type ListPublicProductsInput struct {
	Page     int32
	PageSize int32
	Keyword  string
	ShopID   int64
}

type ListPublicProductsResult struct {
	Products []*domainproduct.Summary
	Total    int64
	Page     int32
	PageSize int32
}

type ListPublicProductsService struct {
	repo domainproduct.Repository
}

type GetPublicProductService struct {
	repo domainproduct.Repository
}

func NewListPublicProductsService(repo domainproduct.Repository) *ListPublicProductsService {
	return &ListPublicProductsService{repo: repo}
}

func NewGetPublicProductService(repo domainproduct.Repository) *GetPublicProductService {
	return &GetPublicProductService{repo: repo}
}

func (s *ListPublicProductsService) Execute(ctx context.Context, input ListPublicProductsInput) (*ListPublicProductsResult, error) {
	page := input.Page
	if page <= 0 {
		page = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if input.ShopID < 0 {
		return nil, appErrors.InvalidArgument("shop id is invalid")
	}

	products, total, err := s.repo.ListPublicProducts(ctx, domainproduct.ListPublicProductsQuery{
		Page:     page,
		PageSize: pageSize,
		Keyword:  strings.TrimSpace(input.Keyword),
		ShopID:   input.ShopID,
	})
	if err != nil {
		return nil, err
	}

	return &ListPublicProductsResult{Products: products, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *GetPublicProductService) Execute(ctx context.Context, id int64) (*domainproduct.Detail, error) {
	if id <= 0 {
		return nil, appErrors.InvalidArgument("product id is required")
	}

	return s.repo.GetPublicProduct(ctx, id)
}

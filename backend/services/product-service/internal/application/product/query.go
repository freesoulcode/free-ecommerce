package product

import (
	"context"
	"sort"
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

type ListAdminProductsInput struct {
	Page         int32
	PageSize     int32
	Keyword      string
	ShopID       int64
	ReviewStatus string
	SaleStatus   string
}

type ListPublicProductsResult struct {
	Products []*domainproduct.Summary
	Total    int64
	Page     int32
	PageSize int32
}

type ListAdminProductsResult struct {
	Products []*domainproduct.Summary
	Total    int64
	Page     int32
	PageSize int32
}

type ListPublicProductsService struct {
	repo domainproduct.Repository
}

type ListAdminProductsService struct {
	repo domainproduct.Repository
}

type GetPublicProductService struct {
	repo domainproduct.Repository
}

type GetAdminProductService struct {
	repo domainproduct.Repository
}

type BatchGetSkuBriefsService struct {
	repo domainproduct.Repository
}

type ReviewProductService struct {
	repo domainproduct.Repository
}

func NewListPublicProductsService(repo domainproduct.Repository) *ListPublicProductsService {
	return &ListPublicProductsService{repo: repo}
}

func NewListAdminProductsService(repo domainproduct.Repository) *ListAdminProductsService {
	return &ListAdminProductsService{repo: repo}
}

func NewGetPublicProductService(repo domainproduct.Repository) *GetPublicProductService {
	return &GetPublicProductService{repo: repo}
}

func NewGetAdminProductService(repo domainproduct.Repository) *GetAdminProductService {
	return &GetAdminProductService{repo: repo}
}

func NewBatchGetSkuBriefsService(repo domainproduct.Repository) *BatchGetSkuBriefsService {
	return &BatchGetSkuBriefsService{repo: repo}
}

func NewReviewProductService(repo domainproduct.Repository) *ReviewProductService {
	return &ReviewProductService{repo: repo}
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

func (s *ListAdminProductsService) Execute(ctx context.Context, input ListAdminProductsInput) (*ListAdminProductsResult, error) {
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

	products, total, err := s.repo.ListAdminProducts(ctx, domainproduct.ListAdminProductsQuery{
		Page:         page,
		PageSize:     pageSize,
		Keyword:      strings.TrimSpace(input.Keyword),
		ShopID:       input.ShopID,
		ReviewStatus: strings.TrimSpace(input.ReviewStatus),
		SaleStatus:   strings.TrimSpace(input.SaleStatus),
	})
	if err != nil {
		return nil, err
	}

	return &ListAdminProductsResult{Products: products, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *GetPublicProductService) Execute(ctx context.Context, id int64) (*domainproduct.Detail, error) {
	if id <= 0 {
		return nil, appErrors.InvalidArgument("product id is required")
	}

	return s.repo.GetPublicProduct(ctx, id)
}

func (s *GetAdminProductService) Execute(ctx context.Context, id int64) (*domainproduct.Detail, error) {
	if id <= 0 {
		return nil, appErrors.InvalidArgument("product id is required")
	}

	return s.repo.GetAdminProduct(ctx, id)
}

func (s *BatchGetSkuBriefsService) Execute(ctx context.Context, skuIDs []int64) ([]*domainproduct.SkuBrief, error) {
	if len(skuIDs) == 0 {
		return []*domainproduct.SkuBrief{}, nil
	}

	unique := make(map[int64]struct{}, len(skuIDs))
	ordered := make([]int64, 0, len(skuIDs))
	for _, skuID := range skuIDs {
		if skuID <= 0 {
			return nil, appErrors.InvalidArgument("sku id is required")
		}
		if _, exists := unique[skuID]; exists {
			continue
		}
		unique[skuID] = struct{}{}
		ordered = append(ordered, skuID)
	}

	briefs, err := s.repo.BatchGetSkuBriefs(ctx, ordered)
	if err != nil {
		return nil, err
	}

	byID := make(map[int64]*domainproduct.SkuBrief, len(briefs))
	for _, brief := range briefs {
		byID[brief.SKUID] = brief
	}

	result := make([]*domainproduct.SkuBrief, 0, len(briefs))
	for _, skuID := range ordered {
		if brief, ok := byID[skuID]; ok {
			result = append(result, brief)
		}
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].SKUID < result[j].SKUID
	})

	return result, nil
}

func (s *ReviewProductService) Execute(ctx context.Context, id int64, reviewStatus string) (*domainproduct.Detail, error) {
	if id <= 0 {
		return nil, appErrors.InvalidArgument("product id is required")
	}
	reviewStatus = strings.TrimSpace(reviewStatus)
	switch reviewStatus {
	case domainproduct.ReviewStatusApproved, domainproduct.ReviewStatusRejected:
	default:
		return nil, appErrors.InvalidArgument("review status is invalid")
	}

	return s.repo.ReviewProduct(ctx, id, reviewStatus)
}

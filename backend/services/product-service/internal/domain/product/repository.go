package product

import "context"

type ListPublicProductsQuery struct {
	Page     int32
	PageSize int32
	Keyword  string
	ShopID   int64
}

type Repository interface {
	ListPublicProducts(ctx context.Context, query ListPublicProductsQuery) ([]*Summary, int64, error)
	GetPublicProduct(ctx context.Context, id int64) (*Detail, error)
	BatchGetSkuBriefs(ctx context.Context, skuIDs []int64) ([]*SkuBrief, error)
}

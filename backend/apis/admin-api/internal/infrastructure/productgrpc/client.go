package productgrpc

import (
	"context"

	applicationadmin "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/application/admin"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	productv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/product/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn   *grpc.ClientConn
	client productv1.ProductServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, client: productv1.NewProductServiceClient(conn)}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) ListAdminProducts(ctx context.Context, input applicationadmin.ListProductsInput) (*applicationadmin.ListProductsResult, error) {
	resp, err := c.client.ListAdminProducts(ctx, &productv1.ListAdminProductsRequest{
		Page:         input.Page,
		PageSize:     input.PageSize,
		Keyword:      input.Keyword,
		ShopId:       input.ShopID,
		ReviewStatus: input.ReviewStatus,
		SaleStatus:   input.SaleStatus,
	})
	if err != nil {
		return nil, toAppError(err)
	}
	items := make([]*applicationadmin.ProductSummary, 0, len(resp.GetProducts()))
	for _, product := range resp.GetProducts() {
		items = append(items, &applicationadmin.ProductSummary{
			ID:             product.GetId(),
			ShopID:         product.GetShopId(),
			ShopName:       product.GetShopName(),
			Title:          product.GetTitle(),
			SubTitle:       product.GetSubTitle(),
			MainImageURL:   product.GetMainImageUrl(),
			MinPriceAmount: product.GetMinPriceAmount(),
			MaxPriceAmount: product.GetMaxPriceAmount(),
			Currency:       product.GetCurrency(),
			TotalStock:     product.GetTotalStock(),
			ReviewStatus:   product.GetReviewStatus(),
			SaleStatus:     product.GetSaleStatus(),
			CreatedAt:      product.GetCreatedAt(),
			UpdatedAt:      product.GetUpdatedAt(),
		})
	}
	return &applicationadmin.ListProductsResult{Products: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

func (c *Client) GetAdminProduct(ctx context.Context, id int64) (*applicationadmin.ProductDetail, error) {
	resp, err := c.client.GetAdminProduct(ctx, &productv1.GetAdminProductRequest{Id: id})
	if err != nil {
		return nil, toAppError(err)
	}
	return toAppProductDetail(resp.GetProduct()), nil
}

func (c *Client) ReviewProduct(ctx context.Context, id int64, reviewStatus string) (*applicationadmin.ProductDetail, error) {
	resp, err := c.client.ReviewProduct(ctx, &productv1.ReviewProductRequest{Id: id, ReviewStatus: reviewStatus})
	if err != nil {
		return nil, toAppError(err)
	}
	return toAppProductDetail(resp.GetProduct()), nil
}

func toAppProductDetail(product *productv1.ProductDetail) *applicationadmin.ProductDetail {
	if product == nil {
		return nil
	}
	skus := make([]*applicationadmin.ProductSKU, 0, len(product.GetSkus()))
	for _, sku := range product.GetSkus() {
		skus = append(skus, &applicationadmin.ProductSKU{ID: sku.GetId(), Name: sku.GetName(), PriceAmount: sku.GetPriceAmount(), Currency: sku.GetCurrency(), Stock: sku.GetStock(), SaleStatus: sku.GetSaleStatus()})
	}
	return &applicationadmin.ProductDetail{
		ID:             product.GetId(),
		ShopID:         product.GetShopId(),
		ShopName:       product.GetShopName(),
		Title:          product.GetTitle(),
		SubTitle:       product.GetSubTitle(),
		MainImageURL:   product.GetMainImageUrl(),
		Description:    product.GetDescription(),
		ReviewStatus:   product.GetReviewStatus(),
		SaleStatus:     product.GetSaleStatus(),
		MinPriceAmount: product.GetMinPriceAmount(),
		MaxPriceAmount: product.GetMaxPriceAmount(),
		Currency:       product.GetCurrency(),
		TotalStock:     product.GetTotalStock(),
		SKUs:           skus,
	}
}

func toAppError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return appErrors.Internal("call product service failed")
	}
	switch st.Code() {
	case codes.InvalidArgument:
		return appErrors.InvalidArgument(st.Message())
	case codes.NotFound:
		return appErrors.NotFound(st.Message())
	default:
		return appErrors.Internal(st.Message())
	}
}

package productgrpc

import (
	"context"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domaincart "github.com/freesoulcode/free-ecommerce/backend/services/cart-service/internal/domain/cart"
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

func (c *Client) BatchGetSkuBriefs(ctx context.Context, skuIDs []int64) (map[int64]*domaincart.SkuBrief, error) {
	resp, err := c.client.BatchGetSkuBriefs(ctx, &productv1.BatchGetSkuBriefsRequest{SkuIds: skuIDs})
	if err != nil {
		return nil, toAppError(err)
	}

	result := make(map[int64]*domaincart.SkuBrief, len(resp.GetSkus()))
	for _, sku := range resp.GetSkus() {
		result[sku.GetSkuId()] = &domaincart.SkuBrief{
			SKUID:             sku.GetSkuId(),
			ProductID:         sku.GetProductId(),
			ShopID:            sku.GetShopId(),
			ShopName:          sku.GetShopName(),
			ProductTitle:      sku.GetProductTitle(),
			ProductSubTitle:   sku.GetProductSubTitle(),
			MainImageURL:      sku.GetMainImageUrl(),
			SKUName:           sku.GetSkuName(),
			PriceAmount:       sku.GetPriceAmount(),
			Currency:          sku.GetCurrency(),
			Stock:             sku.GetStock(),
			ReviewStatus:      sku.GetReviewStatus(),
			ProductSaleStatus: sku.GetProductSaleStatus(),
			SKUSaleStatus:     sku.GetSkuSaleStatus(),
		}
	}

	return result, nil
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

package cartgrpc

import (
	"context"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	applicationorder "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/application/order"
	cartv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/cart/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn   *grpc.ClientConn
	client cartv1.CartServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, client: cartv1.NewCartServiceClient(conn)}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) ListCartItems(ctx context.Context, userID int64) ([]*applicationorder.CartItem, error) {
	resp, err := c.client.ListCartItems(ctx, &cartv1.ListCartItemsRequest{UserId: userID})
	if err != nil {
		return nil, toAppError(err)
	}
	items := make([]*applicationorder.CartItem, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		items = append(items, &applicationorder.CartItem{
			ID:                item.GetId(),
			UserID:            item.GetUserId(),
			SKUID:             item.GetSkuId(),
			ProductID:         item.GetProductId(),
			ShopID:            item.GetShopId(),
			ShopName:          item.GetShopName(),
			ProductTitle:      item.GetProductTitle(),
			ProductSubTitle:   item.GetProductSubTitle(),
			MainImageURL:      item.GetMainImageUrl(),
			SKUName:           item.GetSkuName(),
			PriceAmount:       item.GetPriceAmount(),
			Currency:          item.GetCurrency(),
			Stock:             item.GetStock(),
			Quantity:          item.GetQuantity(),
			Selected:          item.GetSelected(),
			ReviewStatus:      item.GetReviewStatus(),
			ProductSaleStatus: item.GetProductSaleStatus(),
			SKUSaleStatus:     item.GetSkuSaleStatus(),
			Available:         item.GetAvailable(),
		})
	}
	return items, nil
}

func toAppError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return appErrors.Internal("call cart service failed")
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

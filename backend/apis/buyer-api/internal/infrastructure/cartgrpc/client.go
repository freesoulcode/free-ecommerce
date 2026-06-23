package cartgrpc

import (
	"context"

	applicationbuyer "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/application/buyer"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
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

func (c *Client) AddCartItem(ctx context.Context, input applicationbuyer.AddCartItemInput) (*applicationbuyer.CartItem, error) {
	resp, err := c.client.AddCartItem(ctx, &cartv1.AddCartItemRequest{UserId: input.UserID, SkuId: input.SKUID, Quantity: input.Quantity})
	if err != nil {
		return nil, toAppError(err)
	}
	return toAppItem(resp.GetItem()), nil
}

func (c *Client) UpdateCartItem(ctx context.Context, input applicationbuyer.UpdateCartItemInput) (*applicationbuyer.CartItem, error) {
	resp, err := c.client.UpdateCartItem(ctx, &cartv1.UpdateCartItemRequest{Id: input.ID, UserId: input.UserID, Quantity: input.Quantity, Selected: input.Selected})
	if err != nil {
		return nil, toAppError(err)
	}
	return toAppItem(resp.GetItem()), nil
}

func (c *Client) DeleteCartItem(ctx context.Context, input applicationbuyer.DeleteCartItemInput) error {
	_, err := c.client.DeleteCartItem(ctx, &cartv1.DeleteCartItemRequest{Id: input.ID, UserId: input.UserID})
	if err != nil {
		return toAppError(err)
	}
	return nil
}

func (c *Client) ListCartItems(ctx context.Context, userID int64) ([]*applicationbuyer.CartItem, error) {
	resp, err := c.client.ListCartItems(ctx, &cartv1.ListCartItemsRequest{UserId: userID})
	if err != nil {
		return nil, toAppError(err)
	}

	items := make([]*applicationbuyer.CartItem, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		items = append(items, toAppItem(item))
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

func toAppItem(item *cartv1.CartItem) *applicationbuyer.CartItem {
	if item == nil {
		return nil
	}

	return &applicationbuyer.CartItem{
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
		CreatedAt:         item.GetCreatedAt(),
		UpdatedAt:         item.GetUpdatedAt(),
	}
}

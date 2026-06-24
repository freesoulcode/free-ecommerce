package paymentgrpc

import (
	"context"

	applicationadmin "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/application/admin"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	paymentv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn   *grpc.ClientConn
	client paymentv1.PaymentServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, client: paymentv1.NewPaymentServiceClient(conn)}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) ListAdminPaymentOrders(ctx context.Context, input applicationadmin.ListPaymentOrdersInput) (*applicationadmin.ListPaymentOrdersResult, error) {
	resp, err := c.client.ListAdminPaymentOrders(ctx, &paymentv1.ListAdminPaymentOrdersRequest{
		Page:         input.Page,
		PageSize:     input.PageSize,
		Status:       input.Status,
		UserId:       input.UserID,
		OrderGroupId: input.OrderGroupID,
		Channel:      input.Channel,
	})
	if err != nil {
		return nil, toAppError(err)
	}
	items := make([]*applicationadmin.PaymentOrder, 0, len(resp.GetPaymentOrders()))
	for _, item := range resp.GetPaymentOrders() {
		items = append(items, toPaymentOrder(item))
	}
	return &applicationadmin.ListPaymentOrdersResult{PaymentOrders: items, Total: resp.GetTotal(), Page: resp.GetPage(), PageSize: resp.GetPageSize()}, nil
}

func (c *Client) GetAdminPaymentOrder(ctx context.Context, id int64) (*applicationadmin.PaymentOrder, error) {
	resp, err := c.client.GetAdminPaymentOrder(ctx, &paymentv1.GetAdminPaymentOrderRequest{Id: id})
	if err != nil {
		return nil, toAppError(err)
	}
	return toPaymentOrder(resp.GetPaymentOrder()), nil
}

func (c *Client) GetPaymentOrderByOrderGroup(ctx context.Context, userID, orderGroupID int64) (*applicationadmin.PaymentOrder, error) {
	resp, err := c.client.GetPaymentOrderByOrderGroup(ctx, &paymentv1.GetPaymentOrderByOrderGroupRequest{UserId: userID, OrderGroupId: orderGroupID})
	if err != nil {
		return nil, toAppError(err)
	}
	return toPaymentOrder(resp.GetPaymentOrder()), nil
}

func (c *Client) SimulatePay(ctx context.Context, userID, orderGroupID int64) (*applicationadmin.PaymentOrder, error) {
	resp, err := c.client.SimulatePay(ctx, &paymentv1.SimulatePayRequest{UserId: userID, OrderGroupId: orderGroupID})
	if err != nil {
		return nil, toAppError(err)
	}
	return toPaymentOrder(resp.GetPaymentOrder()), nil
}

func toPaymentOrder(order *paymentv1.PaymentOrder) *applicationadmin.PaymentOrder {
	if order == nil {
		return nil
	}
	return &applicationadmin.PaymentOrder{
		ID:           order.GetId(),
		UserID:       order.GetUserId(),
		OrderGroupID: order.GetOrderGroupId(),
		Status:       order.GetStatus(),
		Channel:      order.GetChannel(),
		PayAmount:    order.GetPayAmount(),
		Currency:     order.GetCurrency(),
		ExpireAt:     order.GetExpireAt(),
		PaidAt:       order.GetPaidAt(),
		CreatedAt:    order.GetCreatedAt(),
		UpdatedAt:    order.GetUpdatedAt(),
	}
}

func toAppError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return appErrors.Internal("call payment service failed")
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return appErrors.InvalidArgument(st.Message())
	case codes.NotFound:
		return appErrors.NotFound(st.Message())
	case codes.Unauthenticated:
		return appErrors.Unauthorized(st.Message())
	default:
		return appErrors.Internal(st.Message())
	}
}

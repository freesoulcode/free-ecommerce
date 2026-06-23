package ordergrpc

import (
	"context"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	applicationpayment "github.com/freesoulcode/free-ecommerce/backend/services/payment-service/internal/application/payment"
	orderv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/order/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn   *grpc.ClientConn
	client orderv1.OrderServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, client: orderv1.NewOrderServiceClient(conn)}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) GetOrderGroupPaymentInfo(ctx context.Context, userID, orderGroupID int64) (*applicationpayment.OrderPaymentInfo, error) {
	resp, err := c.client.GetOrderGroupPaymentInfo(ctx, &orderv1.GetOrderGroupPaymentInfoRequest{UserId: userID, OrderGroupId: orderGroupID})
	if err != nil {
		return nil, toAppError(err)
	}
	return toPaymentInfo(resp.GetPaymentInfo()), nil
}

func (c *Client) MarkOrderGroupPaid(ctx context.Context, userID, orderGroupID int64) (*applicationpayment.OrderPaymentInfo, error) {
	resp, err := c.client.MarkOrderGroupPaid(ctx, &orderv1.MarkOrderGroupPaidRequest{UserId: userID, OrderGroupId: orderGroupID})
	if err != nil {
		return nil, toAppError(err)
	}
	return toPaymentInfo(resp.GetPaymentInfo()), nil
}

func (c *Client) CloseOrderGroupByPaymentTimeout(ctx context.Context, userID, orderGroupID int64) (*applicationpayment.OrderPaymentInfo, error) {
	resp, err := c.client.CloseOrderGroupByPaymentTimeout(ctx, &orderv1.CloseOrderGroupByPaymentTimeoutRequest{UserId: userID, OrderGroupId: orderGroupID})
	if err != nil {
		return nil, toAppError(err)
	}
	return toPaymentInfo(resp.GetPaymentInfo()), nil
}

func toPaymentInfo(info *orderv1.OrderGroupPaymentInfo) *applicationpayment.OrderPaymentInfo {
	if info == nil {
		return nil
	}
	var paidAt *time.Time
	if ts := info.GetPaidAt(); ts > 0 {
		value := time.Unix(ts, 0).UTC()
		paidAt = &value
	}
	return &applicationpayment.OrderPaymentInfo{OrderGroupID: info.GetOrderGroupId(), UserID: info.GetUserId(), Status: info.GetStatus(), TotalPayAmount: info.GetTotalPayAmount(), Currency: info.GetCurrency(), PaymentDeadlineAt: time.Unix(info.GetPaymentDeadlineAt(), 0).UTC(), PaidAt: paidAt}
}

func toAppError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return appErrors.Internal("call order service failed")
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

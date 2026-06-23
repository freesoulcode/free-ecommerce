package usergrpc

import (
	"context"
	"strings"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	applicationorder "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/application/order"
	userv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn   *grpc.ClientConn
	client userv1.UserServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, client: userv1.NewUserServiceClient(conn)}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) GetAddress(ctx context.Context, userID, addressID int64) (*applicationorder.Address, error) {
	resp, err := c.client.GetAddress(ctx, &userv1.GetAddressRequest{Id: addressID, UserId: userID})
	if err != nil {
		return nil, toAppError(err)
	}
	address := resp.GetAddress()
	if address == nil {
		return nil, appErrors.NotFound("address not found")
	}
	return &applicationorder.Address{
		ID:            address.GetId(),
		UserID:        address.GetUserId(),
		ReceiverName:  address.GetReceiverName(),
		ReceiverPhone: address.GetReceiverPhone(),
		CountryCode:   address.GetCountryCode(),
		Province:      address.GetProvince(),
		City:          address.GetCity(),
		District:      address.GetDistrict(),
		AddressLine1:  address.GetAddressLine1(),
		AddressLine2:  address.GetAddressLine2(),
		PostalCode:    address.GetPostalCode(),
		Tag:           address.GetTag(),
	}, nil
}

func toAppError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return appErrors.Internal("call user service failed")
	}
	message := st.Message()
	switch st.Code() {
	case codes.InvalidArgument:
		if strings.Contains(message, "address") || strings.Contains(message, "user") {
			return appErrors.InvalidArgument(message)
		}
		return appErrors.InvalidArgument(message)
	case codes.NotFound:
		return appErrors.NotFound(message)
	case codes.Unauthenticated:
		return appErrors.Unauthorized(message)
	default:
		return appErrors.Internal(message)
	}
}

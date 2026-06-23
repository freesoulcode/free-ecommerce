package usergrpc

import (
	"context"
	"strings"

	applicationbuyer "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/application/buyer"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
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

	return &Client{
		conn:   conn,
		client: userv1.NewUserServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) CreateUser(ctx context.Context, input applicationbuyer.CreateBuyerInput) (*applicationbuyer.Buyer, error) {
	resp, err := c.client.CreateUser(ctx, &userv1.CreateUserRequest{
		Email:    input.Email,
		Phone:    input.Phone,
		Nickname: input.Nickname,
	})
	if err != nil {
		return nil, toAppError(err)
	}

	return &applicationbuyer.Buyer{
		ID:            resp.GetId(),
		Email:         resp.GetEmail(),
		Phone:         resp.GetPhone(),
		Nickname:      resp.GetNickname(),
		Status:        resp.GetStatus(),
		EmailVerified: resp.GetEmailVerified(),
		PhoneVerified: resp.GetPhoneVerified(),
	}, nil
}

func (c *Client) DeleteUser(ctx context.Context, id int64) error {
	_, err := c.client.DeleteUser(ctx, &userv1.DeleteUserRequest{Id: id})
	if err != nil {
		return toAppError(err)
	}

	return nil
}

func (c *Client) GetUser(ctx context.Context, id int64) (*applicationbuyer.Buyer, error) {
	resp, err := c.client.GetUser(ctx, &userv1.GetUserRequest{Id: id})
	if err != nil {
		return nil, toAppError(err)
	}

	return &applicationbuyer.Buyer{
		ID:            resp.GetId(),
		Email:         resp.GetEmail(),
		Phone:         resp.GetPhone(),
		Nickname:      resp.GetNickname(),
		Status:        resp.GetStatus(),
		EmailVerified: resp.GetEmailVerified(),
		PhoneVerified: resp.GetPhoneVerified(),
	}, nil
}

func (c *Client) CreateAddress(ctx context.Context, input applicationbuyer.CreateAddressInput) (*applicationbuyer.Address, error) {
	resp, err := c.client.CreateAddress(ctx, &userv1.CreateAddressRequest{
		UserId:        input.UserID,
		ReceiverName:  input.ReceiverName,
		ReceiverPhone: input.ReceiverPhone,
		CountryCode:   input.CountryCode,
		Province:      input.Province,
		City:          input.City,
		District:      input.District,
		AddressLine1:  input.AddressLine1,
		AddressLine2:  input.AddressLine2,
		PostalCode:    input.PostalCode,
		Tag:           input.Tag,
		IsDefault:     input.IsDefault,
	})
	if err != nil {
		return nil, toAppError(err)
	}

	return toAppAddress(resp.GetAddress()), nil
}

func (c *Client) UpdateAddress(ctx context.Context, input applicationbuyer.UpdateAddressInput) (*applicationbuyer.Address, error) {
	resp, err := c.client.UpdateAddress(ctx, &userv1.UpdateAddressRequest{
		Id:            input.ID,
		UserId:        input.UserID,
		ReceiverName:  input.ReceiverName,
		ReceiverPhone: input.ReceiverPhone,
		CountryCode:   input.CountryCode,
		Province:      input.Province,
		City:          input.City,
		District:      input.District,
		AddressLine1:  input.AddressLine1,
		AddressLine2:  input.AddressLine2,
		PostalCode:    input.PostalCode,
		Tag:           input.Tag,
		IsDefault:     input.IsDefault,
	})
	if err != nil {
		return nil, toAppError(err)
	}

	return toAppAddress(resp.GetAddress()), nil
}

func (c *Client) DeleteAddress(ctx context.Context, input applicationbuyer.DeleteAddressInput) error {
	_, err := c.client.DeleteAddress(ctx, &userv1.DeleteAddressRequest{Id: input.ID, UserId: input.UserID})
	if err != nil {
		return toAppError(err)
	}

	return nil
}

func (c *Client) ListAddresses(ctx context.Context, userID int64) ([]*applicationbuyer.Address, error) {
	resp, err := c.client.ListAddresses(ctx, &userv1.ListAddressesRequest{UserId: userID})
	if err != nil {
		return nil, toAppError(err)
	}

	addresses := make([]*applicationbuyer.Address, 0, len(resp.GetAddresses()))
	for _, address := range resp.GetAddresses() {
		addresses = append(addresses, toAppAddress(address))
	}

	return addresses, nil
}

func toAppError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return appErrors.Internal("call user service failed")
	}

	message := st.Message()
	switch st.Code() {
	case codes.InvalidArgument:
		if strings.Contains(message, "nickname") {
			return appErrors.New(appErrors.Code("BUYER_NICKNAME_REQUIRED"), message, 400)
		}
		return appErrors.InvalidArgument(message)
	case codes.AlreadyExists:
		if strings.Contains(message, "email") {
			return appErrors.New(appErrors.Code("USER_EMAIL_ALREADY_EXISTS"), message, 400)
		}
		if strings.Contains(message, "phone") {
			return appErrors.New(appErrors.Code("USER_PHONE_ALREADY_EXISTS"), message, 400)
		}
		return appErrors.New(appErrors.Code("BUYER_ALREADY_EXISTS"), message, 400)
	case codes.NotFound:
		return appErrors.NotFound(message)
	case codes.Unauthenticated:
		return appErrors.Unauthorized(message)
	default:
		return appErrors.Internal(message)
	}
}

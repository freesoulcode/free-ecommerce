package authgrpc

import (
	"context"
	"strings"

	applicationbuyer "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/application/buyer"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	authv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn   *grpc.ClientConn
	client authv1.AuthServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, client: authv1.NewAuthServiceClient(conn)}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) CreatePasswordCredential(ctx context.Context, input applicationbuyer.CreateBuyerInput, userID int64) error {
	_, err := c.client.CreatePasswordCredential(ctx, &authv1.CreatePasswordCredentialRequest{
		UserId:   userID,
		Email:    input.Email,
		Phone:    input.Phone,
		Password: input.Password,
	})
	if err != nil {
		return toAppError(err)
	}
	return nil
}

func toAppError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return appErrors.Internal("call auth service failed")
	}

	message := st.Message()
	switch st.Code() {
	case codes.InvalidArgument:
		if strings.Contains(message, "password") {
			return appErrors.New(appErrors.Code("AUTH_PASSWORD_INVALID"), message, 400)
		}
		return appErrors.InvalidArgument(message)
	case codes.AlreadyExists:
		if strings.Contains(message, "email") {
			return appErrors.New(appErrors.Code("AUTH_EMAIL_ALREADY_EXISTS"), message, 400)
		}
		if strings.Contains(message, "phone") {
			return appErrors.New(appErrors.Code("AUTH_PHONE_ALREADY_EXISTS"), message, 400)
		}
		return appErrors.New(appErrors.Code("AUTH_CREDENTIAL_ALREADY_EXISTS"), message, 400)
	default:
		return appErrors.Internal(message)
	}
}

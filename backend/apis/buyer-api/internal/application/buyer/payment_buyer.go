package buyer

import (
	"context"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
)

type PaymentOrder struct {
	ID           int64
	UserID       int64
	OrderGroupID int64
	Status       string
	Channel      string
	PayAmount    int64
	Currency     string
	ExpireAt     int64
	PaidAt       int64
	CreatedAt    int64
	UpdatedAt    int64
}

type CreatePaymentOrderInput struct {
	UserID       int64
	OrderGroupID int64
	Channel      string
}

type PaymentServiceClient interface {
	CreatePaymentOrder(ctx context.Context, input CreatePaymentOrderInput) (*PaymentOrder, error)
	GetPaymentOrderByOrderGroup(ctx context.Context, userID, orderGroupID int64) (*PaymentOrder, error)
	SimulatePay(ctx context.Context, userID, orderGroupID int64) (*PaymentOrder, error)
}

type PaymentBuyerService struct{ client PaymentServiceClient }

func NewPaymentBuyerService(client PaymentServiceClient) *PaymentBuyerService {
	return &PaymentBuyerService{client: client}
}

func (s *PaymentBuyerService) Create(ctx context.Context, input CreatePaymentOrderInput) (*PaymentOrder, error) {
	if s.client == nil {
		return nil, appErrors.Internal("payment service is not configured")
	}
	return s.client.CreatePaymentOrder(ctx, input)
}

func (s *PaymentBuyerService) Get(ctx context.Context, userID, orderGroupID int64) (*PaymentOrder, error) {
	if s.client == nil {
		return nil, appErrors.Internal("payment service is not configured")
	}
	return s.client.GetPaymentOrderByOrderGroup(ctx, userID, orderGroupID)
}

func (s *PaymentBuyerService) SimulatePay(ctx context.Context, userID, orderGroupID int64) (*PaymentOrder, error) {
	if s.client == nil {
		return nil, appErrors.Internal("payment service is not configured")
	}
	return s.client.SimulatePay(ctx, userID, orderGroupID)
}

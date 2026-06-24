package admin

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

type ListPaymentOrdersInput struct {
	Page         int32
	PageSize     int32
	Status       string
	UserID       int64
	OrderGroupID int64
	Channel      string
}

type ListPaymentOrdersResult struct {
	PaymentOrders []*PaymentOrder
	Total         int64
	Page          int32
	PageSize      int32
}

type PaymentAdminClient interface {
	ListAdminPaymentOrders(ctx context.Context, input ListPaymentOrdersInput) (*ListPaymentOrdersResult, error)
	GetAdminPaymentOrder(ctx context.Context, id int64) (*PaymentOrder, error)
	GetPaymentOrderByOrderGroup(ctx context.Context, userID, orderGroupID int64) (*PaymentOrder, error)
	SimulatePay(ctx context.Context, userID, orderGroupID int64) (*PaymentOrder, error)
}

type AdminOrderGroupLookup interface {
	GetAdminOrderGroupDetail(ctx context.Context, orderGroupID int64) (*OrderGroupDetail, error)
}

type PaymentAdminService struct {
	client      PaymentAdminClient
	orderLookup AdminOrderGroupLookup
}

func NewPaymentAdminService(client PaymentAdminClient, orderLookup AdminOrderGroupLookup) *PaymentAdminService {
	return &PaymentAdminService{client: client, orderLookup: orderLookup}
}

func (s *PaymentAdminService) List(ctx context.Context, input ListPaymentOrdersInput) (*ListPaymentOrdersResult, error) {
	if s.client == nil {
		return nil, appErrors.Internal("payment service is not configured")
	}
	return s.client.ListAdminPaymentOrders(ctx, input)
}

func (s *PaymentAdminService) Get(ctx context.Context, id int64) (*PaymentOrder, error) {
	if s.client == nil {
		return nil, appErrors.Internal("payment service is not configured")
	}
	return s.client.GetAdminPaymentOrder(ctx, id)
}

func (s *PaymentAdminService) GetByOrderGroup(ctx context.Context, orderGroupID int64) (*PaymentOrder, error) {
	if s.client == nil {
		return nil, appErrors.Internal("payment service is not configured")
	}
	group, err := s.loadOrderGroup(ctx, orderGroupID)
	if err != nil {
		return nil, err
	}
	return s.client.GetPaymentOrderByOrderGroup(ctx, group.UserID, orderGroupID)
}

func (s *PaymentAdminService) MockPay(ctx context.Context, orderGroupID int64) (*PaymentOrder, error) {
	if s.client == nil {
		return nil, appErrors.Internal("payment service is not configured")
	}
	group, err := s.loadOrderGroup(ctx, orderGroupID)
	if err != nil {
		return nil, err
	}
	return s.client.SimulatePay(ctx, group.UserID, orderGroupID)
}

func (s *PaymentAdminService) loadOrderGroup(ctx context.Context, orderGroupID int64) (*OrderGroupDetail, error) {
	if orderGroupID <= 0 {
		return nil, appErrors.InvalidArgument("order group id is required")
	}
	if s.orderLookup == nil {
		return nil, appErrors.Internal("order service is not configured")
	}
	return s.orderLookup.GetAdminOrderGroupDetail(ctx, orderGroupID)
}

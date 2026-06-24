package payment

import (
	"context"
	"time"
)

type ListAdminPaymentOrdersQuery struct {
	Page         int32
	PageSize     int32
	Status       string
	UserID       int64
	OrderGroupID int64
	Channel      string
}

type Repository interface {
	FindByOrderGroup(ctx context.Context, userID, orderGroupID int64) (*Order, error)
	GetByID(ctx context.Context, id int64) (*Order, error)
	ListAdminPaymentOrders(ctx context.Context, query ListAdminPaymentOrdersQuery) ([]*Order, int64, error)
	Create(ctx context.Context, order *Order) error
	MarkPaid(ctx context.Context, userID, orderGroupID int64, paidAt time.Time) (*Order, error)
	MarkExpired(ctx context.Context, userID, orderGroupID int64, expiredAt time.Time) (*Order, error)
}

package payment

import (
	"context"
	"time"
)

type Repository interface {
	FindByOrderGroup(ctx context.Context, userID, orderGroupID int64) (*Order, error)
	Create(ctx context.Context, order *Order) error
	MarkPaid(ctx context.Context, userID, orderGroupID int64, paidAt time.Time) (*Order, error)
	MarkExpired(ctx context.Context, userID, orderGroupID int64, expiredAt time.Time) (*Order, error)
}

package order

import (
	"context"
	"time"
)

type ListBuyerOrderGroupsQuery struct {
	UserID   int64
	Page     int32
	PageSize int32
	Status   string
}

type Repository interface {
	SubmitOrder(ctx context.Context, group *Group) error
	ListBuyerOrderGroups(ctx context.Context, query ListBuyerOrderGroupsQuery) ([]*GroupSummary, int64, error)
	GetBuyerOrderGroupDetail(ctx context.Context, userID, orderGroupID int64) (*Group, error)
	GetOrderGroupPaymentInfo(ctx context.Context, userID, orderGroupID int64) (*PaymentInfo, error)
	MarkOrderGroupPaid(ctx context.Context, userID, orderGroupID int64, paidAt time.Time) (*PaymentInfo, error)
	CloseOrderGroupByPaymentTimeout(ctx context.Context, userID, orderGroupID int64, closedAt time.Time) (*PaymentInfo, error)
}

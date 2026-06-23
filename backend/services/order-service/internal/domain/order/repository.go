package order

import "context"

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
}

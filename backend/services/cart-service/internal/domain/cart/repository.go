package cart

import "context"

type Repository interface {
	CreateItem(ctx context.Context, item *Item) error
	UpdateItem(ctx context.Context, item *Item) error
	DeleteItem(ctx context.Context, userID, itemID int64) error
	FindItemByID(ctx context.Context, userID, itemID int64) (*Item, error)
	FindItemByUserAndSKU(ctx context.Context, userID, skuID int64) (*Item, error)
	ListItemsByUserID(ctx context.Context, userID int64) ([]*Item, error)
}

package buyer

import (
	"context"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
)

type CartItem struct {
	ID                int64
	UserID            int64
	SKUID             int64
	ProductID         int64
	ShopID            int64
	ShopName          string
	ProductTitle      string
	ProductSubTitle   string
	MainImageURL      string
	SKUName           string
	PriceAmount       int64
	Currency          string
	Stock             int32
	Quantity          int32
	Selected          bool
	ReviewStatus      string
	ProductSaleStatus string
	SKUSaleStatus     string
	Available         bool
	CreatedAt         int64
	UpdatedAt         int64
}

type AddCartItemInput struct {
	UserID   int64
	SKUID    int64
	Quantity int32
}

type UpdateCartItemInput struct {
	ID       int64
	UserID   int64
	Quantity int32
	Selected bool
}

type DeleteCartItemInput struct {
	ID     int64
	UserID int64
}

type CartServiceClient interface {
	AddCartItem(ctx context.Context, input AddCartItemInput) (*CartItem, error)
	UpdateCartItem(ctx context.Context, input UpdateCartItemInput) (*CartItem, error)
	DeleteCartItem(ctx context.Context, input DeleteCartItemInput) error
	ListCartItems(ctx context.Context, userID int64) ([]*CartItem, error)
}

type CartBuyerService struct {
	client CartServiceClient
}

func NewCartBuyerService(client CartServiceClient) *CartBuyerService {
	return &CartBuyerService{client: client}
}

func (s *CartBuyerService) Add(ctx context.Context, input AddCartItemInput) (*CartItem, error) {
	if s.client == nil {
		return nil, appErrors.Internal("cart service is not configured")
	}

	return s.client.AddCartItem(ctx, input)
}

func (s *CartBuyerService) Update(ctx context.Context, input UpdateCartItemInput) (*CartItem, error) {
	if s.client == nil {
		return nil, appErrors.Internal("cart service is not configured")
	}

	return s.client.UpdateCartItem(ctx, input)
}

func (s *CartBuyerService) Delete(ctx context.Context, input DeleteCartItemInput) error {
	if s.client == nil {
		return appErrors.Internal("cart service is not configured")
	}

	return s.client.DeleteCartItem(ctx, input)
}

func (s *CartBuyerService) List(ctx context.Context, userID int64) ([]*CartItem, error) {
	if s.client == nil {
		return nil, appErrors.Internal("cart service is not configured")
	}

	return s.client.ListCartItems(ctx, userID)
}

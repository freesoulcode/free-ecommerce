package buyer

import (
	"context"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
)

type OrderAddressSnapshot struct {
	ID            int64
	OrderGroupID  int64
	UserID        int64
	ReceiverName  string
	ReceiverPhone string
	CountryCode   string
	Province      string
	City          string
	District      string
	AddressLine1  string
	AddressLine2  string
	PostalCode    string
	Tag           string
	CreatedAt     int64
	UpdatedAt     int64
}

type OrderItem struct {
	ID                        int64
	OrderGroupID              int64
	ShopOrderID               int64
	UserID                    int64
	ShopID                    int64
	ProductID                 int64
	SKUID                     int64
	ProductTitle              string
	ProductSubTitle           string
	MainImageURL              string
	SKUName                   string
	PriceAmount               int64
	Currency                  string
	Quantity                  int32
	ItemAmount                int64
	ReviewStatusSnapshot      string
	ProductSaleStatusSnapshot string
	SKUSaleStatusSnapshot     string
	CreatedAt                 int64
	UpdatedAt                 int64
}

type ShopOrderSummary struct {
	ID             int64
	OrderGroupID   int64
	UserID         int64
	ShopID         int64
	ShopName       string
	Status         string
	ItemAmount     int64
	ShippingAmount int64
	PayAmount      int64
	Currency       string
	ItemCount      int32
	PaidAt         int64
	CreatedAt      int64
	UpdatedAt      int64
}

type ShopOrder struct {
	ID             int64
	OrderGroupID   int64
	UserID         int64
	ShopID         int64
	ShopName       string
	Status         string
	ItemAmount     int64
	ShippingAmount int64
	PayAmount      int64
	Currency       string
	ItemCount      int32
	Items          []*OrderItem
	PaidAt         int64
	CreatedAt      int64
	UpdatedAt      int64
}

type OrderGroupSummary struct {
	ID                  int64
	UserID              int64
	Status              string
	Source              string
	TotalItemAmount     int64
	TotalShippingAmount int64
	TotalPayAmount      int64
	Currency            string
	ShopOrderCount      int32
	ItemCount           int32
	PaymentDeadlineAt   int64
	PaidAt              int64
	ShopOrders          []*ShopOrderSummary
	CreatedAt           int64
	UpdatedAt           int64
}

type OrderGroupDetail struct {
	ID                  int64
	UserID              int64
	Status              string
	Source              string
	TotalItemAmount     int64
	TotalShippingAmount int64
	TotalPayAmount      int64
	Currency            string
	ShopOrderCount      int32
	ItemCount           int32
	PaymentDeadlineAt   int64
	PaidAt              int64
	Address             *OrderAddressSnapshot
	ShopOrders          []*ShopOrder
	CreatedAt           int64
	UpdatedAt           int64
}

type MerchantShopOrderSummary struct {
	ID                int64
	OrderGroupID      int64
	UserID            int64
	ShopID            int64
	ShopName          string
	Status            string
	ItemAmount        int64
	ShippingAmount    int64
	PayAmount         int64
	Currency          string
	ItemCount         int32
	PaidAt            int64
	CreatedAt         int64
	UpdatedAt         int64
	OrderGroupStatus  string
	PaymentDeadlineAt int64
}

type MerchantShopOrderDetail struct {
	OrderGroupID      int64
	UserID            int64
	OrderGroupStatus  string
	Source            string
	PaymentDeadlineAt int64
	PaidAt            int64
	Address           *OrderAddressSnapshot
	ShopOrder         *ShopOrder
}

type SubmitOrderInput struct {
	UserID      int64
	AddressID   int64
	CartItemIDs []int64
	Source      string
}

type ListOrdersInput struct {
	UserID   int64
	Page     int32
	PageSize int32
	Status   string
}

type ListOrdersResult struct {
	OrderGroups []*OrderGroupSummary
	Total       int64
	Page        int32
	PageSize    int32
}

type ListMerchantShopOrdersInput struct {
	ShopID   int64
	Page     int32
	PageSize int32
	Status   string
}

type ListMerchantShopOrdersResult struct {
	ShopOrders []*MerchantShopOrderSummary
	Total      int64
	Page       int32
	PageSize   int32
}

type OrderServiceClient interface {
	SubmitOrder(ctx context.Context, input SubmitOrderInput) (*OrderGroupDetail, error)
	ListBuyerOrderGroups(ctx context.Context, input ListOrdersInput) (*ListOrdersResult, error)
	GetBuyerOrderGroupDetail(ctx context.Context, userID, orderGroupID int64) (*OrderGroupDetail, error)
	MarkBuyerShopOrderReceived(ctx context.Context, userID, shopOrderID int64) (*OrderGroupDetail, error)
	ListMerchantShopOrders(ctx context.Context, input ListMerchantShopOrdersInput) (*ListMerchantShopOrdersResult, error)
	GetMerchantShopOrderDetail(ctx context.Context, shopID, shopOrderID int64) (*MerchantShopOrderDetail, error)
	MarkMerchantShopOrderProcessing(ctx context.Context, shopID, shopOrderID int64) (*MerchantShopOrderDetail, error)
	MarkMerchantShopOrderShipped(ctx context.Context, shopID, shopOrderID int64) (*MerchantShopOrderDetail, error)
}

type OrderBuyerService struct{ client OrderServiceClient }

type MerchantOrderBuyerService struct{ client OrderServiceClient }

func NewOrderBuyerService(client OrderServiceClient) *OrderBuyerService {
	return &OrderBuyerService{client: client}
}

func NewMerchantOrderBuyerService(client OrderServiceClient) *MerchantOrderBuyerService {
	return &MerchantOrderBuyerService{client: client}
}

func (s *OrderBuyerService) Submit(ctx context.Context, input SubmitOrderInput) (*OrderGroupDetail, error) {
	if s.client == nil {
		return nil, appErrors.Internal("order service is not configured")
	}
	return s.client.SubmitOrder(ctx, input)
}

func (s *OrderBuyerService) List(ctx context.Context, input ListOrdersInput) (*ListOrdersResult, error) {
	if s.client == nil {
		return nil, appErrors.Internal("order service is not configured")
	}
	return s.client.ListBuyerOrderGroups(ctx, input)
}

func (s *OrderBuyerService) Detail(ctx context.Context, userID, orderGroupID int64) (*OrderGroupDetail, error) {
	if s.client == nil {
		return nil, appErrors.Internal("order service is not configured")
	}
	return s.client.GetBuyerOrderGroupDetail(ctx, userID, orderGroupID)
}

func (s *OrderBuyerService) MarkReceived(ctx context.Context, userID, shopOrderID int64) (*OrderGroupDetail, error) {
	if s.client == nil {
		return nil, appErrors.Internal("order service is not configured")
	}
	return s.client.MarkBuyerShopOrderReceived(ctx, userID, shopOrderID)
}

func (s *MerchantOrderBuyerService) List(ctx context.Context, input ListMerchantShopOrdersInput) (*ListMerchantShopOrdersResult, error) {
	if s.client == nil {
		return nil, appErrors.Internal("order service is not configured")
	}
	return s.client.ListMerchantShopOrders(ctx, input)
}

func (s *MerchantOrderBuyerService) Detail(ctx context.Context, shopID, shopOrderID int64) (*MerchantShopOrderDetail, error) {
	if s.client == nil {
		return nil, appErrors.Internal("order service is not configured")
	}
	return s.client.GetMerchantShopOrderDetail(ctx, shopID, shopOrderID)
}

func (s *MerchantOrderBuyerService) MarkProcessing(ctx context.Context, shopID, shopOrderID int64) (*MerchantShopOrderDetail, error) {
	if s.client == nil {
		return nil, appErrors.Internal("order service is not configured")
	}
	return s.client.MarkMerchantShopOrderProcessing(ctx, shopID, shopOrderID)
}

func (s *MerchantOrderBuyerService) MarkShipped(ctx context.Context, shopID, shopOrderID int64) (*MerchantShopOrderDetail, error) {
	if s.client == nil {
		return nil, appErrors.Internal("order service is not configured")
	}
	return s.client.MarkMerchantShopOrderShipped(ctx, shopID, shopOrderID)
}

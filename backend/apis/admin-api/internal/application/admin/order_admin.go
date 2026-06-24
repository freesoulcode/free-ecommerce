package admin

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
	ShopOrders          []*ShopOrder
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

type ListShopOrdersInput struct {
	ShopID   int64
	Page     int32
	PageSize int32
	Status   string
}

type ListShopOrdersResult struct {
	ShopOrders []*MerchantShopOrderSummary
	Total      int64
	Page       int32
	PageSize   int32
}

type ListAdminOrderGroupsInput struct {
	UserID   int64
	ShopID   int64
	Page     int32
	PageSize int32
	Status   string
}

type ListAdminOrderGroupsResult struct {
	OrderGroups []*OrderGroupSummary
	Total       int64
	Page        int32
	PageSize    int32
}

type AdminOrderGroupReader interface {
	ListAdminOrderGroups(ctx context.Context, input ListAdminOrderGroupsInput) (*ListAdminOrderGroupsResult, error)
	GetAdminOrderGroupDetail(ctx context.Context, orderGroupID int64) (*OrderGroupDetail, error)
}

type ShopOrderReader interface {
	ListMerchantShopOrders(ctx context.Context, input ListShopOrdersInput) (*ListShopOrdersResult, error)
	GetMerchantShopOrderDetail(ctx context.Context, shopID, shopOrderID int64) (*MerchantShopOrderDetail, error)
	MarkMerchantShopOrderProcessing(ctx context.Context, shopID, shopOrderID int64) (*MerchantShopOrderDetail, error)
	MarkMerchantShopOrderShipped(ctx context.Context, shopID, shopOrderID int64) (*MerchantShopOrderDetail, error)
}

type ShopOrderAdminService struct {
	reader ShopOrderReader
}

type OrderGroupAdminService struct {
	reader AdminOrderGroupReader
}

func NewShopOrderAdminService(reader ShopOrderReader) *ShopOrderAdminService {
	return &ShopOrderAdminService{reader: reader}
}

func NewOrderGroupAdminService(reader AdminOrderGroupReader) *OrderGroupAdminService {
	return &OrderGroupAdminService{reader: reader}
}

func (s *OrderGroupAdminService) List(ctx context.Context, input ListAdminOrderGroupsInput) (*ListAdminOrderGroupsResult, error) {
	if s.reader == nil {
		return nil, appErrors.Internal("order service is not configured")
	}

	return s.reader.ListAdminOrderGroups(ctx, input)
}

func (s *OrderGroupAdminService) Detail(ctx context.Context, orderGroupID int64) (*OrderGroupDetail, error) {
	if s.reader == nil {
		return nil, appErrors.Internal("order service is not configured")
	}

	return s.reader.GetAdminOrderGroupDetail(ctx, orderGroupID)
}

func (s *ShopOrderAdminService) List(ctx context.Context, input ListShopOrdersInput) (*ListShopOrdersResult, error) {
	if s.reader == nil {
		return nil, appErrors.Internal("order service is not configured")
	}

	return s.reader.ListMerchantShopOrders(ctx, input)
}

func (s *ShopOrderAdminService) Detail(ctx context.Context, shopID, shopOrderID int64) (*MerchantShopOrderDetail, error) {
	if s.reader == nil {
		return nil, appErrors.Internal("order service is not configured")
	}

	return s.reader.GetMerchantShopOrderDetail(ctx, shopID, shopOrderID)
}

func (s *ShopOrderAdminService) MarkProcessing(ctx context.Context, shopID, shopOrderID int64) (*MerchantShopOrderDetail, error) {
	if s.reader == nil {
		return nil, appErrors.Internal("order service is not configured")
	}

	return s.reader.MarkMerchantShopOrderProcessing(ctx, shopID, shopOrderID)
}

func (s *ShopOrderAdminService) MarkShipped(ctx context.Context, shopID, shopOrderID int64) (*MerchantShopOrderDetail, error) {
	if s.reader == nil {
		return nil, appErrors.Internal("order service is not configured")
	}

	return s.reader.MarkMerchantShopOrderShipped(ctx, shopID, shopOrderID)
}

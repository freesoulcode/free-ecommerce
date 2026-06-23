package order

import "time"

const (
	StatusPendingPayment     = "pending_payment"
	StatusPaid               = "paid"
	StatusMerchantProcessing = "merchant_processing"
	StatusShipped            = "shipped"
	StatusReceived           = "received"
	StatusCompleted          = "completed"
	StatusCancelled          = "cancelled"
	StatusClosed             = "closed"

	SourceCart      = "cart"
	SourceDirectBuy = "direct_buy"
)

type AddressSnapshot struct {
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
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Item struct {
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
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
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
	Items          []*Item
	PaidAt         *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
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
	PaidAt         *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Group struct {
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
	PaymentDeadlineAt   time.Time
	PaidAt              *time.Time
	Address             *AddressSnapshot
	ShopOrders          []*ShopOrder
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type GroupSummary struct {
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
	PaymentDeadlineAt   time.Time
	PaidAt              *time.Time
	ShopOrders          []*ShopOrderSummary
	CreatedAt           time.Time
	UpdatedAt           time.Time
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
	OrderGroupStatus  string
	PaymentDeadlineAt time.Time
	PaidAt            *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type MerchantShopOrderDetail struct {
	OrderGroupID      int64
	UserID            int64
	OrderGroupStatus  string
	Source            string
	PaymentDeadlineAt time.Time
	PaidAt            *time.Time
	Address           *AddressSnapshot
	ShopOrder         *ShopOrder
}

type PaymentInfo struct {
	OrderGroupID      int64
	UserID            int64
	Status            string
	TotalPayAmount    int64
	Currency          string
	PaymentDeadlineAt time.Time
	PaidAt            *time.Time
}

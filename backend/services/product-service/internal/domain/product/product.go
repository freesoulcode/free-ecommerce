package product

import "time"

const (
	ReviewStatusDraft         = "draft"
	ReviewStatusPendingReview = "pending_review"
	ReviewStatusRejected      = "review_rejected"
	ReviewStatusApproved      = "approved"
	SaleStatusOnSale          = "on_sale"
	SaleStatusOffSale         = "off_sale"
)

type Summary struct {
	ID             int64
	ShopID         int64
	ShopName       string
	Title          string
	SubTitle       string
	MainImageURL   string
	MinPriceAmount int64
	MaxPriceAmount int64
	Currency       string
	TotalStock     int32
	ReviewStatus   string
	SaleStatus     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type SKU struct {
	ID          int64
	Name        string
	PriceAmount int64
	Currency    string
	Stock       int32
	SaleStatus  string
}

type SkuBrief struct {
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
	ReviewStatus      string
	ProductSaleStatus string
	SKUSaleStatus     string
}

type Detail struct {
	ID             int64
	ShopID         int64
	ShopName       string
	Title          string
	SubTitle       string
	MainImageURL   string
	Description    string
	ReviewStatus   string
	SaleStatus     string
	MinPriceAmount int64
	MaxPriceAmount int64
	Currency       string
	TotalStock     int32
	SKUs           []*SKU
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

package product

import "time"

const (
	ReviewStatusApproved = "approved"
	SaleStatusOnSale     = "on_sale"
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
}

type SKU struct {
	ID          int64
	Name        string
	PriceAmount int64
	Currency    string
	Stock       int32
	SaleStatus  string
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

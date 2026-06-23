package cart

import (
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
)

const (
	ReviewStatusApproved = "approved"
	SaleStatusOnSale     = "on_sale"
)

type Item struct {
	ID                int64
	UserID            int64
	SKUID             int64
	Quantity          int32
	Selected          bool
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
	Available         bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
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

func NewItem(id, userID, skuID int64, quantity int32, selected bool, now time.Time) (*Item, error) {
	if userID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}
	if skuID <= 0 {
		return nil, appErrors.InvalidArgument("sku id is required")
	}
	if quantity <= 0 {
		return nil, appErrors.InvalidArgument("quantity must be greater than zero")
	}

	return &Item{
		ID:        id,
		UserID:    userID,
		SKUID:     skuID,
		Quantity:  quantity,
		Selected:  selected,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (i *Item) Update(quantity int32, selected bool, now time.Time) error {
	if i == nil {
		return appErrors.Internal("cart item is required")
	}
	if quantity <= 0 {
		return appErrors.InvalidArgument("quantity must be greater than zero")
	}

	i.Quantity = quantity
	i.Selected = selected
	i.UpdatedAt = now
	return nil
}

func (i *Item) AttachBrief(brief *SkuBrief) {
	if i == nil {
		return
	}
	if brief == nil {
		i.Available = false
		return
	}

	i.ProductID = brief.ProductID
	i.ShopID = brief.ShopID
	i.ShopName = strings.TrimSpace(brief.ShopName)
	i.ProductTitle = strings.TrimSpace(brief.ProductTitle)
	i.ProductSubTitle = strings.TrimSpace(brief.ProductSubTitle)
	i.MainImageURL = strings.TrimSpace(brief.MainImageURL)
	i.SKUName = strings.TrimSpace(brief.SKUName)
	i.PriceAmount = brief.PriceAmount
	i.Currency = strings.TrimSpace(brief.Currency)
	i.Stock = brief.Stock
	i.ReviewStatus = strings.TrimSpace(brief.ReviewStatus)
	i.ProductSaleStatus = strings.TrimSpace(brief.ProductSaleStatus)
	i.SKUSaleStatus = strings.TrimSpace(brief.SKUSaleStatus)
	i.Available = brief.IsAvailable()
}

func (b *SkuBrief) IsAvailable() bool {
	if b == nil {
		return false
	}

	return b.ReviewStatus == ReviewStatusApproved &&
		b.ProductSaleStatus == SaleStatusOnSale &&
		b.SKUSaleStatus == SaleStatusOnSale &&
		b.Stock > 0
}

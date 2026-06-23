package order

import (
	"context"
	"sort"
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domainorder "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/domain/order"
)

const fixedShippingAmount int64 = 800

const defaultPaymentTimeout = 30 * time.Minute

type Address struct {
	ID            int64
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
}

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
}

type AddressService interface {
	GetAddress(ctx context.Context, userID, addressID int64) (*Address, error)
}

type CartService interface {
	ListCartItems(ctx context.Context, userID int64) ([]*CartItem, error)
}

type SubmitOrderInput struct {
	UserID      int64
	AddressID   int64
	CartItemIDs []int64
	Source      string
}

type ListBuyerOrderGroupsInput struct {
	UserID   int64
	Page     int32
	PageSize int32
	Status   string
}

type ListBuyerOrderGroupsResult struct {
	OrderGroups []*domainorder.GroupSummary
	Total       int64
	Page        int32
	PageSize    int32
}

type SubmitOrderService struct {
	repo        domainorder.Repository
	idGenerator IDGenerator
	addressSvc  AddressService
	cartSvc     CartService
	paymentTTL  time.Duration
	now         func() time.Time
}

type ListBuyerOrderGroupsService struct {
	repo domainorder.Repository
}

type GetBuyerOrderGroupDetailService struct {
	repo domainorder.Repository
}

type GetOrderGroupPaymentInfoService struct {
	repo domainorder.Repository
	now  func() time.Time
}

type MarkOrderGroupPaidService struct {
	repo domainorder.Repository
	now  func() time.Time
}

type CloseOrderGroupByPaymentTimeoutService struct {
	repo domainorder.Repository
	now  func() time.Time
}

func NewSubmitOrderService(repo domainorder.Repository, idGenerator IDGenerator, addressSvc AddressService, cartSvc CartService, paymentTTL time.Duration, now func() time.Time) *SubmitOrderService {
	if now == nil {
		now = time.Now
	}
	if paymentTTL <= 0 {
		paymentTTL = defaultPaymentTimeout
	}
	return &SubmitOrderService{repo: repo, idGenerator: idGenerator, addressSvc: addressSvc, cartSvc: cartSvc, paymentTTL: paymentTTL, now: now}
}

func NewListBuyerOrderGroupsService(repo domainorder.Repository) *ListBuyerOrderGroupsService {
	return &ListBuyerOrderGroupsService{repo: repo}
}

func NewGetBuyerOrderGroupDetailService(repo domainorder.Repository) *GetBuyerOrderGroupDetailService {
	return &GetBuyerOrderGroupDetailService{repo: repo}
}

func NewGetOrderGroupPaymentInfoService(repo domainorder.Repository, now func() time.Time) *GetOrderGroupPaymentInfoService {
	if now == nil {
		now = time.Now
	}
	return &GetOrderGroupPaymentInfoService{repo: repo, now: now}
}

func NewMarkOrderGroupPaidService(repo domainorder.Repository, now func() time.Time) *MarkOrderGroupPaidService {
	if now == nil {
		now = time.Now
	}
	return &MarkOrderGroupPaidService{repo: repo, now: now}
}

func NewCloseOrderGroupByPaymentTimeoutService(repo domainorder.Repository, now func() time.Time) *CloseOrderGroupByPaymentTimeoutService {
	if now == nil {
		now = time.Now
	}
	return &CloseOrderGroupByPaymentTimeoutService{repo: repo, now: now}
}

func (s *SubmitOrderService) Execute(ctx context.Context, input SubmitOrderInput) (*domainorder.Group, error) {
	if input.UserID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}
	if input.AddressID <= 0 {
		return nil, appErrors.InvalidArgument("address id is required")
	}
	source := strings.TrimSpace(input.Source)
	if source == "" {
		source = domainorder.SourceCart
	}
	if source != domainorder.SourceCart {
		return nil, appErrors.InvalidArgument("source is invalid")
	}

	address, err := s.addressSvc.GetAddress(ctx, input.UserID, input.AddressID)
	if err != nil {
		return nil, err
	}
	cartItems, err := s.cartSvc.ListCartItems(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	selected, err := filterSubmitItems(cartItems, input.CartItemIDs)
	if err != nil {
		return nil, err
	}
	now := s.now().UTC()
	groupID, err := s.idGenerator.NextID()
	if err != nil {
		return nil, appErrors.Internal("generate order group id failed")
	}
	addressSnapshotID, err := s.idGenerator.NextID()
	if err != nil {
		return nil, appErrors.Internal("generate order address id failed")
	}

	group := &domainorder.Group{
		ID:                groupID,
		UserID:            input.UserID,
		Status:            domainorder.StatusPendingPayment,
		Source:            source,
		PaymentDeadlineAt: now.Add(s.paymentTTL),
		CreatedAt:         now,
		UpdatedAt:         now,
		Address: &domainorder.AddressSnapshot{
			ID:            addressSnapshotID,
			OrderGroupID:  groupID,
			UserID:        input.UserID,
			ReceiverName:  address.ReceiverName,
			ReceiverPhone: address.ReceiverPhone,
			CountryCode:   address.CountryCode,
			Province:      address.Province,
			City:          address.City,
			District:      address.District,
			AddressLine1:  address.AddressLine1,
			AddressLine2:  address.AddressLine2,
			PostalCode:    address.PostalCode,
			Tag:           address.Tag,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	shopBuckets := make(map[int64]*domainorder.ShopOrder)
	shopOrderIDs := make([]int64, 0)
	for _, cartItem := range selected {
		shopOrder, ok := shopBuckets[cartItem.ShopID]
		if !ok {
			shopOrderID, idErr := s.idGenerator.NextID()
			if idErr != nil {
				return nil, appErrors.Internal("generate shop order id failed")
			}
			shopOrder = &domainorder.ShopOrder{
				ID:             shopOrderID,
				OrderGroupID:   groupID,
				UserID:         input.UserID,
				ShopID:         cartItem.ShopID,
				ShopName:       cartItem.ShopName,
				Status:         domainorder.StatusPendingPayment,
				ShippingAmount: fixedShippingAmount,
				Currency:       cartItem.Currency,
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			shopBuckets[cartItem.ShopID] = shopOrder
			shopOrderIDs = append(shopOrderIDs, cartItem.ShopID)
		}

		if shopOrder.Currency != cartItem.Currency {
			return nil, appErrors.InvalidArgument("cross-currency order is not supported")
		}

		itemID, idErr := s.idGenerator.NextID()
		if idErr != nil {
			return nil, appErrors.Internal("generate order item id failed")
		}
		itemAmount := cartItem.PriceAmount * int64(cartItem.Quantity)
		item := &domainorder.Item{
			ID:                        itemID,
			OrderGroupID:              groupID,
			ShopOrderID:               shopOrder.ID,
			UserID:                    input.UserID,
			ShopID:                    cartItem.ShopID,
			ProductID:                 cartItem.ProductID,
			SKUID:                     cartItem.SKUID,
			ProductTitle:              cartItem.ProductTitle,
			ProductSubTitle:           cartItem.ProductSubTitle,
			MainImageURL:              cartItem.MainImageURL,
			SKUName:                   cartItem.SKUName,
			PriceAmount:               cartItem.PriceAmount,
			Currency:                  cartItem.Currency,
			Quantity:                  cartItem.Quantity,
			ItemAmount:                itemAmount,
			ReviewStatusSnapshot:      cartItem.ReviewStatus,
			ProductSaleStatusSnapshot: cartItem.ProductSaleStatus,
			SKUSaleStatusSnapshot:     cartItem.SKUSaleStatus,
			CreatedAt:                 now,
			UpdatedAt:                 now,
		}
		shopOrder.Items = append(shopOrder.Items, item)
		shopOrder.ItemAmount += itemAmount
		shopOrder.ItemCount += cartItem.Quantity
		shopOrder.PayAmount = shopOrder.ItemAmount + shopOrder.ShippingAmount
		group.TotalItemAmount += itemAmount
		group.ItemCount += cartItem.Quantity
	}

	sort.SliceStable(shopOrderIDs, func(i, j int) bool { return shopOrderIDs[i] < shopOrderIDs[j] })
	for _, shopID := range shopOrderIDs {
		shopOrder := shopBuckets[shopID]
		group.ShopOrders = append(group.ShopOrders, shopOrder)
		group.TotalShippingAmount += shopOrder.ShippingAmount
	}
	group.TotalPayAmount = group.TotalItemAmount + group.TotalShippingAmount
	group.ShopOrderCount = int32(len(group.ShopOrders))
	if len(group.ShopOrders) > 0 {
		group.Currency = group.ShopOrders[0].Currency
	}

	if err := s.repo.SubmitOrder(ctx, group); err != nil {
		return nil, err
	}

	return group, nil
}

func (s *ListBuyerOrderGroupsService) Execute(ctx context.Context, input ListBuyerOrderGroupsInput) (*ListBuyerOrderGroupsResult, error) {
	if input.UserID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}
	page := input.Page
	if page <= 0 {
		page = 1
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	groups, total, err := s.repo.ListBuyerOrderGroups(ctx, domainorder.ListBuyerOrderGroupsQuery{UserID: input.UserID, Page: page, PageSize: pageSize, Status: strings.TrimSpace(input.Status)})
	if err != nil {
		return nil, err
	}

	return &ListBuyerOrderGroupsResult{OrderGroups: groups, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *GetBuyerOrderGroupDetailService) Execute(ctx context.Context, userID, orderGroupID int64) (*domainorder.Group, error) {
	if userID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}
	if orderGroupID <= 0 {
		return nil, appErrors.InvalidArgument("order group id is required")
	}

	return s.repo.GetBuyerOrderGroupDetail(ctx, userID, orderGroupID)
}

func (s *GetOrderGroupPaymentInfoService) Execute(ctx context.Context, userID, orderGroupID int64) (*domainorder.PaymentInfo, error) {
	if userID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}
	if orderGroupID <= 0 {
		return nil, appErrors.InvalidArgument("order group id is required")
	}

	info, err := s.repo.GetOrderGroupPaymentInfo(ctx, userID, orderGroupID)
	if err != nil {
		return nil, err
	}
	if info.Status == domainorder.StatusPendingPayment && !info.PaymentDeadlineAt.IsZero() && s.now().UTC().After(info.PaymentDeadlineAt) {
		return s.repo.CloseOrderGroupByPaymentTimeout(ctx, userID, orderGroupID, s.now().UTC())
	}
	return info, nil
}

func (s *MarkOrderGroupPaidService) Execute(ctx context.Context, userID, orderGroupID int64) (*domainorder.PaymentInfo, error) {
	if userID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}
	if orderGroupID <= 0 {
		return nil, appErrors.InvalidArgument("order group id is required")
	}
	return s.repo.MarkOrderGroupPaid(ctx, userID, orderGroupID, s.now().UTC())
}

func (s *CloseOrderGroupByPaymentTimeoutService) Execute(ctx context.Context, userID, orderGroupID int64) (*domainorder.PaymentInfo, error) {
	if userID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}
	if orderGroupID <= 0 {
		return nil, appErrors.InvalidArgument("order group id is required")
	}
	return s.repo.CloseOrderGroupByPaymentTimeout(ctx, userID, orderGroupID, s.now().UTC())
}

func filterSubmitItems(items []*CartItem, cartItemIDs []int64) ([]*CartItem, error) {
	allowedIDs := make(map[int64]struct{}, len(cartItemIDs))
	if len(cartItemIDs) > 0 {
		for _, id := range cartItemIDs {
			if id <= 0 {
				return nil, appErrors.InvalidArgument("cart item id is required")
			}
			allowedIDs[id] = struct{}{}
		}
	}

	selected := make([]*CartItem, 0)
	for _, item := range items {
		if item == nil || !item.Selected {
			continue
		}
		if len(allowedIDs) > 0 {
			if _, ok := allowedIDs[item.ID]; !ok {
				continue
			}
		}
		if !item.Available {
			return nil, appErrors.InvalidArgument("selected cart item is not available")
		}
		if item.Quantity <= 0 {
			return nil, appErrors.InvalidArgument("selected cart item quantity is invalid")
		}
		if item.Quantity > item.Stock {
			return nil, appErrors.InvalidArgument("selected cart item exceeds stock")
		}
		selected = append(selected, item)
	}

	if len(selected) == 0 {
		return nil, appErrors.InvalidArgument("no selected cart items to submit")
	}

	sort.SliceStable(selected, func(i, j int) bool {
		if selected[i].ShopID == selected[j].ShopID {
			return selected[i].ID < selected[j].ID
		}
		return selected[i].ShopID < selected[j].ShopID
	})
	return selected, nil
}

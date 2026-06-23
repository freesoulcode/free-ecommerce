package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domainorder "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/domain/order"
	"gorm.io/gorm"
)

type OrderGroupModel struct {
	ID                  int64      `gorm:"column:id;primaryKey"`
	UserID              int64      `gorm:"column:user_id"`
	Status              string     `gorm:"column:status"`
	Source              string     `gorm:"column:source"`
	TotalItemAmount     int64      `gorm:"column:total_item_amount"`
	TotalShippingAmount int64      `gorm:"column:total_shipping_amount"`
	TotalPayAmount      int64      `gorm:"column:total_pay_amount"`
	Currency            string     `gorm:"column:currency"`
	ShopOrderCount      int32      `gorm:"column:shop_order_count"`
	ItemCount           int32      `gorm:"column:item_count"`
	PaymentDeadlineAt   time.Time  `gorm:"column:payment_deadline_at"`
	PaidAt              *time.Time `gorm:"column:paid_at"`
	CreatedAt           time.Time  `gorm:"column:created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at"`
}

func (OrderGroupModel) TableName() string { return "order_groups" }

type OrderGroupAddressModel struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	OrderGroupID  int64     `gorm:"column:order_group_id"`
	UserID        int64     `gorm:"column:user_id"`
	ReceiverName  string    `gorm:"column:receiver_name"`
	ReceiverPhone string    `gorm:"column:receiver_phone"`
	CountryCode   string    `gorm:"column:country_code"`
	Province      string    `gorm:"column:province"`
	City          string    `gorm:"column:city"`
	District      string    `gorm:"column:district"`
	AddressLine1  string    `gorm:"column:address_line1"`
	AddressLine2  *string   `gorm:"column:address_line2"`
	PostalCode    *string   `gorm:"column:postal_code"`
	Tag           *string   `gorm:"column:tag"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (OrderGroupAddressModel) TableName() string { return "order_group_addresses" }

type ShopOrderModel struct {
	ID             int64      `gorm:"column:id;primaryKey"`
	OrderGroupID   int64      `gorm:"column:order_group_id"`
	UserID         int64      `gorm:"column:user_id"`
	ShopID         int64      `gorm:"column:shop_id"`
	ShopName       string     `gorm:"column:shop_name"`
	Status         string     `gorm:"column:status"`
	ItemAmount     int64      `gorm:"column:item_amount"`
	ShippingAmount int64      `gorm:"column:shipping_amount"`
	PayAmount      int64      `gorm:"column:pay_amount"`
	Currency       string     `gorm:"column:currency"`
	ItemCount      int32      `gorm:"column:item_count"`
	PaidAt         *time.Time `gorm:"column:paid_at"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
}

func (ShopOrderModel) TableName() string { return "shop_orders" }

type OrderItemModel struct {
	ID                        int64     `gorm:"column:id;primaryKey"`
	OrderGroupID              int64     `gorm:"column:order_group_id"`
	ShopOrderID               int64     `gorm:"column:shop_order_id"`
	UserID                    int64     `gorm:"column:user_id"`
	ShopID                    int64     `gorm:"column:shop_id"`
	ProductID                 int64     `gorm:"column:product_id"`
	SKUID                     int64     `gorm:"column:sku_id"`
	ProductTitle              string    `gorm:"column:product_title"`
	ProductSubTitle           string    `gorm:"column:product_sub_title"`
	MainImageURL              string    `gorm:"column:main_image_url"`
	SKUName                   string    `gorm:"column:sku_name"`
	PriceAmount               int64     `gorm:"column:price_amount"`
	Currency                  string    `gorm:"column:currency"`
	Quantity                  int32     `gorm:"column:quantity"`
	ItemAmount                int64     `gorm:"column:item_amount"`
	ReviewStatusSnapshot      string    `gorm:"column:review_status_snapshot"`
	ProductSaleStatusSnapshot string    `gorm:"column:product_sale_status_snapshot"`
	SKUSaleStatusSnapshot     string    `gorm:"column:sku_sale_status_snapshot"`
	CreatedAt                 time.Time `gorm:"column:created_at"`
	UpdatedAt                 time.Time `gorm:"column:updated_at"`
}

func (OrderItemModel) TableName() string { return "order_items" }

type OrderRepository struct{ db *gorm.DB }

type merchantShopOrderRow struct {
	ID                int64      `gorm:"column:id"`
	OrderGroupID      int64      `gorm:"column:order_group_id"`
	UserID            int64      `gorm:"column:user_id"`
	ShopID            int64      `gorm:"column:shop_id"`
	ShopName          string     `gorm:"column:shop_name"`
	Status            string     `gorm:"column:status"`
	ItemAmount        int64      `gorm:"column:item_amount"`
	ShippingAmount    int64      `gorm:"column:shipping_amount"`
	PayAmount         int64      `gorm:"column:pay_amount"`
	Currency          string     `gorm:"column:currency"`
	ItemCount         int32      `gorm:"column:item_count"`
	PaidAt            *time.Time `gorm:"column:paid_at"`
	CreatedAt         time.Time  `gorm:"column:created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at"`
	OrderGroupStatus  string     `gorm:"column:order_group_status"`
	PaymentDeadlineAt time.Time  `gorm:"column:payment_deadline_at"`
}

func NewOrderRepository(db *gorm.DB) *OrderRepository { return &OrderRepository{db: db} }

func (r *OrderRepository) SubmitOrder(ctx context.Context, group *domainorder.Group) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		groupModel := toOrderGroupModel(group)
		if err := tx.Create(&groupModel).Error; err != nil {
			return fmt.Errorf("create order group: %w", err)
		}
		if group.Address != nil {
			addressModel := toAddressModel(group.Address)
			if err := tx.Create(&addressModel).Error; err != nil {
				return fmt.Errorf("create order address: %w", err)
			}
		}
		for _, shopOrder := range group.ShopOrders {
			shopOrderModel := toShopOrderModel(shopOrder)
			if err := tx.Create(&shopOrderModel).Error; err != nil {
				return fmt.Errorf("create shop order: %w", err)
			}
			for _, item := range shopOrder.Items {
				itemModel := toOrderItemModel(item)
				if err := tx.Create(&itemModel).Error; err != nil {
					return fmt.Errorf("create order item: %w", err)
				}
			}
		}
		return nil
	})
}

func (r *OrderRepository) ListBuyerOrderGroups(ctx context.Context, query domainorder.ListBuyerOrderGroupsQuery) ([]*domainorder.GroupSummary, int64, error) {
	db := r.db.WithContext(ctx).Model(&OrderGroupModel{}).Where("user_id = ?", query.UserID)
	if status := strings.TrimSpace(query.Status); status != "" {
		db = db.Where("status = ?", status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count order groups: %w", err)
	}

	var groupModels []OrderGroupModel
	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("created_at DESC, id DESC").Limit(int(query.PageSize)).Offset(int(offset)).Find(&groupModels).Error; err != nil {
		return nil, 0, fmt.Errorf("list order groups: %w", err)
	}
	if len(groupModels) == 0 {
		return []*domainorder.GroupSummary{}, total, nil
	}

	groupIDs := make([]int64, 0, len(groupModels))
	for _, model := range groupModels {
		groupIDs = append(groupIDs, model.ID)
	}
	var shopOrderModels []ShopOrderModel
	if err := r.db.WithContext(ctx).Where("order_group_id IN ?", groupIDs).Order("created_at ASC, id ASC").Find(&shopOrderModels).Error; err != nil {
		return nil, 0, fmt.Errorf("list shop orders for summary: %w", err)
	}
	shopOrdersByGroup := make(map[int64][]*domainorder.ShopOrderSummary, len(groupIDs))
	for _, model := range shopOrderModels {
		shopOrdersByGroup[model.OrderGroupID] = append(shopOrdersByGroup[model.OrderGroupID], toShopOrderSummary(model))
	}

	groups := make([]*domainorder.GroupSummary, 0, len(groupModels))
	for _, model := range groupModels {
		groups = append(groups, &domainorder.GroupSummary{
			ID:                  model.ID,
			UserID:              model.UserID,
			Status:              model.Status,
			Source:              model.Source,
			TotalItemAmount:     model.TotalItemAmount,
			TotalShippingAmount: model.TotalShippingAmount,
			TotalPayAmount:      model.TotalPayAmount,
			Currency:            model.Currency,
			ShopOrderCount:      model.ShopOrderCount,
			ItemCount:           model.ItemCount,
			PaymentDeadlineAt:   model.PaymentDeadlineAt,
			PaidAt:              model.PaidAt,
			ShopOrders:          shopOrdersByGroup[model.ID],
			CreatedAt:           model.CreatedAt,
			UpdatedAt:           model.UpdatedAt,
		})
	}

	return groups, total, nil
}

func (r *OrderRepository) GetBuyerOrderGroupDetail(ctx context.Context, userID, orderGroupID int64) (*domainorder.Group, error) {
	var groupModel OrderGroupModel
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", orderGroupID, userID).First(&groupModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("order group not found")
		}
		return nil, fmt.Errorf("find order group detail: %w", err)
	}

	group := &domainorder.Group{
		ID:                  groupModel.ID,
		UserID:              groupModel.UserID,
		Status:              groupModel.Status,
		Source:              groupModel.Source,
		TotalItemAmount:     groupModel.TotalItemAmount,
		TotalShippingAmount: groupModel.TotalShippingAmount,
		TotalPayAmount:      groupModel.TotalPayAmount,
		Currency:            groupModel.Currency,
		ShopOrderCount:      groupModel.ShopOrderCount,
		ItemCount:           groupModel.ItemCount,
		PaymentDeadlineAt:   groupModel.PaymentDeadlineAt,
		PaidAt:              groupModel.PaidAt,
		CreatedAt:           groupModel.CreatedAt,
		UpdatedAt:           groupModel.UpdatedAt,
	}

	var addressModel OrderGroupAddressModel
	if err := r.db.WithContext(ctx).Where("order_group_id = ?", orderGroupID).First(&addressModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("order address not found")
		}
		return nil, fmt.Errorf("find order address detail: %w", err)
	}
	group.Address = toAddressSnapshot(addressModel)

	var shopOrderModels []ShopOrderModel
	if err := r.db.WithContext(ctx).Where("order_group_id = ?", orderGroupID).Order("created_at ASC, id ASC").Find(&shopOrderModels).Error; err != nil {
		return nil, fmt.Errorf("list shop orders detail: %w", err)
	}
	var itemModels []OrderItemModel
	if err := r.db.WithContext(ctx).Where("order_group_id = ?", orderGroupID).Order("created_at ASC, id ASC").Find(&itemModels).Error; err != nil {
		return nil, fmt.Errorf("list order items detail: %w", err)
	}
	itemsByShopOrder := make(map[int64][]*domainorder.Item, len(shopOrderModels))
	for _, model := range itemModels {
		itemsByShopOrder[model.ShopOrderID] = append(itemsByShopOrder[model.ShopOrderID], toOrderItem(model))
	}
	for _, model := range shopOrderModels {
		shopOrder := &domainorder.ShopOrder{
			ID:             model.ID,
			OrderGroupID:   model.OrderGroupID,
			UserID:         model.UserID,
			ShopID:         model.ShopID,
			ShopName:       model.ShopName,
			Status:         model.Status,
			ItemAmount:     model.ItemAmount,
			ShippingAmount: model.ShippingAmount,
			PayAmount:      model.PayAmount,
			Currency:       model.Currency,
			ItemCount:      model.ItemCount,
			Items:          itemsByShopOrder[model.ID],
			PaidAt:         model.PaidAt,
			CreatedAt:      model.CreatedAt,
			UpdatedAt:      model.UpdatedAt,
		}
		group.ShopOrders = append(group.ShopOrders, shopOrder)
	}

	return group, nil
}

func (r *OrderRepository) ListMerchantShopOrders(ctx context.Context, query domainorder.ListMerchantShopOrdersQuery) ([]*domainorder.MerchantShopOrderSummary, int64, error) {
	db := r.db.WithContext(ctx).
		Model(&ShopOrderModel{}).
		Joins("JOIN order_groups ON order_groups.id = shop_orders.order_group_id").
		Where("shop_orders.shop_id = ?", query.ShopID)
	if status := strings.TrimSpace(query.Status); status != "" {
		db = db.Where("shop_orders.status = ?", status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count merchant shop orders: %w", err)
	}

	var rows []merchantShopOrderRow
	offset := (query.Page - 1) * query.PageSize
	if err := db.Select("shop_orders.*, order_groups.status AS order_group_status, order_groups.payment_deadline_at AS payment_deadline_at").Order("shop_orders.created_at DESC, shop_orders.id DESC").Limit(int(query.PageSize)).Offset(int(offset)).Scan(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list merchant shop orders: %w", err)
	}
	if len(rows) == 0 {
		return []*domainorder.MerchantShopOrderSummary{}, total, nil
	}

	items := make([]*domainorder.MerchantShopOrderSummary, 0, len(rows))
	for _, row := range rows {
		items = append(items, toMerchantShopOrderSummary(row))
	}
	return items, total, nil
}

func (r *OrderRepository) GetMerchantShopOrderDetail(ctx context.Context, shopID, shopOrderID int64) (*domainorder.MerchantShopOrderDetail, error) {
	shopOrderModel, err := r.getShopOrderModel(ctx, shopID, shopOrderID)
	if err != nil {
		return nil, err
	}
	groupModel, err := r.getOrderGroupModelByID(ctx, shopOrderModel.OrderGroupID)
	if err != nil {
		return nil, err
	}
	var addressModel OrderGroupAddressModel
	if err := r.db.WithContext(ctx).Where("order_group_id = ?", shopOrderModel.OrderGroupID).First(&addressModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("order address not found")
		}
		return nil, fmt.Errorf("find merchant order address detail: %w", err)
	}
	var itemModels []OrderItemModel
	if err := r.db.WithContext(ctx).Where("shop_order_id = ?", shopOrderID).Order("created_at ASC, id ASC").Find(&itemModels).Error; err != nil {
		return nil, fmt.Errorf("list merchant order items detail: %w", err)
	}
	items := make([]*domainorder.Item, 0, len(itemModels))
	for _, model := range itemModels {
		items = append(items, toOrderItem(model))
	}
	shopOrder := toShopOrder(shopOrderModel)
	shopOrder.Items = items
	return &domainorder.MerchantShopOrderDetail{
		OrderGroupID:      groupModel.ID,
		UserID:            groupModel.UserID,
		OrderGroupStatus:  groupModel.Status,
		Source:            groupModel.Source,
		PaymentDeadlineAt: groupModel.PaymentDeadlineAt,
		PaidAt:            groupModel.PaidAt,
		Address:           toAddressSnapshot(addressModel),
		ShopOrder:         shopOrder,
	}, nil
}

func (r *OrderRepository) MarkMerchantShopOrderProcessing(ctx context.Context, shopID, shopOrderID int64, updatedAt time.Time) (*domainorder.MerchantShopOrderDetail, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		shopOrderModel, err := r.getShopOrderModelForUpdate(tx, shopID, shopOrderID)
		if err != nil {
			return err
		}
		groupModel, err := r.getOrderGroupModelByIDForUpdate(tx, shopOrderModel.OrderGroupID)
		if err != nil {
			return err
		}

		switch shopOrderModel.Status {
		case domainorder.StatusMerchantProcessing:
			return nil
		case domainorder.StatusPaid:
		default:
			return appErrors.InvalidArgument("shop order status does not allow merchant processing")
		}

		if groupModel.Status != domainorder.StatusPaid && groupModel.Status != domainorder.StatusMerchantProcessing {
			return appErrors.InvalidArgument("order group status does not allow merchant processing")
		}

		if shopOrderModel.Status == domainorder.StatusPaid {
			if err := tx.Model(&ShopOrderModel{}).Where("id = ? AND shop_id = ?", shopOrderID, shopID).Updates(map[string]any{"status": domainorder.StatusMerchantProcessing, "updated_at": updatedAt}).Error; err != nil {
				return fmt.Errorf("mark merchant shop order processing: %w", err)
			}
		}
		if groupModel.Status == domainorder.StatusPaid {
			if err := tx.Model(&OrderGroupModel{}).Where("id = ?", groupModel.ID).Updates(map[string]any{"status": domainorder.StatusMerchantProcessing, "updated_at": updatedAt}).Error; err != nil {
				return fmt.Errorf("mark order group merchant processing: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetMerchantShopOrderDetail(ctx, shopID, shopOrderID)
}

func (r *OrderRepository) MarkMerchantShopOrderCompleted(ctx context.Context, shopID, shopOrderID int64, updatedAt time.Time) (*domainorder.MerchantShopOrderDetail, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		shopOrderModel, err := r.getShopOrderModelForUpdate(tx, shopID, shopOrderID)
		if err != nil {
			return err
		}
		groupModel, err := r.getOrderGroupModelByIDForUpdate(tx, shopOrderModel.OrderGroupID)
		if err != nil {
			return err
		}

		switch shopOrderModel.Status {
		case domainorder.StatusCompleted:
			return nil
		case domainorder.StatusMerchantProcessing:
		default:
			return appErrors.InvalidArgument("shop order status does not allow completion")
		}

		if groupModel.Status != domainorder.StatusMerchantProcessing && groupModel.Status != domainorder.StatusCompleted {
			return appErrors.InvalidArgument("order group status does not allow completion")
		}

		if shopOrderModel.Status == domainorder.StatusMerchantProcessing {
			if err := tx.Model(&ShopOrderModel{}).Where("id = ? AND shop_id = ?", shopOrderID, shopID).Updates(map[string]any{"status": domainorder.StatusCompleted, "updated_at": updatedAt}).Error; err != nil {
				return fmt.Errorf("mark merchant shop order completed: %w", err)
			}
		}

		var remaining int64
		if err := tx.Model(&ShopOrderModel{}).Where("order_group_id = ? AND status <> ?", shopOrderModel.OrderGroupID, domainorder.StatusCompleted).Count(&remaining).Error; err != nil {
			return fmt.Errorf("count remaining shop orders: %w", err)
		}
		groupStatus := domainorder.StatusMerchantProcessing
		if remaining == 0 {
			groupStatus = domainorder.StatusCompleted
		}
		if groupModel.Status != groupStatus {
			if err := tx.Model(&OrderGroupModel{}).Where("id = ?", groupModel.ID).Updates(map[string]any{"status": groupStatus, "updated_at": updatedAt}).Error; err != nil {
				return fmt.Errorf("update order group completion status: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.GetMerchantShopOrderDetail(ctx, shopID, shopOrderID)
}

func (r *OrderRepository) GetOrderGroupPaymentInfo(ctx context.Context, userID, orderGroupID int64) (*domainorder.PaymentInfo, error) {
	groupModel, err := r.getOrderGroupModel(ctx, userID, orderGroupID)
	if err != nil {
		return nil, err
	}
	info := toPaymentInfo(groupModel)
	return &info, nil
}

func (r *OrderRepository) MarkOrderGroupPaid(ctx context.Context, userID, orderGroupID int64, paidAt time.Time) (*domainorder.PaymentInfo, error) {
	var result *domainorder.PaymentInfo
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		groupModel, err := r.getOrderGroupModelForUpdate(tx, userID, orderGroupID)
		if err != nil {
			return err
		}

		switch groupModel.Status {
		case domainorder.StatusPaid:
			info := toPaymentInfo(groupModel)
			result = &info
			return nil
		case domainorder.StatusPendingPayment:
			if !groupModel.PaymentDeadlineAt.IsZero() && paidAt.After(groupModel.PaymentDeadlineAt) {
				return appErrors.InvalidArgument("order group payment deadline has expired")
			}
		default:
			return appErrors.InvalidArgument("order group status does not allow payment")
		}

		updates := map[string]any{
			"status":     domainorder.StatusPaid,
			"paid_at":    paidAt,
			"updated_at": paidAt,
		}
		if err := tx.Model(&OrderGroupModel{}).Where("id = ? AND user_id = ?", orderGroupID, userID).Updates(updates).Error; err != nil {
			return fmt.Errorf("mark order group paid: %w", err)
		}
		if err := tx.Model(&ShopOrderModel{}).Where("order_group_id = ?", orderGroupID).Updates(map[string]any{"status": domainorder.StatusPaid, "paid_at": paidAt, "updated_at": paidAt}).Error; err != nil {
			return fmt.Errorf("mark shop orders paid: %w", err)
		}

		groupModel.Status = domainorder.StatusPaid
		groupModel.PaidAt = &paidAt
		groupModel.UpdatedAt = paidAt
		info := toPaymentInfo(groupModel)
		result = &info
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *OrderRepository) CloseOrderGroupByPaymentTimeout(ctx context.Context, userID, orderGroupID int64, closedAt time.Time) (*domainorder.PaymentInfo, error) {
	var result *domainorder.PaymentInfo
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		groupModel, err := r.getOrderGroupModelForUpdate(tx, userID, orderGroupID)
		if err != nil {
			return err
		}

		switch groupModel.Status {
		case domainorder.StatusClosed, domainorder.StatusPaid:
			info := toPaymentInfo(groupModel)
			result = &info
			return nil
		case domainorder.StatusPendingPayment:
			if !groupModel.PaymentDeadlineAt.IsZero() && closedAt.Before(groupModel.PaymentDeadlineAt) {
				return appErrors.InvalidArgument("order group payment is not expired yet")
			}
		default:
			return appErrors.InvalidArgument("order group status does not allow close by payment timeout")
		}

		if err := tx.Model(&OrderGroupModel{}).Where("id = ? AND user_id = ?", orderGroupID, userID).Updates(map[string]any{"status": domainorder.StatusClosed, "updated_at": closedAt}).Error; err != nil {
			return fmt.Errorf("close order group by payment timeout: %w", err)
		}
		if err := tx.Model(&ShopOrderModel{}).Where("order_group_id = ?", orderGroupID).Updates(map[string]any{"status": domainorder.StatusClosed, "updated_at": closedAt}).Error; err != nil {
			return fmt.Errorf("close shop orders by payment timeout: %w", err)
		}

		groupModel.Status = domainorder.StatusClosed
		groupModel.UpdatedAt = closedAt
		info := toPaymentInfo(groupModel)
		result = &info
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *OrderRepository) getOrderGroupModel(ctx context.Context, userID, orderGroupID int64) (OrderGroupModel, error) {
	var groupModel OrderGroupModel
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", orderGroupID, userID).First(&groupModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return OrderGroupModel{}, appErrors.NotFound("order group not found")
		}
		return OrderGroupModel{}, fmt.Errorf("find order group: %w", err)
	}
	return groupModel, nil
}

func (r *OrderRepository) getOrderGroupModelForUpdate(tx *gorm.DB, userID, orderGroupID int64) (OrderGroupModel, error) {
	var groupModel OrderGroupModel
	if err := tx.Where("id = ? AND user_id = ?", orderGroupID, userID).Take(&groupModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return OrderGroupModel{}, appErrors.NotFound("order group not found")
		}
		return OrderGroupModel{}, fmt.Errorf("find order group for update: %w", err)
	}
	return groupModel, nil
}

func (r *OrderRepository) getOrderGroupModelByID(ctx context.Context, orderGroupID int64) (OrderGroupModel, error) {
	var groupModel OrderGroupModel
	if err := r.db.WithContext(ctx).Where("id = ?", orderGroupID).First(&groupModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return OrderGroupModel{}, appErrors.NotFound("order group not found")
		}
		return OrderGroupModel{}, fmt.Errorf("find order group by id: %w", err)
	}
	return groupModel, nil
}

func (r *OrderRepository) getOrderGroupModelByIDForUpdate(tx *gorm.DB, orderGroupID int64) (OrderGroupModel, error) {
	var groupModel OrderGroupModel
	if err := tx.Where("id = ?", orderGroupID).Take(&groupModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return OrderGroupModel{}, appErrors.NotFound("order group not found")
		}
		return OrderGroupModel{}, fmt.Errorf("find order group by id for update: %w", err)
	}
	return groupModel, nil
}

func (r *OrderRepository) getShopOrderModel(ctx context.Context, shopID, shopOrderID int64) (ShopOrderModel, error) {
	var shopOrderModel ShopOrderModel
	if err := r.db.WithContext(ctx).Where("id = ? AND shop_id = ?", shopOrderID, shopID).First(&shopOrderModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ShopOrderModel{}, appErrors.NotFound("shop order not found")
		}
		return ShopOrderModel{}, fmt.Errorf("find shop order: %w", err)
	}
	return shopOrderModel, nil
}

func (r *OrderRepository) getShopOrderModelForUpdate(tx *gorm.DB, shopID, shopOrderID int64) (ShopOrderModel, error) {
	var shopOrderModel ShopOrderModel
	if err := tx.Where("id = ? AND shop_id = ?", shopOrderID, shopID).Take(&shopOrderModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ShopOrderModel{}, appErrors.NotFound("shop order not found")
		}
		return ShopOrderModel{}, fmt.Errorf("find shop order for update: %w", err)
	}
	return shopOrderModel, nil
}

func nullableString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func toOrderGroupModel(group *domainorder.Group) OrderGroupModel {
	return OrderGroupModel{ID: group.ID, UserID: group.UserID, Status: group.Status, Source: group.Source, TotalItemAmount: group.TotalItemAmount, TotalShippingAmount: group.TotalShippingAmount, TotalPayAmount: group.TotalPayAmount, Currency: group.Currency, ShopOrderCount: group.ShopOrderCount, ItemCount: group.ItemCount, PaymentDeadlineAt: group.PaymentDeadlineAt, PaidAt: group.PaidAt, CreatedAt: group.CreatedAt, UpdatedAt: group.UpdatedAt}
}

func toAddressModel(address *domainorder.AddressSnapshot) OrderGroupAddressModel {
	return OrderGroupAddressModel{ID: address.ID, OrderGroupID: address.OrderGroupID, UserID: address.UserID, ReceiverName: address.ReceiverName, ReceiverPhone: address.ReceiverPhone, CountryCode: address.CountryCode, Province: address.Province, City: address.City, District: address.District, AddressLine1: address.AddressLine1, AddressLine2: nullableString(address.AddressLine2), PostalCode: nullableString(address.PostalCode), Tag: nullableString(address.Tag), CreatedAt: address.CreatedAt, UpdatedAt: address.UpdatedAt}
}

func toShopOrderModel(shopOrder *domainorder.ShopOrder) ShopOrderModel {
	return ShopOrderModel{ID: shopOrder.ID, OrderGroupID: shopOrder.OrderGroupID, UserID: shopOrder.UserID, ShopID: shopOrder.ShopID, ShopName: shopOrder.ShopName, Status: shopOrder.Status, ItemAmount: shopOrder.ItemAmount, ShippingAmount: shopOrder.ShippingAmount, PayAmount: shopOrder.PayAmount, Currency: shopOrder.Currency, ItemCount: shopOrder.ItemCount, PaidAt: shopOrder.PaidAt, CreatedAt: shopOrder.CreatedAt, UpdatedAt: shopOrder.UpdatedAt}
}

func toOrderItemModel(item *domainorder.Item) OrderItemModel {
	return OrderItemModel{ID: item.ID, OrderGroupID: item.OrderGroupID, ShopOrderID: item.ShopOrderID, UserID: item.UserID, ShopID: item.ShopID, ProductID: item.ProductID, SKUID: item.SKUID, ProductTitle: item.ProductTitle, ProductSubTitle: item.ProductSubTitle, MainImageURL: item.MainImageURL, SKUName: item.SKUName, PriceAmount: item.PriceAmount, Currency: item.Currency, Quantity: item.Quantity, ItemAmount: item.ItemAmount, ReviewStatusSnapshot: item.ReviewStatusSnapshot, ProductSaleStatusSnapshot: item.ProductSaleStatusSnapshot, SKUSaleStatusSnapshot: item.SKUSaleStatusSnapshot, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}
}

func toAddressSnapshot(model OrderGroupAddressModel) *domainorder.AddressSnapshot {
	return &domainorder.AddressSnapshot{ID: model.ID, OrderGroupID: model.OrderGroupID, UserID: model.UserID, ReceiverName: model.ReceiverName, ReceiverPhone: model.ReceiverPhone, CountryCode: model.CountryCode, Province: model.Province, City: model.City, District: model.District, AddressLine1: model.AddressLine1, AddressLine2: derefString(model.AddressLine2), PostalCode: derefString(model.PostalCode), Tag: derefString(model.Tag), CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func toShopOrderSummary(model ShopOrderModel) *domainorder.ShopOrderSummary {
	return &domainorder.ShopOrderSummary{ID: model.ID, OrderGroupID: model.OrderGroupID, UserID: model.UserID, ShopID: model.ShopID, ShopName: model.ShopName, Status: model.Status, ItemAmount: model.ItemAmount, ShippingAmount: model.ShippingAmount, PayAmount: model.PayAmount, Currency: model.Currency, ItemCount: model.ItemCount, PaidAt: model.PaidAt, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func toMerchantShopOrderSummary(row merchantShopOrderRow) *domainorder.MerchantShopOrderSummary {
	return &domainorder.MerchantShopOrderSummary{ID: row.ID, OrderGroupID: row.OrderGroupID, UserID: row.UserID, ShopID: row.ShopID, ShopName: row.ShopName, Status: row.Status, ItemAmount: row.ItemAmount, ShippingAmount: row.ShippingAmount, PayAmount: row.PayAmount, Currency: row.Currency, ItemCount: row.ItemCount, OrderGroupStatus: row.OrderGroupStatus, PaymentDeadlineAt: row.PaymentDeadlineAt, PaidAt: row.PaidAt, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt}
}

func toShopOrder(model ShopOrderModel) *domainorder.ShopOrder {
	return &domainorder.ShopOrder{ID: model.ID, OrderGroupID: model.OrderGroupID, UserID: model.UserID, ShopID: model.ShopID, ShopName: model.ShopName, Status: model.Status, ItemAmount: model.ItemAmount, ShippingAmount: model.ShippingAmount, PayAmount: model.PayAmount, Currency: model.Currency, ItemCount: model.ItemCount, PaidAt: model.PaidAt, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func toOrderItem(model OrderItemModel) *domainorder.Item {
	return &domainorder.Item{ID: model.ID, OrderGroupID: model.OrderGroupID, ShopOrderID: model.ShopOrderID, UserID: model.UserID, ShopID: model.ShopID, ProductID: model.ProductID, SKUID: model.SKUID, ProductTitle: model.ProductTitle, ProductSubTitle: model.ProductSubTitle, MainImageURL: model.MainImageURL, SKUName: model.SKUName, PriceAmount: model.PriceAmount, Currency: model.Currency, Quantity: model.Quantity, ItemAmount: model.ItemAmount, ReviewStatusSnapshot: model.ReviewStatusSnapshot, ProductSaleStatusSnapshot: model.ProductSaleStatusSnapshot, SKUSaleStatusSnapshot: model.SKUSaleStatusSnapshot, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func toPaymentInfo(model OrderGroupModel) domainorder.PaymentInfo {
	return domainorder.PaymentInfo{OrderGroupID: model.ID, UserID: model.UserID, Status: model.Status, TotalPayAmount: model.TotalPayAmount, Currency: model.Currency, PaymentDeadlineAt: model.PaymentDeadlineAt, PaidAt: model.PaidAt}
}

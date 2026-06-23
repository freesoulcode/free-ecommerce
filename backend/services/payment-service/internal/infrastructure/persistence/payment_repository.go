package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domainpayment "github.com/freesoulcode/free-ecommerce/backend/services/payment-service/internal/domain/payment"
	"gorm.io/gorm"
)

type PaymentOrderModel struct {
	ID           int64      `gorm:"column:id;primaryKey"`
	UserID       int64      `gorm:"column:user_id"`
	OrderGroupID int64      `gorm:"column:order_group_id"`
	Status       string     `gorm:"column:status"`
	Channel      string     `gorm:"column:channel"`
	PayAmount    int64      `gorm:"column:pay_amount"`
	Currency     string     `gorm:"column:currency"`
	ExpireAt     time.Time  `gorm:"column:expire_at"`
	PaidAt       *time.Time `gorm:"column:paid_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
}

func (PaymentOrderModel) TableName() string { return "payment_orders" }

type PaymentRepository struct{ db *gorm.DB }

func NewPaymentRepository(db *gorm.DB) *PaymentRepository { return &PaymentRepository{db: db} }

func (r *PaymentRepository) FindByOrderGroup(ctx context.Context, userID, orderGroupID int64) (*domainpayment.Order, error) {
	var model PaymentOrderModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND order_group_id = ?", userID, orderGroupID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("payment order not found")
		}
		return nil, fmt.Errorf("find payment order: %w", err)
	}
	return toDomainOrder(model), nil
}

func (r *PaymentRepository) Create(ctx context.Context, order *domainpayment.Order) error {
	model := toModel(order)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create payment order: %w", err)
	}
	return nil
}

func (r *PaymentRepository) MarkPaid(ctx context.Context, userID, orderGroupID int64, paidAt time.Time) (*domainpayment.Order, error) {
	var result *domainpayment.Order
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model PaymentOrderModel
		if err := tx.Where("user_id = ? AND order_group_id = ?", userID, orderGroupID).First(&model).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return appErrors.NotFound("payment order not found")
			}
			return fmt.Errorf("find payment order for mark paid: %w", err)
		}
		if model.Status == domainpayment.StatusPaid {
			result = toDomainOrder(model)
			return nil
		}
		if model.Status == domainpayment.StatusExpired {
			return appErrors.InvalidArgument("payment order expired")
		}
		if err := tx.Model(&PaymentOrderModel{}).Where("user_id = ? AND order_group_id = ?", userID, orderGroupID).Updates(map[string]any{"status": domainpayment.StatusPaid, "paid_at": paidAt, "updated_at": paidAt}).Error; err != nil {
			return fmt.Errorf("mark payment order paid: %w", err)
		}
		model.Status = domainpayment.StatusPaid
		model.PaidAt = &paidAt
		model.UpdatedAt = paidAt
		result = toDomainOrder(model)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PaymentRepository) MarkExpired(ctx context.Context, userID, orderGroupID int64, expiredAt time.Time) (*domainpayment.Order, error) {
	var result *domainpayment.Order
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model PaymentOrderModel
		if err := tx.Where("user_id = ? AND order_group_id = ?", userID, orderGroupID).First(&model).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return appErrors.NotFound("payment order not found")
			}
			return fmt.Errorf("find payment order for mark expired: %w", err)
		}
		if model.Status == domainpayment.StatusExpired {
			result = toDomainOrder(model)
			return nil
		}
		if model.Status == domainpayment.StatusPaid {
			result = toDomainOrder(model)
			return nil
		}
		if err := tx.Model(&PaymentOrderModel{}).Where("user_id = ? AND order_group_id = ?", userID, orderGroupID).Updates(map[string]any{"status": domainpayment.StatusExpired, "updated_at": expiredAt}).Error; err != nil {
			return fmt.Errorf("mark payment order expired: %w", err)
		}
		model.Status = domainpayment.StatusExpired
		model.UpdatedAt = expiredAt
		result = toDomainOrder(model)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func toModel(order *domainpayment.Order) PaymentOrderModel {
	return PaymentOrderModel{ID: order.ID, UserID: order.UserID, OrderGroupID: order.OrderGroupID, Status: order.Status, Channel: order.Channel, PayAmount: order.PayAmount, Currency: order.Currency, ExpireAt: order.ExpireAt, PaidAt: order.PaidAt, CreatedAt: order.CreatedAt, UpdatedAt: order.UpdatedAt}
}

func toDomainOrder(model PaymentOrderModel) *domainpayment.Order {
	return &domainpayment.Order{ID: model.ID, UserID: model.UserID, OrderGroupID: model.OrderGroupID, Status: model.Status, Channel: model.Channel, PayAmount: model.PayAmount, Currency: model.Currency, ExpireAt: model.ExpireAt, PaidAt: model.PaidAt, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

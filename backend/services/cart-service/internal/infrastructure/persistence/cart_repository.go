package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domaincart "github.com/freesoulcode/free-ecommerce/backend/services/cart-service/internal/domain/cart"
	"gorm.io/gorm"
)

type CartItemModel struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	UserID    int64     `gorm:"column:user_id"`
	SKUID     int64     `gorm:"column:sku_id"`
	Quantity  int32     `gorm:"column:quantity"`
	Selected  bool      `gorm:"column:selected"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (CartItemModel) TableName() string {
	return "cart_items"
}

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) CreateItem(ctx context.Context, item *domaincart.Item) error {
	model := toModel(item)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		if isDuplicateKeyError(err) {
			return appErrors.New(appErrors.Code("CART_ITEM_ALREADY_EXISTS"), "cart item already exists", 400)
		}
		return fmt.Errorf("create cart item: %w", err)
	}

	return nil
}

func (r *CartRepository) UpdateItem(ctx context.Context, item *domaincart.Item) error {
	result := r.db.WithContext(ctx).Model(&CartItemModel{}).
		Where("id = ? AND user_id = ?", item.ID, item.UserID).
		Updates(map[string]any{"quantity": item.Quantity, "selected": item.Selected, "updated_at": item.UpdatedAt})
	if result.Error != nil {
		return fmt.Errorf("update cart item: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return appErrors.NotFound("cart item not found")
	}

	return nil
}

func (r *CartRepository) DeleteItem(ctx context.Context, userID, itemID int64) error {
	result := r.db.WithContext(ctx).Delete(&CartItemModel{}, "id = ? AND user_id = ?", itemID, userID)
	if result.Error != nil {
		return fmt.Errorf("delete cart item: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return appErrors.NotFound("cart item not found")
	}

	return nil
}

func (r *CartRepository) FindItemByID(ctx context.Context, userID, itemID int64) (*domaincart.Item, error) {
	var model CartItemModel
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", itemID, userID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("cart item not found")
		}
		return nil, fmt.Errorf("find cart item by id: %w", err)
	}

	return toDomain(model), nil
}

func (r *CartRepository) FindItemByUserAndSKU(ctx context.Context, userID, skuID int64) (*domaincart.Item, error) {
	var model CartItemModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND sku_id = ?", userID, skuID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("cart item not found")
		}
		return nil, fmt.Errorf("find cart item by user and sku: %w", err)
	}

	return toDomain(model), nil
}

func (r *CartRepository) ListItemsByUserID(ctx context.Context, userID int64) ([]*domaincart.Item, error) {
	var models []CartItemModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC, id DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list cart items by user id: %w", err)
	}

	items := make([]*domaincart.Item, 0, len(models))
	for _, model := range models {
		items = append(items, toDomain(model))
	}

	return items, nil
}

func toModel(item *domaincart.Item) CartItemModel {
	return CartItemModel{ID: item.ID, UserID: item.UserID, SKUID: item.SKUID, Quantity: item.Quantity, Selected: item.Selected, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt}
}

func toDomain(model CartItemModel) *domaincart.Item {
	return &domaincart.Item{ID: model.ID, UserID: model.UserID, SKUID: model.SKUID, Quantity: model.Quantity, Selected: model.Selected, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(strings.ToLower(err.Error()), "duplicate")
}

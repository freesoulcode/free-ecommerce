package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domainuser "github.com/freesoulcode/free-ecommerce/backend/services/user-service/internal/domain/user"
	"gorm.io/gorm"
)

type UserModel struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	Email         *string   `gorm:"column:email"`
	Phone         *string   `gorm:"column:phone"`
	Nickname      string    `gorm:"column:nickname"`
	Status        string    `gorm:"column:status"`
	EmailVerified bool      `gorm:"column:email_verified"`
	PhoneVerified bool      `gorm:"column:phone_verified"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (UserModel) TableName() string {
	return "users"
}

type AddressModel struct {
	ID            int64     `gorm:"column:id;primaryKey"`
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
	IsDefault     bool      `gorm:"column:is_default"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (AddressModel) TableName() string {
	return "user_addresses"
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domainuser.User) error {
	model := UserModel{
		ID:            user.ID,
		Email:         nullableString(user.Email),
		Phone:         nullableString(user.Phone),
		Nickname:      user.Nickname,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		if isDuplicateKeyError(err) {
			return appErrors.New(appErrors.Code("USER_ALREADY_EXISTS"), "user already exists", 400)
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *UserRepository) DeleteByID(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&UserModel{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("delete user by id: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return appErrors.NotFound("user not found")
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domainuser.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("user not found")
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &domainuser.User{
		ID:            model.ID,
		Email:         derefString(model.Email),
		Phone:         derefString(model.Phone),
		Nickname:      model.Nickname,
		Status:        model.Status,
		EmailVerified: model.EmailVerified,
		PhoneVerified: model.PhoneVerified,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if strings.TrimSpace(email) == "" {
		return false, nil
	}

	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("count users by email: %w", err)
	}

	return count > 0, nil
}

func (r *UserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	if strings.TrimSpace(phone) == "" {
		return false, nil
	}

	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("phone = ?", phone).Count(&count).Error; err != nil {
		return false, fmt.Errorf("count users by phone: %w", err)
	}

	return count > 0, nil
}

func (r *UserRepository) CreateAddress(ctx context.Context, address *domainuser.Address) error {
	model := toAddressModel(address)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create address: %w", err)
	}

	return nil
}

func (r *UserRepository) UpdateAddress(ctx context.Context, address *domainuser.Address) error {
	model := toAddressModel(address)
	result := r.db.WithContext(ctx).Model(&AddressModel{}).
		Where("id = ? AND user_id = ?", address.ID, address.UserID).
		Updates(map[string]any{
			"receiver_name":   model.ReceiverName,
			"receiver_phone":  model.ReceiverPhone,
			"country_code":    model.CountryCode,
			"province":        model.Province,
			"city":            model.City,
			"district":        model.District,
			"address_line1":   model.AddressLine1,
			"address_line2":   model.AddressLine2,
			"postal_code":     model.PostalCode,
			"tag":             model.Tag,
			"is_default":      model.IsDefault,
			"updated_at":      model.UpdatedAt,
		})
	if result.Error != nil {
		return fmt.Errorf("update address: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return appErrors.NotFound("address not found")
	}

	return nil
}

func (r *UserRepository) DeleteAddress(ctx context.Context, userID, addressID int64) error {
	result := r.db.WithContext(ctx).Delete(&AddressModel{}, "id = ? AND user_id = ?", addressID, userID)
	if result.Error != nil {
		return fmt.Errorf("delete address: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return appErrors.NotFound("address not found")
	}

	return nil
}

func (r *UserRepository) FindAddressByID(ctx context.Context, userID, addressID int64) (*domainuser.Address, error) {
	var model AddressModel
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", addressID, userID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("address not found")
		}
		return nil, fmt.Errorf("find address by id: %w", err)
	}

	return toDomainAddress(model), nil
}

func (r *UserRepository) ListAddressesByUserID(ctx context.Context, userID int64) ([]*domainuser.Address, error) {
	var models []AddressModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, updated_at DESC, id DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list addresses by user id: %w", err)
	}

	addresses := make([]*domainuser.Address, 0, len(models))
	for _, model := range models {
		addresses = append(addresses, toDomainAddress(model))
	}

	return addresses, nil
}

func (r *UserRepository) CountAddressesByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&AddressModel{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count addresses by user id: %w", err)
	}

	return count, nil
}

func (r *UserRepository) SetDefaultAddress(ctx context.Context, userID, addressID int64, updatedAt time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&AddressModel{}).Where("user_id = ?", userID).Updates(map[string]any{
			"is_default": false,
			"updated_at": updatedAt,
		}).Error; err != nil {
			return fmt.Errorf("clear default address: %w", err)
		}

		result := tx.Model(&AddressModel{}).Where("id = ? AND user_id = ?", addressID, userID).Updates(map[string]any{
			"is_default": true,
			"updated_at": updatedAt,
		})
		if result.Error != nil {
			return fmt.Errorf("set default address: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return appErrors.NotFound("address not found")
		}

		return nil
	})
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(strings.ToLower(err.Error()), "duplicate")
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

func toAddressModel(address *domainuser.Address) AddressModel {
	return AddressModel{
		ID:            address.ID,
		UserID:        address.UserID,
		ReceiverName:  address.ReceiverName,
		ReceiverPhone: address.ReceiverPhone,
		CountryCode:   address.CountryCode,
		Province:      address.Province,
		City:          address.City,
		District:      address.District,
		AddressLine1:  address.AddressLine1,
		AddressLine2:  nullableString(address.AddressLine2),
		PostalCode:    nullableString(address.PostalCode),
		Tag:           nullableString(address.Tag),
		IsDefault:     address.IsDefault,
		CreatedAt:     address.CreatedAt,
		UpdatedAt:     address.UpdatedAt,
	}
}

func toDomainAddress(model AddressModel) *domainuser.Address {
	return &domainuser.Address{
		ID:            model.ID,
		UserID:        model.UserID,
		ReceiverName:  model.ReceiverName,
		ReceiverPhone: model.ReceiverPhone,
		CountryCode:   model.CountryCode,
		Province:      model.Province,
		City:          model.City,
		District:      model.District,
		AddressLine1:  model.AddressLine1,
		AddressLine2:  derefString(model.AddressLine2),
		PostalCode:    derefString(model.PostalCode),
		Tag:           derefString(model.Tag),
		IsDefault:     model.IsDefault,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}

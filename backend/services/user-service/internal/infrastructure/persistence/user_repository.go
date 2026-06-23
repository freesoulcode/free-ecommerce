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

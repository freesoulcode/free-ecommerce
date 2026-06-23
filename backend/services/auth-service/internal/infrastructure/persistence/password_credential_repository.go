package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domaincredential "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/domain/credential"
	"gorm.io/gorm"
)

type PasswordCredentialModel struct {
	UserID       int64     `gorm:"column:user_id;primaryKey"`
	Email        *string   `gorm:"column:email"`
	Phone        *string   `gorm:"column:phone"`
	PasswordHash string    `gorm:"column:password_hash"`
	PasswordAlgo string    `gorm:"column:password_algo"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (PasswordCredentialModel) TableName() string { return "password_credentials" }

type RefreshSessionModel struct {
	ID                  int64      `gorm:"column:id;primaryKey"`
	UserID              int64      `gorm:"column:user_id"`
	TokenHash           string     `gorm:"column:token_hash"`
	DeviceID            *string    `gorm:"column:device_id"`
	UserAgent           *string    `gorm:"column:user_agent"`
	ClientIP            *string    `gorm:"column:client_ip"`
	ExpiresAt           time.Time  `gorm:"column:expires_at"`
	RevokedAt           *time.Time `gorm:"column:revoked_at"`
	ReplacedBySessionID *int64     `gorm:"column:replaced_by_session_id"`
	CreatedAt           time.Time  `gorm:"column:created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at"`
}

func (RefreshSessionModel) TableName() string { return "refresh_sessions" }

type PasswordCredentialRepository struct{ db *gorm.DB }

func NewPasswordCredentialRepository(db *gorm.DB) *PasswordCredentialRepository {
	return &PasswordCredentialRepository{db: db}
}

func (r *PasswordCredentialRepository) CreatePasswordCredential(ctx context.Context, credential *domaincredential.PasswordCredential) error {
	model := PasswordCredentialModel{
		UserID:       credential.UserID,
		Email:        nullableString(credential.Email),
		Phone:        nullableString(credential.Phone),
		PasswordHash: credential.PasswordHash,
		PasswordAlgo: credential.PasswordAlgo,
		CreatedAt:    credential.CreatedAt,
		UpdatedAt:    credential.UpdatedAt,
	}
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create password credential: %w", err)
	}
	return nil
}

func (r *PasswordCredentialRepository) ExistsByUserID(ctx context.Context, userID int64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&PasswordCredentialModel{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("count credentials by user_id: %w", err)
	}
	return count > 0, nil
}

func (r *PasswordCredentialRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if strings.TrimSpace(email) == "" {
		return false, nil
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&PasswordCredentialModel{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("count credentials by email: %w", err)
	}
	return count > 0, nil
}

func (r *PasswordCredentialRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	if strings.TrimSpace(phone) == "" {
		return false, nil
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&PasswordCredentialModel{}).Where("phone = ?", phone).Count(&count).Error; err != nil {
		return false, fmt.Errorf("count credentials by phone: %w", err)
	}
	return count > 0, nil
}

func (r *PasswordCredentialRepository) FindByEmail(ctx context.Context, email string) (*domaincredential.PasswordCredential, error) {
	if strings.TrimSpace(email) == "" {
		return nil, appErrors.InvalidArgument("email is required")
	}

	var model PasswordCredentialModel
	if err := r.db.WithContext(ctx).Where("email = ?", strings.TrimSpace(strings.ToLower(email))).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.Unauthorized("invalid email or password")
		}
		return nil, fmt.Errorf("find credential by email: %w", err)
	}

	return toDomainPasswordCredential(model), nil
}

func (r *PasswordCredentialRepository) FindByUserID(ctx context.Context, userID int64) (*domaincredential.PasswordCredential, error) {
	var model PasswordCredentialModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.Unauthorized("credential not found")
		}
		return nil, fmt.Errorf("find credential by user_id: %w", err)
	}

	return toDomainPasswordCredential(model), nil
}

func (r *PasswordCredentialRepository) CreateRefreshSession(ctx context.Context, session *domaincredential.RefreshSession) error {
	model := RefreshSessionModel{
		ID:                  session.ID,
		UserID:              session.UserID,
		TokenHash:           session.TokenHash,
		DeviceID:            nullableString(session.DeviceID),
		UserAgent:           nullableString(session.UserAgent),
		ClientIP:            nullableString(session.ClientIP),
		ExpiresAt:           session.ExpiresAt,
		RevokedAt:           session.RevokedAt,
		ReplacedBySessionID: session.ReplacedBySessionID,
		CreatedAt:           session.CreatedAt,
		UpdatedAt:           session.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("create refresh session: %w", err)
	}

	return nil
}

func (r *PasswordCredentialRepository) FindRefreshSessionByTokenHash(ctx context.Context, tokenHash string) (*domaincredential.RefreshSession, error) {
	var model RefreshSessionModel
	if err := r.db.WithContext(ctx).Where("token_hash = ?", strings.TrimSpace(tokenHash)).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.Unauthorized("refresh token is invalid")
		}
		return nil, fmt.Errorf("find refresh session by token hash: %w", err)
	}

	return toDomainRefreshSession(model), nil
}

func (r *PasswordCredentialRepository) RotateRefreshSession(ctx context.Context, oldSessionID int64, revokedAt time.Time, replacedBySessionID int64, newSession *domaincredential.RefreshSession) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&RefreshSessionModel{}).
			Where("id = ? AND revoked_at IS NULL", oldSessionID).
			Updates(map[string]any{
				"revoked_at":              revokedAt,
				"replaced_by_session_id":  replacedBySessionID,
				"updated_at":              revokedAt,
			})
		if result.Error != nil {
			return fmt.Errorf("revoke old refresh session: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return appErrors.Unauthorized("refresh token is invalid")
		}

		model := RefreshSessionModel{
			ID:                  newSession.ID,
			UserID:              newSession.UserID,
			TokenHash:           newSession.TokenHash,
			DeviceID:            nullableString(newSession.DeviceID),
			UserAgent:           nullableString(newSession.UserAgent),
			ClientIP:            nullableString(newSession.ClientIP),
			ExpiresAt:           newSession.ExpiresAt,
			RevokedAt:           newSession.RevokedAt,
			ReplacedBySessionID: newSession.ReplacedBySessionID,
			CreatedAt:           newSession.CreatedAt,
			UpdatedAt:           newSession.UpdatedAt,
		}
		if err := tx.Create(&model).Error; err != nil {
			return fmt.Errorf("create rotated refresh session: %w", err)
		}

		return nil
	})
}

func (r *PasswordCredentialRepository) RevokeRefreshSession(ctx context.Context, sessionID int64, revokedAt time.Time) error {
	result := r.db.WithContext(ctx).Model(&RefreshSessionModel{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Updates(map[string]any{
			"revoked_at": revokedAt,
			"updated_at": revokedAt,
		})
	if result.Error != nil {
		return fmt.Errorf("revoke refresh session: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil
	}

	return nil
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

func toDomainPasswordCredential(model PasswordCredentialModel) *domaincredential.PasswordCredential {
	return &domaincredential.PasswordCredential{
		UserID:       model.UserID,
		Email:        derefString(model.Email),
		Phone:        derefString(model.Phone),
		PasswordHash: model.PasswordHash,
		PasswordAlgo: model.PasswordAlgo,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

func toDomainRefreshSession(model RefreshSessionModel) *domaincredential.RefreshSession {
	return &domaincredential.RefreshSession{
		ID:                  model.ID,
		UserID:              model.UserID,
		TokenHash:           model.TokenHash,
		DeviceID:            derefString(model.DeviceID),
		UserAgent:           derefString(model.UserAgent),
		ClientIP:            derefString(model.ClientIP),
		ExpiresAt:           model.ExpiresAt,
		RevokedAt:           model.RevokedAt,
		ReplacedBySessionID: model.ReplacedBySessionID,
		CreatedAt:           model.CreatedAt,
		UpdatedAt:           model.UpdatedAt,
	}
}

func toDuplicateError(message string) error {
	return appErrors.New(appErrors.Code("AUTH_CREDENTIAL_ALREADY_EXISTS"), message, 400)
}

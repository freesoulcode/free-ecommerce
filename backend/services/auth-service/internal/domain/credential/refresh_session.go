package credential

import (
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
)

type RefreshSession struct {
	ID                  int64
	UserID              int64
	TokenHash           string
	DeviceID            string
	UserAgent           string
	ClientIP            string
	ExpiresAt           time.Time
	RevokedAt           *time.Time
	ReplacedBySessionID *int64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func NewRefreshSession(id, userID int64, tokenHash, deviceID, userAgent, clientIP string, expiresAt, now time.Time) (*RefreshSession, error) {
	tokenHash = strings.TrimSpace(tokenHash)
	deviceID = strings.TrimSpace(deviceID)
	userAgent = strings.TrimSpace(userAgent)
	clientIP = strings.TrimSpace(clientIP)
	now = now.UTC()
	expiresAt = expiresAt.UTC()

	if id <= 0 {
		return nil, appErrors.New(appErrors.Code("AUTH_REFRESH_SESSION_ID_REQUIRED"), "refresh session id is required", 500)
	}
	if userID <= 0 {
		return nil, appErrors.New(appErrors.Code("AUTH_USER_ID_REQUIRED"), "user_id is required", 400)
	}
	if tokenHash == "" {
		return nil, appErrors.New(appErrors.Code("AUTH_REFRESH_TOKEN_HASH_REQUIRED"), "refresh token hash is required", 500)
	}
	if expiresAt.Before(now) || expiresAt.Equal(now) {
		return nil, appErrors.New(appErrors.Code("AUTH_REFRESH_TOKEN_EXPIRES_AT_INVALID"), "refresh token expires_at is invalid", 500)
	}

	return &RefreshSession{
		ID:        id,
		UserID:    userID,
		TokenHash: tokenHash,
		DeviceID:  deviceID,
		UserAgent: userAgent,
		ClientIP:  clientIP,
		ExpiresAt: expiresAt,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *RefreshSession) IsExpired(now time.Time) bool {
	return !s.ExpiresAt.After(now.UTC())
}

func (s *RefreshSession) IsRevoked() bool {
	return s.RevokedAt != nil
}

func (s *RefreshSession) CanRefresh(now time.Time) bool {
	return !s.IsRevoked() && !s.IsExpired(now)
}

func (s *RefreshSession) Revoke(now time.Time, replacedBySessionID *int64) {
	if s == nil {
		return
	}
	revokedAt := now.UTC()
	s.RevokedAt = &revokedAt
	s.ReplacedBySessionID = replacedBySessionID
	s.UpdatedAt = revokedAt
}

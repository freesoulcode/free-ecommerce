package credential

import (
	"context"
	"time"
)

type Repository interface {
	CreatePasswordCredential(ctx context.Context, credential *PasswordCredential) error
	ExistsByUserID(ctx context.Context, userID int64) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	FindByEmail(ctx context.Context, email string) (*PasswordCredential, error)
	FindByUserID(ctx context.Context, userID int64) (*PasswordCredential, error)
}

type RefreshSessionRepository interface {
	CreateRefreshSession(ctx context.Context, session *RefreshSession) error
	FindRefreshSessionByTokenHash(ctx context.Context, tokenHash string) (*RefreshSession, error)
	RotateRefreshSession(ctx context.Context, oldSessionID int64, revokedAt time.Time, replacedBySessionID int64, newSession *RefreshSession) error
	RevokeRefreshSession(ctx context.Context, sessionID int64, revokedAt time.Time) error
}

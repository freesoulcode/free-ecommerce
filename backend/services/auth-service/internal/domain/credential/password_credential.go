package credential

import (
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
)

const PasswordAlgoArgon2id = "argon2id"

type PasswordCredential struct {
	UserID       int64
	Email        string
	Phone        string
	PasswordHash string
	PasswordAlgo string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewPasswordCredential(userID int64, email, phone, passwordHash string, now time.Time) (*PasswordCredential, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	phone = strings.TrimSpace(phone)
	passwordHash = strings.TrimSpace(passwordHash)

	if userID <= 0 {
		return nil, appErrors.New(appErrors.Code("AUTH_USER_ID_REQUIRED"), "user_id is required", 400)
	}
	if email == "" && phone == "" {
		return nil, appErrors.New(appErrors.Code("AUTH_IDENTIFIER_REQUIRED"), "email or phone is required", 400)
	}
	if passwordHash == "" {
		return nil, appErrors.New(appErrors.Code("AUTH_PASSWORD_HASH_REQUIRED"), "password hash is required", 500)
	}

	return &PasswordCredential{
		UserID:       userID,
		Email:        email,
		Phone:        phone,
		PasswordHash: passwordHash,
		PasswordAlgo: PasswordAlgoArgon2id,
		CreatedAt:    now.UTC(),
		UpdatedAt:    now.UTC(),
	}, nil
}

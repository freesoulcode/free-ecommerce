package user

import (
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
)

const (
	StatusActive = "active"
)

type User struct {
	ID            int64
	Email         string
	Phone         string
	Nickname      string
	Status        string
	EmailVerified bool
	PhoneVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func New(id int64, email, phone, nickname string, now time.Time) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	phone = strings.TrimSpace(phone)
	nickname = strings.TrimSpace(nickname)

	if email == "" && phone == "" {
		return nil, appErrors.New(appErrors.Code("USER_IDENTIFIER_REQUIRED"), "email or phone is required", 400)
	}
	if nickname == "" {
		return nil, appErrors.New(appErrors.Code("USER_NICKNAME_REQUIRED"), "nickname is required", 400)
	}

	return &User{
		ID:            id,
		Email:         email,
		Phone:         phone,
		Nickname:      nickname,
		Status:        StatusActive,
		EmailVerified: false,
		PhoneVerified: false,
		CreatedAt:     now.UTC(),
		UpdatedAt:     now.UTC(),
	}, nil
}

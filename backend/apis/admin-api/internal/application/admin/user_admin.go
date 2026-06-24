package admin

import (
	"context"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
)

type User struct {
	ID            int64
	Email         string
	Phone         string
	Nickname      string
	Status        string
	EmailVerified bool
	PhoneVerified bool
}

type Address struct {
	ID            int64
	UserID        int64
	ReceiverName  string
	ReceiverPhone string
	CountryCode   string
	Province      string
	City          string
	District      string
	AddressLine1  string
	AddressLine2  string
	PostalCode    string
	Tag           string
	IsDefault     bool
	CreatedAt     int64
	UpdatedAt     int64
}

type UserReader interface {
	GetUser(ctx context.Context, id int64) (*User, error)
	ListAddresses(ctx context.Context, userID int64) ([]*Address, error)
}

type UserAdminService struct {
	reader UserReader
}

func NewUserAdminService(reader UserReader) *UserAdminService {
	return &UserAdminService{reader: reader}
}

func (s *UserAdminService) Get(ctx context.Context, userID int64) (*User, error) {
	if s.reader == nil {
		return nil, appErrors.Internal("user service is not configured")
	}

	return s.reader.GetUser(ctx, userID)
}

func (s *UserAdminService) ListAddresses(ctx context.Context, userID int64) ([]*Address, error) {
	if s.reader == nil {
		return nil, appErrors.Internal("user service is not configured")
	}

	return s.reader.ListAddresses(ctx, userID)
}

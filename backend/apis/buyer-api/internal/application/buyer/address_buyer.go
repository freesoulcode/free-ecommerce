package buyer

import (
	"context"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
)

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

type AddressProfileService interface {
	CreateAddress(ctx context.Context, input CreateAddressInput) (*Address, error)
	UpdateAddress(ctx context.Context, input UpdateAddressInput) (*Address, error)
	DeleteAddress(ctx context.Context, input DeleteAddressInput) error
	ListAddresses(ctx context.Context, userID int64) ([]*Address, error)
}

type CreateAddressInput struct {
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
}

type UpdateAddressInput struct {
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
}

type DeleteAddressInput struct {
	ID     int64
	UserID int64
}

type AddressBuyerService struct {
	profileService AddressProfileService
}

func NewAddressBuyerService(profileService AddressProfileService) *AddressBuyerService {
	return &AddressBuyerService{profileService: profileService}
}

func (s *AddressBuyerService) Create(ctx context.Context, input CreateAddressInput) (*Address, error) {
	if s.profileService == nil {
		return nil, appErrors.Internal("user service is not configured")
	}

	return s.profileService.CreateAddress(ctx, input)
}

func (s *AddressBuyerService) Update(ctx context.Context, input UpdateAddressInput) (*Address, error) {
	if s.profileService == nil {
		return nil, appErrors.Internal("user service is not configured")
	}

	return s.profileService.UpdateAddress(ctx, input)
}

func (s *AddressBuyerService) Delete(ctx context.Context, input DeleteAddressInput) error {
	if s.profileService == nil {
		return appErrors.Internal("user service is not configured")
	}

	return s.profileService.DeleteAddress(ctx, input)
}

func (s *AddressBuyerService) List(ctx context.Context, userID int64) ([]*Address, error) {
	if s.profileService == nil {
		return nil, appErrors.Internal("user service is not configured")
	}

	return s.profileService.ListAddresses(ctx, userID)
}

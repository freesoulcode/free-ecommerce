package user

import (
	"context"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domainuser "github.com/freesoulcode/free-ecommerce/backend/services/user-service/internal/domain/user"
)

type Address = domainuser.Address

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

type CreateAddressService struct {
	repo        domainuser.Repository
	idGenerator IDGenerator
	now         func() time.Time
}

type UpdateAddressService struct {
	repo domainuser.Repository
	now  func() time.Time
}

type DeleteAddressService struct {
	repo domainuser.Repository
	now  func() time.Time
}

type ListAddressesService struct {
	repo domainuser.Repository
}

func NewCreateAddressService(repo domainuser.Repository, idGenerator IDGenerator, now func() time.Time) *CreateAddressService {
	if now == nil {
		now = time.Now
	}

	return &CreateAddressService{repo: repo, idGenerator: idGenerator, now: now}
}

func NewUpdateAddressService(repo domainuser.Repository, now func() time.Time) *UpdateAddressService {
	if now == nil {
		now = time.Now
	}

	return &UpdateAddressService{repo: repo, now: now}
}

func NewDeleteAddressService(repo domainuser.Repository, now func() time.Time) *DeleteAddressService {
	if now == nil {
		now = time.Now
	}

	return &DeleteAddressService{repo: repo, now: now}
}

func NewListAddressesService(repo domainuser.Repository) *ListAddressesService {
	return &ListAddressesService{repo: repo}
}

func (s *CreateAddressService) Execute(ctx context.Context, input CreateAddressInput) (*domainuser.Address, error) {
	if input.UserID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}
	if _, err := s.repo.FindByID(ctx, input.UserID); err != nil {
		return nil, err
	}

	count, err := s.repo.CountAddressesByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	addressID, err := s.idGenerator.NextID()
	if err != nil {
		return nil, appErrors.Internal("generate address id failed")
	}

	now := s.now().UTC()
	address, err := domainuser.NewAddress(
		addressID,
		input.UserID,
		input.ReceiverName,
		input.ReceiverPhone,
		input.CountryCode,
		input.Province,
		input.City,
		input.District,
		input.AddressLine1,
		input.AddressLine2,
		input.PostalCode,
		input.Tag,
		false,
		now,
	)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateAddress(ctx, address); err != nil {
		return nil, err
	}

	if input.IsDefault || count == 0 {
		if err := s.repo.SetDefaultAddress(ctx, input.UserID, address.ID, now); err != nil {
			return nil, err
		}
		address.IsDefault = true
		address.UpdatedAt = now
	}

	return address, nil
}

func (s *UpdateAddressService) Execute(ctx context.Context, input UpdateAddressInput) (*domainuser.Address, error) {
	if input.ID <= 0 {
		return nil, appErrors.InvalidArgument("address id is required")
	}
	if input.UserID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}

	address, err := s.repo.FindAddressByID(ctx, input.UserID, input.ID)
	if err != nil {
		return nil, err
	}

	now := s.now().UTC()
	keepDefault := address.IsDefault || input.IsDefault
	if err := address.Update(
		input.ReceiverName,
		input.ReceiverPhone,
		input.CountryCode,
		input.Province,
		input.City,
		input.District,
		input.AddressLine1,
		input.AddressLine2,
		input.PostalCode,
		input.Tag,
		keepDefault,
		now,
	); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateAddress(ctx, address); err != nil {
		return nil, err
	}

	if keepDefault {
		if err := s.repo.SetDefaultAddress(ctx, input.UserID, address.ID, now); err != nil {
			return nil, err
		}
		address.IsDefault = true
		address.UpdatedAt = now
	}

	return address, nil
}

func (s *DeleteAddressService) Execute(ctx context.Context, input DeleteAddressInput) error {
	if input.ID <= 0 {
		return appErrors.InvalidArgument("address id is required")
	}
	if input.UserID <= 0 {
		return appErrors.InvalidArgument("user id is required")
	}

	address, err := s.repo.FindAddressByID(ctx, input.UserID, input.ID)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteAddress(ctx, input.UserID, input.ID); err != nil {
		return err
	}

	if !address.IsDefault {
		return nil
	}

	remaining, err := s.repo.ListAddressesByUserID(ctx, input.UserID)
	if err != nil {
		return err
	}
	if len(remaining) == 0 {
		return nil
	}

	return s.repo.SetDefaultAddress(ctx, input.UserID, remaining[0].ID, s.now().UTC())
}

func (s *ListAddressesService) Execute(ctx context.Context, userID int64) ([]*domainuser.Address, error) {
	if userID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}
	if _, err := s.repo.FindByID(ctx, userID); err != nil {
		return nil, err
	}

	return s.repo.ListAddressesByUserID(ctx, userID)
}

package usergrpc

import (
	applicationbuyer "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/application/buyer"
	userv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/user/v1"
)

func toAppAddress(address *userv1.Address) *applicationbuyer.Address {
	if address == nil {
		return nil
	}

	return &applicationbuyer.Address{
		ID:            address.GetId(),
		UserID:        address.GetUserId(),
		ReceiverName:  address.GetReceiverName(),
		ReceiverPhone: address.GetReceiverPhone(),
		CountryCode:   address.GetCountryCode(),
		Province:      address.GetProvince(),
		City:          address.GetCity(),
		District:      address.GetDistrict(),
		AddressLine1:  address.GetAddressLine1(),
		AddressLine2:  address.GetAddressLine2(),
		PostalCode:    address.GetPostalCode(),
		Tag:           address.GetTag(),
		IsDefault:     address.GetIsDefault(),
		CreatedAt:     address.GetCreatedAt(),
		UpdatedAt:     address.GetUpdatedAt(),
	}
}

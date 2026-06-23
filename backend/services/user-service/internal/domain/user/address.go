package user

import (
	"strings"
	"time"

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
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewAddress(
	id int64,
	userID int64,
	receiverName string,
	receiverPhone string,
	countryCode string,
	province string,
	city string,
	district string,
	addressLine1 string,
	addressLine2 string,
	postalCode string,
	tag string,
	isDefault bool,
	now time.Time,
) (*Address, error) {
	if id <= 0 {
		return nil, appErrors.InvalidArgument("address id is required")
	}
	if userID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}

	address := &Address{
		ID:            id,
		UserID:        userID,
		ReceiverName:  strings.TrimSpace(receiverName),
		ReceiverPhone: strings.TrimSpace(receiverPhone),
		CountryCode:   normalizeCountryCode(countryCode),
		Province:      strings.TrimSpace(province),
		City:          strings.TrimSpace(city),
		District:      strings.TrimSpace(district),
		AddressLine1:  strings.TrimSpace(addressLine1),
		AddressLine2:  strings.TrimSpace(addressLine2),
		PostalCode:    strings.TrimSpace(postalCode),
		Tag:           strings.TrimSpace(tag),
		IsDefault:     isDefault,
		CreatedAt:     now.UTC(),
		UpdatedAt:     now.UTC(),
	}

	if err := address.Validate(); err != nil {
		return nil, err
	}

	return address, nil
}

func (a *Address) Validate() error {
	if a == nil {
		return appErrors.InvalidArgument("address is required")
	}
	if a.UserID <= 0 {
		return appErrors.InvalidArgument("user id is required")
	}
	if a.ReceiverName == "" {
		return appErrors.New(appErrors.Code("USER_ADDRESS_RECEIVER_NAME_REQUIRED"), "receiver name is required", 400)
	}
	if a.ReceiverPhone == "" {
		return appErrors.New(appErrors.Code("USER_ADDRESS_RECEIVER_PHONE_REQUIRED"), "receiver phone is required", 400)
	}
	if a.CountryCode == "" {
		return appErrors.New(appErrors.Code("USER_ADDRESS_COUNTRY_CODE_REQUIRED"), "country code is required", 400)
	}
	if a.Province == "" {
		return appErrors.New(appErrors.Code("USER_ADDRESS_PROVINCE_REQUIRED"), "province is required", 400)
	}
	if a.City == "" {
		return appErrors.New(appErrors.Code("USER_ADDRESS_CITY_REQUIRED"), "city is required", 400)
	}
	if a.District == "" {
		return appErrors.New(appErrors.Code("USER_ADDRESS_DISTRICT_REQUIRED"), "district is required", 400)
	}
	if a.AddressLine1 == "" {
		return appErrors.New(appErrors.Code("USER_ADDRESS_LINE1_REQUIRED"), "address line1 is required", 400)
	}

	return nil
}

func (a *Address) Update(
	receiverName string,
	receiverPhone string,
	countryCode string,
	province string,
	city string,
	district string,
	addressLine1 string,
	addressLine2 string,
	postalCode string,
	tag string,
	isDefault bool,
	now time.Time,
) error {
	a.ReceiverName = strings.TrimSpace(receiverName)
	a.ReceiverPhone = strings.TrimSpace(receiverPhone)
	a.CountryCode = normalizeCountryCode(countryCode)
	a.Province = strings.TrimSpace(province)
	a.City = strings.TrimSpace(city)
	a.District = strings.TrimSpace(district)
	a.AddressLine1 = strings.TrimSpace(addressLine1)
	a.AddressLine2 = strings.TrimSpace(addressLine2)
	a.PostalCode = strings.TrimSpace(postalCode)
	a.Tag = strings.TrimSpace(tag)
	a.IsDefault = isDefault
	a.UpdatedAt = now.UTC()

	return a.Validate()
}

func normalizeCountryCode(countryCode string) string {
	return strings.ToUpper(strings.TrimSpace(countryCode))
}

package grpc

import (
	"context"
	stderrors "errors"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	applicationuser "github.com/freesoulcode/free-ecommerce/backend/services/user-service/internal/application/user"
	userv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	userv1.UnimplementedUserServiceServer

	createUserService    *applicationuser.CreateUserService
	deleteUserService    *applicationuser.DeleteUserService
	getUserService       *applicationuser.GetUserService
	getAddressService    *applicationuser.GetAddressService
	createAddressService *applicationuser.CreateAddressService
	updateAddressService *applicationuser.UpdateAddressService
	deleteAddressService *applicationuser.DeleteAddressService
	listAddressesService *applicationuser.ListAddressesService
}

func NewUserServiceServer(
	createUserService *applicationuser.CreateUserService,
	deleteUserService *applicationuser.DeleteUserService,
	getUserService *applicationuser.GetUserService,
	getAddressService *applicationuser.GetAddressService,
	createAddressService *applicationuser.CreateAddressService,
	updateAddressService *applicationuser.UpdateAddressService,
	deleteAddressService *applicationuser.DeleteAddressService,
	listAddressesService *applicationuser.ListAddressesService,
) *UserServiceServer {
	return &UserServiceServer{
		createUserService:    createUserService,
		deleteUserService:    deleteUserService,
		getUserService:       getUserService,
		getAddressService:    getAddressService,
		createAddressService: createAddressService,
		updateAddressService: updateAddressService,
		deleteAddressService: deleteAddressService,
		listAddressesService: listAddressesService,
	}
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	user, err := s.createUserService.Execute(ctx, applicationuser.CreateUserInput{
		Email:    req.GetEmail(),
		Phone:    req.GetPhone(),
		Nickname: req.GetNickname(),
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userv1.CreateUserResponse{
		Id:            user.ID,
		Email:         user.Email,
		Phone:         user.Phone,
		Nickname:      user.Nickname,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
	}, nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if err := s.deleteUserService.Execute(ctx, req.GetId()); err != nil {
		return nil, toGRPCError(err)
	}

	return &userv1.DeleteUserResponse{}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	user, err := s.getUserService.Execute(ctx, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userv1.GetUserResponse{
		Id:            user.ID,
		Email:         user.Email,
		Phone:         user.Phone,
		Nickname:      user.Nickname,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
	}, nil
}

func (s *UserServiceServer) CreateAddress(ctx context.Context, req *userv1.CreateAddressRequest) (*userv1.CreateAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	address, err := s.createAddressService.Execute(ctx, applicationuser.CreateAddressInput{
		UserID:        req.GetUserId(),
		ReceiverName:  req.GetReceiverName(),
		ReceiverPhone: req.GetReceiverPhone(),
		CountryCode:   req.GetCountryCode(),
		Province:      req.GetProvince(),
		City:          req.GetCity(),
		District:      req.GetDistrict(),
		AddressLine1:  req.GetAddressLine1(),
		AddressLine2:  req.GetAddressLine2(),
		PostalCode:    req.GetPostalCode(),
		Tag:           req.GetTag(),
		IsDefault:     req.GetIsDefault(),
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userv1.CreateAddressResponse{Address: toAddressPB(address)}, nil
}

func (s *UserServiceServer) GetAddress(ctx context.Context, req *userv1.GetAddressRequest) (*userv1.GetAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	address, err := s.getAddressService.Execute(ctx, req.GetUserId(), req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userv1.GetAddressResponse{Address: toAddressPB(address)}, nil
}

func (s *UserServiceServer) UpdateAddress(ctx context.Context, req *userv1.UpdateAddressRequest) (*userv1.UpdateAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	address, err := s.updateAddressService.Execute(ctx, applicationuser.UpdateAddressInput{
		ID:            req.GetId(),
		UserID:        req.GetUserId(),
		ReceiverName:  req.GetReceiverName(),
		ReceiverPhone: req.GetReceiverPhone(),
		CountryCode:   req.GetCountryCode(),
		Province:      req.GetProvince(),
		City:          req.GetCity(),
		District:      req.GetDistrict(),
		AddressLine1:  req.GetAddressLine1(),
		AddressLine2:  req.GetAddressLine2(),
		PostalCode:    req.GetPostalCode(),
		Tag:           req.GetTag(),
		IsDefault:     req.GetIsDefault(),
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userv1.UpdateAddressResponse{Address: toAddressPB(address)}, nil
}

func (s *UserServiceServer) DeleteAddress(ctx context.Context, req *userv1.DeleteAddressRequest) (*userv1.DeleteAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if err := s.deleteAddressService.Execute(ctx, applicationuser.DeleteAddressInput{ID: req.GetId(), UserID: req.GetUserId()}); err != nil {
		return nil, toGRPCError(err)
	}

	return &userv1.DeleteAddressResponse{}, nil
}

func (s *UserServiceServer) ListAddresses(ctx context.Context, req *userv1.ListAddressesRequest) (*userv1.ListAddressesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	addresses, err := s.listAddressesService.Execute(ctx, req.GetUserId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	items := make([]*userv1.Address, 0, len(addresses))
	for _, address := range addresses {
		items = append(items, toAddressPB(address))
	}

	return &userv1.ListAddressesResponse{Addresses: items}, nil
}

func toGRPCError(err error) error {
	var appErr *appErrors.Error
	if !stderrors.As(err, &appErr) {
		return status.Error(codes.Internal, "internal server error")
	}

	switch appErr.Code {
	case appErrors.CodeInvalidArgument,
		appErrors.Code("USER_IDENTIFIER_REQUIRED"),
		appErrors.Code("USER_NICKNAME_REQUIRED"),
		appErrors.Code("USER_ADDRESS_RECEIVER_NAME_REQUIRED"),
		appErrors.Code("USER_ADDRESS_RECEIVER_PHONE_REQUIRED"),
		appErrors.Code("USER_ADDRESS_COUNTRY_CODE_REQUIRED"),
		appErrors.Code("USER_ADDRESS_PROVINCE_REQUIRED"),
		appErrors.Code("USER_ADDRESS_CITY_REQUIRED"),
		appErrors.Code("USER_ADDRESS_DISTRICT_REQUIRED"),
		appErrors.Code("USER_ADDRESS_LINE1_REQUIRED"):
		return status.Error(codes.InvalidArgument, appErr.Message)
	case appErrors.Code("USER_EMAIL_ALREADY_EXISTS"),
		appErrors.Code("USER_PHONE_ALREADY_EXISTS"),
		appErrors.Code("USER_ALREADY_EXISTS"):
		return status.Error(codes.AlreadyExists, appErr.Message)
	case appErrors.CodeUnauthorized:
		return status.Error(codes.Unauthenticated, appErr.Message)
	case appErrors.CodeNotFound:
		return status.Error(codes.NotFound, appErr.Message)
	default:
		return status.Error(codes.Internal, appErr.Message)
	}
}

func toAddressPB(address *applicationuser.Address) *userv1.Address {
	if address == nil {
		return nil
	}

	return &userv1.Address{
		Id:            address.ID,
		UserId:        address.UserID,
		ReceiverName:  address.ReceiverName,
		ReceiverPhone: address.ReceiverPhone,
		CountryCode:   address.CountryCode,
		Province:      address.Province,
		City:          address.City,
		District:      address.District,
		AddressLine1:  address.AddressLine1,
		AddressLine2:  address.AddressLine2,
		PostalCode:    address.PostalCode,
		Tag:           address.Tag,
		IsDefault:     address.IsDefault,
		CreatedAt:     address.CreatedAt.Unix(),
		UpdatedAt:     address.UpdatedAt.Unix(),
	}
}

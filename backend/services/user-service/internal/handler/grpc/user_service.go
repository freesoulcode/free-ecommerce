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

	createUserService *applicationuser.CreateUserService
}

func NewUserServiceServer(createUserService *applicationuser.CreateUserService) *UserServiceServer {
	return &UserServiceServer{createUserService: createUserService}
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

func toGRPCError(err error) error {
	var appErr *appErrors.Error
	if !stderrors.As(err, &appErr) {
		return status.Error(codes.Internal, "internal server error")
	}

	switch appErr.Code {
	case appErrors.CodeInvalidArgument,
		appErrors.Code("USER_IDENTIFIER_REQUIRED"),
		appErrors.Code("USER_NICKNAME_REQUIRED"):
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

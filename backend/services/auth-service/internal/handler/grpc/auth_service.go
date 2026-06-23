package grpc

import (
	"context"
	stderrors "errors"

	authv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/auth/v1"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	applicationcredential "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/application/credential"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceServer struct {
	authv1.UnimplementedAuthServiceServer
	createPasswordCredentialService *applicationcredential.CreatePasswordCredentialService
}

func NewAuthServiceServer(createPasswordCredentialService *applicationcredential.CreatePasswordCredentialService) *AuthServiceServer {
	return &AuthServiceServer{createPasswordCredentialService: createPasswordCredentialService}
}

func (s *AuthServiceServer) CreatePasswordCredential(ctx context.Context, req *authv1.CreatePasswordCredentialRequest) (*authv1.CreatePasswordCredentialResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	credential, err := s.createPasswordCredentialService.Execute(ctx, applicationcredential.CreatePasswordCredentialInput{
		UserID:   req.GetUserId(),
		Email:    req.GetEmail(),
		Phone:    req.GetPhone(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.CreatePasswordCredentialResponse{
		UserId:       credential.UserID,
		Email:        credential.Email,
		Phone:        credential.Phone,
		PasswordAlgo: credential.PasswordAlgo,
	}, nil
}

func toGRPCError(err error) error {
	var appErr *appErrors.Error
	if !stderrors.As(err, &appErr) {
		return status.Error(codes.Internal, "internal server error")
	}

	switch appErr.Code {
	case appErrors.CodeInvalidArgument,
		appErrors.Code("AUTH_USER_ID_REQUIRED"),
		appErrors.Code("AUTH_IDENTIFIER_REQUIRED"),
		appErrors.Code("AUTH_PASSWORD_TOO_SHORT"):
		return status.Error(codes.InvalidArgument, appErr.Message)
	case appErrors.Code("AUTH_CREDENTIAL_ALREADY_EXISTS"),
		appErrors.Code("AUTH_EMAIL_ALREADY_EXISTS"),
		appErrors.Code("AUTH_PHONE_ALREADY_EXISTS"):
		return status.Error(codes.AlreadyExists, appErr.Message)
	default:
		return status.Error(codes.Internal, appErr.Message)
	}
}

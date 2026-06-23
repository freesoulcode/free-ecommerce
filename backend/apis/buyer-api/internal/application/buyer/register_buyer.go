package buyer

import (
	"context"
	stderrors "errors"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
)

type UserService interface {
	CreateUser(ctx context.Context, input CreateBuyerInput) (*Buyer, error)
	DeleteUser(ctx context.Context, id int64) error
}

type AuthService interface {
	CreatePasswordCredential(ctx context.Context, input CreateBuyerInput, userID int64) error
}

type CreateBuyerInput struct {
	Email    string
	Phone    string
	Nickname string
	Password string
}

type Buyer struct {
	ID            int64
	Email         string
	Phone         string
	Nickname      string
	Status        string
	EmailVerified bool
	PhoneVerified bool
}

type RegisterBuyerService struct {
	userService UserService
	authService AuthService
}

func NewRegisterBuyerService(userService UserService, authService AuthService) *RegisterBuyerService {
	return &RegisterBuyerService{userService: userService, authService: authService}
}

func (s *RegisterBuyerService) Execute(ctx context.Context, input CreateBuyerInput) (*Buyer, error) {
	if s.userService == nil {
		return nil, appErrors.Internal("user service is not configured")
	}
	if s.authService == nil {
		return nil, appErrors.Internal("auth service is not configured")
	}

	buyer, err := s.userService.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}

	if err := s.authService.CreatePasswordCredential(ctx, input, buyer.ID); err != nil {
		if rollbackErr := s.userService.DeleteUser(ctx, buyer.ID); rollbackErr != nil && !isNotFoundError(rollbackErr) {
			return nil, appErrors.Internal("register buyer failed and rollback user failed")
		}
		return nil, err
	}

	return buyer, nil
}

func isNotFoundError(err error) bool {
	var appErr *appErrors.Error
	if !stderrors.As(err, &appErr) {
		return false
	}

	return appErr.Code == appErrors.CodeNotFound
}

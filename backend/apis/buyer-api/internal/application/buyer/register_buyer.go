package buyer

import (
	"context"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
)

type UserService interface {
	CreateUser(ctx context.Context, input CreateBuyerInput) (*Buyer, error)
}

type CreateBuyerInput struct {
	Email    string
	Phone    string
	Nickname string
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
}

func NewRegisterBuyerService(userService UserService) *RegisterBuyerService {
	return &RegisterBuyerService{userService: userService}
}

func (s *RegisterBuyerService) Execute(ctx context.Context, input CreateBuyerInput) (*Buyer, error) {
	if s.userService == nil {
		return nil, appErrors.Internal("user service is not configured")
	}

	return s.userService.CreateUser(ctx, input)
}

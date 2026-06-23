package buyer

import (
	"context"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
)

type LoginAuthService interface {
	Login(ctx context.Context, input LoginBuyerInput) (*LoginAuthResult, error)
	RefreshToken(ctx context.Context, input RefreshBuyerTokenInput) (*LoginAuthResult, error)
	Logout(ctx context.Context, input LogoutBuyerInput) (*LogoutBuyerResult, error)
}

type BuyerProfileService interface {
	GetUser(ctx context.Context, id int64) (*Buyer, error)
}

type LoginBuyerInput struct {
	Email     string
	Password  string
	DeviceID  string
	UserAgent string
	ClientIP  string
}

type RefreshBuyerTokenInput struct {
	RefreshToken string
	DeviceID     string
	UserAgent    string
	ClientIP     string
}

type LogoutBuyerInput struct {
	RefreshToken string
}

type LoginAuthResult struct {
	UserID                int64
	Email                 string
	Phone                 string
	AccessToken           string
	RefreshToken          string
	TokenType             string
	AccessTokenExpiresAt  int64
	RefreshTokenExpiresAt int64
	RefreshSessionID      int64
}

type LoginBuyerResult struct {
	Buyer                 *Buyer
	AccessToken           string
	RefreshToken          string
	TokenType             string
	AccessTokenExpiresAt  int64
	RefreshTokenExpiresAt int64
	RefreshSessionID      int64
}

type LogoutBuyerResult struct {
	RefreshSessionID int64
}

type LoginBuyerService struct {
	authService    LoginAuthService
	profileService BuyerProfileService
}

func NewLoginBuyerService(authService LoginAuthService, profileService BuyerProfileService) *LoginBuyerService {
	return &LoginBuyerService{authService: authService, profileService: profileService}
}

func (s *LoginBuyerService) Execute(ctx context.Context, input LoginBuyerInput) (*LoginBuyerResult, error) {
	if s.authService == nil {
		return nil, appErrors.Internal("auth service is not configured")
	}
	if s.profileService == nil {
		return nil, appErrors.Internal("user service is not configured")
	}

	authResult, err := s.authService.Login(ctx, input)
	if err != nil {
		return nil, err
	}

	buyer, err := s.profileService.GetUser(ctx, authResult.UserID)
	if err != nil {
		return nil, err
	}

	return &LoginBuyerResult{
		Buyer:                 buyer,
		AccessToken:           authResult.AccessToken,
		RefreshToken:          authResult.RefreshToken,
		TokenType:             authResult.TokenType,
		AccessTokenExpiresAt:  authResult.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: authResult.RefreshTokenExpiresAt,
		RefreshSessionID:      authResult.RefreshSessionID,
	}, nil
}

func (s *LoginBuyerService) Refresh(ctx context.Context, input RefreshBuyerTokenInput) (*LoginBuyerResult, error) {
	if s.authService == nil {
		return nil, appErrors.Internal("auth service is not configured")
	}
	if s.profileService == nil {
		return nil, appErrors.Internal("user service is not configured")
	}

	authResult, err := s.authService.RefreshToken(ctx, input)
	if err != nil {
		return nil, err
	}

	buyer, err := s.profileService.GetUser(ctx, authResult.UserID)
	if err != nil {
		return nil, err
	}

	return &LoginBuyerResult{
		Buyer:                 buyer,
		AccessToken:           authResult.AccessToken,
		RefreshToken:          authResult.RefreshToken,
		TokenType:             authResult.TokenType,
		AccessTokenExpiresAt:  authResult.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: authResult.RefreshTokenExpiresAt,
		RefreshSessionID:      authResult.RefreshSessionID,
	}, nil
}

func (s *LoginBuyerService) Logout(ctx context.Context, input LogoutBuyerInput) (*LogoutBuyerResult, error) {
	if s.authService == nil {
		return nil, appErrors.Internal("auth service is not configured")
	}

	result, err := s.authService.Logout(ctx, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

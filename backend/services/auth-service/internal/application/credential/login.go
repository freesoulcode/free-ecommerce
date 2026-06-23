package credential

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domaincredential "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/domain/credential"
)

type PasswordVerifier interface {
	Verify(password, encodedHash string) (bool, error)
}

type SessionIDGenerator interface {
	NextID() (int64, error)
}

type AccessTokenSigner interface {
	SignAccessToken(input AccessTokenClaims) (string, error)
}

type RefreshTokenGenerator interface {
	Generate() (string, error)
}

type AccessTokenClaims struct {
	UserID     int64
	SessionID  int64
	IssuedAt   time.Time
	ExpiresAt  time.Time
	TokenType  string
	Audience   string
	Issuer     string
	Subject    string
	Identifier string
}

type LoginInput struct {
	Email     string
	Password  string
	DeviceID  string
	UserAgent string
	ClientIP  string
}

type LoginResult struct {
	UserID                int64
	Email                 string
	Phone                 string
	AccessToken           string
	RefreshToken          string
	TokenType             string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	RefreshSessionID      int64
}

type LoginService struct {
	credentialRepo   domaincredential.Repository
	refreshRepo      domaincredential.RefreshSessionRepository
	verifier         PasswordVerifier
	idGenerator      SessionIDGenerator
	accessSigner     AccessTokenSigner
	refreshGenerator RefreshTokenGenerator
	issuer           string
	audience         string
	accessTTL        time.Duration
	refreshTTL       time.Duration
	now              func() time.Time
}

func NewLoginService(
	credentialRepo domaincredential.Repository,
	refreshRepo domaincredential.RefreshSessionRepository,
	verifier PasswordVerifier,
	idGenerator SessionIDGenerator,
	accessSigner AccessTokenSigner,
	refreshGenerator RefreshTokenGenerator,
	issuer string,
	audience string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	now func() time.Time,
) *LoginService {
	if now == nil {
		now = time.Now
	}

	return &LoginService{
		credentialRepo:   credentialRepo,
		refreshRepo:      refreshRepo,
		verifier:         verifier,
		idGenerator:      idGenerator,
		accessSigner:     accessSigner,
		refreshGenerator: refreshGenerator,
		issuer:           issuer,
		audience:         audience,
		accessTTL:        accessTTL,
		refreshTTL:       refreshTTL,
		now:              now,
	}
}

func (s *LoginService) Execute(ctx context.Context, input LoginInput) (*LoginResult, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" {
		return nil, appErrors.New(appErrors.Code("AUTH_EMAIL_REQUIRED"), "email is required", 400)
	}
	if input.Password == "" {
		return nil, appErrors.New(appErrors.Code("AUTH_PASSWORD_REQUIRED"), "password is required", 400)
	}

	credential, err := s.credentialRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	matched, err := s.verifier.Verify(input.Password, credential.PasswordHash)
	if err != nil {
		return nil, appErrors.Internal("verify password failed")
	}
	if !matched {
		return nil, appErrors.Unauthorized("invalid email or password")
	}

	sessionID, err := s.idGenerator.NextID()
	if err != nil {
		return nil, appErrors.Internal("generate refresh session id failed")
	}

	now := s.now().UTC()
	accessExpiresAt := now.Add(s.accessTTL)
	refreshExpiresAt := now.Add(s.refreshTTL)

	accessToken, err := s.accessSigner.SignAccessToken(AccessTokenClaims{
		UserID:     credential.UserID,
		SessionID:  sessionID,
		IssuedAt:   now,
		ExpiresAt:  accessExpiresAt,
		TokenType:  "access",
		Audience:   s.audience,
		Issuer:     s.issuer,
		Subject:    strings.TrimSpace(email),
		Identifier: email,
	})
	if err != nil {
		return nil, appErrors.Internal("sign access token failed")
	}

	refreshToken, err := s.refreshGenerator.Generate()
	if err != nil {
		return nil, appErrors.Internal("generate refresh token failed")
	}

	refreshSession, err := domaincredential.NewRefreshSession(
		sessionID,
		credential.UserID,
		hashToken(refreshToken),
		input.DeviceID,
		input.UserAgent,
		input.ClientIP,
		refreshExpiresAt,
		now,
	)
	if err != nil {
		return nil, err
	}

	if err := s.refreshRepo.CreateRefreshSession(ctx, refreshSession); err != nil {
		return nil, err
	}

	return &LoginResult{
		UserID:                credential.UserID,
		Email:                 credential.Email,
		Phone:                 credential.Phone,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		TokenType:             "Bearer",
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshTokenExpiresAt: refreshExpiresAt,
		RefreshSessionID:      sessionID,
	}, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

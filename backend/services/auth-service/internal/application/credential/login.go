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

type RefreshTokenInput struct {
	RefreshToken string
	DeviceID     string
	UserAgent    string
	ClientIP     string
}

type LogoutInput struct {
	RefreshToken string
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

type RefreshTokenResult = LoginResult

type LogoutResult struct {
	RefreshSessionID int64
}

func issueSession(
	now time.Time,
	credential *domaincredential.PasswordCredential,
	sessionID int64,
	deviceID string,
	userAgent string,
	clientIP string,
	issuer string,
	audience string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	accessSigner AccessTokenSigner,
	refreshGenerator RefreshTokenGenerator,
) (*RefreshTokenResult, *domaincredential.RefreshSession, error) {
	accessExpiresAt := now.Add(accessTTL)
	refreshExpiresAt := now.Add(refreshTTL)

	accessToken, err := accessSigner.SignAccessToken(AccessTokenClaims{
		UserID:     credential.UserID,
		SessionID:  sessionID,
		IssuedAt:   now,
		ExpiresAt:  accessExpiresAt,
		TokenType:  "access",
		Audience:   audience,
		Issuer:     issuer,
		Subject:    strings.TrimSpace(credential.Email),
		Identifier: credential.Email,
	})
	if err != nil {
		return nil, nil, appErrors.Internal("sign access token failed")
	}

	refreshToken, err := refreshGenerator.Generate()
	if err != nil {
		return nil, nil, appErrors.Internal("generate refresh token failed")
	}

	refreshSession, err := domaincredential.NewRefreshSession(
		sessionID,
		credential.UserID,
		hashToken(refreshToken),
		deviceID,
		userAgent,
		clientIP,
		refreshExpiresAt,
		now,
	)
	if err != nil {
		return nil, nil, err
	}

	return &RefreshTokenResult{
		UserID:                credential.UserID,
		Email:                 credential.Email,
		Phone:                 credential.Phone,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		TokenType:             "Bearer",
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshTokenExpiresAt: refreshExpiresAt,
		RefreshSessionID:      sessionID,
	}, refreshSession, nil
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
	result, refreshSession, err := issueSession(
		now,
		credential,
		sessionID,
		input.DeviceID,
		input.UserAgent,
		input.ClientIP,
		s.issuer,
		s.audience,
		s.accessTTL,
		s.refreshTTL,
		s.accessSigner,
		s.refreshGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := s.refreshRepo.CreateRefreshSession(ctx, refreshSession); err != nil {
		return nil, err
	}

	return (*LoginResult)(result), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

type RefreshTokenService struct {
	credentialRepo    domaincredential.Repository
	refreshRepo       domaincredential.RefreshSessionRepository
	idGenerator       SessionIDGenerator
	accessSigner      AccessTokenSigner
	refreshGenerator  RefreshTokenGenerator
	issuer            string
	audience          string
	accessTTL         time.Duration
	refreshTTL        time.Duration
	now               func() time.Time
}

func NewRefreshTokenService(
	credentialRepo domaincredential.Repository,
	refreshRepo domaincredential.RefreshSessionRepository,
	idGenerator SessionIDGenerator,
	accessSigner AccessTokenSigner,
	refreshGenerator RefreshTokenGenerator,
	issuer string,
	audience string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	now func() time.Time,
) *RefreshTokenService {
	if now == nil {
		now = time.Now
	}

	return &RefreshTokenService{
		credentialRepo:   credentialRepo,
		refreshRepo:      refreshRepo,
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

func (s *RefreshTokenService) Execute(ctx context.Context, input RefreshTokenInput) (*RefreshTokenResult, error) {
	refreshToken := strings.TrimSpace(input.RefreshToken)
	if refreshToken == "" {
		return nil, appErrors.New(appErrors.Code("AUTH_REFRESH_TOKEN_REQUIRED"), "refresh token is required", 400)
	}

	now := s.now().UTC()
	currentSession, err := s.refreshRepo.FindRefreshSessionByTokenHash(ctx, hashToken(refreshToken))
	if err != nil {
		return nil, err
	}
	if !currentSession.CanRefresh(now) {
		return nil, appErrors.Unauthorized("refresh token is invalid")
	}

	credential, err := s.credentialRepo.FindByUserID(ctx, currentSession.UserID)
	if err != nil {
		return nil, err
	}

	newSessionID, err := s.idGenerator.NextID()
	if err != nil {
		return nil, appErrors.Internal("generate refresh session id failed")
	}

	accessExpiresAt := now.Add(s.accessTTL)
	refreshExpiresAt := now.Add(s.refreshTTL)
	accessToken, err := s.accessSigner.SignAccessToken(AccessTokenClaims{
		UserID:     credential.UserID,
		SessionID:  newSessionID,
		IssuedAt:   now,
		ExpiresAt:  accessExpiresAt,
		TokenType:  "access",
		Audience:   s.audience,
		Issuer:     s.issuer,
		Subject:    strings.TrimSpace(credential.Email),
		Identifier: credential.Email,
	})
	if err != nil {
		return nil, appErrors.Internal("sign access token failed")
	}

	newRefreshToken, err := s.refreshGenerator.Generate()
	if err != nil {
		return nil, appErrors.Internal("generate refresh token failed")
	}

	newSession, err := domaincredential.NewRefreshSession(
		newSessionID,
		credential.UserID,
		hashToken(newRefreshToken),
		fallbackString(input.DeviceID, currentSession.DeviceID),
		fallbackString(input.UserAgent, currentSession.UserAgent),
		fallbackString(input.ClientIP, currentSession.ClientIP),
		refreshExpiresAt,
		now,
	)
	if err != nil {
		return nil, err
	}

	if err := s.refreshRepo.RotateRefreshSession(ctx, currentSession.ID, now, newSessionID, newSession); err != nil {
		return nil, err
	}

	return &RefreshTokenResult{
		UserID:                credential.UserID,
		Email:                 credential.Email,
		Phone:                 credential.Phone,
		AccessToken:           accessToken,
		RefreshToken:          newRefreshToken,
		TokenType:             "Bearer",
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshTokenExpiresAt: refreshExpiresAt,
		RefreshSessionID:      newSessionID,
	}, nil
}

type LogoutService struct {
	refreshRepo domaincredential.RefreshSessionRepository
	now         func() time.Time
}

func NewLogoutService(refreshRepo domaincredential.RefreshSessionRepository, now func() time.Time) *LogoutService {
	if now == nil {
		now = time.Now
	}

	return &LogoutService{refreshRepo: refreshRepo, now: now}
}

func (s *LogoutService) Execute(ctx context.Context, input LogoutInput) (*LogoutResult, error) {
	refreshToken := strings.TrimSpace(input.RefreshToken)
	if refreshToken == "" {
		return nil, appErrors.New(appErrors.Code("AUTH_REFRESH_TOKEN_REQUIRED"), "refresh token is required", 400)
	}

	session, err := s.refreshRepo.FindRefreshSessionByTokenHash(ctx, hashToken(refreshToken))
	if err != nil {
		return nil, err
	}
	if session.IsRevoked() {
		return &LogoutResult{RefreshSessionID: session.ID}, nil
	}

	if err := s.refreshRepo.RevokeRefreshSession(ctx, session.ID, s.now().UTC()); err != nil {
		return nil, err
	}

	return &LogoutResult{RefreshSessionID: session.ID}, nil
}

func fallbackString(preferred, fallback string) string {
	preferred = strings.TrimSpace(preferred)
	if preferred != "" {
		return preferred
	}

	return strings.TrimSpace(fallback)
}

func issueSessionID(idGenerator SessionIDGenerator) (int64, error) {
	sessionID, err := idGenerator.NextID()
	if err != nil {
		return 0, appErrors.Internal("generate refresh session id failed")
	}

	return sessionID, nil
}

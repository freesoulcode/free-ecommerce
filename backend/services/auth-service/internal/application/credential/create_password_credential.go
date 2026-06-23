package credential

import (
	"context"
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domaincredential "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/domain/credential"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type CreatePasswordCredentialInput struct {
	UserID   int64
	Email    string
	Phone    string
	Password string
}

type CreatePasswordCredentialService struct {
	repo   domaincredential.Repository
	hasher PasswordHasher
	now    func() time.Time
}

func NewCreatePasswordCredentialService(repo domaincredential.Repository, hasher PasswordHasher, now func() time.Time) *CreatePasswordCredentialService {
	if now == nil {
		now = time.Now
	}

	return &CreatePasswordCredentialService{repo: repo, hasher: hasher, now: now}
}

func (s *CreatePasswordCredentialService) Execute(ctx context.Context, input CreatePasswordCredentialInput) (*domaincredential.PasswordCredential, error) {
	password := strings.TrimSpace(input.Password)
	if len(password) < 8 {
		return nil, appErrors.New(appErrors.Code("AUTH_PASSWORD_TOO_SHORT"), "password must be at least 8 characters", 400)
	}

	exists, err := s.repo.ExistsByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, appErrors.New(appErrors.Code("AUTH_CREDENTIAL_ALREADY_EXISTS"), "credential already exists", 400)
	}

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email != "" {
		exists, err = s.repo.ExistsByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, appErrors.New(appErrors.Code("AUTH_EMAIL_ALREADY_EXISTS"), "email credential already exists", 400)
		}
	}

	phone := strings.TrimSpace(input.Phone)
	if phone != "" {
		exists, err = s.repo.ExistsByPhone(ctx, phone)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, appErrors.New(appErrors.Code("AUTH_PHONE_ALREADY_EXISTS"), "phone credential already exists", 400)
		}
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, appErrors.Internal("hash password failed")
	}

	credential, err := domaincredential.NewPasswordCredential(input.UserID, email, phone, hash, s.now())
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreatePasswordCredential(ctx, credential); err != nil {
		return nil, err
	}

	return credential, nil
}

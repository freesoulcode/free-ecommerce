package user

import (
	"context"
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domainuser "github.com/freesoulcode/free-ecommerce/backend/services/user-service/internal/domain/user"
)

type IDGenerator interface {
	NextID() (int64, error)
}

type CreateUserInput struct {
	Email    string
	Phone    string
	Nickname string
}

type CreateUserService struct {
	repo        domainuser.Repository
	idGenerator IDGenerator
	now         func() time.Time
}

func NewCreateUserService(repo domainuser.Repository, idGenerator IDGenerator, now func() time.Time) *CreateUserService {
	if now == nil {
		now = time.Now
	}

	return &CreateUserService{
		repo:        repo,
		idGenerator: idGenerator,
		now:         now,
	}
}

func (s *CreateUserService) Execute(ctx context.Context, input CreateUserInput) (*domainuser.User, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	phone := strings.TrimSpace(input.Phone)

	if email != "" {
		exists, err := s.repo.ExistsByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, appErrors.New(appErrors.Code("USER_EMAIL_ALREADY_EXISTS"), "email already exists", 400)
		}
	}

	if phone != "" {
		exists, err := s.repo.ExistsByPhone(ctx, phone)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, appErrors.New(appErrors.Code("USER_PHONE_ALREADY_EXISTS"), "phone already exists", 400)
		}
	}

	id, err := s.idGenerator.NextID()
	if err != nil {
		return nil, appErrors.Internal("generate user id failed")
	}

	user, err := domainuser.New(id, email, phone, input.Nickname, s.now())
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

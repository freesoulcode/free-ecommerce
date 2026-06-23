package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	DeleteByID(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
}

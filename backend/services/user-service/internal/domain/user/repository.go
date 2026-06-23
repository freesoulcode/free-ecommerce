package user

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	DeleteByID(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	CreateAddress(ctx context.Context, address *Address) error
	UpdateAddress(ctx context.Context, address *Address) error
	DeleteAddress(ctx context.Context, userID, addressID int64) error
	FindAddressByID(ctx context.Context, userID, addressID int64) (*Address, error)
	ListAddressesByUserID(ctx context.Context, userID int64) ([]*Address, error)
	CountAddressesByUserID(ctx context.Context, userID int64) (int64, error)
	SetDefaultAddress(ctx context.Context, userID, addressID int64, updatedAt time.Time) error
}

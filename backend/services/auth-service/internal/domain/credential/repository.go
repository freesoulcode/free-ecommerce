package credential

import "context"

type Repository interface {
	CreatePasswordCredential(ctx context.Context, credential *PasswordCredential) error
	ExistsByUserID(ctx context.Context, userID int64) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
}

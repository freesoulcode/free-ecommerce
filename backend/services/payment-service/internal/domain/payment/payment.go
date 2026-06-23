package payment

import "time"

const (
	StatusPending = "pending"
	StatusPaid    = "paid"
	StatusExpired = "expired"

	ChannelMock = "mock"
)

type Order struct {
	ID           int64
	UserID       int64
	OrderGroupID int64
	Status       string
	Channel      string
	PayAmount    int64
	Currency     string
	ExpireAt     time.Time
	PaidAt       *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

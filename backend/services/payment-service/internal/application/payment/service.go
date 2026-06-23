package payment

import (
	"context"
	"strings"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domainpayment "github.com/freesoulcode/free-ecommerce/backend/services/payment-service/internal/domain/payment"
)

type OrderPaymentInfo struct {
	OrderGroupID      int64
	UserID            int64
	Status            string
	TotalPayAmount    int64
	Currency          string
	PaymentDeadlineAt time.Time
	PaidAt            *time.Time
}

type OrderService interface {
	GetOrderGroupPaymentInfo(ctx context.Context, userID, orderGroupID int64) (*OrderPaymentInfo, error)
	MarkOrderGroupPaid(ctx context.Context, userID, orderGroupID int64) (*OrderPaymentInfo, error)
	CloseOrderGroupByPaymentTimeout(ctx context.Context, userID, orderGroupID int64) (*OrderPaymentInfo, error)
}

type CreatePaymentOrderInput struct {
	UserID       int64
	OrderGroupID int64
	Channel      string
}

type CreatePaymentOrderService struct {
	repo        domainpayment.Repository
	idGenerator IDGenerator
	orderSvc    OrderService
	now         func() time.Time
}

type GetPaymentOrderService struct {
	repo        domainpayment.Repository
	idGenerator IDGenerator
	orderSvc    OrderService
	now         func() time.Time
}

type SimulatePayService struct {
	repo        domainpayment.Repository
	idGenerator IDGenerator
	orderSvc    OrderService
	now         func() time.Time
}

func NewCreatePaymentOrderService(repo domainpayment.Repository, idGenerator IDGenerator, orderSvc OrderService, now func() time.Time) *CreatePaymentOrderService {
	if now == nil {
		now = time.Now
	}
	return &CreatePaymentOrderService{repo: repo, idGenerator: idGenerator, orderSvc: orderSvc, now: now}
}

func NewGetPaymentOrderService(repo domainpayment.Repository, idGenerator IDGenerator, orderSvc OrderService, now func() time.Time) *GetPaymentOrderService {
	if now == nil {
		now = time.Now
	}
	return &GetPaymentOrderService{repo: repo, idGenerator: idGenerator, orderSvc: orderSvc, now: now}
}

func NewSimulatePayService(repo domainpayment.Repository, idGenerator IDGenerator, orderSvc OrderService, now func() time.Time) *SimulatePayService {
	if now == nil {
		now = time.Now
	}
	return &SimulatePayService{repo: repo, idGenerator: idGenerator, orderSvc: orderSvc, now: now}
}

func (s *CreatePaymentOrderService) Execute(ctx context.Context, input CreatePaymentOrderInput) (*domainpayment.Order, error) {
	channel := strings.TrimSpace(input.Channel)
	if channel == "" {
		channel = domainpayment.ChannelMock
	}
	if channel != domainpayment.ChannelMock {
		return nil, appErrors.InvalidArgument("payment channel is invalid")
	}
	if input.UserID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}
	if input.OrderGroupID <= 0 {
		return nil, appErrors.InvalidArgument("order group id is required")
	}

	existing, err := s.repo.FindByOrderGroup(ctx, input.UserID, input.OrderGroupID)
	if err == nil {
		return existing, nil
	}
	if !isNotFound(err) {
		return nil, err
	}

	info, err := s.orderSvc.GetOrderGroupPaymentInfo(ctx, input.UserID, input.OrderGroupID)
	if err != nil {
		return nil, err
	}
	now := s.now().UTC()
	if info.Status == "closed" || info.Status == "cancelled" || info.Status == "completed" || info.Status == "merchant_processing" {
		return nil, appErrors.InvalidArgument("order group status does not allow payment")
	}
	if !info.PaymentDeadlineAt.IsZero() && now.After(info.PaymentDeadlineAt) {
		_, _ = s.orderSvc.CloseOrderGroupByPaymentTimeout(ctx, input.UserID, input.OrderGroupID)
		return nil, appErrors.InvalidArgument("payment order expired")
	}
	orderID, err := s.idGenerator.NextID()
	if err != nil {
		return nil, appErrors.Internal("generate payment order id failed")
	}
	order := &domainpayment.Order{
		ID:           orderID,
		UserID:       input.UserID,
		OrderGroupID: input.OrderGroupID,
		Status:       domainpayment.StatusPending,
		Channel:      channel,
		PayAmount:    info.TotalPayAmount,
		Currency:     info.Currency,
		ExpireAt:     info.PaymentDeadlineAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if info.Status == "paid" {
		order.Status = domainpayment.StatusPaid
		order.PaidAt = info.PaidAt
	}
	if err := s.repo.Create(ctx, order); err != nil {
		if existing, getErr := s.repo.FindByOrderGroup(ctx, input.UserID, input.OrderGroupID); getErr == nil {
			return existing, nil
		}
		return nil, err
	}
	return order, nil
}

func (s *GetPaymentOrderService) Execute(ctx context.Context, userID, orderGroupID int64) (*domainpayment.Order, error) {
	if userID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}
	if orderGroupID <= 0 {
		return nil, appErrors.InvalidArgument("order group id is required")
	}
	order, err := s.repo.FindByOrderGroup(ctx, userID, orderGroupID)
	if err == nil {
		if order.Status == domainpayment.StatusPending && !order.ExpireAt.IsZero() && s.now().UTC().After(order.ExpireAt) {
			_, _ = s.orderSvc.CloseOrderGroupByPaymentTimeout(ctx, userID, orderGroupID)
			return s.repo.MarkExpired(ctx, userID, orderGroupID, s.now().UTC())
		}
		return order, nil
	}
	if !isNotFound(err) {
		return nil, err
	}
	info, infoErr := s.orderSvc.GetOrderGroupPaymentInfo(ctx, userID, orderGroupID)
	if infoErr != nil {
		return nil, infoErr
	}
	if info.Status == "paid" {
		return s.rebuildPaidOrder(ctx, info)
	}
	return nil, err
}

func (s *SimulatePayService) Execute(ctx context.Context, userID, orderGroupID int64) (*domainpayment.Order, error) {
	if userID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}
	if orderGroupID <= 0 {
		return nil, appErrors.InvalidArgument("order group id is required")
	}
	order, err := s.repo.FindByOrderGroup(ctx, userID, orderGroupID)
	if err != nil {
		if !isNotFound(err) {
			return nil, err
		}
		creator := NewCreatePaymentOrderService(s.repo, s.idGenerator, s.orderSvc, s.now)
		order, err = creator.Execute(ctx, CreatePaymentOrderInput{UserID: userID, OrderGroupID: orderGroupID, Channel: domainpayment.ChannelMock})
		if err != nil {
			return nil, err
		}
	}
	if order.Status == domainpayment.StatusPaid {
		return order, nil
	}
	if order.Status == domainpayment.StatusExpired {
		return nil, appErrors.InvalidArgument("payment order expired")
	}
	if !order.ExpireAt.IsZero() && s.now().UTC().After(order.ExpireAt) {
		_, _ = s.orderSvc.CloseOrderGroupByPaymentTimeout(ctx, userID, orderGroupID)
		_, _ = s.repo.MarkExpired(ctx, userID, orderGroupID, s.now().UTC())
		return nil, appErrors.InvalidArgument("payment order expired")
	}
	info, err := s.orderSvc.GetOrderGroupPaymentInfo(ctx, userID, orderGroupID)
	if err != nil {
		return nil, err
	}
	if info.Status == "paid" {
		return s.repo.MarkPaid(ctx, userID, orderGroupID, paidAtOrNow(info.PaidAt, s.now().UTC()))
	}
	if info.Status != "pending_payment" {
		return nil, appErrors.InvalidArgument("order group status does not allow payment")
	}
	if !info.PaymentDeadlineAt.IsZero() && s.now().UTC().After(info.PaymentDeadlineAt) {
		_, _ = s.orderSvc.CloseOrderGroupByPaymentTimeout(ctx, userID, orderGroupID)
		_, _ = s.repo.MarkExpired(ctx, userID, orderGroupID, s.now().UTC())
		return nil, appErrors.InvalidArgument("payment order expired")
	}
	marked, err := s.orderSvc.MarkOrderGroupPaid(ctx, userID, orderGroupID)
	if err != nil {
		return nil, err
	}
	return s.repo.MarkPaid(ctx, userID, orderGroupID, paidAtOrNow(marked.PaidAt, s.now().UTC()))
}

func (s *GetPaymentOrderService) rebuildPaidOrder(ctx context.Context, info *OrderPaymentInfo) (*domainpayment.Order, error) {
	orderID, err := s.idGenerator.NextID()
	if err != nil {
		return nil, appErrors.Internal("generate payment order id failed")
	}
	now := s.now().UTC()
	order := &domainpayment.Order{
		ID:           orderID,
		UserID:       info.UserID,
		OrderGroupID: info.OrderGroupID,
		Status:       domainpayment.StatusPaid,
		Channel:      domainpayment.ChannelMock,
		PayAmount:    info.TotalPayAmount,
		Currency:     info.Currency,
		ExpireAt:     info.PaymentDeadlineAt,
		PaidAt:       info.PaidAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.repo.Create(ctx, order); err != nil {
		if existing, getErr := s.repo.FindByOrderGroup(ctx, info.UserID, info.OrderGroupID); getErr == nil {
			return existing, nil
		}
		return nil, err
	}
	return order, nil
}

func isNotFound(err error) bool {
	appErr, ok := err.(*appErrors.Error)
	return ok && appErr.Code == appErrors.CodeNotFound
}

func paidAtOrNow(paidAt *time.Time, fallback time.Time) time.Time {
	if paidAt != nil {
		return paidAt.UTC()
	}
	return fallback.UTC()
}

package grpc

import (
	"context"
	stderrors "errors"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	applicationpayment "github.com/freesoulcode/free-ecommerce/backend/services/payment-service/internal/application/payment"
	domainpayment "github.com/freesoulcode/free-ecommerce/backend/services/payment-service/internal/domain/payment"
	paymentv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServiceServer struct {
	paymentv1.UnimplementedPaymentServiceServer
	createService   *applicationpayment.CreatePaymentOrderService
	getService      *applicationpayment.GetPaymentOrderService
	simulatePay     *applicationpayment.SimulatePayService
	listAdminOrders *applicationpayment.ListAdminPaymentOrdersService
	getAdminOrder   *applicationpayment.GetAdminPaymentOrderService
}

func NewPaymentServiceServer(createService *applicationpayment.CreatePaymentOrderService, getService *applicationpayment.GetPaymentOrderService, simulatePay *applicationpayment.SimulatePayService, listAdminOrders *applicationpayment.ListAdminPaymentOrdersService, getAdminOrder *applicationpayment.GetAdminPaymentOrderService) *PaymentServiceServer {
	return &PaymentServiceServer{createService: createService, getService: getService, simulatePay: simulatePay, listAdminOrders: listAdminOrders, getAdminOrder: getAdminOrder}
}

func (s *PaymentServiceServer) CreatePaymentOrder(ctx context.Context, req *paymentv1.CreatePaymentOrderRequest) (*paymentv1.CreatePaymentOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	order, err := s.createService.Execute(ctx, applicationpayment.CreatePaymentOrderInput{UserID: req.GetUserId(), OrderGroupID: req.GetOrderGroupId(), Channel: req.GetChannel()})
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &paymentv1.CreatePaymentOrderResponse{PaymentOrder: toPaymentOrderPB(order)}, nil
}

func (s *PaymentServiceServer) GetPaymentOrderByOrderGroup(ctx context.Context, req *paymentv1.GetPaymentOrderByOrderGroupRequest) (*paymentv1.GetPaymentOrderByOrderGroupResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	order, err := s.getService.Execute(ctx, req.GetUserId(), req.GetOrderGroupId())
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &paymentv1.GetPaymentOrderByOrderGroupResponse{PaymentOrder: toPaymentOrderPB(order)}, nil
}

func (s *PaymentServiceServer) SimulatePay(ctx context.Context, req *paymentv1.SimulatePayRequest) (*paymentv1.SimulatePayResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	order, err := s.simulatePay.Execute(ctx, req.GetUserId(), req.GetOrderGroupId())
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &paymentv1.SimulatePayResponse{PaymentOrder: toPaymentOrderPB(order)}, nil
}

func (s *PaymentServiceServer) ListAdminPaymentOrders(ctx context.Context, req *paymentv1.ListAdminPaymentOrdersRequest) (*paymentv1.ListAdminPaymentOrdersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	result, err := s.listAdminOrders.Execute(ctx, applicationpayment.ListAdminPaymentOrdersInput{
		Page:         req.GetPage(),
		PageSize:     req.GetPageSize(),
		Status:       req.GetStatus(),
		UserID:       req.GetUserId(),
		OrderGroupID: req.GetOrderGroupId(),
		Channel:      req.GetChannel(),
	})
	if err != nil {
		return nil, toGRPCError(err)
	}
	items := make([]*paymentv1.PaymentOrder, 0, len(result.PaymentOrders))
	for _, item := range result.PaymentOrders {
		items = append(items, toPaymentOrderPB(item))
	}
	return &paymentv1.ListAdminPaymentOrdersResponse{PaymentOrders: items, Total: result.Total, Page: result.Page, PageSize: result.PageSize}, nil
}

func (s *PaymentServiceServer) GetAdminPaymentOrder(ctx context.Context, req *paymentv1.GetAdminPaymentOrderRequest) (*paymentv1.GetAdminPaymentOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	order, err := s.getAdminOrder.Execute(ctx, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &paymentv1.GetAdminPaymentOrderResponse{PaymentOrder: toPaymentOrderPB(order)}, nil
}

func toGRPCError(err error) error {
	var appErr *appErrors.Error
	if !stderrors.As(err, &appErr) {
		return status.Error(codes.Internal, "internal server error")
	}
	switch appErr.Code {
	case appErrors.CodeInvalidArgument:
		return status.Error(codes.InvalidArgument, appErr.Message)
	case appErrors.CodeNotFound:
		return status.Error(codes.NotFound, appErr.Message)
	case appErrors.CodeUnauthorized:
		return status.Error(codes.Unauthenticated, appErr.Message)
	default:
		return status.Error(codes.Internal, appErr.Message)
	}
}

func toPaymentOrderPB(order *domainpayment.Order) *paymentv1.PaymentOrder {
	if order == nil {
		return nil
	}
	return &paymentv1.PaymentOrder{Id: order.ID, UserId: order.UserID, OrderGroupId: order.OrderGroupID, Status: order.Status, Channel: order.Channel, PayAmount: order.PayAmount, Currency: order.Currency, ExpireAt: order.ExpireAt.Unix(), PaidAt: unixOrZero(order.PaidAt), CreatedAt: order.CreatedAt.Unix(), UpdatedAt: order.UpdatedAt.Unix()}
}

func unixOrZero(value *time.Time) int64 {
	if value == nil {
		return 0
	}
	return value.Unix()
}

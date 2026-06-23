package grpc

import (
	"context"
	stderrors "errors"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	applicationorder "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/application/order"
	domainorder "github.com/freesoulcode/free-ecommerce/backend/services/order-service/internal/domain/order"
	orderv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/order/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServiceServer struct {
	orderv1.UnimplementedOrderServiceServer
	submitService       *applicationorder.SubmitOrderService
	listService         *applicationorder.ListBuyerOrderGroupsService
	getService          *applicationorder.GetBuyerOrderGroupDetailService
	getPaymentInfo      *applicationorder.GetOrderGroupPaymentInfoService
	markPaidService     *applicationorder.MarkOrderGroupPaidService
	closeTimeoutService *applicationorder.CloseOrderGroupByPaymentTimeoutService
}

func NewOrderServiceServer(
	submitService *applicationorder.SubmitOrderService,
	listService *applicationorder.ListBuyerOrderGroupsService,
	getService *applicationorder.GetBuyerOrderGroupDetailService,
	getPaymentInfo *applicationorder.GetOrderGroupPaymentInfoService,
	markPaidService *applicationorder.MarkOrderGroupPaidService,
	closeTimeoutService *applicationorder.CloseOrderGroupByPaymentTimeoutService,
) *OrderServiceServer {
	return &OrderServiceServer{
		submitService:       submitService,
		listService:         listService,
		getService:          getService,
		getPaymentInfo:      getPaymentInfo,
		markPaidService:     markPaidService,
		closeTimeoutService: closeTimeoutService,
	}
}

func (s *OrderServiceServer) SubmitOrder(ctx context.Context, req *orderv1.SubmitOrderRequest) (*orderv1.SubmitOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	group, err := s.submitService.Execute(ctx, applicationorder.SubmitOrderInput{
		UserID:      req.GetUserId(),
		AddressID:   req.GetAddressId(),
		CartItemIDs: req.GetCartItemIds(),
		Source:      req.GetSource(),
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &orderv1.SubmitOrderResponse{OrderGroup: toOrderGroupDetailPB(group)}, nil
}

func (s *OrderServiceServer) ListBuyerOrderGroups(ctx context.Context, req *orderv1.ListBuyerOrderGroupsRequest) (*orderv1.ListBuyerOrderGroupsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	result, err := s.listService.Execute(ctx, applicationorder.ListBuyerOrderGroupsInput{
		UserID:   req.GetUserId(),
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
		Status:   req.GetStatus(),
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	items := make([]*orderv1.OrderGroupSummary, 0, len(result.OrderGroups))
	for _, group := range result.OrderGroups {
		items = append(items, toOrderGroupSummaryPB(group))
	}

	return &orderv1.ListBuyerOrderGroupsResponse{OrderGroups: items, Total: result.Total, Page: result.Page, PageSize: result.PageSize}, nil
}

func (s *OrderServiceServer) GetBuyerOrderGroupDetail(ctx context.Context, req *orderv1.GetBuyerOrderGroupDetailRequest) (*orderv1.GetBuyerOrderGroupDetailResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	group, err := s.getService.Execute(ctx, req.GetUserId(), req.GetOrderGroupId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &orderv1.GetBuyerOrderGroupDetailResponse{OrderGroup: toOrderGroupDetailPB(group)}, nil
}

func (s *OrderServiceServer) GetOrderGroupPaymentInfo(ctx context.Context, req *orderv1.GetOrderGroupPaymentInfoRequest) (*orderv1.GetOrderGroupPaymentInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	info, err := s.getPaymentInfo.Execute(ctx, req.GetUserId(), req.GetOrderGroupId())
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &orderv1.GetOrderGroupPaymentInfoResponse{PaymentInfo: toPaymentInfoPB(info)}, nil
}

func (s *OrderServiceServer) MarkOrderGroupPaid(ctx context.Context, req *orderv1.MarkOrderGroupPaidRequest) (*orderv1.MarkOrderGroupPaidResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	info, err := s.markPaidService.Execute(ctx, req.GetUserId(), req.GetOrderGroupId())
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &orderv1.MarkOrderGroupPaidResponse{PaymentInfo: toPaymentInfoPB(info)}, nil
}

func (s *OrderServiceServer) CloseOrderGroupByPaymentTimeout(ctx context.Context, req *orderv1.CloseOrderGroupByPaymentTimeoutRequest) (*orderv1.CloseOrderGroupByPaymentTimeoutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	info, err := s.closeTimeoutService.Execute(ctx, req.GetUserId(), req.GetOrderGroupId())
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &orderv1.CloseOrderGroupByPaymentTimeoutResponse{PaymentInfo: toPaymentInfoPB(info)}, nil
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

func toAddressSnapshotPB(address *domainorder.AddressSnapshot) *orderv1.AddressSnapshot {
	if address == nil {
		return nil
	}

	return &orderv1.AddressSnapshot{
		Id:            address.ID,
		OrderGroupId:  address.OrderGroupID,
		UserId:        address.UserID,
		ReceiverName:  address.ReceiverName,
		ReceiverPhone: address.ReceiverPhone,
		CountryCode:   address.CountryCode,
		Province:      address.Province,
		City:          address.City,
		District:      address.District,
		AddressLine1:  address.AddressLine1,
		AddressLine2:  address.AddressLine2,
		PostalCode:    address.PostalCode,
		Tag:           address.Tag,
		CreatedAt:     address.CreatedAt.Unix(),
		UpdatedAt:     address.UpdatedAt.Unix(),
	}
}

func toOrderItemPB(item *domainorder.Item) *orderv1.OrderItem {
	if item == nil {
		return nil
	}

	return &orderv1.OrderItem{
		Id:                        item.ID,
		OrderGroupId:              item.OrderGroupID,
		ShopOrderId:               item.ShopOrderID,
		UserId:                    item.UserID,
		ShopId:                    item.ShopID,
		ProductId:                 item.ProductID,
		SkuId:                     item.SKUID,
		ProductTitle:              item.ProductTitle,
		ProductSubTitle:           item.ProductSubTitle,
		MainImageUrl:              item.MainImageURL,
		SkuName:                   item.SKUName,
		PriceAmount:               item.PriceAmount,
		Currency:                  item.Currency,
		Quantity:                  item.Quantity,
		ItemAmount:                item.ItemAmount,
		ReviewStatusSnapshot:      item.ReviewStatusSnapshot,
		ProductSaleStatusSnapshot: item.ProductSaleStatusSnapshot,
		SkuSaleStatusSnapshot:     item.SKUSaleStatusSnapshot,
		CreatedAt:                 item.CreatedAt.Unix(),
		UpdatedAt:                 item.UpdatedAt.Unix(),
	}
}

func toShopOrderSummaryPB(shopOrder *domainorder.ShopOrderSummary) *orderv1.ShopOrderSummary {
	if shopOrder == nil {
		return nil
	}

	return &orderv1.ShopOrderSummary{
		Id:             shopOrder.ID,
		OrderGroupId:   shopOrder.OrderGroupID,
		UserId:         shopOrder.UserID,
		ShopId:         shopOrder.ShopID,
		ShopName:       shopOrder.ShopName,
		Status:         shopOrder.Status,
		ItemAmount:     shopOrder.ItemAmount,
		ShippingAmount: shopOrder.ShippingAmount,
		PayAmount:      shopOrder.PayAmount,
		Currency:       shopOrder.Currency,
		ItemCount:      shopOrder.ItemCount,
		CreatedAt:      shopOrder.CreatedAt.Unix(),
		UpdatedAt:      shopOrder.UpdatedAt.Unix(),
		PaidAt:         unixOrZero(shopOrder.PaidAt),
	}
}

func toShopOrderPB(shopOrder *domainorder.ShopOrder) *orderv1.ShopOrder {
	if shopOrder == nil {
		return nil
	}

	items := make([]*orderv1.OrderItem, 0, len(shopOrder.Items))
	for _, item := range shopOrder.Items {
		items = append(items, toOrderItemPB(item))
	}

	return &orderv1.ShopOrder{
		Id:             shopOrder.ID,
		OrderGroupId:   shopOrder.OrderGroupID,
		UserId:         shopOrder.UserID,
		ShopId:         shopOrder.ShopID,
		ShopName:       shopOrder.ShopName,
		Status:         shopOrder.Status,
		ItemAmount:     shopOrder.ItemAmount,
		ShippingAmount: shopOrder.ShippingAmount,
		PayAmount:      shopOrder.PayAmount,
		Currency:       shopOrder.Currency,
		ItemCount:      shopOrder.ItemCount,
		Items:          items,
		CreatedAt:      shopOrder.CreatedAt.Unix(),
		UpdatedAt:      shopOrder.UpdatedAt.Unix(),
		PaidAt:         unixOrZero(shopOrder.PaidAt),
	}
}

func toOrderGroupSummaryPB(group *domainorder.GroupSummary) *orderv1.OrderGroupSummary {
	if group == nil {
		return nil
	}

	shopOrders := make([]*orderv1.ShopOrderSummary, 0, len(group.ShopOrders))
	for _, shopOrder := range group.ShopOrders {
		shopOrders = append(shopOrders, toShopOrderSummaryPB(shopOrder))
	}

	return &orderv1.OrderGroupSummary{
		Id:                  group.ID,
		UserId:              group.UserID,
		Status:              group.Status,
		Source:              group.Source,
		TotalItemAmount:     group.TotalItemAmount,
		TotalShippingAmount: group.TotalShippingAmount,
		TotalPayAmount:      group.TotalPayAmount,
		Currency:            group.Currency,
		ShopOrderCount:      group.ShopOrderCount,
		ItemCount:           group.ItemCount,
		ShopOrders:          shopOrders,
		CreatedAt:           group.CreatedAt.Unix(),
		UpdatedAt:           group.UpdatedAt.Unix(),
		PaymentDeadlineAt:   group.PaymentDeadlineAt.Unix(),
		PaidAt:              unixOrZero(group.PaidAt),
	}
}

func toOrderGroupDetailPB(group *domainorder.Group) *orderv1.OrderGroupDetail {
	if group == nil {
		return nil
	}

	shopOrders := make([]*orderv1.ShopOrder, 0, len(group.ShopOrders))
	for _, shopOrder := range group.ShopOrders {
		shopOrders = append(shopOrders, toShopOrderPB(shopOrder))
	}

	return &orderv1.OrderGroupDetail{
		Id:                  group.ID,
		UserId:              group.UserID,
		Status:              group.Status,
		Source:              group.Source,
		TotalItemAmount:     group.TotalItemAmount,
		TotalShippingAmount: group.TotalShippingAmount,
		TotalPayAmount:      group.TotalPayAmount,
		Currency:            group.Currency,
		ShopOrderCount:      group.ShopOrderCount,
		ItemCount:           group.ItemCount,
		Address:             toAddressSnapshotPB(group.Address),
		ShopOrders:          shopOrders,
		CreatedAt:           group.CreatedAt.Unix(),
		UpdatedAt:           group.UpdatedAt.Unix(),
		PaymentDeadlineAt:   group.PaymentDeadlineAt.Unix(),
		PaidAt:              unixOrZero(group.PaidAt),
	}
}

func toPaymentInfoPB(info *domainorder.PaymentInfo) *orderv1.OrderGroupPaymentInfo {
	if info == nil {
		return nil
	}
	return &orderv1.OrderGroupPaymentInfo{
		OrderGroupId:      info.OrderGroupID,
		UserId:            info.UserID,
		Status:            info.Status,
		TotalPayAmount:    info.TotalPayAmount,
		Currency:          info.Currency,
		PaymentDeadlineAt: info.PaymentDeadlineAt.Unix(),
		PaidAt:            unixOrZero(info.PaidAt),
	}
}

func unixOrZero(value *time.Time) int64 {
	if value == nil {
		return 0
	}
	return value.Unix()
}

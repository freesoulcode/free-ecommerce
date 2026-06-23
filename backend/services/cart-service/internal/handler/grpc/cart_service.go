package grpc

import (
	"context"
	stderrors "errors"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	applicationcart "github.com/freesoulcode/free-ecommerce/backend/services/cart-service/internal/application/cart"
	domaincart "github.com/freesoulcode/free-ecommerce/backend/services/cart-service/internal/domain/cart"
	cartv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/cart/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartServiceServer struct {
	cartv1.UnimplementedCartServiceServer
	addService    *applicationcart.AddCartItemService
	updateService *applicationcart.UpdateCartItemService
	deleteService *applicationcart.DeleteCartItemService
	listService   *applicationcart.ListCartItemsService
}

func NewCartServiceServer(addService *applicationcart.AddCartItemService, updateService *applicationcart.UpdateCartItemService, deleteService *applicationcart.DeleteCartItemService, listService *applicationcart.ListCartItemsService) *CartServiceServer {
	return &CartServiceServer{addService: addService, updateService: updateService, deleteService: deleteService, listService: listService}
}

func (s *CartServiceServer) AddCartItem(ctx context.Context, req *cartv1.AddCartItemRequest) (*cartv1.AddCartItemResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	item, err := s.addService.Execute(ctx, applicationcart.AddCartItemInput{UserID: req.GetUserId(), SKUID: req.GetSkuId(), Quantity: req.GetQuantity()})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &cartv1.AddCartItemResponse{Item: toItemPB(item)}, nil
}

func (s *CartServiceServer) UpdateCartItem(ctx context.Context, req *cartv1.UpdateCartItemRequest) (*cartv1.UpdateCartItemResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	item, err := s.updateService.Execute(ctx, applicationcart.UpdateCartItemInput{ID: req.GetId(), UserID: req.GetUserId(), Quantity: req.GetQuantity(), Selected: req.GetSelected()})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &cartv1.UpdateCartItemResponse{Item: toItemPB(item)}, nil
}

func (s *CartServiceServer) DeleteCartItem(ctx context.Context, req *cartv1.DeleteCartItemRequest) (*cartv1.DeleteCartItemResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if err := s.deleteService.Execute(ctx, applicationcart.DeleteCartItemInput{ID: req.GetId(), UserID: req.GetUserId()}); err != nil {
		return nil, toGRPCError(err)
	}

	return &cartv1.DeleteCartItemResponse{}, nil
}

func (s *CartServiceServer) ListCartItems(ctx context.Context, req *cartv1.ListCartItemsRequest) (*cartv1.ListCartItemsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	items, err := s.listService.Execute(ctx, req.GetUserId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	result := make([]*cartv1.CartItem, 0, len(items))
	for _, item := range items {
		result = append(result, toItemPB(item))
	}

	return &cartv1.ListCartItemsResponse{Items: result}, nil
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
	default:
		return status.Error(codes.Internal, appErr.Message)
	}
}

func toItemPB(item *domaincart.Item) *cartv1.CartItem {
	if item == nil {
		return nil
	}

	return &cartv1.CartItem{
		Id:                item.ID,
		UserId:            item.UserID,
		SkuId:             item.SKUID,
		ProductId:         item.ProductID,
		ShopId:            item.ShopID,
		ShopName:          item.ShopName,
		ProductTitle:      item.ProductTitle,
		ProductSubTitle:   item.ProductSubTitle,
		MainImageUrl:      item.MainImageURL,
		SkuName:           item.SKUName,
		PriceAmount:       item.PriceAmount,
		Currency:          item.Currency,
		Stock:             item.Stock,
		Quantity:          item.Quantity,
		Selected:          item.Selected,
		ReviewStatus:      item.ReviewStatus,
		ProductSaleStatus: item.ProductSaleStatus,
		SkuSaleStatus:     item.SKUSaleStatus,
		Available:         item.Available,
		CreatedAt:         item.CreatedAt.Unix(),
		UpdatedAt:         item.UpdatedAt.Unix(),
	}
}

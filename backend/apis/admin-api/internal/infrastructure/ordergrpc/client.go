package ordergrpc

import (
	"context"

	applicationadmin "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/application/admin"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	orderv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/order/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn   *grpc.ClientConn
	client orderv1.OrderServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, client: orderv1.NewOrderServiceClient(conn)}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

func (c *Client) ListMerchantShopOrders(ctx context.Context, input applicationadmin.ListShopOrdersInput) (*applicationadmin.ListShopOrdersResult, error) {
	resp, err := c.client.ListMerchantShopOrders(ctx, &orderv1.ListMerchantShopOrdersRequest{
		ShopId:   input.ShopID,
		Page:     input.Page,
		PageSize: input.PageSize,
		Status:   input.Status,
	})
	if err != nil {
		return nil, toAppError(err)
	}

	items := make([]*applicationadmin.MerchantShopOrderSummary, 0, len(resp.GetShopOrders()))
	for _, shopOrder := range resp.GetShopOrders() {
		items = append(items, toMerchantShopOrderSummary(shopOrder))
	}

	return &applicationadmin.ListShopOrdersResult{
		ShopOrders: items,
		Total:      resp.GetTotal(),
		Page:       resp.GetPage(),
		PageSize:   resp.GetPageSize(),
	}, nil
}

func (c *Client) GetMerchantShopOrderDetail(ctx context.Context, shopID, shopOrderID int64) (*applicationadmin.MerchantShopOrderDetail, error) {
	resp, err := c.client.GetMerchantShopOrderDetail(ctx, &orderv1.GetMerchantShopOrderDetailRequest{ShopId: shopID, ShopOrderId: shopOrderID})
	if err != nil {
		return nil, toAppError(err)
	}

	return toMerchantShopOrderDetail(resp.GetShopOrder()), nil
}

func (c *Client) MarkMerchantShopOrderProcessing(ctx context.Context, shopID, shopOrderID int64) (*applicationadmin.MerchantShopOrderDetail, error) {
	resp, err := c.client.MarkMerchantShopOrderProcessing(ctx, &orderv1.MarkMerchantShopOrderProcessingRequest{ShopId: shopID, ShopOrderId: shopOrderID})
	if err != nil {
		return nil, toAppError(err)
	}

	return toMerchantShopOrderDetail(resp.GetShopOrder()), nil
}

func (c *Client) MarkMerchantShopOrderShipped(ctx context.Context, shopID, shopOrderID int64) (*applicationadmin.MerchantShopOrderDetail, error) {
	resp, err := c.client.MarkMerchantShopOrderShipped(ctx, &orderv1.MarkMerchantShopOrderShippedRequest{ShopId: shopID, ShopOrderId: shopOrderID})
	if err != nil {
		return nil, toAppError(err)
	}

	return toMerchantShopOrderDetail(resp.GetShopOrder()), nil
}

func toAppError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return appErrors.Internal("call order service failed")
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return appErrors.InvalidArgument(st.Message())
	case codes.NotFound:
		return appErrors.NotFound(st.Message())
	case codes.Unauthenticated:
		return appErrors.Unauthorized(st.Message())
	default:
		return appErrors.Internal(st.Message())
	}
}

func toAddressSnapshot(address *orderv1.AddressSnapshot) *applicationadmin.OrderAddressSnapshot {
	if address == nil {
		return nil
	}

	return &applicationadmin.OrderAddressSnapshot{
		ID:            address.GetId(),
		OrderGroupID:  address.GetOrderGroupId(),
		UserID:        address.GetUserId(),
		ReceiverName:  address.GetReceiverName(),
		ReceiverPhone: address.GetReceiverPhone(),
		CountryCode:   address.GetCountryCode(),
		Province:      address.GetProvince(),
		City:          address.GetCity(),
		District:      address.GetDistrict(),
		AddressLine1:  address.GetAddressLine1(),
		AddressLine2:  address.GetAddressLine2(),
		PostalCode:    address.GetPostalCode(),
		Tag:           address.GetTag(),
		CreatedAt:     address.GetCreatedAt(),
		UpdatedAt:     address.GetUpdatedAt(),
	}
}

func toOrderItem(item *orderv1.OrderItem) *applicationadmin.OrderItem {
	if item == nil {
		return nil
	}

	return &applicationadmin.OrderItem{
		ID:                        item.GetId(),
		OrderGroupID:              item.GetOrderGroupId(),
		ShopOrderID:               item.GetShopOrderId(),
		UserID:                    item.GetUserId(),
		ShopID:                    item.GetShopId(),
		ProductID:                 item.GetProductId(),
		SKUID:                     item.GetSkuId(),
		ProductTitle:              item.GetProductTitle(),
		ProductSubTitle:           item.GetProductSubTitle(),
		MainImageURL:              item.GetMainImageUrl(),
		SKUName:                   item.GetSkuName(),
		PriceAmount:               item.GetPriceAmount(),
		Currency:                  item.GetCurrency(),
		Quantity:                  item.GetQuantity(),
		ItemAmount:                item.GetItemAmount(),
		ReviewStatusSnapshot:      item.GetReviewStatusSnapshot(),
		ProductSaleStatusSnapshot: item.GetProductSaleStatusSnapshot(),
		SKUSaleStatusSnapshot:     item.GetSkuSaleStatusSnapshot(),
		CreatedAt:                 item.GetCreatedAt(),
		UpdatedAt:                 item.GetUpdatedAt(),
	}
}

func toShopOrder(shopOrder *orderv1.ShopOrder) *applicationadmin.ShopOrder {
	if shopOrder == nil {
		return nil
	}

	items := make([]*applicationadmin.OrderItem, 0, len(shopOrder.GetItems()))
	for _, item := range shopOrder.GetItems() {
		items = append(items, toOrderItem(item))
	}

	return &applicationadmin.ShopOrder{
		ID:             shopOrder.GetId(),
		OrderGroupID:   shopOrder.GetOrderGroupId(),
		UserID:         shopOrder.GetUserId(),
		ShopID:         shopOrder.GetShopId(),
		ShopName:       shopOrder.GetShopName(),
		Status:         shopOrder.GetStatus(),
		ItemAmount:     shopOrder.GetItemAmount(),
		ShippingAmount: shopOrder.GetShippingAmount(),
		PayAmount:      shopOrder.GetPayAmount(),
		Currency:       shopOrder.GetCurrency(),
		ItemCount:      shopOrder.GetItemCount(),
		Items:          items,
		PaidAt:         shopOrder.GetPaidAt(),
		CreatedAt:      shopOrder.GetCreatedAt(),
		UpdatedAt:      shopOrder.GetUpdatedAt(),
	}
}

func toMerchantShopOrderSummary(shopOrder *orderv1.MerchantShopOrderSummary) *applicationadmin.MerchantShopOrderSummary {
	if shopOrder == nil {
		return nil
	}

	return &applicationadmin.MerchantShopOrderSummary{
		ID:                shopOrder.GetId(),
		OrderGroupID:      shopOrder.GetOrderGroupId(),
		UserID:            shopOrder.GetUserId(),
		ShopID:            shopOrder.GetShopId(),
		ShopName:          shopOrder.GetShopName(),
		Status:            shopOrder.GetStatus(),
		ItemAmount:        shopOrder.GetItemAmount(),
		ShippingAmount:    shopOrder.GetShippingAmount(),
		PayAmount:         shopOrder.GetPayAmount(),
		Currency:          shopOrder.GetCurrency(),
		ItemCount:         shopOrder.GetItemCount(),
		PaidAt:            shopOrder.GetPaidAt(),
		CreatedAt:         shopOrder.GetCreatedAt(),
		UpdatedAt:         shopOrder.GetUpdatedAt(),
		OrderGroupStatus:  shopOrder.GetOrderGroupStatus(),
		PaymentDeadlineAt: shopOrder.GetPaymentDeadlineAt(),
	}
}

func toMerchantShopOrderDetail(shopOrder *orderv1.MerchantShopOrderDetail) *applicationadmin.MerchantShopOrderDetail {
	if shopOrder == nil {
		return nil
	}

	return &applicationadmin.MerchantShopOrderDetail{
		OrderGroupID:      shopOrder.GetOrderGroupId(),
		UserID:            shopOrder.GetUserId(),
		OrderGroupStatus:  shopOrder.GetOrderGroupStatus(),
		Source:            shopOrder.GetSource(),
		PaymentDeadlineAt: shopOrder.GetPaymentDeadlineAt(),
		PaidAt:            shopOrder.GetPaidAt(),
		Address:           toAddressSnapshot(shopOrder.GetAddress()),
		ShopOrder:         toShopOrder(shopOrder.GetShopOrder()),
	}
}

package grpc

import (
	"context"
	stderrors "errors"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	applicationproduct "github.com/freesoulcode/free-ecommerce/backend/services/product-service/internal/application/product"
	domainproduct "github.com/freesoulcode/free-ecommerce/backend/services/product-service/internal/domain/product"
	productv1 "github.com/freesoulcode/free-ecommerce/gen/go/proto/product/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductServiceServer struct {
	productv1.UnimplementedProductServiceServer
	listService      *applicationproduct.ListPublicProductsService
	getService       *applicationproduct.GetPublicProductService
	batchSkuService  *applicationproduct.BatchGetSkuBriefsService
	listAdminService *applicationproduct.ListAdminProductsService
	getAdminService  *applicationproduct.GetAdminProductService
	reviewService    *applicationproduct.ReviewProductService
}

func NewProductServiceServer(listService *applicationproduct.ListPublicProductsService, getService *applicationproduct.GetPublicProductService, batchSkuService *applicationproduct.BatchGetSkuBriefsService, listAdminService *applicationproduct.ListAdminProductsService, getAdminService *applicationproduct.GetAdminProductService, reviewService *applicationproduct.ReviewProductService) *ProductServiceServer {
	return &ProductServiceServer{listService: listService, getService: getService, batchSkuService: batchSkuService, listAdminService: listAdminService, getAdminService: getAdminService, reviewService: reviewService}
}

func (s *ProductServiceServer) ListPublicProducts(ctx context.Context, req *productv1.ListPublicProductsRequest) (*productv1.ListPublicProductsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	result, err := s.listService.Execute(ctx, applicationproduct.ListPublicProductsInput{
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
		Keyword:  req.GetKeyword(),
		ShopID:   req.GetShopId(),
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	items := make([]*productv1.ProductSummary, 0, len(result.Products))
	for _, product := range result.Products {
		items = append(items, &productv1.ProductSummary{
			Id:             product.ID,
			ShopId:         product.ShopID,
			ShopName:       product.ShopName,
			Title:          product.Title,
			SubTitle:       product.SubTitle,
			MainImageUrl:   product.MainImageURL,
			MinPriceAmount: product.MinPriceAmount,
			MaxPriceAmount: product.MaxPriceAmount,
			Currency:       product.Currency,
			TotalStock:     product.TotalStock,
		})
	}

	return &productv1.ListPublicProductsResponse{Products: items, Total: result.Total, Page: result.Page, PageSize: result.PageSize}, nil
}

func (s *ProductServiceServer) GetPublicProduct(ctx context.Context, req *productv1.GetPublicProductRequest) (*productv1.GetPublicProductResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	product, err := s.getService.Execute(ctx, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	skus := make([]*productv1.ProductSku, 0, len(product.SKUs))
	for _, sku := range product.SKUs {
		skus = append(skus, &productv1.ProductSku{Id: sku.ID, Name: sku.Name, PriceAmount: sku.PriceAmount, Currency: sku.Currency, Stock: sku.Stock, SaleStatus: sku.SaleStatus})
	}

	return &productv1.GetPublicProductResponse{Product: &productv1.ProductDetail{
		Id:             product.ID,
		ShopId:         product.ShopID,
		ShopName:       product.ShopName,
		Title:          product.Title,
		SubTitle:       product.SubTitle,
		MainImageUrl:   product.MainImageURL,
		Description:    product.Description,
		ReviewStatus:   product.ReviewStatus,
		SaleStatus:     product.SaleStatus,
		MinPriceAmount: product.MinPriceAmount,
		MaxPriceAmount: product.MaxPriceAmount,
		Currency:       product.Currency,
		TotalStock:     product.TotalStock,
		Skus:           skus,
	}}, nil
}

func (s *ProductServiceServer) BatchGetSkuBriefs(ctx context.Context, req *productv1.BatchGetSkuBriefsRequest) (*productv1.BatchGetSkuBriefsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	briefs, err := s.batchSkuService.Execute(ctx, req.GetSkuIds())
	if err != nil {
		return nil, toGRPCError(err)
	}

	items := make([]*productv1.SkuBrief, 0, len(briefs))
	for _, brief := range briefs {
		items = append(items, &productv1.SkuBrief{
			SkuId:             brief.SKUID,
			ProductId:         brief.ProductID,
			ShopId:            brief.ShopID,
			ShopName:          brief.ShopName,
			ProductTitle:      brief.ProductTitle,
			ProductSubTitle:   brief.ProductSubTitle,
			MainImageUrl:      brief.MainImageURL,
			SkuName:           brief.SKUName,
			PriceAmount:       brief.PriceAmount,
			Currency:          brief.Currency,
			Stock:             brief.Stock,
			ReviewStatus:      brief.ReviewStatus,
			ProductSaleStatus: brief.ProductSaleStatus,
			SkuSaleStatus:     brief.SKUSaleStatus,
		})
	}

	return &productv1.BatchGetSkuBriefsResponse{Skus: items}, nil
}

func (s *ProductServiceServer) ListAdminProducts(ctx context.Context, req *productv1.ListAdminProductsRequest) (*productv1.ListAdminProductsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	result, err := s.listAdminService.Execute(ctx, applicationproduct.ListAdminProductsInput{
		Page:         req.GetPage(),
		PageSize:     req.GetPageSize(),
		Keyword:      req.GetKeyword(),
		ShopID:       req.GetShopId(),
		ReviewStatus: req.GetReviewStatus(),
		SaleStatus:   req.GetSaleStatus(),
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	items := make([]*productv1.AdminProductSummary, 0, len(result.Products))
	for _, product := range result.Products {
		items = append(items, &productv1.AdminProductSummary{
			Id:             product.ID,
			ShopId:         product.ShopID,
			ShopName:       product.ShopName,
			Title:          product.Title,
			SubTitle:       product.SubTitle,
			MainImageUrl:   product.MainImageURL,
			MinPriceAmount: product.MinPriceAmount,
			MaxPriceAmount: product.MaxPriceAmount,
			Currency:       product.Currency,
			TotalStock:     product.TotalStock,
			ReviewStatus:   product.ReviewStatus,
			SaleStatus:     product.SaleStatus,
			CreatedAt:      product.CreatedAt.Unix(),
			UpdatedAt:      product.UpdatedAt.Unix(),
		})
	}

	return &productv1.ListAdminProductsResponse{Products: items, Total: result.Total, Page: result.Page, PageSize: result.PageSize}, nil
}

func (s *ProductServiceServer) GetAdminProduct(ctx context.Context, req *productv1.GetAdminProductRequest) (*productv1.GetAdminProductResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	product, err := s.getAdminService.Execute(ctx, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &productv1.GetAdminProductResponse{Product: toProductDetailPB(product)}, nil
}

func (s *ProductServiceServer) ReviewProduct(ctx context.Context, req *productv1.ReviewProductRequest) (*productv1.ReviewProductResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	product, err := s.reviewService.Execute(ctx, req.GetId(), req.GetReviewStatus())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &productv1.ReviewProductResponse{Product: toProductDetailPB(product)}, nil
}

func toProductDetailPB(product *domainproduct.Detail) *productv1.ProductDetail {
	if product == nil {
		return nil
	}
	skus := make([]*productv1.ProductSku, 0, len(product.SKUs))
	for _, sku := range product.SKUs {
		skus = append(skus, &productv1.ProductSku{Id: sku.ID, Name: sku.Name, PriceAmount: sku.PriceAmount, Currency: sku.Currency, Stock: sku.Stock, SaleStatus: sku.SaleStatus})
	}
	return &productv1.ProductDetail{
		Id:             product.ID,
		ShopId:         product.ShopID,
		ShopName:       product.ShopName,
		Title:          product.Title,
		SubTitle:       product.SubTitle,
		MainImageUrl:   product.MainImageURL,
		Description:    product.Description,
		ReviewStatus:   product.ReviewStatus,
		SaleStatus:     product.SaleStatus,
		MinPriceAmount: product.MinPriceAmount,
		MaxPriceAmount: product.MaxPriceAmount,
		Currency:       product.Currency,
		TotalStock:     product.TotalStock,
		Skus:           skus,
	}
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

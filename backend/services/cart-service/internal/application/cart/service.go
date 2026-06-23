package cart

import (
	"context"
	stderrors "errors"
	"time"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	domaincart "github.com/freesoulcode/free-ecommerce/backend/services/cart-service/internal/domain/cart"
)

type ProductCatalog interface {
	BatchGetSkuBriefs(ctx context.Context, skuIDs []int64) (map[int64]*domaincart.SkuBrief, error)
}

type AddCartItemInput struct {
	UserID   int64
	SKUID    int64
	Quantity int32
}

type UpdateCartItemInput struct {
	ID       int64
	UserID   int64
	Quantity int32
	Selected bool
}

type DeleteCartItemInput struct {
	ID     int64
	UserID int64
}

type AddCartItemService struct {
	repo        domaincart.Repository
	idGenerator IDGenerator
	product     ProductCatalog
	now         func() time.Time
}

type UpdateCartItemService struct {
	repo    domaincart.Repository
	product ProductCatalog
	now     func() time.Time
}

type DeleteCartItemService struct {
	repo domaincart.Repository
}

type ListCartItemsService struct {
	repo    domaincart.Repository
	product ProductCatalog
}

func NewAddCartItemService(repo domaincart.Repository, idGenerator IDGenerator, product ProductCatalog, now func() time.Time) *AddCartItemService {
	if now == nil {
		now = time.Now
	}

	return &AddCartItemService{repo: repo, idGenerator: idGenerator, product: product, now: now}
}

func NewUpdateCartItemService(repo domaincart.Repository, product ProductCatalog, now func() time.Time) *UpdateCartItemService {
	if now == nil {
		now = time.Now
	}

	return &UpdateCartItemService{repo: repo, product: product, now: now}
}

func NewDeleteCartItemService(repo domaincart.Repository) *DeleteCartItemService {
	return &DeleteCartItemService{repo: repo}
}

func NewListCartItemsService(repo domaincart.Repository, product ProductCatalog) *ListCartItemsService {
	return &ListCartItemsService{repo: repo, product: product}
}

func (s *AddCartItemService) Execute(ctx context.Context, input AddCartItemInput) (*domaincart.Item, error) {
	if input.UserID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}
	if input.SKUID <= 0 {
		return nil, appErrors.InvalidArgument("sku id is required")
	}
	if input.Quantity <= 0 {
		return nil, appErrors.InvalidArgument("quantity must be greater than zero")
	}

	brief, err := s.getRequiredAvailableBrief(ctx, input.SKUID)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindItemByUserAndSKU(ctx, input.UserID, input.SKUID)
	if err != nil && !isNotFound(err) {
		return nil, err
	}

	now := s.now().UTC()
	if existing != nil {
		newQuantity := existing.Quantity + input.Quantity
		if newQuantity > brief.Stock {
			return nil, appErrors.InvalidArgument("quantity exceeds stock")
		}
		if err := existing.Update(newQuantity, existing.Selected, now); err != nil {
			return nil, err
		}
		if err := s.repo.UpdateItem(ctx, existing); err != nil {
			return nil, err
		}
		existing.AttachBrief(brief)
		return existing, nil
	}

	if input.Quantity > brief.Stock {
		return nil, appErrors.InvalidArgument("quantity exceeds stock")
	}

	itemID, err := s.idGenerator.NextID()
	if err != nil {
		return nil, appErrors.Internal("generate cart item id failed")
	}
	item, err := domaincart.NewItem(itemID, input.UserID, input.SKUID, input.Quantity, true, now)
	if err != nil {
		return nil, err
	}
	if err := s.repo.CreateItem(ctx, item); err != nil {
		return nil, err
	}
	item.AttachBrief(brief)
	return item, nil
}

func (s *UpdateCartItemService) Execute(ctx context.Context, input UpdateCartItemInput) (*domaincart.Item, error) {
	if input.ID <= 0 {
		return nil, appErrors.InvalidArgument("cart item id is required")
	}
	if input.UserID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}

	item, err := s.repo.FindItemByID(ctx, input.UserID, input.ID)
	if err != nil {
		return nil, err
	}

	briefs, err := s.product.BatchGetSkuBriefs(ctx, []int64{item.SKUID})
	if err != nil {
		return nil, err
	}
	brief := briefs[item.SKUID]
	if brief != nil && brief.IsAvailable() {
		if input.Quantity > brief.Stock {
			return nil, appErrors.InvalidArgument("quantity exceeds stock")
		}
	} else if input.Selected {
		return nil, appErrors.InvalidArgument("sku is not available")
	}

	if err := item.Update(input.Quantity, input.Selected, s.now().UTC()); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateItem(ctx, item); err != nil {
		return nil, err
	}
	item.AttachBrief(brief)
	return item, nil
}

func (s *DeleteCartItemService) Execute(ctx context.Context, input DeleteCartItemInput) error {
	if input.ID <= 0 {
		return appErrors.InvalidArgument("cart item id is required")
	}
	if input.UserID <= 0 {
		return appErrors.InvalidArgument("user id is required")
	}

	return s.repo.DeleteItem(ctx, input.UserID, input.ID)
}

func (s *ListCartItemsService) Execute(ctx context.Context, userID int64) ([]*domaincart.Item, error) {
	if userID <= 0 {
		return nil, appErrors.InvalidArgument("user id is required")
	}

	items, err := s.repo.ListItemsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	skuIDs := make([]int64, 0, len(items))
	for _, item := range items {
		skuIDs = append(skuIDs, item.SKUID)
	}
	briefs, err := s.product.BatchGetSkuBriefs(ctx, skuIDs)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		item.AttachBrief(briefs[item.SKUID])
	}
	return items, nil
}

func (s *AddCartItemService) getRequiredAvailableBrief(ctx context.Context, skuID int64) (*domaincart.SkuBrief, error) {
	briefs, err := s.product.BatchGetSkuBriefs(ctx, []int64{skuID})
	if err != nil {
		return nil, err
	}
	brief := briefs[skuID]
	if brief == nil {
		return nil, appErrors.NotFound("sku not found")
	}
	if !brief.IsAvailable() {
		return nil, appErrors.InvalidArgument("sku is not available")
	}
	return brief, nil
}

func isNotFound(err error) bool {
	var appErr *appErrors.Error
	if !stderrors.As(err, &appErr) {
		return false
	}

	return appErr.Code == appErrors.CodeNotFound
}

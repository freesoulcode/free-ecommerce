package http

import (
	"strconv"

	applicationadmin "github.com/freesoulcode/free-ecommerce/backend/apis/admin-api/internal/application/admin"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type AdminUserHandler struct {
	service *applicationadmin.UserAdminService
}

func NewAdminUserHandler(service *applicationadmin.UserAdminService) *AdminUserHandler {
	return &AdminUserHandler{service: service}
}

func (h *AdminUserHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/v1/admin/users/:userID", h.get)
	router.GET("/api/v1/admin/users/:userID/addresses", h.listAddresses)
}

func (h *AdminUserHandler) get(c *gin.Context) {
	userID, ok := parseInt64PathParam(c, "userID")
	if !ok {
		return
	}

	user, err := h.service.Get(c.Request.Context(), userID)
	if err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{"user": userResponse(user)})
}

func (h *AdminUserHandler) listAddresses(c *gin.Context) {
	userID, ok := parseInt64PathParam(c, "userID")
	if !ok {
		return
	}

	addresses, err := h.service.ListAddresses(c.Request.Context(), userID)
	if err != nil {
		httpx.Error(c, err)
		return
	}

	items := make([]gin.H, 0, len(addresses))
	for _, address := range addresses {
		items = append(items, addressResponse(address))
	}

	httpx.OK(c, gin.H{"addresses": items})
}

func parseInt64PathParam(c *gin.Context, key string) (int64, bool) {
	value, err := strconv.ParseInt(c.Param(key), 10, 64)
	if err != nil || value <= 0 {
		httpx.Error(c, appErrors.InvalidArgument("invalid path parameter"))
		return 0, false
	}

	return value, true
}

func userResponse(user *applicationadmin.User) gin.H {
	if user == nil {
		return nil
	}

	return gin.H{
		"id":             user.ID,
		"email":          user.Email,
		"phone":          user.Phone,
		"nickname":       user.Nickname,
		"status":         user.Status,
		"email_verified": user.EmailVerified,
		"phone_verified": user.PhoneVerified,
	}
}

func addressResponse(address *applicationadmin.Address) gin.H {
	if address == nil {
		return nil
	}

	return gin.H{
		"id":             address.ID,
		"user_id":        address.UserID,
		"receiver_name":  address.ReceiverName,
		"receiver_phone": address.ReceiverPhone,
		"country_code":   address.CountryCode,
		"province":       address.Province,
		"city":           address.City,
		"district":       address.District,
		"address_line1":  address.AddressLine1,
		"address_line2":  address.AddressLine2,
		"postal_code":    address.PostalCode,
		"tag":            address.Tag,
		"is_default":     address.IsDefault,
		"created_at":     address.CreatedAt,
		"updated_at":     address.UpdatedAt,
	}
}

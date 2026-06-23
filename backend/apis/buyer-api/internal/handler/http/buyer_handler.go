package http

import (
	"net/http"
	"strconv"

	applicationbuyer "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/application/buyer"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type BuyerHandler struct {
	registerBuyerService *applicationbuyer.RegisterBuyerService
	loginBuyerService    *applicationbuyer.LoginBuyerService
	addressBuyerService  *applicationbuyer.AddressBuyerService
}

func NewBuyerHandler(registerBuyerService *applicationbuyer.RegisterBuyerService, loginBuyerService *applicationbuyer.LoginBuyerService, addressBuyerService *applicationbuyer.AddressBuyerService) *BuyerHandler {
	return &BuyerHandler{registerBuyerService: registerBuyerService, loginBuyerService: loginBuyerService, addressBuyerService: addressBuyerService}
}

type registerBuyerRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginBuyerRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	DeviceID string `json:"device_id"`
}

type refreshBuyerTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	DeviceID     string `json:"device_id"`
}

type logoutBuyerRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type createBuyerAddressRequest struct {
	ReceiverName  string `json:"receiver_name" binding:"required"`
	ReceiverPhone string `json:"receiver_phone" binding:"required"`
	CountryCode   string `json:"country_code" binding:"required"`
	Province      string `json:"province" binding:"required"`
	City          string `json:"city" binding:"required"`
	District      string `json:"district" binding:"required"`
	AddressLine1  string `json:"address_line1" binding:"required"`
	AddressLine2  string `json:"address_line2"`
	PostalCode    string `json:"postal_code"`
	Tag           string `json:"tag"`
	IsDefault     bool   `json:"is_default"`
}

type updateBuyerAddressRequest = createBuyerAddressRequest

func (h *BuyerHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/api/v1/buyers/register", h.register)
	router.POST("/api/v1/buyers/login", h.login)
	router.POST("/api/v1/buyers/token/refresh", h.refreshToken)
	router.POST("/api/v1/buyers/logout", h.logout)
	router.GET("/api/v1/buyers/:buyerID/addresses", h.listAddresses)
	router.POST("/api/v1/buyers/:buyerID/addresses", h.createAddress)
	router.PUT("/api/v1/buyers/:buyerID/addresses/:addressID", h.updateAddress)
	router.DELETE("/api/v1/buyers/:buyerID/addresses/:addressID", h.deleteAddress)
}

func (h *BuyerHandler) register(c *gin.Context) {
	var req registerBuyerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, appErrors.InvalidArgument("invalid request body"))
		return
	}

	buyer, err := h.registerBuyerService.Execute(c.Request.Context(), applicationbuyer.CreateBuyerInput{
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Password: req.Password,
	})
	if err != nil {
		httpx.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, httpx.Response{
		Code:    "OK",
		Message: "success",
		Data: gin.H{
			"id":             buyer.ID,
			"email":          buyer.Email,
			"phone":          buyer.Phone,
			"nickname":       buyer.Nickname,
			"status":         buyer.Status,
			"email_verified": buyer.EmailVerified,
			"phone_verified": buyer.PhoneVerified,
		},
		RequestID: c.GetHeader(httpx.HeaderRequestID),
		TraceID:   c.GetHeader(httpx.HeaderTraceID),
	})
}

func (h *BuyerHandler) login(c *gin.Context) {
	var req loginBuyerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, appErrors.InvalidArgument("invalid request body"))
		return
	}

	result, err := h.loginBuyerService.Execute(c.Request.Context(), applicationbuyer.LoginBuyerInput{
		Email:     req.Email,
		Password:  req.Password,
		DeviceID:  req.DeviceID,
		UserAgent: c.GetHeader("User-Agent"),
		ClientIP:  c.ClientIP(),
	})
	if err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{
		"buyer": gin.H{
			"id":             result.Buyer.ID,
			"email":          result.Buyer.Email,
			"phone":          result.Buyer.Phone,
			"nickname":       result.Buyer.Nickname,
			"status":         result.Buyer.Status,
			"email_verified": result.Buyer.EmailVerified,
			"phone_verified": result.Buyer.PhoneVerified,
		},
		"access_token":             result.AccessToken,
		"refresh_token":            result.RefreshToken,
		"token_type":               result.TokenType,
		"access_token_expires_at":  result.AccessTokenExpiresAt,
		"refresh_token_expires_at": result.RefreshTokenExpiresAt,
		"refresh_session_id":       result.RefreshSessionID,
	})
}

func (h *BuyerHandler) refreshToken(c *gin.Context) {
	var req refreshBuyerTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, appErrors.InvalidArgument("invalid request body"))
		return
	}

	result, err := h.loginBuyerService.Refresh(c.Request.Context(), applicationbuyer.RefreshBuyerTokenInput{
		RefreshToken: req.RefreshToken,
		DeviceID:     req.DeviceID,
		UserAgent:    c.GetHeader("User-Agent"),
		ClientIP:     c.ClientIP(),
	})
	if err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{
		"buyer": gin.H{
			"id":             result.Buyer.ID,
			"email":          result.Buyer.Email,
			"phone":          result.Buyer.Phone,
			"nickname":       result.Buyer.Nickname,
			"status":         result.Buyer.Status,
			"email_verified": result.Buyer.EmailVerified,
			"phone_verified": result.Buyer.PhoneVerified,
		},
		"access_token":             result.AccessToken,
		"refresh_token":            result.RefreshToken,
		"token_type":               result.TokenType,
		"access_token_expires_at":  result.AccessTokenExpiresAt,
		"refresh_token_expires_at": result.RefreshTokenExpiresAt,
		"refresh_session_id":       result.RefreshSessionID,
	})
}

func (h *BuyerHandler) logout(c *gin.Context) {
	var req logoutBuyerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, appErrors.InvalidArgument("invalid request body"))
		return
	}

	result, err := h.loginBuyerService.Logout(c.Request.Context(), applicationbuyer.LogoutBuyerInput{RefreshToken: req.RefreshToken})
	if err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{"refresh_session_id": result.RefreshSessionID})
}

func (h *BuyerHandler) listAddresses(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}

	addresses, err := h.addressBuyerService.List(c.Request.Context(), buyerID)
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

func (h *BuyerHandler) createAddress(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}

	var req createBuyerAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, appErrors.InvalidArgument("invalid request body"))
		return
	}

	address, err := h.addressBuyerService.Create(c.Request.Context(), applicationbuyer.CreateAddressInput{
		UserID:        buyerID,
		ReceiverName:  req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		CountryCode:   req.CountryCode,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		AddressLine1:  req.AddressLine1,
		AddressLine2:  req.AddressLine2,
		PostalCode:    req.PostalCode,
		Tag:           req.Tag,
		IsDefault:     req.IsDefault,
	})
	if err != nil {
		httpx.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, httpx.Response{Code: "OK", Message: "success", Data: gin.H{"address": addressResponse(address)}, RequestID: c.GetHeader(httpx.HeaderRequestID), TraceID: c.GetHeader(httpx.HeaderTraceID)})
}

func (h *BuyerHandler) updateAddress(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}
	addressID, ok := parseInt64PathParam(c, "addressID")
	if !ok {
		return
	}

	var req updateBuyerAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, appErrors.InvalidArgument("invalid request body"))
		return
	}

	address, err := h.addressBuyerService.Update(c.Request.Context(), applicationbuyer.UpdateAddressInput{
		ID:            addressID,
		UserID:        buyerID,
		ReceiverName:  req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		CountryCode:   req.CountryCode,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		AddressLine1:  req.AddressLine1,
		AddressLine2:  req.AddressLine2,
		PostalCode:    req.PostalCode,
		Tag:           req.Tag,
		IsDefault:     req.IsDefault,
	})
	if err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{"address": addressResponse(address)})
}

func (h *BuyerHandler) deleteAddress(c *gin.Context) {
	buyerID, ok := parseInt64PathParam(c, "buyerID")
	if !ok {
		return
	}
	addressID, ok := parseInt64PathParam(c, "addressID")
	if !ok {
		return
	}

	if err := h.addressBuyerService.Delete(c.Request.Context(), applicationbuyer.DeleteAddressInput{ID: addressID, UserID: buyerID}); err != nil {
		httpx.Error(c, err)
		return
	}

	httpx.OK(c, gin.H{"deleted": true})
}

func parseInt64PathParam(c *gin.Context, key string) (int64, bool) {
	value, err := strconv.ParseInt(c.Param(key), 10, 64)
	if err != nil || value <= 0 {
		httpx.Error(c, appErrors.InvalidArgument("invalid path parameter"))
		return 0, false
	}

	return value, true
}

func addressResponse(address *applicationbuyer.Address) gin.H {
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

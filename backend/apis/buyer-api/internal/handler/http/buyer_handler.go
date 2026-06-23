package http

import (
	"net/http"

	applicationbuyer "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/application/buyer"
	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type BuyerHandler struct {
	registerBuyerService *applicationbuyer.RegisterBuyerService
	loginBuyerService    *applicationbuyer.LoginBuyerService
}

func NewBuyerHandler(registerBuyerService *applicationbuyer.RegisterBuyerService, loginBuyerService *applicationbuyer.LoginBuyerService) *BuyerHandler {
	return &BuyerHandler{registerBuyerService: registerBuyerService, loginBuyerService: loginBuyerService}
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

func (h *BuyerHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/api/v1/buyers/register", h.register)
	router.POST("/api/v1/buyers/login", h.login)
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

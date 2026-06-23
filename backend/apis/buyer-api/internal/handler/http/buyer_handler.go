package http

import (
	"net/http"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"
	"github.com/freesoulcode/free-ecommerce/backend/pkg/httpx"
	applicationbuyer "github.com/freesoulcode/free-ecommerce/backend/apis/buyer-api/internal/application/buyer"
	"github.com/gin-gonic/gin"
)

type BuyerHandler struct {
	registerBuyerService *applicationbuyer.RegisterBuyerService
}

func NewBuyerHandler(registerBuyerService *applicationbuyer.RegisterBuyerService) *BuyerHandler {
	return &BuyerHandler{registerBuyerService: registerBuyerService}
}

type registerBuyerRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname" binding:"required"`
}

func (h *BuyerHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/api/v1/buyers/register", h.register)
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

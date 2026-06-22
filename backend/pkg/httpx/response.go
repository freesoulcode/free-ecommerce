package httpx

import (
	"errors"
	"net/http"

	appErrors "github.com/freesoulcode/free-ecommerce/backend/pkg/errors"

	"github.com/gin-gonic/gin"
)

const (
	HeaderRequestID = "X-Request-Id"
	HeaderTraceID   = "X-Trace-Id"
)

type Response struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	TraceID   string `json:"trace_id,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:      string(appErrors.CodeOK),
		Message:   "success",
		Data:      data,
		RequestID: c.GetHeader(HeaderRequestID),
		TraceID:   c.GetHeader(HeaderTraceID),
	})
}

func Error(c *gin.Context, err error) {
	var appErr *appErrors.Error
	if !errors.As(err, &appErr) {
		appErr = appErrors.Internal("internal server error")
	}
	c.JSON(appErr.HTTPStatus, Response{
		Code:      string(appErr.Code),
		Message:   appErr.Message,
		RequestID: c.GetHeader(HeaderRequestID),
		TraceID:   c.GetHeader(HeaderTraceID),
	})
}

package errors

import "net/http"

// Code is a stable business error code exposed by HTTP APIs.
type Code string

const (
	CodeOK              Code = "OK"
	CodeInternal        Code = "INTERNAL_ERROR"
	CodeInvalidArgument Code = "INVALID_ARGUMENT"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeNotFound        Code = "NOT_FOUND"
)

// Error represents a business error with HTTP mapping.
type Error struct {
	Code       Code
	Message    string
	HTTPStatus int
}

func (e *Error) Error() string {
	return e.Message
}

func New(code Code, message string, httpStatus int) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

func Internal(message string) *Error {
	return New(CodeInternal, message, http.StatusInternalServerError)
}

func InvalidArgument(message string) *Error {
	return New(CodeInvalidArgument, message, http.StatusBadRequest)
}

func Unauthorized(message string) *Error {
	return New(CodeUnauthorized, message, http.StatusUnauthorized)
}

func NotFound(message string) *Error {
	return New(CodeNotFound, message, http.StatusNotFound)
}

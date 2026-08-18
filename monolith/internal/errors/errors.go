package errors

import (
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(status int, code, msg string) *AppError {
	return &AppError{
		StatusCode: status,
		Code:       code,
		Message:    msg,
	}
}

// Define the error statndard that mostly use
var (
	ErrInternalFailure = NewAppError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Some went wrong, please try again..")
	ErrBadRequest      = NewAppError(http.StatusInternalServerError, "BAD_REQUEST", "Bad request, please check the request payload")
	ErrUnauthorized    = NewAppError(http.StatusInternalServerError, "UNAUTHORIZED", "Unauthorized, please check the credentials")
	ErrNotFound        = NewAppError(http.StatusInternalServerError, "NOT_FOUND", "Resource not found")
)

package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error with HTTP status code
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error constructors

// BadRequest creates a 400 error
func BadRequest(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, err)
}

// Unauthorized creates a 401 error
func Unauthorized(message string, err error) *AppError {
	return NewAppError(http.StatusUnauthorized, message, err)
}

// Forbidden creates a 403 error
func Forbidden(message string, err error) *AppError {
	return NewAppError(http.StatusForbidden, message, err)
}

// NotFound creates a 404 error
func NotFound(message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, message, err)
}

// InternalServerError creates a 500 error
func InternalServerError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, err)
}

// Conflict creates a 409 error
func Conflict(message string, err error) *AppError {
	return NewAppError(http.StatusConflict, message, err)
}

// ValidationError creates a 422 error
func ValidationError(message string, err error) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, message, err)
}

// InsufficientStorage creates a 507 error
func InsufficientStorage(message string, err error) *AppError {
	return NewAppError(http.StatusInsufficientStorage, message, err)
}

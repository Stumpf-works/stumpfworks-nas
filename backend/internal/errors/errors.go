package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	statusCode int
	code       string
	message    string
	err        error
}

func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.message, e.err)
	}
	return e.message
}

func (e *AppError) StatusCode() int {
	return e.statusCode
}

func (e *AppError) Code() string {
	return e.code
}

func (e *AppError) Message() string {
	return e.message
}

// New creates a new AppError
func New(statusCode int, code, message string, err error) *AppError {
	return &AppError{
		statusCode: statusCode,
		code:       code,
		message:    message,
		err:        err,
	}
}

// BadRequest creates a 400 error
func BadRequest(message string, err error) *AppError {
	return New(http.StatusBadRequest, "bad_request", message, err)
}

// NotFound creates a 404 error
func NotFound(message string, err error) *AppError {
	return New(http.StatusNotFound, "not_found", message, err)
}

// InternalServerError creates a 500 error
func InternalServerError(message string, err error) *AppError {
	return New(http.StatusInternalServerError, "internal_server_error", message, err)
}

// Unauthorized creates a 401 error
func Unauthorized(message string, err error) *AppError {
	return New(http.StatusUnauthorized, "unauthorized", message, err)
}

// Forbidden creates a 403 error
func Forbidden(message string, err error) *AppError {
	return New(http.StatusForbidden, "forbidden", message, err)
}

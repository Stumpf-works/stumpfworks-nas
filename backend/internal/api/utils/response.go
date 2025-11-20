package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RespondSuccess sends a successful JSON response
func RespondSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
	})
}

// RespondError sends an error JSON response
func RespondError(w http.ResponseWriter, err interface{}) {
	// Support both *errors.AppError and error types
	var statusCode int
	var code, message string

	switch e := err.(type) {
	case *AppError:
		statusCode = e.statusCode
		code = e.code
		message = e.message
	case error:
		statusCode = http.StatusInternalServerError
		code = "internal_server_error"
		message = e.Error()
	default:
		statusCode = http.StatusInternalServerError
		code = "internal_server_error"
		message = "Unknown error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

// AppError represents an application error
type AppError struct {
	statusCode int
	code       string
	message    string
	err        error
}

func (e *AppError) Error() string {
	if e.err != nil {
		return e.err.Error()
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

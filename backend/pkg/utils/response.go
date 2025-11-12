package utils

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents error information in the response
type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RespondJSON writes a JSON response
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: statusCode >= 200 && statusCode < 300,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode JSON response", zap.Error(err))
	}
}

// RespondError writes an error JSON response
func RespondError(w http.ResponseWriter, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		appErr = errors.InternalServerError("Internal server error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)

	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		logger.Error("Failed to encode error response", zap.Error(encodeErr))
	}

	// Log the original error
	if appErr.Code >= 500 {
		logger.Error("Server error", zap.Error(appErr), zap.Int("status", appErr.Code))
	} else {
		logger.Warn("Client error", zap.String("message", appErr.Message), zap.Int("status", appErr.Code))
	}
}

// RespondSuccess writes a success JSON response with data
func RespondSuccess(w http.ResponseWriter, data interface{}) {
	RespondJSON(w, http.StatusOK, data)
}

// RespondCreated writes a 201 created response
func RespondCreated(w http.ResponseWriter, data interface{}) {
	RespondJSON(w, http.StatusCreated, data)
}

// RespondNoContent writes a 204 no content response
func RespondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

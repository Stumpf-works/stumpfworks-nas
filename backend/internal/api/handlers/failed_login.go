// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Stumpf-works/stumpfworks-nas/internal/auth"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

// FailedLoginHandler handles failed login-related HTTP requests
type FailedLoginHandler struct {
	service *auth.FailedLoginService
}

// NewFailedLoginHandler creates a new failed login handler
func NewFailedLoginHandler() *FailedLoginHandler {
	return &FailedLoginHandler{
		service: auth.GetFailedLoginService(),
	}
}

// ListFailedAttempts retrieves failed login attempts
func (h *FailedLoginHandler) ListFailedAttempts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Pagination
	limit := 50
	offset := 0

	if limitStr := query.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsedOffset
		}
	}

	// Get attempts
	attempts, total, err := h.service.GetRecentFailedAttempts(limit, offset)
	if err != nil {
		logger.Error("Failed to retrieve failed attempts", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to retrieve failed attempts", err))
		return
	}

	// Return response with pagination metadata
	utils.RespondSuccess(w, map[string]interface{}{
		"attempts": attempts,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetBlockedIPs retrieves all currently blocked IPs
func (h *FailedLoginHandler) GetBlockedIPs(w http.ResponseWriter, r *http.Request) {
	blocks, err := h.service.GetBlockedIPs()
	if err != nil {
		logger.Error("Failed to retrieve blocked IPs", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to retrieve blocked IPs", err))
		return
	}

	utils.RespondSuccess(w, blocks)
}

// UnblockIP removes the block on an IP address
func (h *FailedLoginHandler) UnblockIP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IPAddress string `json:"ipAddress"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.IPAddress == "" {
		utils.RespondError(w, errors.BadRequest("IP address is required", nil))
		return
	}

	if err := h.service.UnblockIP(r.Context(), req.IPAddress); err != nil {
		logger.Error("Failed to unblock IP",
			zap.Error(err),
			zap.String("ip", req.IPAddress))
		utils.RespondError(w, errors.InternalServerError("Failed to unblock IP", err))
		return
	}

	logger.Info("IP address unblocked", zap.String("ip", req.IPAddress))
	utils.RespondSuccess(w, map[string]string{
		"message": "IP address unblocked successfully",
	})
}

// GetStats retrieves failed login statistics
func (h *FailedLoginHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats()
	if err != nil {
		logger.Error("Failed to retrieve failed login stats", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to retrieve statistics", err))
		return
	}

	utils.RespondSuccess(w, stats)
}

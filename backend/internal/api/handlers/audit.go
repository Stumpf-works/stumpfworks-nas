// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/audit"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

// AuditHandler handles audit log-related HTTP requests
type AuditHandler struct {
	service *audit.Service
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler() *AuditHandler {
	return &AuditHandler{
		service: audit.GetService(),
	}
}

// ListAuditLogs retrieves audit logs with filtering and pagination
func (h *AuditHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Parse query parameters
	params := &audit.QueryParams{}

	// User ID filter
	if userIDStr := query.Get("userId"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uid := uint(userID)
			params.UserID = &uid
		}
	}

	// Username filter
	params.Username = query.Get("username")

	// Action filter
	params.Action = query.Get("action")

	// Status filter
	params.Status = query.Get("status")

	// Severity filter
	params.Severity = query.Get("severity")

	// Date range filters
	if startDateStr := query.Get("startDate"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			params.StartDate = &startDate
		}
	}

	if endDateStr := query.Get("endDate"); endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			params.EndDate = &endDate
		}
	}

	// Pagination
	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			params.Limit = limit
		}
	} else {
		params.Limit = 100 // Default limit
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			params.Offset = offset
		}
	}

	// Query audit logs
	logs, total, err := h.service.Query(r.Context(), params)
	if err != nil {
		logger.Error("Failed to query audit logs", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to retrieve audit logs", err))
		return
	}

	// Return response with pagination metadata
	utils.RespondSuccess(w, map[string]interface{}{
		"logs":   logs,
		"total":  total,
		"limit":  params.Limit,
		"offset": params.Offset,
	})
}

// GetAuditLog retrieves a specific audit log by ID
func (h *AuditHandler) GetAuditLog(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path (assuming chi router)
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.RespondError(w, errors.BadRequest("Audit log ID is required", nil))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid audit log ID", err))
		return
	}

	log, err := h.service.GetByID(r.Context(), uint(id))
	if err != nil {
		logger.Error("Failed to retrieve audit log",
			zap.Error(err),
			zap.Uint64("id", id))
		utils.RespondError(w, errors.NotFound("Audit log not found", err))
		return
	}

	utils.RespondSuccess(w, log)
}

// GetRecentAuditLogs retrieves the most recent audit logs
func (h *AuditHandler) GetRecentAuditLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit := 50 // Default limit

	if limitStr := query.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs, err := h.service.GetRecent(r.Context(), limit)
	if err != nil {
		logger.Error("Failed to retrieve recent audit logs", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to retrieve recent audit logs", err))
		return
	}

	utils.RespondSuccess(w, logs)
}

// GetAuditStats retrieves audit log statistics
func (h *AuditHandler) GetAuditStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		logger.Error("Failed to retrieve audit stats", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to retrieve audit statistics", err))
		return
	}

	utils.RespondSuccess(w, stats)
}

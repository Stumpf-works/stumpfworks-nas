package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Stumpf-works/stumpfworks-nas/internal/alerts"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

// AlertHandler handles alert-related API requests
type AlertHandler struct {
	alertService *alerts.Service
}

// NewAlertHandler creates a new alert handler
func NewAlertHandler() *AlertHandler {
	return &AlertHandler{
		alertService: alerts.GetService(),
	}
}

// GetConfig retrieves the alert configuration
func (h *AlertHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	config, err := h.alertService.GetConfig(ctx)
	if err != nil {
		logger.Error("Failed to get alert config", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get alert config", err))
		return
	}

	utils.RespondSuccess(w, config)
}

// UpdateConfig updates the alert configuration
func (h *AlertHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var config models.AlertConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.alertService.UpdateConfig(ctx, &config); err != nil {
		logger.Error("Failed to update alert config", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to update alert config", err))
		return
	}

	// Fetch the updated config to return
	updatedConfig, err := h.alertService.GetConfig(ctx)
	if err != nil {
		logger.Error("Failed to fetch updated config", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to fetch updated config", err))
		return
	}

	utils.RespondSuccess(w, updatedConfig)
}

// TestEmail sends a test email
func (h *AlertHandler) TestEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var config models.AlertConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.alertService.TestEmail(ctx, &config); err != nil {
		logger.Error("Failed to send test email", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to send test email", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Test email sent successfully",
	})
}

// TestWebhook sends a test webhook
func (h *AlertHandler) TestWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var config models.AlertConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.alertService.TestWebhook(ctx, &config); err != nil {
		logger.Error("Failed to send test webhook", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to send test webhook", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Test webhook sent successfully",
	})
}

// GetAlertLogs retrieves recent alert logs
func (h *AlertHandler) GetAlertLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get limit from query params (default 50)
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs, err := h.alertService.GetAlertLogs(ctx, limit)
	if err != nil {
		logger.Error("Failed to get alert logs", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get alert logs", err))
		return
	}

	utils.RespondSuccess(w, logs)
}

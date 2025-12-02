// Revision: 2025-12-02 | Author: Claude | Version: 1.2.0
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/ups"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

// UPSHandler handles UPS-related HTTP requests
type UPSHandler struct {
	service *ups.Service
}

// NewUPSHandler creates a new UPS handler
func NewUPSHandler() *UPSHandler {
	return &UPSHandler{
		service: ups.GetService(),
	}
}

// CheckAvailability middleware to check if UPS service is available
func (h *UPSHandler) CheckAvailability(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.service == nil {
			utils.RespondError(w, errors.NewAppError(503, "UPS service is not available", nil))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetConfig returns the current UPS configuration
func (h *UPSHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	config := h.service.GetConfig()
	if config == nil {
		utils.RespondError(w, errors.NotFound("UPS configuration not found", nil))
		return
	}

	utils.RespondSuccess(w, config)
}

// UpdateConfig updates the UPS configuration
func (h *UPSHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	var config models.UPSConfig

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate configuration
	if config.UPSName == "" {
		utils.RespondError(w, errors.BadRequest("UPS name is required", nil))
		return
	}

	if config.UPSHost == "" {
		utils.RespondError(w, errors.BadRequest("UPS host is required", nil))
		return
	}

	if config.UPSPort <= 0 || config.UPSPort > 65535 {
		utils.RespondError(w, errors.BadRequest("Invalid UPS port", nil))
		return
	}

	if config.PollInterval < 10 || config.PollInterval > 300 {
		utils.RespondError(w, errors.BadRequest("Poll interval must be between 10 and 300 seconds", nil))
		return
	}

	if config.LowBatteryThreshold < 5 || config.LowBatteryThreshold > 50 {
		utils.RespondError(w, errors.BadRequest("Low battery threshold must be between 5 and 50 percent", nil))
		return
	}

	if config.ShutdownDelay < 0 || config.ShutdownDelay > 600 {
		utils.RespondError(w, errors.BadRequest("Shutdown delay must be between 0 and 600 seconds", nil))
		return
	}

	// Save configuration
	if err := h.service.SaveConfig(&config); err != nil {
		logger.Error("Failed to save UPS config", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to save configuration", err))
		return
	}

	logger.Info("UPS configuration updated",
		zap.String("ups_name", config.UPSName),
		zap.Bool("enabled", config.Enabled))

	utils.RespondSuccess(w, config)
}

// GetStatus returns the current UPS status
func (h *UPSHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.service.GetStatus()
	if err != nil {
		logger.Error("Failed to get UPS status", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get UPS status", err))
		return
	}

	utils.RespondSuccess(w, status)
}

// TestConnection tests the UPS connection with provided config
func (h *UPSHandler) TestConnection(w http.ResponseWriter, r *http.Request) {
	var config models.UPSConfig

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	status, err := h.service.TestUPS(&config)
	if err != nil {
		logger.Error("UPS test connection failed", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Connection test failed", err))
		return
	}

	logger.Info("UPS test connection successful", zap.String("ups_name", config.UPSName))
	utils.RespondSuccess(w, status)
}

// GetEvents returns UPS event history
func (h *UPSHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 100 // Default limit
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	events, err := h.service.GetEvents(limit, offset)
	if err != nil {
		logger.Error("Failed to get UPS events", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get events", err))
		return
	}

	utils.RespondSuccess(w, events)
}

// StartMonitoring starts UPS monitoring
func (h *UPSHandler) StartMonitoring(w http.ResponseWriter, r *http.Request) {
	if err := h.service.StartMonitoring(); err != nil {
		logger.Error("Failed to start UPS monitoring", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to start monitoring", err))
		return
	}

	logger.Info("UPS monitoring started")
	utils.RespondSuccess(w, map[string]string{"message": "Monitoring started successfully"})
}

// StopMonitoring stops UPS monitoring
func (h *UPSHandler) StopMonitoring(w http.ResponseWriter, r *http.Request) {
	h.service.StopMonitoring()

	logger.Info("UPS monitoring stopped")
	utils.RespondSuccess(w, map[string]string{"message": "Monitoring stopped successfully"})
}

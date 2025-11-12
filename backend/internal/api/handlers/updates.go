package handlers

import (
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/updates"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

// UpdateHandler handles update-related HTTP requests
type UpdateHandler struct {
	service *updates.UpdateService
}

// NewUpdateHandler creates a new update handler
func NewUpdateHandler() *UpdateHandler {
	return &UpdateHandler{
		service: updates.GetService(),
	}
}

// CheckForUpdates checks for available updates
func (h *UpdateHandler) CheckForUpdates(w http.ResponseWriter, r *http.Request) {
	// Check if force refresh is requested
	forceCheck := r.URL.Query().Get("force") == "true"

	result, err := h.service.CheckForUpdates(r.Context(), forceCheck)
	if err != nil {
		logger.Error("Failed to check for updates", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to check for updates", err))
		return
	}

	logger.Info("Update check completed",
		zap.Bool("update_available", result.UpdateAvailable),
		zap.String("current", result.CurrentVersion),
		zap.String("latest", result.LatestVersion))

	utils.RespondSuccess(w, result)
}

// GetCurrentVersion returns the current version
func (h *UpdateHandler) GetCurrentVersion(w http.ResponseWriter, r *http.Request) {
	version := h.service.GetCurrentVersion()
	utils.RespondSuccess(w, map[string]string{
		"version": version,
	})
}

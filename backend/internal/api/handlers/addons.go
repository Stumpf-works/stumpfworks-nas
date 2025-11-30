package handlers

import (
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/addons"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var addonManager *addons.Manager

// InitAddonManager initializes the addon manager
func InitAddonManager(manager *addons.Manager) {
	addonManager = manager
	logger.Info("Addon manager initialized in handlers")
}

// ListAddons lists all available addons with their installation status
func ListAddons(w http.ResponseWriter, r *http.Request) {
	if addonManager == nil {
		utils.RespondError(w, errors.InternalServerError("Addon manager not initialized", nil))
		return
	}

	addonsWithStatus, err := addonManager.GetAllAddonsWithStatus()
	if err != nil {
		logger.Error("Failed to get addons with status", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get addons", err))
		return
	}

	utils.RespondSuccess(w, addonsWithStatus)
}

// GetAddon gets details of a specific addon
func GetAddon(w http.ResponseWriter, r *http.Request) {
	if addonManager == nil {
		utils.RespondError(w, errors.InternalServerError("Addon manager not initialized", nil))
		return
	}

	addonID := chi.URLParam(r, "id")
	if addonID == "" {
		utils.RespondError(w, errors.BadRequest("Addon ID is required", nil))
		return
	}

	addon, err := addonManager.GetAddon(addonID)
	if err != nil {
		logger.Error("Failed to get addon", zap.Error(err), zap.String("addon_id", addonID))
		utils.RespondError(w, errors.NotFound("Addon not found", err))
		return
	}

	// Get status
	status, err := addonManager.GetAddonStatus(addonID)
	if err != nil {
		logger.Warn("Failed to get addon status", zap.Error(err), zap.String("addon_id", addonID))
		status = &addons.InstallationStatus{
			AddonID:   addonID,
			Installed: false,
		}
	}

	response := map[string]interface{}{
		"manifest": addon,
		"status":   status,
	}

	utils.RespondSuccess(w, response)
}

// GetAddonStatus gets the installation status of an addon
func GetAddonStatus(w http.ResponseWriter, r *http.Request) {
	if addonManager == nil {
		utils.RespondError(w, errors.InternalServerError("Addon manager not initialized", nil))
		return
	}

	addonID := chi.URLParam(r, "id")
	if addonID == "" {
		utils.RespondError(w, errors.BadRequest("Addon ID is required", nil))
		return
	}

	status, err := addonManager.GetAddonStatus(addonID)
	if err != nil {
		logger.Error("Failed to get addon status", zap.Error(err), zap.String("addon_id", addonID))
		utils.RespondError(w, errors.InternalServerError("Failed to get addon status", err))
		return
	}

	utils.RespondSuccess(w, status)
}

// InstallAddon installs an addon
func InstallAddon(w http.ResponseWriter, r *http.Request) {
	if addonManager == nil {
		utils.RespondError(w, errors.InternalServerError("Addon manager not initialized", nil))
		return
	}

	addonID := chi.URLParam(r, "id")
	if addonID == "" {
		utils.RespondError(w, errors.BadRequest("Addon ID is required", nil))
		return
	}

	logger.Info("Installing addon via API", zap.String("addon_id", addonID))

	// Get addon manifest to check if restart is required
	addon, err := addonManager.GetAddon(addonID)
	if err != nil {
		logger.Error("Failed to get addon", zap.Error(err), zap.String("addon_id", addonID))
		utils.RespondError(w, errors.InternalServerError("Failed to get addon", err))
		return
	}

	if err := addonManager.InstallAddon(addonID); err != nil {
		logger.Error("Failed to install addon", zap.Error(err), zap.String("addon_id", addonID))
		utils.RespondError(w, errors.InternalServerError("Failed to install addon", err))
		return
	}

	logger.Info("Addon installed successfully via API", zap.String("addon_id", addonID))

	// Schedule service restart if addon requires it
	if addon.RequiresRestart {
		logger.Info("Addon requires service restart, scheduling restart", zap.String("addon_id", addonID))
		addonManager.ScheduleServiceRestart()

		utils.RespondSuccess(w, map[string]string{
			"message":           "Addon installed successfully. Service will restart in 3 seconds to initialize addon.",
			"addon_id":          addonID,
			"restart_scheduled": "true",
		})
	} else {
		utils.RespondSuccess(w, map[string]string{
			"message":  "Addon installed successfully",
			"addon_id": addonID,
		})
	}
}

// UninstallAddon uninstalls an addon
func UninstallAddon(w http.ResponseWriter, r *http.Request) {
	if addonManager == nil {
		utils.RespondError(w, errors.InternalServerError("Addon manager not initialized", nil))
		return
	}

	addonID := chi.URLParam(r, "id")
	if addonID == "" {
		utils.RespondError(w, errors.BadRequest("Addon ID is required", nil))
		return
	}

	logger.Info("Uninstalling addon via API", zap.String("addon_id", addonID))

	if err := addonManager.UninstallAddon(addonID); err != nil {
		logger.Error("Failed to uninstall addon", zap.Error(err), zap.String("addon_id", addonID))
		utils.RespondError(w, errors.InternalServerError("Failed to uninstall addon", err))
		return
	}

	logger.Info("Addon uninstalled successfully via API", zap.String("addon_id", addonID))
	utils.RespondSuccess(w, map[string]string{
		"message":  "Addon uninstalled successfully",
		"addon_id": addonID,
	})
}

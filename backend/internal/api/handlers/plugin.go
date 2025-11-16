package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/plugins"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// PluginHandler handles plugin-related HTTP requests
type PluginHandler struct {
	service *plugins.Service
	runtime *plugins.Runtime
}

// NewPluginHandler creates a new plugin handler
func NewPluginHandler() *PluginHandler {
	svc := plugins.GetService()
	return &PluginHandler{
		service: svc,
		runtime: plugins.NewRuntime(svc),
	}
}

// CheckAvailability middleware checks if plugin service is available
func (h *PluginHandler) CheckAvailability(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.service == nil {
			utils.RespondError(w, errors.NewAppError(
				http.StatusServiceUnavailable,
				"Plugin service is not available",
				nil,
			))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ListPlugins lists all plugins
func (h *PluginHandler) ListPlugins(w http.ResponseWriter, r *http.Request) {
	plugins, err := h.service.ListPlugins(r.Context())
	if err != nil {
		logger.Error("Failed to list plugins", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list plugins", err))
		return
	}

	utils.RespondSuccess(w, plugins)
}

// GetPlugin gets a specific plugin by ID
func (h *PluginHandler) GetPlugin(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "id")

	plugin, err := h.service.GetPlugin(r.Context(), pluginID)
	if err != nil {
		logger.Error("Failed to get plugin", zap.Error(err), zap.String("pluginID", pluginID))
		utils.RespondError(w, errors.NotFound("Plugin not found", err))
		return
	}

	utils.RespondSuccess(w, plugin)
}

// InstallPlugin installs a new plugin
func (h *PluginHandler) InstallPlugin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SourcePath string `json:"sourcePath"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.SourcePath == "" {
		utils.RespondError(w, errors.BadRequest("Source path is required", nil))
		return
	}

	plugin, err := h.service.InstallPlugin(r.Context(), req.SourcePath)
	if err != nil {
		logger.Error("Failed to install plugin", zap.Error(err), zap.String("sourcePath", req.SourcePath))
		utils.RespondError(w, errors.InternalServerError("Failed to install plugin", err))
		return
	}

	logger.Info("Plugin installed", zap.String("pluginID", plugin.ID), zap.String("name", plugin.Name))
	utils.RespondSuccess(w, plugin)
}

// UninstallPlugin uninstalls a plugin
func (h *PluginHandler) UninstallPlugin(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "id")

	if err := h.service.UninstallPlugin(r.Context(), pluginID); err != nil {
		logger.Error("Failed to uninstall plugin", zap.Error(err), zap.String("pluginID", pluginID))
		utils.RespondError(w, errors.InternalServerError("Failed to uninstall plugin", err))
		return
	}

	logger.Info("Plugin uninstalled", zap.String("pluginID", pluginID))
	utils.RespondSuccess(w, map[string]string{"message": "Plugin uninstalled successfully"})
}

// EnablePlugin enables a plugin
func (h *PluginHandler) EnablePlugin(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "id")

	if err := h.service.EnablePlugin(r.Context(), pluginID); err != nil {
		logger.Error("Failed to enable plugin", zap.Error(err), zap.String("pluginID", pluginID))
		utils.RespondError(w, errors.InternalServerError("Failed to enable plugin", err))
		return
	}

	logger.Info("Plugin enabled", zap.String("pluginID", pluginID))
	utils.RespondSuccess(w, map[string]string{"message": "Plugin enabled successfully"})
}

// DisablePlugin disables a plugin
func (h *PluginHandler) DisablePlugin(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "id")

	if err := h.service.DisablePlugin(r.Context(), pluginID); err != nil {
		logger.Error("Failed to disable plugin", zap.Error(err), zap.String("pluginID", pluginID))
		utils.RespondError(w, errors.InternalServerError("Failed to disable plugin", err))
		return
	}

	logger.Info("Plugin disabled", zap.String("pluginID", pluginID))
	utils.RespondSuccess(w, map[string]string{"message": "Plugin disabled successfully"})
}

// UpdatePluginConfig updates a plugin's configuration
func (h *PluginHandler) UpdatePluginConfig(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "id")

	var req struct {
		Config map[string]interface{} `json:"config"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.service.UpdatePluginConfig(r.Context(), pluginID, req.Config); err != nil {
		logger.Error("Failed to update plugin config", zap.Error(err), zap.String("pluginID", pluginID))
		utils.RespondError(w, errors.InternalServerError("Failed to update plugin config", err))
		return
	}

	logger.Info("Plugin config updated", zap.String("pluginID", pluginID))
	utils.RespondSuccess(w, map[string]string{"message": "Plugin config updated successfully"})
}

// StartPlugin starts a plugin's runtime execution
func (h *PluginHandler) StartPlugin(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "id")

	if err := h.runtime.StartPlugin(r.Context(), pluginID); err != nil {
		logger.Error("Failed to start plugin", zap.Error(err), zap.String("pluginID", pluginID))
		utils.RespondError(w, errors.InternalServerError("Failed to start plugin", err))
		return
	}

	logger.Info("Plugin started", zap.String("pluginID", pluginID))
	utils.RespondSuccess(w, map[string]string{"message": "Plugin started successfully"})
}

// StopPlugin stops a running plugin
func (h *PluginHandler) StopPlugin(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "id")

	if err := h.runtime.StopPlugin(r.Context(), pluginID); err != nil {
		logger.Error("Failed to stop plugin", zap.Error(err), zap.String("pluginID", pluginID))
		utils.RespondError(w, errors.InternalServerError("Failed to stop plugin", err))
		return
	}

	logger.Info("Plugin stopped", zap.String("pluginID", pluginID))
	utils.RespondSuccess(w, map[string]string{"message": "Plugin stopped successfully"})
}

// RestartPlugin restarts a plugin
func (h *PluginHandler) RestartPlugin(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "id")

	if err := h.runtime.RestartPlugin(r.Context(), pluginID); err != nil {
		logger.Error("Failed to restart plugin", zap.Error(err), zap.String("pluginID", pluginID))
		utils.RespondError(w, errors.InternalServerError("Failed to restart plugin", err))
		return
	}

	logger.Info("Plugin restarted", zap.String("pluginID", pluginID))
	utils.RespondSuccess(w, map[string]string{"message": "Plugin restarted successfully"})
}

// GetPluginStatus returns the runtime status of a plugin
func (h *PluginHandler) GetPluginStatus(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "id")

	status, err := h.runtime.GetPluginStatus(pluginID)
	if err != nil {
		// Plugin not running, return basic status
		utils.RespondSuccess(w, map[string]interface{}{
			"pluginID": pluginID,
			"status":   "stopped",
			"running":  false,
		})
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"pluginID":  status.PluginID,
		"status":    status.Status,
		"running":   status.Status == "running",
		"startedAt": status.StartedAt,
		"error":     status.LastError,
	})
}

// ListRunningPlugins returns all currently running plugins
func (h *PluginHandler) ListRunningPlugins(w http.ResponseWriter, r *http.Request) {
	procs := h.runtime.ListRunningPlugins()

	running := make([]map[string]interface{}, len(procs))
	for i, proc := range procs {
		running[i] = map[string]interface{}{
			"pluginID":  proc.PluginID,
			"status":    proc.Status,
			"startedAt": proc.StartedAt,
		}
	}

	utils.RespondSuccess(w, running)
}

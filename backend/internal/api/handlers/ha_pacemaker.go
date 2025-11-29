package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/ha"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var pacemakerManager *ha.PacemakerManager

// InitPacemakerManager initializes the Pacemaker manager
func InitPacemakerManager(manager *ha.PacemakerManager) {
	pacemakerManager = manager
	logger.Info("Pacemaker manager initialized in handlers")
}

// GetClusterStatus gets the current cluster status
func GetClusterStatus(w http.ResponseWriter, r *http.Request) {
	if pacemakerManager == nil || !pacemakerManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Pacemaker service not available",
			nil,
		))
		return
	}

	status, err := pacemakerManager.GetClusterStatus()
	if err != nil {
		logger.Error("Failed to get cluster status", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get cluster status", err))
		return
	}

	utils.RespondSuccess(w, status)
}

// CreateClusterResource creates a new cluster resource
func CreateClusterResource(w http.ResponseWriter, r *http.Request) {
	if pacemakerManager == nil || !pacemakerManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Pacemaker service not available",
			nil,
		))
		return
	}

	var config ha.ResourceConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate required fields
	if config.ID == "" || config.Agent == "" {
		utils.RespondError(w, errors.BadRequest("ID and agent are required", nil))
		return
	}

	if err := pacemakerManager.CreateResource(config); err != nil {
		logger.Error("Failed to create cluster resource", zap.Error(err), zap.String("id", config.ID))
		utils.RespondError(w, errors.InternalServerError("Failed to create cluster resource", err))
		return
	}

	logger.Info("Cluster resource created", zap.String("id", config.ID))
	utils.RespondSuccess(w, map[string]string{
		"message": "Cluster resource created successfully",
		"id":      config.ID,
	})
}

// DeleteClusterResource deletes a cluster resource
func DeleteClusterResource(w http.ResponseWriter, r *http.Request) {
	if pacemakerManager == nil || !pacemakerManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Pacemaker service not available",
			nil,
		))
		return
	}

	resourceID := chi.URLParam(r, "id")
	if resourceID == "" {
		utils.RespondError(w, errors.BadRequest("Resource ID is required", nil))
		return
	}

	if err := pacemakerManager.DeleteResource(resourceID); err != nil {
		logger.Error("Failed to delete cluster resource", zap.Error(err), zap.String("id", resourceID))
		utils.RespondError(w, errors.InternalServerError("Failed to delete cluster resource", err))
		return
	}

	logger.Info("Cluster resource deleted", zap.String("id", resourceID))
	utils.RespondSuccess(w, map[string]string{
		"message": "Cluster resource deleted successfully",
		"id":      resourceID,
	})
}

// EnableClusterResource enables a cluster resource
func EnableClusterResource(w http.ResponseWriter, r *http.Request) {
	if pacemakerManager == nil || !pacemakerManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Pacemaker service not available",
			nil,
		))
		return
	}

	resourceID := chi.URLParam(r, "id")
	if resourceID == "" {
		utils.RespondError(w, errors.BadRequest("Resource ID is required", nil))
		return
	}

	if err := pacemakerManager.EnableResource(resourceID); err != nil {
		logger.Error("Failed to enable cluster resource", zap.Error(err), zap.String("id", resourceID))
		utils.RespondError(w, errors.InternalServerError("Failed to enable cluster resource", err))
		return
	}

	logger.Info("Cluster resource enabled", zap.String("id", resourceID))
	utils.RespondSuccess(w, map[string]string{
		"message": "Resource enabled successfully",
		"id":      resourceID,
	})
}

// DisableClusterResource disables a cluster resource
func DisableClusterResource(w http.ResponseWriter, r *http.Request) {
	if pacemakerManager == nil || !pacemakerManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Pacemaker service not available",
			nil,
		))
		return
	}

	resourceID := chi.URLParam(r, "id")
	if resourceID == "" {
		utils.RespondError(w, errors.BadRequest("Resource ID is required", nil))
		return
	}

	if err := pacemakerManager.DisableResource(resourceID); err != nil {
		logger.Error("Failed to disable cluster resource", zap.Error(err), zap.String("id", resourceID))
		utils.RespondError(w, errors.InternalServerError("Failed to disable cluster resource", err))
		return
	}

	logger.Info("Cluster resource disabled", zap.String("id", resourceID))
	utils.RespondSuccess(w, map[string]string{
		"message": "Resource disabled successfully",
		"id":      resourceID,
	})
}

// MoveClusterResource moves a resource to a specific node
func MoveClusterResource(w http.ResponseWriter, r *http.Request) {
	if pacemakerManager == nil || !pacemakerManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Pacemaker service not available",
			nil,
		))
		return
	}

	resourceID := chi.URLParam(r, "id")
	if resourceID == "" {
		utils.RespondError(w, errors.BadRequest("Resource ID is required", nil))
		return
	}

	var req struct {
		TargetNode string `json:"target_node"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.TargetNode == "" {
		utils.RespondError(w, errors.BadRequest("Target node is required", nil))
		return
	}

	if err := pacemakerManager.MoveResource(resourceID, req.TargetNode); err != nil {
		logger.Error("Failed to move cluster resource", zap.Error(err), zap.String("id", resourceID))
		utils.RespondError(w, errors.InternalServerError("Failed to move cluster resource", err))
		return
	}

	logger.Info("Cluster resource moved", zap.String("id", resourceID), zap.String("node", req.TargetNode))
	utils.RespondSuccess(w, map[string]string{
		"message": "Resource moved successfully",
		"id":      resourceID,
		"node":    req.TargetNode,
	})
}

// ClearClusterResource clears the failed state of a resource
func ClearClusterResource(w http.ResponseWriter, r *http.Request) {
	if pacemakerManager == nil || !pacemakerManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Pacemaker service not available",
			nil,
		))
		return
	}

	resourceID := chi.URLParam(r, "id")
	if resourceID == "" {
		utils.RespondError(w, errors.BadRequest("Resource ID is required", nil))
		return
	}

	if err := pacemakerManager.ClearResource(resourceID); err != nil {
		logger.Error("Failed to clear cluster resource", zap.Error(err), zap.String("id", resourceID))
		utils.RespondError(w, errors.InternalServerError("Failed to clear cluster resource", err))
		return
	}

	logger.Info("Cluster resource cleared", zap.String("id", resourceID))
	utils.RespondSuccess(w, map[string]string{
		"message": "Resource cleared successfully",
		"id":      resourceID,
	})
}

// SetMaintenanceMode enables or disables maintenance mode
func SetMaintenanceMode(w http.ResponseWriter, r *http.Request) {
	if pacemakerManager == nil || !pacemakerManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Pacemaker service not available",
			nil,
		))
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := pacemakerManager.SetMaintenanceMode(req.Enabled); err != nil {
		logger.Error("Failed to set maintenance mode", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to set maintenance mode", err))
		return
	}

	logger.Info("Maintenance mode changed", zap.Bool("enabled", req.Enabled))
	utils.RespondSuccess(w, map[string]interface{}{
		"message": "Maintenance mode updated successfully",
		"enabled": req.Enabled,
	})
}

// StandbyNode puts a node in standby mode
func StandbyNode(w http.ResponseWriter, r *http.Request) {
	if pacemakerManager == nil || !pacemakerManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Pacemaker service not available",
			nil,
		))
		return
	}

	nodeName := chi.URLParam(r, "name")
	if nodeName == "" {
		utils.RespondError(w, errors.BadRequest("Node name is required", nil))
		return
	}

	if err := pacemakerManager.StandbyNode(nodeName); err != nil {
		logger.Error("Failed to put node in standby", zap.Error(err), zap.String("node", nodeName))
		utils.RespondError(w, errors.InternalServerError("Failed to put node in standby", err))
		return
	}

	logger.Info("Node put in standby", zap.String("node", nodeName))
	utils.RespondSuccess(w, map[string]string{
		"message": "Node put in standby successfully",
		"node":    nodeName,
	})
}

// UnstandbyNode removes a node from standby mode
func UnstandbyNode(w http.ResponseWriter, r *http.Request) {
	if pacemakerManager == nil || !pacemakerManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Pacemaker service not available",
			nil,
		))
		return
	}

	nodeName := chi.URLParam(r, "name")
	if nodeName == "" {
		utils.RespondError(w, errors.BadRequest("Node name is required", nil))
		return
	}

	if err := pacemakerManager.UnstandbyNode(nodeName); err != nil {
		logger.Error("Failed to remove node from standby", zap.Error(err), zap.String("node", nodeName))
		utils.RespondError(w, errors.InternalServerError("Failed to remove node from standby", err))
		return
	}

	logger.Info("Node removed from standby", zap.String("node", nodeName))
	utils.RespondSuccess(w, map[string]string{
		"message": "Node removed from standby successfully",
		"node":    nodeName,
	})
}

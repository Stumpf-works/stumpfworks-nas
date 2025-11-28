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

var drbdManager *ha.DRBDManager

// InitDRBDManager initializes the DRBD manager
func InitDRBDManager(manager *ha.DRBDManager) {
	drbdManager = manager
	logger.Info("DRBD manager initialized in handlers")
}

// ListDRBDResources lists all DRBD resources
func ListDRBDResources(w http.ResponseWriter, r *http.Request) {
	if drbdManager == nil || !drbdManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"DRBD service not available",
			nil,
		))
		return
	}

	resources, err := drbdManager.ListResources()
	if err != nil {
		logger.Error("Failed to list DRBD resources", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list DRBD resources", err))
		return
	}

	utils.RespondSuccess(w, resources)
}

// GetDRBDResourceStatus gets the status of a specific DRBD resource
func GetDRBDResourceStatus(w http.ResponseWriter, r *http.Request) {
	if drbdManager == nil || !drbdManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"DRBD service not available",
			nil,
		))
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		utils.RespondError(w, errors.BadRequest("Resource name is required", nil))
		return
	}

	status, err := drbdManager.GetResourceStatus(name)
	if err != nil {
		logger.Error("Failed to get DRBD resource status", zap.Error(err), zap.String("name", name))
		utils.RespondError(w, errors.InternalServerError("Failed to get resource status", err))
		return
	}

	utils.RespondSuccess(w, status)
}

// CreateDRBDResource creates a new DRBD resource
func CreateDRBDResource(w http.ResponseWriter, r *http.Request) {
	if drbdManager == nil || !drbdManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"DRBD service not available",
			nil,
		))
		return
	}

	var resource ha.DRBDResource
	if err := json.NewDecoder(r.Body).Decode(&resource); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate required fields
	if resource.Name == "" || resource.Device == "" || resource.Disk == "" {
		utils.RespondError(w, errors.BadRequest("Name, device, and disk are required", nil))
		return
	}

	if err := drbdManager.CreateResource(resource); err != nil {
		logger.Error("Failed to create DRBD resource", zap.Error(err), zap.String("name", resource.Name))
		utils.RespondError(w, errors.InternalServerError("Failed to create DRBD resource", err))
		return
	}

	logger.Info("DRBD resource created", zap.String("name", resource.Name))
	utils.RespondSuccess(w, map[string]string{
		"message": "DRBD resource created successfully",
		"name":    resource.Name,
	})
}

// DeleteDRBDResource deletes a DRBD resource
func DeleteDRBDResource(w http.ResponseWriter, r *http.Request) {
	if drbdManager == nil || !drbdManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"DRBD service not available",
			nil,
		))
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		utils.RespondError(w, errors.BadRequest("Resource name is required", nil))
		return
	}

	if err := drbdManager.DeleteResource(name); err != nil {
		logger.Error("Failed to delete DRBD resource", zap.Error(err), zap.String("name", name))
		utils.RespondError(w, errors.InternalServerError("Failed to delete DRBD resource", err))
		return
	}

	logger.Info("DRBD resource deleted", zap.String("name", name))
	utils.RespondSuccess(w, map[string]string{
		"message": "DRBD resource deleted successfully",
		"name":    name,
	})
}

// PromoteDRBDResource promotes a DRBD resource to primary
func PromoteDRBDResource(w http.ResponseWriter, r *http.Request) {
	if drbdManager == nil || !drbdManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"DRBD service not available",
			nil,
		))
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		utils.RespondError(w, errors.BadRequest("Resource name is required", nil))
		return
	}

	if err := drbdManager.PromoteToPrimary(name); err != nil {
		logger.Error("Failed to promote DRBD resource", zap.Error(err), zap.String("name", name))
		utils.RespondError(w, errors.InternalServerError("Failed to promote to primary", err))
		return
	}

	logger.Info("DRBD resource promoted to primary", zap.String("name", name))
	utils.RespondSuccess(w, map[string]string{
		"message": "Resource promoted to primary",
		"name":    name,
	})
}

// DemoteDRBDResource demotes a DRBD resource to secondary
func DemoteDRBDResource(w http.ResponseWriter, r *http.Request) {
	if drbdManager == nil || !drbdManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"DRBD service not available",
			nil,
		))
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		utils.RespondError(w, errors.BadRequest("Resource name is required", nil))
		return
	}

	if err := drbdManager.DemoteToSecondary(name); err != nil {
		logger.Error("Failed to demote DRBD resource", zap.Error(err), zap.String("name", name))
		utils.RespondError(w, errors.InternalServerError("Failed to demote to secondary", err))
		return
	}

	logger.Info("DRBD resource demoted to secondary", zap.String("name", name))
	utils.RespondSuccess(w, map[string]string{
		"message": "Resource demoted to secondary",
		"name":    name,
	})
}

// ForcePrimaryDRBDResource forces a DRBD resource to become primary
func ForcePrimaryDRBDResource(w http.ResponseWriter, r *http.Request) {
	if drbdManager == nil || !drbdManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"DRBD service not available",
			nil,
		))
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		utils.RespondError(w, errors.BadRequest("Resource name is required", nil))
		return
	}

	if err := drbdManager.ForcePrimary(name); err != nil {
		logger.Error("Failed to force primary DRBD resource", zap.Error(err), zap.String("name", name))
		utils.RespondError(w, errors.InternalServerError("Failed to force primary", err))
		return
	}

	logger.Warn("DRBD resource forced to primary", zap.String("name", name))
	utils.RespondSuccess(w, map[string]string{
		"message": "Resource forced to primary - check for split-brain",
		"name":    name,
	})
}

// DisconnectDRBDResource disconnects a DRBD resource from its peer
func DisconnectDRBDResource(w http.ResponseWriter, r *http.Request) {
	if drbdManager == nil || !drbdManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"DRBD service not available",
			nil,
		))
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		utils.RespondError(w, errors.BadRequest("Resource name is required", nil))
		return
	}

	if err := drbdManager.Disconnect(name); err != nil {
		logger.Error("Failed to disconnect DRBD resource", zap.Error(err), zap.String("name", name))
		utils.RespondError(w, errors.InternalServerError("Failed to disconnect", err))
		return
	}

	logger.Info("DRBD resource disconnected", zap.String("name", name))
	utils.RespondSuccess(w, map[string]string{
		"message": "Resource disconnected",
		"name":    name,
	})
}

// ConnectDRBDResource connects a DRBD resource to its peer
func ConnectDRBDResource(w http.ResponseWriter, r *http.Request) {
	if drbdManager == nil || !drbdManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"DRBD service not available",
			nil,
		))
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		utils.RespondError(w, errors.BadRequest("Resource name is required", nil))
		return
	}

	if err := drbdManager.Connect(name); err != nil {
		logger.Error("Failed to connect DRBD resource", zap.Error(err), zap.String("name", name))
		utils.RespondError(w, errors.InternalServerError("Failed to connect", err))
		return
	}

	logger.Info("DRBD resource connected", zap.String("name", name))
	utils.RespondSuccess(w, map[string]string{
		"message": "Resource connected",
		"name":    name,
	})
}

// StartDRBDSync starts synchronization for a DRBD resource
func StartDRBDSync(w http.ResponseWriter, r *http.Request) {
	if drbdManager == nil || !drbdManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"DRBD service not available",
			nil,
		))
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		utils.RespondError(w, errors.BadRequest("Resource name is required", nil))
		return
	}

	if err := drbdManager.StartSync(name); err != nil {
		logger.Error("Failed to start DRBD sync", zap.Error(err), zap.String("name", name))
		utils.RespondError(w, errors.InternalServerError("Failed to start sync", err))
		return
	}

	logger.Info("DRBD synchronization started", zap.String("name", name))
	utils.RespondSuccess(w, map[string]string{
		"message": "Synchronization started",
		"name":    name,
	})
}

// VerifyDRBDData verifies data integrity of a DRBD resource
func VerifyDRBDData(w http.ResponseWriter, r *http.Request) {
	if drbdManager == nil || !drbdManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"DRBD service not available",
			nil,
		))
		return
	}

	name := chi.URLParam(r, "name")
	if name == "" {
		utils.RespondError(w, errors.BadRequest("Resource name is required", nil))
		return
	}

	if err := drbdManager.VerifyData(name); err != nil {
		logger.Error("Failed to verify DRBD data", zap.Error(err), zap.String("name", name))
		utils.RespondError(w, errors.InternalServerError("Failed to verify data", err))
		return
	}

	logger.Info("DRBD data verification started", zap.String("name", name))
	utils.RespondSuccess(w, map[string]string{
		"message": "Data verification started",
		"name":    name,
	})
}

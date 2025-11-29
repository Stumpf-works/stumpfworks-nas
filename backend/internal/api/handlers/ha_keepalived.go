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

var keepalivedManager *ha.KeepalivedManager

// InitKeepalivedManager initializes the Keepalived manager
func InitKeepalivedManager(manager *ha.KeepalivedManager) {
	keepalivedManager = manager
	logger.Info("Keepalived manager initialized in handlers")
}

// ListVIPs lists all configured Virtual IPs
func ListVIPs(w http.ResponseWriter, r *http.Request) {
	if keepalivedManager == nil || !keepalivedManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Keepalived service not available",
			nil,
		))
		return
	}

	vips, err := keepalivedManager.ListVIPs()
	if err != nil {
		logger.Error("Failed to list VIPs", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list VIPs", err))
		return
	}

	utils.RespondSuccess(w, vips)
}

// GetVIPStatus gets the status of a specific VIP
func GetVIPStatus(w http.ResponseWriter, r *http.Request) {
	if keepalivedManager == nil || !keepalivedManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Keepalived service not available",
			nil,
		))
		return
	}

	vipID := chi.URLParam(r, "id")
	if vipID == "" {
		utils.RespondError(w, errors.BadRequest("VIP ID is required", nil))
		return
	}

	status, err := keepalivedManager.GetVIPStatus(vipID)
	if err != nil {
		logger.Error("Failed to get VIP status", zap.Error(err), zap.String("id", vipID))
		utils.RespondError(w, errors.InternalServerError("Failed to get VIP status", err))
		return
	}

	utils.RespondSuccess(w, status)
}

// CreateVIP creates a new Virtual IP
func CreateVIP(w http.ResponseWriter, r *http.Request) {
	if keepalivedManager == nil || !keepalivedManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Keepalived service not available",
			nil,
		))
		return
	}

	var config ha.VIPConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate required fields
	if config.VirtualIP == "" || config.Interface == "" {
		utils.RespondError(w, errors.BadRequest("virtual_ip and interface are required", nil))
		return
	}

	// Set default ID if not provided
	if config.ID == "" {
		config.ID = "VIP_1"
	}

	if err := keepalivedManager.CreateVIP(config); err != nil {
		logger.Error("Failed to create VIP", zap.Error(err), zap.String("vip", config.VirtualIP))
		utils.RespondError(w, errors.InternalServerError("Failed to create VIP", err))
		return
	}

	logger.Info("VIP created", zap.String("vip", config.VirtualIP), zap.String("interface", config.Interface))
	utils.RespondSuccess(w, map[string]string{
		"message":    "VIP created successfully",
		"id":         config.ID,
		"virtual_ip": config.VirtualIP,
	})
}

// DeleteVIP deletes a Virtual IP
func DeleteVIP(w http.ResponseWriter, r *http.Request) {
	if keepalivedManager == nil || !keepalivedManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Keepalived service not available",
			nil,
		))
		return
	}

	vipID := chi.URLParam(r, "id")
	if vipID == "" {
		utils.RespondError(w, errors.BadRequest("VIP ID is required", nil))
		return
	}

	if err := keepalivedManager.DeleteVIP(vipID); err != nil {
		logger.Error("Failed to delete VIP", zap.Error(err), zap.String("id", vipID))
		utils.RespondError(w, errors.InternalServerError("Failed to delete VIP", err))
		return
	}

	logger.Info("VIP deleted", zap.String("id", vipID))
	utils.RespondSuccess(w, map[string]string{
		"message": "VIP deleted successfully",
		"id":      vipID,
	})
}

// PromoteVIPToMaster promotes this node to MASTER for the VIP
func PromoteVIPToMaster(w http.ResponseWriter, r *http.Request) {
	if keepalivedManager == nil || !keepalivedManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Keepalived service not available",
			nil,
		))
		return
	}

	vipID := chi.URLParam(r, "id")
	if vipID == "" {
		utils.RespondError(w, errors.BadRequest("VIP ID is required", nil))
		return
	}

	if err := keepalivedManager.PromoteToMaster(vipID); err != nil {
		logger.Error("Failed to promote VIP to MASTER", zap.Error(err), zap.String("id", vipID))
		utils.RespondError(w, errors.InternalServerError("Failed to promote VIP to MASTER", err))
		return
	}

	logger.Info("VIP promoted to MASTER", zap.String("id", vipID))
	utils.RespondSuccess(w, map[string]string{
		"message": "VIP promoted to MASTER successfully",
		"id":      vipID,
	})
}

// DemoteVIPToBackup demotes this node to BACKUP for the VIP
func DemoteVIPToBackup(w http.ResponseWriter, r *http.Request) {
	if keepalivedManager == nil || !keepalivedManager.IsEnabled() {
		utils.RespondError(w, errors.NewAppError(
			http.StatusServiceUnavailable,
			"Keepalived service not available",
			nil,
		))
		return
	}

	vipID := chi.URLParam(r, "id")
	if vipID == "" {
		utils.RespondError(w, errors.BadRequest("VIP ID is required", nil))
		return
	}

	if err := keepalivedManager.DemoteToBackup(vipID); err != nil {
		logger.Error("Failed to demote VIP to BACKUP", zap.Error(err), zap.String("id", vipID))
		utils.RespondError(w, errors.InternalServerError("Failed to demote VIP to BACKUP", err))
		return
	}

	logger.Info("VIP demoted to BACKUP", zap.String("id", vipID))
	utils.RespondSuccess(w, map[string]string{
		"message": "VIP demoted to BACKUP successfully",
		"id":      vipID,
	})
}

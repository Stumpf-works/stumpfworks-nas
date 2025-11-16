// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/Stumpf-works/stumpfworks-nas/internal/api/middleware"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/storage"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

// ===== Disk Handlers =====

// ListDisks lists all available disks
func ListDisks(w http.ResponseWriter, r *http.Request) {
	disks, err := storage.ListDisks()
	if err != nil {
		logger.Error("Failed to list disks", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list disks", err))
		return
	}

	utils.RespondSuccess(w, disks)
}

// GetDisk retrieves information about a specific disk
func GetDisk(w http.ResponseWriter, r *http.Request) {
	diskName := chi.URLParam(r, "name")

	disk, err := storage.GetDiskInfo(diskName)
	if err != nil {
		logger.Error("Failed to get disk info", zap.String("disk", diskName), zap.Error(err))
		utils.RespondError(w, errors.NotFound("Disk not found", err))
		return
	}

	utils.RespondSuccess(w, disk)
}

// FormatDisk formats a disk with the specified filesystem
func FormatDisk(w http.ResponseWriter, r *http.Request) {
	var req storage.FormatDiskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	if err := storage.FormatDisk(&req); err != nil {
		logger.Error("Failed to format disk", zap.String("disk", req.Disk), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to format disk", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Disk formatted successfully",
	})
}

// GetDiskSMART retrieves SMART data for a disk
func GetDiskSMART(w http.ResponseWriter, r *http.Request) {
	diskName := chi.URLParam(r, "name")

	smart, err := storage.GetSMARTData(diskName)
	if err != nil {
		logger.Error("Failed to get SMART data", zap.String("disk", diskName), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get SMART data", err))
		return
	}

	utils.RespondSuccess(w, smart)
}

// GetDiskHealth retrieves health assessment for a disk
func GetDiskHealth(w http.ResponseWriter, r *http.Request) {
	diskName := chi.URLParam(r, "name")

	health, err := storage.AssessDiskHealth(diskName)
	if err != nil {
		logger.Error("Failed to assess disk health", zap.String("disk", diskName), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to assess disk health", err))
		return
	}

	utils.RespondSuccess(w, health)
}

// ===== Volume Handlers =====

// ListVolumes lists all storage volumes
func ListVolumes(w http.ResponseWriter, r *http.Request) {
	volumes, err := storage.ListVolumes()
	if err != nil {
		logger.Error("Failed to list volumes", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list volumes", err))
		return
	}

	utils.RespondSuccess(w, volumes)
}

// GetVolume retrieves information about a specific volume
func GetVolume(w http.ResponseWriter, r *http.Request) {
	volumeID := chi.URLParam(r, "id")

	volume, err := storage.GetVolume(volumeID)
	if err != nil {
		logger.Error("Failed to get volume", zap.String("id", volumeID), zap.Error(err))
		utils.RespondError(w, errors.NotFound("Volume not found", err))
		return
	}

	utils.RespondSuccess(w, volume)
}

// CreateVolume creates a new storage volume
func CreateVolume(w http.ResponseWriter, r *http.Request) {
	var req storage.CreateVolumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	volume, err := storage.CreateVolume(&req)
	if err != nil {
		logger.Error("Failed to create volume", zap.String("name", req.Name), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create volume", err))
		return
	}

	utils.RespondSuccess(w, volume)
}

// DeleteVolume deletes a storage volume
func DeleteVolume(w http.ResponseWriter, r *http.Request) {
	volumeID := chi.URLParam(r, "id")

	if err := storage.DeleteVolume(volumeID); err != nil {
		logger.Error("Failed to delete volume", zap.String("id", volumeID), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to delete volume", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Volume deleted successfully",
	})
}

// ===== Share Handlers =====

// ListShares lists all network shares (filtered by user permissions)
func ListShares(w http.ResponseWriter, r *http.Request) {
	allShares, err := storage.ListShares()
	if err != nil {
		logger.Error("Failed to list shares", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list shares", err))
		return
	}

	// Get user from context
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		utils.RespondError(w, errors.Unauthorized("User not authenticated", nil))
		return
	}

	// Admins see all shares
	if user.IsAdmin() {
		utils.RespondSuccess(w, allShares)
		return
	}

	// Regular users only see shares they have access to
	filteredShares := filterSharesForUser(allShares, user)
	utils.RespondSuccess(w, filteredShares)
}

// filterSharesForUser filters shares based on user access permissions
func filterSharesForUser(shares []storage.Share, user *models.User) []storage.Share {
	var filtered []storage.Share

	for _, share := range shares {
		// Skip disabled shares
		if !share.Enabled {
			continue
		}

		// Include shares that are open to guests
		if share.GuestOK {
			filtered = append(filtered, share)
			continue
		}

		// Include shares where user is in ValidUsers list
		if len(share.ValidUsers) > 0 {
			for _, validUser := range share.ValidUsers {
				if validUser == user.Username {
					filtered = append(filtered, share)
					break
				}
			}
		}
	}

	return filtered
}

// GetShare retrieves information about a specific share
func GetShare(w http.ResponseWriter, r *http.Request) {
	shareID := chi.URLParam(r, "id")

	share, err := storage.GetShare(shareID)
	if err != nil {
		logger.Error("Failed to get share", zap.String("id", shareID), zap.Error(err))
		utils.RespondError(w, errors.NotFound("Share not found", err))
		return
	}

	utils.RespondSuccess(w, share)
}

// CreateShare creates a new network share
func CreateShare(w http.ResponseWriter, r *http.Request) {
	var req storage.CreateShareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	share, err := storage.CreateShare(&req)
	if err != nil {
		logger.Error("Failed to create share", zap.String("name", req.Name), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create share", err))
		return
	}

	utils.RespondSuccess(w, share)
}

// UpdateShare updates an existing share
func UpdateShare(w http.ResponseWriter, r *http.Request) {
	shareID := chi.URLParam(r, "id")

	var req storage.CreateShareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	share, err := storage.UpdateShare(shareID, &req)
	if err != nil {
		logger.Error("Failed to update share", zap.String("id", shareID), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to update share", err))
		return
	}

	utils.RespondSuccess(w, share)
}

// DeleteShare deletes a network share
func DeleteShare(w http.ResponseWriter, r *http.Request) {
	shareID := chi.URLParam(r, "id")

	if err := storage.DeleteShare(shareID); err != nil {
		logger.Error("Failed to delete share", zap.String("id", shareID), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to delete share", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Share deleted successfully",
	})
}

// EnableShare enables a network share
func EnableShare(w http.ResponseWriter, r *http.Request) {
	shareID := chi.URLParam(r, "id")

	if err := storage.EnableShare(shareID); err != nil {
		logger.Error("Failed to enable share", zap.String("id", shareID), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to enable share", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Share enabled successfully",
	})
}

// DisableShare disables a network share
func DisableShare(w http.ResponseWriter, r *http.Request) {
	shareID := chi.URLParam(r, "id")

	if err := storage.DisableShare(shareID); err != nil {
		logger.Error("Failed to disable share", zap.String("id", shareID), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to disable share", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Share disabled successfully",
	})
}

// ===== Storage Statistics Handlers =====

// GetStorageStats retrieves overall storage statistics
func GetStorageStats(w http.ResponseWriter, r *http.Request) {
	stats, err := storage.GetStorageStats()
	if err != nil {
		logger.Error("Failed to get storage stats", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get storage stats", err))
		return
	}

	utils.RespondSuccess(w, stats)
}

// GetDiskIOStats retrieves I/O statistics for all disks
func GetDiskIOStats(w http.ResponseWriter, r *http.Request) {
	stats, err := storage.GetDiskIOStats()
	if err != nil {
		logger.Error("Failed to get disk I/O stats", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get disk I/O stats", err))
		return
	}

	utils.RespondSuccess(w, stats)
}

// GetDiskIOStatsForDisk retrieves I/O statistics for a specific disk
func GetDiskIOStatsForDisk(w http.ResponseWriter, r *http.Request) {
	diskName := chi.URLParam(r, "name")

	stats, err := storage.GetDiskIOStatsForDisk(diskName)
	if err != nil {
		logger.Error("Failed to get disk I/O stats", zap.String("disk", diskName), zap.Error(err))
		utils.RespondError(w, errors.NotFound("Disk not found", err))
		return
	}

	utils.RespondSuccess(w, stats)
}

// GetAllDiskHealth retrieves health assessment for all disks
func GetAllDiskHealth(w http.ResponseWriter, r *http.Request) {
	healthList, err := storage.GetAllDiskHealth()
	if err != nil {
		logger.Error("Failed to get disk health", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get disk health", err))
		return
	}

	utils.RespondSuccess(w, healthList)
}

// ===== WebSocket for Real-time Monitoring =====

var ioMonitor *storage.DiskIOMonitor

// StartIOMonitoring starts real-time I/O monitoring
func StartIOMonitoring() {
	if ioMonitor != nil {
		return
	}

	ioMonitor = storage.NewDiskIOMonitor(time.Second * 2)
	ioMonitor.Start()

	logger.Info("Storage I/O monitoring started")
}

// StopIOMonitoring stops real-time I/O monitoring
func StopIOMonitoring() {
	if ioMonitor != nil {
		ioMonitor.Stop()
		ioMonitor = nil
		logger.Info("Storage I/O monitoring stopped")
	}
}

// GetIOMonitorStats retrieves the latest I/O monitoring stats
func GetIOMonitorStats(w http.ResponseWriter, r *http.Request) {
	if ioMonitor == nil {
		utils.RespondError(w, errors.NewAppError(http.StatusServiceUnavailable, "I/O monitoring not started", nil))
		return
	}

	// Get stats with timeout
	select {
	case stats := <-ioMonitor.Stats():
		utils.RespondSuccess(w, stats)
	case <-time.After(5 * time.Second):
		utils.RespondError(w, errors.NewAppError(http.StatusRequestTimeout, "Timeout waiting for stats", nil))
	}
}

// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package handlers

import (
	"net/http"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/cache"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
)

// System metrics cache with 5s TTL (frequently polled, needs to be fresh)
var systemMetricsCache = cache.New(5 * time.Second)

// GetSystemInfo returns basic system information
func GetSystemInfo(w http.ResponseWriter, r *http.Request) {
	info, err := system.GetSystemInfo()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to get system information", err))
		return
	}

	utils.RespondSuccess(w, info)
}

// GetSystemMetrics returns real-time system metrics
func GetSystemMetrics(w http.ResponseWriter, r *http.Request) {
	// Try cache first (5s TTL to keep metrics relatively fresh)
	if cached, ok := systemMetricsCache.Get("metrics"); ok {
		utils.RespondSuccess(w, cached)
		return
	}

	// Cache miss - fetch realtime metrics
	metrics, err := system.GetRealtimeSystemMetrics()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to get system metrics", err))
		return
	}

	// Cache for 5 seconds
	systemMetricsCache.Set("metrics", metrics)

	utils.RespondSuccess(w, metrics)
}

// CheckForUpdates checks if system updates are available
func CheckForUpdates(w http.ResponseWriter, r *http.Request) {
	updateInfo, err := system.CheckForUpdates()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to check for updates", err))
		return
	}

	utils.RespondSuccess(w, updateInfo)
}

// ApplyUpdates applies available system updates (admin only)
func ApplyUpdates(w http.ResponseWriter, r *http.Request) {
	err := system.PerformUpdate()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to apply updates", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "System updated successfully. Please restart the server to apply changes.",
	})
}

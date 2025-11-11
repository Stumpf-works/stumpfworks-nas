package handlers

import (
	"net/http"

	"github.com/stumpfworks/nas/internal/system"
	"github.com/stumpfworks/nas/pkg/errors"
	"github.com/stumpfworks/nas/pkg/utils"
)

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
	metrics, err := system.GetSystemMetrics()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to get system metrics", err))
		return
	}

	utils.RespondSuccess(w, metrics)
}

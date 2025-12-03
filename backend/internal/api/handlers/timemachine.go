package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/timemachine"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
)

var timeMachineManager *timemachine.Manager

// InitTimeMachineHandler initializes the Time Machine handler
func InitTimeMachineHandler(manager *timemachine.Manager) {
	timeMachineManager = manager
	logger.Info("Time Machine handler initialized")
}

// GetTimeMachineConfig returns the Time Machine configuration
func GetTimeMachineConfig(w http.ResponseWriter, r *http.Request) {
	if timeMachineManager == nil {
		utils.RespondError(w, errors.InternalServerError("Time Machine not initialized", nil))
		return
	}

	config, err := timeMachineManager.GetConfig()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to get configuration", err))
		return
	}

	utils.RespondSuccess(w, config)
}

// UpdateTimeMachineConfig updates the Time Machine configuration
func UpdateTimeMachineConfig(w http.ResponseWriter, r *http.Request) {
	if timeMachineManager == nil {
		utils.RespondError(w, errors.InternalServerError("Time Machine not initialized", nil))
		return
	}

	var config models.TimeMachineConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := timeMachineManager.UpdateConfig(&config); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to update configuration", err))
		return
	}

	utils.RespondSuccess(w, config)
}

// EnableTimeMachine enables the Time Machine service
func EnableTimeMachine(w http.ResponseWriter, r *http.Request) {
	if timeMachineManager == nil {
		utils.RespondError(w, errors.InternalServerError("Time Machine not initialized", nil))
		return
	}

	if err := timeMachineManager.Enable(); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to enable Time Machine", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"success": true,
		"message": "Time Machine enabled successfully",
	})
}

// DisableTimeMachine disables the Time Machine service
func DisableTimeMachine(w http.ResponseWriter, r *http.Request) {
	if timeMachineManager == nil {
		utils.RespondError(w, errors.InternalServerError("Time Machine not initialized", nil))
		return
	}

	if err := timeMachineManager.Disable(); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to disable Time Machine", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"success": true,
		"message": "Time Machine disabled successfully",
	})
}

// ListTimeMachineDevices returns all Time Machine devices
func ListTimeMachineDevices(w http.ResponseWriter, r *http.Request) {
	if timeMachineManager == nil {
		utils.RespondError(w, errors.InternalServerError("Time Machine not initialized", nil))
		return
	}

	devices, err := timeMachineManager.ListDevices()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to list devices", err))
		return
	}

	utils.RespondSuccess(w, devices)
}

// GetTimeMachineDevice returns a specific Time Machine device
func GetTimeMachineDevice(w http.ResponseWriter, r *http.Request) {
	if timeMachineManager == nil {
		utils.RespondError(w, errors.InternalServerError("Time Machine not initialized", nil))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid device ID", err))
		return
	}

	device, err := timeMachineManager.GetDevice(uint(id))
	if err != nil {
		utils.RespondError(w, errors.NotFound("Device not found", err))
		return
	}

	utils.RespondSuccess(w, device)
}

// CreateTimeMachineDevice creates a new Time Machine device
func CreateTimeMachineDevice(w http.ResponseWriter, r *http.Request) {
	if timeMachineManager == nil {
		utils.RespondError(w, errors.InternalServerError("Time Machine not initialized", nil))
		return
	}

	var device models.TimeMachineDevice
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate required fields
	if device.DeviceName == "" {
		utils.RespondError(w, errors.BadRequest("Device name is required", nil))
		return
	}

	if err := timeMachineManager.CreateDevice(&device); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to create device", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	utils.RespondSuccess(w, device)
}

// UpdateTimeMachineDevice updates a Time Machine device
func UpdateTimeMachineDevice(w http.ResponseWriter, r *http.Request) {
	if timeMachineManager == nil {
		utils.RespondError(w, errors.InternalServerError("Time Machine not initialized", nil))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid device ID", err))
		return
	}

	var device models.TimeMachineDevice
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	device.ID = uint(id)

	if err := timeMachineManager.UpdateDevice(&device); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to update device", err))
		return
	}

	utils.RespondSuccess(w, device)
}

// DeleteTimeMachineDevice removes a Time Machine device
func DeleteTimeMachineDevice(w http.ResponseWriter, r *http.Request) {
	if timeMachineManager == nil {
		utils.RespondError(w, errors.InternalServerError("Time Machine not initialized", nil))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid device ID", err))
		return
	}

	if err := timeMachineManager.DeleteDevice(uint(id)); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to delete device", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"success": true,
		"message": "Device deleted successfully",
	})
}

// UpdateTimeMachineDeviceUsage updates the used space for a device
func UpdateTimeMachineDeviceUsage(w http.ResponseWriter, r *http.Request) {
	if timeMachineManager == nil {
		utils.RespondError(w, errors.InternalServerError("Time Machine not initialized", nil))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid device ID", err))
		return
	}

	if err := timeMachineManager.UpdateDeviceUsage(uint(id)); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to update device usage", err))
		return
	}

	device, err := timeMachineManager.GetDevice(uint(id))
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to get updated device", err))
		return
	}

	utils.RespondSuccess(w, device)
}

// UpdateAllTimeMachineDeviceUsages updates usage for all devices
func UpdateAllTimeMachineDeviceUsages(w http.ResponseWriter, r *http.Request) {
	if timeMachineManager == nil {
		utils.RespondError(w, errors.InternalServerError("Time Machine not initialized", nil))
		return
	}

	if err := timeMachineManager.UpdateAllDeviceUsages(); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to update device usages", err))
		return
	}

	devices, err := timeMachineManager.ListDevices()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to list devices", err))
		return
	}

	utils.RespondSuccess(w, devices)
}

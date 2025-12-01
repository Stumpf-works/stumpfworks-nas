package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/vm"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var vmManager *vm.LibvirtManager

// InitVMManager initializes the VM manager
func InitVMManager(manager *vm.LibvirtManager) {
	vmManager = manager
	logger.Info("VM manager initialized in handlers")
}

// ListVMs lists all virtual machines
func ListVMs(w http.ResponseWriter, r *http.Request) {
	if vmManager == nil {
		utils.RespondError(w, errors.InternalServerError("VM manager not initialized", nil))
		return
	}

	vms, err := vmManager.ListVMs()
	if err != nil {
		logger.Error("Failed to list VMs", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list VMs", err))
		return
	}

	utils.RespondSuccess(w, vms)
}

// GetVM gets details of a specific VM
func GetVM(w http.ResponseWriter, r *http.Request) {
	if vmManager == nil {
		utils.RespondError(w, errors.InternalServerError("VM manager not initialized", nil))
		return
	}

	vmID := chi.URLParam(r, "id")
	if vmID == "" {
		utils.RespondError(w, errors.BadRequest("VM ID is required", nil))
		return
	}

	vmDetails, err := vmManager.GetVM(vmID)
	if err != nil {
		logger.Error("Failed to get VM", zap.Error(err), zap.String("vm_id", vmID))
		utils.RespondError(w, errors.NotFound("VM not found", err))
		return
	}

	utils.RespondSuccess(w, vmDetails)
}

// CreateVM creates a new virtual machine
func CreateVM(w http.ResponseWriter, r *http.Request) {
	if vmManager == nil {
		utils.RespondError(w, errors.InternalServerError("VM manager not initialized", nil))
		return
	}

	var req vm.VMCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	logger.Info("Creating VM via API", zap.String("vm_name", req.Name))

	if err := vmManager.CreateVM(req); err != nil {
		logger.Error("Failed to create VM", zap.Error(err), zap.String("vm_name", req.Name))
		utils.RespondError(w, errors.InternalServerError("Failed to create VM", err))
		return
	}

	logger.Info("VM created successfully via API", zap.String("vm_name", req.Name))
	utils.RespondSuccess(w, map[string]string{
		"message": "VM created successfully",
		"name":    req.Name,
	})
}

// StartVM starts a virtual machine
func StartVM(w http.ResponseWriter, r *http.Request) {
	if vmManager == nil {
		utils.RespondError(w, errors.InternalServerError("VM manager not initialized", nil))
		return
	}

	vmID := chi.URLParam(r, "id")
	if vmID == "" {
		utils.RespondError(w, errors.BadRequest("VM ID is required", nil))
		return
	}

	logger.Info("Starting VM via API", zap.String("vm_id", vmID))

	if err := vmManager.StartVM(vmID); err != nil {
		logger.Error("Failed to start VM", zap.Error(err), zap.String("vm_id", vmID))
		utils.RespondError(w, errors.InternalServerError("Failed to start VM", err))
		return
	}

	logger.Info("VM started successfully via API", zap.String("vm_id", vmID))
	utils.RespondSuccess(w, map[string]string{
		"message": "VM started successfully",
		"vm_id":   vmID,
	})
}

// StopVM stops a virtual machine
func StopVM(w http.ResponseWriter, r *http.Request) {
	if vmManager == nil {
		utils.RespondError(w, errors.InternalServerError("VM manager not initialized", nil))
		return
	}

	vmID := chi.URLParam(r, "id")
	if vmID == "" {
		utils.RespondError(w, errors.BadRequest("VM ID is required", nil))
		return
	}

	// Check if force shutdown is requested
	force := r.URL.Query().Get("force") == "true"

	logger.Info("Stopping VM via API", zap.String("vm_id", vmID), zap.Bool("force", force))

	if err := vmManager.StopVM(vmID, force); err != nil {
		logger.Error("Failed to stop VM", zap.Error(err), zap.String("vm_id", vmID))
		utils.RespondError(w, errors.InternalServerError("Failed to stop VM", err))
		return
	}

	logger.Info("VM stopped successfully via API", zap.String("vm_id", vmID))
	utils.RespondSuccess(w, map[string]string{
		"message": "VM stopped successfully",
		"vm_id":   vmID,
	})
}

// DeleteVM deletes a virtual machine
func DeleteVM(w http.ResponseWriter, r *http.Request) {
	if vmManager == nil {
		utils.RespondError(w, errors.InternalServerError("VM manager not initialized", nil))
		return
	}

	vmID := chi.URLParam(r, "id")
	if vmID == "" {
		utils.RespondError(w, errors.BadRequest("VM ID is required", nil))
		return
	}

	// Check if disk deletion is requested
	deleteDisks := r.URL.Query().Get("delete_disks") == "true"

	logger.Info("Deleting VM via API", zap.String("vm_id", vmID), zap.Bool("delete_disks", deleteDisks))

	if err := vmManager.DeleteVM(vmID, deleteDisks); err != nil {
		logger.Error("Failed to delete VM", zap.Error(err), zap.String("vm_id", vmID))
		utils.RespondError(w, errors.InternalServerError("Failed to delete VM", err))
		return
	}

	logger.Info("VM deleted successfully via API", zap.String("vm_id", vmID))
	utils.RespondSuccess(w, map[string]string{
		"message": "VM deleted successfully",
		"vm_id":   vmID,
	})
}

// GetVMVNCPort gets the VNC port for a VM
func GetVMVNCPort(w http.ResponseWriter, r *http.Request) {
	if vmManager == nil {
		utils.RespondError(w, errors.InternalServerError("VM manager not initialized", nil))
		return
	}

	vmID := chi.URLParam(r, "id")
	if vmID == "" {
		utils.RespondError(w, errors.BadRequest("VM ID is required", nil))
		return
	}

	port, err := vmManager.GetVNCPort(vmID)
	if err != nil {
		logger.Error("Failed to get VNC port", zap.Error(err), zap.String("vm_id", vmID))
		utils.RespondError(w, errors.InternalServerError("Failed to get VNC port", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"vm_id": vmID,
		"port":  port,
	})
}

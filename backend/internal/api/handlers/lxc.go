package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/lxc"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var lxcManager *lxc.LXCManager

// InitLXCManager initializes the LXC manager
func InitLXCManager(manager *lxc.LXCManager) {
	lxcManager = manager
	logger.Info("LXC manager initialized in handlers")
}

// ListContainers lists all LXC containers
func ListContainers(w http.ResponseWriter, r *http.Request) {
	if lxcManager == nil {
		utils.RespondError(w, errors.InternalServerError("LXC manager not initialized", nil))
		return
	}

	containers, err := lxcManager.ListContainers()
	if err != nil {
		logger.Error("Failed to list containers", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list containers", err))
		return
	}

	utils.RespondSuccess(w, containers)
}

// GetContainer gets details of a specific container
func GetContainer(w http.ResponseWriter, r *http.Request) {
	if lxcManager == nil {
		utils.RespondError(w, errors.InternalServerError("LXC manager not initialized", nil))
		return
	}

	containerName := chi.URLParam(r, "name")
	if containerName == "" {
		utils.RespondError(w, errors.BadRequest("Container name is required", nil))
		return
	}

	container, err := lxcManager.GetContainer(containerName)
	if err != nil {
		logger.Error("Failed to get container", zap.Error(err), zap.String("container", containerName))
		utils.RespondError(w, errors.NotFound("Container not found", err))
		return
	}

	utils.RespondSuccess(w, container)
}

// CreateContainer creates a new LXC container
func CreateContainer(w http.ResponseWriter, r *http.Request) {
	if lxcManager == nil {
		utils.RespondError(w, errors.InternalServerError("LXC manager not initialized", nil))
		return
	}

	var req lxc.ContainerCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	logger.Info("Creating container via API", zap.String("container_name", req.Name))

	if err := lxcManager.CreateContainer(req); err != nil {
		logger.Error("Failed to create container", zap.Error(err), zap.String("container_name", req.Name))
		utils.RespondError(w, errors.InternalServerError("Failed to create container", err))
		return
	}

	logger.Info("Container created successfully via API", zap.String("container_name", req.Name))
	utils.RespondSuccess(w, map[string]string{
		"message": "Container created successfully",
		"name":    req.Name,
	})
}

// StartContainer starts an LXC container
func StartContainer(w http.ResponseWriter, r *http.Request) {
	if lxcManager == nil {
		utils.RespondError(w, errors.InternalServerError("LXC manager not initialized", nil))
		return
	}

	containerName := chi.URLParam(r, "name")
	if containerName == "" {
		utils.RespondError(w, errors.BadRequest("Container name is required", nil))
		return
	}

	logger.Info("Starting container via API", zap.String("container", containerName))

	if err := lxcManager.StartContainer(containerName); err != nil {
		logger.Error("Failed to start container", zap.Error(err), zap.String("container", containerName))
		utils.RespondError(w, errors.InternalServerError("Failed to start container", err))
		return
	}

	logger.Info("Container started successfully via API", zap.String("container", containerName))
	utils.RespondSuccess(w, map[string]string{
		"message": "Container started successfully",
		"name":    containerName,
	})
}

// StopContainer stops an LXC container
func StopContainer(w http.ResponseWriter, r *http.Request) {
	if lxcManager == nil {
		utils.RespondError(w, errors.InternalServerError("LXC manager not initialized", nil))
		return
	}

	containerName := chi.URLParam(r, "name")
	if containerName == "" {
		utils.RespondError(w, errors.BadRequest("Container name is required", nil))
		return
	}

	// Check if force shutdown is requested
	force := r.URL.Query().Get("force") == "true"

	logger.Info("Stopping container via API", zap.String("container", containerName), zap.Bool("force", force))

	if err := lxcManager.StopContainer(containerName, force); err != nil {
		logger.Error("Failed to stop container", zap.Error(err), zap.String("container", containerName))
		utils.RespondError(w, errors.InternalServerError("Failed to stop container", err))
		return
	}

	logger.Info("Container stopped successfully via API", zap.String("container", containerName))
	utils.RespondSuccess(w, map[string]string{
		"message": "Container stopped successfully",
		"name":    containerName,
	})
}

// DeleteContainer deletes an LXC container
func DeleteContainer(w http.ResponseWriter, r *http.Request) {
	if lxcManager == nil {
		utils.RespondError(w, errors.InternalServerError("LXC manager not initialized", nil))
		return
	}

	containerName := chi.URLParam(r, "name")
	if containerName == "" {
		utils.RespondError(w, errors.BadRequest("Container name is required", nil))
		return
	}

	logger.Info("Deleting container via API", zap.String("container", containerName))

	if err := lxcManager.DeleteContainer(containerName); err != nil {
		logger.Error("Failed to delete container", zap.Error(err), zap.String("container", containerName))
		utils.RespondError(w, errors.InternalServerError("Failed to delete container", err))
		return
	}

	logger.Info("Container deleted successfully via API", zap.String("container", containerName))
	utils.RespondSuccess(w, map[string]string{
		"message": "Container deleted successfully",
		"name":    containerName,
	})
}

// ListLXCTemplates lists available LXC templates
func ListLXCTemplates(w http.ResponseWriter, r *http.Request) {
	if lxcManager == nil {
		utils.RespondError(w, errors.InternalServerError("LXC manager not initialized", nil))
		return
	}

	templates, err := lxcManager.GetAvailableTemplates()
	if err != nil {
		logger.Error("Failed to list templates", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list templates", err))
		return
	}

	utils.RespondSuccess(w, templates)
}

// ExecContainerCommand executes a command in a container
func ExecContainerCommand(w http.ResponseWriter, r *http.Request) {
	if lxcManager == nil {
		utils.RespondError(w, errors.InternalServerError("LXC manager not initialized", nil))
		return
	}

	containerName := chi.URLParam(r, "name")
	if containerName == "" {
		utils.RespondError(w, errors.BadRequest("Container name is required", nil))
		return
	}

	var req struct {
		Command string `json:"command"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Command == "" {
		utils.RespondError(w, errors.BadRequest("Command is required", nil))
		return
	}

	logger.Info("Executing command in container", zap.String("container", containerName), zap.String("command", req.Command))

	result, err := lxcManager.ExecCommand(containerName, req.Command)
	if err != nil {
		logger.Error("Failed to execute command", zap.Error(err), zap.String("container", containerName))
		utils.RespondError(w, errors.InternalServerError("Failed to execute command", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"stdout":    result.Stdout,
		"stderr":    result.Stderr,
		"exit_code": result.ExitCode,
	})
}

// GetContainerConsole gets console access information for a container
func GetContainerConsole(w http.ResponseWriter, r *http.Request) {
	if lxcManager == nil {
		utils.RespondError(w, errors.InternalServerError("LXC manager not initialized", nil))
		return
	}

	containerName := chi.URLParam(r, "name")
	if containerName == "" {
		utils.RespondError(w, errors.BadRequest("Container name is required", nil))
		return
	}

	consoleCmd, err := lxcManager.GetConsoleURL(containerName)
	if err != nil {
		logger.Error("Failed to get console access", zap.Error(err), zap.String("container", containerName))
		utils.RespondError(w, errors.InternalServerError("Failed to get console access", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"console_command": consoleCmd,
		"container_name":  containerName,
	})
}

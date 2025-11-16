// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/docker"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/docker/docker/api/types/container"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// DockerHandler handles Docker-related HTTP requests
type DockerHandler struct {
	service *docker.Service
}

// NewDockerHandler creates a new Docker handler
func NewDockerHandler() *DockerHandler {
	return &DockerHandler{
		service: docker.GetService(),
	}
}

// CheckAvailability middleware to check if Docker is available
func (h *DockerHandler) CheckAvailability(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.service == nil || !h.service.IsAvailable() {
			utils.RespondError(w, errors.NewAppError(503, "Docker is not available on this system", nil))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Container Handlers

// ListContainers lists all containers
func (h *DockerHandler) ListContainers(w http.ResponseWriter, r *http.Request) {
	all := r.URL.Query().Get("all") == "true"

	containers, err := h.service.ListContainers(r.Context(), all)
	if err != nil {
		logger.Error("Failed to list containers", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list containers", err))
		return
	}

	utils.RespondSuccess(w, containers)
}

// InspectContainer gets detailed container information
func (h *DockerHandler) InspectContainer(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "id")

	containerInfo, err := h.service.InspectContainer(r.Context(), containerID)
	if err != nil {
		logger.Error("Failed to inspect container", zap.Error(err), zap.String("container", containerID))
		utils.RespondError(w, errors.InternalServerError("Failed to inspect container", err))
		return
	}

	utils.RespondSuccess(w, containerInfo)
}

// GetContainerStats gets container resource usage
func (h *DockerHandler) GetContainerStats(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "id")

	stats, err := h.service.GetContainerStats(r.Context(), containerID)
	if err != nil {
		logger.Error("Failed to get container stats", zap.Error(err), zap.String("container", containerID))
		utils.RespondError(w, errors.InternalServerError("Failed to get container stats", err))
		return
	}

	utils.RespondSuccess(w, stats)
}

// StartContainer starts a container
func (h *DockerHandler) StartContainer(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "id")

	if err := h.service.StartContainer(r.Context(), containerID); err != nil {
		logger.Error("Failed to start container", zap.Error(err), zap.String("container", containerID))
		utils.RespondError(w, errors.InternalServerError("Failed to start container", err))
		return
	}

	logger.Info("Container started", zap.String("container", containerID))
	utils.RespondSuccess(w, map[string]string{"message": "Container started successfully"})
}

// StopContainer stops a container
func (h *DockerHandler) StopContainer(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "id")

	if err := h.service.StopContainer(r.Context(), containerID); err != nil {
		logger.Error("Failed to stop container", zap.Error(err), zap.String("container", containerID))
		utils.RespondError(w, errors.InternalServerError("Failed to stop container", err))
		return
	}

	logger.Info("Container stopped", zap.String("container", containerID))
	utils.RespondSuccess(w, map[string]string{"message": "Container stopped successfully"})
}

// RestartContainer restarts a container
func (h *DockerHandler) RestartContainer(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "id")

	if err := h.service.RestartContainer(r.Context(), containerID); err != nil {
		logger.Error("Failed to restart container", zap.Error(err), zap.String("container", containerID))
		utils.RespondError(w, errors.InternalServerError("Failed to restart container", err))
		return
	}

	logger.Info("Container restarted", zap.String("container", containerID))
	utils.RespondSuccess(w, map[string]string{"message": "Container restarted successfully"})
}

// PauseContainer pauses a container
func (h *DockerHandler) PauseContainer(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "id")

	if err := h.service.PauseContainer(r.Context(), containerID); err != nil {
		logger.Error("Failed to pause container", zap.Error(err), zap.String("container", containerID))
		utils.RespondError(w, errors.InternalServerError("Failed to pause container", err))
		return
	}

	logger.Info("Container paused", zap.String("container", containerID))
	utils.RespondSuccess(w, map[string]string{"message": "Container paused successfully"})
}

// UnpauseContainer unpauses a container
func (h *DockerHandler) UnpauseContainer(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "id")

	if err := h.service.UnpauseContainer(r.Context(), containerID); err != nil {
		logger.Error("Failed to unpause container", zap.Error(err), zap.String("container", containerID))
		utils.RespondError(w, errors.InternalServerError("Failed to unpause container", err))
		return
	}

	logger.Info("Container unpaused", zap.String("container", containerID))
	utils.RespondSuccess(w, map[string]string{"message": "Container unpaused successfully"})
}

// RemoveContainer removes a container
func (h *DockerHandler) RemoveContainer(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "id")

	if err := h.service.RemoveContainer(r.Context(), containerID); err != nil {
		logger.Error("Failed to remove container", zap.Error(err), zap.String("container", containerID))
		utils.RespondError(w, errors.InternalServerError("Failed to remove container", err))
		return
	}

	logger.Info("Container removed", zap.String("container", containerID))
	utils.RespondSuccess(w, map[string]string{"message": "Container removed successfully"})
}

// GetContainerLogs gets container logs
func (h *DockerHandler) GetContainerLogs(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "id")

	logs, err := h.service.GetContainerLogs(r.Context(), containerID)
	if err != nil {
		logger.Error("Failed to get container logs", zap.Error(err), zap.String("container", containerID))
		utils.RespondError(w, errors.InternalServerError("Failed to get container logs", err))
		return
	}

	utils.RespondSuccess(w, logs)
}

// CreateContainer creates a new container
func (h *DockerHandler) CreateContainer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name          string            `json:"name"`
		Image         string            `json:"image"`
		Env           []string          `json:"env"`
		Cmd           []string          `json:"cmd"`
		Ports         map[string]string `json:"ports"`
		Volumes       map[string]string `json:"volumes"`
		RestartPolicy string            `json:"restartPolicy"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	config := &container.Config{
		Image: req.Image,
		Env:   req.Env,
		Cmd:   req.Cmd,
	}

	hostConfig := &container.HostConfig{}
	if req.RestartPolicy != "" {
		hostConfig.RestartPolicy = container.RestartPolicy{Name: container.RestartPolicyMode(req.RestartPolicy)}
	}

	resp, err := h.service.CreateContainer(r.Context(), config, hostConfig, nil, req.Name)
	if err != nil {
		logger.Error("Failed to create container", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create container", err))
		return
	}

	logger.Info("Container created", zap.String("id", resp.ID), zap.String("name", req.Name))
	utils.RespondSuccess(w, resp)
}

// Image Handlers

// ListImages lists all images
func (h *DockerHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	images, err := h.service.ListImages(r.Context())
	if err != nil {
		logger.Error("Failed to list images", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list images", err))
		return
	}

	utils.RespondSuccess(w, images)
}

// InspectImage gets detailed image information
func (h *DockerHandler) InspectImage(w http.ResponseWriter, r *http.Request) {
	imageID := chi.URLParam(r, "id")

	imageInfo, err := h.service.InspectImage(r.Context(), imageID)
	if err != nil {
		logger.Error("Failed to inspect image", zap.Error(err), zap.String("image", imageID))
		utils.RespondError(w, errors.InternalServerError("Failed to inspect image", err))
		return
	}

	utils.RespondSuccess(w, imageInfo)
}

// PullImage pulls an image
func (h *DockerHandler) PullImage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Image string `json:"image"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Image == "" {
		utils.RespondError(w, errors.BadRequest("Image name is required", nil))
		return
	}

	output, err := h.service.PullImage(r.Context(), req.Image)
	if err != nil {
		logger.Error("Failed to pull image", zap.Error(err), zap.String("image", req.Image))
		utils.RespondError(w, errors.InternalServerError("Failed to pull image", err))
		return
	}

	logger.Info("Image pulled", zap.String("image", req.Image))
	utils.RespondSuccess(w, output)
}

// SearchImages searches for images on Docker Hub
func (h *DockerHandler) SearchImages(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	if term == "" {
		utils.RespondError(w, errors.BadRequest("Search term is required", nil))
		return
	}

	results, err := h.service.SearchImages(r.Context(), term)
	if err != nil {
		logger.Error("Failed to search images", zap.Error(err), zap.String("term", term))
		utils.RespondError(w, errors.InternalServerError("Failed to search images", err))
		return
	}

	utils.RespondSuccess(w, results)
}

// RemoveImage removes an image
func (h *DockerHandler) RemoveImage(w http.ResponseWriter, r *http.Request) {
	imageID := chi.URLParam(r, "id")

	if err := h.service.RemoveImage(r.Context(), imageID); err != nil {
		logger.Error("Failed to remove image", zap.Error(err), zap.String("image", imageID))
		utils.RespondError(w, errors.InternalServerError("Failed to remove image", err))
		return
	}

	logger.Info("Image removed", zap.String("image", imageID))
	utils.RespondSuccess(w, map[string]string{"message": "Image removed successfully"})
}

// Volume Handlers

// ListVolumes lists all volumes
func (h *DockerHandler) ListVolumes(w http.ResponseWriter, r *http.Request) {
	volumes, err := h.service.ListVolumes(r.Context())
	if err != nil {
		logger.Error("Failed to list volumes", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list volumes", err))
		return
	}

	utils.RespondSuccess(w, volumes)
}

// InspectVolume gets detailed volume information
func (h *DockerHandler) InspectVolume(w http.ResponseWriter, r *http.Request) {
	volumeName := chi.URLParam(r, "name")

	volumeInfo, err := h.service.InspectVolume(r.Context(), volumeName)
	if err != nil {
		logger.Error("Failed to inspect volume", zap.Error(err), zap.String("volume", volumeName))
		utils.RespondError(w, errors.InternalServerError("Failed to inspect volume", err))
		return
	}

	utils.RespondSuccess(w, volumeInfo)
}

// CreateVolume creates a new volume
func (h *DockerHandler) CreateVolume(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string            `json:"name"`
		Driver string            `json:"driver"`
		Labels map[string]string `json:"labels"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Driver == "" {
		req.Driver = "local"
	}

	volume, err := h.service.CreateVolume(r.Context(), req.Name, req.Driver, req.Labels)
	if err != nil {
		logger.Error("Failed to create volume", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create volume", err))
		return
	}

	logger.Info("Volume created", zap.String("name", volume.Name))
	utils.RespondSuccess(w, volume)
}

// RemoveVolume removes a volume
func (h *DockerHandler) RemoveVolume(w http.ResponseWriter, r *http.Request) {
	volumeName := chi.URLParam(r, "name")

	if err := h.service.RemoveVolume(r.Context(), volumeName); err != nil {
		logger.Error("Failed to remove volume", zap.Error(err), zap.String("volume", volumeName))
		utils.RespondError(w, errors.InternalServerError("Failed to remove volume", err))
		return
	}

	logger.Info("Volume removed", zap.String("volume", volumeName))
	utils.RespondSuccess(w, map[string]string{"message": "Volume removed successfully"})
}

// Network Handlers

// ListNetworks lists all networks
func (h *DockerHandler) ListNetworks(w http.ResponseWriter, r *http.Request) {
	networks, err := h.service.ListNetworks(r.Context())
	if err != nil {
		logger.Error("Failed to list networks", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list networks", err))
		return
	}

	utils.RespondSuccess(w, networks)
}

// InspectNetwork gets detailed network information
func (h *DockerHandler) InspectNetwork(w http.ResponseWriter, r *http.Request) {
	networkID := chi.URLParam(r, "id")

	networkInfo, err := h.service.InspectNetwork(r.Context(), networkID)
	if err != nil {
		logger.Error("Failed to inspect network", zap.Error(err), zap.String("network", networkID))
		utils.RespondError(w, errors.InternalServerError("Failed to inspect network", err))
		return
	}

	utils.RespondSuccess(w, networkInfo)
}

// CreateNetwork creates a new network
func (h *DockerHandler) CreateNetwork(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string `json:"name"`
		Driver string `json:"driver"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Name == "" {
		utils.RespondError(w, errors.BadRequest("Network name is required", nil))
		return
	}

	if req.Driver == "" {
		req.Driver = "bridge"
	}

	resp, err := h.service.CreateNetwork(r.Context(), req.Name, req.Driver)
	if err != nil {
		logger.Error("Failed to create network", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create network", err))
		return
	}

	logger.Info("Network created", zap.String("name", req.Name), zap.String("id", resp.ID))
	utils.RespondSuccess(w, resp)
}

// RemoveNetwork removes a network
func (h *DockerHandler) RemoveNetwork(w http.ResponseWriter, r *http.Request) {
	networkID := chi.URLParam(r, "id")

	if err := h.service.RemoveNetwork(r.Context(), networkID); err != nil {
		logger.Error("Failed to remove network", zap.Error(err), zap.String("network", networkID))
		utils.RespondError(w, errors.InternalServerError("Failed to remove network", err))
		return
	}

	logger.Info("Network removed", zap.String("network", networkID))
	utils.RespondSuccess(w, map[string]string{"message": "Network removed successfully"})
}

// ConnectContainerToNetwork connects a container to a network
func (h *DockerHandler) ConnectContainerToNetwork(w http.ResponseWriter, r *http.Request) {
	networkID := chi.URLParam(r, "id")

	var req struct {
		ContainerID string `json:"container"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.service.ConnectContainerToNetwork(r.Context(), networkID, req.ContainerID); err != nil {
		logger.Error("Failed to connect container to network", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to connect container to network", err))
		return
	}

	logger.Info("Container connected to network", zap.String("network", networkID), zap.String("container", req.ContainerID))
	utils.RespondSuccess(w, map[string]string{"message": "Container connected to network successfully"})
}

// DisconnectContainerFromNetwork disconnects a container from a network
func (h *DockerHandler) DisconnectContainerFromNetwork(w http.ResponseWriter, r *http.Request) {
	networkID := chi.URLParam(r, "id")

	var req struct {
		ContainerID string `json:"container"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.service.DisconnectContainerFromNetwork(r.Context(), networkID, req.ContainerID); err != nil {
		logger.Error("Failed to disconnect container from network", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to disconnect container from network", err))
		return
	}

	logger.Info("Container disconnected from network", zap.String("network", networkID), zap.String("container", req.ContainerID))
	utils.RespondSuccess(w, map[string]string{"message": "Container disconnected from network successfully"})
}

// System Handlers

// GetDockerInfo gets Docker system information
func (h *DockerHandler) GetDockerInfo(w http.ResponseWriter, r *http.Request) {
	info, err := h.service.GetDockerInfo(r.Context())
	if err != nil {
		logger.Error("Failed to get Docker info", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get Docker info", err))
		return
	}

	utils.RespondSuccess(w, info)
}

// GetDockerVersion gets Docker version information
func (h *DockerHandler) GetDockerVersion(w http.ResponseWriter, r *http.Request) {
	version, err := h.service.GetDockerVersion(r.Context())
	if err != nil {
		logger.Error("Failed to get Docker version", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get Docker version", err))
		return
	}

	utils.RespondSuccess(w, version)
}

// PruneSystem prunes unused Docker objects
func (h *DockerHandler) PruneSystem(w http.ResponseWriter, r *http.Request) {
	usage, err := h.service.PruneSystem(r.Context())
	if err != nil {
		logger.Error("Failed to prune system", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to prune system", err))
		return
	}

	logger.Info("System pruned successfully")
	utils.RespondSuccess(w, usage)
}

// UpdateContainerResources updates resource limits for a container
func (h *DockerHandler) UpdateContainerResources(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "id")

	var req struct {
		Memory     int64 `json:"memory"`      // Memory limit in bytes
		MemorySwap int64 `json:"memorySwap"`  // Memory + Swap limit
		CPUShares  int64 `json:"cpuShares"`   // CPU shares (relative weight)
		CPUQuota   int64 `json:"cpuQuota"`    // CPU quota in microseconds
		CPUPeriod  int64 `json:"cpuPeriod"`   // CPU period in microseconds
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	resources := container.Resources{
		Memory:     req.Memory,
		MemorySwap: req.MemorySwap,
		CPUShares:  req.CPUShares,
		CPUQuota:   req.CPUQuota,
		CPUPeriod:  req.CPUPeriod,
	}

	if err := h.service.UpdateContainerResources(r.Context(), containerID, resources); err != nil {
		logger.Error("Failed to update container resources", zap.Error(err), zap.String("containerID", containerID))
		utils.RespondError(w, errors.InternalServerError("Failed to update container resources", err))
		return
	}

	logger.Info("Container resources updated", zap.String("containerID", containerID))
	utils.RespondSuccess(w, map[string]string{"message": "Container resources updated successfully"})
}

// ExecContainer executes a command in a container
func (h *DockerHandler) ExecContainer(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "id")

	var req struct {
		Command []string `json:"command"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if len(req.Command) == 0 {
		utils.RespondError(w, errors.BadRequest("Command is required", nil))
		return
	}

	output, err := h.service.ExecContainer(r.Context(), containerID, req.Command)
	if err != nil {
		logger.Error("Failed to execute command in container", zap.Error(err), zap.String("containerID", containerID))
		utils.RespondError(w, errors.InternalServerError("Failed to execute command", err))
		return
	}

	logger.Info("Command executed in container", zap.String("containerID", containerID))
	utils.RespondSuccess(w, map[string]string{"output": output})
}

// GetContainerTop gets running processes in a container
func (h *DockerHandler) GetContainerTop(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "id")

	top, err := h.service.GetContainerTop(r.Context(), containerID)
	if err != nil {
		logger.Error("Failed to get container processes", zap.Error(err), zap.String("containerID", containerID))
		utils.RespondError(w, errors.InternalServerError("Failed to get container processes", err))
		return
	}

	utils.RespondSuccess(w, top)
}

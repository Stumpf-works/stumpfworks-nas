package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/docker"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
)

// SimpleDockerHandler handles Docker-related requests
type SimpleDockerHandler struct {
	client *docker.SimpleClient
}

// NewSimpleDockerHandler creates a new simple Docker handler
func NewSimpleDockerHandler() (*SimpleDockerHandler, error) {
	client, err := docker.NewSimpleClient()
	if err != nil {
		return nil, err
	}
	return &SimpleDockerHandler{client: client}, nil
}

// ListContainers handles GET /api/docker/containers
func (h *SimpleDockerHandler) ListContainers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	all := r.URL.Query().Get("all") == "true"

	containers, err := h.client.ListContainers(ctx, all)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to list containers", err))
		return
	}

	utils.RespondSuccess(w, containers)
}

// StartContainer handles POST /api/docker/containers/{id}/start
func (h *SimpleDockerHandler) StartContainer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.client.StartContainer(ctx, id); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to start container", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Container started"})
}

// StopContainer handles POST /api/docker/containers/{id}/stop
func (h *SimpleDockerHandler) StopContainer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.client.StopContainer(ctx, id); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to stop container", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Container stopped"})
}

// RestartContainer handles POST /api/docker/containers/{id}/restart
func (h *SimpleDockerHandler) RestartContainer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.client.RestartContainer(ctx, id); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to restart container", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Container restarted"})
}

// RemoveContainer handles DELETE /api/docker/containers/{id}
func (h *SimpleDockerHandler) RemoveContainer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	force := r.URL.Query().Get("force") == "true"

	if err := h.client.RemoveContainer(ctx, id, force); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to remove container", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Container removed"})
}

// GetContainerLogs handles GET /api/docker/containers/{id}/logs
func (h *SimpleDockerHandler) GetContainerLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	logs, err := h.client.GetContainerLogs(ctx, id)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to get container logs", err))
		return
	}
	defer logs.Close()

	w.Header().Set("Content-Type", "text/plain")
	io.Copy(w, logs)
}

// ListImages handles GET /api/docker/images
func (h *SimpleDockerHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	images, err := h.client.ListImages(ctx)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to list images", err))
		return
	}

	utils.RespondSuccess(w, images)
}

// PullImage handles POST /api/docker/images/pull
func (h *SimpleDockerHandler) PullImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Image string `json:"image"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request", err))
		return
	}

	reader, err := h.client.PullImage(ctx, req.Image)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to pull image", err))
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", "text/plain")
	io.Copy(w, reader)
}

// RemoveImage handles DELETE /api/docker/images/{id}
func (h *SimpleDockerHandler) RemoveImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	force := r.URL.Query().Get("force") == "true"

	if err := h.client.RemoveImage(ctx, id, force); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to remove image", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Image removed"})
}

// ListVolumes handles GET /api/docker/volumes
func (h *SimpleDockerHandler) ListVolumes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	volumes, err := h.client.ListVolumes(ctx)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to list volumes", err))
		return
	}

	utils.RespondSuccess(w, volumes)
}

// RemoveVolume handles DELETE /api/docker/volumes/{id}
func (h *SimpleDockerHandler) RemoveVolume(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	force := r.URL.Query().Get("force") == "true"

	if err := h.client.RemoveVolume(ctx, id, force); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to remove volume", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Volume removed"})
}

// ListNetworks handles GET /api/docker/networks
func (h *SimpleDockerHandler) ListNetworks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	networks, err := h.client.ListNetworks(ctx)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to list networks", err))
		return
	}

	utils.RespondSuccess(w, networks)
}

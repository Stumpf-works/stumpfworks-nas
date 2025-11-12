package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Stumpf-works/stumpfworks-nas/internal/backup"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// BackupHandler handles backup-related HTTP requests
type BackupHandler struct {
	service *backup.Service
}

// NewBackupHandler creates a new backup handler
func NewBackupHandler() *BackupHandler {
	return &BackupHandler{
		service: backup.GetService(),
	}
}

// CheckAvailability middleware checks if backup service is available
func (h *BackupHandler) CheckAvailability(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.service == nil {
			utils.RespondError(w, errors.NewAppError(
				http.StatusServiceUnavailable,
				"Backup service is not available",
				nil,
			))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ListJobs lists all backup jobs
func (h *BackupHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.service.ListJobs(r.Context())
	if err != nil {
		logger.Error("Failed to list backup jobs", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list backup jobs", err))
		return
	}

	utils.RespondSuccess(w, jobs)
}

// GetJob gets a specific backup job
func (h *BackupHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")

	job, err := h.service.GetJob(r.Context(), jobID)
	if err != nil {
		logger.Error("Failed to get backup job", zap.Error(err), zap.String("jobID", jobID))
		utils.RespondError(w, errors.NotFound("Backup job not found", err))
		return
	}

	utils.RespondSuccess(w, job)
}

// CreateJob creates a new backup job
func (h *BackupHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var job backup.BackupJob

	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.service.CreateJob(r.Context(), &job); err != nil {
		logger.Error("Failed to create backup job", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create backup job", err))
		return
	}

	logger.Info("Backup job created", zap.String("jobID", job.ID), zap.String("name", job.Name))
	utils.RespondSuccess(w, job)
}

// UpdateJob updates an existing backup job
func (h *BackupHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")

	var updates backup.BackupJob
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.service.UpdateJob(r.Context(), jobID, &updates); err != nil {
		logger.Error("Failed to update backup job", zap.Error(err), zap.String("jobID", jobID))
		utils.RespondError(w, errors.InternalServerError("Failed to update backup job", err))
		return
	}

	logger.Info("Backup job updated", zap.String("jobID", jobID))
	utils.RespondSuccess(w, map[string]string{"message": "Backup job updated successfully"})
}

// DeleteJob deletes a backup job
func (h *BackupHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")

	if err := h.service.DeleteJob(r.Context(), jobID); err != nil {
		logger.Error("Failed to delete backup job", zap.Error(err), zap.String("jobID", jobID))
		utils.RespondError(w, errors.InternalServerError("Failed to delete backup job", err))
		return
	}

	logger.Info("Backup job deleted", zap.String("jobID", jobID))
	utils.RespondSuccess(w, map[string]string{"message": "Backup job deleted successfully"})
}

// RunJob executes a backup job
func (h *BackupHandler) RunJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")

	history, err := h.service.RunJob(r.Context(), jobID)
	if err != nil {
		logger.Error("Failed to run backup job", zap.Error(err), zap.String("jobID", jobID))
		utils.RespondError(w, errors.InternalServerError("Failed to run backup job", err))
		return
	}

	logger.Info("Backup job started", zap.String("jobID", jobID))
	utils.RespondSuccess(w, history)
}

// GetHistory gets backup history
func (h *BackupHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("jobId")
	limitStr := r.URL.Query().Get("limit")

	limit := 50 // default
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	history, err := h.service.GetHistory(r.Context(), jobID, limit)
	if err != nil {
		logger.Error("Failed to get backup history", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get backup history", err))
		return
	}

	utils.RespondSuccess(w, history)
}

// ListSnapshots lists all snapshots
func (h *BackupHandler) ListSnapshots(w http.ResponseWriter, r *http.Request) {
	snapshots, err := h.service.ListSnapshots(r.Context())
	if err != nil {
		logger.Error("Failed to list snapshots", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list snapshots", err))
		return
	}

	utils.RespondSuccess(w, snapshots)
}

// CreateSnapshot creates a new snapshot
func (h *BackupHandler) CreateSnapshot(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Filesystem string `json:"filesystem"`
		Name       string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Filesystem == "" || req.Name == "" {
		utils.RespondError(w, errors.BadRequest("Filesystem and name are required", nil))
		return
	}

	snapshot, err := h.service.CreateSnapshot(r.Context(), req.Filesystem, req.Name)
	if err != nil {
		logger.Error("Failed to create snapshot", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create snapshot", err))
		return
	}

	logger.Info("Snapshot created", zap.String("snapshotID", snapshot.ID))
	utils.RespondSuccess(w, snapshot)
}

// DeleteSnapshot deletes a snapshot
func (h *BackupHandler) DeleteSnapshot(w http.ResponseWriter, r *http.Request) {
	snapshotID := chi.URLParam(r, "id")

	if err := h.service.DeleteSnapshot(r.Context(), snapshotID); err != nil {
		logger.Error("Failed to delete snapshot", zap.Error(err), zap.String("snapshotID", snapshotID))
		utils.RespondError(w, errors.InternalServerError("Failed to delete snapshot", err))
		return
	}

	logger.Info("Snapshot deleted", zap.String("snapshotID", snapshotID))
	utils.RespondSuccess(w, map[string]string{"message": "Snapshot deleted successfully"})
}

// RestoreSnapshot restores from a snapshot
func (h *BackupHandler) RestoreSnapshot(w http.ResponseWriter, r *http.Request) {
	snapshotID := chi.URLParam(r, "id")

	var req struct {
		Destination string `json:"destination"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.service.RestoreSnapshot(r.Context(), snapshotID, req.Destination); err != nil {
		logger.Error("Failed to restore snapshot", zap.Error(err), zap.String("snapshotID", snapshotID))
		utils.RespondError(w, errors.InternalServerError("Failed to restore snapshot", err))
		return
	}

	logger.Info("Snapshot restored", zap.String("snapshotID", snapshotID))
	utils.RespondSuccess(w, map[string]string{"message": "Snapshot restored successfully"})
}

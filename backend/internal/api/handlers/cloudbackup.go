// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Stumpf-works/stumpfworks-nas/internal/cloudbackup"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// CloudBackupHandler handles cloud backup HTTP requests
type CloudBackupHandler struct {
	service *cloudbackup.Service
}

// NewCloudBackupHandler creates a new cloud backup handler
func NewCloudBackupHandler() *CloudBackupHandler {
	return &CloudBackupHandler{
		service: cloudbackup.GetService(),
	}
}

// CheckAvailability middleware checks if cloud backup service is available
func (h *CloudBackupHandler) CheckAvailability(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.service == nil {
			utils.RespondError(w, errors.NewAppError(503, "Cloud backup service is not available", nil))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Provider Handlers

// ListProviders lists all cloud providers
func (h *CloudBackupHandler) ListProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := h.service.ListProviders()
	if err != nil {
		logger.Error("Failed to list cloud providers", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list providers", err))
		return
	}

	utils.RespondSuccess(w, providers)
}

// GetProvider gets a specific cloud provider
func (h *CloudBackupHandler) GetProvider(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid provider ID", err))
		return
	}

	provider, err := h.service.GetProvider(uint(id))
	if err != nil {
		logger.Error("Failed to get cloud provider", zap.Error(err), zap.Uint64("id", id))
		utils.RespondError(w, errors.NotFound("Provider not found", err))
		return
	}

	utils.RespondSuccess(w, provider)
}

// CreateProvider creates a new cloud provider
func (h *CloudBackupHandler) CreateProvider(w http.ResponseWriter, r *http.Request) {
	var provider models.CloudProvider

	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate input
	if provider.Name == "" {
		utils.RespondError(w, errors.BadRequest("Provider name is required", nil))
		return
	}
	if provider.Type == "" {
		utils.RespondError(w, errors.BadRequest("Provider type is required", nil))
		return
	}
	if provider.Config == "" {
		utils.RespondError(w, errors.BadRequest("Provider configuration is required", nil))
		return
	}

	if err := h.service.CreateProvider(&provider); err != nil {
		logger.Error("Failed to create cloud provider", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create provider", err))
		return
	}

	logger.Info("Cloud provider created", zap.String("name", provider.Name), zap.String("type", provider.Type))
	utils.RespondSuccess(w, provider)
}

// UpdateProvider updates an existing cloud provider
func (h *CloudBackupHandler) UpdateProvider(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid provider ID", err))
		return
	}

	var updates models.CloudProvider
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.service.UpdateProvider(uint(id), &updates); err != nil {
		logger.Error("Failed to update cloud provider", zap.Error(err), zap.Uint64("id", id))
		utils.RespondError(w, errors.InternalServerError("Failed to update provider", err))
		return
	}

	logger.Info("Cloud provider updated", zap.Uint64("id", id))
	utils.RespondSuccess(w, map[string]string{"message": "Provider updated successfully"})
}

// DeleteProvider deletes a cloud provider
func (h *CloudBackupHandler) DeleteProvider(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid provider ID", err))
		return
	}

	if err := h.service.DeleteProvider(uint(id)); err != nil {
		logger.Error("Failed to delete cloud provider", zap.Error(err), zap.Uint64("id", id))
		utils.RespondError(w, errors.InternalServerError("Failed to delete provider", err))
		return
	}

	logger.Info("Cloud provider deleted", zap.Uint64("id", id))
	utils.RespondSuccess(w, map[string]string{"message": "Provider deleted successfully"})
}

// TestProvider tests connectivity to a cloud provider
func (h *CloudBackupHandler) TestProvider(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid provider ID", err))
		return
	}

	if err := h.service.TestProvider(uint(id)); err != nil {
		logger.Error("Provider test failed", zap.Error(err), zap.Uint64("id", id))
		utils.RespondError(w, errors.InternalServerError("Connection test failed", err))
		return
	}

	logger.Info("Provider test successful", zap.Uint64("id", id))
	utils.RespondSuccess(w, map[string]string{"message": "Connection test successful"})
}

// GetProviderTypes returns supported provider types
func (h *CloudBackupHandler) GetProviderTypes(w http.ResponseWriter, r *http.Request) {
	types := h.service.GetProviderTypes()
	utils.RespondSuccess(w, types)
}

// Job Handlers

// ListJobs lists all cloud sync jobs
func (h *CloudBackupHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.service.ListJobs()
	if err != nil {
		logger.Error("Failed to list cloud sync jobs", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list jobs", err))
		return
	}

	utils.RespondSuccess(w, jobs)
}

// GetJob gets a specific cloud sync job
func (h *CloudBackupHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid job ID", err))
		return
	}

	job, err := h.service.GetJob(uint(id))
	if err != nil {
		logger.Error("Failed to get cloud sync job", zap.Error(err), zap.Uint64("id", id))
		utils.RespondError(w, errors.NotFound("Job not found", err))
		return
	}

	utils.RespondSuccess(w, job)
}

// CreateJob creates a new cloud sync job
func (h *CloudBackupHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var job models.CloudSyncJob

	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate input
	if job.Name == "" {
		utils.RespondError(w, errors.BadRequest("Job name is required", nil))
		return
	}
	if job.ProviderID == 0 {
		utils.RespondError(w, errors.BadRequest("Provider ID is required", nil))
		return
	}
	if job.LocalPath == "" {
		utils.RespondError(w, errors.BadRequest("Local path is required", nil))
		return
	}
	if job.RemotePath == "" {
		utils.RespondError(w, errors.BadRequest("Remote path is required", nil))
		return
	}
	if job.Direction == "" {
		utils.RespondError(w, errors.BadRequest("Sync direction is required", nil))
		return
	}

	if err := h.service.CreateJob(&job); err != nil {
		logger.Error("Failed to create cloud sync job", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create job", err))
		return
	}

	logger.Info("Cloud sync job created", zap.String("name", job.Name))
	utils.RespondSuccess(w, job)
}

// UpdateJob updates an existing cloud sync job
func (h *CloudBackupHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid job ID", err))
		return
	}

	var updates models.CloudSyncJob
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.service.UpdateJob(uint(id), &updates); err != nil {
		logger.Error("Failed to update cloud sync job", zap.Error(err), zap.Uint64("id", id))
		utils.RespondError(w, errors.InternalServerError("Failed to update job", err))
		return
	}

	logger.Info("Cloud sync job updated", zap.Uint64("id", id))
	utils.RespondSuccess(w, map[string]string{"message": "Job updated successfully"})
}

// DeleteJob deletes a cloud sync job
func (h *CloudBackupHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid job ID", err))
		return
	}

	if err := h.service.DeleteJob(uint(id)); err != nil {
		logger.Error("Failed to delete cloud sync job", zap.Error(err), zap.Uint64("id", id))
		utils.RespondError(w, errors.InternalServerError("Failed to delete job", err))
		return
	}

	logger.Info("Cloud sync job deleted", zap.Uint64("id", id))
	utils.RespondSuccess(w, map[string]string{"message": "Job deleted successfully"})
}

// RunJob executes a cloud sync job
func (h *CloudBackupHandler) RunJob(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid job ID", err))
		return
	}

	logEntry, err := h.service.RunJob(r.Context(), uint(id), "manual")
	if err != nil {
		logger.Error("Failed to run cloud sync job", zap.Error(err), zap.Uint64("id", id))
		utils.RespondError(w, errors.InternalServerError("Failed to run job", err))
		return
	}

	logger.Info("Cloud sync job started", zap.Uint64("id", id))
	utils.RespondSuccess(w, logEntry)
}

// Log Handlers

// GetLogs gets cloud sync logs
func (h *CloudBackupHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	jobIDStr := r.URL.Query().Get("jobId")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var jobID uint
	if jobIDStr != "" {
		id, err := strconv.ParseUint(jobIDStr, 10, 32)
		if err != nil {
			utils.RespondError(w, errors.BadRequest("Invalid job ID", err))
			return
		}
		jobID = uint(id)
	}

	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	logs, total, err := h.service.GetLogs(jobID, limit, offset)
	if err != nil {
		logger.Error("Failed to get cloud sync logs", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get logs", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"logs":  logs,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

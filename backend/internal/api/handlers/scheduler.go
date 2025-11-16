// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/scheduler"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// SchedulerHandler handles scheduled task-related HTTP requests
type SchedulerHandler struct {
	service *scheduler.Service
}

// NewSchedulerHandler creates a new scheduler handler
func NewSchedulerHandler() *SchedulerHandler {
	return &SchedulerHandler{
		service: scheduler.GetService(),
	}
}

// ListTasks retrieves all scheduled tasks
func (h *SchedulerHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse pagination
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 50
	}

	tasks, total, err := h.service.ListTasks(ctx, offset, limit)
	if err != nil {
		logger.Error("Failed to list tasks", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list tasks", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"tasks":  tasks,
		"total":  total,
		"offset": offset,
		"limit":  limit,
	})
}

// GetTask retrieves a specific task
func (h *SchedulerHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid task ID", err))
		return
	}

	task, err := h.service.GetTask(ctx, uint(id))
	if err != nil {
		logger.Error("Failed to get task", zap.Error(err))
		utils.RespondError(w, errors.NotFound("Task not found", err))
		return
	}

	utils.RespondSuccess(w, task)
}

// CreateTask creates a new scheduled task
func (h *SchedulerHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var task models.ScheduledTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate required fields
	if task.Name == "" {
		utils.RespondError(w, errors.BadRequest("Task name is required", nil))
		return
	}

	if task.CronExpression == "" {
		utils.RespondError(w, errors.BadRequest("Cron expression is required", nil))
		return
	}

	if task.TaskType == "" {
		utils.RespondError(w, errors.BadRequest("Task type is required", nil))
		return
	}

	if err := h.service.CreateTask(ctx, &task); err != nil {
		logger.Error("Failed to create task", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create task", err))
		return
	}

	utils.RespondSuccess(w, task)
}

// UpdateTask updates a scheduled task
func (h *SchedulerHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid task ID", err))
		return
	}

	var task models.ScheduledTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	task.ID = uint(id)

	if err := h.service.UpdateTask(ctx, &task); err != nil {
		logger.Error("Failed to update task", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to update task", err))
		return
	}

	utils.RespondSuccess(w, task)
}

// DeleteTask deletes a scheduled task
func (h *SchedulerHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid task ID", err))
		return
	}

	if err := h.service.DeleteTask(ctx, uint(id)); err != nil {
		logger.Error("Failed to delete task", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to delete task", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Task deleted successfully",
	})
}

// RunTaskNow executes a task immediately
func (h *SchedulerHandler) RunTaskNow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid task ID", err))
		return
	}

	if err := h.service.RunTaskNow(ctx, uint(id)); err != nil {
		logger.Error("Failed to run task", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to run task", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Task execution started",
	})
}

// GetTaskExecutions retrieves execution history for a task
func (h *SchedulerHandler) GetTaskExecutions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid task ID", err))
		return
	}

	// Parse pagination
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 50
	}

	executions, total, err := h.service.GetTaskExecutions(ctx, uint(id), offset, limit)
	if err != nil {
		logger.Error("Failed to get task executions", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get task executions", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"executions": executions,
		"total":      total,
		"offset":     offset,
		"limit":      limit,
	})
}

// ValidateCron validates a cron expression
func (h *SchedulerHandler) ValidateCron(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := scheduler.ValidateCronExpression(req.Expression); err != nil {
		utils.RespondSuccess(w, map[string]interface{}{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	// Calculate next run times
	schedule, _ := scheduler.ParseCronExpression(req.Expression)
	now := time.Now()
	nextRuns := make([]string, 5)
	current := now
	for i := 0; i < 5; i++ {
		current = schedule.Next(current)
		nextRuns[i] = current.Format("2006-01-02 15:04:05")
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"valid":     true,
		"nextRuns":  nextRuns,
	})
}

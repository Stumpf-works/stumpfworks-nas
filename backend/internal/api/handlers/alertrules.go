// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Stumpf-works/stumpfworks-nas/internal/alertrules"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// AlertRulesHandler handles alert rules API requests
type AlertRulesHandler struct {
	service *alertrules.Service
}

// NewAlertRulesHandler creates a new alert rules handler
func NewAlertRulesHandler() *AlertRulesHandler {
	return &AlertRulesHandler{
		service: alertrules.GetService(),
	}
}

// CreateRule creates a new alert rule
func (h *AlertRulesHandler) CreateRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var rule models.AlertRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate rule
	if rule.Name == "" {
		utils.RespondError(w, errors.BadRequest("Rule name is required", nil))
		return
	}
	if rule.MetricType == "" {
		utils.RespondError(w, errors.BadRequest("Metric type is required", nil))
		return
	}
	if rule.Condition == "" {
		utils.RespondError(w, errors.BadRequest("Condition is required", nil))
		return
	}

	if err := h.service.CreateRule(ctx, &rule); err != nil {
		logger.Error("Failed to create alert rule", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create alert rule", err))
		return
	}

	utils.RespondSuccess(w, rule)
}

// UpdateRule updates an existing alert rule
func (h *AlertRulesHandler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ruleID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid rule ID", err))
		return
	}

	var rule models.AlertRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	rule.ID = uint(ruleID)

	if err := h.service.UpdateRule(ctx, &rule); err != nil {
		logger.Error("Failed to update alert rule", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to update alert rule", err))
		return
	}

	utils.RespondSuccess(w, rule)
}

// DeleteRule deletes an alert rule
func (h *AlertRulesHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ruleID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid rule ID", err))
		return
	}

	if err := h.service.DeleteRule(ctx, uint(ruleID)); err != nil {
		logger.Error("Failed to delete alert rule", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to delete alert rule", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Alert rule deleted successfully",
	})
}

// GetRule retrieves a single alert rule
func (h *AlertRulesHandler) GetRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ruleID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid rule ID", err))
		return
	}

	rule, err := h.service.GetRule(ctx, uint(ruleID))
	if err != nil {
		logger.Error("Failed to get alert rule", zap.Error(err))
		utils.RespondError(w, errors.NotFound("Alert rule not found", err))
		return
	}

	utils.RespondSuccess(w, rule)
}

// ListRules retrieves all alert rules
func (h *AlertRulesHandler) ListRules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rules, err := h.service.ListRules(ctx)
	if err != nil {
		logger.Error("Failed to list alert rules", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list alert rules", err))
		return
	}

	utils.RespondSuccess(w, rules)
}

// GetExecutions retrieves executions for a specific rule
func (h *AlertRulesHandler) GetExecutions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ruleID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid rule ID", err))
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50 // Default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	executions, err := h.service.GetExecutions(ctx, uint(ruleID), limit)
	if err != nil {
		logger.Error("Failed to get rule executions", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get rule executions", err))
		return
	}

	utils.RespondSuccess(w, executions)
}

// GetRecentExecutions retrieves recent executions across all rules
func (h *AlertRulesHandler) GetRecentExecutions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limitStr := r.URL.Query().Get("limit")
	limit := 100 // Default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	executions, err := h.service.GetRecentExecutions(ctx, limit)
	if err != nil {
		logger.Error("Failed to get recent executions", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get recent executions", err))
		return
	}

	utils.RespondSuccess(w, executions)
}

// AcknowledgeExecution acknowledges an alert execution
func (h *AlertRulesHandler) AcknowledgeExecution(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	executionID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid execution ID", err))
		return
	}

	var input struct {
		AcknowledgedBy string `json:"acknowledgedBy"`
		Note           string `json:"note"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if input.AcknowledgedBy == "" {
		utils.RespondError(w, errors.BadRequest("acknowledgedBy is required", nil))
		return
	}

	if err := h.service.AcknowledgeExecution(ctx, uint(executionID), input.AcknowledgedBy, input.Note); err != nil {
		logger.Error("Failed to acknowledge execution", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to acknowledge execution", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Execution acknowledged successfully",
	})
}

// Revision: 2025-11-28 | Author: Claude | Version: 1.0.0
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/filesystem"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

var quotaManager *filesystem.QuotaManager

// InitQuotaManager initializes the quota manager
func InitQuotaManager(qm *filesystem.QuotaManager) {
	quotaManager = qm
	logger.Info("Quota manager initialized")
}

// ===== Request/Response Structures =====

// GetQuotaRequest represents the request for getting quota info
type GetQuotaRequest struct {
	Name       string                 `json:"name"`       // username or groupname
	Type       filesystem.QuotaType   `json:"type"`       // user or group
	Filesystem string                 `json:"filesystem"` // filesystem path
}

// SetQuotaRequest represents the request for setting quota
type SetQuotaRequest struct {
	Name       string                  `json:"name"`       // username or groupname
	Type       filesystem.QuotaType    `json:"type"`       // user or group
	Filesystem string                  `json:"filesystem"` // filesystem path
	Limits     filesystem.QuotaLimits  `json:"limits"`     // quota limits
}

// RemoveQuotaRequest represents the request for removing quota
type RemoveQuotaRequest struct {
	Name       string                 `json:"name"`       // username or groupname
	Type       filesystem.QuotaType   `json:"type"`       // user or group
	Filesystem string                 `json:"filesystem"` // filesystem path
}

// ListQuotasRequest represents the request for listing quotas
type ListQuotasRequest struct {
	Filesystem string                 `json:"filesystem"` // filesystem path
	Type       filesystem.QuotaType   `json:"type"`       // user or group
}

// ===== Quota Handlers =====

// GetUserQuota retrieves quota information for a user
// GET /api/v1/quotas/user?name=username&filesystem=/path
func GetUserQuota(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("name")
	filesystem := r.URL.Query().Get("filesystem")

	if username == "" {
		utils.RespondError(w, errors.BadRequest("Missing name parameter", nil))
		return
	}

	if filesystem == "" {
		utils.RespondError(w, errors.BadRequest("Missing filesystem parameter", nil))
		return
	}

	if quotaManager == nil || !quotaManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("Quota support not available", nil))
		return
	}

	quota, err := quotaManager.GetUserQuota(username, filesystem)
	if err != nil {
		logger.Error("Failed to get user quota",
			zap.String("username", username),
			zap.String("filesystem", filesystem),
			zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get user quota", err))
		return
	}

	utils.RespondSuccess(w, quota)
}

// GetGroupQuota retrieves quota information for a group
// GET /api/v1/quotas/group?name=groupname&filesystem=/path
func GetGroupQuota(w http.ResponseWriter, r *http.Request) {
	groupname := r.URL.Query().Get("name")
	filesystem := r.URL.Query().Get("filesystem")

	if groupname == "" {
		utils.RespondError(w, errors.BadRequest("Missing name parameter", nil))
		return
	}

	if filesystem == "" {
		utils.RespondError(w, errors.BadRequest("Missing filesystem parameter", nil))
		return
	}

	if quotaManager == nil || !quotaManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("Quota support not available", nil))
		return
	}

	quota, err := quotaManager.GetGroupQuota(groupname, filesystem)
	if err != nil {
		logger.Error("Failed to get group quota",
			zap.String("groupname", groupname),
			zap.String("filesystem", filesystem),
			zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get group quota", err))
		return
	}

	utils.RespondSuccess(w, quota)
}

// SetUserQuota sets quota limits for a user
// POST /api/v1/quotas/user
// Body: { "name": "username", "filesystem": "/path", "limits": {...} }
func SetUserQuota(w http.ResponseWriter, r *http.Request) {
	var req SetQuotaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Name == "" {
		utils.RespondError(w, errors.BadRequest("Missing name in request", nil))
		return
	}

	if req.Filesystem == "" {
		utils.RespondError(w, errors.BadRequest("Missing filesystem in request", nil))
		return
	}

	if quotaManager == nil || !quotaManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("Quota support not available", nil))
		return
	}

	if err := quotaManager.SetUserQuota(req.Name, req.Filesystem, req.Limits); err != nil {
		logger.Error("Failed to set user quota",
			zap.String("username", req.Name),
			zap.String("filesystem", req.Filesystem),
			zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to set user quota", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "User quota set successfully",
		"name":    req.Name,
	})
}

// SetGroupQuota sets quota limits for a group
// POST /api/v1/quotas/group
// Body: { "name": "groupname", "filesystem": "/path", "limits": {...} }
func SetGroupQuota(w http.ResponseWriter, r *http.Request) {
	var req SetQuotaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Name == "" {
		utils.RespondError(w, errors.BadRequest("Missing name in request", nil))
		return
	}

	if req.Filesystem == "" {
		utils.RespondError(w, errors.BadRequest("Missing filesystem in request", nil))
		return
	}

	if quotaManager == nil || !quotaManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("Quota support not available", nil))
		return
	}

	if err := quotaManager.SetGroupQuota(req.Name, req.Filesystem, req.Limits); err != nil {
		logger.Error("Failed to set group quota",
			zap.String("groupname", req.Name),
			zap.String("filesystem", req.Filesystem),
			zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to set group quota", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Group quota set successfully",
		"name":    req.Name,
	})
}

// RemoveUserQuota removes quota limits for a user
// DELETE /api/v1/quotas/user
// Body: { "name": "username", "filesystem": "/path" }
func RemoveUserQuota(w http.ResponseWriter, r *http.Request) {
	var req RemoveQuotaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Name == "" {
		utils.RespondError(w, errors.BadRequest("Missing name in request", nil))
		return
	}

	if req.Filesystem == "" {
		utils.RespondError(w, errors.BadRequest("Missing filesystem in request", nil))
		return
	}

	if quotaManager == nil || !quotaManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("Quota support not available", nil))
		return
	}

	if err := quotaManager.RemoveUserQuota(req.Name, req.Filesystem); err != nil {
		logger.Error("Failed to remove user quota",
			zap.String("username", req.Name),
			zap.String("filesystem", req.Filesystem),
			zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to remove user quota", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "User quota removed successfully",
		"name":    req.Name,
	})
}

// RemoveGroupQuota removes quota limits for a group
// DELETE /api/v1/quotas/group
// Body: { "name": "groupname", "filesystem": "/path" }
func RemoveGroupQuota(w http.ResponseWriter, r *http.Request) {
	var req RemoveQuotaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Name == "" {
		utils.RespondError(w, errors.BadRequest("Missing name in request", nil))
		return
	}

	if req.Filesystem == "" {
		utils.RespondError(w, errors.BadRequest("Missing filesystem in request", nil))
		return
	}

	if quotaManager == nil || !quotaManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("Quota support not available", nil))
		return
	}

	if err := quotaManager.RemoveGroupQuota(req.Name, req.Filesystem); err != nil {
		logger.Error("Failed to remove group quota",
			zap.String("groupname", req.Name),
			zap.String("filesystem", req.Filesystem),
			zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to remove group quota", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Group quota removed successfully",
		"name":    req.Name,
	})
}

// ListUserQuotas lists all user quotas on a filesystem
// GET /api/v1/quotas/users?filesystem=/path
func ListUserQuotas(w http.ResponseWriter, r *http.Request) {
	filesystem := r.URL.Query().Get("filesystem")

	if filesystem == "" {
		utils.RespondError(w, errors.BadRequest("Missing filesystem parameter", nil))
		return
	}

	if quotaManager == nil || !quotaManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("Quota support not available", nil))
		return
	}

	quotas, err := quotaManager.ListUserQuotas(filesystem)
	if err != nil {
		logger.Error("Failed to list user quotas",
			zap.String("filesystem", filesystem),
			zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list user quotas", err))
		return
	}

	utils.RespondSuccess(w, quotas)
}

// ListGroupQuotas lists all group quotas on a filesystem
// GET /api/v1/quotas/groups?filesystem=/path
func ListGroupQuotas(w http.ResponseWriter, r *http.Request) {
	filesystem := r.URL.Query().Get("filesystem")

	if filesystem == "" {
		utils.RespondError(w, errors.BadRequest("Missing filesystem parameter", nil))
		return
	}

	if quotaManager == nil || !quotaManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("Quota support not available", nil))
		return
	}

	quotas, err := quotaManager.ListGroupQuotas(filesystem)
	if err != nil {
		logger.Error("Failed to list group quotas",
			zap.String("filesystem", filesystem),
			zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list group quotas", err))
		return
	}

	utils.RespondSuccess(w, quotas)
}

// GetFilesystemQuotaStatus checks if quotas are enabled on a filesystem
// GET /api/v1/quotas/status?filesystem=/path
func GetFilesystemQuotaStatus(w http.ResponseWriter, r *http.Request) {
	filesystem := r.URL.Query().Get("filesystem")

	if filesystem == "" {
		utils.RespondError(w, errors.BadRequest("Missing filesystem parameter", nil))
		return
	}

	if quotaManager == nil || !quotaManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("Quota support not available", nil))
		return
	}

	status, err := quotaManager.GetFilesystemQuotaStatus(filesystem)
	if err != nil {
		logger.Error("Failed to get filesystem quota status",
			zap.String("filesystem", filesystem),
			zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get quota status", err))
		return
	}

	utils.RespondSuccess(w, status)
}

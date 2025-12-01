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

var aclManager *filesystem.ACLManager

// InitACLManager initializes the ACL manager
func InitACLManager(acl *filesystem.ACLManager) {
	aclManager = acl
	logger.Info("ACL manager initialized")
}

// ===== Request/Response Structures =====

// GetACLRequest represents the request for getting ACLs
type GetACLRequest struct {
	Path string `json:"path"`
}

// SetACLRequest represents the request for setting ACLs
type SetACLRequest struct {
	Path    string                    `json:"path"`
	Entries []filesystem.ACLEntry     `json:"entries"`
}

// RemoveACLRequest represents the request for removing an ACL entry
type RemoveACLRequest struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// SetDefaultACLRequest represents the request for setting default ACLs
type SetDefaultACLRequest struct {
	DirPath string                    `json:"dir_path"`
	Entries []filesystem.ACLEntry     `json:"entries"`
}

// ApplyRecursiveRequest represents the request for applying ACLs recursively
type ApplyRecursiveRequest struct {
	DirPath string                    `json:"dir_path"`
	Entries []filesystem.ACLEntry     `json:"entries"`
}

// ===== ACL Handlers =====

// GetACL retrieves ACL entries for a file or directory
// GET /api/v1/filesystem/acl?path=/path/to/file
func GetACL(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		utils.RespondError(w, errors.BadRequest("Missing path parameter", nil))
		return
	}

	if aclManager == nil || !aclManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("ACL support not available", nil))
		return
	}

	entries, err := aclManager.GetACL(path)
	if err != nil {
		logger.Error("Failed to get ACL", zap.String("path", path), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get ACL", err))
		return
	}

	response := filesystem.ACLInfo{
		Path:    path,
		Entries: entries,
	}

	utils.RespondSuccess(w, response)
}

// SetACL sets ACL entries on a file or directory
// POST /api/v1/filesystem/acl
// Body: { "path": "/path/to/file", "entries": [...] }
func SetACL(w http.ResponseWriter, r *http.Request) {
	var req SetACLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Path == "" {
		utils.RespondError(w, errors.BadRequest("Missing path in request", nil))
		return
	}

	if len(req.Entries) == 0 {
		utils.RespondError(w, errors.BadRequest("No ACL entries provided", nil))
		return
	}

	if aclManager == nil || !aclManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("ACL support not available", nil))
		return
	}

	if err := aclManager.SetACL(req.Path, req.Entries); err != nil {
		logger.Error("Failed to set ACL", zap.String("path", req.Path), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to set ACL", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "ACL entries set successfully",
		"path":    req.Path,
	})
}

// RemoveACL removes a specific ACL entry from a file or directory
// DELETE /api/v1/filesystem/acl
// Body: { "path": "/path/to/file", "type": "user", "name": "alice" }
func RemoveACL(w http.ResponseWriter, r *http.Request) {
	var req RemoveACLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Path == "" {
		utils.RespondError(w, errors.BadRequest("Missing path in request", nil))
		return
	}

	if req.Type == "" {
		utils.RespondError(w, errors.BadRequest("Missing ACL type in request", nil))
		return
	}

	if aclManager == nil || !aclManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("ACL support not available", nil))
		return
	}

	if err := aclManager.RemoveACL(req.Path, req.Type, req.Name); err != nil {
		logger.Error("Failed to remove ACL entry",
			zap.String("path", req.Path),
			zap.String("type", req.Type),
			zap.String("name", req.Name),
			zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to remove ACL entry", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "ACL entry removed successfully",
		"path":    req.Path,
	})
}

// SetDefaultACL sets default ACL entries for new files created in a directory
// POST /api/v1/filesystem/acl/default
// Body: { "dir_path": "/path/to/dir", "entries": [...] }
func SetDefaultACL(w http.ResponseWriter, r *http.Request) {
	var req SetDefaultACLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.DirPath == "" {
		utils.RespondError(w, errors.BadRequest("Missing dir_path in request", nil))
		return
	}

	if len(req.Entries) == 0 {
		utils.RespondError(w, errors.BadRequest("No ACL entries provided", nil))
		return
	}

	if aclManager == nil || !aclManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("ACL support not available", nil))
		return
	}

	if err := aclManager.SetDefaultACL(req.DirPath, req.Entries); err != nil {
		logger.Error("Failed to set default ACL", zap.String("dir_path", req.DirPath), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to set default ACL", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Default ACL entries set successfully",
		"path":    req.DirPath,
	})
}

// ApplyRecursive applies ACL entries recursively to a directory and all its contents
// POST /api/v1/filesystem/acl/recursive
// Body: { "dir_path": "/path/to/dir", "entries": [...] }
func ApplyRecursive(w http.ResponseWriter, r *http.Request) {
	var req ApplyRecursiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.DirPath == "" {
		utils.RespondError(w, errors.BadRequest("Missing dir_path in request", nil))
		return
	}

	if len(req.Entries) == 0 {
		utils.RespondError(w, errors.BadRequest("No ACL entries provided", nil))
		return
	}

	if aclManager == nil || !aclManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("ACL support not available", nil))
		return
	}

	if err := aclManager.ApplyRecursive(req.DirPath, req.Entries); err != nil {
		logger.Error("Failed to apply ACLs recursively", zap.String("dir_path", req.DirPath), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to apply ACLs recursively", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "ACL entries applied recursively",
		"path":    req.DirPath,
	})
}

// RemoveAllACLs removes all ACL entries from a file or directory
// DELETE /api/v1/filesystem/acl/all
// Body: { "path": "/path/to/file" }
func RemoveAllACLs(w http.ResponseWriter, r *http.Request) {
	var req GetACLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Path == "" {
		utils.RespondError(w, errors.BadRequest("Missing path in request", nil))
		return
	}

	if aclManager == nil || !aclManager.IsEnabled() {
		utils.RespondError(w, errors.InternalServerError("ACL support not available", nil))
		return
	}

	if err := aclManager.RemoveAllACLs(req.Path); err != nil {
		logger.Error("Failed to remove all ACLs", zap.String("path", req.Path), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to remove all ACLs", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "All ACL entries removed successfully",
		"path":    req.Path,
	})
}

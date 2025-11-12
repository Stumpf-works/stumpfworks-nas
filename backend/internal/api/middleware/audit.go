package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/audit"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// AuditMiddleware logs important actions to the audit log
func AuditMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip audit logging for certain paths
		if shouldSkipAuditLog(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Get user from context
		user := GetUserFromContext(r.Context())

		// Create a response writer wrapper to capture status code
		wrappedWriter := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Serve the request
		next.ServeHTTP(wrappedWriter, r)

		// Determine if we should log this request
		if shouldLogRequest(r.Method, r.URL.Path, wrappedWriter.statusCode) {
			// Log the action asynchronously
			go logAuditEntry(r, user, wrappedWriter.statusCode)
		}
	})
}

// responseWriterWrapper wraps http.ResponseWriter to capture the status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// shouldSkipAuditLog determines if a path should be skipped from audit logging
func shouldSkipAuditLog(path string) bool {
	// Skip health checks, metrics, and audit log endpoints themselves
	skipPrefixes := []string{
		"/health",
		"/api/v1/system/metrics",
		"/api/v1/audit/logs",
		"/ws",
	}

	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

// shouldLogRequest determines if a request should be logged based on method, path, and status
func shouldLogRequest(method, path string, statusCode int) bool {
	// Log all POST, PUT, DELETE, PATCH requests (write operations)
	if method == http.MethodPost || method == http.MethodPut ||
		method == http.MethodDelete || method == http.MethodPatch {
		return true
	}

	// Log failed GET requests (4xx, 5xx)
	if method == http.MethodGet && statusCode >= 400 {
		return true
	}

	return false
}

// logAuditEntry creates an audit log entry for the request
func logAuditEntry(r *http.Request, user *models.User, statusCode int) {
	auditService := audit.GetService()
	if auditService == nil {
		return
	}

	// Determine action and resource from path and method
	action, resource := inferActionAndResource(r.Method, r.URL.Path)

	// Determine status
	status := models.StatusSuccess
	if statusCode >= 400 {
		status = models.StatusFailure
	}
	if statusCode >= 500 {
		status = models.StatusError
	}

	// Determine severity
	severity := models.SeverityInfo
	if statusCode >= 400 && statusCode < 500 {
		severity = models.SeverityWarning
	}
	if statusCode >= 500 {
		severity = models.SeverityCritical
	}

	// Get user info
	var userID *uint
	username := "anonymous"
	if user != nil {
		userID = &user.ID
		username = user.Username
	}

	// Create audit log entry
	entry := &audit.LogEntry{
		UserID:    userID,
		Username:  username,
		Action:    action,
		Resource:  resource,
		Status:    status,
		Severity:  severity,
		IPAddress: getClientIP(r),
		UserAgent: r.UserAgent(),
		Message:   generateAuditMessage(action, resource, status),
	}

	// Log the entry
	if err := auditService.Log(context.Background(), entry); err != nil {
		logger.Error("Failed to create audit log entry",
			zap.Error(err),
			zap.String("action", action),
			zap.String("resource", resource))
	}
}

// inferActionAndResource attempts to infer the action and resource from the HTTP method and path
func inferActionAndResource(method, path string) (action, resource string) {
	// Remove /api/v1 prefix
	path = strings.TrimPrefix(path, "/api/v1")

	// Split path into segments
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) == 0 {
		return "unknown", path
	}

	mainResource := segments[0]

	// Map common patterns
	switch {
	// Authentication
	case strings.HasPrefix(path, "/auth/login"):
		return models.ActionAuthLogin, "auth"
	case strings.HasPrefix(path, "/auth/logout"):
		return models.ActionAuthLogout, "auth"
	case strings.HasPrefix(path, "/auth/refresh"):
		return models.ActionAuthTokenRefresh, "auth"

	// Users
	case strings.HasPrefix(path, "/users"):
		switch method {
		case http.MethodPost:
			return models.ActionUserCreate, path
		case http.MethodPut, http.MethodPatch:
			return models.ActionUserUpdate, path
		case http.MethodDelete:
			return models.ActionUserDelete, path
		}

	// Files
	case strings.HasPrefix(path, "/files"):
		switch {
		case strings.Contains(path, "/upload"):
			return models.ActionFileUpload, path
		case strings.Contains(path, "/delete"):
			return models.ActionFileDelete, path
		case strings.Contains(path, "/rename"):
			return models.ActionFileRename, path
		case strings.Contains(path, "/move"):
			return models.ActionFileMove, path
		case strings.Contains(path, "/copy"):
			return models.ActionFileCopy, path
		}

	// Storage
	case strings.HasPrefix(path, "/storage"):
		switch {
		case strings.Contains(path, "/volumes") && method == http.MethodPost:
			return models.ActionStorageVolumeCreate, path
		case strings.Contains(path, "/volumes") && method == http.MethodDelete:
			return models.ActionStorageVolumeDelete, path
		case strings.Contains(path, "/shares") && method == http.MethodPost:
			return models.ActionStorageShareCreate, path
		case strings.Contains(path, "/shares") && method == http.MethodPut:
			return models.ActionStorageShareUpdate, path
		case strings.Contains(path, "/shares") && method == http.MethodDelete:
			return models.ActionStorageShareDelete, path
		}

	// Docker
	case strings.HasPrefix(path, "/docker/containers"):
		switch {
		case strings.Contains(path, "/start"):
			return models.ActionDockerContainerStart, path
		case strings.Contains(path, "/stop"):
			return models.ActionDockerContainerStop, path
		case strings.Contains(path, "/remove") || method == http.MethodDelete:
			return models.ActionDockerContainerRemove, path
		}

	// Active Directory
	case strings.HasPrefix(path, "/ad"):
		switch {
		case strings.Contains(path, "/config"):
			return models.ActionADConfigUpdate, "ad_config"
		case strings.Contains(path, "/sync"):
			return models.ActionADSync, "ad_users"
		}
	}

	// Default: construct action from method and resource
	actionPrefix := strings.ToLower(mainResource)
	actionSuffix := strings.ToLower(method)

	switch method {
	case http.MethodPost:
		actionSuffix = "create"
	case http.MethodPut, http.MethodPatch:
		actionSuffix = "update"
	case http.MethodDelete:
		actionSuffix = "delete"
	case http.MethodGet:
		actionSuffix = "read"
	}

	return actionPrefix + "." + actionSuffix, path
}

// generateAuditMessage generates a human-readable message for the audit log
func generateAuditMessage(action, resource, status string) string {
	statusText := "completed"
	if status == models.StatusFailure {
		statusText = "failed"
	} else if status == models.StatusError {
		statusText = "errored"
	}

	return action + " " + statusText + " for " + resource
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP if multiple are present
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

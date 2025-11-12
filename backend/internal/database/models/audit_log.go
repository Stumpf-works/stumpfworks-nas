package models

import (
	"time"
)

// AuditLog represents a security audit log entry
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`

	// User information
	UserID   *uint  `gorm:"index" json:"userId,omitempty"` // Nullable for anonymous/system actions
	Username string `gorm:"size:100;index" json:"username"`

	// Action details
	Action   string `gorm:"size:100;not null;index" json:"action"` // e.g., "auth.login", "file.delete"
	Resource string `gorm:"size:255" json:"resource,omitempty"`     // e.g., "file:/path", "user:123"
	Status   string `gorm:"size:20;not null" json:"status"`         // success, failure, error

	// Severity level
	Severity string `gorm:"size:20;not null;index" json:"severity"` // info, warning, critical

	// Request context
	IPAddress string `gorm:"size:45" json:"ipAddress,omitempty"` // IPv6 max length
	UserAgent string `gorm:"size:500" json:"userAgent,omitempty"`

	// Additional data
	Details string `gorm:"type:text" json:"details,omitempty"` // JSON string for additional context
	Message string `gorm:"size:500" json:"message,omitempty"`  // Human-readable message
}

// TableName specifies the table name for AuditLog model
func (AuditLog) TableName() string {
	return "audit_logs"
}

// Common action constants
const (
	// Authentication actions
	ActionAuthLogin        = "auth.login"
	ActionAuthLogout       = "auth.logout"
	ActionAuthLoginFailed  = "auth.login_failed"
	ActionAuthTokenRefresh = "auth.token_refresh"

	// User management actions
	ActionUserCreate = "user.create"
	ActionUserUpdate = "user.update"
	ActionUserDelete = "user.delete"

	// File actions
	ActionFileUpload = "file.upload"
	ActionFileDelete = "file.delete"
	ActionFileRename = "file.rename"
	ActionFileMove   = "file.move"
	ActionFileCopy   = "file.copy"

	// System actions
	ActionSystemConfigUpdate = "system.config_update"
	ActionSystemRestart      = "system.restart"
	ActionSystemShutdown     = "system.shutdown"

	// Storage actions
	ActionStorageVolumeCreate = "storage.volume_create"
	ActionStorageVolumeDelete = "storage.volume_delete"
	ActionStorageShareCreate  = "storage.share_create"
	ActionStorageShareUpdate  = "storage.share_update"
	ActionStorageShareDelete  = "storage.share_delete"

	// Docker actions
	ActionDockerContainerStart  = "docker.container_start"
	ActionDockerContainerStop   = "docker.container_stop"
	ActionDockerContainerRemove = "docker.container_remove"

	// AD actions
	ActionADConfigUpdate = "ad.config_update"
	ActionADSync         = "ad.sync"
)

// Severity levels
const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityCritical = "critical"
)

// Status values
const (
	StatusSuccess = "success"
	StatusFailure = "failure"
	StatusError   = "error"
)

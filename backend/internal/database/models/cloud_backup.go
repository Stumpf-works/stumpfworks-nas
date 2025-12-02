// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package models

import (
	"time"

	"gorm.io/gorm"
)

// CloudProvider represents a cloud storage provider configuration
type CloudProvider struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Type        string         `gorm:"type:varchar(50);not null" json:"type"` // s3, b2, gdrive, dropbox, onedrive, azureblob, sftp
	Description string         `gorm:"type:text" json:"description"`
	Config      string         `gorm:"type:text;not null" json:"config"` // JSON-encoded provider-specific configuration
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	TestStatus  string         `gorm:"type:varchar(50)" json:"test_status"` // untested, success, failed
	TestedAt    *time.Time     `json:"tested_at,omitempty"`
}

// CloudSyncJob represents a scheduled cloud synchronization job
type CloudSyncJob struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	Name              string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Description       string         `gorm:"type:text" json:"description"`
	Enabled           bool           `gorm:"default:true" json:"enabled"`
	ProviderID        uint           `gorm:"not null;index" json:"provider_id"`
	Provider          *CloudProvider `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	Direction         string         `gorm:"type:varchar(20);not null" json:"direction"` // upload, download, sync
	LocalPath         string         `gorm:"type:varchar(500);not null" json:"local_path"`
	RemotePath        string         `gorm:"type:varchar(500);not null" json:"remote_path"`
	Schedule          string         `gorm:"type:varchar(100)" json:"schedule"` // cron expression
	ScheduleEnabled   bool           `gorm:"default:false" json:"schedule_enabled"`
	BandwidthLimit    string         `gorm:"type:varchar(20)" json:"bandwidth_limit"` // e.g., "10M", "1G"
	EncryptionEnabled bool           `gorm:"default:false" json:"encryption_enabled"`
	EncryptionKey     string         `gorm:"type:text" json:"-"` // Encrypted encryption key
	DeleteAfterUpload bool           `gorm:"default:false" json:"delete_after_upload"`
	Retention         int            `gorm:"default:0" json:"retention"` // days to keep backups (0 = forever)
	Filters           string         `gorm:"type:text" json:"filters"` // JSON array of include/exclude patterns
	LastRunAt         *time.Time     `json:"last_run_at,omitempty"`
	NextRunAt         *time.Time     `json:"next_run_at,omitempty"`
	LastStatus        string         `gorm:"type:varchar(50)" json:"last_status"` // idle, running, success, failed
	FailureCount      int            `gorm:"default:0" json:"failure_count"`
}

// CloudSyncLog represents a cloud sync execution log entry
type CloudSyncLog struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	JobID           uint           `gorm:"not null;index" json:"job_id"`
	Job             *CloudSyncJob  `gorm:"foreignKey:JobID" json:"job,omitempty"`
	JobName         string         `gorm:"type:varchar(100)" json:"job_name"`
	StartedAt       time.Time      `json:"started_at"`
	CompletedAt     *time.Time     `json:"completed_at,omitempty"`
	Status          string         `gorm:"type:varchar(50);not null" json:"status"` // running, success, failed, cancelled
	Direction       string         `gorm:"type:varchar(20)" json:"direction"` // upload, download, sync
	BytesTransferred int64         `gorm:"default:0" json:"bytes_transferred"`
	FilesTransferred int           `gorm:"default:0" json:"files_transferred"`
	FilesDeleted     int           `gorm:"default:0" json:"files_deleted"`
	FilesFailed      int           `gorm:"default:0" json:"files_failed"`
	Duration         int64         `gorm:"default:0" json:"duration"` // seconds
	ErrorMessage     string        `gorm:"type:text" json:"error_message,omitempty"`
	Output           string        `gorm:"type:text" json:"output,omitempty"` // rclone output
	TriggeredBy      string        `gorm:"type:varchar(50)" json:"triggered_by"` // manual, schedule, api
}

// TableName overrides the table name
func (CloudProvider) TableName() string {
	return "cloud_providers"
}

func (CloudSyncJob) TableName() string {
	return "cloud_sync_jobs"
}

func (CloudSyncLog) TableName() string {
	return "cloud_sync_logs"
}

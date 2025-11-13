package models

import "time"

// ScheduledTask represents a scheduled task configuration
type ScheduledTask struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Task identification
	Name        string `gorm:"size:255;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	TaskType    string `gorm:"size:100;not null;index" json:"taskType"` // cleanup, backup, maintenance, custom

	// Scheduling
	CronExpression string `gorm:"size:255;not null" json:"cronExpression"` // Standard cron format
	Enabled        bool   `gorm:"default:true;index" json:"enabled"`

	// Execution tracking
	LastRun     *time.Time `json:"lastRun,omitempty"`
	NextRun     *time.Time `json:"nextRun,omitempty"`
	LastStatus  string     `gorm:"size:50" json:"lastStatus,omitempty"` // success, failed, running
	LastError   string     `gorm:"type:text" json:"lastError,omitempty"`
	RunCount    int        `gorm:"default:0" json:"runCount"`

	// Task configuration (JSON)
	Config string `gorm:"type:text" json:"config,omitempty"` // Task-specific config as JSON

	// Timeout and retry
	TimeoutSeconds int  `gorm:"default:300" json:"timeoutSeconds"` // 5 minutes default
	RetryOnFailure bool `gorm:"default:false" json:"retryOnFailure"`
	MaxRetries     int  `gorm:"default:3" json:"maxRetries"`
}

// TaskExecution represents a single execution of a scheduled task
type TaskExecution struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`

	TaskID      uint      `gorm:"not null;index" json:"taskId"`
	Task        *ScheduledTask `gorm:"foreignKey:TaskID" json:"task,omitempty"`

	// Execution details
	StartedAt   time.Time  `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	Duration    int64      `json:"duration"` // Duration in milliseconds

	Status string `gorm:"size:50;not null;index" json:"status"` // running, success, failed, timeout
	Output string `gorm:"type:text" json:"output,omitempty"`
	Error  string `gorm:"type:text" json:"error,omitempty"`

	// Metadata
	TriggeredBy string `gorm:"size:100" json:"triggeredBy"` // scheduler, manual, api
	RetryCount  int    `gorm:"default:0" json:"retryCount"`
}

// Task types
const (
	TaskTypeCleanup     = "cleanup"
	TaskTypeBackup      = "backup"
	TaskTypeMaintenance = "maintenance"
	TaskTypeCustom      = "custom"
	TaskTypeLogRotation = "log_rotation"
	TaskTypeMetrics     = "metrics"
)

// Task status
const (
	TaskStatusRunning = "running"
	TaskStatusSuccess = "success"
	TaskStatusFailed  = "failed"
	TaskStatusTimeout = "timeout"
)

// Triggered by
const (
	TriggerScheduler = "scheduler"
	TriggerManual    = "manual"
	TriggerAPI       = "api"
)

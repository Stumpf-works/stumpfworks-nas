package models

import (
	"time"
)

// SystemMetric stores historical system metrics
type SystemMetric struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Timestamp time.Time `gorm:"not null;index" json:"timestamp"`

	// CPU metrics
	CPUUsage       float64 `json:"cpuUsage"`       // Percentage (0-100)
	CPULoadAvg1    float64 `json:"cpuLoadAvg1"`    // 1-minute load average
	CPULoadAvg5    float64 `json:"cpuLoadAvg5"`    // 5-minute load average
	CPULoadAvg15   float64 `json:"cpuLoadAvg15"`   // 15-minute load average
	CPUTemperature float64 `json:"cpuTemperature"` // Celsius

	// Memory metrics
	MemoryUsedBytes  uint64  `json:"memoryUsedBytes"`
	MemoryTotalBytes uint64  `json:"memoryTotalBytes"`
	MemoryUsage      float64 `json:"memoryUsage"`      // Percentage (0-100)
	SwapUsedBytes    uint64  `json:"swapUsedBytes"`
	SwapTotalBytes   uint64  `json:"swapTotalBytes"`
	SwapUsage        float64 `json:"swapUsage"`        // Percentage (0-100)

	// Disk metrics (aggregated across all disks)
	DiskUsedBytes      uint64  `json:"diskUsedBytes"`
	DiskTotalBytes     uint64  `json:"diskTotalBytes"`
	DiskUsage          float64 `json:"diskUsage"`          // Percentage (0-100)
	DiskReadBytesPerSec  uint64  `json:"diskReadBytesPerSec"`
	DiskWriteBytesPerSec uint64  `json:"diskWriteBytesPerSec"`
	DiskIOPS           uint64  `json:"diskIOPS"`           // IO operations per second

	// Network metrics (aggregated across all interfaces)
	NetworkRxBytesPerSec uint64 `json:"networkRxBytesPerSec"` // Bytes received per second
	NetworkTxBytesPerSec uint64 `json:"networkTxBytesPerSec"` // Bytes transmitted per second
	NetworkRxPacketsPerSec uint64 `json:"networkRxPacketsPerSec"`
	NetworkTxPacketsPerSec uint64 `json:"networkTxPacketsPerSec"`

	// Process metrics
	ProcessCount  int `json:"processCount"`
	ThreadCount   int `json:"threadCount"`

	CreatedAt time.Time `json:"createdAt"`
}

// HealthScore stores calculated system health scores
type HealthScore struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Timestamp time.Time `gorm:"not null;index" json:"timestamp"`

	// Overall health score (0-100, higher is better)
	Score int `gorm:"not null" json:"score"`

	// Component scores (0-100 each)
	CPUScore     int `json:"cpuScore"`
	MemoryScore  int `json:"memoryScore"`
	DiskScore    int `json:"diskScore"`
	NetworkScore int `json:"networkScore"`

	// Issues detected
	Issues string `gorm:"type:text" json:"issues,omitempty"` // JSON array of issue descriptions

	CreatedAt time.Time `json:"createdAt"`
}

// MetricsTrend represents trend data for a specific metric
type MetricsTrend struct {
	MetricName    string    `json:"metricName"`
	CurrentValue  float64   `json:"currentValue"`
	PreviousValue float64   `json:"previousValue"`
	Change        float64   `json:"change"`        // Absolute change
	ChangePercent float64   `json:"changePercent"` // Percentage change
	Direction     string    `json:"direction"`     // "up", "down", "stable"
	Timestamp     time.Time `json:"timestamp"`
}

// TableName specifies the table name for SystemMetric
func (SystemMetric) TableName() string {
	return "system_metrics"
}

// TableName specifies the table name for HealthScore
func (HealthScore) TableName() string {
	return "health_scores"
}

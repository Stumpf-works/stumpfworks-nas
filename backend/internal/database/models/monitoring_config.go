// Revision: 2025-11-17 | Author: Claude | Version: 1.0.0
package models

import (
	"time"

	"gorm.io/gorm"
)

// MonitoringConfig stores monitoring and observability configuration
type MonitoringConfig struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Prometheus settings
	PrometheusEnabled bool `json:"prometheus_enabled" gorm:"default:true"`

	// Grafana settings
	GrafanaURL string `json:"grafana_url" gorm:"default:'http://localhost:3000'"`

	// Datadog settings
	DatadogEnabled bool   `json:"datadog_enabled" gorm:"default:false"`
	DatadogAPIKey  string `json:"datadog_api_key,omitempty"`
}

// TableName specifies the table name for MonitoringConfig
func (MonitoringConfig) TableName() string {
	return "monitoring_config"
}

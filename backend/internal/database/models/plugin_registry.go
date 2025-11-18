package models

import (
	"time"
)

// PluginRegistry represents a plugin from the registry
type PluginRegistry struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Category    string    `json:"category"`

	// Repository info
	RepositoryURL string `json:"repository_url"`
	DownloadURL   string `json:"download_url"`
	Homepage      string `json:"homepage"`

	// Requirements
	MinNasVersion string   `json:"min_nas_version"`
	RequireDocker bool     `json:"require_docker"`
	RequiredPorts []int    `json:"required_ports" gorm:"serializer:json"`

	// Stats
	Downloads     int       `json:"downloads"`
	Rating        float64   `json:"rating"`
	LastUpdated   time.Time `json:"last_updated"`

	// Installation status (local)
	Installed     bool      `json:"installed" gorm:"-"`
	InstalledVersion string `json:"installed_version,omitempty" gorm:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (PluginRegistry) TableName() string {
	return "plugin_registry"
}

// InstalledPlugin represents a locally installed plugin
type InstalledPlugin struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	Version       string    `json:"version"`
	InstallPath   string    `json:"install_path"`
	Enabled       bool      `json:"enabled"`
	AutoUpdate    bool      `json:"auto_update"`
	InstallDate   time.Time `json:"install_date"`
	LastStarted   time.Time `json:"last_started,omitempty"`

	// Runtime info
	Status        string    `json:"status"` // running, stopped, crashed, updating
	PID           int       `json:"pid,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (InstalledPlugin) TableName() string {
	return "installed_plugins"
}

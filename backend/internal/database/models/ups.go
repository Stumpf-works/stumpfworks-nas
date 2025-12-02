// Revision: 2025-12-02 | Author: Claude | Version: 1.2.0
package models

import (
	"time"

	"gorm.io/gorm"
)

// UPSConfig stores UPS configuration
type UPSConfig struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// UPS Connection
	Enabled        bool   `json:"enabled" gorm:"default:false"`
	UPSName        string `json:"ups_name" gorm:"default:'ups'"`        // NUT UPS name
	UPSHost        string `json:"ups_host" gorm:"default:'localhost'"` // NUT server host
	UPSPort        int    `json:"ups_port" gorm:"default:3493"`        // NUT server port
	UPSUsername    string `json:"ups_username,omitempty"`              // Optional username
	UPSPassword    string `json:"ups_password,omitempty"`              // Optional password
	PollInterval   int    `json:"poll_interval" gorm:"default:30"`     // Poll interval in seconds

	// Shutdown Settings
	LowBatteryShutdown    bool `json:"low_battery_shutdown" gorm:"default:true"`    // Enable automatic shutdown
	LowBatteryThreshold   int  `json:"low_battery_threshold" gorm:"default:20"`    // Battery percentage threshold
	ShutdownDelay         int  `json:"shutdown_delay" gorm:"default:120"`          // Delay in seconds before shutdown
	ShutdownCommand       string `json:"shutdown_command" gorm:"default:'shutdown -h now'"` // Custom shutdown command

	// Notifications
	NotifyOnPowerLoss     bool `json:"notify_on_power_loss" gorm:"default:true"`
	NotifyOnBatteryLow    bool `json:"notify_on_battery_low" gorm:"default:true"`
	NotifyOnPowerRestored bool `json:"notify_on_power_restored" gorm:"default:true"`
}

// TableName specifies the table name for UPSConfig
func (UPSConfig) TableName() string {
	return "ups_config"
}

// UPSEvent stores power event history
type UPSEvent struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	EventType   string `json:"event_type" gorm:"index"` // POWER_LOSS, BATTERY_LOW, POWER_RESTORED, SHUTDOWN_INITIATED, etc.
	Description string `json:"description"`
	BatteryLevel int   `json:"battery_level,omitempty"` // Battery percentage at time of event
	Runtime     int    `json:"runtime,omitempty"`       // Estimated runtime in seconds
	LoadPercent int    `json:"load_percent,omitempty"`  // UPS load percentage
	Voltage     float64 `json:"voltage,omitempty"`      // Input/output voltage
	Severity    string `json:"severity" gorm:"default:'info'"` // info, warning, critical
}

// TableName specifies the table name for UPSEvent
func (UPSEvent) TableName() string {
	return "ups_events"
}

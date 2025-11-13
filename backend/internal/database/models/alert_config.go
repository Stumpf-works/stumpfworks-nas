package models

import "time"

// AlertConfig represents the alerting configuration
type AlertConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Email settings
	Enabled        bool   `gorm:"default:false" json:"enabled"`
	SMTPHost       string `gorm:"size:255" json:"smtpHost"`
	SMTPPort       int    `gorm:"default:587" json:"smtpPort"`
	SMTPUsername   string `gorm:"size:255" json:"smtpUsername"`
	SMTPPassword   string `gorm:"size:255" json:"-"` // Never expose in JSON
	SMTPFromEmail  string `gorm:"size:255" json:"smtpFromEmail"`
	SMTPFromName   string `gorm:"size:255" json:"smtpFromName"`
	SMTPUseTLS     bool   `gorm:"default:true" json:"smtpUseTLS"`
	AlertRecipient string `gorm:"size:255" json:"alertRecipient"`

	// Alert triggers
	OnFailedLogin     bool `gorm:"default:true" json:"onFailedLogin"`
	OnIPBlock         bool `gorm:"default:true" json:"onIPBlock"`
	OnCriticalEvent   bool `gorm:"default:true" json:"onCriticalEvent"`
	FailedLoginThreshold int `gorm:"default:3" json:"failedLoginThreshold"` // Alert after N failed logins

	// Rate limiting for alerts (minutes)
	RateLimitMinutes int `gorm:"default:15" json:"rateLimitMinutes"`
}

// AlertLog represents a sent alert
type AlertLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`

	AlertType string `gorm:"size:100;not null;index" json:"alertType"`
	Subject   string `gorm:"size:255;not null" json:"subject"`
	Body      string `gorm:"type:text" json:"body"`
	Recipient string `gorm:"size:255;not null" json:"recipient"`
	Status    string `gorm:"size:20;not null" json:"status"` // sent, failed
	Error     string `gorm:"type:text" json:"error,omitempty"`
}

// Alert types
const (
	AlertTypeFailedLogin   = "failed_login"
	AlertTypeIPBlock       = "ip_block"
	AlertTypeCriticalEvent = "critical_event"
	AlertTypeSystemError   = "system_error"
)

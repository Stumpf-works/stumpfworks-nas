// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package models

import (
	"time"
)

// AlertRule represents a custom alert rule
type AlertRule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name        string `gorm:"size:255;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Enabled     bool   `gorm:"default:true" json:"enabled"`

	// Rule configuration
	MetricType    string  `gorm:"size:50;not null" json:"metricType"`    // cpu, memory, disk, network, health
	Condition     string  `gorm:"size:20;not null" json:"condition"`     // gt, lt, eq, gte, lte
	Threshold     float64 `gorm:"not null" json:"threshold"`             // The threshold value
	Duration      int     `gorm:"default:0" json:"duration"`             // Duration in seconds (0 = instant)
	CooldownMins  int     `gorm:"default:15" json:"cooldownMins"`        // Cooldown period in minutes

	// Notification settings
	Severity         string `gorm:"size:20;default:'warning'" json:"severity"` // info, warning, critical
	NotifyEmail      bool   `gorm:"default:true" json:"notifyEmail"`
	NotifyWebhook    bool   `gorm:"default:false" json:"notifyWebhook"`
	NotifyChannels   string `gorm:"type:text" json:"notifyChannels"`           // JSON array of additional channels

	// State tracking
	LastTriggered *time.Time `json:"lastTriggered,omitempty"`
	TriggerCount  int        `gorm:"default:0" json:"triggerCount"`
	IsActive      bool       `gorm:"default:false" json:"isActive"` // Currently in alert state
	ActivatedAt   *time.Time `json:"activatedAt,omitempty"`
}

// AlertRuleExecution represents a single execution/trigger of an alert rule
type AlertRuleExecution struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	RuleID       uint    `gorm:"not null;index" json:"ruleId"`
	Rule         AlertRule `gorm:"foreignKey:RuleID" json:"rule,omitempty"`

	MetricValue  float64 `json:"metricValue"`
	Threshold    float64 `json:"threshold"`
	Triggered    bool    `json:"triggered"`
	Acknowledged bool    `gorm:"default:false" json:"acknowledged"`
	AcknowledgedAt   *time.Time `json:"acknowledgedAt,omitempty"`
	AcknowledgedBy   string     `gorm:"size:255" json:"acknowledgedBy,omitempty"`
	AcknowledgeNote  string     `gorm:"type:text" json:"acknowledgeNote,omitempty"`

	NotificationsSent bool   `gorm:"default:false" json:"notificationsSent"`
	Message           string `gorm:"type:text" json:"message"`
}

// Metric types
const (
	MetricTypeCPU     = "cpu"
	MetricTypeMemory  = "memory"
	MetricTypeDisk    = "disk"
	MetricTypeNetwork = "network"
	MetricTypeHealth  = "health"
	MetricTypeTemp    = "temperature"
	MetricTypeIOPS    = "iops"
)

// Condition types
const (
	ConditionGreaterThan        = "gt"
	ConditionLessThan           = "lt"
	ConditionEqual              = "eq"
	ConditionGreaterThanOrEqual = "gte"
	ConditionLessThanOrEqual    = "lte"
)

// Note: Severity levels (SeverityInfo, SeverityWarning, SeverityCritical) are already defined in audit_log.go

// EvaluateCondition checks if the condition is met
func (r *AlertRule) EvaluateCondition(value float64) bool {
	switch r.Condition {
	case ConditionGreaterThan:
		return value > r.Threshold
	case ConditionLessThan:
		return value < r.Threshold
	case ConditionEqual:
		return value == r.Threshold
	case ConditionGreaterThanOrEqual:
		return value >= r.Threshold
	case ConditionLessThanOrEqual:
		return value <= r.Threshold
	default:
		return false
	}
}

// ShouldTrigger checks if the alert should trigger based on cooldown
func (r *AlertRule) ShouldTrigger() bool {
	if !r.Enabled {
		return false
	}

	if r.LastTriggered == nil {
		return true
	}

	cooldownDuration := time.Duration(r.CooldownMins) * time.Minute
	return time.Since(*r.LastTriggered) >= cooldownDuration
}

// TableName specifies the table name for AlertRule
func (AlertRule) TableName() string {
	return "alert_rules"
}

// TableName specifies the table name for AlertRuleExecution
func (AlertRuleExecution) TableName() string {
	return "alert_rule_executions"
}

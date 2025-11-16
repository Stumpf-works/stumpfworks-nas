// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package models

import (
	"time"
)

// FailedLoginAttempt represents a failed login attempt
type FailedLoginAttempt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`

	// Attempt details
	Username  string `gorm:"size:100;not null;index" json:"username"`
	IPAddress string `gorm:"size:45;not null;index" json:"ipAddress"` // IPv6 max length
	UserAgent string `gorm:"size:500" json:"userAgent,omitempty"`

	// Failure reason
	Reason string `gorm:"size:255" json:"reason"` // e.g., "invalid_password", "user_not_found", "account_disabled"

	// Blocking information
	Blocked   bool       `gorm:"default:false;index" json:"blocked"`
	BlockedAt *time.Time `json:"blockedAt,omitempty"`
}

// TableName specifies the table name for FailedLoginAttempt model
func (FailedLoginAttempt) TableName() string {
	return "failed_login_attempts"
}

// IPBlock represents a blocked IP address
type IPBlock struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`
	ExpiresAt time.Time `gorm:"index" json:"expiresAt"`

	// Block details
	IPAddress   string `gorm:"size:45;not null;uniqueIndex" json:"ipAddress"`
	Reason      string `gorm:"size:255" json:"reason"`
	Attempts    int    `gorm:"default:0" json:"attempts"` // Number of failed attempts that triggered the block
	IsActive    bool   `gorm:"default:true;index" json:"isActive"`
	IsPermanent bool   `gorm:"default:false" json:"isPermanent"` // Manual permanent blocks by admin
}

// TableName specifies the table name for IPBlock model
func (IPBlock) TableName() string {
	return "ip_blocks"
}

// IsExpired checks if the IP block has expired
func (b *IPBlock) IsExpired() bool {
	if b.IsPermanent {
		return false
	}
	return time.Now().UTC().After(b.ExpiresAt)
}

// Failure reasons
const (
	FailureReasonInvalidPassword  = "invalid_password"
	FailureReasonUserNotFound     = "user_not_found"
	FailureReasonAccountDisabled  = "account_disabled"
	FailureReasonAccountLocked    = "account_locked"
	FailureReasonTooManyAttempts  = "too_many_attempts"
	FailureReasonIPBlocked        = "ip_blocked"
	FailureReasonInvalidToken     = "invalid_token"
	FailureReasonTokenExpired     = "token_expired"
)

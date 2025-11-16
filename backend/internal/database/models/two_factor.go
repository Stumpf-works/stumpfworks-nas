// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package models

import (
	"time"
)

// TwoFactorAuth stores 2FA configuration for a user
type TwoFactorAuth struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex" json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Enabled   bool      `gorm:"default:false" json:"enabled"`
	Secret    string    `gorm:"size:255;not null" json:"-"` // TOTP secret (encrypted)
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TwoFactorBackupCode stores backup codes for account recovery
type TwoFactorBackupCode struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"userId"`
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	Code      string     `gorm:"size:255;not null" json:"-"` // Hashed backup code
	Used      bool       `gorm:"default:false;index" json:"used"`
	UsedAt    *time.Time `json:"usedAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

// TwoFactorAttempt tracks failed 2FA attempts for rate limiting
type TwoFactorAttempt struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"userId"`
	User       User      `gorm:"foreignKey:UserID" json:"-"`
	IPAddress  string    `gorm:"size:45;not null;index" json:"ipAddress"`
	Success    bool      `gorm:"not null;index" json:"success"`
	AttemptedAt time.Time `gorm:"not null;index" json:"attemptedAt"`
}

// TableName specifies the table name for TwoFactorAuth
func (TwoFactorAuth) TableName() string {
	return "two_factor_auth"
}

// TableName specifies the table name for TwoFactorBackupCode
func (TwoFactorBackupCode) TableName() string {
	return "two_factor_backup_codes"
}

// TableName specifies the table name for TwoFactorAttempt
func (TwoFactorAttempt) TableName() string {
	return "two_factor_attempts"
}

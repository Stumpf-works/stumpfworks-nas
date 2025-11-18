// Revision: 2025-11-18 | Author: Claude | Version: 1.0.0
package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

// PasswordResetToken represents a password reset token for secure password recovery
type PasswordResetToken struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Token     string     `gorm:"uniqueIndex;size:64;not null" json:"token"`
	UserID    uint       `gorm:"not null;index" json:"userId"`
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expiresAt"`
	UsedAt    *time.Time `json:"usedAt,omitempty"`
	IPAddress string     `gorm:"size:45" json:"ipAddress,omitempty"` // IPv4 or IPv6
}

// TableName specifies the table name for PasswordResetToken model
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// GenerateToken creates a cryptographically secure random token
func GenerateToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// IsExpired checks if the token has expired
func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed checks if the token has already been used
func (t *PasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}

// IsValid checks if the token is valid (not expired and not used)
func (t *PasswordResetToken) IsValid() bool {
	return !t.IsExpired() && !t.IsUsed()
}

// MarkAsUsed marks the token as used
func (t *PasswordResetToken) MarkAsUsed(db *gorm.DB) error {
	now := time.Now()
	t.UsedAt = &now
	return db.Model(t).Update("used_at", now).Error
}

// CreatePasswordResetToken creates a new password reset token for a user
func CreatePasswordResetToken(db *gorm.DB, userID uint, validDuration time.Duration) (*PasswordResetToken, error) {
	token, err := GenerateToken()
	if err != nil {
		return nil, err
	}

	resetToken := &PasswordResetToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(validDuration),
	}

	if err := db.Create(resetToken).Error; err != nil {
		return nil, err
	}

	return resetToken, nil
}

// FindValidToken finds a valid (not expired, not used) token
func FindValidToken(db *gorm.DB, token string) (*PasswordResetToken, error) {
	var resetToken PasswordResetToken
	err := db.Preload("User").
		Where("token = ? AND expires_at > ? AND used_at IS NULL", token, time.Now()).
		First(&resetToken).Error

	if err != nil {
		return nil, err
	}

	return &resetToken, nil
}

// CleanupExpiredTokens removes expired and used tokens older than the specified duration
func CleanupExpiredTokens(db *gorm.DB, olderThan time.Duration) error {
	cutoffTime := time.Now().Add(-olderThan)
	return db.Unscoped().
		Where("expires_at < ? OR used_at < ?", time.Now(), cutoffTime).
		Delete(&PasswordResetToken{}).Error
}

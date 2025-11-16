// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a system user
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username     string `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Email        string `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	FullName     string `gorm:"size:255" json:"fullName"`

	Role        string `gorm:"size:50;not null;default:'user'" json:"role"` // admin, user, guest
	IsActive    bool   `gorm:"default:true" json:"isActive"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword verifies if the provided password matches the user's password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// IsAdmin returns true if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// UpdateLastLogin updates the user's last login timestamp
func (u *User) UpdateLastLogin(db *gorm.DB) error {
	now := time.Now()
	u.LastLoginAt = &now
	return db.Model(u).Update("last_login_at", now).Error
}

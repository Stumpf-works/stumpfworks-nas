package core

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/pkg/database"
)

// UserManager manages VPN users and their permissions
type UserManager struct {
	db *gorm.DB
}

// NewUserManager creates a new user manager
func NewUserManager(db *gorm.DB) *UserManager {
	return &UserManager{db: db}
}

// CreateUser creates a new VPN user
func (um *UserManager) CreateUser(username, email, password string) (*database.VPNUser, error) {
	// Check if user already exists
	var existing database.VPNUser
	if err := um.db.Where("username = ?", username).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &database.VPNUser{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Enabled:      true,
	}

	if err := um.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create default protocol permissions (all disabled)
	protocols := []string{"wireguard", "openvpn", "pptp", "l2tp"}
	for _, protocol := range protocols {
		permission := &database.VPNUserProtocol{
			UserID:   user.ID,
			Protocol: protocol,
			Enabled:  false,
		}
		um.db.Create(permission)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (um *UserManager) GetUser(id string) (*database.VPNUser, error) {
	var user database.VPNUser
	if err := um.db.Preload("Protocols").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (um *UserManager) GetUserByUsername(username string) (*database.VPNUser, error) {
	var user database.VPNUser
	if err := um.db.Preload("Protocols").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers retrieves all users
func (um *UserManager) GetAllUsers() ([]database.VPNUser, error) {
	var users []database.VPNUser
	if err := um.db.Preload("Protocols").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser updates a user
func (um *UserManager) UpdateUser(user *database.VPNUser) error {
	return um.db.Save(user).Error
}

// DeleteUser deletes a user and all associated data
func (um *UserManager) DeleteUser(id string) error {
	return um.db.Transaction(func(tx *gorm.DB) error {
		// Delete user (cascade will handle related records)
		if err := tx.Delete(&database.VPNUser{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

// EnableUser enables a user account
func (um *UserManager) EnableUser(id string) error {
	return um.db.Model(&database.VPNUser{}).Where("id = ?", id).Update("enabled", true).Error
}

// DisableUser disables a user account
func (um *UserManager) DisableUser(id string) error {
	return um.db.Model(&database.VPNUser{}).Where("id = ?", id).Update("enabled", false).Error
}

// UpdatePassword updates a user's password
func (um *UserManager) UpdatePassword(id, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return um.db.Model(&database.VPNUser{}).
		Where("id = ?", id).
		Update("password_hash", string(hashedPassword)).Error
}

// VerifyPassword verifies a user's password
func (um *UserManager) VerifyPassword(username, password string) (*database.VPNUser, error) {
	user, err := um.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	if !user.Enabled {
		return nil, fmt.Errorf("user account is disabled")
	}

	return user, nil
}

// UpdateProtocolAccess updates protocol access permissions for a user
func (um *UserManager) UpdateProtocolAccess(userID string, protocols map[string]bool) error {
	return um.db.Transaction(func(tx *gorm.DB) error {
		for protocol, enabled := range protocols {
			permission := &database.VPNUserProtocol{
				UserID:   userID,
				Protocol: protocol,
				Enabled:  enabled,
			}

			// Upsert
			if err := tx.Where("user_id = ? AND protocol = ?", userID, protocol).
				Assign(database.VPNUserProtocol{Enabled: enabled}).
				FirstOrCreate(permission).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetProtocolAccess retrieves protocol access permissions for a user
func (um *UserManager) GetProtocolAccess(userID string) (map[string]bool, error) {
	var protocols []database.VPNUserProtocol
	if err := um.db.Where("user_id = ?", userID).Find(&protocols).Error; err != nil {
		return nil, err
	}

	access := make(map[string]bool)
	for _, p := range protocols {
		access[p.Protocol] = p.Enabled
	}

	return access, nil
}

// HasProtocolAccess checks if a user has access to a specific protocol
func (um *UserManager) HasProtocolAccess(userID, protocol string) (bool, error) {
	var permission database.VPNUserProtocol
	if err := um.db.Where("user_id = ? AND protocol = ?", userID, protocol).First(&permission).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return permission.Enabled, nil
}

// UpdateLastConnection updates the last connection time for a user
func (um *UserManager) UpdateLastConnection(userID string) error {
	now := time.Now()
	return um.db.Model(&database.VPNUser{}).
		Where("id = ?", userID).
		Update("last_connection", now).Error
}

// GetUserStats retrieves statistics about users
func (um *UserManager) GetUserStats() (*UserStats, error) {
	var totalUsers int64
	var activeUsers int64

	um.db.Model(&database.VPNUser{}).Count(&totalUsers)
	um.db.Model(&database.VPNUser{}).Where("enabled = ?", true).Count(&activeUsers)

	// Count users with recent connections (last 24 hours)
	var recentlyConnected int64
	yesterday := time.Now().Add(-24 * time.Hour)
	um.db.Model(&database.VPNUser{}).
		Where("last_connection > ?", yesterday).
		Count(&recentlyConnected)

	return &UserStats{
		TotalUsers:        int(totalUsers),
		ActiveUsers:       int(activeUsers),
		RecentlyConnected: int(recentlyConnected),
	}, nil
}

// SearchUsers searches for users by username or email
func (um *UserManager) SearchUsers(query string) ([]database.VPNUser, error) {
	var users []database.VPNUser
	searchPattern := fmt.Sprintf("%%%s%%", query)

	if err := um.db.Preload("Protocols").
		Where("username LIKE ? OR email LIKE ?", searchPattern, searchPattern).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// UserStats represents user statistics
type UserStats struct {
	TotalUsers        int `json:"totalUsers"`
	ActiveUsers       int `json:"activeUsers"`
	RecentlyConnected int `json:"recentlyConnected"`
}

// ProtocolAccessRequest represents a protocol access update request
type ProtocolAccessRequest struct {
	WireGuard bool `json:"wireguard"`
	OpenVPN   bool `json:"openvpn"`
	PPTP      bool `json:"pptp"`
	L2TP      bool `json:"l2tp"`
}

// ToMap converts ProtocolAccessRequest to map
func (r *ProtocolAccessRequest) ToMap() map[string]bool {
	return map[string]bool{
		"wireguard": r.WireGuard,
		"openvpn":   r.OpenVPN,
		"pptp":      r.PPTP,
		"l2tp":      r.L2TP,
	}
}

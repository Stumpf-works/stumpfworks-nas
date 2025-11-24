// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package users

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/sysutil"
	"go.uber.org/zap"
)

// SambaUserManager handles synchronization between web users and Samba users
type SambaUserManager struct {
	enabled bool
}

var sambaManager *SambaUserManager

// InitSambaUserManager initializes the Samba user manager
func InitSambaUserManager() *SambaUserManager {
	manager := &SambaUserManager{
		enabled: isSambaInstalled(),
	}

	if manager.enabled {
		logger.Info("Samba user synchronization enabled")
	} else {
		logger.Warn("Samba not installed - users will only work for web access",
			zap.String("note", "Install Samba to enable Windows network drives: apt install samba"))
	}

	sambaManager = manager
	return manager
}

// GetSambaManager returns the global Samba manager instance
func GetSambaManager() *SambaUserManager {
	if sambaManager == nil {
		return InitSambaUserManager()
	}
	return sambaManager
}

// isSambaInstalled checks if Samba is installed on the system
func isSambaInstalled() bool {
	// Check for smbpasswd command
	if !sysutil.CommandExists("smbpasswd") {
		return false
	}

	// Check for pdbedit command
	if !sysutil.CommandExists("pdbedit") {
		return false
	}

	return true
}

// CreateSambaUser creates a Samba user synchronized with a web user
// This allows the user to access SMB shares from Windows/Mac/Linux clients
func (m *SambaUserManager) CreateSambaUser(username, password string) error {
	if !m.enabled {
		logger.Debug("Samba not enabled, skipping user creation", zap.String("username", username))
		return nil // Not an error - just not available
	}

	logger.Info("Creating Samba user", zap.String("username", username))

	// Step 1: Create Linux system user (required for Samba)
	// We create a "system" user without home directory and no shell access
	if err := m.createLinuxUser(username); err != nil {
		return fmt.Errorf("failed to create Linux user: %w", err)
	}

	// Small delay to ensure /etc/passwd locks are fully released
	// This prevents lock contention between useradd and smbpasswd
	time.Sleep(50 * time.Millisecond)
	logger.Info("Linux user created, proceeding to Samba password setup", zap.String("username", username))

	// Step 2: Add user to Samba database with password
	if err := m.addSambaPassword(username, password); err != nil {
		// Cleanup: remove Linux user if Samba setup failed
		m.deleteLinuxUser(username)
		return fmt.Errorf("failed to add Samba password: %w", err)
	}

	// Step 3: Enable the Samba user
	if err := m.enableSambaUser(username); err != nil {
		logger.Warn("Failed to enable Samba user (user created but disabled)",
			zap.String("username", username),
			zap.Error(err))
		// Don't fail - user is created, just disabled
	}

	logger.Info("Samba user created successfully", zap.String("username", username))
	return nil
}

// UpdateSambaPassword updates the password for an existing Samba user
func (m *SambaUserManager) UpdateSambaPassword(username, newPassword string) error {
	if !m.enabled {
		logger.Debug("Samba not enabled, skipping password update", zap.String("username", username))
		return nil
	}

	logger.Info("Updating Samba password", zap.String("username", username))

	// Check if user exists in Samba
	exists, err := m.sambaUserExists(username)
	if err != nil {
		return fmt.Errorf("failed to check Samba user existence: %w", err)
	}

	if !exists {
		// User doesn't exist in Samba yet - create it
		logger.Info("Samba user doesn't exist, creating", zap.String("username", username))
		return m.CreateSambaUser(username, newPassword)
	}

	// Update password using smbpasswd
	if err := m.addSambaPassword(username, newPassword); err != nil {
		return fmt.Errorf("failed to update Samba password: %w", err)
	}

	logger.Info("Samba password updated successfully", zap.String("username", username))
	return nil
}

// DeleteSambaUser removes a Samba user
func (m *SambaUserManager) DeleteSambaUser(username string) error {
	if !m.enabled {
		logger.Debug("Samba not enabled, skipping user deletion", zap.String("username", username))
		return nil
	}

	logger.Info("Deleting Samba user", zap.String("username", username))

	// Step 1: Remove from Samba database
	if err := m.removeSambaUser(username); err != nil {
		logger.Warn("Failed to remove Samba user", zap.String("username", username), zap.Error(err))
		// Continue anyway - try to remove Linux user
	}

	// Step 2: Remove Linux system user
	if err := m.deleteLinuxUser(username); err != nil {
		logger.Warn("Failed to remove Linux user", zap.String("username", username), zap.Error(err))
		// Not critical - user is gone from Samba
	}

	logger.Info("Samba user deleted successfully", zap.String("username", username))
	return nil
}

// createLinuxUser creates a Linux system user for Samba
func (m *SambaUserManager) createLinuxUser(username string) error {
	// Check if user already exists
	cmd := exec.Command(sysutil.FindCommand("id"), username)
	if err := cmd.Run(); err == nil {
		logger.Debug("Linux user already exists", zap.String("username", username))
		return nil // User exists, that's fine
	}

	// Create user without home directory (-M) and with no shell access (-s /bin/false)
	// This is a "system user" only for Samba authentication
	useraddPath := sysutil.FindCommand("useradd")

	// Retry logic for /etc/passwd lock contention
	maxRetries := 5
	baseDelay := 100 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 100ms, 200ms, 400ms, 800ms, 1600ms
			delay := baseDelay * time.Duration(1<<uint(attempt-1))
			logger.Info("Retrying useradd after delay",
				zap.String("username", username),
				zap.Int("attempt", attempt+1),
				zap.Duration("delay", delay))
			time.Sleep(delay)
		}

		cmd = exec.Command(useraddPath,
			"-M",                  // No home directory
			"-s", "/bin/false",    // No shell access (security)
			"-c", "Stumpf.Works NAS User", // Comment
			username)

		output, err := cmd.CombinedOutput()
		if err == nil {
			logger.Info("Linux user created successfully",
				zap.String("username", username),
				zap.String("useradd_path", useraddPath),
				zap.Int("attempts", attempt+1))
			return nil
		}

		// Check if error is due to /etc/passwd lock contention
		outputStr := string(output)
		isLockError := strings.Contains(outputStr, "konnte nicht gesperrt werden") ||
			strings.Contains(outputStr, "cannot lock") ||
			strings.Contains(outputStr, "unable to lock") ||
			strings.Contains(outputStr, "temporarily unavailable")

		// If it's not a lock error, or we're on the last attempt, return the error
		if !isLockError || attempt == maxRetries-1 {
			return fmt.Errorf("useradd failed: %s: %w", outputStr, err)
		}

		logger.Info("useradd lock contention detected, will retry",
			zap.String("username", username),
			zap.Int("attempt", attempt+1),
			zap.Int("max_retries", maxRetries),
			zap.String("error", outputStr))
	}

	return fmt.Errorf("useradd failed after %d attempts", maxRetries)
}

// deleteLinuxUser removes a Linux system user
func (m *SambaUserManager) deleteLinuxUser(username string) error {
	// Check if user exists
	cmd := exec.Command(sysutil.FindCommand("id"), username)
	if err := cmd.Run(); err != nil {
		logger.Debug("Linux user doesn't exist", zap.String("username", username))
		return nil // Already gone
	}

	// Delete user (but keep home directory if any - just in case)
	cmd = exec.Command(sysutil.FindCommand("userdel"), username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("userdel failed: %s: %w", string(output), err)
	}

	logger.Debug("Linux user deleted", zap.String("username", username))
	return nil
}

// addSambaPassword adds or updates a password for a Samba user
func (m *SambaUserManager) addSambaPassword(username, password string) error {
	smbpasswdPath := sysutil.FindCommand("smbpasswd")

	// Retry logic for /etc/passwd lock contention
	// smbpasswd needs to read /etc/passwd to get user UID
	maxRetries := 5
	baseDelay := 100 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 100ms, 200ms, 400ms, 800ms, 1600ms
			delay := baseDelay * time.Duration(1<<uint(attempt-1))
			logger.Info("Retrying smbpasswd after delay",
				zap.String("username", username),
				zap.Int("attempt", attempt+1),
				zap.Duration("delay", delay))
			time.Sleep(delay)
		}

		// Use smbpasswd to set password
		// -a = add user (or update if exists)
		// -s = silent mode (read password from stdin)
		cmd := exec.Command(smbpasswdPath, "-a", "-s", username)

		// Pass password via stdin (format: password\npassword\n)
		cmd.Stdin = strings.NewReader(password + "\n" + password + "\n")

		output, err := cmd.CombinedOutput()
		if err == nil {
			logger.Info("Samba password set successfully",
				zap.String("username", username),
				zap.String("smbpasswd_path", smbpasswdPath),
				zap.Int("attempts", attempt+1))
			return nil
		}

		// Check if error is due to /etc/passwd lock contention
		outputStr := string(output)
		isLockError := strings.Contains(outputStr, "konnte nicht gesperrt werden") ||
			strings.Contains(outputStr, "cannot lock") ||
			strings.Contains(outputStr, "unable to lock") ||
			strings.Contains(outputStr, "temporarily unavailable") ||
			strings.Contains(outputStr, "passwd") && strings.Contains(outputStr, "lock")

		// If it's not a lock error, or we're on the last attempt, return the error
		if !isLockError || attempt == maxRetries-1 {
			return fmt.Errorf("smbpasswd failed: %s: %w", outputStr, err)
		}

		logger.Info("smbpasswd lock contention detected, will retry",
			zap.String("username", username),
			zap.Int("attempt", attempt+1),
			zap.Int("max_retries", maxRetries),
			zap.String("error", outputStr))
	}

	return fmt.Errorf("smbpasswd failed after %d attempts", maxRetries)
}

// enableSambaUser enables a Samba user account
func (m *SambaUserManager) enableSambaUser(username string) error {
	// Use smbpasswd -e to enable
	cmd := exec.Command(sysutil.FindCommand("smbpasswd"), "-e", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("smbpasswd enable failed: %s: %w", string(output), err)
	}

	logger.Debug("Samba user enabled", zap.String("username", username))
	return nil
}

// removeSambaUser removes a user from Samba database
func (m *SambaUserManager) removeSambaUser(username string) error {
	// Use smbpasswd -x to remove
	cmd := exec.Command(sysutil.FindCommand("smbpasswd"), "-x", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if error is just "user doesn't exist"
		if strings.Contains(string(output), "Failed to find entry") {
			logger.Debug("Samba user doesn't exist", zap.String("username", username))
			return nil // Not an error
		}
		return fmt.Errorf("smbpasswd remove failed: %s: %w", string(output), err)
	}

	logger.Debug("Samba user removed", zap.String("username", username))
	return nil
}

// sambaUserExists checks if a user exists in Samba database
func (m *SambaUserManager) sambaUserExists(username string) (bool, error) {
	// Use pdbedit to list users and check if username exists
	cmd := exec.Command(sysutil.FindCommand("pdbedit"), "-L", "-u", username)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// pdbedit returns various messages when user doesn't exist
		outputStr := string(output)
		if strings.Contains(outputStr, "Failed to find entry") ||
		   strings.Contains(outputStr, "Username not found") ||
		   strings.Contains(outputStr, "user not found") {
			logger.Debug("Samba user does not exist", zap.String("username", username))
			return false, nil
		}
		return false, fmt.Errorf("pdbedit failed: %s: %w", outputStr, err)
	}

	// If output contains username, user exists
	return strings.Contains(string(output), username), nil
}

// ListSambaUsers returns all Samba users (for debugging/admin)
func (m *SambaUserManager) ListSambaUsers() ([]string, error) {
	if !m.enabled {
		return []string{}, nil
	}

	cmd := exec.Command(sysutil.FindCommand("pdbedit"), "-L")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("pdbedit failed: %s: %w", string(output), err)
	}

	// Parse output (format: "username:uid:...")
	lines := strings.Split(string(output), "\n")
	users := []string{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) > 0 {
			users = append(users, parts[0])
		}
	}

	return users, nil
}

// IsEnabled returns whether Samba sync is enabled
func (m *SambaUserManager) IsEnabled() bool {
	return m.enabled
}

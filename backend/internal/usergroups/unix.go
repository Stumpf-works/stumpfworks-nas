package usergroups

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/sysutil"
	"go.uber.org/zap"
)

// UnixGroupManager handles synchronization between database user groups and Unix system groups
type UnixGroupManager struct {
	enabled bool
}

var unixGroupManager *UnixGroupManager

// InitUnixGroupManager initializes the Unix group manager
func InitUnixGroupManager() *UnixGroupManager {
	if unixGroupManager != nil {
		return unixGroupManager
	}

	manager := &UnixGroupManager{
		enabled: checkUnixGroupCommandsAvailable(),
	}

	unixGroupManager = manager
	return manager
}

// GetUnixGroupManager returns the singleton Unix group manager instance
func GetUnixGroupManager() *UnixGroupManager {
	if unixGroupManager == nil {
		return InitUnixGroupManager()
	}
	return unixGroupManager
}

// checkUnixGroupCommandsAvailable checks if required Unix group commands are available
func checkUnixGroupCommandsAvailable() bool {
	requiredCommands := []string{"groupadd", "groupdel", "usermod", "getent"}
	for _, cmd := range requiredCommands {
		if sysutil.FindCommand(cmd) == "" {
			logger.Warn("Unix group command not found - group sync disabled",
				zap.String("command", cmd))
			return false
		}
	}
	return true
}

// IsEnabled returns whether the Unix group manager is enabled
func (m *UnixGroupManager) IsEnabled() bool {
	return m.enabled
}

// CreateUnixGroup creates a Unix system group for a database user group
func (m *UnixGroupManager) CreateUnixGroup(group *models.UserGroup) error {
	if !m.enabled {
		return fmt.Errorf("unix group manager not enabled")
	}

	groupName := group.UnixGroupName()

	// Check if group already exists
	if m.groupExists(groupName) {
		logger.Debug("Unix group already exists", zap.String("group", groupName))
		return nil
	}

	// Create the group
	groupaddPath := sysutil.FindCommand("groupadd")
	cmd := exec.Command(groupaddPath, groupName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create Unix group %s: %s: %w", groupName, string(output), err)
	}

	logger.Info("Created Unix group", zap.String("group", groupName))
	return nil
}

// DeleteUnixGroup deletes a Unix system group
func (m *UnixGroupManager) DeleteUnixGroup(group *models.UserGroup) error {
	if !m.enabled {
		return fmt.Errorf("unix group manager not enabled")
	}

	groupName := group.UnixGroupName()

	// Check if group exists
	if !m.groupExists(groupName) {
		logger.Debug("Unix group doesn't exist, nothing to delete", zap.String("group", groupName))
		return nil
	}

	// Delete the group
	groupdelPath := sysutil.FindCommand("groupdel")
	cmd := exec.Command(groupdelPath, groupName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to delete Unix group %s: %s: %w", groupName, string(output), err)
	}

	logger.Info("Deleted Unix group", zap.String("group", groupName))
	return nil
}

// AddUserToUnixGroup adds a user to a Unix system group
func (m *UnixGroupManager) AddUserToUnixGroup(username string, group *models.UserGroup) error {
	if !m.enabled {
		return fmt.Errorf("unix group manager not enabled")
	}

	groupName := group.UnixGroupName()

	// Ensure the group exists first
	if !m.groupExists(groupName) {
		if err := m.CreateUnixGroup(group); err != nil {
			return fmt.Errorf("failed to create Unix group before adding user: %w", err)
		}
	}

	// Check if user is already in the group
	if m.userInGroup(username, groupName) {
		logger.Debug("User already in Unix group",
			zap.String("user", username),
			zap.String("group", groupName))
		return nil
	}

	// Add user to group
	usermodPath := sysutil.FindCommand("usermod")
	cmd := exec.Command(usermodPath, "-aG", groupName, username)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add user %s to Unix group %s: %s: %w",
			username, groupName, string(output), err)
	}

	logger.Info("Added user to Unix group",
		zap.String("user", username),
		zap.String("group", groupName))
	return nil
}

// RemoveUserFromUnixGroup removes a user from a Unix system group
func (m *UnixGroupManager) RemoveUserFromUnixGroup(username string, group *models.UserGroup) error {
	if !m.enabled {
		return fmt.Errorf("unix group manager not enabled")
	}

	groupName := group.UnixGroupName()

	// Check if user is in the group
	if !m.userInGroup(username, groupName) {
		logger.Debug("User not in Unix group, nothing to remove",
			zap.String("user", username),
			zap.String("group", groupName))
		return nil
	}

	// Remove user from group using gpasswd -d
	gpasswdPath := sysutil.FindCommand("gpasswd")
	if gpasswdPath == "" {
		// Fallback: use deluser if available
		deluserPath := sysutil.FindCommand("deluser")
		if deluserPath == "" {
			return fmt.Errorf("neither gpasswd nor deluser found - cannot remove user from group")
		}
		cmd := exec.Command(deluserPath, username, groupName)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to remove user %s from Unix group %s: %s: %w",
				username, groupName, string(output), err)
		}
	} else {
		cmd := exec.Command(gpasswdPath, "-d", username, groupName)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to remove user %s from Unix group %s: %s: %w",
				username, groupName, string(output), err)
		}
	}

	logger.Info("Removed user from Unix group",
		zap.String("user", username),
		zap.String("group", groupName))
	return nil
}

// groupExists checks if a Unix group exists
func (m *UnixGroupManager) groupExists(groupName string) bool {
	getentPath := sysutil.FindCommand("getent")
	cmd := exec.Command(getentPath, "group", groupName)
	err := cmd.Run()
	return err == nil
}

// userInGroup checks if a user is a member of a Unix group
func (m *UnixGroupManager) userInGroup(username, groupName string) bool {
	idPath := sysutil.FindCommand("id")
	cmd := exec.Command(idPath, "-nG", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	groups := strings.Fields(string(output))
	for _, group := range groups {
		if group == groupName {
			return true
		}
	}
	return false
}

// SyncGroupMembers synchronizes all members of a group to the Unix system
func (m *UnixGroupManager) SyncGroupMembers(group *models.UserGroup) error {
	if !m.enabled {
		return fmt.Errorf("unix group manager not enabled")
	}

	// Ensure the Unix group exists
	if !m.groupExists(group.UnixGroupName()) {
		if err := m.CreateUnixGroup(group); err != nil {
			return fmt.Errorf("failed to create Unix group: %w", err)
		}
	}

	// Add all members to the Unix group
	for _, member := range group.Members {
		if err := m.AddUserToUnixGroup(member.Username, group); err != nil {
			logger.Warn("Failed to add member to Unix group",
				zap.String("user", member.Username),
				zap.String("group", group.UnixGroupName()),
				zap.Error(err))
			// Continue with other members even if one fails
		}
	}

	return nil
}

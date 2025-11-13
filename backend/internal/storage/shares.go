package storage

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// findSmbdPath searches for smbd binary in common locations
func findSmbdPath() (string, error) {
	// Try exec.LookPath first (checks PATH)
	if path, err := exec.LookPath("smbd"); err == nil {
		return path, nil
	}

	// Check common installation paths
	commonPaths := []string{
		"/usr/sbin/smbd",
		"/usr/bin/smbd",
		"/sbin/smbd",
		"/usr/local/sbin/smbd",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("smbd not found in PATH or common locations")
}

// findExportfsPath searches for exportfs binary in common locations
func findExportfsPath() (string, error) {
	// Try exec.LookPath first (checks PATH)
	if path, err := exec.LookPath("exportfs"); err == nil {
		return path, nil
	}

	// Check common installation paths
	commonPaths := []string{
		"/usr/sbin/exportfs",
		"/usr/bin/exportfs",
		"/sbin/exportfs",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("exportfs not found in PATH or common locations")
}

// toShare converts models.Share to Share
func toShare(s *models.Share) *Share {
	var validUsers []string
	if s.ValidUsers != "" {
		validUsers = strings.Split(s.ValidUsers, ",")
	}

	return &Share{
		ID:          fmt.Sprintf("%d", s.ID),
		Name:        s.Name,
		Path:        s.Path,
		Type:        ShareType(s.Type),
		Description: s.Description,
		Enabled:     s.Enabled,
		ReadOnly:    s.ReadOnly,
		Browseable:  s.Browseable,
		GuestOK:     s.GuestOK,
		ValidUsers:  validUsers,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

// ListShares lists all network shares
func ListShares() ([]Share, error) {
	var models []models.Share
	if err := database.DB.Find(&models).Error; err != nil {
		return nil, err
	}

	shares := make([]Share, len(models))
	for i, model := range models {
		shares[i] = *toShare(&model)
	}

	return shares, nil
}

// GetShare retrieves a specific share by ID
func GetShare(id string) (*Share, error) {
	var model models.Share
	if err := database.DB.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("share not found")
		}
		return nil, err
	}

	return toShare(&model), nil
}

// CreateShare creates a new network share
func CreateShare(req *CreateShareRequest) (*Share, error) {
	logger.Info("Creating share",
		zap.String("name", req.Name),
		zap.String("type", string(req.Type)),
		zap.String("path", req.Path))

	// Validate path exists
	if _, err := os.Stat(req.Path); os.IsNotExist(err) {
		return nil, fmt.Errorf("path does not exist: %s", req.Path)
	}

	// Create database record
	model := &models.Share{
		Name:        req.Name,
		Path:        req.Path,
		Type:        string(req.Type),
		Description: req.Description,
		Enabled:     true,
		ReadOnly:    req.ReadOnly,
		Browseable:  req.Browseable,
		GuestOK:     req.GuestOK,
		ValidUsers:  strings.Join(req.ValidUsers, ","),
	}

	// Check if share with this name already exists
	var existingShare models.Share
	if err := database.DB.Where("name = ?", req.Name).First(&existingShare).Error; err == nil {
		return nil, fmt.Errorf("a share with the name '%s' already exists", req.Name)
	}

	if err := database.DB.Create(model).Error; err != nil {
		// Check if it's a duplicate key error (in case of race condition)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") ||
		   strings.Contains(err.Error(), "duplicate key") {
			return nil, fmt.Errorf("a share with the name '%s' already exists", req.Name)
		}
		return nil, fmt.Errorf("failed to create share in database: %w", err)
	}

	// Configure the share based on type
	switch req.Type {
	case ShareTypeSMB:
		if err := configureSMBShare(model); err != nil {
			database.DB.Delete(model)
			return nil, fmt.Errorf("failed to configure SMB share: %w", err)
		}
	case ShareTypeNFS:
		if err := configureNFSShare(model); err != nil {
			database.DB.Delete(model)
			return nil, fmt.Errorf("failed to configure NFS share: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported share type: %s", req.Type)
	}

	logger.Info("Share created successfully", zap.String("name", req.Name))

	return toShare(model), nil
}

// UpdateShare updates an existing share
func UpdateShare(id string, req *CreateShareRequest) (*Share, error) {
	var model models.Share
	if err := database.DB.First(&model, id).Error; err != nil {
		return nil, err
	}

	// Update fields
	model.Name = req.Name
	model.Path = req.Path
	model.Description = req.Description
	model.ReadOnly = req.ReadOnly
	model.Browseable = req.Browseable
	model.GuestOK = req.GuestOK
	model.ValidUsers = strings.Join(req.ValidUsers, ",")

	if err := database.DB.Save(&model).Error; err != nil {
		return nil, err
	}

	// Reconfigure the share
	switch ShareType(model.Type) {
	case ShareTypeSMB:
		if err := configureSMBShare(&model); err != nil {
			return nil, err
		}
	case ShareTypeNFS:
		if err := configureNFSShare(&model); err != nil {
			return nil, err
		}
	}

	return toShare(&model), nil
}

// DeleteShare deletes a network share
func DeleteShare(id string) error {
	var model models.Share
	if err := database.DB.First(&model, id).Error; err != nil {
		return err
	}

	// Remove configuration
	switch ShareType(model.Type) {
	case ShareTypeSMB:
		if err := removeSMBShare(&model); err != nil {
			return err
		}
	case ShareTypeNFS:
		if err := removeNFSShare(&model); err != nil {
			return err
		}
	}

	// Delete from database
	if err := database.DB.Delete(&model).Error; err != nil {
		return err
	}

	logger.Info("Share deleted successfully", zap.String("name", model.Name))

	return nil
}

// configureSMBShare configures a Samba share
func configureSMBShare(share *models.Share) error {
	// Check if Samba is installed
	smbdPath, err := findSmbdPath()
	if err != nil {
		logger.Warn("Samba not installed - share created but network access disabled",
			zap.String("share", share.Name),
			zap.String("note", "Install Samba to enable network access: apt install samba"),
			zap.Error(err))
		return nil // Don't fail - share will work locally for File Manager
	}

	logger.Info("Found Samba", zap.String("path", smbdPath))

	// Set up share permissions (group, ownership, etc.)
	if err := setupSharePermissions(share); err != nil {
		logger.Warn("Failed to set share permissions",
			zap.String("share", share.Name),
			zap.Error(err))
		// Don't fail - share config can still be written
	}

	// Build Samba configuration
	config := fmt.Sprintf(`
[%s]
   path = %s
   comment = %s
   browseable = %s
   read only = %s
   guest ok = %s
`, share.Name, share.Path, share.Description,
		boolToYesNo(share.Browseable),
		boolToYesNo(share.ReadOnly),
		boolToYesNo(share.GuestOK))

	if share.ValidUsers != "" {
		config += fmt.Sprintf("   valid users = %s\n", strings.ReplaceAll(share.ValidUsers, ",", " "))
	}

	// Write to Samba config directory
	configPath := filepath.Join("/etc/samba/shares.d", share.Name+".conf")
	if err := os.MkdirAll("/etc/samba/shares.d", 0755); err != nil {
		return err
	}

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return err
	}

	// Ensure main smb.conf includes shares.d
	ensureSambaInclude()

	// Reload Samba
	cmd := exec.Command("systemctl", "reload", "smbd")
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Warn("Failed to reload Samba", zap.String("output", string(output)))
	}

	return nil
}

// removeSMBShare removes a Samba share configuration
func removeSMBShare(share *models.Share) error {
	configPath := filepath.Join("/etc/samba/shares.d", share.Name+".conf")
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Reload Samba
	cmd := exec.Command("systemctl", "reload", "smbd")
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Warn("Failed to reload Samba", zap.String("output", string(output)))
	}

	return nil
}

// configureNFSShare configures an NFS export
func configureNFSShare(share *models.Share) error {
	// Check if NFS is installed
	exportfsPath, err := findExportfsPath()
	if err != nil {
		logger.Warn("NFS not installed - share created but network access disabled",
			zap.String("share", share.Name),
			zap.String("note", "Install NFS to enable network access: apt install nfs-kernel-server"),
			zap.Error(err))
		return nil // Don't fail - share will work locally for File Manager
	}

	logger.Info("Found NFS", zap.String("path", exportfsPath))

	// Build export entry
	export := fmt.Sprintf("%s *(rw,sync,no_subtree_check)\n", share.Path)
	if share.ReadOnly {
		export = fmt.Sprintf("%s *(ro,sync,no_subtree_check)\n", share.Path)
	}

	// Append to /etc/exports
	file, err := os.OpenFile("/etc/exports", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(export); err != nil {
		return err
	}

	// Reload NFS exports
	cmd := exec.Command("exportfs", "-ra")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload exports: %s: %w", string(output), err)
	}

	return nil
}

// removeNFSShare removes an NFS export
func removeNFSShare(share *models.Share) error {
	// This is a simplified version
	// In production, you'd want to parse and rewrite /etc/exports properly

	// Unexport
	cmd := exec.Command("exportfs", "-u", "*:"+share.Path)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Warn("Failed to unexport", zap.String("output", string(output)))
	}

	return nil
}

// ensureSambaInclude ensures the main smb.conf includes shares.d directory
func ensureSambaInclude() error {
	smbConfPath := "/etc/samba/smb.conf"
	includeDirective := "include = /etc/samba/shares.d/*.conf"

	// Read current config
	data, err := os.ReadFile(smbConfPath)
	if err != nil {
		return err
	}

	content := string(data)

	// Check if include directive already exists
	if strings.Contains(content, includeDirective) {
		return nil
	}

	// Append include directive to [global] section
	lines := strings.Split(content, "\n")
	var newLines []string
	addedInclude := false

	for _, line := range lines {
		newLines = append(newLines, line)

		// Add include after [global] section
		if strings.Contains(line, "[global]") && !addedInclude {
			newLines = append(newLines, "   "+includeDirective)
			addedInclude = true
		}
	}

	// Write back
	return os.WriteFile(smbConfPath, []byte(strings.Join(newLines, "\n")), 0644)
}

// boolToYesNo converts a boolean to yes/no string for Samba config
func boolToYesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// EnableShare enables a share
func EnableShare(id string) error {
	return updateShareStatus(id, true)
}

// DisableShare disables a share
func DisableShare(id string) error {
	return updateShareStatus(id, false)
}

// updateShareStatus updates the enabled status of a share
func updateShareStatus(id string, enabled bool) error {
	var model models.Share
	if err := database.DB.First(&model, id).Error; err != nil {
		return err
	}

	model.Enabled = enabled
	if err := database.DB.Save(&model).Error; err != nil {
		return err
	}

	// If disabling, remove the configuration
	if !enabled {
		switch ShareType(model.Type) {
		case ShareTypeSMB:
			return removeSMBShare(&model)
		case ShareTypeNFS:
			return removeNFSShare(&model)
		}
	} else {
		// If enabling, reconfigure
		switch ShareType(model.Type) {
		case ShareTypeSMB:
			return configureSMBShare(&model)
		case ShareTypeNFS:
			return configureNFSShare(&model)
		}
	}

	return nil
}

// setupSharePermissions sets up proper permissions for a share directory
// Creates smbusers group, sets group ownership, and configures permissions
func setupSharePermissions(share *models.Share) error {
	const smbGroup = "smbusers"

	// Ensure the smbusers group exists
	if err := ensureSMBGroup(smbGroup); err != nil {
		return fmt.Errorf("failed to ensure SMB group: %w", err)
	}

	// Set group ownership on the share path
	if err := setShareGroupOwnership(share.Path, smbGroup); err != nil {
		return fmt.Errorf("failed to set group ownership: %w", err)
	}

	// Set permissions (775 = rwxrwxr-x)
	if err := os.Chmod(share.Path, 0775); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Add valid users to the smbusers group
	if share.ValidUsers != "" {
		users := strings.Split(share.ValidUsers, ",")
		for _, username := range users {
			username = strings.TrimSpace(username)
			if username == "" {
				continue
			}
			if err := addUserToGroup(username, smbGroup); err != nil {
				logger.Warn("Failed to add user to SMB group",
					zap.String("user", username),
					zap.String("group", smbGroup),
					zap.Error(err))
				// Don't fail the whole operation if one user fails
			}
		}
	}

	logger.Info("Share permissions configured",
		zap.String("share", share.Name),
		zap.String("path", share.Path),
		zap.String("group", smbGroup),
		zap.String("permissions", "775"))

	return nil
}

// ensureSMBGroup ensures the smbusers group exists, creates it if not
func ensureSMBGroup(groupName string) error {
	// Check if group exists
	cmd := exec.Command("getent", "group", groupName)
	if err := cmd.Run(); err == nil {
		// Group exists
		return nil
	}

	// Group doesn't exist, create it
	cmd = exec.Command("groupadd", groupName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create group %s: %s: %w", groupName, string(output), err)
	}

	logger.Info("Created SMB group", zap.String("group", groupName))
	return nil
}

// setShareGroupOwnership sets the group ownership of a path
func setShareGroupOwnership(path, groupName string) error {
	// Use chgrp to set group ownership
	cmd := exec.Command("chgrp", groupName, path)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("chgrp failed: %s: %w", string(output), err)
	}
	return nil
}

// addUserToGroup adds a user to a group
func addUserToGroup(username, groupName string) error {
	// Check if user already in group
	cmd := exec.Command("id", "-nG", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to check user groups: %w", err)
	}

	groups := strings.Fields(string(output))
	for _, group := range groups {
		if group == groupName {
			// User already in group
			return nil
		}
	}

	// Add user to group
	cmd = exec.Command("usermod", "-aG", groupName, username)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("usermod failed: %s: %w", string(output), err)
	}

	logger.Info("Added user to SMB group",
		zap.String("user", username),
		zap.String("group", groupName))

	return nil
}

// FixExistingSharePermissions fixes permissions for all existing shares
// Should be called once at server startup to ensure all shares have correct permissions
func FixExistingSharePermissions() error {
	var shares []models.Share
	if err := database.DB.Find(&shares).Error; err != nil {
		return fmt.Errorf("failed to list shares: %w", err)
	}

	logger.Info("Fixing permissions for existing shares", zap.Int("count", len(shares)))

	fixedCount := 0
	errorCount := 0

	for _, share := range shares {
		// Only fix enabled shares
		if !share.Enabled {
			continue
		}

		// Check if path exists
		if _, err := os.Stat(share.Path); os.IsNotExist(err) {
			logger.Warn("Share path does not exist, skipping permission fix",
				zap.String("share", share.Name),
				zap.String("path", share.Path))
			continue
		}

		// Fix permissions
		if err := setupSharePermissions(&share); err != nil {
			logger.Error("Failed to fix share permissions",
				zap.String("share", share.Name),
				zap.String("path", share.Path),
				zap.Error(err))
			errorCount++
		} else {
			fixedCount++
		}
	}

	logger.Info("Share permission fix completed",
		zap.Int("fixed", fixedCount),
		zap.Int("errors", errorCount),
		zap.Int("total", len(shares)))

	return nil
}

// GetShareStats returns statistics about shares
func GetShareStats() (int, error) {
	var count int64
	if err := database.DB.Model(&models.Share{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// MigrateShares runs database migrations for shares
func MigrateShares() error {
	return database.DB.AutoMigrate(&models.Share{})
}

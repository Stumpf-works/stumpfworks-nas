// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package storage

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/users"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/sysutil"
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

	var validGroups []string
	if s.ValidGroups != "" {
		validGroups = strings.Split(s.ValidGroups, ",")
	}

	return &Share{
		ID:          fmt.Sprintf("%d", s.ID),
		Name:        s.Name,
		Path:        s.Path,
		VolumeID:    s.VolumeID,
		Type:        ShareType(s.Type),
		Description: s.Description,
		Enabled:     s.Enabled,
		ReadOnly:    s.ReadOnly,
		Browseable:  s.Browseable,
		GuestOK:     s.GuestOK,
		ValidUsers:  validUsers,
		ValidGroups: validGroups,
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
		zap.String("volumeId", req.VolumeID),
		zap.String("path", req.Path))

	// Validate that either VolumeID or Path is provided
	if req.VolumeID == "" && req.Path == "" {
		return nil, fmt.Errorf("either volumeId or path must be provided")
	}

	// Resolve the actual path
	sharePath := req.Path
	volumeID := req.VolumeID

	// If VolumeID is provided, resolve the volume's mount point
	if req.VolumeID != "" {
		volume, err := GetVolume(req.VolumeID)
		if err != nil {
			return nil, fmt.Errorf("volume not found: %s", req.VolumeID)
		}
		if volume.Status != VolumeStatusOnline {
			return nil, fmt.Errorf("volume '%s' is not online (status: %s)", req.VolumeID, volume.Status)
		}
		// Use the volume's mount point as the share path
		sharePath = volume.MountPoint
		logger.Info("Resolved volume to mount point",
			zap.String("volumeId", req.VolumeID),
			zap.String("mountPoint", sharePath))
	}

	// Validate path exists
	if _, err := os.Stat(sharePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("path does not exist: %s", sharePath)
	}

	// Validate that all users in ValidUsers exist
	for _, username := range req.ValidUsers {
		if username == "" {
			continue // Skip empty usernames
		}
		if _, err := users.GetUserByUsername(username); err != nil {
			return nil, fmt.Errorf("user '%s' does not exist - cannot add to valid users list", username)
		}
	}

	// Validate that all groups in ValidGroups exist (system groups)
	for _, groupname := range req.ValidGroups {
		if groupname == "" {
			continue // Skip empty group names
		}
		if _, err := user.LookupGroup(groupname); err != nil {
			return nil, fmt.Errorf("group '%s' does not exist - cannot add to valid groups list", groupname)
		}
	}

	// Create database record
	model := &models.Share{
		Name:        req.Name,
		Path:        sharePath, // Use resolved path (from volume or manual)
		VolumeID:    volumeID,  // Store volume reference if provided
		Type:        string(req.Type),
		Description: req.Description,
		Enabled:     true,
		ReadOnly:    req.ReadOnly,
		Browseable:  req.Browseable,
		GuestOK:     req.GuestOK,
		ValidUsers:  strings.Join(req.ValidUsers, ","),
		ValidGroups: strings.Join(req.ValidGroups, ","),
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

	// Validate that all users in ValidUsers exist
	for _, username := range req.ValidUsers {
		if username == "" {
			continue // Skip empty usernames
		}
		if _, err := users.GetUserByUsername(username); err != nil {
			return nil, fmt.Errorf("user '%s' does not exist - cannot add to valid users list", username)
		}
	}

	// Validate that all groups in ValidGroups exist (system groups)
	for _, groupname := range req.ValidGroups {
		if groupname == "" {
			continue // Skip empty group names
		}
		if _, err := user.LookupGroup(groupname); err != nil {
			return nil, fmt.Errorf("group '%s' does not exist - cannot add to valid groups list", groupname)
		}
	}

	// Update fields
	model.Name = req.Name
	model.Path = req.Path
	model.Description = req.Description
	model.ReadOnly = req.ReadOnly
	model.Browseable = req.Browseable
	model.GuestOK = req.GuestOK
	model.ValidUsers = strings.Join(req.ValidUsers, ",")
	model.ValidGroups = strings.Join(req.ValidGroups, ",")

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

// configureSMBShare configures a Samba share by writing it directly to smb.conf
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

	// Build Samba share configuration
	shareConfig := buildSambaShareConfig(share)

	// Write share directly to smb.conf
	if err := addShareToSmbConf(share.Name, shareConfig); err != nil {
		return fmt.Errorf("failed to add share to smb.conf: %w", err)
	}

	// Reload Samba
	reloadSamba()

	return nil
}

// buildSambaShareConfig builds the configuration text for a share
func buildSambaShareConfig(share *models.Share) string {
	config := fmt.Sprintf(`[%s]
   path = %s
   comment = %s
   browseable = %s
   read only = %s
   guest ok = %s`,
		share.Name,
		share.Path,
		share.Description,
		boolToYesNo(share.Browseable),
		boolToYesNo(share.ReadOnly),
		boolToYesNo(share.GuestOK))

	// Build valid users list (combining individual users and groups)
	var validEntries []string

	// Add individual users
	if share.ValidUsers != "" {
		users := strings.Split(share.ValidUsers, ",")
		for _, user := range users {
			user = strings.TrimSpace(user)
			if user != "" {
				validEntries = append(validEntries, user)
			}
		}
	}

	// Add groups (prefixed with @ for Samba group syntax)
	if share.ValidGroups != "" {
		groups := strings.Split(share.ValidGroups, ",")
		for _, group := range groups {
			group = strings.TrimSpace(group)
			if group != "" {
				validEntries = append(validEntries, "@"+group)
			}
		}
	}

	// Add valid users directive if we have any entries
	if len(validEntries) > 0 {
		config += fmt.Sprintf("\n   valid users = %s", strings.Join(validEntries, " "))
	}

	return config
}

// addShareToSmbConf adds or updates a share in smb.conf
func addShareToSmbConf(shareName, shareConfig string) error {
	smbConfPath := "/etc/samba/smb.conf"

	// Read current smb.conf
	data, err := os.ReadFile(smbConfPath)
	if err != nil {
		return fmt.Errorf("failed to read smb.conf: %w", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	// Remove existing share with this name if it exists
	lines = removeShareFromLines(lines, shareName)

	// Add the new share at the end
	marker := fmt.Sprintf("# Share '%s' - Managed by Stumpf.Works NAS", shareName)

	// Add newline before marker if file doesn't end with one
	if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) != "" {
		lines = append(lines, "")
	}

	lines = append(lines, marker)
	for _, line := range strings.Split(shareConfig, "\n") {
		lines = append(lines, line)
	}
	lines = append(lines, "") // Empty line after share

	// Write back to smb.conf
	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(smbConfPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write smb.conf: %w", err)
	}

	logger.Info("Share added to smb.conf", zap.String("share", shareName))
	return nil
}

// removeShareFromLines removes a share section from smb.conf lines
func removeShareFromLines(lines []string, shareName string) []string {
	var newLines []string
	skipUntilNextSection := false
	shareMarker := fmt.Sprintf("# Share '%s' - Managed by Stumpf.Works NAS", shareName)
	shareSection := fmt.Sprintf("[%s]", shareName)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if this is our managed share marker
		if trimmed == shareMarker {
			skipUntilNextSection = true
			continue
		}

		// Also detect share section by name (for backward compatibility)
		if trimmed == shareSection {
			// Check if previous line is our marker
			if i > 0 && strings.TrimSpace(lines[i-1]) == shareMarker {
				// Already handled by marker check above
			} else {
				// This is an unmanaged share with the same name - remove it anyway
				logger.Warn("Removing unmanaged share with same name", zap.String("share", shareName))
				skipUntilNextSection = true
				continue
			}
		}

		// If we're skipping, check if we've reached the next section
		if skipUntilNextSection {
			if strings.HasPrefix(trimmed, "[") && trimmed != shareSection {
				// New section started, stop skipping
				skipUntilNextSection = false
				newLines = append(newLines, line)
			}
			// Skip this line (it's part of the share we're removing)
			continue
		}

		newLines = append(newLines, line)
	}

	return newLines
}

// removeSMBShare removes a Samba share from smb.conf
func removeSMBShare(share *models.Share) error {
	smbConfPath := "/etc/samba/smb.conf"

	// Read current smb.conf
	data, err := os.ReadFile(smbConfPath)
	if err != nil {
		return fmt.Errorf("failed to read smb.conf: %w", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	// Remove the share
	newLines := removeShareFromLines(lines, share.Name)

	// Write back to smb.conf
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(smbConfPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write smb.conf: %w", err)
	}

	logger.Info("Share removed from smb.conf", zap.String("share", share.Name))

	// Reload Samba
	reloadSamba()

	return nil
}

// reloadSamba reloads the Samba service to apply configuration changes
func reloadSamba() {
	// Try systemctl first
	cmd := exec.Command("systemctl", "reload", "smbd")
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Warn("Failed to reload smbd via systemctl",
			zap.String("output", string(output)),
			zap.Error(err))
		// Try service command as fallback
		cmd = exec.Command("service", "smbd", "reload")
		if output, err := cmd.CombinedOutput(); err != nil {
			logger.Warn("Failed to reload smbd via service",
				zap.String("output", string(output)),
				zap.Error(err))
		}
	}

	// Also reload nmbd
	cmd = exec.Command("systemctl", "reload", "nmbd")
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Debug("Failed to reload nmbd", zap.String("output", string(output)))
	}
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

// RepairSambaConfig repairs common issues in smb.conf
// This function is called on startup to migrate from old include-based config to direct shares in smb.conf
func RepairSambaConfig() error {
	smbConfPath := "/etc/samba/smb.conf"
	sharesDir := "/etc/samba/shares.d"

	// Check if smb.conf exists
	if _, err := os.Stat(smbConfPath); os.IsNotExist(err) {
		logger.Debug("Samba config not found, skipping repair", zap.String("path", smbConfPath))
		return nil // Not an error - Samba might not be installed
	}

	// Read current config
	data, err := os.ReadFile(smbConfPath)
	if err != nil {
		return err
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	// Step 1: Remove any include directives (old system)
	var cleanedLines []string
	removedInclude := false

	for _, line := range lines {
		// Skip include directives and their comments
		if strings.Contains(line, "include = /etc/samba/shares.d") ||
		   strings.Contains(line, "Include dynamic share configurations") {
			removedInclude = true
			logger.Info("Removing obsolete include directive", zap.String("line", strings.TrimSpace(line)))
			continue
		}
		cleanedLines = append(cleanedLines, line)
	}

	// Step 2: Migrate shares from shares.d/*.conf to smb.conf (if shares.d exists)
	migratedShares := 0
	if _, err := os.Stat(sharesDir); err == nil {
		entries, err := os.ReadDir(sharesDir)
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".conf") {
					shareName := strings.TrimSuffix(entry.Name(), ".conf")
					shareFilePath := filepath.Join(sharesDir, entry.Name())

					// Read the share config
					shareData, err := os.ReadFile(shareFilePath)
					if err != nil {
						logger.Warn("Failed to read share file for migration",
							zap.String("file", shareFilePath),
							zap.Error(err))
						continue
					}

					shareConfig := string(shareData)

					// Add marker comment
					marker := fmt.Sprintf("# Share '%s' - Managed by Stumpf.Works NAS", shareName)

					// Add to cleaned lines
					if len(cleanedLines) > 0 && strings.TrimSpace(cleanedLines[len(cleanedLines)-1]) != "" {
						cleanedLines = append(cleanedLines, "")
					}
					cleanedLines = append(cleanedLines, marker)
					for _, line := range strings.Split(strings.TrimSpace(shareConfig), "\n") {
						cleanedLines = append(cleanedLines, line)
					}
					cleanedLines = append(cleanedLines, "")

					migratedShares++
					logger.Info("Migrated share from shares.d to smb.conf",
						zap.String("share", shareName))
				}
			}
		}
	}

	// Only write back if we made changes
	if removedInclude || migratedShares > 0 {
		newContent := strings.Join(cleanedLines, "\n")
		if err := os.WriteFile(smbConfPath, []byte(newContent), 0644); err != nil {
			return fmt.Errorf("failed to write repaired smb.conf: %w", err)
		}

		logger.Info("Samba configuration repaired",
			zap.Bool("removedInclude", removedInclude),
			zap.Int("migratedShares", migratedShares))

		// Reload Samba to apply changes
		reloadSamba()
	} else {
		logger.Debug("Samba config is correct, no repair needed")
	}

	return nil
}

// ensureSambaInclude is deprecated - we now write shares directly to smb.conf
// This function is kept for backward compatibility but does nothing
func ensureSambaInclude() error {
	// No longer needed - shares are written directly to smb.conf
	return nil
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
	getentPath := sysutil.FindCommand("getent")
	cmd := exec.Command(getentPath, "group", groupName)
	if err := cmd.Run(); err == nil {
		// Group exists
		return nil
	}

	// Group doesn't exist, create it with retry logic
	groupaddPath := sysutil.FindCommand("groupadd")

	// Retry logic for /etc/group lock contention
	// Increased retries due to severe lock contention during service startup
	maxRetries := 10
	baseDelay := 150 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 150ms, 300ms, 600ms, 1200ms, 2400ms, 4800ms, 9600ms, 19200ms, 38400ms
			delay := baseDelay * time.Duration(1<<uint(attempt-1))
			logger.Info("Retrying groupadd after delay",
				zap.String("group", groupName),
				zap.Int("attempt", attempt+1),
				zap.Duration("delay", delay))
			time.Sleep(delay)
		}

		cmd = exec.Command(groupaddPath, groupName)
		output, err := cmd.CombinedOutput()
		if err == nil {
			logger.Info("Created SMB group successfully",
				zap.String("group", groupName),
				zap.String("groupadd_path", groupaddPath),
				zap.Int("attempts", attempt+1))
			return nil
		}

		// Check if error is due to /etc/group lock contention
		outputStr := string(output)
		isLockError := strings.Contains(outputStr, "konnte nicht gesperrt werden") ||
			strings.Contains(outputStr, "cannot lock") ||
			strings.Contains(outputStr, "unable to lock") ||
			strings.Contains(outputStr, "temporarily unavailable") ||
			strings.Contains(outputStr, "group") && strings.Contains(outputStr, "lock")

		// If group already exists (race condition), that's fine
		if strings.Contains(outputStr, "already exists") {
			logger.Info("SMB group already exists (race condition resolved)",
				zap.String("group", groupName))
			return nil
		}

		// If it's not a lock error, or we're on the last attempt, return the error
		if !isLockError || attempt == maxRetries-1 {
			return fmt.Errorf("failed to create group %s: %s: %w", groupName, outputStr, err)
		}

		logger.Info("groupadd lock contention detected, will retry",
			zap.String("group", groupName),
			zap.Int("attempt", attempt+1),
			zap.Int("max_retries", maxRetries),
			zap.String("error", outputStr))
	}

	return fmt.Errorf("groupadd failed after %d attempts for group %s", maxRetries, groupName)
}

// setShareGroupOwnership sets the group ownership of a path
func setShareGroupOwnership(path, groupName string) error {
	// Use chgrp to set group ownership
	chgrpPath := sysutil.FindCommand("chgrp")
	cmd := exec.Command(chgrpPath, groupName, path)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("chgrp failed: %s: %w", string(output), err)
	}
	return nil
}

// addUserToGroup adds a user to a group
func addUserToGroup(username, groupName string) error {
	// Check if user already in group
	idPath := sysutil.FindCommand("id")
	cmd := exec.Command(idPath, "-nG", username)
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
	usermodPath := sysutil.FindCommand("usermod")
	cmd = exec.Command(usermodPath, "-aG", groupName, username)
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

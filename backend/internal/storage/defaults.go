package storage

import (
	"os"
	"path/filepath"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// EnsureDefaultShares ensures that at least one share exists for users to access
// This prevents "Access Denied" errors when no shares are configured
func EnsureDefaultShares() error {
	logger.Info("Checking for default shares...")

	// Check if any shares exist
	var count int64
	if err := database.DB.Model(&models.Share{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		logger.Info("Shares already exist, skipping default share creation", zap.Int64("count", count))
		return nil
	}

	logger.Info("No shares found, creating default share...")

	// Define default storage path
	defaultPath := "/mnt/stumpfworks-nas/storage"

	// Try common alternative paths if default doesn't exist
	alternativePaths := []string{
		defaultPath,
		"/mnt/storage",
		"/data/storage",
		"/storage",
		"/home/storage",
		"/tmp/stumpfworks-nas-storage", // Fallback for testing
	}

	var selectedPath string
	for _, path := range alternativePaths {
		if err := os.MkdirAll(path, 0755); err == nil {
			selectedPath = path
			logger.Info("Using storage path", zap.String("path", path))
			break
		}
	}

	if selectedPath == "" {
		logger.Warn("Could not create storage directory, file access may be limited")
		return nil // Don't fail - system can still work without shares
	}

	// Create subdirectories for organization
	subdirs := []string{"files", "media", "documents", "backups"}
	for _, subdir := range subdirs {
		subdirPath := filepath.Join(selectedPath, subdir)
		if err := os.MkdirAll(subdirPath, 0755); err != nil {
			logger.Warn("Failed to create subdirectory", zap.String("path", subdirPath), zap.Error(err))
		}
	}

	// Create default SMB share for general file storage
	defaultShare := &models.Share{
		Name:        "Files",
		Path:        filepath.Join(selectedPath, "files"),
		Type:        "smb",
		Description: "Default file storage share",
		Enabled:     true,
		ReadOnly:    false,
		Browseable:  true,
		GuestOK:     false, // Require authentication
		ValidUsers:  "",    // Empty = all authenticated users can access
	}

	if err := database.DB.Create(defaultShare).Error; err != nil {
		logger.Error("Failed to create default share", zap.Error(err))
		return err
	}

	// Configure Samba for the share
	if err := configureSMBShare(defaultShare); err != nil {
		logger.Warn("Failed to configure Samba for default share",
			zap.Error(err),
			zap.String("note", "Share created but network access may not work"))
		// Don't fail - share exists in DB for File Manager
	}

	logger.Info("Default share created successfully",
		zap.String("name", defaultShare.Name),
		zap.String("path", defaultShare.Path))

	// Create additional media share
	mediaShare := &models.Share{
		Name:        "Media",
		Path:        filepath.Join(selectedPath, "media"),
		Type:        "smb",
		Description: "Media files (videos, music, photos)",
		Enabled:     true,
		ReadOnly:    false,
		Browseable:  true,
		GuestOK:     false,
		ValidUsers:  "",
	}

	if err := database.DB.Create(mediaShare).Error; err != nil {
		logger.Warn("Failed to create media share", zap.Error(err))
		// Not critical - Files share is enough
	} else {
		configureSMBShare(mediaShare)
		logger.Info("Media share created successfully")
	}

	// Create README in storage directory
	readmePath := filepath.Join(selectedPath, "README.txt")
	readmeContent := `Stumpf.Works NAS - Default Storage

This directory contains the default file storage for your NAS system.

Directory Structure:
  files/      - General file storage (default share)
  media/      - Media files (videos, music, photos)
  documents/  - Document storage
  backups/    - Backup files

You can access these folders via:
  - Web UI: Files app
  - Windows: \\<nas-ip>\Files or \\<nas-ip>\Media
  - macOS: smb://<nas-ip>/Files
  - Linux: smb://<nas-ip>/Files

To add more shares, use the Storage app in the web interface.
`
	os.WriteFile(readmePath, []byte(readmeContent), 0644)

	return nil
}

package addons

import (
	"fmt"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/system"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// Manager manages addon installation and lifecycle
type Manager struct {
	packageInstaller *PackageInstaller
}

// NewManager creates a new addon manager
func NewManager(shell *system.ShellExecutor) *Manager {
	return &Manager{
		packageInstaller: NewPackageInstaller(shell),
	}
}

// ListAvailableAddons returns all available addons
func (m *Manager) ListAvailableAddons() []Manifest {
	return BuiltinAddons
}

// GetAddon returns a specific addon by ID
func (m *Manager) GetAddon(addonID string) (*Manifest, error) {
	for _, addon := range BuiltinAddons {
		if addon.ID == addonID {
			return &addon, nil
		}
	}
	return nil, fmt.Errorf("addon not found: %s", addonID)
}

// GetAddonStatus returns the installation status of an addon
func (m *Manager) GetAddonStatus(addonID string) (*InstallationStatus, error) {
	// Get addon manifest
	addon, err := m.GetAddon(addonID)
	if err != nil {
		return nil, err
	}

	// Get installation record from database
	var installation models.AddonInstallation
	result := database.DB.Where("addon_id = ?", addonID).First(&installation)

	status := &InstallationStatus{
		AddonID:   addonID,
		Installed: false,
	}

	if result.Error == nil {
		// Found installation record
		status.Installed = installation.Installed
		status.Version = installation.Version
		status.InstallDate = installation.InstallDate.Format(time.RFC3339)
		status.Error = installation.Error
	}

	// Check if packages are installed
	status.PackagesOK = m.packageInstaller.AreAllPackagesInstalled(addon.SystemPackages)

	// Check if services are running
	status.ServicesOK = m.packageInstaller.AreAllServicesRunning(addon.Services)

	return status, nil
}

// InstallAddon installs an addon
func (m *Manager) InstallAddon(addonID string) error {
	logger.Info("Installing addon", zap.String("addon_id", addonID))

	// Get addon manifest
	addon, err := m.GetAddon(addonID)
	if err != nil {
		return err
	}

	// Check if already installed
	status, err := m.GetAddonStatus(addonID)
	if err == nil && status.Installed && status.PackagesOK {
		return fmt.Errorf("addon already installed: %s", addonID)
	}

	// Create or update installation record
	var installation models.AddonInstallation
	result := database.DB.Where("addon_id = ?", addonID).First(&installation)
	if result.Error != nil {
		// Create new record
		installation = models.AddonInstallation{
			AddonID:     addonID,
			Version:     addon.Version,
			Installed:   false,
			InstallDate: time.Now(),
		}
	}

	// Install system packages
	if len(addon.SystemPackages) > 0 {
		logger.Info("Installing system packages", zap.Strings("packages", addon.SystemPackages))
		if err := m.packageInstaller.InstallPackages(addon.SystemPackages); err != nil {
			installation.Error = err.Error()
			database.DB.Save(&installation)
			return fmt.Errorf("failed to install packages: %w", err)
		}
	}

	// Enable and start services
	if len(addon.Services) > 0 {
		logger.Info("Enabling services", zap.Strings("services", addon.Services))
		for _, service := range addon.Services {
			if err := m.packageInstaller.EnableService(service); err != nil {
				installation.Error = err.Error()
				database.DB.Save(&installation)
				return fmt.Errorf("failed to enable service %s: %w", service, err)
			}
		}
	}

	// Run install script if specified
	if addon.InstallScript != "" {
		logger.Info("Running install script", zap.String("addon_id", addonID))
		// TODO: Execute install script
	}

	// Mark as installed
	installation.Installed = true
	installation.Version = addon.Version
	installation.InstallDate = time.Now()
	installation.Error = ""

	if err := database.DB.Save(&installation).Error; err != nil {
		return fmt.Errorf("failed to save installation record: %w", err)
	}

	logger.Info("Addon installed successfully", zap.String("addon_id", addonID))
	return nil
}

// UninstallAddon uninstalls an addon
func (m *Manager) UninstallAddon(addonID string) error {
	logger.Info("Uninstalling addon", zap.String("addon_id", addonID))

	// Get addon manifest
	addon, err := m.GetAddon(addonID)
	if err != nil {
		return err
	}

	// Check if installed
	var installation models.AddonInstallation
	result := database.DB.Where("addon_id = ?", addonID).First(&installation)
	if result.Error != nil || !installation.Installed {
		return fmt.Errorf("addon not installed: %s", addonID)
	}

	// Run uninstall script if specified
	if addon.UninstallScript != "" {
		logger.Info("Running uninstall script", zap.String("addon_id", addonID))
		// TODO: Execute uninstall script
	}

	// Disable and stop services
	if len(addon.Services) > 0 {
		logger.Info("Disabling services", zap.Strings("services", addon.Services))
		for _, service := range addon.Services {
			if err := m.packageInstaller.DisableService(service); err != nil {
				logger.Warn("Failed to disable service", zap.String("service", service), zap.Error(err))
				// Continue anyway
			}
		}
	}

	// Uninstall system packages
	if len(addon.SystemPackages) > 0 {
		logger.Info("Uninstalling system packages", zap.Strings("packages", addon.SystemPackages))
		if err := m.packageInstaller.UninstallPackages(addon.SystemPackages); err != nil {
			return fmt.Errorf("failed to uninstall packages: %w", err)
		}
	}

	// Mark as uninstalled
	installation.Installed = false
	if err := database.DB.Save(&installation).Error; err != nil {
		return fmt.Errorf("failed to save installation record: %w", err)
	}

	logger.Info("Addon uninstalled successfully", zap.String("addon_id", addonID))
	return nil
}

// GetAllAddonsWithStatus returns all addons with their installation status
func (m *Manager) GetAllAddonsWithStatus() ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}

	for _, addon := range BuiltinAddons {
		status, err := m.GetAddonStatus(addon.ID)
		if err != nil {
			logger.Warn("Failed to get addon status", zap.String("addon_id", addon.ID), zap.Error(err))
			status = &InstallationStatus{
				AddonID:   addon.ID,
				Installed: false,
			}
		}

		addonWithStatus := map[string]interface{}{
			"manifest": addon,
			"status":   status,
		}
		result = append(result, addonWithStatus)
	}

	return result, nil
}

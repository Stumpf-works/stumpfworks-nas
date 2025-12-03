// Package timemachine provides macOS Time Machine backup server management
package timemachine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/system"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Manager manages Time Machine backup service
type Manager struct {
	shell *system.ShellExecutor
}

// NewManager creates a new Time Machine manager
func NewManager(shell *system.ShellExecutor) *Manager {
	return &Manager{
		shell: shell,
	}
}

// GetConfig returns the current Time Machine configuration
func (m *Manager) GetConfig() (*models.TimeMachineConfig, error) {
	var config models.TimeMachineConfig
	result := database.DB.First(&config)

	if result.Error == gorm.ErrRecordNotFound {
		// Create default configuration
		config = models.TimeMachineConfig{
			Enabled:        false,
			ShareName:      "TimeMachine",
			BasePath:       "/mnt/storage/timemachine",
			DefaultQuotaGB: 500,
			AutoDiscovery:  true,
			UseAFP:         false,
			UseSMB:         true,
			SMBVersion:     "3",
		}
		if err := database.DB.Create(&config).Error; err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &config, nil
}

// UpdateConfig updates the Time Machine configuration
func (m *Manager) UpdateConfig(config *models.TimeMachineConfig) error {
	if err := database.DB.Save(config).Error; err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	// Apply configuration changes
	if config.Enabled {
		if err := m.Enable(); err != nil {
			return fmt.Errorf("failed to enable Time Machine: %w", err)
		}
	} else {
		if err := m.Disable(); err != nil {
			return fmt.Errorf("failed to disable Time Machine: %w", err)
		}
	}

	return nil
}

// Enable enables the Time Machine service
func (m *Manager) Enable() error {
	config, err := m.GetConfig()
	if err != nil {
		return err
	}

	// Create base directory if it doesn't exist
	if err := os.MkdirAll(config.BasePath, 0755); err != nil {
		return fmt.Errorf("failed to create base directory: %w", err)
	}

	// Configure Samba
	if config.UseSMB {
		if err := m.configureSamba(); err != nil {
			return fmt.Errorf("failed to configure Samba: %w", err)
		}
	}

	// Enable service discovery (Avahi)
	if config.AutoDiscovery {
		if err := m.enableAvahi(); err != nil {
			logger.Warn("Failed to enable Avahi discovery", zap.Error(err))
		}
	}

	config.Enabled = true
	return database.DB.Save(config).Error
}

// Disable disables the Time Machine service
func (m *Manager) Disable() error {
	config, err := m.GetConfig()
	if err != nil {
		return err
	}

	// Remove Samba configuration
	if err := m.removeSambaConfig(); err != nil {
		logger.Warn("Failed to remove Samba config", zap.Error(err))
	}

	config.Enabled = false
	return database.DB.Save(config).Error
}

// ListDevices returns all registered Time Machine devices
func (m *Manager) ListDevices() ([]models.TimeMachineDevice, error) {
	var devices []models.TimeMachineDevice
	if err := database.DB.Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

// GetDevice returns a specific device by ID
func (m *Manager) GetDevice(id uint) (*models.TimeMachineDevice, error) {
	var device models.TimeMachineDevice
	if err := database.DB.First(&device, id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

// CreateDevice creates a new Time Machine device registration
func (m *Manager) CreateDevice(device *models.TimeMachineDevice) error {
	config, err := m.GetConfig()
	if err != nil {
		return err
	}

	// Set default quota if not specified
	if device.QuotaGB == 0 {
		device.QuotaGB = config.DefaultQuotaGB
	}

	// Set share path if not specified
	if device.SharePath == "" {
		device.SharePath = filepath.Join(config.BasePath, sanitizeDeviceName(device.DeviceName))
	}

	// Create device directory
	if err := os.MkdirAll(device.SharePath, 0700); err != nil {
		return fmt.Errorf("failed to create device directory: %w", err)
	}

	// Set proper ownership (stumpfs user)
	if err := m.shell.Execute("chown", "-R", "stumpfs:stumpfs", device.SharePath); err != nil {
		logger.Warn("Failed to set directory ownership", zap.Error(err))
	}

	// Save to database
	if err := database.DB.Create(device).Error; err != nil {
		return err
	}

	// Update Samba configuration
	if config.UseSMB && config.Enabled {
		if err := m.configureSamba(); err != nil {
			return fmt.Errorf("failed to update Samba config: %w", err)
		}
	}

	return nil
}

// UpdateDevice updates a Time Machine device
func (m *Manager) UpdateDevice(device *models.TimeMachineDevice) error {
	if err := database.DB.Save(device).Error; err != nil {
		return err
	}

	// Update Samba configuration
	config, err := m.GetConfig()
	if err != nil {
		return err
	}

	if config.UseSMB && config.Enabled {
		if err := m.configureSamba(); err != nil {
			return fmt.Errorf("failed to update Samba config: %w", err)
		}
	}

	return nil
}

// DeleteDevice removes a Time Machine device registration
func (m *Manager) DeleteDevice(id uint) error {
	var device models.TimeMachineDevice
	if err := database.DB.First(&device, id).Error; err != nil {
		return err
	}

	// Delete from database (soft delete)
	if err := database.DB.Delete(&device).Error; err != nil {
		return err
	}

	// Update Samba configuration
	config, err := m.GetConfig()
	if err != nil {
		return err
	}

	if config.UseSMB && config.Enabled {
		if err := m.configureSamba(); err != nil {
			return fmt.Errorf("failed to update Samba config: %w", err)
		}
	}

	return nil
}

// UpdateDeviceUsage updates the used space for a device
func (m *Manager) UpdateDeviceUsage(id uint) error {
	device, err := m.GetDevice(id)
	if err != nil {
		return err
	}

	// Calculate directory size using du
	cmd := exec.Command("du", "-sb", device.SharePath)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to calculate directory size: %w", err)
	}

	// Parse output: "12345678\t/path/to/dir"
	parts := strings.Fields(string(output))
	if len(parts) < 1 {
		return fmt.Errorf("invalid du output")
	}

	var bytes int64
	if _, err := fmt.Sscanf(parts[0], "%d", &bytes); err != nil {
		return fmt.Errorf("failed to parse size: %w", err)
	}

	// Convert bytes to GB
	device.UsedGB = float64(bytes) / (1024 * 1024 * 1024)

	// Update last seen timestamp
	now := time.Now()
	device.LastSeen = &now

	return database.DB.Save(device).Error
}

// UpdateAllDeviceUsages updates usage for all devices
func (m *Manager) UpdateAllDeviceUsages() error {
	devices, err := m.ListDevices()
	if err != nil {
		return err
	}

	for _, device := range devices {
		if err := m.UpdateDeviceUsage(device.ID); err != nil {
			logger.Warn("Failed to update device usage",
				zap.Uint("device_id", device.ID),
				zap.String("device_name", device.DeviceName),
				zap.Error(err))
		}
	}

	return nil
}

// sanitizeDeviceName creates a safe directory name from a device name
func sanitizeDeviceName(name string) string {
	// Replace spaces and special characters with underscores
	safe := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, name)
	return safe
}

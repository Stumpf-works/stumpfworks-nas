package timemachine

import (
	"fmt"
	"os"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

const (
	sambaConfigFile = "/etc/samba/smb.conf"
	timeMachineMarkerStart = "# === Time Machine Configuration - Start (Managed by StumpfWorks) ==="
	timeMachineMarkerEnd   = "# === Time Machine Configuration - End ==="
)

// configureSamba configures Samba for Time Machine support
func (m *Manager) configureSamba() error {
	logger.Info("Configuring Samba for Time Machine")

	config, err := m.GetConfig()
	if err != nil {
		return err
	}

	devices, err := m.ListDevices()
	if err != nil {
		return err
	}

	// Read current smb.conf
	content, err := os.ReadFile(sambaConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read smb.conf: %w", err)
	}

	configStr := string(content)

	// Remove existing Time Machine configuration
	configStr = removeTimeMachineConfig(configStr)

	// Build new Time Machine configuration
	var tmConfig strings.Builder
	tmConfig.WriteString("\n" + timeMachineMarkerStart + "\n\n")

	// Add share for each enabled device
	for _, device := range devices {
		if !device.Enabled {
			continue
		}

		shareName := fmt.Sprintf("TimeMachine-%s", sanitizeDeviceName(device.DeviceName))
		tmConfig.WriteString(fmt.Sprintf("[%s]\n", shareName))
		tmConfig.WriteString(fmt.Sprintf("  path = %s\n", device.SharePath))
		tmConfig.WriteString("  browseable = yes\n")
		tmConfig.WriteString("  writeable = yes\n")
		tmConfig.WriteString("  create mask = 0600\n")
		tmConfig.WriteString("  directory mask = 0700\n")

		// Time Machine specific settings
		tmConfig.WriteString("  vfs objects = catia fruit streams_xattr\n")
		tmConfig.WriteString("  fruit:aapl = yes\n")
		tmConfig.WriteString("  fruit:time machine = yes\n")
		tmConfig.WriteString("  fruit:model = MacSamba\n")

		// Set quota if specified
		if device.QuotaGB > 0 {
			// Quota in MB for Samba
			quotaMB := device.QuotaGB * 1024
			tmConfig.WriteString(fmt.Sprintf("  fruit:time machine max size = %dM\n", quotaMB))
		}

		// Access control
		if device.Username != "" {
			tmConfig.WriteString(fmt.Sprintf("  valid users = %s\n", device.Username))
		} else {
			tmConfig.WriteString("  valid users = stumpfs\n")
		}

		tmConfig.WriteString("  force user = stumpfs\n")
		tmConfig.WriteString("  force group = stumpfs\n")

		tmConfig.WriteString("\n")
	}

	tmConfig.WriteString(timeMachineMarkerEnd + "\n")

	// Append Time Machine configuration
	configStr += tmConfig.String()

	// Write updated configuration
	if err := os.WriteFile(sambaConfigFile, []byte(configStr), 0644); err != nil {
		return fmt.Errorf("failed to write smb.conf: %w", err)
	}

	// Validate configuration
	if err := m.validateSambaConfig(); err != nil {
		logger.Error("Samba configuration validation failed", zap.Error(err))
		return err
	}

	// Reload Samba
	if err := m.reloadSamba(); err != nil {
		return fmt.Errorf("failed to reload Samba: %w", err)
	}

	logger.Info("Samba configured successfully for Time Machine")
	return nil
}

// removeSambaConfig removes Time Machine configuration from Samba
func (m *Manager) removeSambaConfig() error {
	logger.Info("Removing Time Machine configuration from Samba")

	// Read current smb.conf
	content, err := os.ReadFile(sambaConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read smb.conf: %w", err)
	}

	configStr := string(content)

	// Remove Time Machine configuration
	configStr = removeTimeMachineConfig(configStr)

	// Write updated configuration
	if err := os.WriteFile(sambaConfigFile, []byte(configStr), 0644); err != nil {
		return fmt.Errorf("failed to write smb.conf: %w", err)
	}

	// Reload Samba
	if err := m.reloadSamba(); err != nil {
		return fmt.Errorf("failed to reload Samba: %w", err)
	}

	logger.Info("Time Machine configuration removed from Samba")
	return nil
}

// removeTimeMachineConfig removes the Time Machine section from smb.conf
func removeTimeMachineConfig(config string) string {
	startIdx := strings.Index(config, timeMachineMarkerStart)
	if startIdx == -1 {
		return config
	}

	endIdx := strings.Index(config, timeMachineMarkerEnd)
	if endIdx == -1 {
		return config
	}

	// Remove everything between markers, including the end marker line
	endLineIdx := strings.Index(config[endIdx:], "\n")
	if endLineIdx != -1 {
		endIdx += endLineIdx + 1
	} else {
		endIdx += len(timeMachineMarkerEnd)
	}

	return config[:startIdx] + config[endIdx:]
}

// validateSambaConfig validates the Samba configuration using testparm
func (m *Manager) validateSambaConfig() error {
	result, err := m.shell.Execute("testparm", "-s", "--suppress-prompt")
	if err != nil {
		return fmt.Errorf("testparm validation failed: %w (output: %s)", err, result.Stderr)
	}
	return nil
}

// reloadSamba reloads the Samba service
func (m *Manager) reloadSamba() error {
	logger.Info("Reloading Samba service")

	// Try reload first
	result, err := m.shell.Execute("systemctl", "reload", "smbd")
	if err != nil {
		logger.Warn("Failed to reload Samba, trying restart", zap.Error(err))
		// If reload fails, try restart
		result, err = m.shell.Execute("systemctl", "restart", "smbd")
		if err != nil {
			return fmt.Errorf("failed to restart Samba: %w (output: %s)", err, result.Stderr)
		}
	}

	// Also reload nmbd for NetBIOS name service
	if _, err := m.shell.Execute("systemctl", "reload", "nmbd"); err != nil {
		logger.Warn("Failed to reload nmbd", zap.Error(err))
	}

	return nil
}

// GetSambaStatus returns the status of the Samba service
func (m *Manager) GetSambaStatus() (bool, error) {
	result, err := m.shell.Execute("systemctl", "is-active", "smbd")
	if err != nil {
		return false, nil
	}
	return strings.TrimSpace(result.Stdout) == "active", nil
}

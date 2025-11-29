// Package ha provides High Availability features for Stumpf.Works NAS
package ha

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// KeepalivedManager manages Keepalived for Virtual IP (VIP) management
type KeepalivedManager struct {
	shell   executor.ShellExecutor
	enabled bool
}

// VIPConfig represents a Virtual IP configuration
type VIPConfig struct {
	ID            string   `json:"id"`              // Unique identifier
	VirtualIP     string   `json:"virtual_ip"`      // Virtual IP address (e.g., 192.168.1.100)
	Interface     string   `json:"interface"`       // Network interface (e.g., eth0)
	RouterID      int      `json:"router_id"`       // VRRP router ID (1-255)
	Priority      int      `json:"priority"`        // Priority (1-255, higher = master)
	State         string   `json:"state"`           // MASTER or BACKUP
	AuthPass      string   `json:"auth_pass"`       // Authentication password
	VirtualRoutes []string `json:"virtual_routes"`  // Optional virtual routes
	TrackScripts  []string `json:"track_scripts"`   // Optional tracking scripts
}

// VIPStatus represents the status of a VIP
type VIPStatus struct {
	ID          string `json:"id"`
	VirtualIP   string `json:"virtual_ip"`
	Interface   string `json:"interface"`
	State       string `json:"state"`        // MASTER, BACKUP, FAULT
	IsMaster    bool   `json:"is_master"`
	Priority    int    `json:"priority"`
	IsActive    bool   `json:"is_active"`    // Is VIP currently assigned to this node?
}

// NewKeepalivedManager creates a new Keepalived manager
func NewKeepalivedManager(shell executor.ShellExecutor) (*KeepalivedManager, error) {
	manager := &KeepalivedManager{
		shell:   shell,
		enabled: false,
	}

	// Check if keepalived is available
	result, err := shell.Execute("which", "keepalived")
	if err != nil || result.Stdout == "" {
		logger.Warn("keepalived not found, VIP features will be disabled")
		return manager, fmt.Errorf("keepalived not available: install keepalived package")
	}

	manager.enabled = true
	logger.Info("Keepalived manager initialized successfully")
	return manager, nil
}

// IsEnabled returns whether Keepalived is available
func (km *KeepalivedManager) IsEnabled() bool {
	return km.enabled
}

// CreateVIP creates a new Virtual IP configuration
func (km *KeepalivedManager) CreateVIP(config VIPConfig) error {
	if !km.enabled {
		return fmt.Errorf("Keepalived is not enabled")
	}

	// Validate required fields
	if config.VirtualIP == "" || config.Interface == "" {
		return fmt.Errorf("virtual_ip and interface are required")
	}

	// Set defaults
	if config.RouterID == 0 {
		config.RouterID = 51
	}
	if config.Priority == 0 {
		config.Priority = 100
	}
	if config.State == "" {
		config.State = "BACKUP"
	}
	if config.AuthPass == "" {
		config.AuthPass = "StumpfWorks"
	}

	// Generate keepalived configuration
	configContent := km.generateConfig(config)

	// Write to /etc/keepalived/keepalived.conf
	// Note: In production, you'd want to append or manage multiple VIPs
	configPath := "/etc/keepalived/keepalived.conf"

	// Backup existing config if it exists
	km.shell.Execute("sudo", "cp", configPath, configPath+".bak")

	// Write new config
	writeCmd := fmt.Sprintf("echo '%s' | sudo tee %s", configContent, configPath)
	result, err := km.shell.Execute("sh", "-c", writeCmd)
	if err != nil {
		return fmt.Errorf("failed to write keepalived config: %s: %w", result.Stderr, err)
	}

	// Restart keepalived service
	result, err = km.shell.Execute("sudo", "systemctl", "restart", "keepalived")
	if err != nil {
		logger.Error("Failed to restart keepalived", zap.Error(err), zap.String("stderr", result.Stderr))
		return fmt.Errorf("failed to restart keepalived: %s: %w", result.Stderr, err)
	}

	// Enable keepalived to start on boot
	km.shell.Execute("sudo", "systemctl", "enable", "keepalived")

	logger.Info("Keepalived VIP created", zap.String("vip", config.VirtualIP), zap.String("interface", config.Interface))
	return nil
}

// generateConfig generates keepalived configuration file content
func (km *KeepalivedManager) generateConfig(config VIPConfig) string {
	var sb strings.Builder

	sb.WriteString("! Configuration File for keepalived\n\n")
	sb.WriteString("global_defs {\n")
	sb.WriteString("   router_id StumpfWorks_NAS\n")
	sb.WriteString("   enable_script_security\n")
	sb.WriteString("   script_user root\n")
	sb.WriteString("}\n\n")

	// VRRP instance
	sb.WriteString(fmt.Sprintf("vrrp_instance %s {\n", config.ID))
	sb.WriteString(fmt.Sprintf("    state %s\n", config.State))
	sb.WriteString(fmt.Sprintf("    interface %s\n", config.Interface))
	sb.WriteString(fmt.Sprintf("    virtual_router_id %d\n", config.RouterID))
	sb.WriteString(fmt.Sprintf("    priority %d\n", config.Priority))
	sb.WriteString("    advert_int 1\n")

	sb.WriteString("    authentication {\n")
	sb.WriteString("        auth_type PASS\n")
	sb.WriteString(fmt.Sprintf("        auth_pass %s\n", config.AuthPass))
	sb.WriteString("    }\n\n")

	sb.WriteString("    virtual_ipaddress {\n")
	sb.WriteString(fmt.Sprintf("        %s\n", config.VirtualIP))
	sb.WriteString("    }\n")

	// Add virtual routes if specified
	if len(config.VirtualRoutes) > 0 {
		sb.WriteString("\n    virtual_routes {\n")
		for _, route := range config.VirtualRoutes {
			sb.WriteString(fmt.Sprintf("        %s\n", route))
		}
		sb.WriteString("    }\n")
	}

	// Add track scripts if specified
	if len(config.TrackScripts) > 0 {
		sb.WriteString("\n    track_script {\n")
		for _, script := range config.TrackScripts {
			sb.WriteString(fmt.Sprintf("        %s\n", script))
		}
		sb.WriteString("    }\n")
	}

	sb.WriteString("}\n")

	return sb.String()
}

// GetVIPStatus gets the status of a VIP
func (km *KeepalivedManager) GetVIPStatus(vipID string) (*VIPStatus, error) {
	if !km.enabled {
		return nil, fmt.Errorf("Keepalived is not enabled")
	}

	status := &VIPStatus{
		ID:       vipID,
		State:    "Unknown",
		IsMaster: false,
		IsActive: false,
	}

	// Check if keepalived is running
	result, err := km.shell.Execute("systemctl", "is-active", "keepalived")
	if err != nil || strings.TrimSpace(result.Stdout) != "active" {
		status.State = "FAULT"
		return status, nil
	}

	// Read keepalived config to get VIP details
	configPath := "/etc/keepalived/keepalived.conf"
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keepalived config: %w", err)
	}

	configContent := string(data)

	// Parse VIP from config
	vipRegex := regexp.MustCompile(`virtual_ipaddress\s*{\s*([0-9.]+)`)
	if matches := vipRegex.FindStringSubmatch(configContent); len(matches) > 1 {
		status.VirtualIP = matches[1]
	}

	// Parse interface from config
	interfaceRegex := regexp.MustCompile(`interface\s+(\S+)`)
	if matches := interfaceRegex.FindStringSubmatch(configContent); len(matches) > 1 {
		status.Interface = matches[1]
	}

	// Parse priority from config
	priorityRegex := regexp.MustCompile(`priority\s+(\d+)`)
	if matches := priorityRegex.FindStringSubmatch(configContent); len(matches) > 1 {
		fmt.Sscanf(matches[1], "%d", &status.Priority)
	}

	// Check if VIP is assigned to this interface
	if status.VirtualIP != "" && status.Interface != "" {
		result, err := km.shell.Execute("ip", "addr", "show", status.Interface)
		if err == nil && strings.Contains(result.Stdout, status.VirtualIP) {
			status.IsActive = true
			status.IsMaster = true
			status.State = "MASTER"
		} else {
			status.State = "BACKUP"
		}
	}

	return status, nil
}

// ListVIPs lists all configured VIPs
func (km *KeepalivedManager) ListVIPs() ([]VIPStatus, error) {
	if !km.enabled {
		return nil, fmt.Errorf("Keepalived is not enabled")
	}

	vips := []VIPStatus{}

	// Check if config file exists
	configPath := "/etc/keepalived/keepalived.conf"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return vips, nil // No VIPs configured
	}

	// Read keepalived config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keepalived config: %w", err)
	}

	configContent := string(data)

	// Parse VRRP instances
	instanceRegex := regexp.MustCompile(`vrrp_instance\s+(\S+)\s*{`)
	matches := instanceRegex.FindAllStringSubmatch(configContent, -1)

	for _, match := range matches {
		if len(match) > 1 {
			vipID := match[1]
			status, err := km.GetVIPStatus(vipID)
			if err == nil {
				vips = append(vips, *status)
			}
		}
	}

	return vips, nil
}

// DeleteVIP deletes a VIP configuration
func (km *KeepalivedManager) DeleteVIP(vipID string) error {
	if !km.enabled {
		return fmt.Errorf("Keepalived is not enabled")
	}

	// Stop keepalived
	result, err := km.shell.Execute("sudo", "systemctl", "stop", "keepalived")
	if err != nil {
		logger.Warn("Failed to stop keepalived", zap.Error(err), zap.String("stderr", result.Stderr))
	}

	// Remove configuration file
	configPath := "/etc/keepalived/keepalived.conf"
	result, err = km.shell.Execute("sudo", "rm", "-f", configPath)
	if err != nil {
		return fmt.Errorf("failed to delete keepalived config: %s: %w", result.Stderr, err)
	}

	// Disable keepalived service
	km.shell.Execute("sudo", "systemctl", "disable", "keepalived")

	logger.Info("Keepalived VIP deleted", zap.String("id", vipID))
	return nil
}

// PromoteToMaster promotes this node to MASTER (increase priority)
func (km *KeepalivedManager) PromoteToMaster(vipID string) error {
	if !km.enabled {
		return fmt.Errorf("Keepalived is not enabled")
	}

	// Read current config
	configPath := "/etc/keepalived/keepalived.conf"
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read keepalived config: %w", err)
	}

	configContent := string(data)

	// Change state to MASTER and increase priority
	configContent = regexp.MustCompile(`state\s+\S+`).ReplaceAllString(configContent, "state MASTER")
	configContent = regexp.MustCompile(`priority\s+(\d+)`).ReplaceAllStringFunc(configContent, func(s string) string {
		var priority int
		fmt.Sscanf(s, "priority %d", &priority)
		return fmt.Sprintf("priority %d", priority+50) // Increase priority to become master
	})

	// Write updated config
	writeCmd := fmt.Sprintf("echo '%s' | sudo tee %s", configContent, configPath)
	result, err := km.shell.Execute("sh", "-c", writeCmd)
	if err != nil {
		return fmt.Errorf("failed to write keepalived config: %s: %w", result.Stderr, err)
	}

	// Restart keepalived
	result, err = km.shell.Execute("sudo", "systemctl", "restart", "keepalived")
	if err != nil {
		return fmt.Errorf("failed to restart keepalived: %s: %w", result.Stderr, err)
	}

	logger.Info("Keepalived promoted to MASTER", zap.String("id", vipID))
	return nil
}

// DemoteToBackup demotes this node to BACKUP (decrease priority)
func (km *KeepalivedManager) DemoteToBackup(vipID string) error {
	if !km.enabled {
		return fmt.Errorf("Keepalived is not enabled")
	}

	// Read current config
	configPath := "/etc/keepalived/keepalived.conf"
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read keepalived config: %w", err)
	}

	configContent := string(data)

	// Change state to BACKUP and decrease priority
	configContent = regexp.MustCompile(`state\s+\S+`).ReplaceAllString(configContent, "state BACKUP")
	configContent = regexp.MustCompile(`priority\s+(\d+)`).ReplaceAllStringFunc(configContent, func(s string) string {
		var priority int
		fmt.Sscanf(s, "priority %d", &priority)
		newPriority := priority - 50
		if newPriority < 1 {
			newPriority = 1
		}
		return fmt.Sprintf("priority %d", newPriority)
	})

	// Write updated config
	writeCmd := fmt.Sprintf("echo '%s' | sudo tee %s", configContent, configPath)
	result, err := km.shell.Execute("sh", "-c", writeCmd)
	if err != nil {
		return fmt.Errorf("failed to write keepalived config: %s: %w", result.Stderr, err)
	}

	// Restart keepalived
	result, err = km.shell.Execute("sudo", "systemctl", "restart", "keepalived")
	if err != nil {
		return fmt.Errorf("failed to restart keepalived: %s: %w", result.Stderr, err)
	}

	logger.Info("Keepalived demoted to BACKUP", zap.String("id", vipID))
	return nil
}

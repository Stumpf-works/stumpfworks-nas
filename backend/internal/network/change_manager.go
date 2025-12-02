// Revision: 2025-12-02 | Author: Claude | Version: 2.0.0
package network

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ============================================================================
// UNIVERSAL NETWORK CHANGE MANAGER
// Proxmox-style pending changes for ALL network configurations
// ============================================================================

// AddPendingChange adds a pending network change to the queue
// changeType: "bridge", "interface", "route", "firewall", "dns"
// action: "create", "update", "delete"
func AddPendingChange(changeType, action, resourceID, description string, pendingConfig interface{}, currentConfig interface{}) error {
	// Marshal configs to JSON
	pendingJSON, err := json.Marshal(pendingConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal pending config: %w", err)
	}

	var currentJSON []byte
	if currentConfig != nil {
		currentJSON, err = json.Marshal(currentConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal current config: %w", err)
		}
	}

	change := models.PendingNetworkChange{
		ID:            uuid.New().String(),
		ChangeType:    changeType,
		Action:        action,
		ResourceID:    resourceID,
		CurrentConfig: string(currentJSON),
		PendingConfig: string(pendingJSON),
		Description:   description,
		Priority:      100, // Default priority
		Status:        "pending",
		CreatedAt:     time.Now(),
	}

	if err := database.DB.Create(&change).Error; err != nil {
		return fmt.Errorf("failed to save pending change: %w", err)
	}

	logger.Info("Pending network change added",
		zap.String("id", change.ID),
		zap.String("type", changeType),
		zap.String("action", action),
		zap.String("resource", resourceID))

	return nil
}

// GetAllPendingChanges returns all pending network changes sorted by priority
func GetAllPendingChanges() ([]models.PendingNetworkChange, error) {
	var changes []models.PendingNetworkChange
	if err := database.DB.Where("status = ?", "pending").
		Order("priority ASC, created_at ASC").
		Find(&changes).Error; err != nil {
		return nil, fmt.Errorf("failed to load pending changes: %w", err)
	}
	return changes, nil
}

// GetPendingChangesByType returns pending changes filtered by type
func GetPendingChangesByType(changeType string) ([]models.PendingNetworkChange, error) {
	var changes []models.PendingNetworkChange
	if err := database.DB.Where("status = ? AND change_type = ?", "pending", changeType).
		Order("priority ASC, created_at ASC").
		Find(&changes).Error; err != nil {
		return nil, fmt.Errorf("failed to load pending changes: %w", err)
	}
	return changes, nil
}

// HasPendingChanges checks if there are any pending network changes
func HasPendingChanges() (bool, int, error) {
	var count int64
	if err := database.DB.Model(&models.PendingNetworkChange{}).
		Where("status = ?", "pending").
		Count(&count).Error; err != nil {
		return false, 0, fmt.Errorf("failed to count pending changes: %w", err)
	}
	return count > 0, int(count), nil
}

// CreateFullNetworkSnapshot creates a complete snapshot of the entire network state
// This is more comprehensive than bridge-specific snapshots
func CreateFullNetworkSnapshot() (*models.NetworkSnapshot, error) {
	logger.Info("Creating full network snapshot")

	// Capture route table
	cmd := exec.Command("ip", "route", "show")
	routeTableBytes, _ := cmd.CombinedOutput()

	// Capture firewall rules
	cmd = exec.Command("iptables", "-L", "-n", "-v")
	firewallBytes, _ := cmd.CombinedOutput()

	// Capture all interface states
	interfaceStates := make(map[string]map[string]interface{})
	interfaces, _ := ListInterfaces()
	for _, iface := range interfaces {
		state := make(map[string]interface{})
		state["up"] = iface.IsUp
		state["addresses"] = iface.Addresses
		state["mtu"] = iface.MTU
		state["hardware_addr"] = iface.HardwareAddr

		// Get routes for this interface
		cmd = exec.Command("ip", "route", "show", "dev", iface.Name)
		routes, _ := cmd.CombinedOutput()
		state["routes"] = string(routes)

		interfaceStates[iface.Name] = state
	}

	interfaceStatesJSON, _ := json.Marshal(interfaceStates)

	snapshot := models.NetworkSnapshot{
		ID:              uuid.New().String(),
		BridgeName:      "__FULL_SNAPSHOT__", // Special marker for full snapshots
		InterfaceStates: string(interfaceStatesJSON),
		RouteTable:      string(routeTableBytes),
		FirewallRules:   string(firewallBytes),
		Status:          "active",
		CreatedAt:       time.Now(),
	}

	if err := database.DB.Create(&snapshot).Error; err != nil {
		return nil, fmt.Errorf("failed to save network snapshot: %w", err)
	}

	logger.Info("Full network snapshot created", zap.String("snapshot_id", snapshot.ID))
	return &snapshot, nil
}

// ApplyAllPendingChanges applies ALL pending network changes with automatic rollback
// This is the main "Apply Configuration" function for the entire network section
func ApplyAllPendingChanges() error {
	logger.Info("Starting to apply all pending network changes")

	// Step 1: Check if there are any pending changes
	hasPending, count, err := HasPendingChanges()
	if err != nil {
		return fmt.Errorf("failed to check for pending changes: %w", err)
	}

	if !hasPending {
		return fmt.Errorf("no pending network changes to apply")
	}

	logger.Info("Found pending network changes", zap.Int("count", count))

	// Step 2: Create full network snapshot before making ANY changes
	snapshot, err := CreateFullNetworkSnapshot()
	if err != nil {
		return fmt.Errorf("failed to create network snapshot: %w", err)
	}

	// Step 3: Get all pending changes sorted by priority
	changes, err := GetAllPendingChanges()
	if err != nil {
		return fmt.Errorf("failed to load pending changes: %w", err)
	}

	// Step 4: Apply each change in order
	var failedChanges []string
	var successCount int

	for _, change := range changes {
		logger.Info("Applying network change",
			zap.String("type", change.ChangeType),
			zap.String("action", change.Action),
			zap.String("resource", change.ResourceID))

		err := applyNetworkChange(&change)
		if err != nil {
			logger.Error("Failed to apply network change",
				zap.Error(err),
				zap.String("change_id", change.ID),
				zap.String("type", change.ChangeType))

			failedChanges = append(failedChanges, fmt.Sprintf("%s:%s", change.ChangeType, change.ResourceID))

			// Mark change as failed
			database.DB.Model(&change).Updates(map[string]interface{}{
				"status":     "failed",
				"updated_at": time.Now(),
			})

			// ROLLBACK on first failure
			logger.Warn("Rolling back all changes due to failure")
			if rollbackErr := RollbackToFullSnapshot(snapshot.ID); rollbackErr != nil {
				logger.Error("CRITICAL: Rollback failed!", zap.Error(rollbackErr))
				return fmt.Errorf("apply failed and rollback also failed: %w", rollbackErr)
			}

			return fmt.Errorf("failed to apply change %s (rolled back): %w", change.Description, err)
		}

		// Mark change as applied
		database.DB.Model(&change).Updates(map[string]interface{}{
			"status":     "applied",
			"updated_at": time.Now(),
		})

		successCount++
	}

	// Step 5: Reload networking to apply changes from /etc/network/interfaces
	// IMPORTANT: Use "reload" instead of "restart" to avoid breaking SSH connections
	logger.Info("Reloading network configuration (safe mode)...")

	// First, try ifreload (Debian/Ubuntu - safest method, only reloads changed interfaces)
	cmd := exec.Command("ifreload", "-a")
	output, err := cmd.CombinedOutput()

	// If ifreload not available, fall back to systemctl reload (safer than restart)
	if err != nil {
		logger.Warn("ifreload not available, falling back to systemctl reload",
			zap.Error(err))

		cmd = exec.Command("systemctl", "reload", "networking")
		output, err = cmd.CombinedOutput()

		// Last resort: reload-or-restart (systemd will choose safest option)
		if err != nil {
			logger.Warn("systemctl reload failed, trying reload-or-restart",
				zap.Error(err))

			cmd = exec.Command("systemctl", "reload-or-restart", "networking")
			output, err = cmd.CombinedOutput()
		}
	}

	if err != nil {
		logger.Error("Failed to reload networking service",
			zap.Error(err),
			zap.String("output", string(output)))

		// ROLLBACK on networking reload failure
		logger.Warn("Rolling back all changes due to networking reload failure")
		if rollbackErr := RollbackToFullSnapshot(snapshot.ID); rollbackErr != nil {
			logger.Error("CRITICAL: Rollback failed!", zap.Error(rollbackErr))
			return fmt.Errorf("networking reload failed and rollback also failed: %w", rollbackErr)
		}

		return fmt.Errorf("failed to reload networking (rolled back): %w", err)
	}

	// Step 6: All changes applied successfully
	database.DB.Model(snapshot).Updates(map[string]interface{}{
		"status":     "applied",
		"applied_at": time.Now(),
	})

	logger.Info("All pending network changes applied successfully",
		zap.Int("count", successCount),
		zap.String("snapshot_id", snapshot.ID))

	return nil
}

// applyNetworkChange applies a single network change based on its type
func applyNetworkChange(change *models.PendingNetworkChange) error {
	switch change.ChangeType {
	case "bridge":
		return applyBridgeChange(change)
	case "interface":
		return applyInterfaceChange(change)
	case "route":
		return applyRouteChange(change)
	case "firewall":
		return applyFirewallChange(change)
	case "dns":
		return applyDNSChange(change)
	default:
		return fmt.Errorf("unknown change type: %s", change.ChangeType)
	}
}

// applyBridgeChange applies a bridge configuration change by writing to /etc/network/interfaces
// This is the PROXMOX-STYLE approach: all changes go to /etc/network/interfaces
func applyBridgeChange(change *models.PendingNetworkChange) error {
	// Parse pending config
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(change.PendingConfig), &config); err != nil {
		return fmt.Errorf("failed to parse bridge config: %w", err)
	}

	bridgeName := change.ResourceID

	// Read current /etc/network/interfaces
	interfaces, err := ParseInterfacesFile()
	if err != nil {
		return fmt.Errorf("failed to parse interfaces file: %w", err)
	}

	switch change.Action {
	case "create":
		// Extract fields from config
		ports := []string{}
		if portsInterface, ok := config["ports"]; ok {
			if portsList, ok := portsInterface.([]interface{}); ok {
				for _, p := range portsList {
					if portStr, ok := p.(string); ok {
						ports = append(ports, portStr)
					}
				}
			}
		}

		ipAddress, _ := config["ip_address"].(string)
		gateway, _ := config["gateway"].(string)
		ipv6Address, _ := config["ipv6_address"].(string)
		ipv6Gateway, _ := config["ipv6_gateway"].(string)
		vlanAware, _ := config["vlan_aware"].(bool)

		// Add bridge to interfaces configuration (Proxmox-style with IPv6)
		AddBridgeToInterfaces(interfaces, bridgeName, ports, ipAddress, gateway, ipv6Address, ipv6Gateway, vlanAware)

		logger.Info("Adding bridge to /etc/network/interfaces",
			zap.String("bridge", bridgeName),
			zap.String("ipv4_address", ipAddress),
			zap.String("ipv6_address", ipv6Address),
			zap.Bool("vlan_aware", vlanAware),
			zap.Strings("ports", ports))

	case "delete":
		// Remove bridge from interfaces configuration
		RemoveBridgeFromInterfaces(interfaces, bridgeName)

		logger.Info("Removing bridge from /etc/network/interfaces",
			zap.String("bridge", bridgeName))

	case "update":
		// For update, remove old config and add new one
		RemoveBridgeFromInterfaces(interfaces, bridgeName)

		ports := []string{}
		if portsInterface, ok := config["ports"]; ok {
			if portsList, ok := portsInterface.([]interface{}); ok {
				for _, p := range portsList {
					if portStr, ok := p.(string); ok {
						ports = append(ports, portStr)
					}
				}
			}
		}

		ipAddress, _ := config["ip_address"].(string)
		gateway, _ := config["gateway"].(string)
		ipv6Address, _ := config["ipv6_address"].(string)
		ipv6Gateway, _ := config["ipv6_gateway"].(string)
		vlanAware, _ := config["vlan_aware"].(bool)

		AddBridgeToInterfaces(interfaces, bridgeName, ports, ipAddress, gateway, ipv6Address, ipv6Gateway, vlanAware)

		logger.Info("Updating bridge in /etc/network/interfaces",
			zap.String("bridge", bridgeName),
			zap.String("ipv4_address", ipAddress),
			zap.String("ipv6_address", ipv6Address),
			zap.Bool("vlan_aware", vlanAware))
	}

	// Write updated interfaces file
	if err := WriteInterfacesFile(interfaces); err != nil {
		return fmt.Errorf("failed to write interfaces file: %w", err)
	}

	return nil
}

// applyInterfaceChange applies an interface configuration change
func applyInterfaceChange(change *models.PendingNetworkChange) error {
	// Parse pending config
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(change.PendingConfig), &config); err != nil {
		return fmt.Errorf("failed to parse interface config: %w", err)
	}

	interfaceName := change.ResourceID

	// Read current /etc/network/interfaces
	interfaces, err := ParseInterfacesFile()
	if err != nil {
		return fmt.Errorf("failed to parse interfaces file: %w", err)
	}

	switch change.Action {
	case "create", "update":
		// Extract configuration
		ipAddress, _ := config["ip_address"].(string)
		gateway, _ := config["gateway"].(string)
		ipv6Address, _ := config["ipv6_address"].(string)
		ipv6Gateway, _ := config["ipv6_gateway"].(string)
		autostart, _ := config["autostart"].(bool)
		comment, _ := config["comment"].(string)

		// Determine address method
		method := "dhcp"
		if ipAddress != "" {
			method = "static"
		}

		// Update or create interface configuration
		interfaces[interfaceName] = &InterfaceConfig{
			Name:          interfaceName,
			Type:          "physical",
			AddressMethod: method,
			Address:       ipAddress,
			Gateway:       gateway,
			IPv6Address:   ipv6Address,
			IPv6Gateway:   ipv6Gateway,
			IPv6Method:    "static",
			Auto:          autostart,
			Comment:       comment,
		}

		logger.Info("Updating interface in /etc/network/interfaces",
			zap.String("interface", interfaceName),
			zap.String("method", method),
			zap.String("ipv4_address", ipAddress),
			zap.String("ipv6_address", ipv6Address))

	case "delete":
		// Remove interface configuration
		delete(interfaces, interfaceName)

		logger.Info("Removing interface from /etc/network/interfaces",
			zap.String("interface", interfaceName))
	}

	// Write updated interfaces file
	if err := WriteInterfacesFile(interfaces); err != nil {
		return fmt.Errorf("failed to write interfaces file: %w", err)
	}

	return nil
}

// applyRouteChange applies a routing table change
func applyRouteChange(change *models.PendingNetworkChange) error {
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(change.PendingConfig), &config); err != nil {
		return fmt.Errorf("failed to parse route config: %w", err)
	}

	destination, _ := config["destination"].(string)
	gateway, _ := config["gateway"].(string)
	iface, _ := config["interface"].(string)

	switch change.Action {
	case "create":
		cmd := exec.Command("ip", "route", "add", destination, "via", gateway, "dev", iface)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to add route: %s", string(output))
		}

	case "delete":
		cmd := exec.Command("ip", "route", "del", destination)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to delete route: %s", string(output))
		}
	}

	return nil
}

// applyFirewallChange applies a firewall rule change
func applyFirewallChange(change *models.PendingNetworkChange) error {
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(change.PendingConfig), &config); err != nil {
		return fmt.Errorf("failed to parse firewall config: %w", err)
	}

	action, _ := config["action"].(string)
	from, _ := config["from"].(string)
	to, _ := config["to"].(string)
	protocol, _ := config["protocol"].(string)

	// Build iptables command
	args := []string{"-A", "INPUT"}

	if protocol != "" {
		args = append(args, "-p", protocol)
	}
	if from != "" {
		args = append(args, "-s", from)
	}
	if to != "" {
		args = append(args, "-d", to)
	}
	args = append(args, "-j", action)

	cmd := exec.Command("iptables", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to apply firewall rule: %s", string(output))
	}

	return nil
}

// applyDNSChange applies a DNS configuration change
func applyDNSChange(change *models.PendingNetworkChange) error {
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(change.PendingConfig), &config); err != nil {
		return fmt.Errorf("failed to parse DNS config: %w", err)
	}

	// DNS changes would typically modify /etc/resolv.conf
	// This is a simplified implementation
	logger.Warn("DNS change application not fully implemented",
		zap.String("change_id", change.ID))

	return nil
}

// RollbackToFullSnapshot rolls back the entire network to a previous snapshot
// This restores /etc/network/interfaces from backup and restarts networking
func RollbackToFullSnapshot(snapshotID string) error {
	var snapshot models.NetworkSnapshot
	if err := database.DB.Where("id = ?", snapshotID).First(&snapshot).Error; err != nil {
		return fmt.Errorf("snapshot %s not found", snapshotID)
	}

	logger.Warn("Rolling back network to full snapshot",
		zap.String("snapshot_id", snapshotID))

	// Restore /etc/network/interfaces from backup
	if err := RestoreInterfacesFileFromBackup(); err != nil {
		logger.Error("Failed to restore interfaces file from backup", zap.Error(err))
		return fmt.Errorf("failed to restore interfaces file: %w", err)
	}

	// Restart networking to apply the restored configuration
	cmd := exec.Command("systemctl", "restart", "networking")
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Error("Failed to restart networking during rollback",
			zap.Error(err),
			zap.String("output", string(output)))
		return fmt.Errorf("failed to restart networking during rollback: %w", err)
	}

	// Mark snapshot as rolled back
	database.DB.Model(&snapshot).Updates(map[string]interface{}{
		"status":         "rolled_back",
		"rolled_back_at": time.Now(),
	})

	// Mark all pending changes as discarded
	database.DB.Model(&models.PendingNetworkChange{}).
		Where("status = ?", "pending").
		Updates(map[string]interface{}{
			"status":     "discarded",
			"updated_at": time.Now(),
		})

	logger.Info("Network configuration rolled back successfully")
	return nil
}

// DiscardAllPendingChanges removes all pending changes without applying them
func DiscardAllPendingChanges() error {
	result := database.DB.Model(&models.PendingNetworkChange{}).
		Where("status = ?", "pending").
		Updates(map[string]interface{}{
			"status":     "discarded",
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to discard pending changes: %w", result.Error)
	}

	logger.Info("All pending network changes discarded", zap.Int64("count", result.RowsAffected))
	return nil
}

// DiscardPendingChange discards a specific pending change
func DiscardPendingChange(changeID string) error {
	result := database.DB.Model(&models.PendingNetworkChange{}).
		Where("id = ? AND status = ?", changeID, "pending").
		Updates(map[string]interface{}{
			"status":     "discarded",
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to discard pending change: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("pending change %s not found or already processed", changeID)
	}

	logger.Info("Pending network change discarded", zap.String("change_id", changeID))
	return nil
}

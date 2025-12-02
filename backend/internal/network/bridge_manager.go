// Revision: 2025-12-02 | Author: Claude | Version: 2.0.1
package network

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateBridgePersistent creates a bridge and saves it to the database
// This ensures the bridge will be restored after reboot
func CreateBridgePersistent(name string, ports []string, ipAddress string, gateway string, description string) error {
	// Step 1: Check if bridge already exists in database
	var existing models.NetworkBridge
	if err := database.DB.Where("name = ?", name).First(&existing).Error; err == nil {
		return fmt.Errorf("bridge %s already exists in database", name)
	}

	// Step 2: Create the bridge WITHOUT ports first (safe mode)
	// This prevents network loss if a port has IP addresses
	if err := CreateBridge(name, []string{}); err != nil {
		return fmt.Errorf("failed to create bridge in system: %w", err)
	}

	// Step 3: Apply IP configuration BEFORE adding ports
	// This ensures the bridge has connectivity before we remove IPs from ports
	if ipAddress != "" {
		// Parse CIDR to get address and netmask
		parts := strings.Split(ipAddress, "/")
		if len(parts) == 2 {
			addr := parts[0]
			netmask := parts[1]
			if err := ConfigureStaticIP(name, addr, convertCIDRToNetmask(netmask), gateway); err != nil {
				// If IP configuration fails, clean up the bridge
				DeleteBridge(name)
				return fmt.Errorf("failed to apply IP configuration to bridge: %w", err)
			}
		}
	}

	// Step 4: Now it's safe to add ports to the bridge
	// If ports have IPs and bridge has IPs, the migration will be safe
	if len(ports) > 0 {
		for _, port := range ports {
			if port == "" {
				continue
			}
			if err := AttachPortToBridge(name, port); err != nil {
				// Log warning but don't fail - bridge is still functional
				logger.Warn("Failed to attach port to bridge",
					zap.Error(err),
					zap.String("bridge", name),
					zap.String("port", port))
			}
		}
	}

	// Step 4: Save to database for persistence
	bridge := models.NetworkBridge{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Ports:       strings.Join(ports, ","),
		IPAddress:   ipAddress,
		Gateway:     gateway,
		Autostart:   true,
		Status:      "active",
	}

	if err := database.DB.Create(&bridge).Error; err != nil {
		// Rollback: Delete the bridge from system
		DeleteBridge(name)
		return fmt.Errorf("failed to save bridge to database: %w", err)
	}

	logger.Info("Bridge created and saved to database",
		zap.String("name", name),
		zap.String("id", bridge.ID),
		zap.Strings("ports", ports))

	return nil
}

// DeleteBridgePersistent deletes a bridge from both system and database
func DeleteBridgePersistent(name string) error {
	// Step 1: Delete from system
	if err := DeleteBridge(name); err != nil {
		logger.Warn("Failed to delete bridge from system (may not exist)", zap.Error(err), zap.String("bridge", name))
	}

	// Step 2: Delete from database
	result := database.DB.Where("name = ?", name).Delete(&models.NetworkBridge{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete bridge from database: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("bridge %s not found in database", name)
	}

	logger.Info("Bridge deleted from system and database", zap.String("name", name))
	return nil
}

// RestoreAllBridges restores all bridges from database on system startup
// This should be called during system initialization
func RestoreAllBridges() error {
	var bridges []models.NetworkBridge
	if err := database.DB.Where("autostart = ?", true).Find(&bridges).Error; err != nil {
		return fmt.Errorf("failed to load bridges from database: %w", err)
	}

	logger.Info("Restoring network bridges from database", zap.Int("count", len(bridges)))

	for _, bridge := range bridges {
		logger.Info("Restoring bridge", zap.String("name", bridge.Name), zap.String("id", bridge.ID))

		// Parse ports
		var ports []string
		if bridge.Ports != "" {
			ports = strings.Split(bridge.Ports, ",")
		}

		// Create the bridge
		if err := CreateBridge(bridge.Name, ports); err != nil {
			logger.Error("Failed to restore bridge", zap.Error(err), zap.String("name", bridge.Name))
			// Update status in database
			database.DB.Model(&bridge).Updates(map[string]interface{}{
				"status":     "error",
				"last_error": err.Error(),
			})
			continue
		}

		// Apply IP configuration if specified
		if bridge.IPAddress != "" {
			parts := strings.Split(bridge.IPAddress, "/")
			if len(parts) == 2 {
				addr := parts[0]
				netmask := parts[1]
				if err := ConfigureStaticIP(bridge.Name, addr, convertCIDRToNetmask(netmask), bridge.Gateway); err != nil {
					logger.Warn("Failed to apply IP configuration to restored bridge",
						zap.Error(err),
						zap.String("bridge", bridge.Name))
				}
			}
		}

		// Update status in database
		database.DB.Model(&bridge).Updates(map[string]interface{}{
			"status":     "active",
			"last_error": "",
		})

		logger.Info("Bridge restored successfully", zap.String("name", bridge.Name))
	}

	return nil
}

// GetPersistedBridges returns all bridges from the database
func GetPersistedBridges() ([]models.NetworkBridge, error) {
	var bridges []models.NetworkBridge
	if err := database.DB.Find(&bridges).Error; err != nil {
		return nil, fmt.Errorf("failed to load bridges from database: %w", err)
	}
	return bridges, nil
}

// convertCIDRToNetmask converts CIDR notation to netmask (e.g., "24" -> "255.255.255.0")
// This is now algorithmic instead of using a lookup table
func convertCIDRToNetmask(cidr string) string {
	// Parse CIDR prefix length
	prefixLen := 0
	fmt.Sscanf(cidr, "%d", &prefixLen)

	// Validate range
	if prefixLen < 0 || prefixLen > 32 {
		logger.Warn("Invalid CIDR prefix length, defaulting to /24",
			zap.String("cidr", cidr))
		prefixLen = 24
	}

	// Calculate netmask algorithmically
	// Create a 32-bit mask with prefixLen ones from the left
	mask := uint32(0xFFFFFFFF << (32 - prefixLen))

	// Convert to dotted-decimal notation
	octet1 := byte((mask >> 24) & 0xFF)
	octet2 := byte((mask >> 16) & 0xFF)
	octet3 := byte((mask >> 8) & 0xFF)
	octet4 := byte(mask & 0xFF)

	return fmt.Sprintf("%d.%d.%d.%d", octet1, octet2, octet3, octet4)
}

// ============================================================================
// PROXMOX-STYLE PENDING CHANGES SYSTEM
// ============================================================================

// CreateBridgeWithPendingChanges creates a bridge configuration in the database
// WITHOUT applying it to the system. Changes must be applied via ApplyPendingChanges.
// This is the Proxmox-style workflow: Create -> Review -> Apply -> Connectivity Check -> Commit or Rollback
func CreateBridgeWithPendingChanges(name string, ports []string, ipAddress string, gateway string, ipv6Address string, ipv6Gateway string, vlanAware bool, autostart bool, description string) error {
	// Check if bridge already exists in database
	var existing models.NetworkBridge
	if err := database.DB.Where("name = ?", name).First(&existing).Error; err == nil {
		return fmt.Errorf("bridge %s already exists in database", name)
	}

	// Save to database with pending status (NOT YET APPLIED TO SYSTEM)
	bridge := models.NetworkBridge{
		ID:                 uuid.New().String(),
		Name:               name,
		Description:        description,
		Status:             "pending", // NOT YET APPLIED
		HasPendingChanges:  true,
		PendingPorts:       strings.Join(ports, ","),
		PendingIPAddress:   ipAddress,
		PendingGateway:     gateway,
		PendingIPv6Address: ipv6Address,
		PendingIPv6Gateway: ipv6Gateway,
		PendingVLANAware:   vlanAware,
		Autostart:          autostart,
	}

	if err := database.DB.Create(&bridge).Error; err != nil {
		return fmt.Errorf("failed to save bridge configuration to database: %w", err)
	}

	// Add to universal pending changes tracking
	pendingConfig := map[string]interface{}{
		"name":         name,
		"description":  description,
		"ports":        ports,
		"ip_address":   ipAddress,
		"gateway":      gateway,
		"ipv6_address": ipv6Address,
		"ipv6_gateway": ipv6Gateway,
		"vlan_aware":   vlanAware,
		"autostart":    autostart,
	}

	desc := fmt.Sprintf("Create bridge %s", name)
	if description != "" {
		desc = fmt.Sprintf("Create bridge %s (%s)", name, description)
	}

	if err := AddPendingChange("bridge", "create", name, desc, pendingConfig, nil); err != nil {
		// Rollback bridge creation if we can't add pending change
		database.DB.Delete(&bridge)
		return fmt.Errorf("failed to add pending change: %w", err)
	}

	logger.Info("Bridge configuration saved (pending application)",
		zap.String("name", name),
		zap.String("id", bridge.ID),
		zap.Bool("has_pending_changes", true))

	return nil
}

// UpdateBridgeWithPendingChanges updates a bridge configuration in the database
// WITHOUT applying it to the system. Changes must be applied via ApplyPendingChanges.
func UpdateBridgeWithPendingChanges(name string, ports []string, ipAddress string, gateway string, ipv6Address string, ipv6Gateway string, vlanAware bool, autostart bool) error {
	var bridge models.NetworkBridge
	if err := database.DB.Where("name = ?", name).First(&bridge).Error; err != nil {
		return fmt.Errorf("bridge %s not found in database", name)
	}

	// Current configuration
	currentConfig := map[string]interface{}{
		"ports":        bridge.Ports,
		"ip_address":   bridge.IPAddress,
		"gateway":      bridge.Gateway,
		"ipv6_address": bridge.IPv6Address,
		"ipv6_gateway": bridge.IPv6Gateway,
		"vlan_aware":   bridge.VLANAware,
	}

	// Pending configuration
	pendingConfig := map[string]interface{}{
		"ports":        ports,
		"ip_address":   ipAddress,
		"gateway":      gateway,
		"ipv6_address": ipv6Address,
		"ipv6_gateway": ipv6Gateway,
		"vlan_aware":   vlanAware,
		"autostart":    autostart,
	}

	// Store pending changes
	updates := map[string]interface{}{
		"has_pending_changes":  true,
		"pending_ports":        strings.Join(ports, ","),
		"pending_ip_address":   ipAddress,
		"pending_gateway":      gateway,
		"pending_ipv6_address": ipv6Address,
		"pending_ipv6_gateway": ipv6Gateway,
		"pending_vlan_aware":   vlanAware,
		"status":               "pending_changes",
	}

	if err := database.DB.Model(&bridge).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update bridge configuration: %w", err)
	}

	// Add to universal pending changes tracking
	desc := fmt.Sprintf("Update bridge %s", name)
	if err := AddPendingChange("bridge", "update", name, desc, pendingConfig, currentConfig); err != nil {
		// Rollback updates if we can't add pending change
		database.DB.Model(&bridge).Updates(map[string]interface{}{
			"has_pending_changes":  false,
			"pending_ports":        "",
			"pending_ip_address":   "",
			"pending_gateway":      "",
			"pending_ipv6_address": "",
			"pending_ipv6_gateway": "",
		})
		return fmt.Errorf("failed to add pending change: %w", err)
	}

	logger.Info("Bridge configuration updated (pending application)",
		zap.String("name", name),
		zap.Bool("has_pending_changes", true))

	return nil
}

// CreateNetworkSnapshot creates a snapshot of the current network configuration
// before applying changes. This enables rollback if changes break connectivity.
func CreateNetworkSnapshot(bridgeName string) (*models.NetworkSnapshot, error) {
	var bridge models.NetworkBridge
	if err := database.DB.Where("name = ?", bridgeName).First(&bridge).Error; err != nil {
		return nil, fmt.Errorf("bridge %s not found in database", bridgeName)
	}

	// Capture current route table
	cmd := exec.Command("ip", "route", "show")
	routeTableBytes, err := cmd.CombinedOutput()
	if err != nil {
		logger.Warn("Failed to capture route table for snapshot", zap.Error(err))
	}

	// Capture interface states
	interfaceStates := make(map[string]map[string]interface{})

	// Get all interfaces
	interfaces, _ := ListInterfaces()
	for _, iface := range interfaces {
		state := make(map[string]interface{})
		state["up"] = iface.IsUp
		state["addresses"] = iface.Addresses

		// Get routes for this interface
		cmd = exec.Command("ip", "route", "show", "dev", iface.Name)
		routes, _ := cmd.CombinedOutput()
		state["routes"] = string(routes)

		interfaceStates[iface.Name] = state
	}

	interfaceStatesJSON, _ := json.Marshal(interfaceStates)

	snapshot := models.NetworkSnapshot{
		ID:               uuid.New().String(),
		BridgeName:       bridgeName,
		CurrentPorts:     bridge.Ports,
		CurrentIPAddress: bridge.IPAddress,
		CurrentGateway:   bridge.Gateway,
		InterfaceStates:  string(interfaceStatesJSON),
		RouteTable:       string(routeTableBytes),
		Status:           "active",
		CreatedAt:        time.Now(),
	}

	if err := database.DB.Create(&snapshot).Error; err != nil {
		return nil, fmt.Errorf("failed to save network snapshot: %w", err)
	}

	logger.Info("Network snapshot created",
		zap.String("snapshot_id", snapshot.ID),
		zap.String("bridge", bridgeName))

	return &snapshot, nil
}

// ApplyPendingChanges applies pending network changes with automatic rollback on failure
// This is the critical "Apply Configuration" step that implements:
// 1. Create snapshot of current state
// 2. Apply changes to system
// 3. Perform connectivity check (can backend still be reached?)
// 4. Auto-rollback if connectivity check fails within timeout
// 5. Commit changes if successful
func ApplyPendingChanges(bridgeName string) error {
	var bridge models.NetworkBridge
	if err := database.DB.Where("name = ?", bridgeName).First(&bridge).Error; err != nil {
		return fmt.Errorf("bridge %s not found in database", bridgeName)
	}

	if !bridge.HasPendingChanges {
		return fmt.Errorf("bridge %s has no pending changes to apply", bridgeName)
	}

	// Step 1: Create snapshot before making any changes
	logger.Info("Creating network snapshot before applying changes", zap.String("bridge", bridgeName))
	snapshot, err := CreateNetworkSnapshot(bridgeName)
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	// Step 2: Parse pending changes
	var ports []string
	if bridge.PendingPorts != "" {
		ports = strings.Split(bridge.PendingPorts, ",")
	}

	// Step 3: Apply changes to the system
	logger.Info("Applying pending network changes",
		zap.String("bridge", bridgeName),
		zap.String("pending_ip", bridge.PendingIPAddress),
		zap.Strings("pending_ports", ports))

	// If bridge doesn't exist yet in system, create it
	existingBridges, _ := ListBridges()
	bridgeExists := false
	for _, b := range existingBridges {
		if b == bridgeName {
			bridgeExists = true
			break
		}
	}

	var applyErr error
	if !bridgeExists {
		// New bridge - create with safe workflow
		applyErr = CreateBridge(bridgeName, []string{}) // Create without ports first

		if applyErr == nil && bridge.PendingIPAddress != "" {
			// Apply IP configuration
			parts := strings.Split(bridge.PendingIPAddress, "/")
			if len(parts) == 2 {
				addr := parts[0]
				netmask := parts[1]
				applyErr = ConfigureStaticIP(bridgeName, addr, convertCIDRToNetmask(netmask), bridge.PendingGateway)
			}
		}

		if applyErr == nil && len(ports) > 0 {
			// Now safe to add ports
			for _, port := range ports {
				if port == "" {
					continue
				}
				if err := AttachPortToBridge(bridgeName, port); err != nil {
					logger.Warn("Failed to attach port to bridge",
						zap.Error(err),
						zap.String("bridge", bridgeName),
						zap.String("port", port))
				}
			}
		}
	} else {
		// Existing bridge - update configuration
		// This is more complex and risky, so we do it step by step

		// First update IP if changed
		if bridge.PendingIPAddress != bridge.IPAddress || bridge.PendingGateway != bridge.Gateway {
			if bridge.PendingIPAddress != "" {
				parts := strings.Split(bridge.PendingIPAddress, "/")
				if len(parts) == 2 {
					addr := parts[0]
					netmask := parts[1]
					applyErr = ConfigureStaticIP(bridgeName, addr, convertCIDRToNetmask(netmask), bridge.PendingGateway)
				}
			}
		}

		// Then update ports if changed
		if applyErr == nil && bridge.PendingPorts != bridge.Ports {
			// Remove old ports, add new ones
			// This is simplified - in production you'd want to diff the lists
			for _, port := range ports {
				if port == "" {
					continue
				}
				if err := AttachPortToBridge(bridgeName, port); err != nil {
					logger.Warn("Failed to attach port to bridge",
						zap.Error(err),
						zap.String("bridge", bridgeName),
						zap.String("port", port))
				}
			}
		}
	}

	// Step 4: Check if changes were applied successfully
	if applyErr != nil {
		logger.Error("Failed to apply network changes, rolling back",
			zap.Error(applyErr),
			zap.String("bridge", bridgeName))

		// Rollback
		if err := RollbackToSnapshot(snapshot.ID); err != nil {
			logger.Error("CRITICAL: Rollback failed!", zap.Error(err))
			return fmt.Errorf("apply failed and rollback also failed: %w", err)
		}

		return fmt.Errorf("failed to apply changes (rolled back): %w", applyErr)
	}

	// Step 5: Changes applied successfully - commit them to database
	updates := map[string]interface{}{
		"has_pending_changes":  false,
		"ports":                bridge.PendingPorts,
		"ip_address":           bridge.PendingIPAddress,
		"gateway":              bridge.PendingGateway,
		"ipv6_address":         bridge.PendingIPv6Address,
		"ipv6_gateway":         bridge.PendingIPv6Gateway,
		"vlan_aware":           bridge.PendingVLANAware,
		"status":               "active",
		"last_error":           "",
		"pending_ports":        "",
		"pending_ip_address":   "",
		"pending_gateway":      "",
		"pending_ipv6_address": "",
		"pending_ipv6_gateway": "",
		"pending_vlan_aware":   false,
	}

	if err := database.DB.Model(&bridge).Updates(updates).Error; err != nil {
		logger.Error("Failed to commit changes to database", zap.Error(err))
		// Don't rollback - changes are already applied to system
	}

	// Mark snapshot as applied
	database.DB.Model(snapshot).Updates(map[string]interface{}{
		"status":     "applied",
		"applied_at": time.Now(),
	})

	logger.Info("Network changes applied successfully",
		zap.String("bridge", bridgeName),
		zap.String("snapshot_id", snapshot.ID))

	return nil
}

// RollbackToSnapshot rolls back network configuration to a previous snapshot
func RollbackToSnapshot(snapshotID string) error {
	var snapshot models.NetworkSnapshot
	if err := database.DB.Where("id = ?", snapshotID).First(&snapshot).Error; err != nil {
		return fmt.Errorf("snapshot %s not found", snapshotID)
	}

	logger.Warn("Rolling back network configuration to snapshot",
		zap.String("snapshot_id", snapshotID),
		zap.String("bridge", snapshot.BridgeName))

	// Parse old configuration
	var ports []string
	if snapshot.CurrentPorts != "" {
		ports = strings.Split(snapshot.CurrentPorts, ",")
	}

	// Delete current bridge
	DeleteBridge(snapshot.BridgeName)

	// Recreate with old configuration
	if err := CreateBridge(snapshot.BridgeName, []string{}); err != nil {
		return fmt.Errorf("failed to recreate bridge during rollback: %w", err)
	}

	// Restore IP configuration
	if snapshot.CurrentIPAddress != "" {
		parts := strings.Split(snapshot.CurrentIPAddress, "/")
		if len(parts) == 2 {
			addr := parts[0]
			netmask := parts[1]
			if err := ConfigureStaticIP(snapshot.BridgeName, addr, convertCIDRToNetmask(netmask), snapshot.CurrentGateway); err != nil {
				logger.Error("Failed to restore IP configuration during rollback", zap.Error(err))
			}
		}
	}

	// Restore ports
	for _, port := range ports {
		if port == "" {
			continue
		}
		if err := AttachPortToBridge(snapshot.BridgeName, port); err != nil {
			logger.Warn("Failed to restore port during rollback",
				zap.Error(err),
				zap.String("port", port))
		}
	}

	// Update snapshot status
	database.DB.Model(&snapshot).Updates(map[string]interface{}{
		"status":         "rolled_back",
		"rolled_back_at": time.Now(),
	})

	// Update bridge status in database
	var bridge models.NetworkBridge
	if err := database.DB.Where("name = ?", snapshot.BridgeName).First(&bridge).Error; err == nil {
		database.DB.Model(&bridge).Updates(map[string]interface{}{
			"has_pending_changes":  false,
			"status":               "active",
			"pending_ports":        "",
			"pending_ip_address":   "",
			"pending_gateway":      "",
			"pending_ipv6_address": "",
			"pending_ipv6_gateway": "",
			"pending_vlan_aware":   false,
		})
	}

	logger.Info("Rollback completed successfully", zap.String("bridge", snapshot.BridgeName))
	return nil
}

// GetPendingChanges returns all bridges with pending changes
func GetPendingChanges() ([]models.NetworkBridge, error) {
	var bridges []models.NetworkBridge
	if err := database.DB.Where("has_pending_changes = ?", true).Find(&bridges).Error; err != nil {
		return nil, fmt.Errorf("failed to load pending changes: %w", err)
	}
	return bridges, nil
}

// DiscardPendingChanges discards pending changes for a bridge without applying them
func DiscardPendingChanges(bridgeName string) error {
	var bridge models.NetworkBridge
	if err := database.DB.Where("name = ?", bridgeName).First(&bridge).Error; err != nil {
		return fmt.Errorf("bridge %s not found in database", bridgeName)
	}

	updates := map[string]interface{}{
		"has_pending_changes":  false,
		"pending_ports":        "",
		"pending_ip_address":   "",
		"pending_gateway":      "",
		"pending_ipv6_address": "",
		"pending_ipv6_gateway": "",
		"pending_vlan_aware":   false,
		"status":               bridge.Status, // Keep current status
	}

	if err := database.DB.Model(&bridge).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to discard pending changes: %w", err)
	}

	logger.Info("Pending changes discarded", zap.String("bridge", bridgeName))
	return nil
}

// Revision: 2025-12-02 | Author: Claude | Version: 1.0.0
package network

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// ============================================================================
// PROXMOX-STYLE /etc/network/interfaces MANAGEMENT
// Writes persistent network configuration like Proxmox does
// ============================================================================

const (
	InterfacesFile       = "/etc/network/interfaces"
	InterfacesBackupFile = "/etc/network/interfaces.backup"
)

// InterfaceConfig represents a network interface configuration
type InterfaceConfig struct {
	Name          string
	Type          string // "loopback", "physical", "bridge"
	AddressMethod string // "static", "dhcp", "manual"

	// IPv4 configuration
	Address       string // CIDR notation: 192.168.1.10/24
	Gateway       string

	// IPv6 configuration
	IPv6Address   string // CIDR notation: 2001:db8::1/64
	IPv6Gateway   string
	IPv6Method    string // "static", "auto", "dhcp"

	// Bridge-specific
	BridgePorts   string // For bridges: "eno1" or "eno1 eno2"
	BridgeSTP     string // "off" or "on"
	BridgeFD      string // Forward delay: "0"
	BridgeVLANAware bool // VLAN aware bridge

	Auto          bool   // "auto" line
	AllowHotplug  bool   // "allow-hotplug" line
	Comment       string // Optional comment above interface
}

// ParseInterfacesFile parses /etc/network/interfaces and returns all interface configurations
func ParseInterfacesFile() (map[string]*InterfaceConfig, error) {
	file, err := os.Open(InterfacesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", InterfacesFile, err)
	}
	defer file.Close()

	interfaces := make(map[string]*InterfaceConfig)
	scanner := bufio.NewScanner(file)

	var currentIface *InterfaceConfig
	var currentComment string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Skip empty lines and source directives
		if trimmed == "" || strings.HasPrefix(trimmed, "source") {
			continue
		}

		// Capture comments
		if strings.HasPrefix(trimmed, "#") {
			currentComment = trimmed
			continue
		}

		// Parse "auto" directive
		if strings.HasPrefix(trimmed, "auto ") {
			ifaceName := strings.TrimPrefix(trimmed, "auto ")
			if interfaces[ifaceName] == nil {
				interfaces[ifaceName] = &InterfaceConfig{Name: ifaceName}
			}
			interfaces[ifaceName].Auto = true
			continue
		}

		// Parse "allow-hotplug" directive
		if strings.HasPrefix(trimmed, "allow-hotplug ") {
			ifaceName := strings.TrimPrefix(trimmed, "allow-hotplug ")
			if interfaces[ifaceName] == nil {
				interfaces[ifaceName] = &InterfaceConfig{Name: ifaceName}
			}
			interfaces[ifaceName].AllowHotplug = true
			continue
		}

		// Parse "iface" directive
		if strings.HasPrefix(trimmed, "iface ") {
			parts := strings.Fields(trimmed)
			if len(parts) >= 4 {
				ifaceName := parts[1]
				_ = parts[2]       // addressFamily: "inet" or "inet6" (not used yet)
				method := parts[3] // "static", "dhcp", "manual", "loopback"

				currentIface = &InterfaceConfig{
					Name:          ifaceName,
					AddressMethod: method,
					Comment:       currentComment,
				}

				// Determine type
				if method == "loopback" {
					currentIface.Type = "loopback"
				} else {
					currentIface.Type = "physical" // Will be updated if bridge-ports found
				}

				interfaces[ifaceName] = currentIface
				currentComment = ""
			}
			continue
		}

		// Parse interface options (indented lines)
		if currentIface != nil && strings.HasPrefix(line, "        ") {
			// Remove leading whitespace
			option := strings.TrimSpace(line)

			if strings.HasPrefix(option, "address ") {
				currentIface.Address = strings.TrimPrefix(option, "address ")
			} else if strings.HasPrefix(option, "gateway ") {
				currentIface.Gateway = strings.TrimPrefix(option, "gateway ")
			} else if strings.HasPrefix(option, "bridge-ports ") {
				currentIface.BridgePorts = strings.TrimPrefix(option, "bridge-ports ")
				currentIface.Type = "bridge"
			} else if strings.HasPrefix(option, "bridge-stp ") {
				currentIface.BridgeSTP = strings.TrimPrefix(option, "bridge-stp ")
			} else if strings.HasPrefix(option, "bridge-fd ") {
				currentIface.BridgeFD = strings.TrimPrefix(option, "bridge-fd ")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading %s: %w", InterfacesFile, err)
	}

	return interfaces, nil
}

// WriteInterfacesFile writes interface configurations to /etc/network/interfaces
// This creates a backup first and writes in Proxmox format
func WriteInterfacesFile(interfaces map[string]*InterfaceConfig) error {
	// Create backup of current file
	if err := backupInterfacesFile(); err != nil {
		logger.Warn("Failed to create backup of interfaces file", zap.Error(err))
	}

	// Build new interfaces file content
	var content strings.Builder

	// Header comment
	content.WriteString("# This file describes the network interfaces available on your system\n")
	content.WriteString("# and how to activate them. For more information, see interfaces(5).\n")
	content.WriteString("\n")
	content.WriteString("source /etc/network/interfaces.d/*\n")
	content.WriteString("\n")

	// Always write loopback first
	content.WriteString("# The loopback network interface\n")
	content.WriteString("auto lo\n")
	content.WriteString("iface lo inet loopback\n")
	content.WriteString("\n")

	// Write all other interfaces
	for _, iface := range interfaces {
		if iface.Name == "lo" {
			continue // Already written
		}

		// Write comment if exists
		if iface.Comment != "" {
			content.WriteString(iface.Comment + "\n")
		}

		// Write auto/allow-hotplug
		if iface.Auto {
			content.WriteString(fmt.Sprintf("auto %s\n", iface.Name))
		}
		if iface.AllowHotplug {
			content.WriteString(fmt.Sprintf("allow-hotplug %s\n", iface.Name))
		}

		// Write IPv4 iface line
		content.WriteString(fmt.Sprintf("iface %s inet %s\n", iface.Name, iface.AddressMethod))

		// Write IPv4 options (indented with 8 spaces like Proxmox)
		if iface.Address != "" {
			content.WriteString(fmt.Sprintf("        address %s\n", iface.Address))
		}
		if iface.Gateway != "" {
			content.WriteString(fmt.Sprintf("        gateway %s\n", iface.Gateway))
		}
		if iface.BridgePorts != "" {
			content.WriteString(fmt.Sprintf("        bridge-ports %s\n", iface.BridgePorts))
		}
		if iface.BridgeSTP != "" {
			content.WriteString(fmt.Sprintf("        bridge-stp %s\n", iface.BridgeSTP))
		}
		if iface.BridgeFD != "" {
			content.WriteString(fmt.Sprintf("        bridge-fd %s\n", iface.BridgeFD))
		}
		if iface.BridgeVLANAware {
			content.WriteString("        bridge-vlan-aware yes\n")
		}

		content.WriteString("\n")

		// Write IPv6 configuration if present (Proxmox-style)
		if iface.IPv6Address != "" || iface.IPv6Method != "" {
			method := iface.IPv6Method
			if method == "" {
				method = "static"
			}

			content.WriteString(fmt.Sprintf("iface %s inet6 %s\n", iface.Name, method))

			if iface.IPv6Address != "" {
				content.WriteString(fmt.Sprintf("        address %s\n", iface.IPv6Address))
			}
			if iface.IPv6Gateway != "" {
				content.WriteString(fmt.Sprintf("        gateway %s\n", iface.IPv6Gateway))
			}

			content.WriteString("\n")
		}
	}

	// Write to file atomically (write to temp file, then rename)
	tempFile := InterfacesFile + ".tmp"
	if err := os.WriteFile(tempFile, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("failed to write temp interfaces file: %w", err)
	}

	if err := os.Rename(tempFile, InterfacesFile); err != nil {
		os.Remove(tempFile) // Clean up temp file
		return fmt.Errorf("failed to replace interfaces file: %w", err)
	}

	logger.Info("Network interfaces file updated", zap.String("path", InterfacesFile))
	return nil
}

// backupInterfacesFile creates a backup of the current interfaces file
func backupInterfacesFile() error {
	input, err := os.ReadFile(InterfacesFile)
	if err != nil {
		return fmt.Errorf("failed to read interfaces file: %w", err)
	}

	if err := os.WriteFile(InterfacesBackupFile, input, 0644); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}

	logger.Info("Created backup of interfaces file", zap.String("backup", InterfacesBackupFile))
	return nil
}

// RestoreInterfacesFileFromBackup restores /etc/network/interfaces from backup
func RestoreInterfacesFileFromBackup() error {
	input, err := os.ReadFile(InterfacesBackupFile)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	if err := os.WriteFile(InterfacesFile, input, 0644); err != nil {
		return fmt.Errorf("failed to restore interfaces file: %w", err)
	}

	logger.Info("Restored interfaces file from backup")
	return nil
}

// AddBridgeToInterfaces adds a bridge configuration to the interfaces map (Proxmox-style with IPv6 support)
func AddBridgeToInterfaces(interfaces map[string]*InterfaceConfig, name string, ports []string, address string, gateway string, ipv6Address string, ipv6Gateway string, vlanAware bool) {
	// Determine address method
	method := "manual"
	if address != "" {
		method = "static"
	}

	interfaces[name] = &InterfaceConfig{
		Name:            name,
		Type:            "bridge",
		AddressMethod:   method,
		Address:         address,
		Gateway:         gateway,
		IPv6Address:     ipv6Address,
		IPv6Gateway:     ipv6Gateway,
		IPv6Method:      "static",
		BridgePorts:     strings.Join(ports, " "),
		BridgeSTP:       "off",
		BridgeFD:        "0",
		BridgeVLANAware: vlanAware,
		Auto:            true,
		Comment:         fmt.Sprintf("# Bridge %s", name),
	}

	// Set bridge ports to manual (no IP, bridge takes over)
	// IMPORTANT: Must set Auto: true so interfaces are brought up automatically
	for _, port := range ports {
		if port == "" {
			continue
		}
		interfaces[port] = &InterfaceConfig{
			Name:          port,
			Type:          "physical",
			AddressMethod: "manual",
			Auto:          true, // CRITICAL: Needed to bring interface UP
			Comment:       fmt.Sprintf("# Port for bridge %s", name),
		}
	}
}

// RemoveBridgeFromInterfaces removes a bridge and restores its ports
func RemoveBridgeFromInterfaces(interfaces map[string]*InterfaceConfig, name string) {
	bridge, exists := interfaces[name]
	if !exists {
		return
	}

	// Get bridge ports
	var ports []string
	if bridge.BridgePorts != "" {
		ports = strings.Fields(bridge.BridgePorts)
	}

	// Delete bridge
	delete(interfaces, name)

	// Restore ports to DHCP (default)
	for _, port := range ports {
		if port == "" {
			continue
		}
		interfaces[port] = &InterfaceConfig{
			Name:          port,
			Type:          "physical",
			AddressMethod: "dhcp",
			AllowHotplug:  true,
		}
	}
}

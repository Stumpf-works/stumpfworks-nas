// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package network

import (
	"fmt"
	"net"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
)

// InterfaceManager manages network interfaces
type InterfaceManager struct {
	shell   executor.ShellExecutor
	enabled bool
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
	Name         string   `json:"name"`
	HardwareAddr string   `json:"hardware_addr"`
	IPAddresses  []string `json:"ip_addresses"`
	IPv6Addresses []string `json:"ipv6_addresses"`
	MTU          int      `json:"mtu"`
	State        string   `json:"state"` // up, down
	Speed        string   `json:"speed"`
	Duplex       string   `json:"duplex"`
	Type         string   `json:"type"` // ethernet, wireless, bridge, bond
	Master       string   `json:"master"` // for bonded/bridged interfaces
	Slaves       []string `json:"slaves"` // for bond/bridge masters
}

// InterfaceConfig represents interface configuration
type InterfaceConfig struct {
	Name       string `json:"name"`
	Method     string `json:"method"` // dhcp, static
	Address    string `json:"address"`
	Netmask    string `json:"netmask"`
	Gateway    string `json:"gateway"`
	DNS        []string `json:"dns"`
	MTU        int    `json:"mtu"`
}

// BondConfig represents bonding configuration
type BondConfig struct {
	Name       string   `json:"name"`
	Mode       string   `json:"mode"` // balance-rr, active-backup, 802.3ad, etc.
	Slaves     []string `json:"slaves"`
	Primary    string   `json:"primary"`
	MIIMon     int      `json:"miimon"`
}

// NewInterfaceManager creates a new interface manager
func NewInterfaceManager(shell executor.ShellExecutor) (*InterfaceManager, error) {
	return &InterfaceManager{
		shell:   shell,
		enabled: true,
	}, nil
}

// IsEnabled returns whether interface management is available
func (i *InterfaceManager) IsEnabled() bool {
	return i.enabled
}

// ListInterfaces lists all network interfaces
func (i *InterfaceManager) ListInterfaces() ([]NetworkInterface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to list interfaces: %w", err)
	}

	var result []NetworkInterface
	for _, iface := range interfaces {
		netIface := NetworkInterface{
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr.String(),
			MTU:          iface.MTU,
		}

		// Get state
		if iface.Flags&net.FlagUp != 0 {
			netIface.State = "up"
		} else {
			netIface.State = "down"
		}

		// Get IP addresses
		addrs, err := iface.Addrs()
		if err == nil {
			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if !ok {
					continue
				}

				if ipNet.IP.To4() != nil {
					netIface.IPAddresses = append(netIface.IPAddresses, ipNet.IP.String())
				} else {
					netIface.IPv6Addresses = append(netIface.IPv6Addresses, ipNet.IP.String())
				}
			}
		}

		// Get additional details from ethtool/ip command
		i.getInterfaceDetails(&netIface)

		result = append(result, netIface)
	}

	return result, nil
}

// getInterfaceDetails gets additional interface details
func (i *InterfaceManager) getInterfaceDetails(iface *NetworkInterface) {
	// Try to get speed and duplex from ethtool
	if i.shell.CommandExists("ethtool") {
		result, err := i.shell.Execute("ethtool", iface.Name)
		if err == nil {
			lines := strings.Split(result.Stdout, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "Speed:") {
					iface.Speed = strings.TrimSpace(strings.TrimPrefix(line, "Speed:"))
				}
				if strings.HasPrefix(line, "Duplex:") {
					iface.Duplex = strings.TrimSpace(strings.TrimPrefix(line, "Duplex:"))
				}
			}
		}
	}

	// Get interface type and master info from ip command
	result, err := i.shell.Execute("ip", "-details", "link", "show", iface.Name)
	if err == nil {
		output := result.Stdout

		// Determine type
		if strings.Contains(output, "bond") {
			iface.Type = "bond"
		} else if strings.Contains(output, "bridge") {
			iface.Type = "bridge"
		} else if strings.Contains(output, "vlan") {
			iface.Type = "vlan"
		} else if strings.Contains(output, "wlan") || strings.Contains(output, "wifi") {
			iface.Type = "wireless"
		} else {
			iface.Type = "ethernet"
		}

		// Check if slave
		if strings.Contains(output, "master") {
			parts := strings.Split(output, "master")
			if len(parts) > 1 {
				fields := strings.Fields(parts[1])
				if len(fields) > 0 {
					iface.Master = fields[0]
				}
			}
		}
	}

	// Get bond/bridge slaves if this is a master
	if iface.Type == "bond" || iface.Type == "bridge" {
		i.getSlaves(iface)
	}
}

// getSlaves gets slave interfaces for bond/bridge
func (i *InterfaceManager) getSlaves(iface *NetworkInterface) {
	result, err := i.shell.Execute("ip", "link", "show", "master", iface.Name)
	if err == nil {
		lines := strings.Split(result.Stdout, "\n")
		for _, line := range lines {
			if strings.Contains(line, ":") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					slaveName := strings.TrimSuffix(fields[1], ":")
					if slaveName != iface.Name {
						iface.Slaves = append(iface.Slaves, slaveName)
					}
				}
			}
		}
	}
}

// GetInterface gets details for a specific interface
func (i *InterfaceManager) GetInterface(name string) (*NetworkInterface, error) {
	interfaces, err := i.ListInterfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		if iface.Name == name {
			return &iface, nil
		}
	}

	return nil, fmt.Errorf("interface not found: %s", name)
}

// SetInterfaceUp brings an interface up
func (i *InterfaceManager) SetInterfaceUp(name string) error {
	_, err := i.shell.Execute("ip", "link", "set", name, "up")
	if err != nil {
		return fmt.Errorf("failed to bring interface up: %w", err)
	}
	return nil
}

// SetInterfaceDown brings an interface down
func (i *InterfaceManager) SetInterfaceDown(name string) error {
	_, err := i.shell.Execute("ip", "link", "set", name, "down")
	if err != nil {
		return fmt.Errorf("failed to bring interface down: %w", err)
	}
	return nil
}

// SetIPAddress sets an IP address on an interface
func (i *InterfaceManager) SetIPAddress(name string, address string, netmask string) error {
	cidr := fmt.Sprintf("%s/%s", address, netmask)

	// Remove existing addresses
	_, _ = i.shell.Execute("ip", "addr", "flush", "dev", name)

	// Add new address
	_, err := i.shell.Execute("ip", "addr", "add", cidr, "dev", name)
	if err != nil {
		return fmt.Errorf("failed to set IP address: %w", err)
	}

	return nil
}

// SetDHCP configures interface for DHCP
func (i *InterfaceManager) SetDHCP(name string) error {
	// This is distribution-specific
	// For Debian/Ubuntu with dhclient:
	if i.shell.CommandExists("dhclient") {
		_, err := i.shell.Execute("dhclient", name)
		if err != nil {
			return fmt.Errorf("failed to start DHCP client: %w", err)
		}
		return nil
	}

	// For systems with dhcpcd:
	if i.shell.CommandExists("dhcpcd") {
		_, err := i.shell.Execute("dhcpcd", name)
		if err != nil {
			return fmt.Errorf("failed to start DHCP client: %w", err)
		}
		return nil
	}

	return fmt.Errorf("no DHCP client found")
}

// SetMTU sets the MTU for an interface
func (i *InterfaceManager) SetMTU(name string, mtu int) error {
	_, err := i.shell.Execute("ip", "link", "set", name, "mtu", fmt.Sprintf("%d", mtu))
	if err != nil {
		return fmt.Errorf("failed to set MTU: %w", err)
	}
	return nil
}

// CreateBond creates a bonded interface
func (i *InterfaceManager) CreateBond(config BondConfig) error {
	// Load bonding module if not already loaded
	_, _ = i.shell.Execute("modprobe", "bonding")

	// Create bond interface
	_, err := i.shell.Execute("ip", "link", "add", config.Name, "type", "bond", "mode", config.Mode)
	if err != nil {
		return fmt.Errorf("failed to create bond: %w", err)
	}

	// Set miimon if specified
	if config.MIIMon > 0 {
		bondPath := fmt.Sprintf("/sys/class/net/%s/bonding/miimon", config.Name)
		_, _ = i.shell.Execute("bash", "-c", fmt.Sprintf("echo %d > %s", config.MIIMon, bondPath))
	}

	// Set primary if specified
	if config.Primary != "" {
		bondPath := fmt.Sprintf("/sys/class/net/%s/bonding/primary", config.Name)
		_, _ = i.shell.Execute("bash", "-c", fmt.Sprintf("echo %s > %s", config.Primary, bondPath))
	}

	// Add slaves
	for _, slave := range config.Slaves {
		// Bring slave down first
		_, _ = i.shell.Execute("ip", "link", "set", slave, "down")

		// Add to bond
		_, err := i.shell.Execute("ip", "link", "set", slave, "master", config.Name)
		if err != nil {
			return fmt.Errorf("failed to add slave %s: %w", slave, err)
		}

		// Bring slave up
		_, _ = i.shell.Execute("ip", "link", "set", slave, "up")
	}

	// Bring bond up
	return i.SetInterfaceUp(config.Name)
}

// DeleteBond deletes a bonded interface
func (i *InterfaceManager) DeleteBond(name string) error {
	// Remove all slaves first
	iface, err := i.GetInterface(name)
	if err == nil && iface != nil {
		for _, slave := range iface.Slaves {
			_, _ = i.shell.Execute("ip", "link", "set", slave, "nomaster")
		}
	}

	// Delete bond interface
	_, err = i.shell.Execute("ip", "link", "delete", name, "type", "bond")
	if err != nil {
		return fmt.Errorf("failed to delete bond: %w", err)
	}

	return nil
}

// CreateBridge creates a bridge interface
func (i *InterfaceManager) CreateBridge(name string, slaves []string) error {
	// Create bridge
	_, err := i.shell.Execute("ip", "link", "add", name, "type", "bridge")
	if err != nil {
		return fmt.Errorf("failed to create bridge: %w", err)
	}

	// Add slaves
	for _, slave := range slaves {
		_, err := i.shell.Execute("ip", "link", "set", slave, "master", name)
		if err != nil {
			return fmt.Errorf("failed to add slave %s: %w", slave, err)
		}
	}

	// Bring bridge up
	return i.SetInterfaceUp(name)
}

// DeleteBridge deletes a bridge interface
func (i *InterfaceManager) DeleteBridge(name string) error {
	_, err := i.shell.Execute("ip", "link", "delete", name, "type", "bridge")
	if err != nil {
		return fmt.Errorf("failed to delete bridge: %w", err)
	}
	return nil
}

// CreateVLAN creates a VLAN interface
func (i *InterfaceManager) CreateVLAN(parent string, vlanID int) error {
	vlanName := fmt.Sprintf("%s.%d", parent, vlanID)

	_, err := i.shell.Execute("ip", "link", "add", "link", parent, "name", vlanName,
		"type", "vlan", "id", fmt.Sprintf("%d", vlanID))
	if err != nil {
		return fmt.Errorf("failed to create VLAN: %w", err)
	}

	return i.SetInterfaceUp(vlanName)
}

// DeleteVLAN deletes a VLAN interface
func (i *InterfaceManager) DeleteVLAN(parent string, vlanID int) error {
	vlanName := fmt.Sprintf("%s.%d", parent, vlanID)

	_, err := i.shell.Execute("ip", "link", "delete", vlanName)
	if err != nil {
		return fmt.Errorf("failed to delete VLAN: %w", err)
	}
	return nil
}

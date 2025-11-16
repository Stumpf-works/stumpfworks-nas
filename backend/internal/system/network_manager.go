// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package system

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/network"
)

// NetworkManager manages all network-related operations
type NetworkManager struct {
	shell *ShellExecutor

	// Subsystems
	Interfaces *network.InterfaceManager
	Firewall   *network.FirewallManager
	DNS        *network.DNSManager
}

// NewNetworkManager creates a new network manager
func NewNetworkManager(shell *ShellExecutor) (*NetworkManager, error) {
	nm := &NetworkManager{
		shell: shell,
	}

	// Initialize interface manager
	interfaces, err := network.NewInterfaceManager(shell)
	if err != nil {
		return nil, err
	}
	nm.Interfaces = interfaces

	// Initialize firewall manager
	firewall, err := network.NewFirewallManager(shell)
	if err != nil {
		// Firewall is optional
		nm.Firewall = nil
	} else {
		nm.Firewall = firewall
	}

	// Initialize DNS manager
	dns, err := network.NewDNSManager(shell)
	if err != nil {
		// DNS is optional
		nm.DNS = nil
	} else {
		nm.DNS = dns
	}

	return nm, nil
}

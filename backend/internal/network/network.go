// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package network

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Interface represents a network interface
type Interface struct {
	Name         string   `json:"name"`
	Index        int      `json:"index"`
	HardwareAddr string   `json:"hardwareAddr"`
	Flags        []string `json:"flags"`
	MTU          int      `json:"mtu"`
	Addresses    []string `json:"addresses"`
	IsUp         bool     `json:"isUp"`
	Speed        string   `json:"speed"`
	Type         string   `json:"type"`
}

// InterfaceStats represents network interface statistics
type InterfaceStats struct {
	Name        string `json:"name"`
	RxBytes     uint64 `json:"rxBytes"`
	TxBytes     uint64 `json:"txBytes"`
	RxPackets   uint64 `json:"rxPackets"`
	TxPackets   uint64 `json:"txPackets"`
	RxErrors    uint64 `json:"rxErrors"`
	TxErrors    uint64 `json:"txErrors"`
	RxDropped   uint64 `json:"rxDropped"`
	TxDropped   uint64 `json:"txDropped"`
}

// Route represents a network route
type Route struct {
	Destination string `json:"destination"`
	Gateway     string `json:"gateway"`
	Genmask     string `json:"genmask"`
	Flags       string `json:"flags"`
	Metric      int    `json:"metric"`
	Iface       string `json:"iface"`
}

// DNSConfig represents DNS configuration
type DNSConfig struct {
	Nameservers []string `json:"nameservers"`
	SearchDomains []string `json:"searchDomains"`
}

// FirewallRule represents a firewall rule
type FirewallRule struct {
	Number      int    `json:"number"`
	Action      string `json:"action"`
	From        string `json:"from"`
	To          string `json:"to"`
	Protocol    string `json:"protocol,omitempty"`
	Port        string `json:"port,omitempty"`
	Description string `json:"description,omitempty"`
}

// DiagnosticResult represents diagnostic command results
type DiagnosticResult struct {
	Command string `json:"command"`
	Output  string `json:"output"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// ListInterfaces returns all network interfaces
func ListInterfaces() ([]Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to list interfaces: %w", err)
	}

	var result []Interface
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		var addresses []string
		for _, addr := range addrs {
			addresses = append(addresses, addr.String())
		}

		flags := []string{}
		if iface.Flags&net.FlagUp != 0 {
			flags = append(flags, "UP")
		}
		if iface.Flags&net.FlagBroadcast != 0 {
			flags = append(flags, "BROADCAST")
		}
		if iface.Flags&net.FlagLoopback != 0 {
			flags = append(flags, "LOOPBACK")
		}
		if iface.Flags&net.FlagMulticast != 0 {
			flags = append(flags, "MULTICAST")
		}

		// Determine interface type
		ifaceType := "ethernet"
		if strings.HasPrefix(iface.Name, "wl") || strings.HasPrefix(iface.Name, "wifi") {
			ifaceType = "wireless"
		} else if strings.HasPrefix(iface.Name, "lo") {
			ifaceType = "loopback"
		} else if strings.HasPrefix(iface.Name, "br") {
			ifaceType = "bridge"
		} else if strings.HasPrefix(iface.Name, "veth") || strings.HasPrefix(iface.Name, "docker") {
			ifaceType = "virtual"
		}

		// Get interface speed
		speed := getInterfaceSpeed(iface.Name)

		result = append(result, Interface{
			Name:         iface.Name,
			Index:        iface.Index,
			HardwareAddr: iface.HardwareAddr.String(),
			Flags:        flags,
			MTU:          iface.MTU,
			Addresses:    addresses,
			IsUp:         iface.Flags&net.FlagUp != 0,
			Speed:        speed,
			Type:         ifaceType,
		})
	}

	return result, nil
}

// getInterfaceSpeed tries to get interface speed from sysfs
func getInterfaceSpeed(name string) string {
	speedFile := fmt.Sprintf("/sys/class/net/%s/speed", name)
	data, err := os.ReadFile(speedFile)
	if err != nil {
		return "Unknown"
	}
	speed := strings.TrimSpace(string(data))
	if speedInt, err := strconv.Atoi(speed); err == nil && speedInt > 0 {
		if speedInt >= 1000 {
			return fmt.Sprintf("%d Gbps", speedInt/1000)
		}
		return fmt.Sprintf("%d Mbps", speedInt)
	}
	return "Unknown"
}

// GetInterfaceStats returns statistics for all interfaces
func GetInterfaceStats() ([]InterfaceStats, error) {
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc/net/dev: %w", err)
	}

	var stats []InterfaceStats
	scanner := bufio.NewScanner(bytes.NewReader(data))

	// Skip header lines
	scanner.Scan()
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 17 {
			continue
		}

		name := strings.TrimSuffix(fields[0], ":")

		stat := InterfaceStats{
			Name: name,
		}

		// Parse statistics (format: bytes packets errs drop fifo frame compressed multicast)
		if val, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
			stat.RxBytes = val
		}
		if val, err := strconv.ParseUint(fields[2], 10, 64); err == nil {
			stat.RxPackets = val
		}
		if val, err := strconv.ParseUint(fields[3], 10, 64); err == nil {
			stat.RxErrors = val
		}
		if val, err := strconv.ParseUint(fields[4], 10, 64); err == nil {
			stat.RxDropped = val
		}
		if val, err := strconv.ParseUint(fields[9], 10, 64); err == nil {
			stat.TxBytes = val
		}
		if val, err := strconv.ParseUint(fields[10], 10, 64); err == nil {
			stat.TxPackets = val
		}
		if val, err := strconv.ParseUint(fields[11], 10, 64); err == nil {
			stat.TxErrors = val
		}
		if val, err := strconv.ParseUint(fields[12], 10, 64); err == nil {
			stat.TxDropped = val
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

// SetInterfaceUp brings an interface up
func SetInterfaceUp(name string) error {
	cmd := exec.Command("ip", "link", "set", name, "up")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to bring interface up: %s", string(output))
	}
	return nil
}

// SetInterfaceDown brings an interface down
func SetInterfaceDown(name string) error {
	cmd := exec.Command("ip", "link", "set", name, "down")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to bring interface down: %s", string(output))
	}
	return nil
}

// ConfigureStaticIP configures a static IP address on an interface
func ConfigureStaticIP(name, ipAddress, netmask, gateway string) error {
	// Remove existing IP addresses
	cmd := exec.Command("ip", "addr", "flush", "dev", name)
	cmd.Run()

	// Calculate CIDR notation
	cidr := calculateCIDR(netmask)
	ipWithCIDR := fmt.Sprintf("%s/%d", ipAddress, cidr)

	// Add new IP address
	cmd = exec.Command("ip", "addr", "add", ipWithCIDR, "dev", name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set IP: %s", string(output))
	}

	// Set default gateway if provided
	if gateway != "" {
		// Remove existing default route
		exec.Command("ip", "route", "del", "default").Run()

		// Add new default route
		cmd = exec.Command("ip", "route", "add", "default", "via", gateway, "dev", name)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to set gateway: %s", string(output))
		}
	}

	return nil
}

// ConfigureDHCP configures an interface to use DHCP
func ConfigureDHCP(name string) error {
	// This would typically require dhclient or dhcpcd
	cmd := exec.Command("dhclient", name)
	if _, err := cmd.CombinedOutput(); err != nil {
		// Try dhcpcd as fallback
		cmd = exec.Command("dhcpcd", name)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to configure DHCP: %s", string(output))
		}
	}
	return nil
}

// calculateCIDR converts netmask to CIDR notation
func calculateCIDR(netmask string) int {
	ip := net.ParseIP(netmask).To4()
	if ip == nil {
		return 24 // default
	}

	var cidr int
	for _, octet := range ip {
		for octet != 0 {
			cidr += int(octet & 1)
			octet >>= 1
		}
	}
	return cidr
}

// GetRoutes returns the routing table
func GetRoutes() ([]Route, error) {
	cmd := exec.Command("ip", "route", "show")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get routes: %w", err)
	}

	var routes []Route
	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		route := Route{
			Destination: fields[0],
		}

		// Parse route details
		for i := 1; i < len(fields); i++ {
			switch fields[i] {
			case "via":
				if i+1 < len(fields) {
					route.Gateway = fields[i+1]
					i++
				}
			case "dev":
				if i+1 < len(fields) {
					route.Iface = fields[i+1]
					i++
				}
			case "metric":
				if i+1 < len(fields) {
					if m, err := strconv.Atoi(fields[i+1]); err == nil {
						route.Metric = m
					}
					i++
				}
			}
		}

		routes = append(routes, route)
	}

	return routes, nil
}

// AddRoute adds a static route
func AddRoute(destination, gateway, iface string, metric int) error {
	args := []string{"route", "add"}

	// Add destination
	args = append(args, destination)

	// Add gateway if provided
	if gateway != "" {
		args = append(args, "via", gateway)
	}

	// Add interface if provided
	if iface != "" {
		args = append(args, "dev", iface)
	}

	// Add metric if provided
	if metric > 0 {
		args = append(args, "metric", strconv.Itoa(metric))
	}

	cmd := exec.Command("ip", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add route: %w: %s", err, string(output))
	}

	return nil
}

// DeleteRoute deletes a static route
func DeleteRoute(destination, gateway, iface string) error {
	args := []string{"route", "del", destination}

	// Add gateway if provided (helps identify specific route)
	if gateway != "" {
		args = append(args, "via", gateway)
	}

	// Add interface if provided (helps identify specific route)
	if iface != "" {
		args = append(args, "dev", iface)
	}

	cmd := exec.Command("ip", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete route: %w: %s", err, string(output))
	}

	return nil
}

// GetDNSConfig returns DNS configuration
func GetDNSConfig() (*DNSConfig, error) {
	data, err := os.ReadFile("/etc/resolv.conf")
	if err != nil {
		return nil, fmt.Errorf("failed to read resolv.conf: %w", err)
	}

	config := &DNSConfig{
		Nameservers:   []string{},
		SearchDomains: []string{},
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "nameserver":
			config.Nameservers = append(config.Nameservers, fields[1])
		case "search", "domain":
			config.SearchDomains = append(config.SearchDomains, fields[1:]...)
		}
	}

	return config, nil
}

// SetDNSConfig updates DNS configuration
func SetDNSConfig(nameservers []string, searchDomains []string) error {
	var content strings.Builder

	content.WriteString("# Generated by Stumpf.Works NAS\n")
	content.WriteString(fmt.Sprintf("# %s\n\n", time.Now().Format(time.RFC3339)))

	for _, ns := range nameservers {
		content.WriteString(fmt.Sprintf("nameserver %s\n", ns))
	}

	if len(searchDomains) > 0 {
		content.WriteString(fmt.Sprintf("search %s\n", strings.Join(searchDomains, " ")))
	}

	return os.WriteFile("/etc/resolv.conf", []byte(content.String()), 0644)
}

// Ping executes a ping command
func Ping(host string, count int) (*DiagnosticResult, error) {
	cmd := exec.Command("ping", "-c", strconv.Itoa(count), host)
	output, err := cmd.CombinedOutput()

	result := &DiagnosticResult{
		Command: fmt.Sprintf("ping -c %d %s", count, host),
		Output:  string(output),
		Success: err == nil,
	}

	if err != nil {
		result.Error = err.Error()
	}

	return result, nil
}

// Traceroute executes a traceroute command
func Traceroute(host string) (*DiagnosticResult, error) {
	cmd := exec.Command("traceroute", host)
	output, err := cmd.CombinedOutput()

	result := &DiagnosticResult{
		Command: fmt.Sprintf("traceroute %s", host),
		Output:  string(output),
		Success: err == nil,
	}

	if err != nil {
		result.Error = err.Error()
	}

	return result, nil
}

// Netstat executes netstat command
func Netstat(options string) (*DiagnosticResult, error) {
	args := []string{}
	if options != "" {
		args = strings.Split(options, " ")
	} else {
		args = []string{"-tuln"}
	}

	cmd := exec.Command("netstat", args...)
	output, err := cmd.CombinedOutput()

	// Try ss if netstat is not available
	if err != nil {
		cmd = exec.Command("ss", args...)
		output, err = cmd.CombinedOutput()
	}

	result := &DiagnosticResult{
		Command: fmt.Sprintf("netstat %s", strings.Join(args, " ")),
		Output:  string(output),
		Success: err == nil,
	}

	if err != nil {
		result.Error = err.Error()
	}

	return result, nil
}

// WakeOnLAN sends a magic packet to wake a device
func WakeOnLAN(macAddress string) error {
	// Parse MAC address
	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		return fmt.Errorf("invalid MAC address: %w", err)
	}

	// Create magic packet: 6 bytes of 0xFF followed by 16 repetitions of MAC
	packet := make([]byte, 102)
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	for i := 0; i < 16; i++ {
		copy(packet[6+i*6:], mac)
	}

	// Broadcast on port 9
	conn, err := net.Dial("udp", "255.255.255.255:9")
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	if err != nil {
		return fmt.Errorf("failed to send magic packet: %w", err)
	}

	return nil
}

// CreateBridge creates a new bridge interface with Proxmox-style IP migration
// This safely migrates IP addresses from physical interfaces to the bridge
func CreateBridge(name string, ports []string) error {
	// Create the bridge
	cmd := exec.Command("ip", "link", "add", name, "type", "bridge")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create bridge: %s", string(output))
	}

	// Bring the bridge up immediately
	cmd = exec.Command("ip", "link", "set", name, "up")
	if output, err := cmd.CombinedOutput(); err != nil {
		exec.Command("ip", "link", "delete", name, "type", "bridge").Run()
		return fmt.Errorf("failed to bring bridge up: %s", string(output))
	}

	// Add ports to the bridge if specified
	// NOTE: This is a safe operation ONLY if the bridge already has IP configuration
	// or if the port being added has no IP addresses (is not the primary interface)
	for _, port := range ports {
		if port == "" {
			continue
		}

		// Get current IP addresses and routes from the port
		portAddrs, err := getInterfaceAddresses(port)
		if err != nil {
			continue // Skip if we can't get addresses
		}

		// Get default gateway if this interface has one
		gateway, _ := getDefaultGatewayForInterface(port)

		// SAFETY CHECK: If the port has IP addresses but the bridge doesn't,
		// DO NOT proceed - this could break network connectivity!
		bridgeAddrs, _ := getInterfaceAddresses(name)
		if len(portAddrs) > 0 && len(bridgeAddrs) == 0 {
			// This is dangerous - the port has IPs but the bridge doesn't
			// Return an error to prevent network loss
			return fmt.Errorf("cannot add port %s with IP addresses to bridge %s without IP configuration - this would break network connectivity. Please configure the bridge IP first", port, name)
		}

		// If both port and bridge have IP addresses, this is a safe migration
		// The bridge already has connectivity, so we can safely move the port
		if len(portAddrs) > 0 && len(bridgeAddrs) > 0 {
			// Step 1: If there's a default gateway on the port, ensure bridge has it too
			if gateway != "" {
				// Check if bridge already has default route
				bridgeGateway, _ := getDefaultGatewayForInterface(name)
				if bridgeGateway == "" {
					// Remove old default route via port
					exec.Command("ip", "route", "del", "default", "dev", port).Run()

					// Add new default route via bridge
					cmd = exec.Command("ip", "route", "add", "default", "via", gateway, "dev", name)
					cmd.CombinedOutput()
				}
			}

			// Step 2: Remove IPs from the port (safe because bridge has IPs)
			cmd = exec.Command("ip", "addr", "flush", "dev", port)
			cmd.Run()
		}

		// Step 3: Attach port to bridge
		cmd = exec.Command("ip", "link", "set", port, "master", name)
		if output, err := cmd.CombinedOutput(); err != nil {
			// If attachment fails and we removed IPs, try to restore them
			if len(portAddrs) > 0 && len(bridgeAddrs) > 0 {
				for _, addr := range portAddrs {
					exec.Command("ip", "addr", "add", addr, "dev", port).Run()
				}
				if gateway != "" {
					exec.Command("ip", "route", "del", "default", "dev", name).Run()
					exec.Command("ip", "route", "add", "default", "via", gateway, "dev", port).Run()
				}
			}
			// Clean up the bridge
			exec.Command("ip", "link", "delete", name, "type", "bridge").Run()
			return fmt.Errorf("failed to attach port %s to bridge: %s", port, string(output))
		}

		// Step 4: Ensure port is up as a bridge port
		exec.Command("ip", "link", "set", port, "up").Run()
	}

	// Step 6: Add iptables rules to allow forwarding through the bridge
	// This is essential for containers/VMs to communicate with the external network
	exec.Command("iptables", "-I", "FORWARD", "-i", name, "-o", name, "-j", "ACCEPT").Run()
	exec.Command("iptables", "-I", "FORWARD", "-i", name, "-j", "ACCEPT").Run()
	exec.Command("iptables", "-I", "FORWARD", "-o", name, "-j", "ACCEPT").Run()

	return nil
}

// getInterfaceAddresses retrieves IP addresses configured on an interface
func getInterfaceAddresses(ifaceName string) ([]string, error) {
	cmd := exec.Command("ip", "-o", "addr", "show", ifaceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var addresses []string
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		// Format: index: ifacename inet|inet6 addr/prefix ...
		for i, field := range fields {
			if (field == "inet" || field == "inet6") && i+1 < len(fields) {
				addr := fields[i+1]
				// Skip link-local IPv6 addresses
				if !strings.HasPrefix(addr, "fe80:") {
					addresses = append(addresses, addr)
				}
			}
		}
	}

	return addresses, nil
}

// getDefaultGatewayForInterface finds the default gateway for a specific interface
func getDefaultGatewayForInterface(ifaceName string) (string, error) {
	cmd := exec.Command("ip", "route", "show", "default", "dev", ifaceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// Parse: default via <gateway> dev <iface> ...
	fields := strings.Fields(string(output))
	for i, field := range fields {
		if field == "via" && i+1 < len(fields) {
			return fields[i+1], nil
		}
	}

	return "", nil
}

// DeleteBridge deletes a bridge interface
func DeleteBridge(name string) error {
	// Get all ports attached to this bridge
	cmd := exec.Command("ip", "link", "show", "master", name)
	output, _ := cmd.CombinedOutput()

	// Parse output to find ports
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ":") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				portName := strings.TrimSuffix(fields[1], ":")
				if portName != name {
					// Remove port from bridge
					exec.Command("ip", "link", "set", portName, "nomaster").Run()
				}
			}
		}
	}

	// Bring bridge down
	exec.Command("ip", "link", "set", name, "down").Run()

	// Delete the bridge
	cmd = exec.Command("ip", "link", "delete", name, "type", "bridge")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to delete bridge: %s", string(output))
	}

	return nil
}

// AttachPortToBridge attaches an interface to a bridge with IP migration
func AttachPortToBridge(bridgeName string, portName string) error {
	// Get current IP addresses and routes from the port
	portAddrs, err := getInterfaceAddresses(portName)
	if err != nil {
		portAddrs = []string{} // Continue with empty list if error
	}

	// Get default gateway if this interface has one
	gateway, _ := getDefaultGatewayForInterface(portName)

	// SAFETY CHECK: If the port has IP addresses but the bridge doesn't,
	// DO NOT proceed - this could break network connectivity!
	bridgeAddrs, _ := getInterfaceAddresses(bridgeName)
	if len(portAddrs) > 0 && len(bridgeAddrs) == 0 {
		// This is dangerous - the port has IPs but the bridge doesn't
		// Return an error to prevent network loss
		return fmt.Errorf("cannot add port %s with IP addresses to bridge %s without IP configuration - this would break network connectivity. Please configure the bridge IP first", portName, bridgeName)
	}

	// If both port and bridge have IP addresses, this is a safe migration
	// The bridge already has connectivity, so we can safely move the port
	if len(portAddrs) > 0 && len(bridgeAddrs) > 0 {
		// Step 1: If there's a default gateway on the port, ensure bridge has it too
		if gateway != "" {
			// Check if bridge already has default route
			bridgeGateway, _ := getDefaultGatewayForInterface(bridgeName)
			if bridgeGateway == "" {
				// Remove old default route via port
				exec.Command("ip", "route", "del", "default", "dev", portName).Run()

				// Add new default route via bridge
				cmd := exec.Command("ip", "route", "add", "default", "via", gateway, "dev", bridgeName)
				cmd.CombinedOutput()
			}
		}

		// Step 2: Remove IPs from the port (safe because bridge has IPs)
		cmd := exec.Command("ip", "addr", "flush", "dev", portName)
		cmd.Run()
	}

	// Step 3: Attach port to bridge
	cmd := exec.Command("ip", "link", "set", portName, "master", bridgeName)
	if output, err := cmd.CombinedOutput(); err != nil {
		// If attachment fails and we removed IPs, try to restore them
		if len(portAddrs) > 0 && len(bridgeAddrs) > 0 {
			for _, addr := range portAddrs {
				exec.Command("ip", "addr", "add", addr, "dev", portName).Run()
			}
			if gateway != "" {
				exec.Command("ip", "route", "del", "default", "dev", bridgeName).Run()
				exec.Command("ip", "route", "add", "default", "via", gateway, "dev", portName).Run()
			}
		}
		return fmt.Errorf("failed to attach port to bridge: %s", string(output))
	}

	// Step 4: Ensure port is up as a bridge port
	cmd = exec.Command("ip", "link", "set", portName, "up")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to bring port up: %s", string(output))
	}

	return nil
}

// DetachPortFromBridge detaches an interface from a bridge
func DetachPortFromBridge(portName string) error {
	// Remove port from bridge
	cmd := exec.Command("ip", "link", "set", portName, "nomaster")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to detach port from bridge: %s", string(output))
	}

	return nil
}

// ListBridges returns a list of all bridge interfaces
func ListBridges() ([]string, error) {
	cmd := exec.Command("ip", "-o", "link", "show", "type", "bridge")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If no bridges exist, this is not an error
		return []string{}, nil
	}

	var bridges []string
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		// Format: index: bridge_name: <BROADCAST,MULTICAST,UP,LOWER_UP> ...
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			// Remove the trailing colon from the bridge name
			bridgeName := strings.TrimSuffix(fields[1], ":")
			bridges = append(bridges, bridgeName)
		}
	}

	return bridges, nil
}

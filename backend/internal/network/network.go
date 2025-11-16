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

// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package network

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"fmt"
	"os"
	"strings"
)

// DNSManager manages DNS configuration
type DNSManager struct {
	shell   executor.ShellExecutor
	enabled bool
}

// DNSConfig represents DNS configuration
type DNSConfig struct {
	Nameservers []string `json:"nameservers"`
	Search      []string `json:"search"`
	Domain      string   `json:"domain"`
}

// NewDNSManager creates a new DNS manager
func NewDNSManager(shell executor.ShellExecutor) (*DNSManager, error) {
	return &DNSManager{
		shell:   shell,
		enabled: true,
	}, nil
}

// IsEnabled returns whether DNS management is available
func (d *DNSManager) IsEnabled() bool {
	return d.enabled
}

// GetConfig gets current DNS configuration
func (d *DNSManager) GetConfig() (*DNSConfig, error) {
	config := &DNSConfig{
		Nameservers: make([]string, 0),
		Search:      make([]string, 0),
	}

	// Read /etc/resolv.conf
	data, err := os.ReadFile("/etc/resolv.conf")
	if err != nil {
		return nil, fmt.Errorf("failed to read resolv.conf: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "nameserver":
			config.Nameservers = append(config.Nameservers, fields[1])
		case "search":
			config.Search = append(config.Search, fields[1:]...)
		case "domain":
			config.Domain = fields[1]
		}
	}

	return config, nil
}

// SetConfig sets DNS configuration
func (d *DNSManager) SetConfig(config DNSConfig) error {
	var lines []string

	// Add domain if specified
	if config.Domain != "" {
		lines = append(lines, fmt.Sprintf("domain %s", config.Domain))
	}

	// Add search domains
	if len(config.Search) > 0 {
		lines = append(lines, fmt.Sprintf("search %s", strings.Join(config.Search, " ")))
	}

	// Add nameservers (max 3 typically)
	for i, ns := range config.Nameservers {
		if i >= 3 {
			break // Most systems only use first 3 nameservers
		}
		lines = append(lines, fmt.Sprintf("nameserver %s", ns))
	}

	content := strings.Join(lines, "\n") + "\n"

	// Write to resolv.conf
	err := os.WriteFile("/etc/resolv.conf", []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write resolv.conf: %w", err)
	}

	return nil
}

// AddNameserver adds a nameserver
func (d *DNSManager) AddNameserver(nameserver string) error {
	config, err := d.GetConfig()
	if err != nil {
		return err
	}

	// Check if already exists
	for _, ns := range config.Nameservers {
		if ns == nameserver {
			return nil // Already exists
		}
	}

	config.Nameservers = append(config.Nameservers, nameserver)
	return d.SetConfig(*config)
}

// RemoveNameserver removes a nameserver
func (d *DNSManager) RemoveNameserver(nameserver string) error {
	config, err := d.GetConfig()
	if err != nil {
		return err
	}

	var newNameservers []string
	for _, ns := range config.Nameservers {
		if ns != nameserver {
			newNameservers = append(newNameservers, ns)
		}
	}

	config.Nameservers = newNameservers
	return d.SetConfig(*config)
}

// SetSearchDomains sets search domains
func (d *DNSManager) SetSearchDomains(domains []string) error {
	config, err := d.GetConfig()
	if err != nil {
		return err
	}

	config.Search = domains
	return d.SetConfig(*config)
}

// SetHostname sets the system hostname
func (d *DNSManager) SetHostname(hostname string) error {
	// Set runtime hostname
	_, err := d.shell.Execute("hostname", hostname)
	if err != nil {
		return fmt.Errorf("failed to set runtime hostname: %w", err)
	}

	// Write to /etc/hostname for persistence
	err = os.WriteFile("/etc/hostname", []byte(hostname+"\n"), 0644)
	if err != nil {
		return fmt.Errorf("failed to write /etc/hostname: %w", err)
	}

	return nil
}

// GetHostname gets the system hostname
func (d *DNSManager) GetHostname() (string, error) {
	result, err := d.shell.Execute("hostname")
	if err != nil {
		return "", fmt.Errorf("failed to get hostname: %w", err)
	}

	return strings.TrimSpace(result.Stdout), nil
}

// UpdateHosts updates /etc/hosts file
func (d *DNSManager) UpdateHosts(hostname string, ip string) error {
	// Read current hosts file
	data, err := os.ReadFile("/etc/hosts")
	if err != nil {
		return fmt.Errorf("failed to read /etc/hosts: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	hostnameFound := false

	for _, line := range lines {
		// Skip empty lines and comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			newLines = append(newLines, line)
			continue
		}

		// Check if this line contains our hostname
		if strings.Contains(line, hostname) {
			// Replace with new entry
			newLines = append(newLines, fmt.Sprintf("%s\t%s", ip, hostname))
			hostnameFound = true
		} else {
			newLines = append(newLines, line)
		}
	}

	// If hostname wasn't found, add it
	if !hostnameFound {
		newLines = append(newLines, fmt.Sprintf("%s\t%s", ip, hostname))
	}

	content := strings.Join(newLines, "\n")
	err = os.WriteFile("/etc/hosts", []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write /etc/hosts: %w", err)
	}

	return nil
}

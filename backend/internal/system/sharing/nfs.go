// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package sharing

import (
	"fmt"
	"os"
	"strings"
)

// NFSManager manages NFS exports
type NFSManager struct {
	shell      ShellExecutor
	enabled    bool
	exportsPath string
}

// NFSExport represents an NFS export configuration
type NFSExport struct {
	Path        string   `json:"path"`
	Clients     []string `json:"clients"` // IP/CIDR or * for all
	Options     []string `json:"options"`
	ReadOnly    bool     `json:"read_only"`
	Sync        bool     `json:"sync"`
	NoRootSquash bool    `json:"no_root_squash"`
	Subtree     bool     `json:"subtree"`
}

// NewNFSManager creates a new NFS manager
func NewNFSManager(shell ShellExecutor) (*NFSManager, error) {
	if !shell.CommandExists("exportfs") {
		return nil, fmt.Errorf("nfs-kernel-server not installed")
	}

	return &NFSManager{
		shell:       shell,
		enabled:     true,
		exportsPath: "/etc/exports",
	}, nil
}

// IsEnabled returns whether NFS is available
func (n *NFSManager) IsEnabled() bool {
	return n.enabled
}

// GetStatus gets NFS service status
func (n *NFSManager) GetStatus() (bool, error) {
	result, err := n.shell.Execute("systemctl", "is-active", "nfs-server")
	if err != nil {
		// Try alternative service name
		result, err = n.shell.Execute("systemctl", "is-active", "nfs-kernel-server")
		if err != nil {
			return false, nil
		}
	}

	return strings.TrimSpace(result.Stdout) == "active", nil
}

// Start starts the NFS service
func (n *NFSManager) Start() error {
	// Try nfs-server first (systemd standard name)
	_, err := n.shell.Execute("systemctl", "start", "nfs-server")
	if err != nil {
		// Try nfs-kernel-server (Debian/Ubuntu)
		_, err = n.shell.Execute("systemctl", "start", "nfs-kernel-server")
		if err != nil {
			return fmt.Errorf("failed to start NFS: %w", err)
		}
	}

	return nil
}

// Stop stops the NFS service
func (n *NFSManager) Stop() error {
	_, err := n.shell.Execute("systemctl", "stop", "nfs-server")
	if err != nil {
		_, err = n.shell.Execute("systemctl", "stop", "nfs-kernel-server")
		if err != nil {
			return fmt.Errorf("failed to stop NFS: %w", err)
		}
	}

	return nil
}

// Restart restarts the NFS service
func (n *NFSManager) Restart() error {
	_, err := n.shell.Execute("systemctl", "restart", "nfs-server")
	if err != nil {
		_, err = n.shell.Execute("systemctl", "restart", "nfs-kernel-server")
		if err != nil {
			return fmt.Errorf("failed to restart NFS: %w", err)
		}
	}

	return nil
}

// Reload reloads NFS exports without restarting
func (n *NFSManager) Reload() error {
	_, err := n.shell.Execute("exportfs", "-ra")
	if err != nil {
		return fmt.Errorf("failed to reload exports: %w", err)
	}

	return nil
}

// ListExports lists all NFS exports
func (n *NFSManager) ListExports() ([]NFSExport, error) {
	result, err := n.shell.Execute("exportfs", "-v")
	if err != nil {
		return nil, fmt.Errorf("failed to list exports: %w", err)
	}

	var exports []NFSExport
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Parse export line
		// Format: /path client1(options) client2(options)
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		export := NFSExport{
			Path:    parts[0],
			Clients: make([]string, 0),
			Options: make([]string, 0),
		}

		// Parse clients and options
		for i := 1; i < len(parts); i++ {
			clientOpts := parts[i]

			// Split client and options
			if strings.Contains(clientOpts, "(") {
				clientEnd := strings.Index(clientOpts, "(")
				client := clientOpts[:clientEnd]
				opts := strings.Trim(clientOpts[clientEnd:], "()")

				export.Clients = append(export.Clients, client)

				// Parse options
				optsList := strings.Split(opts, ",")
				for _, opt := range optsList {
					export.Options = append(export.Options, opt)

					switch opt {
					case "ro":
						export.ReadOnly = true
					case "sync":
						export.Sync = true
					case "no_root_squash":
						export.NoRootSquash = true
					case "no_subtree_check":
						export.Subtree = false
					}
				}
			} else {
				export.Clients = append(export.Clients, clientOpts)
			}
		}

		exports = append(exports, export)
	}

	return exports, nil
}

// GetExport gets a specific export
func (n *NFSManager) GetExport(path string) (*NFSExport, error) {
	exports, err := n.ListExports()
	if err != nil {
		return nil, err
	}

	for _, export := range exports {
		if export.Path == path {
			return &export, nil
		}
	}

	return nil, fmt.Errorf("export not found: %s", path)
}

// CreateExport creates a new NFS export
func (n *NFSManager) CreateExport(export NFSExport) error {
	// Read current exports
	data, err := os.ReadFile(n.exportsPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read exports: %w", err)
	}

	config := string(data)

	// Check if export already exists
	if strings.Contains(config, export.Path) {
		return fmt.Errorf("export already exists: %s", export.Path)
	}

	// Build export line
	exportLine := export.Path

	// Build options
	var options []string

	if export.ReadOnly {
		options = append(options, "ro")
	} else {
		options = append(options, "rw")
	}

	if export.Sync {
		options = append(options, "sync")
	} else {
		options = append(options, "async")
	}

	if export.NoRootSquash {
		options = append(options, "no_root_squash")
	} else {
		options = append(options, "root_squash")
	}

	if !export.Subtree {
		options = append(options, "no_subtree_check")
	}

	// Add any additional options
	for _, opt := range export.Options {
		if !contains(options, opt) {
			options = append(options, opt)
		}
	}

	optStr := strings.Join(options, ",")

	// Add clients
	if len(export.Clients) == 0 {
		export.Clients = []string{"*"}
	}

	for _, client := range export.Clients {
		exportLine += fmt.Sprintf(" %s(%s)", client, optStr)
	}

	exportLine += "\n"

	// Append to exports file
	config += exportLine

	err = os.WriteFile(n.exportsPath, []byte(config), 0644)
	if err != nil {
		return fmt.Errorf("failed to write exports: %w", err)
	}

	// Reload exports
	return n.Reload()
}

// UpdateExport updates an existing export
func (n *NFSManager) UpdateExport(export NFSExport) error {
	// Delete and recreate
	if err := n.DeleteExport(export.Path); err != nil {
		return err
	}

	return n.CreateExport(export)
}

// DeleteExport deletes an NFS export
func (n *NFSManager) DeleteExport(path string) error {
	data, err := os.ReadFile(n.exportsPath)
	if err != nil {
		return fmt.Errorf("failed to read exports: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	exportFound := false

	for _, line := range lines {
		// Skip the line with our export path
		if strings.HasPrefix(strings.TrimSpace(line), path) {
			exportFound = true
			continue
		}

		newLines = append(newLines, line)
	}

	if !exportFound {
		return fmt.Errorf("export not found: %s", path)
	}

	// Write back
	config := strings.Join(newLines, "\n")
	err = os.WriteFile(n.exportsPath, []byte(config), 0644)
	if err != nil {
		return fmt.Errorf("failed to write exports: %w", err)
	}

	// Reload exports
	return n.Reload()
}

// AddClient adds a client to an export
func (n *NFSManager) AddClient(path string, client string) error {
	export, err := n.GetExport(path)
	if err != nil {
		return err
	}

	// Check if client already exists
	for _, c := range export.Clients {
		if c == client {
			return nil // Already exists
		}
	}

	export.Clients = append(export.Clients, client)
	return n.UpdateExport(*export)
}

// RemoveClient removes a client from an export
func (n *NFSManager) RemoveClient(path string, client string) error {
	export, err := n.GetExport(path)
	if err != nil {
		return err
	}

	var newClients []string
	for _, c := range export.Clients {
		if c != client {
			newClients = append(newClients, c)
		}
	}

	export.Clients = newClients
	return n.UpdateExport(*export)
}

// GetActiveConnections gets currently connected NFS clients
func (n *NFSManager) GetActiveConnections() ([]string, error) {
	result, err := n.shell.Execute("showmount", "-a")
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %w", err)
	}

	var connections []string
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "All mount") {
			continue
		}

		connections = append(connections, strings.TrimSpace(line))
	}

	return connections, nil
}

// Helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

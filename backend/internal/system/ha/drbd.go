// Package ha provides High Availability features for Stumpf.Works NAS
package ha

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// DRBDManager manages DRBD (Distributed Replicated Block Device) resources
type DRBDManager struct {
	shell   executor.ShellExecutor
	enabled bool
}

// DRBDResource represents a DRBD resource configuration
type DRBDResource struct {
	Name         string `json:"name"`
	Device       string `json:"device"`        // e.g., /dev/drbd0
	Disk         string `json:"disk"`          // e.g., /dev/sda1
	MetaDisk     string `json:"meta_disk"`     // e.g., internal or /dev/sdb1
	LocalAddress string `json:"local_address"` // e.g., 192.168.1.10:7788
	PeerAddress  string `json:"peer_address"`  // e.g., 192.168.1.11:7788
	Protocol     string `json:"protocol"`      // A, B, or C
}

// DRBDStatus represents the current status of a DRBD resource
type DRBDStatus struct {
	Name            string `json:"name"`
	Device          string `json:"device"`
	ConnectionState string `json:"connection_state"` // Connected, Disconnected, StandAlone
	Role            string `json:"role"`             // Primary, Secondary, Unknown
	DiskState       string `json:"disk_state"`       // UpToDate, Inconsistent, DUnknown
	PeerRole        string `json:"peer_role"`        // Primary, Secondary, Unknown
	PeerDiskState   string `json:"peer_disk_state"`  // UpToDate, Inconsistent, DUnknown
	SyncProgress    int    `json:"sync_progress"`    // 0-100 percentage
	Resyncing       bool   `json:"resyncing"`
}

// NewDRBDManager creates a new DRBD manager
func NewDRBDManager(shell executor.ShellExecutor) (*DRBDManager, error) {
	manager := &DRBDManager{
		shell:   shell,
		enabled: false,
	}

	// Check if drbdadm is available
	result, err := shell.Execute("which", "drbdadm")
	if err != nil || result.Stdout == "" {
		logger.Warn("drbdadm not found, DRBD features will be disabled")
		return manager, fmt.Errorf("drbdadm not available: install drbd-utils package")
	}

	manager.enabled = true
	logger.Info("DRBD manager initialized successfully")
	return manager, nil
}

// IsEnabled returns whether DRBD is available
func (dm *DRBDManager) IsEnabled() bool {
	return dm.enabled
}

// CreateResource creates a new DRBD resource
func (dm *DRBDManager) CreateResource(resource DRBDResource) error {
	if !dm.enabled {
		return fmt.Errorf("DRBD is not enabled")
	}

	// Validate required fields
	if resource.Name == "" || resource.Device == "" || resource.Disk == "" {
		return fmt.Errorf("name, device, and disk are required")
	}

	// Set defaults
	if resource.Protocol == "" {
		resource.Protocol = "C" // Default to Protocol C (synchronous)
	}
	if resource.MetaDisk == "" {
		resource.MetaDisk = "internal"
	}

	// Create DRBD resource configuration file
	config := fmt.Sprintf(`resource %s {
  device %s;
  disk %s;
  meta-disk %s;

  on %s {
    address %s;
  }

  on %s {
    address %s;
  }

  net {
    protocol %s;
  }
}`, resource.Name, resource.Device, resource.Disk, resource.MetaDisk,
		"local", resource.LocalAddress,
		"peer", resource.PeerAddress,
		resource.Protocol)

	// Write config to /etc/drbd.d/
	configPath := fmt.Sprintf("/etc/drbd.d/%s.res", resource.Name)
	writeCmd := fmt.Sprintf("echo '%s' | sudo tee %s", config, configPath)
	result, err := dm.shell.Execute("sh", "-c", writeCmd)
	if err != nil {
		return fmt.Errorf("failed to write DRBD config: %s: %w", result.Stderr, err)
	}

	// Create metadata
	result, err = dm.shell.Execute("sudo", "drbdadm", "create-md", resource.Name)
	if err != nil {
		logger.Error("Failed to create DRBD metadata", zap.Error(err), zap.String("stderr", result.Stderr))
		return fmt.Errorf("failed to create DRBD metadata: %s: %w", result.Stderr, err)
	}

	// Bring up the resource
	result, err = dm.shell.Execute("sudo", "drbdadm", "up", resource.Name)
	if err != nil {
		logger.Error("Failed to bring up DRBD resource", zap.Error(err), zap.String("stderr", result.Stderr))
		return fmt.Errorf("failed to bring up DRBD resource: %s: %w", result.Stderr, err)
	}

	logger.Info("DRBD resource created", zap.String("name", resource.Name))
	return nil
}

// DeleteResource deletes a DRBD resource
func (dm *DRBDManager) DeleteResource(name string) error {
	if !dm.enabled {
		return fmt.Errorf("DRBD is not enabled")
	}

	// Bring down the resource
	result, err := dm.shell.Execute("sudo", "drbdadm", "down", name)
	if err != nil {
		logger.Warn("Failed to bring down DRBD resource", zap.Error(err), zap.String("stderr", result.Stderr))
		// Continue anyway to delete config
	}

	// Delete configuration file
	configPath := fmt.Sprintf("/etc/drbd.d/%s.res", name)
	result, err = dm.shell.Execute("sudo", "rm", "-f", configPath)
	if err != nil {
		return fmt.Errorf("failed to delete DRBD config: %s: %w", result.Stderr, err)
	}

	logger.Info("DRBD resource deleted", zap.String("name", name))
	return nil
}

// GetResourceStatus gets the status of a DRBD resource
func (dm *DRBDManager) GetResourceStatus(name string) (*DRBDStatus, error) {
	if !dm.enabled {
		return nil, fmt.Errorf("DRBD is not enabled")
	}

	// Execute drbdadm status
	result, err := dm.shell.Execute("sudo", "drbdadm", "status", name)
	if err != nil {
		return nil, fmt.Errorf("failed to get DRBD status: %s: %w", result.Stderr, err)
	}

	// Parse the status output
	status := &DRBDStatus{
		Name:            name,
		ConnectionState: "Unknown",
		Role:            "Unknown",
		DiskState:       "Unknown",
		PeerRole:        "Unknown",
		PeerDiskState:   "Unknown",
		SyncProgress:    0,
		Resyncing:       false,
	}

	// Parse output (format: name role:disk-state connection-state peer-role:peer-disk-state)
	lines := strings.Split(result.Stdout, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract device
		if strings.Contains(line, "/dev/drbd") {
			deviceRegex := regexp.MustCompile(`(/dev/drbd\d+)`)
			if matches := deviceRegex.FindStringSubmatch(line); len(matches) > 0 {
				status.Device = matches[1]
			}
		}

		// Extract connection state
		if strings.Contains(line, "Connected") {
			status.ConnectionState = "Connected"
		} else if strings.Contains(line, "Disconnected") {
			status.ConnectionState = "Disconnected"
		} else if strings.Contains(line, "StandAlone") {
			status.ConnectionState = "StandAlone"
		}

		// Extract role (Primary/Secondary)
		if strings.Contains(line, "Primary") {
			status.Role = "Primary"
		} else if strings.Contains(line, "Secondary") {
			status.Role = "Secondary"
		}

		// Extract disk state
		if strings.Contains(line, "UpToDate") {
			status.DiskState = "UpToDate"
		} else if strings.Contains(line, "Inconsistent") {
			status.DiskState = "Inconsistent"
		}

		// Extract sync progress
		if strings.Contains(line, "sync'ed") {
			status.Resyncing = true
			progressRegex := regexp.MustCompile(`(\d+\.\d+)%`)
			if matches := progressRegex.FindStringSubmatch(line); len(matches) > 0 {
				fmt.Sscanf(matches[1], "%f", &status.SyncProgress)
			}
		}
	}

	return status, nil
}

// ListResources lists all DRBD resources
func (dm *DRBDManager) ListResources() ([]string, error) {
	if !dm.enabled {
		return nil, fmt.Errorf("DRBD is not enabled")
	}

	// List all .res files in /etc/drbd.d/
	result, err := dm.shell.Execute("sudo", "ls", "/etc/drbd.d/")
	if err != nil {
		return nil, fmt.Errorf("failed to list DRBD resources: %s: %w", result.Stderr, err)
	}

	resources := []string{}
	lines := strings.Split(result.Stdout, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, ".res") {
			resourceName := strings.TrimSuffix(line, ".res")
			resources = append(resources, resourceName)
		}
	}

	return resources, nil
}

// PromoteToPrimary promotes a DRBD resource to primary role
func (dm *DRBDManager) PromoteToPrimary(name string) error {
	if !dm.enabled {
		return fmt.Errorf("DRBD is not enabled")
	}

	result, err := dm.shell.Execute("sudo", "drbdadm", "primary", name)
	if err != nil {
		return fmt.Errorf("failed to promote to primary: %s: %w", result.Stderr, err)
	}

	logger.Info("DRBD resource promoted to primary", zap.String("name", name))
	return nil
}

// DemoteToSecondary demotes a DRBD resource to secondary role
func (dm *DRBDManager) DemoteToSecondary(name string) error {
	if !dm.enabled {
		return fmt.Errorf("DRBD is not enabled")
	}

	result, err := dm.shell.Execute("sudo", "drbdadm", "secondary", name)
	if err != nil {
		return fmt.Errorf("failed to demote to secondary: %s: %w", result.Stderr, err)
	}

	logger.Info("DRBD resource demoted to secondary", zap.String("name", name))
	return nil
}

// ForcePrimary forces a DRBD resource to become primary (split-brain recovery)
func (dm *DRBDManager) ForcePrimary(name string) error {
	if !dm.enabled {
		return fmt.Errorf("DRBD is not enabled")
	}

	// First set the resource to primary with --force
	result, err := dm.shell.Execute("sudo", "drbdadm", "primary", "--force", name)
	if err != nil {
		return fmt.Errorf("failed to force primary: %s: %w", result.Stderr, err)
	}

	logger.Warn("DRBD resource forced to primary - check for split-brain", zap.String("name", name))
	return nil
}

// Disconnect disconnects a DRBD resource from its peer
func (dm *DRBDManager) Disconnect(name string) error {
	if !dm.enabled {
		return fmt.Errorf("DRBD is not enabled")
	}

	result, err := dm.shell.Execute("sudo", "drbdadm", "disconnect", name)
	if err != nil {
		return fmt.Errorf("failed to disconnect: %s: %w", result.Stderr, err)
	}

	logger.Info("DRBD resource disconnected", zap.String("name", name))
	return nil
}

// Connect connects a DRBD resource to its peer
func (dm *DRBDManager) Connect(name string) error {
	if !dm.enabled {
		return fmt.Errorf("DRBD is not enabled")
	}

	result, err := dm.shell.Execute("sudo", "drbdadm", "connect", name)
	if err != nil {
		return fmt.Errorf("failed to connect: %s: %w", result.Stderr, err)
	}

	logger.Info("DRBD resource connected", zap.String("name", name))
	return nil
}

// StartSync starts synchronization for a DRBD resource
func (dm *DRBDManager) StartSync(name string) error {
	if !dm.enabled {
		return fmt.Errorf("DRBD is not enabled")
	}

	// Invalidate to force full sync
	result, err := dm.shell.Execute("sudo", "drbdadm", "invalidate", name)
	if err != nil {
		return fmt.Errorf("failed to start sync: %s: %w", result.Stderr, err)
	}

	logger.Info("DRBD synchronization started", zap.String("name", name))
	return nil
}

// VerifyData verifies data integrity of a DRBD resource
func (dm *DRBDManager) VerifyData(name string) error {
	if !dm.enabled {
		return fmt.Errorf("DRBD is not enabled")
	}

	result, err := dm.shell.Execute("sudo", "drbdadm", "verify", name)
	if err != nil {
		return fmt.Errorf("failed to verify data: %s: %w", result.Stderr, err)
	}

	logger.Info("DRBD data verification started", zap.String("name", name))
	return nil
}

// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package sharing

import (
	"fmt"
)

// ISCSIManager manages iSCSI targets
type ISCSIManager struct {
	shell   ShellExecutor
	enabled bool
}

// ISCSITarget represents an iSCSI target configuration
type ISCSITarget struct {
	Name        string   `json:"name"`
	IQN         string   `json:"iqn"`
	LUN         int      `json:"lun"`
	Path        string   `json:"path"`
	Size        uint64   `json:"size"`
	AllowedIPs  []string `json:"allowed_ips"`
	CHAP        bool     `json:"chap"`
	CHAPUser    string   `json:"chap_user"`
}

// NewISCSIManager creates a new iSCSI manager
func NewISCSIManager(shell ShellExecutor) (*ISCSIManager, error) {
	// Check for targetcli (LIO - modern Linux iSCSI target)
	if !shell.CommandExists("targetcli") {
		return nil, fmt.Errorf("targetcli not installed (install targetcli-fb)")
	}

	return &ISCSIManager{
		shell:   shell,
		enabled: true,
	}, nil
}

// IsEnabled returns whether iSCSI is available
func (i *ISCSIManager) IsEnabled() bool {
	return i.enabled
}

// GetStatus gets iSCSI target service status
func (i *ISCSIManager) GetStatus() (bool, error) {
	result, err := i.shell.Execute("systemctl", "is-active", "target")
	if err != nil {
		return false, nil
	}

	return result.Stdout == "active", nil
}

// Start starts the iSCSI target service
func (i *ISCSIManager) Start() error {
	_, err := i.shell.Execute("systemctl", "start", "target")
	if err != nil {
		return fmt.Errorf("failed to start iSCSI target: %w", err)
	}

	return nil
}

// Stop stops the iSCSI target service
func (i *ISCSIManager) Stop() error {
	_, err := i.shell.Execute("systemctl", "stop", "target")
	if err != nil {
		return fmt.Errorf("failed to stop iSCSI target: %w", err)
	}

	return nil
}

// ListTargets lists all iSCSI targets
func (i *ISCSIManager) ListTargets() ([]ISCSITarget, error) {
	result, err := i.shell.Execute("targetcli", "ls")
	if err != nil {
		return nil, fmt.Errorf("failed to list targets: %w", err)
	}

	// TODO: Parse targetcli output
	// This is a simplified implementation
	_ = result

	return []ISCSITarget{}, nil
}

// CreateTarget creates a new iSCSI target
func (i *ISCSIManager) CreateTarget(target ISCSITarget) error {
	// Example command sequence:
	// targetcli /backstores/fileio create disk01 /mnt/disk01.img 10G
	// targetcli /iscsi create iqn.2023-01.com.example:target01
	// targetcli /iscsi/iqn.2023-01.com.example:target01/tpg1/luns create /backstores/fileio/disk01
	// targetcli /iscsi/iqn.2023-01.com.example:target01/tpg1/acls create iqn.2023-01.com.example:client01

	// Create backing store
	_, err := i.shell.Execute("targetcli", "/backstores/fileio", "create",
		target.Name, target.Path, fmt.Sprintf("%d", target.Size))
	if err != nil {
		return fmt.Errorf("failed to create backstore: %w", err)
	}

	// Create target
	_, err = i.shell.Execute("targetcli", "/iscsi", "create", target.IQN)
	if err != nil {
		return fmt.Errorf("failed to create target: %w", err)
	}

	// Create LUN
	lunPath := fmt.Sprintf("/iscsi/%s/tpg1/luns", target.IQN)
	_, err = i.shell.Execute("targetcli", lunPath, "create", fmt.Sprintf("/backstores/fileio/%s", target.Name))
	if err != nil {
		return fmt.Errorf("failed to create LUN: %w", err)
	}

	// Save configuration
	_, _ = i.shell.Execute("targetcli", "saveconfig")

	return nil
}

// DeleteTarget deletes an iSCSI target
func (i *ISCSIManager) DeleteTarget(iqn string) error {
	_, err := i.shell.Execute("targetcli", "/iscsi", "delete", iqn)
	if err != nil {
		return fmt.Errorf("failed to delete target: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")

	return nil
}

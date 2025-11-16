// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package sharing

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"fmt"
)

// ISCSIManager manages iSCSI targets
type ISCSIManager struct {
	shell      executor.ShellExecutor
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
func NewISCSIManager(shell executor.ShellExecutor) (*ISCSIManager, error) {
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

// ===== Advanced iSCSI Features =====

// ISCSIInitiator represents an iSCSI initiator (client)
type ISCSIInitiator struct {
	IQN         string `json:"iqn"`
	IPAddress   string `json:"ip_address"`
	CHAPEnabled bool   `json:"chap_enabled"`
	CHAPUser    string `json:"chap_user"`
}

// ISCSILUN represents a Logical Unit Number
type ISCSILUN struct {
	Number        int    `json:"number"`
	BackstoreName string `json:"backstore_name"`
	Size          uint64 `json:"size"`
	Path          string `json:"path"`
	WriteBack     bool   `json:"write_back"`
}

// ISCSIPortal represents a network portal
type ISCSIPortal struct {
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
}

// AddACL adds an ACL (allowed initiator) to a target
func (i *ISCSIManager) AddACL(targetIQN string, initiatorIQN string) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	aclPath := fmt.Sprintf("/iscsi/%s/tpg1/acls", targetIQN)
	_, err := i.shell.Execute("targetcli", aclPath, "create", initiatorIQN)
	if err != nil {
		return fmt.Errorf("failed to add ACL: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// RemoveACL removes an ACL from a target
func (i *ISCSIManager) RemoveACL(targetIQN string, initiatorIQN string) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	aclPath := fmt.Sprintf("/iscsi/%s/tpg1/acls", targetIQN)
	_, err := i.shell.Execute("targetcli", aclPath, "delete", initiatorIQN)
	if err != nil {
		return fmt.Errorf("failed to remove ACL: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// SetCHAPAuth configures CHAP authentication for an ACL
func (i *ISCSIManager) SetCHAPAuth(targetIQN string, initiatorIQN string, username string, password string) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	aclPath := fmt.Sprintf("/iscsi/%s/tpg1/acls/%s", targetIQN, initiatorIQN)

	// Set CHAP username
	_, err := i.shell.Execute("targetcli", aclPath, "set", "auth", "userid="+username)
	if err != nil {
		return fmt.Errorf("failed to set CHAP username: %w", err)
	}

	// Set CHAP password
	_, err = i.shell.Execute("targetcli", aclPath, "set", "auth", "password="+password)
	if err != nil {
		return fmt.Errorf("failed to set CHAP password: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// DisableCHAPAuth disables CHAP authentication for an ACL
func (i *ISCSIManager) DisableCHAPAuth(targetIQN string, initiatorIQN string) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	aclPath := fmt.Sprintf("/iscsi/%s/tpg1/acls/%s", targetIQN, initiatorIQN)
	_, err := i.shell.Execute("targetcli", aclPath, "set", "auth", "userid=")
	if err != nil {
		return fmt.Errorf("failed to disable CHAP: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// AddPortal adds a network portal to a target
func (i *ISCSIManager) AddPortal(targetIQN string, ipAddress string, port int) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	portalPath := fmt.Sprintf("/iscsi/%s/tpg1/portals", targetIQN)
	portalSpec := fmt.Sprintf("%s:%d", ipAddress, port)

	_, err := i.shell.Execute("targetcli", portalPath, "create", portalSpec)
	if err != nil {
		return fmt.Errorf("failed to add portal: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// RemovePortal removes a network portal from a target
func (i *ISCSIManager) RemovePortal(targetIQN string, ipAddress string, port int) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	portalPath := fmt.Sprintf("/iscsi/%s/tpg1/portals", targetIQN)
	portalSpec := fmt.Sprintf("%s:%d", ipAddress, port)

	_, err := i.shell.Execute("targetcli", portalPath, "delete", portalSpec)
	if err != nil {
		return fmt.Errorf("failed to remove portal: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// CreateBlockBackstore creates a block device backstore
func (i *ISCSIManager) CreateBlockBackstore(name string, devicePath string) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	_, err := i.shell.Execute("targetcli", "/backstores/block", "create", name, devicePath)
	if err != nil {
		return fmt.Errorf("failed to create block backstore: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// CreateFileIOBackstore creates a file-based backstore
func (i *ISCSIManager) CreateFileIOBackstore(name string, filePath string, size uint64, writeBack bool) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	args := []string{"/backstores/fileio", "create", name, filePath, fmt.Sprintf("%d", size)}
	if writeBack {
		args = append(args, "write_back=true")
	}

	_, err := i.shell.Execute("targetcli", args...)
	if err != nil {
		return fmt.Errorf("failed to create fileio backstore: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// DeleteBackstore deletes a backstore
func (i *ISCSIManager) DeleteBackstore(backstoreType string, name string) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	path := fmt.Sprintf("/backstores/%s", backstoreType)
	_, err := i.shell.Execute("targetcli", path, "delete", name)
	if err != nil {
		return fmt.Errorf("failed to delete backstore: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// AddLUN adds a LUN to a target
func (i *ISCSIManager) AddLUN(targetIQN string, backstorePath string, lunNumber int) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	lunPath := fmt.Sprintf("/iscsi/%s/tpg1/luns", targetIQN)
	args := []string{lunPath, "create", backstorePath}
	if lunNumber >= 0 {
		args = append(args, fmt.Sprintf("lun=%d", lunNumber))
	}

	_, err := i.shell.Execute("targetcli", args...)
	if err != nil {
		return fmt.Errorf("failed to add LUN: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// RemoveLUN removes a LUN from a target
func (i *ISCSIManager) RemoveLUN(targetIQN string, lunNumber int) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	lunPath := fmt.Sprintf("/iscsi/%s/tpg1/luns", targetIQN)
	_, err := i.shell.Execute("targetcli", lunPath, "delete", fmt.Sprintf("lun%d", lunNumber))
	if err != nil {
		return fmt.Errorf("failed to remove LUN: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// SetTargetAttribute sets an attribute on a target
func (i *ISCSIManager) SetTargetAttribute(targetIQN string, attribute string, value string) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	tpgPath := fmt.Sprintf("/iscsi/%s/tpg1", targetIQN)
	_, err := i.shell.Execute("targetcli", tpgPath, "set", "attribute", fmt.Sprintf("%s=%s", attribute, value))
	if err != nil {
		return fmt.Errorf("failed to set attribute: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// EnableTarget enables a target portal group
func (i *ISCSIManager) EnableTarget(targetIQN string) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	tpgPath := fmt.Sprintf("/iscsi/%s/tpg1", targetIQN)
	_, err := i.shell.Execute("targetcli", tpgPath, "enable")
	if err != nil {
		return fmt.Errorf("failed to enable target: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// DisableTarget disables a target portal group
func (i *ISCSIManager) DisableTarget(targetIQN string) error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	tpgPath := fmt.Sprintf("/iscsi/%s/tpg1", targetIQN)
	_, err := i.shell.Execute("targetcli", tpgPath, "disable")
	if err != nil {
		return fmt.Errorf("failed to disable target: %w", err)
	}

	_, _ = i.shell.Execute("targetcli", "saveconfig")
	return nil
}

// GetSessions returns active iSCSI sessions
func (i *ISCSIManager) GetSessions() (string, error) {
	if !i.enabled {
		return "", fmt.Errorf("iSCSI not available")
	}

	result, err := i.shell.Execute("targetcli", "sessions")
	if err != nil {
		return "", fmt.Errorf("failed to get sessions: %w", err)
	}

	return result.Stdout, nil
}

// SaveConfig saves the current iSCSI configuration
func (i *ISCSIManager) SaveConfig() error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	_, err := i.shell.Execute("targetcli", "saveconfig")
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// RestoreConfig restores iSCSI configuration from file
func (i *ISCSIManager) RestoreConfig() error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	_, err := i.shell.Execute("targetcli", "restoreconfig")
	if err != nil {
		return fmt.Errorf("failed to restore config: %w", err)
	}

	return nil
}

// ClearConfig clears all iSCSI configuration
func (i *ISCSIManager) ClearConfig() error {
	if !i.enabled {
		return fmt.Errorf("iSCSI not available")
	}

	_, err := i.shell.Execute("targetcli", "clearconfig", "confirm=True")
	if err != nil {
		return fmt.Errorf("failed to clear config: %w", err)
	}

	return nil
}

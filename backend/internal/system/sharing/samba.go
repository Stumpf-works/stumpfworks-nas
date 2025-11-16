// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package sharing

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"fmt"
	"os"
	"strings"
)



// SambaManager manages Samba/SMB shares
type SambaManager struct {
	shell      executor.ShellExecutor
	enabled    bool
	configPath string
}

// SambaShare represents a Samba share configuration
type SambaShare struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Comment     string   `json:"comment"`
	ReadOnly    bool     `json:"read_only"`
	Browseable  bool     `json:"browseable"`
	GuestOK     bool     `json:"guest_ok"`
	ValidUsers  []string `json:"valid_users"`
	ValidGroups []string `json:"valid_groups"`
	WritableUsers []string `json:"writable_users"`
	CreateMask  string   `json:"create_mask"`
	DirectoryMask string `json:"directory_mask"`
	VetoFiles   []string `json:"veto_files"`
	RecycleBin  bool     `json:"recycle_bin"`
}

// SambaUser represents a Samba user
type SambaUser struct {
	Username string `json:"username"`
	UID      int    `json:"uid"`
	Enabled  bool   `json:"enabled"`
}

// NewSambaManager creates a new Samba manager
func NewSambaManager(shell executor.ShellExecutor) (*SambaManager, error) {
	if !shell.CommandExists("smbd") {
		return nil, fmt.Errorf("samba not installed")
	}

	return &SambaManager{
		shell:      shell,
		enabled:    true,
		configPath: "/etc/samba/smb.conf",
	}, nil
}

// IsEnabled returns whether Samba is available
func (s *SambaManager) IsEnabled() bool {
	return s.enabled
}

// GetStatus gets Samba service status
func (s *SambaManager) GetStatus() (bool, error) {
	result, err := s.shell.Execute("systemctl", "is-active", "smbd")
	if err != nil {
		return false, nil
	}

	return strings.TrimSpace(result.Stdout) == "active", nil
}

// Start starts the Samba service
func (s *SambaManager) Start() error {
	_, err := s.shell.Execute("systemctl", "start", "smbd")
	if err != nil {
		return fmt.Errorf("failed to start samba: %w", err)
	}

	// Also start nmbd for NetBIOS
	_, _ = s.shell.Execute("systemctl", "start", "nmbd")

	return nil
}

// Stop stops the Samba service
func (s *SambaManager) Stop() error {
	_, err := s.shell.Execute("systemctl", "stop", "smbd")
	if err != nil {
		return fmt.Errorf("failed to stop samba: %w", err)
	}

	_, _ = s.shell.Execute("systemctl", "stop", "nmbd")

	return nil
}

// Restart restarts the Samba service
func (s *SambaManager) Restart() error {
	_, err := s.shell.Execute("systemctl", "restart", "smbd")
	if err != nil {
		return fmt.Errorf("failed to restart samba: %w", err)
	}

	_, _ = s.shell.Execute("systemctl", "restart", "nmbd")

	return nil
}

// Reload reloads Samba configuration without restarting
func (s *SambaManager) Reload() error {
	_, err := s.shell.Execute("systemctl", "reload", "smbd")
	if err != nil {
		return fmt.Errorf("failed to reload samba: %w", err)
	}

	return nil
}

// TestConfig tests Samba configuration for errors
func (s *SambaManager) TestConfig() error {
	result, err := s.shell.Execute("testparm", "-s", s.configPath)
	if err != nil {
		return fmt.Errorf("samba config test failed: %s", result.Stderr)
	}

	return nil
}

// ListShares lists all configured Samba shares
func (s *SambaManager) ListShares() ([]SambaShare, error) {
	result, err := s.shell.Execute("testparm", "-s", "--section-name")
	if err != nil {
		return nil, fmt.Errorf("failed to list shares: %w", err)
	}

	var shares []SambaShare
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		shareName := strings.TrimSpace(line)

		// Skip global and special sections
		if shareName == "" || shareName == "global" || shareName == "printers" {
			continue
		}

		share, err := s.GetShare(shareName)
		if err != nil {
			continue
		}

		shares = append(shares, *share)
	}

	return shares, nil
}

// GetShare gets configuration for a specific share
func (s *SambaManager) GetShare(name string) (*SambaShare, error) {
	result, err := s.shell.Execute("testparm", "-s", "-d", "0", "--section-name="+name, s.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get share: %w", err)
	}

	share := &SambaShare{
		Name:       name,
		Browseable: true, // default
	}

	lines := strings.Split(result.Stdout, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "path = ") {
			share.Path = strings.TrimPrefix(line, "path = ")
		} else if strings.HasPrefix(line, "comment = ") {
			share.Comment = strings.TrimPrefix(line, "comment = ")
		} else if strings.HasPrefix(line, "read only = ") {
			share.ReadOnly = strings.TrimPrefix(line, "read only = ") == "Yes"
		} else if strings.HasPrefix(line, "browseable = ") {
			share.Browseable = strings.TrimPrefix(line, "browseable = ") == "Yes"
		} else if strings.HasPrefix(line, "guest ok = ") {
			share.GuestOK = strings.TrimPrefix(line, "guest ok = ") == "Yes"
		} else if strings.HasPrefix(line, "valid users = ") {
			users := strings.TrimPrefix(line, "valid users = ")
			share.ValidUsers = strings.Fields(users)
		} else if strings.HasPrefix(line, "create mask = ") {
			share.CreateMask = strings.TrimPrefix(line, "create mask = ")
		} else if strings.HasPrefix(line, "directory mask = ") {
			share.DirectoryMask = strings.TrimPrefix(line, "directory mask = ")
		}
	}

	return share, nil
}

// CreateShare creates a new Samba share
func (s *SambaManager) CreateShare(share SambaShare) error {
	// Read current config
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	config := string(data)

	// Check if share already exists
	if strings.Contains(config, fmt.Sprintf("[%s]", share.Name)) {
		return fmt.Errorf("share already exists: %s", share.Name)
	}

	// Build share configuration
	shareConfig := fmt.Sprintf("\n[%s]\n", share.Name)
	shareConfig += fmt.Sprintf("  path = %s\n", share.Path)

	if share.Comment != "" {
		shareConfig += fmt.Sprintf("  comment = %s\n", share.Comment)
	}

	if share.ReadOnly {
		shareConfig += "  read only = yes\n"
	} else {
		shareConfig += "  read only = no\n"
	}

	if share.Browseable {
		shareConfig += "  browseable = yes\n"
	} else {
		shareConfig += "  browseable = no\n"
	}

	if share.GuestOK {
		shareConfig += "  guest ok = yes\n"
	} else {
		shareConfig += "  guest ok = no\n"
	}

	if len(share.ValidUsers) > 0 {
		shareConfig += fmt.Sprintf("  valid users = %s\n", strings.Join(share.ValidUsers, " "))
	}

	if len(share.ValidGroups) > 0 {
		groups := make([]string, len(share.ValidGroups))
		for i, g := range share.ValidGroups {
			groups[i] = "@" + g
		}
		shareConfig += fmt.Sprintf("  valid users = %s\n", strings.Join(groups, " "))
	}

	if len(share.WritableUsers) > 0 {
		shareConfig += fmt.Sprintf("  write list = %s\n", strings.Join(share.WritableUsers, " "))
	}

	if share.CreateMask != "" {
		shareConfig += fmt.Sprintf("  create mask = %s\n", share.CreateMask)
	}

	if share.DirectoryMask != "" {
		shareConfig += fmt.Sprintf("  directory mask = %s\n", share.DirectoryMask)
	}

	if share.RecycleBin {
		shareConfig += "  vfs objects = recycle\n"
		shareConfig += "  recycle:repository = .recycle\n"
		shareConfig += "  recycle:keeptree = yes\n"
		shareConfig += "  recycle:versions = yes\n"
	}

	if len(share.VetoFiles) > 0 {
		shareConfig += fmt.Sprintf("  veto files = /%s/\n", strings.Join(share.VetoFiles, "/"))
	}

	// Append to config
	config += shareConfig

	// Write back
	err = os.WriteFile(s.configPath, []byte(config), 0644)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Test config
	if err := s.TestConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Reload Samba
	return s.Reload()
}

// UpdateShare updates an existing share
func (s *SambaManager) UpdateShare(share SambaShare) error {
	// Delete and recreate
	if err := s.DeleteShare(share.Name); err != nil {
		return err
	}

	return s.CreateShare(share)
}

// DeleteShare deletes a Samba share
func (s *SambaManager) DeleteShare(name string) error {
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	inShare := false
	shareFound := false

	for _, line := range lines {
		// Check if this is the start of our share
		if strings.TrimSpace(line) == fmt.Sprintf("[%s]", name) {
			inShare = true
			shareFound = true
			continue
		}

		// Check if this is the start of another share
		if strings.HasPrefix(strings.TrimSpace(line), "[") && strings.HasSuffix(strings.TrimSpace(line), "]") {
			inShare = false
		}

		// Skip lines that are part of our share
		if inShare {
			continue
		}

		newLines = append(newLines, line)
	}

	if !shareFound {
		return fmt.Errorf("share not found: %s", name)
	}

	// Write back
	config := strings.Join(newLines, "\n")
	err = os.WriteFile(s.configPath, []byte(config), 0644)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Reload Samba
	return s.Reload()
}

// AddUser adds a Samba user
func (s *SambaManager) AddUser(username string, password string) error {
	// User must exist in system first
	_, err := s.shell.Execute("id", username)
	if err != nil {
		return fmt.Errorf("system user does not exist: %s", username)
	}

	// Add to Samba
	_, err = s.shell.Execute("bash", "-c",
		fmt.Sprintf("(echo %s; echo %s) | smbpasswd -a -s %s", password, password, username))
	if err != nil {
		return fmt.Errorf("failed to add samba user: %w", err)
	}

	// Enable the user
	return s.EnableUser(username)
}

// EnableUser enables a Samba user
func (s *SambaManager) EnableUser(username string) error {
	_, err := s.shell.Execute("smbpasswd", "-e", username)
	if err != nil {
		return fmt.Errorf("failed to enable user: %w", err)
	}

	return nil
}

// DisableUser disables a Samba user
func (s *SambaManager) DisableUser(username string) error {
	_, err := s.shell.Execute("smbpasswd", "-d", username)
	if err != nil {
		return fmt.Errorf("failed to disable user: %w", err)
	}

	return nil
}

// DeleteUser deletes a Samba user
func (s *SambaManager) DeleteUser(username string) error {
	_, err := s.shell.Execute("smbpasswd", "-x", username)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// SetUserPassword sets a Samba user's password
func (s *SambaManager) SetUserPassword(username string, password string) error {
	_, err := s.shell.Execute("bash", "-c",
		fmt.Sprintf("(echo %s; echo %s) | smbpasswd -s %s", password, password, username))
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	return nil
}

// ListUsers lists all Samba users
func (s *SambaManager) ListUsers() ([]SambaUser, error) {
	result, err := s.shell.Execute("pdbedit", "-L")
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	var users []SambaUser
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Format: username:UID:...
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		user := SambaUser{
			Username: parts[0],
			Enabled:  true,
		}

		users = append(users, user)
	}

	return users, nil
}

// GetConnections gets current Samba connections
func (s *SambaManager) GetConnections() ([]string, error) {
	result, err := s.shell.Execute("smbstatus", "-b")
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %w", err)
	}

	var connections []string
	lines := strings.Split(result.Stdout, "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "Samba") || strings.HasPrefix(line, "PID") {
			continue
		}

		connections = append(connections, strings.TrimSpace(line))
	}

	return connections, nil
}

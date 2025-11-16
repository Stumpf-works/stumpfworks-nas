// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package users

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"fmt"
	"strings"
)

// ADManager manages Active Directory integration
type ADManager struct {
	shell      executor.ShellExecutor
	enabled bool
}

// ADConfig represents Active Directory configuration
type ADConfig struct {
	Domain       string `json:"domain"`
	Server       string `json:"server"`
	Workgroup    string `json:"workgroup"`
	Administrator string `json:"administrator"`
	Password     string `json:"password"`
	OU           string `json:"ou"` // Organizational Unit
}

// ADUser represents an Active Directory user
type ADUser struct {
	SamAccountName string   `json:"sam_account_name"`
	DN             string   `json:"dn"`
	Email          string   `json:"email"`
	DisplayName    string   `json:"display_name"`
	Groups         []string `json:"groups"`
}

// NewADManager creates a new Active Directory manager
func NewADManager(shell executor.ShellExecutor) (*ADManager, error) {
	// Check if required tools are available (Samba + winbind/sssd)
	if !shell.CommandExists("net") {
		return nil, fmt.Errorf("Samba tools not installed")
	}

	return &ADManager{
		shell:   shell,
		enabled: true,
	}, nil
}

// IsEnabled returns whether AD integration is available
func (a *ADManager) IsEnabled() bool {
	return a.enabled
}

// GetStatus gets AD join status
func (a *ADManager) GetStatus() (bool, error) {
	result, err := a.shell.Execute("net", "ads", "testjoin")
	if err != nil {
		return false, nil
	}

	return strings.Contains(result.Stdout, "Join is OK"), nil
}

// JoinDomain joins the Active Directory domain
func (a *ADManager) JoinDomain(config ADConfig) error {
	// This requires:
	// 1. Configured /etc/krb5.conf
	// 2. Configured /etc/samba/smb.conf for AD
	// 3. Time sync with AD (NTP)

	args := []string{
		"ads", "join",
		"-U", config.Administrator,
	}

	if config.OU != "" {
		args = append(args, "-O", config.OU)
	}

	// Execute net ads join
	_, err := a.shell.Execute("bash", "-c",
		fmt.Sprintf("echo '%s' | net %s", config.Password, strings.Join(args, " ")))
	if err != nil {
		return fmt.Errorf("failed to join domain: %w", err)
	}

	// Start winbind or sssd
	_, _ = a.shell.Execute("systemctl", "start", "winbind")
	_, _ = a.shell.Execute("systemctl", "enable", "winbind")

	return nil
}

// LeaveDomain leaves the Active Directory domain
func (a *ADManager) LeaveDomain(config ADConfig) error {
	args := []string{
		"ads", "leave",
		"-U", config.Administrator,
	}

	_, err := a.shell.Execute("bash", "-c",
		fmt.Sprintf("echo '%s' | net %s", config.Password, strings.Join(args, " ")))
	if err != nil {
		return fmt.Errorf("failed to leave domain: %w", err)
	}

	// Stop winbind
	_, _ = a.shell.Execute("systemctl", "stop", "winbind")
	_, _ = a.shell.Execute("systemctl", "disable", "winbind")

	return nil
}

// ListADUsers lists users from Active Directory
func (a *ADManager) ListADUsers() ([]ADUser, error) {
	result, err := a.shell.Execute("wbinfo", "-u")
	if err != nil {
		return nil, fmt.Errorf("failed to list AD users: %w", err)
	}

	var users []ADUser
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Format: DOMAIN\username
		user := ADUser{
			SamAccountName: line,
		}

		users = append(users, user)
	}

	return users, nil
}

// ListADGroups lists groups from Active Directory
func (a *ADManager) ListADGroups() ([]string, error) {
	result, err := a.shell.Execute("wbinfo", "-g")
	if err != nil {
		return nil, fmt.Errorf("failed to list AD groups: %w", err)
	}

	var groups []string
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		if line != "" {
			groups = append(groups, line)
		}
	}

	return groups, nil
}

// GetUserInfo gets information about an AD user
func (a *ADManager) GetUserInfo(username string) (*ADUser, error) {
	result, err := a.shell.Execute("wbinfo", "--user-info", username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Parse wbinfo output
	// Format: DOMAIN\username:*:uid:gid:Full Name:/home/username:/bin/bash
	fields := strings.Split(result.Stdout, ":")
	if len(fields) < 5 {
		return nil, fmt.Errorf("invalid user info format")
	}

	user := &ADUser{
		SamAccountName: fields[0],
		DisplayName:    fields[4],
	}

	return user, nil
}

// TestAuthentication tests authentication against AD
func (a *ADManager) TestAuthentication(username string, password string) error {
	_, err := a.shell.Execute("bash", "-c",
		fmt.Sprintf("echo '%s' | wbinfo -a '%s'", password, username))
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	return nil
}

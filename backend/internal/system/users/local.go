// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package users

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"fmt"
	"os/user"
	"strconv"
	"strings"
)



// LocalManager manages local Unix users and groups
type LocalManager struct {
	shell      executor.ShellExecutor
	enabled bool
}

// LocalUser represents a local Unix user
type LocalUser struct {
	Username string   `json:"username"`
	UID      int      `json:"uid"`
	GID      int      `json:"gid"`
	Home     string   `json:"home"`
	Shell    string   `json:"shell"`
	FullName string   `json:"full_name"`
	Groups   []string `json:"groups"`
	Locked   bool     `json:"locked"`
}

// LocalGroup represents a local Unix group
type LocalGroup struct {
	Name    string   `json:"name"`
	GID     int      `json:"gid"`
	Members []string `json:"members"`
}

// NewLocalManager creates a new local user manager
func NewLocalManager(shell executor.ShellExecutor) (*LocalManager, error) {
	return &LocalManager{
		shell:   shell,
		enabled: true,
	}, nil
}

// IsEnabled returns whether local user management is available
func (l *LocalManager) IsEnabled() bool {
	return l.enabled
}

// ListUsers lists all local users
func (l *LocalManager) ListUsers() ([]LocalUser, error) {
	result, err := l.shell.Execute("getent", "passwd")
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	var users []LocalUser
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Format: username:x:uid:gid:fullname:home:shell
		fields := strings.Split(line, ":")
		if len(fields) < 7 {
			continue
		}

		uid, _ := strconv.Atoi(fields[2])
		gid, _ := strconv.Atoi(fields[3])

		// Skip system users (UID < 1000) unless it's root
		if uid < 1000 && uid != 0 {
			continue
		}

		localUser := LocalUser{
			Username: fields[0],
			UID:      uid,
			GID:      gid,
			FullName: fields[4],
			Home:     fields[5],
			Shell:    fields[6],
		}

		// Get user's groups
		groups, _ := l.GetUserGroups(localUser.Username)
		localUser.Groups = groups

		// Check if locked
		localUser.Locked = l.IsUserLocked(localUser.Username)

		users = append(users, localUser)
	}

	return users, nil
}

// GetUser gets details for a specific user
func (l *LocalManager) GetUser(username string) (*LocalUser, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", username)
	}

	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)

	localUser := &LocalUser{
		Username: u.Username,
		UID:      uid,
		GID:      gid,
		FullName: u.Name,
		Home:     u.HomeDir,
	}

	// Get groups
	groups, _ := l.GetUserGroups(username)
	localUser.Groups = groups

	// Get shell
	result, err := l.shell.Execute("getent", "passwd", username)
	if err == nil {
		fields := strings.Split(result.Stdout, ":")
		if len(fields) >= 7 {
			localUser.Shell = fields[6]
		}
	}

	// Check if locked
	localUser.Locked = l.IsUserLocked(username)

	return localUser, nil
}

// CreateUser creates a new local user
func (l *LocalManager) CreateUser(username string, password string, home string, shell string) error {
	args := []string{"-m"} // Create home directory

	if home != "" {
		args = append(args, "-d", home)
	}

	if shell != "" {
		args = append(args, "-s", shell)
	} else {
		args = append(args, "-s", "/bin/bash")
	}

	args = append(args, username)

	_, err := l.shell.Execute("useradd", args...)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Set password if provided
	if password != "" {
		if err := l.SetUserPassword(username, password); err != nil {
			return fmt.Errorf("user created but failed to set password: %w", err)
		}
	}

	return nil
}

// DeleteUser deletes a local user
func (l *LocalManager) DeleteUser(username string, removeHome bool) error {
	args := []string{}

	if removeHome {
		args = append(args, "-r")
	}

	args = append(args, username)

	_, err := l.shell.Execute("userdel", args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// SetUserPassword sets a user's password
func (l *LocalManager) SetUserPassword(username string, password string) error {
	_, err := l.shell.Execute("bash", "-c",
		fmt.Sprintf("echo '%s:%s' | chpasswd", username, password))
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	return nil
}

// LockUser locks a user account
func (l *LocalManager) LockUser(username string) error {
	_, err := l.shell.Execute("usermod", "-L", username)
	if err != nil {
		return fmt.Errorf("failed to lock user: %w", err)
	}

	return nil
}

// UnlockUser unlocks a user account
func (l *LocalManager) UnlockUser(username string) error {
	_, err := l.shell.Execute("usermod", "-U", username)
	if err != nil {
		return fmt.Errorf("failed to unlock user: %w", err)
	}

	return nil
}

// IsUserLocked checks if a user account is locked
func (l *LocalManager) IsUserLocked(username string) bool {
	result, err := l.shell.Execute("passwd", "-S", username)
	if err != nil {
		return false
	}

	// Output format: username L ... (L = locked, P = password set)
	fields := strings.Fields(result.Stdout)
	if len(fields) >= 2 {
		return fields[1] == "L"
	}

	return false
}

// AddUserToGroup adds a user to a group
func (l *LocalManager) AddUserToGroup(username string, group string) error {
	_, err := l.shell.Execute("usermod", "-aG", group, username)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}

	return nil
}

// RemoveUserFromGroup removes a user from a group
func (l *LocalManager) RemoveUserFromGroup(username string, group string) error {
	_, err := l.shell.Execute("gpasswd", "-d", username, group)
	if err != nil {
		return fmt.Errorf("failed to remove user from group: %w", err)
	}

	return nil
}

// GetUserGroups gets all groups a user belongs to
func (l *LocalManager) GetUserGroups(username string) ([]string, error) {
	result, err := l.shell.Execute("groups", username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	// Output format: username : group1 group2 group3
	parts := strings.Split(result.Stdout, ":")
	if len(parts) < 2 {
		return []string{}, nil
	}

	groups := strings.Fields(parts[1])
	return groups, nil
}

// SetUserShell sets a user's shell
func (l *LocalManager) SetUserShell(username string, shell string) error {
	_, err := l.shell.Execute("usermod", "-s", shell, username)
	if err != nil {
		return fmt.Errorf("failed to set shell: %w", err)
	}

	return nil
}

// SetUserHome sets a user's home directory
func (l *LocalManager) SetUserHome(username string, home string) error {
	_, err := l.shell.Execute("usermod", "-d", home, username)
	if err != nil {
		return fmt.Errorf("failed to set home directory: %w", err)
	}

	return nil
}

// ListGroups lists all local groups
func (l *LocalManager) ListGroups() ([]LocalGroup, error) {
	result, err := l.shell.Execute("getent", "group")
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	var groups []LocalGroup
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Format: groupname:x:gid:member1,member2
		fields := strings.Split(line, ":")
		if len(fields) < 4 {
			continue
		}

		gid, _ := strconv.Atoi(fields[2])

		// Skip system groups (GID < 1000) unless it's root
		if gid < 1000 && gid != 0 {
			continue
		}

		group := LocalGroup{
			Name: fields[0],
			GID:  gid,
		}

		// Parse members
		if fields[3] != "" {
			group.Members = strings.Split(fields[3], ",")
		}

		groups = append(groups, group)
	}

	return groups, nil
}

// GetGroup gets details for a specific group
func (l *LocalManager) GetGroup(name string) (*LocalGroup, error) {
	g, err := user.LookupGroup(name)
	if err != nil {
		return nil, fmt.Errorf("group not found: %s", name)
	}

	gid, _ := strconv.Atoi(g.Gid)

	group := &LocalGroup{
		Name: g.Name,
		GID:  gid,
	}

	// Get members
	result, err := l.shell.Execute("getent", "group", name)
	if err == nil {
		fields := strings.Split(strings.TrimSpace(result.Stdout), ":")
		if len(fields) >= 4 && fields[3] != "" {
			group.Members = strings.Split(fields[3], ",")
		}
	}

	return group, nil
}

// CreateGroup creates a new group
func (l *LocalManager) CreateGroup(name string, gid int) error {
	args := []string{}

	if gid > 0 {
		args = append(args, "-g", strconv.Itoa(gid))
	}

	args = append(args, name)

	_, err := l.shell.Execute("groupadd", args...)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}

	return nil
}

// DeleteGroup deletes a group
func (l *LocalManager) DeleteGroup(name string) error {
	_, err := l.shell.Execute("groupdel", name)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	return nil
}

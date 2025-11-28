// Revision: 2025-11-28 | Author: Claude | Version: 1.0.0
// Package filesystem provides filesystem ACL management
package filesystem

import (
	"fmt"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
)

// ACLManager manages POSIX ACLs on files and directories
type ACLManager struct {
	shell   executor.ShellExecutor
	enabled bool
}

// ACLEntry represents a single ACL entry
type ACLEntry struct {
	Type        string `json:"type"`        // user, group, mask, other
	Name        string `json:"name"`        // username or groupname (empty for owner/other)
	Permissions string `json:"permissions"` // rwx format (e.g., "rwx", "r-x", "---")
}

// ACLInfo represents complete ACL information for a file/directory
type ACLInfo struct {
	Path    string     `json:"path"`
	Entries []ACLEntry `json:"entries"`
}

// NewACLManager creates a new ACL manager
func NewACLManager(shell executor.ShellExecutor) (*ACLManager, error) {
	// Check if ACL tools are available
	if !shell.CommandExists("getfacl") || !shell.CommandExists("setfacl") {
		return nil, fmt.Errorf("ACL tools not installed (install 'acl' package)")
	}

	return &ACLManager{
		shell:   shell,
		enabled: true,
	}, nil
}

// IsEnabled returns whether ACL support is available
func (a *ACLManager) IsEnabled() bool {
	return a.enabled
}

// GetACL retrieves ACL entries for a file or directory
func (a *ACLManager) GetACL(path string) ([]ACLEntry, error) {
	if !a.enabled {
		return nil, fmt.Errorf("ACL support not available")
	}

	result, err := a.shell.Execute("getfacl", "--omit-header", "--numeric", "--absolute-names", path)
	if err != nil {
		return nil, fmt.Errorf("failed to get ACL for %s: %w", path, err)
	}

	return a.parseACLOutput(result.Stdout)
}

// SetACL sets ACL entries on a file or directory
func (a *ACLManager) SetACL(path string, entries []ACLEntry) error {
	if !a.enabled {
		return fmt.Errorf("ACL support not available")
	}

	if len(entries) == 0 {
		return fmt.Errorf("no ACL entries provided")
	}

	// Build setfacl arguments
	args := []string{"-m"}

	var aclStrings []string
	for _, entry := range entries {
		aclStr := fmt.Sprintf("%s:%s:%s", entry.Type, entry.Name, entry.Permissions)
		aclStrings = append(aclStrings, aclStr)
	}

	args = append(args, strings.Join(aclStrings, ","))
	args = append(args, path)

	result, err := a.shell.Execute("setfacl", args...)
	if err != nil {
		return fmt.Errorf("failed to set ACL on %s: %s - %w", path, result.Stderr, err)
	}

	return nil
}

// RemoveACL removes a specific ACL entry from a file or directory
func (a *ACLManager) RemoveACL(path string, entryType string, name string) error {
	if !a.enabled {
		return fmt.Errorf("ACL support not available")
	}

	aclSpec := fmt.Sprintf("%s:%s", entryType, name)

	result, err := a.shell.Execute("setfacl", "-x", aclSpec, path)
	if err != nil {
		return fmt.Errorf("failed to remove ACL entry %s from %s: %s - %w", aclSpec, path, result.Stderr, err)
	}

	return nil
}

// SetDefaultACL sets default ACL entries for new files created in a directory
func (a *ACLManager) SetDefaultACL(dirPath string, entries []ACLEntry) error {
	if !a.enabled {
		return fmt.Errorf("ACL support not available")
	}

	if len(entries) == 0 {
		return fmt.Errorf("no ACL entries provided")
	}

	// Build setfacl arguments with default: prefix
	args := []string{"-m"}

	var aclStrings []string
	for _, entry := range entries {
		aclStr := fmt.Sprintf("default:%s:%s:%s", entry.Type, entry.Name, entry.Permissions)
		aclStrings = append(aclStrings, aclStr)
	}

	args = append(args, strings.Join(aclStrings, ","))
	args = append(args, dirPath)

	result, err := a.shell.Execute("setfacl", args...)
	if err != nil {
		return fmt.Errorf("failed to set default ACL on %s: %s - %w", dirPath, result.Stderr, err)
	}

	return nil
}

// RemoveAllACLs removes all ACL entries from a file or directory
func (a *ACLManager) RemoveAllACLs(path string) error {
	if !a.enabled {
		return fmt.Errorf("ACL support not available")
	}

	result, err := a.shell.Execute("setfacl", "-b", path)
	if err != nil {
		return fmt.Errorf("failed to remove all ACLs from %s: %s - %w", path, result.Stderr, err)
	}

	return nil
}

// ApplyRecursive applies ACL entries recursively to a directory and all its contents
func (a *ACLManager) ApplyRecursive(dirPath string, entries []ACLEntry) error {
	if !a.enabled {
		return fmt.Errorf("ACL support not available")
	}

	if len(entries) == 0 {
		return fmt.Errorf("no ACL entries provided")
	}

	// Build setfacl arguments with -R flag
	args := []string{"-R", "-m"}

	var aclStrings []string
	for _, entry := range entries {
		aclStr := fmt.Sprintf("%s:%s:%s", entry.Type, entry.Name, entry.Permissions)
		aclStrings = append(aclStrings, aclStr)
	}

	args = append(args, strings.Join(aclStrings, ","))
	args = append(args, dirPath)

	result, err := a.shell.Execute("setfacl", args...)
	if err != nil {
		return fmt.Errorf("failed to apply ACLs recursively to %s: %s - %w", dirPath, result.Stderr, err)
	}

	return nil
}

// parseACLOutput parses getfacl output into ACLEntry structs
func (a *ACLManager) parseACLOutput(output string) ([]ACLEntry, error) {
	var entries []ACLEntry

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse ACL line format: type:name:permissions
		parts := strings.SplitN(line, ":", 3)
		if len(parts) != 3 {
			continue
		}

		entry := ACLEntry{
			Type:        parts[0],
			Name:        parts[1],
			Permissions: parts[2],
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

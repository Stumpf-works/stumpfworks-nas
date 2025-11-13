package sysutil

import (
	"os"
	"os/exec"
	"path/filepath"
)

// FindCommand searches for a command in common system paths
// Returns the full path to the executable if found, otherwise returns the original name
//
// Search order:
//  1. exec.LookPath() - checks $PATH environment variable
//  2. Common system paths: /usr/sbin, /sbin, /usr/bin, /bin, /usr/local/sbin, /usr/local/bin
//
// This is useful for finding system administration tools that may not be in $PATH
// for non-root users (e.g., useradd, userdel, smbpasswd, pdbedit)
func FindCommand(name string) string {
	// First try exec.LookPath (searches in PATH)
	if path, err := exec.LookPath(name); err == nil {
		return path
	}

	// Common system paths where admin tools are located
	systemPaths := []string{
		"/usr/sbin",      // System administration binaries (primary)
		"/sbin",          // Essential system binaries
		"/usr/bin",       // User binaries
		"/bin",           // Essential command binaries
		"/usr/local/sbin", // Locally installed system binaries
		"/usr/local/bin",  // Locally installed user binaries
	}

	for _, dir := range systemPaths {
		fullPath := filepath.Join(dir, name)
		if info, err := os.Stat(fullPath); err == nil {
			// Check if executable (has any execute bit set)
			if info.Mode()&0111 != 0 {
				return fullPath
			}
		}
	}

	// Return original name as fallback - will fail with proper error message
	return name
}

// CommandExists checks if a command exists and is executable
func CommandExists(name string) bool {
	path := FindCommand(name)
	if path == name {
		// FindCommand returned original name, try exec.LookPath
		_, err := exec.LookPath(name)
		return err == nil
	}
	// FindCommand found full path
	info, err := os.Stat(path)
	return err == nil && info.Mode()&0111 != 0
}

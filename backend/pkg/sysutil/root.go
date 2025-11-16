// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package sysutil

import (
	"os"
)

// IsRoot checks if the current process is running as root
func IsRoot() bool {
	return os.Geteuid() == 0
}

// RequireRoot returns an error if not running as root
func RequireRoot() error {
	if !IsRoot() {
		return ErrNotRoot
	}
	return nil
}

// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package system

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/users"
)

// UserManager manages all user and group operations
type UserManager struct {
	shell *ShellExecutor

	// Subsystems
	Local *users.LocalManager
	LDAP  *users.LDAPManager
	AD    *users.ADManager
}

// NewUserManager creates a new user manager
func NewUserManager(shell *ShellExecutor) (*UserManager, error) {
	um := &UserManager{
		shell: shell,
	}

	// Initialize local user manager
	local, err := users.NewLocalManager(shell)
	if err != nil {
		return nil, err
	}
	um.Local = local

	// Initialize LDAP manager (optional)
	ldap, err := users.NewLDAPManager(shell)
	if err != nil {
		um.LDAP = nil
	} else {
		um.LDAP = ldap
	}

	// Initialize AD manager (optional)
	ad, err := users.NewADManager(shell)
	if err != nil {
		um.AD = nil
	} else {
		um.AD = ad
	}

	return um, nil
}

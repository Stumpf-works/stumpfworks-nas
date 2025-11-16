// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package system

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/sharing"
)

// SharingManager manages all file sharing operations
type SharingManager struct {
	shell *ShellExecutor

	// Subsystems
	Samba  *sharing.SambaManager
	NFS    *sharing.NFSManager
	ISCSI  *sharing.ISCSIManager
	WebDAV *sharing.WebDAVManager
	FTP    *sharing.FTPManager
}

// NewSharingManager creates a new sharing manager
func NewSharingManager(shell *ShellExecutor) (*SharingManager, error) {
	sm := &SharingManager{
		shell: shell,
	}

	// Initialize Samba manager
	samba, err := sharing.NewSambaManager(shell)
	if err != nil {
		// Samba is optional but recommended
		sm.Samba = nil
	} else {
		sm.Samba = samba
	}

	// Initialize NFS manager
	nfs, err := sharing.NewNFSManager(shell)
	if err != nil {
		sm.NFS = nil
	} else {
		sm.NFS = nfs
	}

	// Initialize iSCSI manager
	iscsi, err := sharing.NewISCSIManager(shell)
	if err != nil {
		sm.ISCSI = nil
	} else {
		sm.ISCSI = iscsi
	}

	// Initialize WebDAV manager
	webdav, err := sharing.NewWebDAVManager(shell)
	if err != nil {
		sm.WebDAV = nil
	} else {
		sm.WebDAV = webdav
	}

	// Initialize FTP manager
	ftp, err := sharing.NewFTPManager(shell)
	if err != nil {
		sm.FTP = nil
	} else {
		sm.FTP = ftp
	}

	return sm, nil
}

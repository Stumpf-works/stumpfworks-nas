// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package sharing

import (
	"fmt"
	"strings"
)

// FTPManager manages FTP/SFTP server
type FTPManager struct {
	shell   ShellExecutor
	enabled bool
	backend string // "vsftpd", "proftpd", or "pure-ftpd"
}

// FTPUser represents an FTP user configuration
type FTPUser struct {
	Username  string `json:"username"`
	HomeDir   string `json:"home_dir"`
	ReadOnly  bool   `json:"read_only"`
	Anonymous bool   `json:"anonymous"`
}

// FTPConfig represents FTP server configuration
type FTPConfig struct {
	Port           int    `json:"port"`
	PasvMinPort    int    `json:"pasv_min_port"`
	PasvMaxPort    int    `json:"pasv_max_port"`
	AnonymousEnable bool  `json:"anonymous_enable"`
	TLSEnable      bool   `json:"tls_enable"`
	ChrootEnable   bool   `json:"chroot_enable"`
}

// NewFTPManager creates a new FTP manager
func NewFTPManager(shell ShellExecutor) (*FTPManager, error) {
	fm := &FTPManager{
		shell: shell,
	}

	// Detect which FTP server is installed
	if shell.CommandExists("vsftpd") {
		fm.backend = "vsftpd"
		fm.enabled = true
	} else if shell.CommandExists("proftpd") {
		fm.backend = "proftpd"
		fm.enabled = true
	} else if shell.CommandExists("pure-ftpd") {
		fm.backend = "pure-ftpd"
		fm.enabled = true
	} else {
		return nil, fmt.Errorf("no FTP server installed (vsftpd, proftpd, or pure-ftpd)")
	}

	return fm, nil
}

// IsEnabled returns whether FTP is available
func (f *FTPManager) IsEnabled() bool {
	return f.enabled
}

// GetBackend returns the FTP server backend being used
func (f *FTPManager) GetBackend() string {
	return f.backend
}

// GetStatus gets FTP service status
func (f *FTPManager) GetStatus() (bool, error) {
	serviceName := f.backend
	if serviceName == "vsftpd" {
		serviceName = "vsftpd"
	}

	result, err := f.shell.Execute("systemctl", "is-active", serviceName)
	if err != nil {
		return false, nil
	}

	return strings.TrimSpace(result.Stdout) == "active", nil
}

// Start starts the FTP service
func (f *FTPManager) Start() error {
	_, err := f.shell.Execute("systemctl", "start", f.backend)
	if err != nil {
		return fmt.Errorf("failed to start FTP: %w", err)
	}

	return nil
}

// Stop stops the FTP service
func (f *FTPManager) Stop() error {
	_, err := f.shell.Execute("systemctl", "stop", f.backend)
	if err != nil {
		return fmt.Errorf("failed to stop FTP: %w", err)
	}

	return nil
}

// Restart restarts the FTP service
func (f *FTPManager) Restart() error {
	_, err := f.shell.Execute("systemctl", "restart", f.backend)
	if err != nil {
		return fmt.Errorf("failed to restart FTP: %w", err)
	}

	return nil
}

// GetConfig gets FTP server configuration
func (f *FTPManager) GetConfig() (*FTPConfig, error) {
	// This would need to parse the FTP server's config file
	// Different for each backend (vsftpd.conf, proftpd.conf, etc.)
	return &FTPConfig{
		Port:            21,
		PasvMinPort:     30000,
		PasvMaxPort:     31000,
		AnonymousEnable: false,
		TLSEnable:       false,
		ChrootEnable:    true,
	}, nil
}

// SetConfig sets FTP server configuration
func (f *FTPManager) SetConfig(config FTPConfig) error {
	// Would need to modify config file based on backend
	return fmt.Errorf("FTP configuration not yet implemented - please configure manually in /etc/%s.conf", f.backend)
}

// EnableTLS enables TLS/SSL for FTP (FTPS)
func (f *FTPManager) EnableTLS(certPath string, keyPath string) error {
	// Would need to configure TLS in the FTP server config
	return fmt.Errorf("FTP TLS configuration not yet implemented")
}

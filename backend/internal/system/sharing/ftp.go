// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package sharing

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"fmt"
	"strings"
)

// FTPManager manages FTP/SFTP server
type FTPManager struct {
	shell      executor.ShellExecutor
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
func NewFTPManager(shell executor.ShellExecutor) (*FTPManager, error) {
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

// ===== Advanced FTP Features =====

// UpdateVsftpdConfig updates vsftpd configuration file
func (f *FTPManager) UpdateVsftpdConfig(key string, value string) error {
	if !f.enabled || f.backend != "vsftpd" {
		return fmt.Errorf("vsftpd not available")
	}

	configFile := "/etc/vsftpd.conf"

	// Check if key exists in config
	result, err := f.shell.Execute("grep", "-q", fmt.Sprintf("^%s=", key), configFile)
	if err == nil && result.ExitCode == 0 {
		// Key exists, update it
		_, err = f.shell.Execute("sed", "-i", fmt.Sprintf("s/^%s=.*/%s=%s/", key, key, value), configFile)
	} else {
		// Key doesn't exist, append it
		_, err = f.shell.Execute("sh", "-c", fmt.Sprintf("echo '%s=%s' >> %s", key, value, configFile))
	}

	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}

// EnableAnonymousFTP enables anonymous FTP access
func (f *FTPManager) EnableAnonymousFTP(enable bool) error {
	if !f.enabled || f.backend != "vsftpd" {
		return fmt.Errorf("vsftpd not available")
	}

	value := "NO"
	if enable {
		value = "YES"
	}

	if err := f.UpdateVsftpdConfig("anonymous_enable", value); err != nil {
		return err
	}

	return f.Restart()
}

// EnableLocalUsers enables local user FTP access
func (f *FTPManager) EnableLocalUsers(enable bool) error {
	if !f.enabled || f.backend != "vsftpd" {
		return fmt.Errorf("vsftpd not available")
	}

	value := "NO"
	if enable {
		value = "YES"
	}

	if err := f.UpdateVsftpdConfig("local_enable", value); err != nil {
		return err
	}

	return f.Restart()
}

// EnableWriteAccess enables write access for FTP users
func (f *FTPManager) EnableWriteAccess(enable bool) error {
	if !f.enabled || f.backend != "vsftpd" {
		return fmt.Errorf("vsftpd not available")
	}

	value := "NO"
	if enable {
		value = "YES"
	}

	if err := f.UpdateVsftpdConfig("write_enable", value); err != nil {
		return err
	}

	return f.Restart()
}

// SetPasvPorts sets the passive mode port range
func (f *FTPManager) SetPasvPorts(minPort int, maxPort int) error {
	if !f.enabled || f.backend != "vsftpd" {
		return fmt.Errorf("vsftpd not available")
	}

	if err := f.UpdateVsftpdConfig("pasv_min_port", fmt.Sprintf("%d", minPort)); err != nil {
		return err
	}

	if err := f.UpdateVsftpdConfig("pasv_max_port", fmt.Sprintf("%d", maxPort)); err != nil {
		return err
	}

	return f.Restart()
}

// EnableChrootJail enables chroot jail for FTP users
func (f *FTPManager) EnableChrootJail(enable bool) error {
	if !f.enabled || f.backend != "vsftpd" {
		return fmt.Errorf("vsftpd not available")
	}

	value := "NO"
	if enable {
		value = "YES"
	}

	if err := f.UpdateVsftpdConfig("chroot_local_user", value); err != nil {
		return err
	}

	// Optionally allow write in chroot
	if err := f.UpdateVsftpdConfig("allow_writeable_chroot", "YES"); err != nil {
		return err
	}

	return f.Restart()
}

// ConfigureVsftpdTLS configures TLS/SSL for vsftpd
func (f *FTPManager) ConfigureVsftpdTLS(certFile string, keyFile string) error {
	if !f.enabled || f.backend != "vsftpd" {
		return fmt.Errorf("vsftpd not available")
	}

	// Enable SSL
	if err := f.UpdateVsftpdConfig("ssl_enable", "YES"); err != nil {
		return err
	}

	// Set certificate paths
	if err := f.UpdateVsftpdConfig("rsa_cert_file", certFile); err != nil {
		return err
	}

	if err := f.UpdateVsftpdConfig("rsa_private_key_file", keyFile); err != nil {
		return err
	}

	// Force SSL for data and login
	if err := f.UpdateVsftpdConfig("force_local_data_ssl", "YES"); err != nil {
		return err
	}

	if err := f.UpdateVsftpdConfig("force_local_logins_ssl", "YES"); err != nil {
		return err
	}

	// SSL protocol options
	if err := f.UpdateVsftpdConfig("ssl_tlsv1", "YES"); err != nil {
		return err
	}

	if err := f.UpdateVsftpdConfig("ssl_sslv2", "NO"); err != nil {
		return err
	}

	if err := f.UpdateVsftpdConfig("ssl_sslv3", "NO"); err != nil {
		return err
	}

	return f.Restart()
}

// AddVirtualUser adds a virtual FTP user
func (f *FTPManager) AddVirtualUser(username string, password string, homeDir string) error {
	if !f.enabled {
		return fmt.Errorf("FTP not available")
	}

	// Create home directory
	_, err := f.shell.Execute("mkdir", "-p", homeDir)
	if err != nil {
		return fmt.Errorf("failed to create home directory: %w", err)
	}

	// Create user (this is simplified - real implementation would use PAM/db)
	// For vsftpd virtual users, you'd typically use a database file
	return fmt.Errorf("virtual user creation requires PAM configuration")
}

// SetBandwidthLimit sets upload/download bandwidth limits
func (f *FTPManager) SetBandwidthLimit(downloadRate int, uploadRate int) error {
	if !f.enabled || f.backend != "vsftpd" {
		return fmt.Errorf("vsftpd not available")
	}

	// Set local (authenticated user) rates in bytes/sec
	if err := f.UpdateVsftpdConfig("local_max_rate", fmt.Sprintf("%d", downloadRate)); err != nil {
		return err
	}

	// Set anonymous rates if needed
	if err := f.UpdateVsftpdConfig("anon_max_rate", fmt.Sprintf("%d", downloadRate)); err != nil {
		return err
	}

	return f.Restart()
}

// SetMaxClients sets the maximum number of concurrent clients
func (f *FTPManager) SetMaxClients(max int) error {
	if !f.enabled || f.backend != "vsftpd" {
		return fmt.Errorf("vsftpd not available")
	}

	if err := f.UpdateVsftpdConfig("max_clients", fmt.Sprintf("%d", max)); err != nil {
		return err
	}

	return f.Restart()
}

// SetMaxPerIP sets the maximum connections per IP address
func (f *FTPManager) SetMaxPerIP(max int) error {
	if !f.enabled || f.backend != "vsftpd" {
		return fmt.Errorf("vsftpd not available")
	}

	if err := f.UpdateVsftpdConfig("max_per_ip", fmt.Sprintf("%d", max)); err != nil {
		return err
	}

	return f.Restart()
}

// EnableLogging enables FTP logging
func (f *FTPManager) EnableLogging(logFile string) error {
	if !f.enabled || f.backend != "vsftpd" {
		return fmt.Errorf("vsftpd not available")
	}

	if err := f.UpdateVsftpdConfig("xferlog_enable", "YES"); err != nil {
		return err
	}

	if err := f.UpdateVsftpdConfig("xferlog_file", logFile); err != nil {
		return err
	}

	// Enable vsftpd style logging
	if err := f.UpdateVsftpdConfig("log_ftp_protocol", "YES"); err != nil {
		return err
	}

	return f.Restart()
}

// GetLogs retrieves FTP logs
func (f *FTPManager) GetLogs(lines int) (string, error) {
	if !f.enabled {
		return "", fmt.Errorf("FTP not available")
	}

	logFile := "/var/log/vsftpd.log"
	if f.backend == "proftpd" {
		logFile = "/var/log/proftpd/proftpd.log"
	}

	result, err := f.shell.Execute("tail", "-n", fmt.Sprintf("%d", lines), logFile)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return result.Stdout, nil
}

// GetActiveConnections returns active FTP connections
func (f *FTPManager) GetActiveConnections() (string, error) {
	if !f.enabled {
		return "", fmt.Errorf("FTP not available")
	}

	// Use netstat to show active FTP connections
	result, err := f.shell.Execute("netstat", "-tn", "|", "grep", ":21")
	if err != nil {
		return "", fmt.Errorf("failed to get connections: %w", err)
	}

	return result.Stdout, nil
}

// BanIP bans an IP address from FTP access
func (f *FTPManager) BanIP(ipAddress string) error {
	if !f.enabled {
		return fmt.Errorf("FTP not available")
	}

	// Use iptables to ban IP
	_, err := f.shell.Execute("iptables", "-A", "INPUT", "-s", ipAddress, "-p", "tcp", "--dport", "21", "-j", "DROP")
	if err != nil {
		return fmt.Errorf("failed to ban IP: %w", err)
	}

	return nil
}

// UnbanIP unbans an IP address
func (f *FTPManager) UnbanIP(ipAddress string) error {
	if !f.enabled {
		return fmt.Errorf("FTP not available")
	}

	// Use iptables to unban IP
	_, err := f.shell.Execute("iptables", "-D", "INPUT", "-s", ipAddress, "-p", "tcp", "--dport", "21", "-j", "DROP")
	if err != nil {
		return fmt.Errorf("failed to unban IP: %w", err)
	}

	return nil
}

// TestConfiguration tests FTP configuration
func (f *FTPManager) TestConfiguration() error {
	if !f.enabled {
		return fmt.Errorf("FTP not available")
	}

	// Try to start in test mode if supported
	// vsftpd doesn't have a built-in config test, so we just verify the file exists
	configFile := fmt.Sprintf("/etc/%s.conf", f.backend)
	_, err := f.shell.Execute("test", "-f", configFile)
	if err != nil {
		return fmt.Errorf("configuration file not found: %s", configFile)
	}

	return nil
}

// BackupConfiguration backs up the FTP configuration
func (f *FTPManager) BackupConfiguration(backupPath string) error {
	if !f.enabled {
		return fmt.Errorf("FTP not available")
	}

	configFile := fmt.Sprintf("/etc/%s.conf", f.backend)
	_, err := f.shell.Execute("cp", configFile, backupPath)
	if err != nil {
		return fmt.Errorf("failed to backup configuration: %w", err)
	}

	return nil
}

// RestoreConfiguration restores the FTP configuration from a backup
func (f *FTPManager) RestoreConfiguration(backupPath string) error {
	if !f.enabled {
		return fmt.Errorf("FTP not available")
	}

	configFile := fmt.Sprintf("/etc/%s.conf", f.backend)
	_, err := f.shell.Execute("cp", backupPath, configFile)
	if err != nil {
		return fmt.Errorf("failed to restore configuration: %w", err)
	}

	return f.Restart()
}

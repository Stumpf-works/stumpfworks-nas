// Revision: 2025-11-23 | Author: Claude | Version: 1.0.0
package ad

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system"
	"github.com/rs/zerolog/log"
)

// DCService manages Active Directory Domain Controller functionality
type DCService struct {
	sambaTool    *SambaTool
	config       *DCConfig
	mu           sync.RWMutex
	provisioned  bool
	domainInfo   map[string]interface{}
}

// DCConfig contains AD DC configuration
type DCConfig struct {
	Enabled        bool   `json:"enabled"`
	Realm          string `json:"realm"`           // e.g., EXAMPLE.COM
	Domain         string `json:"domain"`          // NetBIOS name, e.g., EXAMPLE
	ServerRole     string `json:"server_role"`     // dc, member, standalone
	DNSBackend     string `json:"dns_backend"`     // SAMBA_INTERNAL, BIND9_DLZ, NONE
	DNSForwarder   string `json:"dns_forwarder"`   // Forwarder IP
	FunctionLevel  string `json:"function_level"`  // 2008_R2, 2012, 2012_R2, 2016
	HostIP         string `json:"host_ip"`         // Server IP
	SysvolPath     string `json:"sysvol_path"`     // Path to SYSVOL
	PrivateDirPath string `json:"private_dir_path"` // Path to private dir
}

var (
	globalDCService *DCService
	dcOnce          sync.Once
)

// DefaultDCConfig returns default AD DC configuration
func DefaultDCConfig() *DCConfig {
	return &DCConfig{
		Enabled:        false,
		ServerRole:     "dc",
		DNSBackend:     "SAMBA_INTERNAL",
		FunctionLevel:  "2008_R2",
		SysvolPath:     "/var/lib/samba/sysvol",
		PrivateDirPath: "/var/lib/samba/private",
	}
}

// InitializeDC initializes the AD DC service
func InitializeDC() (*DCService, error) {
	var err error
	dcOnce.Do(func() {
		shellExecutor, shellErr := system.NewShellExecutor(30*time.Second, false)
		if shellErr != nil {
			err = fmt.Errorf("failed to create shell executor: %w", shellErr)
			return
		}
		sambaTool := NewSambaTool(shellExecutor)

		if !sambaTool.IsAvailable() {
			log.Warn().Msg("samba-tool not available, AD DC features disabled")
			err = fmt.Errorf("samba-tool not available")
			return
		}

		globalDCService = &DCService{
			sambaTool:   sambaTool,
			config:      DefaultDCConfig(),
			provisioned: false,
		}

		// Check if already provisioned
		if _, statErr := os.Stat("/var/lib/samba/private/sam.ldb"); statErr == nil {
			globalDCService.provisioned = true
			log.Info().Msg("AD DC already provisioned")

			// Load domain info
			if info, infoErr := sambaTool.GetDomainInfo(); infoErr == nil {
				globalDCService.domainInfo = info
				if realm, ok := info["Forest"].(string); ok {
					globalDCService.config.Realm = realm
				}
				if domain, ok := info["Domain"].(string); ok {
					globalDCService.config.Domain = domain
				}
			}
		}

		log.Info().Msg("AD DC service initialized")
	})

	return globalDCService, err
}

// GetDCService returns the global DC service
func GetDCService() *DCService {
	return globalDCService
}

// ===== Domain Controller Management =====

// IsProvisioned returns whether the DC is provisioned
func (dc *DCService) IsProvisioned() bool {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return dc.provisioned
}

// GetConfig returns the DC configuration
func (dc *DCService) GetConfig() *DCConfig {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return dc.config
}

// UpdateConfig updates the DC configuration
func (dc *DCService) UpdateConfig(config *DCConfig) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	dc.config = config
	return nil
}

// Provision provisions a new AD domain
func (dc *DCService) Provision(opts ProvisionOptions) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if dc.provisioned {
		return fmt.Errorf("domain controller already provisioned")
	}

	log.Info().Str("realm", opts.Realm).Str("domain", opts.Domain).Msg("Provisioning AD domain")

	// Stop Samba if running
	if err := dc.stopSambaService(); err != nil {
		log.Warn().Err(err).Msg("Failed to stop Samba service")
	}

	// Backup existing configuration
	if err := dc.backupConfiguration(); err != nil {
		log.Warn().Err(err).Msg("Failed to backup configuration")
	}

	// Provision domain
	if err := dc.sambaTool.ProvisionDomain(opts); err != nil {
		return fmt.Errorf("failed to provision domain: %w", err)
	}

	// Update configuration
	dc.config.Realm = opts.Realm
	dc.config.Domain = opts.Domain
	dc.config.ServerRole = opts.ServerRole
	dc.config.DNSBackend = opts.DNSBackend
	dc.config.DNSForwarder = opts.DNSForwarder
	dc.config.FunctionLevel = opts.FunctionLevel
	dc.config.HostIP = opts.HostIP
	dc.config.Enabled = true
	dc.provisioned = true

	// Start Samba service
	if err := dc.startSambaService(); err != nil {
		return fmt.Errorf("failed to start Samba service: %w", err)
	}

	// Load domain info
	if info, err := dc.sambaTool.GetDomainInfo(); err == nil {
		dc.domainInfo = info
	}

	log.Info().Msg("AD domain provisioned successfully")
	return nil
}

// Demote demotes the domain controller
func (dc *DCService) Demote() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	log.Info().Msg("Demoting AD domain controller")

	// Stop Samba service
	if err := dc.stopSambaService(); err != nil {
		log.Warn().Err(err).Msg("Failed to stop Samba service")
	}

	// Demote domain
	if err := dc.sambaTool.DemoteDomain(); err != nil {
		return fmt.Errorf("failed to demote domain: %w", err)
	}

	dc.provisioned = false
	dc.config.Enabled = false
	dc.domainInfo = nil

	log.Info().Msg("AD domain controller demoted successfully")
	return nil
}

// GetDomainInfo returns domain information
func (dc *DCService) GetDomainInfo() (map[string]interface{}, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	if !dc.provisioned {
		return nil, fmt.Errorf("domain controller not provisioned")
	}

	return dc.domainInfo, nil
}

// GetDomainLevel returns the domain functional level
func (dc *DCService) GetDomainLevel() (string, error) {
	if !dc.provisioned {
		return "", fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.GetDomainLevel()
}

// RaiseDomainLevel raises the domain functional level
func (dc *DCService) RaiseDomainLevel(level string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.RaiseDomainLevel(level)
}

// ===== User Management =====

// CreateUser creates a new AD user
func (dc *DCService) CreateUser(user ADDCUser, password string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.CreateUser(user, password)
}

// DeleteUser deletes an AD user
func (dc *DCService) DeleteUser(username string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.DeleteUser(username)
}

// EnableUser enables an AD user
func (dc *DCService) EnableUser(username string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.EnableUser(username)
}

// DisableUser disables an AD user
func (dc *DCService) DisableUser(username string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.DisableUser(username)
}

// SetUserPassword sets a user's password
func (dc *DCService) SetUserPassword(username, password string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.SetUserPassword(username, password)
}

// SetUserExpiry sets user account expiry
func (dc *DCService) SetUserExpiry(username string, days int) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.SetUserExpiry(username, days)
}

// ListUsers lists all AD users
func (dc *DCService) ListUsers() ([]string, error) {
	if !dc.provisioned {
		return nil, fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.ListUsers()
}

// ===== Group Management =====

// CreateGroup creates a new AD group
func (dc *DCService) CreateGroup(group ADGroup) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.CreateGroup(group)
}

// DeleteGroup deletes an AD group
func (dc *DCService) DeleteGroup(groupName string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.DeleteGroup(groupName)
}

// AddGroupMember adds a user to a group
func (dc *DCService) AddGroupMember(groupName, username string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.AddGroupMember(groupName, username)
}

// RemoveGroupMember removes a user from a group
func (dc *DCService) RemoveGroupMember(groupName, username string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.RemoveGroupMember(groupName, username)
}

// ListGroups lists all AD groups
func (dc *DCService) ListGroups() ([]string, error) {
	if !dc.provisioned {
		return nil, fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.ListGroups()
}

// ListGroupMembers lists members of a group
func (dc *DCService) ListGroupMembers(groupName string) ([]string, error) {
	if !dc.provisioned {
		return nil, fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.ListGroupMembers(groupName)
}

// ===== Computer Management =====

// CreateComputer creates a new AD computer
func (dc *DCService) CreateComputer(computer ADComputer) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.CreateComputer(computer)
}

// DeleteComputer deletes an AD computer
func (dc *DCService) DeleteComputer(computerName string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.DeleteComputer(computerName)
}

// ListComputers lists all AD computers
func (dc *DCService) ListComputers() ([]string, error) {
	if !dc.provisioned {
		return nil, fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.ListComputers()
}

// ===== Organizational Unit Management =====

// CreateOU creates a new OU
func (dc *DCService) CreateOU(ou ADOU) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.CreateOU(ou)
}

// DeleteOU deletes an OU
func (dc *DCService) DeleteOU(ouDN string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.DeleteOU(ouDN)
}

// ListOUs lists all OUs
func (dc *DCService) ListOUs() ([]string, error) {
	if !dc.provisioned {
		return nil, fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.ListOUs()
}

// ===== Group Policy Management =====

// CreateGPO creates a new GPO
func (dc *DCService) CreateGPO(gpo ADGPO) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.CreateGPO(gpo)
}

// DeleteGPO deletes a GPO
func (dc *DCService) DeleteGPO(gpoName string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.DeleteGPO(gpoName)
}

// ListGPOs lists all GPOs
func (dc *DCService) ListGPOs() ([]string, error) {
	if !dc.provisioned {
		return nil, fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.ListGPOs()
}

// LinkGPO links a GPO to an OU
func (dc *DCService) LinkGPO(gpoName, ouDN string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.LinkGPO(gpoName, ouDN)
}

// UnlinkGPO unlinks a GPO from an OU
func (dc *DCService) UnlinkGPO(gpoName, ouDN string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.UnlinkGPO(gpoName, ouDN)
}

// ===== DNS Management =====

// AddDNSRecord adds a DNS record
func (dc *DCService) AddDNSRecord(zone string, record ADDNSRecord) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.AddDNSRecord(zone, record)
}

// DeleteDNSRecord deletes a DNS record
func (dc *DCService) DeleteDNSRecord(zone, name, recordType, value string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.DeleteDNSRecord(zone, name, recordType, value)
}

// ListDNSRecords lists DNS records in a zone
func (dc *DCService) ListDNSRecords(zone string) ([]string, error) {
	if !dc.provisioned {
		return nil, fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.ListDNSRecords(zone)
}

// CreateDNSZone creates a new DNS zone
func (dc *DCService) CreateDNSZone(zoneName string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.CreateDNSZone(zoneName)
}

// DeleteDNSZone deletes a DNS zone
func (dc *DCService) DeleteDNSZone(zoneName string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.DeleteDNSZone(zoneName)
}

// ListDNSZones lists all DNS zones
func (dc *DCService) ListDNSZones() ([]string, error) {
	if !dc.provisioned {
		return nil, fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.ListDNSZones()
}

// ===== FSMO Roles =====

// TransferFSMORoles transfers FSMO roles
func (dc *DCService) TransferFSMORoles(role, targetDC string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.TransferFSMORoles(role, targetDC)
}

// SeizeFSMORoles seizes FSMO roles
func (dc *DCService) SeizeFSMORoles(role string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.SeizeFSMORoles(role)
}

// ShowFSMORoles shows FSMO roles
func (dc *DCService) ShowFSMORoles() (map[string]string, error) {
	if !dc.provisioned {
		return nil, fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.ShowFSMORoles()
}

// ===== Utility Functions =====

// TestConfiguration tests the configuration
func (dc *DCService) TestConfiguration() error {
	return dc.sambaTool.TestConfiguration()
}

// ShowDBCheck runs database check
func (dc *DCService) ShowDBCheck() (string, error) {
	if !dc.provisioned {
		return "", fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.ShowDBCheck()
}

// BackupOnline performs an online backup
func (dc *DCService) BackupOnline(targetDir string) error {
	if !dc.provisioned {
		return fmt.Errorf("domain controller not provisioned")
	}

	return dc.sambaTool.BackupOnline(targetDir)
}

// ===== Private Helper Methods =====

func (dc *DCService) startSambaService() error {
	shell, _ := system.NewShellExecutor(30*time.Second, false)
	result, err := shell.Execute("systemctl", "start", "samba-ad-dc")
	if err != nil {
		return fmt.Errorf("failed to start samba-ad-dc: %s: %w", result.Stderr, err)
	}
	return nil
}

func (dc *DCService) stopSambaService() error {
	shell, _ := system.NewShellExecutor(30*time.Second, false)
	result, err := shell.Execute("systemctl", "stop", "samba-ad-dc")
	if err != nil {
		return fmt.Errorf("failed to stop samba-ad-dc: %s: %w", result.Stderr, err)
	}
	return nil
}

func (dc *DCService) restartSambaService() error {
	shell, _ := system.NewShellExecutor(30*time.Second, false)
	result, err := shell.Execute("systemctl", "restart", "samba-ad-dc")
	if err != nil {
		return fmt.Errorf("failed to restart samba-ad-dc: %s: %w", result.Stderr, err)
	}
	return nil
}

func (dc *DCService) getSambaServiceStatus() (string, error) {
	shell, _ := system.NewShellExecutor(30*time.Second, false)
	result, err := shell.Execute("systemctl", "is-active", "samba-ad-dc")
	if err != nil {
		return "inactive", nil
	}
	return result.Stdout, nil
}

func (dc *DCService) backupConfiguration() error {
	// Backup smb.conf
	shell, _ := system.NewShellExecutor(30*time.Second, false)
	_, _ = shell.Execute("cp", "/etc/samba/smb.conf", "/etc/samba/smb.conf.bak")

	// Backup other configs if they exist
	_, _ = shell.Execute("cp", "/etc/krb5.conf", "/etc/krb5.conf.bak")

	return nil
}

// GetSambaServiceStatus returns the Samba AD DC service status
func (dc *DCService) GetSambaServiceStatus() (string, error) {
	return dc.getSambaServiceStatus()
}

// RestartSambaService restarts the Samba AD DC service
func (dc *DCService) RestartSambaService() error {
	return dc.restartSambaService()
}

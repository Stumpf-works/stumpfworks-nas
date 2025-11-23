// Revision: 2025-11-23 | Author: Claude | Version: 1.0.0
package ad

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"github.com/rs/zerolog/log"
)

// SambaTool provides a wrapper around samba-tool commands
type SambaTool struct {
	shell executor.ShellExecutor
}

// NewSambaTool creates a new samba-tool wrapper
func NewSambaTool(shell executor.ShellExecutor) *SambaTool {
	return &SambaTool{
		shell: shell,
	}
}

// IsAvailable checks if samba-tool is available
func (st *SambaTool) IsAvailable() bool {
	return st.shell.CommandExists("samba-tool")
}

// GetVersion returns the Samba version
func (st *SambaTool) GetVersion() (string, error) {
	result, err := st.shell.Execute("samba-tool", "--version")
	if err != nil {
		return "", fmt.Errorf("failed to get samba version: %w", err)
	}
	return strings.TrimSpace(result.Stdout), nil
}

// ===== Domain Provisioning =====

// ProvisionOptions contains options for provisioning a new AD domain
type ProvisionOptions struct {
	Realm           string // e.g., EXAMPLE.COM
	Domain          string // NetBIOS domain name, e.g., EXAMPLE
	AdminPassword   string // Administrator password
	DNSBackend      string // SAMBA_INTERNAL, BIND9_DLZ, or NONE
	DNSForwarder    string // Optional DNS forwarder IP
	ServerRole      string // dc, member, standalone
	UseTLS          bool   // Use LDAPS
	FunctionLevel   string // 2008_R2, 2012, 2012_R2, 2016
	HostIP          string // Server IP address
}

// ProvisionDomain provisions a new AD domain
func (st *SambaTool) ProvisionDomain(opts ProvisionOptions) error {
	args := []string{
		"domain", "provision",
		"--realm=" + opts.Realm,
		"--domain=" + opts.Domain,
		"--adminpass=" + opts.AdminPassword,
		"--server-role=" + opts.ServerRole,
	}

	if opts.DNSBackend != "" {
		args = append(args, "--dns-backend="+opts.DNSBackend)
	}

	if opts.DNSForwarder != "" {
		args = append(args, "--option=dns forwarder="+opts.DNSForwarder)
	}

	if opts.FunctionLevel != "" {
		args = append(args, "--function-level="+opts.FunctionLevel)
	}

	if opts.HostIP != "" {
		args = append(args, "--host-ip="+opts.HostIP)
	}

	if opts.UseTLS {
		args = append(args, "--use-rfc2307")
	}

	log.Info().Str("realm", opts.Realm).Str("domain", opts.Domain).Msg("Provisioning AD domain")

	result, err := st.shell.Execute("samba-tool", args...)
	if err != nil {
		return fmt.Errorf("failed to provision domain: %s: %w", result.Stderr, err)
	}

	log.Info().Msg("AD domain provisioned successfully")
	return nil
}

// DemoteDomain demotes the domain controller
func (st *SambaTool) DemoteDomain() error {
	result, err := st.shell.Execute("samba-tool", "domain", "demote")
	if err != nil {
		return fmt.Errorf("failed to demote domain: %s: %w", result.Stderr, err)
	}
	return nil
}

// GetDomainLevel returns the current domain functional level
func (st *SambaTool) GetDomainLevel() (string, error) {
	result, err := st.shell.Execute("samba-tool", "domain", "level", "show")
	if err != nil {
		return "", fmt.Errorf("failed to get domain level: %w", err)
	}
	return strings.TrimSpace(result.Stdout), nil
}

// RaiseDomainLevel raises the domain functional level
func (st *SambaTool) RaiseDomainLevel(level string) error {
	result, err := st.shell.Execute("samba-tool", "domain", "level", "raise", "--domain-level="+level)
	if err != nil {
		return fmt.Errorf("failed to raise domain level: %s: %w", result.Stderr, err)
	}
	return nil
}

// ===== User Management =====

// ADDCUser represents an Active Directory Domain Controller user (extended fields)
type ADDCUser struct {
	Username        string   `json:"username"`
	GivenName       string   `json:"given_name,omitempty"`
	Surname         string   `json:"surname,omitempty"`
	DisplayName     string   `json:"display_name,omitempty"`
	Email           string   `json:"email,omitempty"`
	Description     string   `json:"description,omitempty"`
	Department      string   `json:"department,omitempty"`
	Company         string   `json:"company,omitempty"`
	Title           string   `json:"title,omitempty"`
	Telephone       string   `json:"telephone,omitempty"`
	OU              string   `json:"ou,omitempty"`
	Enabled         bool     `json:"enabled"`
	PasswordExpired bool     `json:"password_expired"`
	MemberOf        []string `json:"member_of,omitempty"`
}

// CreateUser creates a new AD user
func (st *SambaTool) CreateUser(user ADDCUser, password string) error {
	args := []string{"user", "create", user.Username, password}

	if user.GivenName != "" {
		args = append(args, "--given-name="+user.GivenName)
	}
	if user.Surname != "" {
		args = append(args, "--surname="+user.Surname)
	}
	if user.Email != "" {
		args = append(args, "--mail-address="+user.Email)
	}
	if user.Description != "" {
		args = append(args, "--description="+user.Description)
	}
	if user.Department != "" {
		args = append(args, "--department="+user.Department)
	}
	if user.Company != "" {
		args = append(args, "--company="+user.Company)
	}
	if user.Title != "" {
		args = append(args, "--job-title="+user.Title)
	}
	if user.Telephone != "" {
		args = append(args, "--telephone-number="+user.Telephone)
	}
	if user.OU != "" {
		args = append(args, "--userou="+user.OU)
	}

	result, err := st.shell.Execute("samba-tool", args...)
	if err != nil {
		return fmt.Errorf("failed to create user: %s: %w", result.Stderr, err)
	}

	log.Info().Str("username", user.Username).Msg("AD user created")
	return nil
}

// DeleteUser deletes an AD user
func (st *SambaTool) DeleteUser(username string) error {
	result, err := st.shell.Execute("samba-tool", "user", "delete", username)
	if err != nil {
		return fmt.Errorf("failed to delete user: %s: %w", result.Stderr, err)
	}
	log.Info().Str("username", username).Msg("AD user deleted")
	return nil
}

// EnableUser enables an AD user
func (st *SambaTool) EnableUser(username string) error {
	result, err := st.shell.Execute("samba-tool", "user", "enable", username)
	if err != nil {
		return fmt.Errorf("failed to enable user: %s: %w", result.Stderr, err)
	}
	return nil
}

// DisableUser disables an AD user
func (st *SambaTool) DisableUser(username string) error {
	result, err := st.shell.Execute("samba-tool", "user", "disable", username)
	if err != nil {
		return fmt.Errorf("failed to disable user: %s: %w", result.Stderr, err)
	}
	return nil
}

// SetUserPassword sets a user's password
func (st *SambaTool) SetUserPassword(username, password string) error {
	result, err := st.shell.Execute("samba-tool", "user", "setpassword", username, "--newpassword="+password)
	if err != nil {
		return fmt.Errorf("failed to set password: %s: %w", result.Stderr, err)
	}
	return nil
}

// SetUserExpiry sets user account expiry
func (st *SambaTool) SetUserExpiry(username string, days int) error {
	result, err := st.shell.Execute("samba-tool", "user", "setexpiry", username, fmt.Sprintf("--days=%d", days))
	if err != nil {
		return fmt.Errorf("failed to set expiry: %s: %w", result.Stderr, err)
	}
	return nil
}

// ListUsers lists all AD users
func (st *SambaTool) ListUsers() ([]string, error) {
	result, err := st.shell.Execute("samba-tool", "user", "list")
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	var filteredUsers []string
	for _, user := range users {
		user = strings.TrimSpace(user)
		if user != "" {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers, nil
}

// ===== Group Management =====

// ADGroup represents an Active Directory group
type ADGroup struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	OU          string   `json:"ou,omitempty"`
	GroupScope  string   `json:"group_scope,omitempty"` // Domain, Global, Universal
	GroupType   string   `json:"group_type,omitempty"`  // Security, Distribution
	Members     []string `json:"members,omitempty"`
}

// CreateGroup creates a new AD group
func (st *SambaTool) CreateGroup(group ADGroup) error {
	args := []string{"group", "add", group.Name}

	if group.Description != "" {
		args = append(args, "--description="+group.Description)
	}
	if group.OU != "" {
		args = append(args, "--groupou="+group.OU)
	}
	if group.GroupScope != "" {
		args = append(args, "--group-scope="+group.GroupScope)
	}
	if group.GroupType != "" {
		args = append(args, "--group-type="+group.GroupType)
	}

	result, err := st.shell.Execute("samba-tool", args...)
	if err != nil {
		return fmt.Errorf("failed to create group: %s: %w", result.Stderr, err)
	}

	log.Info().Str("group", group.Name).Msg("AD group created")
	return nil
}

// DeleteGroup deletes an AD group
func (st *SambaTool) DeleteGroup(groupName string) error {
	result, err := st.shell.Execute("samba-tool", "group", "delete", groupName)
	if err != nil {
		return fmt.Errorf("failed to delete group: %s: %w", result.Stderr, err)
	}
	log.Info().Str("group", groupName).Msg("AD group deleted")
	return nil
}

// AddGroupMember adds a user to a group
func (st *SambaTool) AddGroupMember(groupName, username string) error {
	result, err := st.shell.Execute("samba-tool", "group", "addmembers", groupName, username)
	if err != nil {
		return fmt.Errorf("failed to add group member: %s: %w", result.Stderr, err)
	}
	return nil
}

// RemoveGroupMember removes a user from a group
func (st *SambaTool) RemoveGroupMember(groupName, username string) error {
	result, err := st.shell.Execute("samba-tool", "group", "removemembers", groupName, username)
	if err != nil {
		return fmt.Errorf("failed to remove group member: %s: %w", result.Stderr, err)
	}
	return nil
}

// ListGroups lists all AD groups
func (st *SambaTool) ListGroups() ([]string, error) {
	result, err := st.shell.Execute("samba-tool", "group", "list")
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	groups := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	var filteredGroups []string
	for _, group := range groups {
		group = strings.TrimSpace(group)
		if group != "" {
			filteredGroups = append(filteredGroups, group)
		}
	}
	return filteredGroups, nil
}

// ListGroupMembers lists members of a group
func (st *SambaTool) ListGroupMembers(groupName string) ([]string, error) {
	result, err := st.shell.Execute("samba-tool", "group", "listmembers", groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to list group members: %w", err)
	}

	members := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	var filteredMembers []string
	for _, member := range members {
		member = strings.TrimSpace(member)
		if member != "" {
			filteredMembers = append(filteredMembers, member)
		}
	}
	return filteredMembers, nil
}

// ===== Computer Management =====

// ADComputer represents an Active Directory computer
type ADComputer struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	OU          string `json:"ou,omitempty"`
	IP          string `json:"ip,omitempty"`
	Enabled     bool   `json:"enabled"`
}

// CreateComputer creates a new AD computer account
func (st *SambaTool) CreateComputer(computer ADComputer) error {
	args := []string{"computer", "create", computer.Name}

	if computer.Description != "" {
		args = append(args, "--description="+computer.Description)
	}
	if computer.OU != "" {
		args = append(args, "--computerou="+computer.OU)
	}
	if computer.IP != "" {
		args = append(args, "--ip-address="+computer.IP)
	}

	result, err := st.shell.Execute("samba-tool", args...)
	if err != nil {
		return fmt.Errorf("failed to create computer: %s: %w", result.Stderr, err)
	}

	log.Info().Str("computer", computer.Name).Msg("AD computer created")
	return nil
}

// DeleteComputer deletes an AD computer account
func (st *SambaTool) DeleteComputer(computerName string) error {
	result, err := st.shell.Execute("samba-tool", "computer", "delete", computerName)
	if err != nil {
		return fmt.Errorf("failed to delete computer: %s: %w", result.Stderr, err)
	}
	log.Info().Str("computer", computerName).Msg("AD computer deleted")
	return nil
}

// ListComputers lists all AD computers
func (st *SambaTool) ListComputers() ([]string, error) {
	result, err := st.shell.Execute("samba-tool", "computer", "list")
	if err != nil {
		return nil, fmt.Errorf("failed to list computers: %w", err)
	}

	computers := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	var filteredComputers []string
	for _, computer := range computers {
		computer = strings.TrimSpace(computer)
		if computer != "" {
			filteredComputers = append(filteredComputers, computer)
		}
	}
	return filteredComputers, nil
}

// ===== Organizational Unit Management =====

// ADOU represents an Active Directory Organizational Unit
type ADOU struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ParentDN    string `json:"parent_dn,omitempty"` // e.g., "DC=example,DC=com"
}

// CreateOU creates a new Organizational Unit
func (st *SambaTool) CreateOU(ou ADOU) error {
	args := []string{"ou", "create", ou.Name}

	if ou.Description != "" {
		args = append(args, "--description="+ou.Description)
	}

	result, err := st.shell.Execute("samba-tool", args...)
	if err != nil {
		return fmt.Errorf("failed to create OU: %s: %w", result.Stderr, err)
	}

	log.Info().Str("ou", ou.Name).Msg("AD OU created")
	return nil
}

// DeleteOU deletes an Organizational Unit
func (st *SambaTool) DeleteOU(ouDN string) error {
	result, err := st.shell.Execute("samba-tool", "ou", "delete", ouDN)
	if err != nil {
		return fmt.Errorf("failed to delete OU: %s: %w", result.Stderr, err)
	}
	log.Info().Str("ou", ouDN).Msg("AD OU deleted")
	return nil
}

// ListOUs lists all Organizational Units
func (st *SambaTool) ListOUs() ([]string, error) {
	result, err := st.shell.Execute("samba-tool", "ou", "list")
	if err != nil {
		return nil, fmt.Errorf("failed to list OUs: %w", err)
	}

	ous := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	var filteredOUs []string
	for _, ou := range ous {
		ou = strings.TrimSpace(ou)
		if ou != "" {
			filteredOUs = append(filteredOUs, ou)
		}
	}
	return filteredOUs, nil
}

// ===== Group Policy Object Management =====

// ADGPO represents a Group Policy Object
type ADGPO struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
}

// CreateGPO creates a new Group Policy Object
func (st *SambaTool) CreateGPO(gpo ADGPO) error {
	args := []string{"gpo", "create", gpo.Name}

	if gpo.DisplayName != "" {
		args = append(args, "--displayname="+gpo.DisplayName)
	}

	result, err := st.shell.Execute("samba-tool", args...)
	if err != nil {
		return fmt.Errorf("failed to create GPO: %s: %w", result.Stderr, err)
	}

	log.Info().Str("gpo", gpo.Name).Msg("GPO created")
	return nil
}

// DeleteGPO deletes a Group Policy Object
func (st *SambaTool) DeleteGPO(gpoName string) error {
	result, err := st.shell.Execute("samba-tool", "gpo", "del", gpoName)
	if err != nil {
		return fmt.Errorf("failed to delete GPO: %s: %w", result.Stderr, err)
	}
	log.Info().Str("gpo", gpoName).Msg("GPO deleted")
	return nil
}

// ListGPOs lists all Group Policy Objects
func (st *SambaTool) ListGPOs() ([]string, error) {
	result, err := st.shell.Execute("samba-tool", "gpo", "listall")
	if err != nil {
		return nil, fmt.Errorf("failed to list GPOs: %w", err)
	}

	gpos := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	var filteredGPOs []string
	for _, gpo := range gpos {
		gpo = strings.TrimSpace(gpo)
		if gpo != "" && !strings.HasPrefix(gpo, "GPO") { // Filter header
			filteredGPOs = append(filteredGPOs, gpo)
		}
	}
	return filteredGPOs, nil
}

// LinkGPO links a GPO to an OU
func (st *SambaTool) LinkGPO(gpoName, ouDN string) error {
	result, err := st.shell.Execute("samba-tool", "gpo", "setlink", ouDN, gpoName)
	if err != nil {
		return fmt.Errorf("failed to link GPO: %s: %w", result.Stderr, err)
	}
	return nil
}

// UnlinkGPO unlinks a GPO from an OU
func (st *SambaTool) UnlinkGPO(gpoName, ouDN string) error {
	result, err := st.shell.Execute("samba-tool", "gpo", "dellink", ouDN, gpoName)
	if err != nil {
		return fmt.Errorf("failed to unlink GPO: %s: %w", result.Stderr, err)
	}
	return nil
}

// ===== DNS Management =====

// ADDNSRecord represents a DNS record in AD
type ADDNSRecord struct {
	Name       string `json:"name"`
	RecordType string `json:"type"` // A, AAAA, CNAME, MX, TXT, SRV
	Value      string `json:"value"`
	TTL        int    `json:"ttl,omitempty"`
}

// AddDNSRecord adds a DNS record
func (st *SambaTool) AddDNSRecord(zone string, record ADDNSRecord) error {
	args := []string{"dns", "add", "localhost", zone, record.Name, record.RecordType, record.Value}

	result, err := st.shell.Execute("samba-tool", args...)
	if err != nil {
		return fmt.Errorf("failed to add DNS record: %s: %w", result.Stderr, err)
	}

	log.Info().Str("zone", zone).Str("record", record.Name).Msg("DNS record added")
	return nil
}

// DeleteDNSRecord deletes a DNS record
func (st *SambaTool) DeleteDNSRecord(zone, name, recordType, value string) error {
	args := []string{"dns", "delete", "localhost", zone, name, recordType, value}

	result, err := st.shell.Execute("samba-tool", args...)
	if err != nil {
		return fmt.Errorf("failed to delete DNS record: %s: %w", result.Stderr, err)
	}

	log.Info().Str("zone", zone).Str("record", name).Msg("DNS record deleted")
	return nil
}

// ListDNSRecords lists DNS records in a zone
func (st *SambaTool) ListDNSRecords(zone string) ([]string, error) {
	result, err := st.shell.Execute("samba-tool", "dns", "query", "localhost", zone, "@", "ALL")
	if err != nil {
		return nil, fmt.Errorf("failed to list DNS records: %w", err)
	}

	records := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	return records, nil
}

// CreateDNSZone creates a new DNS zone
func (st *SambaTool) CreateDNSZone(zoneName string) error {
	result, err := st.shell.Execute("samba-tool", "dns", "zonecreate", "localhost", zoneName)
	if err != nil {
		return fmt.Errorf("failed to create DNS zone: %s: %w", result.Stderr, err)
	}

	log.Info().Str("zone", zoneName).Msg("DNS zone created")
	return nil
}

// DeleteDNSZone deletes a DNS zone
func (st *SambaTool) DeleteDNSZone(zoneName string) error {
	result, err := st.shell.Execute("samba-tool", "dns", "zonedelete", "localhost", zoneName)
	if err != nil {
		return fmt.Errorf("failed to delete DNS zone: %s: %w", result.Stderr, err)
	}

	log.Info().Str("zone", zoneName).Msg("DNS zone deleted")
	return nil
}

// ListDNSZones lists all DNS zones
func (st *SambaTool) ListDNSZones() ([]string, error) {
	result, err := st.shell.Execute("samba-tool", "dns", "zonelist", "localhost")
	if err != nil {
		return nil, fmt.Errorf("failed to list DNS zones: %w", err)
	}

	zones := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	var filteredZones []string
	for _, zone := range zones {
		zone = strings.TrimSpace(zone)
		if zone != "" && !strings.Contains(zone, "pszZoneName") {
			filteredZones = append(filteredZones, zone)
		}
	}
	return filteredZones, nil
}

// ===== FSMO Roles Management =====

// TransferFSMORoles transfers FSMO roles to another DC
func (st *SambaTool) TransferFSMORoles(role, targetDC string) error {
	result, err := st.shell.Execute("samba-tool", "fsmo", "transfer", "--role="+role, "--host="+targetDC)
	if err != nil {
		return fmt.Errorf("failed to transfer FSMO role: %s: %w", result.Stderr, err)
	}
	return nil
}

// SeizeFSMORoles seizes FSMO roles
func (st *SambaTool) SeizeFSMORoles(role string) error {
	result, err := st.shell.Execute("samba-tool", "fsmo", "seize", "--role="+role, "--force")
	if err != nil {
		return fmt.Errorf("failed to seize FSMO role: %s: %w", result.Stderr, err)
	}
	return nil
}

// ShowFSMORoles shows current FSMO role holders
func (st *SambaTool) ShowFSMORoles() (map[string]string, error) {
	result, err := st.shell.Execute("samba-tool", "fsmo", "show")
	if err != nil {
		return nil, fmt.Errorf("failed to show FSMO roles: %w", err)
	}

	roles := make(map[string]string)
	lines := strings.Split(result.Stdout, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				roles[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}
	return roles, nil
}

// ===== Utility Functions =====

// TestConfiguration tests the Samba AD DC configuration
func (st *SambaTool) TestConfiguration() error {
	result, err := st.shell.Execute("samba-tool", "testparm", "--suppress-prompt")
	if err != nil {
		return fmt.Errorf("configuration test failed: %s: %w", result.Stderr, err)
	}
	return nil
}

// ShowDBCheck runs database check
func (st *SambaTool) ShowDBCheck() (string, error) {
	result, err := st.shell.Execute("samba-tool", "dbcheck", "--cross-ncs")
	if err != nil {
		return "", fmt.Errorf("DB check failed: %w", err)
	}
	return result.Stdout, nil
}

// BackupOnline performs an online backup
func (st *SambaTool) BackupOnline(targetDir string) error {
	result, err := st.shell.Execute("samba-tool", "domain", "backup", "online", "--targetdir="+targetDir)
	if err != nil {
		return fmt.Errorf("backup failed: %s: %w", result.Stderr, err)
	}
	return nil
}

// GetNTACL gets NT ACLs for a path
func (st *SambaTool) GetNTACL(path string) (string, error) {
	result, err := st.shell.Execute("samba-tool", "ntacl", "get", path)
	if err != nil {
		return "", fmt.Errorf("failed to get NTACL: %w", err)
	}
	return result.Stdout, nil
}

// SetNTACL sets NT ACLs for a path
func (st *SambaTool) SetNTACL(path, acl string) error {
	result, err := st.shell.Execute("samba-tool", "ntacl", "set", acl, path)
	if err != nil {
		return fmt.Errorf("failed to set NTACL: %s: %w", result.Stderr, err)
	}
	return nil
}

// ExportKeytab exports a keytab file
func (st *SambaTool) ExportKeytab(principal, keytabPath string) error {
	result, err := st.shell.Execute("samba-tool", "domain", "exportkeytab", keytabPath, "--principal="+principal)
	if err != nil {
		return fmt.Errorf("failed to export keytab: %s: %w", result.Stderr, err)
	}
	return nil
}

// GetDomainInfo gets domain information
func (st *SambaTool) GetDomainInfo() (map[string]interface{}, error) {
	result, err := st.shell.Execute("samba-tool", "domain", "info", "localhost")
	if err != nil {
		return nil, fmt.Errorf("failed to get domain info: %w", err)
	}

	info := make(map[string]interface{})
	lines := strings.Split(result.Stdout, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				info[key] = value
			}
		}
	}

	return info, nil
}

// ParseJSON parses JSON output from samba-tool
func (st *SambaTool) ParseJSON(output string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(output), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

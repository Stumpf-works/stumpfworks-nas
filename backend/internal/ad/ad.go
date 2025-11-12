package ad

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"sync"

	"github.com/go-ldap/ldap/v3"
)

// ADConfig holds Active Directory configuration
type ADConfig struct {
	Enabled      bool   `json:"enabled"`
	Server       string `json:"server"`        // AD server address
	Port         int    `json:"port"`          // LDAP port (usually 389 or 636 for LDAPS)
	BaseDN       string `json:"baseDN"`        // Base DN for searches (e.g., "dc=example,dc=com")
	BindUser     string `json:"bindUser"`      // User for binding (e.g., "cn=admin,dc=example,dc=com")
	BindPassword string `json:"bindPassword"`  // Password for bind user
	UserFilter   string `json:"userFilter"`    // LDAP filter for users (e.g., "(&(objectClass=user)(sAMAccountName={username}))")
	GroupFilter  string `json:"groupFilter"`   // LDAP filter for groups
	UseTLS       bool   `json:"useTLS"`        // Use TLS for connection
	SkipVerify   bool   `json:"skipVerify"`    // Skip TLS certificate verification
}

// ADUser represents a user from Active Directory
type ADUser struct {
	Username      string   `json:"username"`
	Email         string   `json:"email"`
	DisplayName   string   `json:"displayName"`
	DistinguishedName string `json:"distinguishedName"`
	Groups        []string `json:"groups"`
	Enabled       bool     `json:"enabled"`
}

// Service handles Active Directory operations
type Service struct {
	config    *ADConfig
	mu        sync.RWMutex
	available bool
}

var (
	globalService *Service
	once          sync.Once
)

// Initialize initializes the AD service
func Initialize(config *ADConfig) (*Service, error) {
	once.Do(func() {
		if config == nil {
			config = &ADConfig{
				Enabled: false,
				Port:    389,
				UserFilter: "(&(objectClass=user)(sAMAccountName=%s))",
				GroupFilter: "(&(objectClass=group)(member=%s))",
			}
		}

		globalService = &Service{
			config:    config,
			available: config.Enabled,
		}
	})

	return globalService, nil
}

// GetService returns the global AD service
func GetService() *Service {
	return globalService
}

// IsAvailable returns whether AD is available
func (s *Service) IsAvailable() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.available && s.config.Enabled
}

// UpdateConfig updates the AD configuration
func (s *Service) UpdateConfig(config *ADConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.config = config
	s.available = config.Enabled

	return nil
}

// GetConfig returns the current configuration (without password)
func (s *Service) GetConfig() *ADConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return config without password
	configCopy := *s.config
	configCopy.BindPassword = "***"
	return &configCopy
}

// TestConnection tests the AD connection
func (s *Service) TestConnection(ctx context.Context) error {
	s.mu.RLock()
	config := s.config
	s.mu.RUnlock()

	if !config.Enabled {
		return fmt.Errorf("AD is not enabled")
	}

	conn, err := s.connect(config)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Try to bind with the configured user
	if err := conn.Bind(config.BindUser, config.BindPassword); err != nil {
		return fmt.Errorf("failed to bind: %w", err)
	}

	return nil
}

// Authenticate authenticates a user against AD
func (s *Service) Authenticate(ctx context.Context, username, password string) (*ADUser, error) {
	s.mu.RLock()
	config := s.config
	s.mu.RUnlock()

	if !config.Enabled {
		return nil, fmt.Errorf("AD is not enabled")
	}

	conn, err := s.connect(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AD: %w", err)
	}
	defer conn.Close()

	// Bind with service account first
	if err := conn.Bind(config.BindUser, config.BindPassword); err != nil {
		return nil, fmt.Errorf("failed to bind with service account: %w", err)
	}

	// Search for the user
	searchFilter := fmt.Sprintf(config.UserFilter, ldap.EscapeFilter(username))
	searchRequest := ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		searchFilter,
		[]string{"dn", "sAMAccountName", "mail", "displayName", "userAccountControl"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search for user: %w", err)
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	if len(result.Entries) > 1 {
		return nil, fmt.Errorf("multiple users found")
	}

	entry := result.Entries[0]
	userDN := entry.DN

	// Try to bind with the user's credentials
	if err := conn.Bind(userDN, password); err != nil {
		return nil, fmt.Errorf("authentication failed")
	}

	// Get user details
	user := &ADUser{
		Username:          entry.GetAttributeValue("sAMAccountName"),
		Email:             entry.GetAttributeValue("mail"),
		DisplayName:       entry.GetAttributeValue("displayName"),
		DistinguishedName: userDN,
		Enabled:           !isUserDisabled(entry.GetAttributeValue("userAccountControl")),
	}

	// Get user groups
	groups, err := s.getUserGroups(conn, userDN, config)
	if err == nil {
		user.Groups = groups
	}

	return user, nil
}

// ListUsers lists all users from AD
func (s *Service) ListUsers(ctx context.Context) ([]*ADUser, error) {
	s.mu.RLock()
	config := s.config
	s.mu.RUnlock()

	if !config.Enabled {
		return nil, fmt.Errorf("AD is not enabled")
	}

	conn, err := s.connect(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AD: %w", err)
	}
	defer conn.Close()

	// Bind with service account
	if err := conn.Bind(config.BindUser, config.BindPassword); err != nil {
		return nil, fmt.Errorf("failed to bind: %w", err)
	}

	// Search for all users
	searchFilter := "(&(objectClass=user)(objectCategory=person))"
	searchRequest := ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		searchFilter,
		[]string{"dn", "sAMAccountName", "mail", "displayName", "userAccountControl"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search for users: %w", err)
	}

	users := make([]*ADUser, 0, len(result.Entries))
	for _, entry := range result.Entries {
		user := &ADUser{
			Username:          entry.GetAttributeValue("sAMAccountName"),
			Email:             entry.GetAttributeValue("mail"),
			DisplayName:       entry.GetAttributeValue("displayName"),
			DistinguishedName: entry.DN,
			Enabled:           !isUserDisabled(entry.GetAttributeValue("userAccountControl")),
		}

		// Get user groups
		groups, err := s.getUserGroups(conn, entry.DN, config)
		if err == nil {
			user.Groups = groups
		}

		users = append(users, user)
	}

	return users, nil
}

// SyncUser synchronizes a user from AD
func (s *Service) SyncUser(ctx context.Context, username string) (*ADUser, error) {
	s.mu.RLock()
	config := s.config
	s.mu.RUnlock()

	if !config.Enabled {
		return nil, fmt.Errorf("AD is not enabled")
	}

	conn, err := s.connect(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AD: %w", err)
	}
	defer conn.Close()

	// Bind with service account
	if err := conn.Bind(config.BindUser, config.BindPassword); err != nil {
		return nil, fmt.Errorf("failed to bind: %w", err)
	}

	// Search for the user
	searchFilter := fmt.Sprintf(config.UserFilter, ldap.EscapeFilter(username))
	searchRequest := ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		searchFilter,
		[]string{"dn", "sAMAccountName", "mail", "displayName", "userAccountControl"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to search for user: %w", err)
	}

	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	entry := result.Entries[0]

	user := &ADUser{
		Username:          entry.GetAttributeValue("sAMAccountName"),
		Email:             entry.GetAttributeValue("mail"),
		DisplayName:       entry.GetAttributeValue("displayName"),
		DistinguishedName: entry.DN,
		Enabled:           !isUserDisabled(entry.GetAttributeValue("userAccountControl")),
	}

	// Get user groups
	groups, err := s.getUserGroups(conn, entry.DN, config)
	if err == nil {
		user.Groups = groups
	}

	return user, nil
}

// connect establishes a connection to the AD server
func (s *Service) connect(config *ADConfig) (*ldap.Conn, error) {
	address := fmt.Sprintf("%s:%d", config.Server, config.Port)

	if config.UseTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: config.SkipVerify,
		}
		return ldap.DialTLS("tcp", address, tlsConfig)
	}

	return ldap.Dial("tcp", address)
}

// getUserGroups retrieves the groups a user belongs to
func (s *Service) getUserGroups(conn *ldap.Conn, userDN string, config *ADConfig) ([]string, error) {
	searchFilter := fmt.Sprintf(config.GroupFilter, ldap.EscapeFilter(userDN))
	searchRequest := ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		searchFilter,
		[]string{"cn"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	groups := make([]string, 0, len(result.Entries))
	for _, entry := range result.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}

	return groups, nil
}

// isUserDisabled checks if a user is disabled based on userAccountControl
func isUserDisabled(userAccountControl string) bool {
	if userAccountControl == "" {
		return false
	}

	// UAC_ACCOUNTDISABLE flag is 0x0002
	// This is a simplified check
	return strings.Contains(userAccountControl, "514") || strings.Contains(userAccountControl, "546")
}

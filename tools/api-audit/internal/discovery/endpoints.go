// Revision: 2025-11-16 | Author: Claude | Version: 1.0.0
package discovery

// Endpoint represents a discovered API endpoint
type Endpoint struct {
	Method       string
	Path         string
	Description  string
	AuthRequired bool
	Destructive  bool
}

// KnownEndpoints returns all known API endpoints
// This is manually curated based on backend/internal/api/router.go
func KnownEndpoints() []Endpoint {
	return []Endpoint{
		// Health & System
		{Method: "GET", Path: "/health", Description: "Health check", AuthRequired: false},
		{Method: "GET", Path: "/api/v1/system/info", Description: "System information", AuthRequired: true},
		{Method: "GET", Path: "/api/v1/system/metrics", Description: "System metrics", AuthRequired: true},

		// Auth
		{Method: "POST", Path: "/api/v1/auth/login", Description: "User login", AuthRequired: false},
		{Method: "POST", Path: "/api/v1/auth/refresh", Description: "Refresh token", AuthRequired: false},
		{Method: "POST", Path: "/api/v1/auth/logout", Description: "User logout", AuthRequired: true},

		// Users
		{Method: "GET", Path: "/api/v1/users", Description: "List users", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/users", Description: "Create user", AuthRequired: true},
		{Method: "GET", Path: "/api/v1/users/1", Description: "Get user by ID", AuthRequired: true},
		{Method: "PUT", Path: "/api/v1/users/1", Description: "Update user", AuthRequired: true},
		{Method: "DELETE", Path: "/api/v1/users/1", Description: "Delete user", AuthRequired: true, Destructive: true},

		// Groups
		{Method: "GET", Path: "/api/v1/groups", Description: "List groups", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/groups", Description: "Create group", AuthRequired: true},
		{Method: "GET", Path: "/api/v1/groups/1", Description: "Get group by ID", AuthRequired: true},
		{Method: "PUT", Path: "/api/v1/groups/1", Description: "Update group", AuthRequired: true},
		{Method: "DELETE", Path: "/api/v1/groups/1", Description: "Delete group", AuthRequired: true, Destructive: true},

		// Storage - ZFS
		{Method: "GET", Path: "/api/v1/syslib/zfs/pools", Description: "List ZFS pools", AuthRequired: true},
		{Method: "GET", Path: "/api/v1/syslib/zfs/pools/tank", Description: "Get ZFS pool details", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/syslib/zfs/pools", Description: "Create ZFS pool", AuthRequired: true, Destructive: true},
		{Method: "DELETE", Path: "/api/v1/syslib/zfs/pools/tank", Description: "Delete ZFS pool", AuthRequired: true, Destructive: true},
		{Method: "GET", Path: "/api/v1/syslib/zfs/snapshots", Description: "List ZFS snapshots", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/syslib/zfs/snapshots", Description: "Create ZFS snapshot", AuthRequired: true},
		{Method: "DELETE", Path: "/api/v1/syslib/zfs/snapshots/tank@snap1", Description: "Delete ZFS snapshot", AuthRequired: true, Destructive: true},

		// Storage - SMART
		{Method: "GET", Path: "/api/v1/syslib/smart/disks", Description: "List SMART disks", AuthRequired: true},
		{Method: "GET", Path: "/api/v1/syslib/smart/disks/sda", Description: "Get SMART disk details", AuthRequired: true},

		// Sharing - Samba
		{Method: "GET", Path: "/api/v1/syslib/samba/status", Description: "Get Samba status", AuthRequired: true},
		{Method: "GET", Path: "/api/v1/syslib/samba/shares", Description: "List Samba shares", AuthRequired: true},
		{Method: "GET", Path: "/api/v1/syslib/samba/shares/public", Description: "Get Samba share details", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/syslib/samba/shares", Description: "Create Samba share", AuthRequired: true},
		{Method: "PUT", Path: "/api/v1/syslib/samba/shares/public", Description: "Update Samba share", AuthRequired: true},
		{Method: "DELETE", Path: "/api/v1/syslib/samba/shares/public", Description: "Delete Samba share", AuthRequired: true, Destructive: true},

		// Sharing - NFS
		{Method: "GET", Path: "/api/v1/syslib/nfs/exports", Description: "List NFS exports", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/syslib/nfs/exports", Description: "Create NFS export", AuthRequired: true},
		{Method: "DELETE", Path: "/api/v1/syslib/nfs/exports/1", Description: "Delete NFS export", AuthRequired: true, Destructive: true},

		// Network
		{Method: "GET", Path: "/api/v1/syslib/network/interfaces", Description: "List network interfaces", AuthRequired: true},
		{Method: "GET", Path: "/api/v1/syslib/network/dns", Description: "Get DNS configuration", AuthRequired: true},
		{Method: "PUT", Path: "/api/v1/syslib/network/dns", Description: "Update DNS configuration", AuthRequired: true},
		{Method: "GET", Path: "/api/v1/syslib/network/firewall", Description: "Get firewall status", AuthRequired: true},

		// Services
		{Method: "GET", Path: "/api/v1/syslib/services/smbd/status", Description: "Get Samba service status", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/syslib/services/smbd/start", Description: "Start Samba service", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/syslib/services/smbd/stop", Description: "Stop Samba service", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/syslib/services/smbd/restart", Description: "Restart Samba service", AuthRequired: true},

		// Docker
		{Method: "GET", Path: "/api/v1/docker/containers", Description: "List Docker containers", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/docker/containers", Description: "Create Docker container", AuthRequired: true},

		// Files
		{Method: "GET", Path: "/api/v1/files/browse", Description: "Browse files", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/files/upload", Description: "Upload file", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/files/download", Description: "Download file", AuthRequired: true},
		{Method: "DELETE", Path: "/api/v1/files", Description: "Delete files", AuthRequired: true, Destructive: true},

		// Backup
		{Method: "GET", Path: "/api/v1/backup/jobs", Description: "List backup jobs", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/backup/jobs", Description: "Create backup job", AuthRequired: true},
		{Method: "GET", Path: "/api/v1/backup/jobs/1", Description: "Get backup job details", AuthRequired: true},

		// Scheduler
		{Method: "GET", Path: "/api/v1/tasks", Description: "List scheduled tasks", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/tasks", Description: "Create scheduled task", AuthRequired: true},
		{Method: "GET", Path: "/api/v1/tasks/1", Description: "Get task details", AuthRequired: true},
		{Method: "PUT", Path: "/api/v1/tasks/1", Description: "Update task", AuthRequired: true},
		{Method: "DELETE", Path: "/api/v1/tasks/1", Description: "Delete task", AuthRequired: true, Destructive: true},

		// Metrics
		{Method: "GET", Path: "/metrics", Description: "Prometheus metrics", AuthRequired: false},

		// Updates
		{Method: "GET", Path: "/api/v1/updates/check", Description: "Check for updates", AuthRequired: true},

		// Alerts
		{Method: "GET", Path: "/api/v1/alerts/configs", Description: "List alert configs", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/alerts/configs", Description: "Create alert config", AuthRequired: true},

		// Audit
		{Method: "GET", Path: "/api/v1/audit/logs", Description: "List audit logs", AuthRequired: true},

		// 2FA
		{Method: "POST", Path: "/api/v1/2fa/setup", Description: "Setup 2FA", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/2fa/verify", Description: "Verify 2FA", AuthRequired: true},
		{Method: "DELETE", Path: "/api/v1/2fa", Description: "Disable 2FA", AuthRequired: true},

		// Plugins
		{Method: "GET", Path: "/api/v1/plugins", Description: "List plugins", AuthRequired: true},
		{Method: "POST", Path: "/api/v1/plugins/upload", Description: "Upload plugin", AuthRequired: true},
	}
}

// FilterEndpoints filters endpoints based on options
func FilterEndpoints(endpoints []Endpoint, includeDestructive bool, category string) []Endpoint {
	var filtered []Endpoint

	for _, ep := range endpoints {
		// Skip destructive endpoints if not included
		if ep.Destructive && !includeDestructive {
			continue
		}

		// Filter by category if specified
		if category != "" && category != "all" {
			if !endpointMatchesCategory(ep, category) {
				continue
			}
		}

		filtered = append(filtered, ep)
	}

	return filtered
}

func endpointMatchesCategory(ep Endpoint, category string) bool {
	switch category {
	case "storage":
		return containsAny(ep.Path, "/syslib/zfs", "/syslib/smart")
	case "sharing":
		return containsAny(ep.Path, "/syslib/samba", "/syslib/nfs")
	case "network":
		return containsAny(ep.Path, "/syslib/network")
	case "users":
		return containsAny(ep.Path, "/users", "/groups")
	case "system":
		return containsAny(ep.Path, "/system")
	default:
		return true
	}
}

func containsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) && s[:len(substr)] == substr ||
		   len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
		   findSubstring(s, substr) {
			return true
		}
	}
	return false
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

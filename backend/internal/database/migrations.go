// Revision: 2025-11-23 | Author: Claude | Version: 1.2.0
package database

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// RunMigrations runs all database migrations
func RunMigrations() error {
	logger.Info("Running database migrations...")

	// Auto-migrate models
	if err := DB.AutoMigrate(
		&models.User{},
		&models.UserGroup{},
		&models.Share{},
		&models.Volume{},
		&models.DiskLabel{},
		&models.AuditLog{},
		&models.FailedLoginAttempt{},
		&models.IPBlock{},
		&models.AlertConfig{},
		&models.AlertLog{},
		&models.ScheduledTask{},
		&models.TaskExecution{},
		&models.TwoFactorAuth{},
		&models.TwoFactorBackupCode{},
		&models.TwoFactorAttempt{},
		&models.SystemMetric{},
		&models.HealthScore{},
		&models.MonitoringConfig{},
		&models.AddonInstallation{},
		// VPN Server models
		&models.VPNProtocolConfig{},
		&models.VPNPeer{},
		&models.VPNCertificate{},
		&models.VPNUser{},
		&models.VPNConnection{},
		&models.VPNRoute{},
		&models.VPNFirewallRule{},
		// Network configuration models
		&models.NetworkBridge{},
		&models.NetworkInterface{},
		&models.NetworkSnapshot{},
		// LXC container models
		&models.LXCContainer{},
		// Add more models here as they are created
	); err != nil {
		return err
	}

	logger.Info("Database migrations completed successfully")

	// Add performance indexes
	if err := AddPerformanceIndexes(); err != nil {
		logger.Warn("Failed to add performance indexes (non-fatal)", zap.Error(err))
	}

	// NOTE: Default admin user creation removed.
	// Users must now use the Setup Wizard on first access to create the initial admin account.

	return nil
}

// AddPerformanceIndexes adds database indexes for improved query performance
func AddPerformanceIndexes() error {
	logger.Info("Adding performance indexes...")

	// Index for user username lookups (login queries)
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)").Error; err != nil {
		return err
	}

	// Indexes for audit log queries
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(created_at DESC)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id)").Error; err != nil {
		return err
	}
	// Composite index for filtered time-based queries
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_user_timestamp ON audit_logs(user_id, created_at DESC)").Error; err != nil {
		return err
	}

	// Index for share name lookups
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_shares_name ON shares(name)").Error; err != nil {
		return err
	}

	// Index for failed login attempts by IP
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_failed_logins_ip ON failed_login_attempts(ip_address)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_failed_logins_timestamp ON failed_login_attempts(attempted_at DESC)").Error; err != nil {
		return err
	}

	// Index for scheduled tasks
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_enabled ON scheduled_tasks(enabled)").Error; err != nil {
		return err
	}

	// Index for system metrics timestamp
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON system_metrics(timestamp DESC)").Error; err != nil {
		return err
	}

	// Index for volume name lookups
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_volumes_name ON volumes(name)").Error; err != nil {
		return err
	}
	// Index for volume mount point lookups
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_volumes_mount_point ON volumes(mount_point)").Error; err != nil {
		return err
	}
	// Index for volume status queries
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_volumes_status ON volumes(status)").Error; err != nil {
		return err
	}

	// VPN Server indexes
	// Index for VPN protocol lookups
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_protocols_protocol ON vpn_protocol_configs(protocol)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_protocols_enabled ON vpn_protocol_configs(enabled)").Error; err != nil {
		return err
	}

	// Index for VPN peer lookups
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_peers_protocol_id ON vpn_peers(protocol_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_peers_enabled ON vpn_peers(enabled)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_peers_public_key ON vpn_peers(public_key)").Error; err != nil {
		return err
	}

	// Index for VPN certificate lookups
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_certs_protocol_id ON vpn_certificates(protocol_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_certs_status ON vpn_certificates(status)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_certs_serial ON vpn_certificates(serial_number)").Error; err != nil {
		return err
	}

	// Index for VPN user lookups
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_users_username ON vpn_users(username)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_users_enabled ON vpn_users(enabled)").Error; err != nil {
		return err
	}

	// Index for VPN connection queries
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_connections_user_id ON vpn_connections(user_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_connections_active ON vpn_connections(active)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_connections_protocol ON vpn_connections(protocol)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_connections_connected_at ON vpn_connections(connected_at DESC)").Error; err != nil {
		return err
	}
	// Composite index for active connections by protocol
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_connections_active_protocol ON vpn_connections(active, protocol)").Error; err != nil {
		return err
	}

	// Index for VPN route lookups
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_routes_protocol_id ON vpn_routes(protocol_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_routes_enabled ON vpn_routes(enabled)").Error; err != nil {
		return err
	}

	// Index for VPN firewall rule lookups
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_firewall_protocol_id ON vpn_firewall_rules(protocol_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_firewall_enabled ON vpn_firewall_rules(enabled)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_vpn_firewall_priority ON vpn_firewall_rules(priority ASC)").Error; err != nil {
		return err
	}

	// Network bridge indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_network_bridges_name ON network_bridges(name)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_network_bridges_autostart ON network_bridges(autostart)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_network_bridges_status ON network_bridges(status)").Error; err != nil {
		return err
	}

	// Network interface indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_network_interfaces_name ON network_interfaces(name)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_network_interfaces_autostart ON network_interfaces(autostart)").Error; err != nil {
		return err
	}

	// LXC container indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_lxc_containers_name ON lxc_containers(name)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_lxc_containers_autostart ON lxc_containers(autostart)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_lxc_containers_status ON lxc_containers(status)").Error; err != nil {
		return err
	}

	logger.Info("Performance indexes added successfully")
	return nil
}

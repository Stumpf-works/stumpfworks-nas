// Revision: 2025-11-28 | Author: Claude | Version: 1.0.0
// Package filesystem provides disk quota management for users and groups
package filesystem

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
)

// QuotaManager manages disk quotas for users and groups
type QuotaManager struct {
	shell   executor.ShellExecutor
	enabled bool
}

// QuotaType represents the type of quota (user or group)
type QuotaType string

const (
	UserQuota  QuotaType = "user"
	GroupQuota QuotaType = "group"
)

// QuotaInfo represents quota information for a user or group
type QuotaInfo struct {
	Name         string    `json:"name"`          // username or groupname
	Type         QuotaType `json:"type"`          // user or group
	Filesystem   string    `json:"filesystem"`    // filesystem path
	BlocksUsed   uint64    `json:"blocks_used"`   // blocks currently used (KB)
	BlocksSoft   uint64    `json:"blocks_soft"`   // soft limit for blocks (KB)
	BlocksHard   uint64    `json:"blocks_hard"`   // hard limit for blocks (KB)
	InodesUsed   uint64    `json:"inodes_used"`   // inodes currently used
	InodesSoft   uint64    `json:"inodes_soft"`   // soft limit for inodes
	InodesHard   uint64    `json:"inodes_hard"`   // hard limit for inodes
	BlocksGrace  string    `json:"blocks_grace"`  // grace period for blocks
	InodesGrace  string    `json:"inodes_grace"`  // grace period for inodes
}

// QuotaLimits represents quota limits to be set
type QuotaLimits struct {
	BlocksSoft uint64 `json:"blocks_soft"` // soft limit in KB
	BlocksHard uint64 `json:"blocks_hard"` // hard limit in KB
	InodesSoft uint64 `json:"inodes_soft"` // soft limit for inodes
	InodesHard uint64 `json:"inodes_hard"` // hard limit for inodes
}

// FilesystemQuotaStatus represents quota status for a filesystem
type FilesystemQuotaStatus struct {
	Filesystem    string `json:"filesystem"`
	QuotasEnabled bool   `json:"quotas_enabled"`
	UserQuotas    bool   `json:"user_quotas"`
	GroupQuotas   bool   `json:"group_quotas"`
}

// NewQuotaManager creates a new quota manager
func NewQuotaManager(shell executor.ShellExecutor) (*QuotaManager, error) {
	// Check if quota tools are available
	if !shell.CommandExists("quota") {
		return nil, fmt.Errorf("quota tools not installed (install 'quota' package)")
	}

	return &QuotaManager{
		shell:   shell,
		enabled: true,
	}, nil
}

// IsEnabled returns whether quota support is available
func (q *QuotaManager) IsEnabled() bool {
	return q.enabled
}

// GetUserQuota retrieves quota information for a user
func (q *QuotaManager) GetUserQuota(username string, filesystem string) (*QuotaInfo, error) {
	if !q.enabled {
		return nil, fmt.Errorf("quota support not available")
	}

	result, err := q.shell.Execute("quota", "-u", "-w", username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user quota: %w", err)
	}

	return q.parseQuotaOutput(result.Stdout, username, UserQuota, filesystem)
}

// GetGroupQuota retrieves quota information for a group
func (q *QuotaManager) GetGroupQuota(groupname string, filesystem string) (*QuotaInfo, error) {
	if !q.enabled {
		return nil, fmt.Errorf("quota support not available")
	}

	result, err := q.shell.Execute("quota", "-g", "-w", groupname)
	if err != nil {
		return nil, fmt.Errorf("failed to get group quota: %w", err)
	}

	return q.parseQuotaOutput(result.Stdout, groupname, GroupQuota, filesystem)
}

// SetUserQuota sets quota limits for a user
func (q *QuotaManager) SetUserQuota(username string, filesystem string, limits QuotaLimits) error {
	if !q.enabled {
		return fmt.Errorf("quota support not available")
	}

	// setquota -u username block-soft block-hard inode-soft inode-hard filesystem
	args := []string{
		"-u",
		username,
		fmt.Sprintf("%d", limits.BlocksSoft),
		fmt.Sprintf("%d", limits.BlocksHard),
		fmt.Sprintf("%d", limits.InodesSoft),
		fmt.Sprintf("%d", limits.InodesHard),
		filesystem,
	}

	result, err := q.shell.Execute("setquota", args...)
	if err != nil {
		return fmt.Errorf("failed to set user quota: %s - %w", result.Stderr, err)
	}

	return nil
}

// SetGroupQuota sets quota limits for a group
func (q *QuotaManager) SetGroupQuota(groupname string, filesystem string, limits QuotaLimits) error {
	if !q.enabled {
		return fmt.Errorf("quota support not available")
	}

	// setquota -g groupname block-soft block-hard inode-soft inode-hard filesystem
	args := []string{
		"-g",
		groupname,
		fmt.Sprintf("%d", limits.BlocksSoft),
		fmt.Sprintf("%d", limits.BlocksHard),
		fmt.Sprintf("%d", limits.InodesSoft),
		fmt.Sprintf("%d", limits.InodesHard),
		filesystem,
	}

	result, err := q.shell.Execute("setquota", args...)
	if err != nil {
		return fmt.Errorf("failed to set group quota: %s - %w", result.Stderr, err)
	}

	return nil
}

// RemoveUserQuota removes quota limits for a user (sets to 0)
func (q *QuotaManager) RemoveUserQuota(username string, filesystem string) error {
	if !q.enabled {
		return fmt.Errorf("quota support not available")
	}

	return q.SetUserQuota(username, filesystem, QuotaLimits{
		BlocksSoft: 0,
		BlocksHard: 0,
		InodesSoft: 0,
		InodesHard: 0,
	})
}

// RemoveGroupQuota removes quota limits for a group (sets to 0)
func (q *QuotaManager) RemoveGroupQuota(groupname string, filesystem string) error {
	if !q.enabled {
		return fmt.Errorf("quota support not available")
	}

	return q.SetGroupQuota(groupname, filesystem, QuotaLimits{
		BlocksSoft: 0,
		BlocksHard: 0,
		InodesSoft: 0,
		InodesHard: 0,
	})
}

// ListUserQuotas lists all user quotas on a filesystem
func (q *QuotaManager) ListUserQuotas(filesystem string) ([]QuotaInfo, error) {
	if !q.enabled {
		return nil, fmt.Errorf("quota support not available")
	}

	result, err := q.shell.Execute("repquota", "-u", "-v", filesystem)
	if err != nil {
		return nil, fmt.Errorf("failed to list user quotas: %w", err)
	}

	return q.parseRepquotaOutput(result.Stdout, UserQuota, filesystem)
}

// ListGroupQuotas lists all group quotas on a filesystem
func (q *QuotaManager) ListGroupQuotas(filesystem string) ([]QuotaInfo, error) {
	if !q.enabled {
		return nil, fmt.Errorf("quota support not available")
	}

	result, err := q.shell.Execute("repquota", "-g", "-v", filesystem)
	if err != nil {
		return nil, fmt.Errorf("failed to list group quotas: %w", err)
	}

	return q.parseRepquotaOutput(result.Stdout, GroupQuota, filesystem)
}

// GetFilesystemQuotaStatus checks if quotas are enabled on a filesystem
func (q *QuotaManager) GetFilesystemQuotaStatus(filesystem string) (*FilesystemQuotaStatus, error) {
	if !q.enabled {
		return nil, fmt.Errorf("quota support not available")
	}

	status := &FilesystemQuotaStatus{
		Filesystem:    filesystem,
		QuotasEnabled: false,
		UserQuotas:    false,
		GroupQuotas:   false,
	}

	// Check user quotas
	result, err := q.shell.Execute("quotaon", "-p", "-u", filesystem)
	if err == nil && strings.Contains(result.Stdout, "is on") {
		status.QuotasEnabled = true
		status.UserQuotas = true
	}

	// Check group quotas
	result, err = q.shell.Execute("quotaon", "-p", "-g", filesystem)
	if err == nil && strings.Contains(result.Stdout, "is on") {
		status.QuotasEnabled = true
		status.GroupQuotas = true
	}

	return status, nil
}

// parseQuotaOutput parses output from 'quota' command
func (q *QuotaManager) parseQuotaOutput(output string, name string, quotaType QuotaType, filesystem string) (*QuotaInfo, error) {
	info := &QuotaInfo{
		Name:       name,
		Type:       quotaType,
		Filesystem: filesystem,
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Disk quotas") || strings.HasPrefix(line, "Filesystem") {
			continue
		}

		// Parse quota line format: filesystem blocks quota limit grace files quota limit grace
		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}

		// Only process the line for our filesystem
		if !strings.Contains(fields[0], filesystem) {
			continue
		}

		// Parse blocks (field 1 = used, 2 = soft, 3 = hard)
		if blocksUsed, err := parseSize(fields[1]); err == nil {
			info.BlocksUsed = blocksUsed
		}
		if blocksSoft, err := parseSize(fields[2]); err == nil {
			info.BlocksSoft = blocksSoft
		}
		if blocksHard, err := parseSize(fields[3]); err == nil {
			info.BlocksHard = blocksHard
		}
		info.BlocksGrace = fields[4]

		// Parse inodes (field 5 = used, 6 = soft, 7 = hard)
		if inodesUsed, err := strconv.ParseUint(fields[5], 10, 64); err == nil {
			info.InodesUsed = inodesUsed
		}
		if inodesSoft, err := strconv.ParseUint(fields[6], 10, 64); err == nil {
			info.InodesSoft = inodesSoft
		}
		if inodesHard, err := strconv.ParseUint(fields[7], 10, 64); err == nil {
			info.InodesHard = inodesHard
		}
		if len(fields) > 8 {
			info.InodesGrace = fields[8]
		}

		break
	}

	return info, nil
}

// parseRepquotaOutput parses output from 'repquota' command
func (q *QuotaManager) parseRepquotaOutput(output string, quotaType QuotaType, filesystem string) ([]QuotaInfo, error) {
	var quotas []QuotaInfo

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip header lines
		if line == "" || strings.HasPrefix(line, "***") || strings.HasPrefix(line, "Block grace") ||
			strings.HasPrefix(line, "Report for") || strings.Contains(line, "Block limits") {
			continue
		}

		// Parse repquota line format: name -- blocks soft hard grace files soft hard grace
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		// Skip if name starts with # (system accounts)
		if strings.HasPrefix(fields[0], "#") {
			continue
		}

		info := QuotaInfo{
			Name:       fields[0],
			Type:       quotaType,
			Filesystem: filesystem,
		}

		// Parse blocks (skip "--" at field 1)
		if blocksUsed, err := parseSize(fields[2]); err == nil {
			info.BlocksUsed = blocksUsed
		}
		if blocksSoft, err := parseSize(fields[3]); err == nil {
			info.BlocksSoft = blocksSoft
		}
		if blocksHard, err := parseSize(fields[4]); err == nil {
			info.BlocksHard = blocksHard
		}
		info.BlocksGrace = fields[5]

		// Parse inodes
		if inodesUsed, err := strconv.ParseUint(fields[6], 10, 64); err == nil {
			info.InodesUsed = inodesUsed
		}
		if inodesSoft, err := strconv.ParseUint(fields[7], 10, 64); err == nil {
			info.InodesSoft = inodesSoft
		}
		if inodesHard, err := strconv.ParseUint(fields[8], 10, 64); err == nil {
			info.InodesHard = inodesHard
		}
		if len(fields) > 9 {
			info.InodesGrace = fields[9]
		}

		// Only add if quota is actually set (not all zeros)
		if info.BlocksSoft > 0 || info.BlocksHard > 0 || info.InodesSoft > 0 || info.InodesHard > 0 {
			quotas = append(quotas, info)
		}
	}

	return quotas, nil
}

// parseSize parses a size string (handles K, M, G suffixes and plain numbers)
func parseSize(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "0" || s == "-" {
		return 0, nil
	}

	// Check for suffix
	multiplier := uint64(1)
	if len(s) > 0 {
		lastChar := s[len(s)-1]
		switch lastChar {
		case 'K', 'k':
			multiplier = 1
			s = s[:len(s)-1]
		case 'M', 'm':
			multiplier = 1024
			s = s[:len(s)-1]
		case 'G', 'g':
			multiplier = 1024 * 1024
			s = s[:len(s)-1]
		case 'T', 't':
			multiplier = 1024 * 1024 * 1024
			s = s[:len(s)-1]
		}
	}

	value, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return value * multiplier, nil
}

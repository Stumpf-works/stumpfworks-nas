package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/audit"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// FailedLoginService handles failed login attempt tracking and IP blocking
type FailedLoginService struct {
	db *gorm.DB
	mu sync.RWMutex

	// Configuration
	maxAttempts      int           // Max attempts before blocking
	blockDuration    time.Duration // How long to block IP
	attemptWindow    time.Duration // Time window for counting attempts
	cleanupInterval  time.Duration // How often to clean old records
	stopCleanup      chan bool
}

var (
	globalFailedLoginService *FailedLoginService
	failedLoginOnce          sync.Once
)

// InitializeFailedLoginService initializes the failed login tracking service
func InitializeFailedLoginService() (*FailedLoginService, error) {
	failedLoginOnce.Do(func() {
		globalFailedLoginService = &FailedLoginService{
			db:               database.GetDB(),
			maxAttempts:      5,                // 5 failed attempts
			blockDuration:    15 * time.Minute, // Block for 15 minutes
			attemptWindow:    15 * time.Minute, // Count attempts in last 15 minutes
			cleanupInterval:  1 * time.Hour,    // Cleanup every hour
			stopCleanup:      make(chan bool),
		}

		// Start background cleanup task
		go globalFailedLoginService.startCleanupTask()
	})

	return globalFailedLoginService, nil
}

// GetFailedLoginService returns the global failed login service
func GetFailedLoginService() *FailedLoginService {
	if globalFailedLoginService == nil {
		globalFailedLoginService, _ = InitializeFailedLoginService()
	}
	return globalFailedLoginService
}

// RecordFailedAttempt records a failed login attempt
func (s *FailedLoginService) RecordFailedAttempt(ctx context.Context, username, ipAddress, userAgent, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create failed login attempt record
	attempt := &models.FailedLoginAttempt{
		Username:  username,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Reason:    reason,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.db.Create(attempt).Error; err != nil {
		logger.Error("Failed to record login attempt",
			zap.Error(err),
			zap.String("username", username),
			zap.String("ip", ipAddress))
		return fmt.Errorf("failed to record attempt: %w", err)
	}

	// Log to audit system
	auditService := audit.GetService()
	if auditService != nil {
		_ = auditService.LogWithDetails(ctx, nil, username, models.ActionAuthLoginFailed, "auth/login",
			models.StatusFailure, models.SeverityWarning, "Failed login attempt",
			map[string]interface{}{
				"reason":     reason,
				"ip_address": ipAddress,
			})
	}

	// Check if we should block this IP
	if err := s.checkAndBlockIP(ctx, ipAddress, username); err != nil {
		logger.Error("Failed to check/block IP", zap.Error(err))
	}

	return nil
}

// checkAndBlockIP checks if an IP should be blocked based on recent failed attempts
func (s *FailedLoginService) checkAndBlockIP(ctx context.Context, ipAddress, username string) error {
	// Count recent failed attempts from this IP
	cutoffTime := time.Now().UTC().Add(-s.attemptWindow)

	var attemptCount int64
	if err := s.db.Model(&models.FailedLoginAttempt{}).
		Where("ip_address = ? AND created_at > ?", ipAddress, cutoffTime).
		Count(&attemptCount).Error; err != nil {
		return fmt.Errorf("failed to count attempts: %w", err)
	}

	// If attempts exceed threshold, block the IP
	if attemptCount >= int64(s.maxAttempts) {
		// Check if already blocked
		var existingBlock models.IPBlock
		err := s.db.Where("ip_address = ? AND is_active = ?", ipAddress, true).First(&existingBlock).Error

		if err == gorm.ErrRecordNotFound {
			// Create new block
			block := &models.IPBlock{
				IPAddress:   ipAddress,
				Reason:      fmt.Sprintf("Too many failed login attempts (%d)", attemptCount),
				Attempts:    int(attemptCount),
				ExpiresAt:   time.Now().UTC().Add(s.blockDuration),
				IsActive:    true,
				IsPermanent: false,
			}

			if err := s.db.Create(block).Error; err != nil {
				return fmt.Errorf("failed to create IP block: %w", err)
			}

			// Mark all attempts as blocked
			if err := s.db.Model(&models.FailedLoginAttempt{}).
				Where("ip_address = ? AND blocked = ?", ipAddress, false).
				Updates(map[string]interface{}{
					"blocked":    true,
					"blocked_at": time.Now().UTC(),
				}).Error; err != nil {
				logger.Error("Failed to mark attempts as blocked", zap.Error(err))
			}

			// Log critical event
			auditService := audit.GetService()
			if auditService != nil {
				_ = auditService.LogWithDetails(ctx, nil, username, "security.ip_blocked", ipAddress,
					models.StatusSuccess, models.SeverityCritical,
					fmt.Sprintf("IP address %s blocked due to %d failed login attempts", ipAddress, attemptCount),
					map[string]interface{}{
						"ip_address":  ipAddress,
						"attempts":    attemptCount,
						"duration":    s.blockDuration.String(),
						"expires_at":  block.ExpiresAt,
					})
			}

			logger.Warn("IP address blocked",
				zap.String("ip", ipAddress),
				zap.Int64("attempts", attemptCount),
				zap.Duration("duration", s.blockDuration))
		}
	}

	return nil
}

// IsIPBlocked checks if an IP address is currently blocked
func (s *FailedLoginService) IsIPBlocked(ipAddress string) (bool, *models.IPBlock, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var block models.IPBlock
	err := s.db.Where("ip_address = ? AND is_active = ?", ipAddress, true).First(&block).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil, nil
	}

	if err != nil {
		return false, nil, fmt.Errorf("failed to check IP block: %w", err)
	}

	// Check if block has expired
	if block.IsExpired() && !block.IsPermanent {
		// Deactivate expired block
		if err := s.db.Model(&block).Update("is_active", false).Error; err != nil {
			logger.Error("Failed to deactivate expired block", zap.Error(err))
		}
		return false, nil, nil
	}

	return true, &block, nil
}

// UnblockIP removes the block on an IP address
func (s *FailedLoginService) UnblockIP(ctx context.Context, ipAddress string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := s.db.Model(&models.IPBlock{}).
		Where("ip_address = ?", ipAddress).
		Update("is_active", false)

	if result.Error != nil {
		return fmt.Errorf("failed to unblock IP: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		logger.Info("IP address unblocked", zap.String("ip", ipAddress))

		// Log audit event
		auditService := audit.GetService()
		if auditService != nil {
			_ = auditService.Log(ctx, &audit.LogEntry{
				Username: "system",
				Action:   "security.ip_unblocked",
				Resource: ipAddress,
				Status:   models.StatusSuccess,
				Severity: models.SeverityInfo,
				Message:  fmt.Sprintf("IP address %s unblocked", ipAddress),
			})
		}
	}

	return nil
}

// GetRecentFailedAttempts retrieves recent failed login attempts
func (s *FailedLoginService) GetRecentFailedAttempts(limit int, offset int) ([]*models.FailedLoginAttempt, int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var attempts []*models.FailedLoginAttempt
	var total int64

	// Get total count
	if err := s.db.Model(&models.FailedLoginAttempt{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count attempts: %w", err)
	}

	// Get paginated results
	if err := s.db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&attempts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve attempts: %w", err)
	}

	return attempts, total, nil
}

// GetBlockedIPs retrieves all currently blocked IPs
func (s *FailedLoginService) GetBlockedIPs() ([]*models.IPBlock, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var blocks []*models.IPBlock
	if err := s.db.Where("is_active = ?", true).
		Order("created_at DESC").
		Find(&blocks).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve blocked IPs: %w", err)
	}

	return blocks, nil
}

// GetStats retrieves failed login statistics
func (s *FailedLoginService) GetStats() (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[string]interface{})

	// Total failed attempts
	var totalAttempts int64
	if err := s.db.Model(&models.FailedLoginAttempt{}).Count(&totalAttempts).Error; err != nil {
		return nil, fmt.Errorf("failed to count total attempts: %w", err)
	}
	stats["total_attempts"] = totalAttempts

	// Attempts in last 24 hours
	last24h := time.Now().UTC().Add(-24 * time.Hour)
	var last24hAttempts int64
	if err := s.db.Model(&models.FailedLoginAttempt{}).
		Where("created_at > ?", last24h).
		Count(&last24hAttempts).Error; err != nil {
		return nil, fmt.Errorf("failed to count 24h attempts: %w", err)
	}
	stats["last_24h_attempts"] = last24hAttempts

	// Currently blocked IPs
	var blockedIPs int64
	if err := s.db.Model(&models.IPBlock{}).
		Where("is_active = ?", true).
		Count(&blockedIPs).Error; err != nil {
		return nil, fmt.Errorf("failed to count blocked IPs: %w", err)
	}
	stats["blocked_ips"] = blockedIPs

	// Top failed usernames
	type UsernameCount struct {
		Username string
		Count    int64
	}
	var topUsernames []UsernameCount
	if err := s.db.Model(&models.FailedLoginAttempt{}).
		Select("username, COUNT(*) as count").
		Group("username").
		Order("count DESC").
		Limit(10).
		Scan(&topUsernames).Error; err != nil {
		return nil, fmt.Errorf("failed to get top usernames: %w", err)
	}
	stats["top_failed_usernames"] = topUsernames

	return stats, nil
}

// startCleanupTask runs periodic cleanup of old records
func (s *FailedLoginService) startCleanupTask() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanup()
		case <-s.stopCleanup:
			return
		}
	}
}

// cleanup removes old failed login attempts and expired blocks
func (s *FailedLoginService) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Delete failed login attempts older than 30 days
	cutoffTime := time.Now().UTC().Add(-30 * 24 * time.Hour)
	result := s.db.Where("created_at < ?", cutoffTime).Delete(&models.FailedLoginAttempt{})
	if result.Error != nil {
		logger.Error("Failed to cleanup old login attempts", zap.Error(result.Error))
	} else if result.RowsAffected > 0 {
		logger.Info("Cleaned up old login attempts", zap.Int64("count", result.RowsAffected))
	}

	// Deactivate expired IP blocks
	now := time.Now().UTC()
	result = s.db.Model(&models.IPBlock{}).
		Where("is_active = ? AND is_permanent = ? AND expires_at < ?", true, false, now).
		Update("is_active", false)
	if result.Error != nil {
		logger.Error("Failed to deactivate expired blocks", zap.Error(result.Error))
	} else if result.RowsAffected > 0 {
		logger.Info("Deactivated expired IP blocks", zap.Int64("count", result.RowsAffected))
	}
}

// Stop stops the cleanup task
func (s *FailedLoginService) Stop() {
	close(s.stopCleanup)
}

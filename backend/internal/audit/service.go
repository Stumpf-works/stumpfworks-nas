// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Service handles audit logging operations
type Service struct {
	db *gorm.DB
	mu sync.RWMutex
}

var (
	globalService *Service
	once          sync.Once
)

// Initialize initializes the audit log service
func Initialize() (*Service, error) {
	once.Do(func() {
		globalService = &Service{
			db: database.GetDB(),
		}
	})

	return globalService, nil
}

// GetService returns the global audit log service
func GetService() *Service {
	if globalService == nil {
		globalService = &Service{
			db: database.GetDB(),
		}
	}
	return globalService
}

// LogEntry represents an audit log entry before persisting
type LogEntry struct {
	UserID    *uint
	Username  string
	Action    string
	Resource  string
	Status    string
	Severity  string
	IPAddress string
	UserAgent string
	Details   map[string]interface{}
	Message   string
}

// Log creates a new audit log entry
func (s *Service) Log(ctx context.Context, entry *LogEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Convert details map to JSON string
	var detailsJSON string
	if entry.Details != nil && len(entry.Details) > 0 {
		detailsBytes, err := json.Marshal(entry.Details)
		if err != nil {
			logger.Error("Failed to marshal audit log details", zap.Error(err))
		} else {
			detailsJSON = string(detailsBytes)
		}
	}

	// Create audit log entry
	auditLog := &models.AuditLog{
		UserID:    entry.UserID,
		Username:  entry.Username,
		Action:    entry.Action,
		Resource:  entry.Resource,
		Status:    entry.Status,
		Severity:  entry.Severity,
		IPAddress: entry.IPAddress,
		UserAgent: entry.UserAgent,
		Details:   detailsJSON,
		Message:   entry.Message,
		CreatedAt: time.Now().UTC(),
	}

	// Persist to database
	if err := s.db.Create(auditLog).Error; err != nil {
		logger.Error("Failed to create audit log entry",
			zap.Error(err),
			zap.String("action", entry.Action),
			zap.String("username", entry.Username))
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	// Log to application logger for immediate visibility
	logFields := []zap.Field{
		zap.Uint("audit_id", auditLog.ID),
		zap.String("action", entry.Action),
		zap.String("username", entry.Username),
		zap.String("status", entry.Status),
		zap.String("severity", entry.Severity),
	}

	if entry.Resource != "" {
		logFields = append(logFields, zap.String("resource", entry.Resource))
	}

	if entry.Message != "" {
		logFields = append(logFields, zap.String("message", entry.Message))
	}

	switch entry.Severity {
	case models.SeverityCritical:
		logger.Warn("Audit log (CRITICAL)", logFields...)
	case models.SeverityWarning:
		logger.Warn("Audit log (WARNING)", logFields...)
	default:
		logger.Info("Audit log", logFields...)
	}

	return nil
}

// LogFromRequest creates an audit log entry from an HTTP request
func (s *Service) LogFromRequest(r *http.Request, userID *uint, username, action, resource, status, severity, message string) error {
	return s.Log(r.Context(), &LogEntry{
		UserID:    userID,
		Username:  username,
		Action:    action,
		Resource:  resource,
		Status:    status,
		Severity:  severity,
		IPAddress: getIPAddress(r),
		UserAgent: r.UserAgent(),
		Message:   message,
	})
}

// LogWithDetails creates an audit log entry with additional details
func (s *Service) LogWithDetails(ctx context.Context, userID *uint, username, action, resource, status, severity, message string, details map[string]interface{}) error {
	return s.Log(ctx, &LogEntry{
		UserID:   userID,
		Username: username,
		Action:   action,
		Resource: resource,
		Status:   status,
		Severity: severity,
		Message:  message,
		Details:  details,
	})
}

// QueryParams represents audit log query parameters
type QueryParams struct {
	UserID    *uint
	Username  string
	Action    string
	Status    string
	Severity  string
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}

// Query retrieves audit logs based on query parameters
func (s *Service) Query(ctx context.Context, params *QueryParams) ([]*models.AuditLog, int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := s.db.Model(&models.AuditLog{})

	// Apply filters
	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}

	if params.Username != "" {
		query = query.Where("username LIKE ?", "%"+params.Username+"%")
	}

	if params.Action != "" {
		query = query.Where("action = ?", params.Action)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Severity != "" {
		query = query.Where("severity = ?", params.Severity)
	}

	if params.StartDate != nil {
		query = query.Where("created_at >= ?", *params.StartDate)
	}

	if params.EndDate != nil {
		query = query.Where("created_at <= ?", *params.EndDate)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Apply pagination and ordering
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	} else {
		query = query.Limit(100) // Default limit
	}

	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	query = query.Order("created_at DESC")

	// Execute query
	var logs []*models.AuditLog
	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to query audit logs: %w", err)
	}

	return logs, total, nil
}

// GetByID retrieves an audit log entry by ID
func (s *Service) GetByID(ctx context.Context, id uint) (*models.AuditLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var log models.AuditLog
	if err := s.db.First(&log, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("audit log not found")
		}
		return nil, fmt.Errorf("failed to retrieve audit log: %w", err)
	}

	return &log, nil
}

// GetRecent retrieves the most recent audit logs
func (s *Service) GetRecent(ctx context.Context, limit int) ([]*models.AuditLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 {
		limit = 50
	}

	var logs []*models.AuditLog
	if err := s.db.Order("created_at DESC").Limit(limit).Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve recent audit logs: %w", err)
	}

	return logs, nil
}

// GetStats retrieves audit log statistics
func (s *Service) GetStats(ctx context.Context) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[string]interface{})

	// Total count
	var total int64
	if err := s.db.Model(&models.AuditLog{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count total logs: %w", err)
	}
	stats["total"] = total

	// Count by severity
	var severityCounts []struct {
		Severity string
		Count    int64
	}
	if err := s.db.Model(&models.AuditLog{}).
		Select("severity, COUNT(*) as count").
		Group("severity").
		Scan(&severityCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to count by severity: %w", err)
	}
	stats["by_severity"] = severityCounts

	// Count by action
	var actionCounts []struct {
		Action string
		Count  int64
	}
	if err := s.db.Model(&models.AuditLog{}).
		Select("action, COUNT(*) as count").
		Group("action").
		Order("count DESC").
		Limit(10).
		Scan(&actionCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to count by action: %w", err)
	}
	stats["top_actions"] = actionCounts

	// Count last 24 hours
	last24h := time.Now().UTC().Add(-24 * time.Hour)
	var last24hCount int64
	if err := s.db.Model(&models.AuditLog{}).
		Where("created_at >= ?", last24h).
		Count(&last24hCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count last 24h: %w", err)
	}
	stats["last_24h"] = last24hCount

	return stats, nil
}

// PurgeOld deletes audit logs older than the specified duration
func (s *Service) PurgeOld(ctx context.Context, olderThan time.Duration) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoffTime := time.Now().UTC().Add(-olderThan)

	result := s.db.Where("created_at < ?", cutoffTime).Delete(&models.AuditLog{})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to purge old audit logs: %w", result.Error)
	}

	logger.Info("Purged old audit logs",
		zap.Int64("count", result.RowsAffected),
		zap.Time("cutoff_time", cutoffTime))

	return result.RowsAffected, nil
}

// getIPAddress extracts the real IP address from the request
func getIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

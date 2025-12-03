// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package alertrules

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/alerts"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/metrics"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Service manages alert rules evaluation and execution
type Service struct {
	db             *gorm.DB
	metricsService *metrics.Service
	alertService   *alerts.Service
	mu             sync.RWMutex
	running        bool
	stop           chan bool

	// Track rule states for duration-based alerts
	ruleStates map[uint]*ruleState
}

type ruleState struct {
	conditionMet   bool
	firstMetTime   time.Time
	lastCheckValue float64
}

var (
	globalService *Service
	once          sync.Once
)

// Initialize initializes the alert rules service
func Initialize() (*Service, error) {
	var initErr error
	once.Do(func() {
		db := database.GetDB()
		if db == nil {
			initErr = fmt.Errorf("database not initialized")
			return
		}

		globalService = &Service{
			db:             db,
			metricsService: metrics.GetService(),
			alertService:   alerts.GetService(),
			stop:           make(chan bool),
			ruleStates:     make(map[uint]*ruleState),
		}

		logger.Info("Alert rules service initialized")
	})

	return globalService, initErr
}

// GetService returns the global alert rules service
func GetService() *Service {
	if globalService == nil {
		globalService, _ = Initialize()
	}
	return globalService
}

// Start starts the alert rules evaluation loop
func (s *Service) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("alert rules service already running")
	}

	s.running = true
	go s.run()

	logger.Info("Alert rules evaluation started")
	return nil
}

// Stop stops the alert rules evaluation loop
func (s *Service) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	s.running = false
	s.stop <- true

	logger.Info("Alert rules evaluation stopped")
}

// run is the main evaluation loop
func (s *Service) run() {
	ticker := time.NewTicker(30 * time.Second) // Evaluate every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.evaluateAllRules()
		case <-s.stop:
			return
		}
	}
}

// evaluateAllRules evaluates all enabled alert rules
func (s *Service) evaluateAllRules() {
	ctx := context.Background()

	// Get latest metric
	metric, err := s.metricsService.GetLatestMetric(ctx)
	if err != nil {
		logger.Error("Failed to get latest metric for alert rule evaluation", zap.Error(err))
		return
	}

	// Get all enabled rules
	var rules []models.AlertRule
	if err := s.db.Where("enabled = ?", true).Find(&rules).Error; err != nil {
		logger.Error("Failed to fetch alert rules", zap.Error(err))
		return
	}

	for _, rule := range rules {
		s.evaluateRule(&rule, metric)
	}
}

// evaluateRule evaluates a single alert rule against the current metric
func (s *Service) evaluateRule(rule *models.AlertRule, metric *models.SystemMetric) {
	ctx := context.Background()

	// Extract the relevant metric value
	value, err := s.extractMetricValue(rule.MetricType, metric)
	if err != nil {
		logger.Warn("Failed to extract metric value",
			zap.String("rule", rule.Name),
			zap.String("metric_type", rule.MetricType),
			zap.Error(err))
		return
	}

	// Check if condition is met
	conditionMet := rule.EvaluateCondition(value)

	// Get or create rule state
	s.mu.Lock()
	state, exists := s.ruleStates[rule.ID]
	if !exists {
		state = &ruleState{}
		s.ruleStates[rule.ID] = state
	}
	s.mu.Unlock()

	state.lastCheckValue = value

	// Handle duration-based alerts
	if rule.Duration > 0 {
		if conditionMet {
			if !state.conditionMet {
				// Condition just started being met
				state.conditionMet = true
				state.firstMetTime = time.Now()
			} else {
				// Check if duration threshold is reached
				if time.Since(state.firstMetTime) >= time.Duration(rule.Duration)*time.Second {
					// Duration threshold met, trigger alert
					if rule.ShouldTrigger() {
						s.triggerAlert(ctx, rule, value, metric)
					}
				}
			}
		} else {
			// Condition no longer met, reset state
			if state.conditionMet {
				state.conditionMet = false
				s.resolveAlert(ctx, rule)
			}
		}
	} else {
		// Instant alert (no duration requirement)
		if conditionMet {
			if !state.conditionMet {
				state.conditionMet = true
				if rule.ShouldTrigger() {
					s.triggerAlert(ctx, rule, value, metric)
				}
			}
		} else {
			if state.conditionMet {
				state.conditionMet = false
				s.resolveAlert(ctx, rule)
			}
		}
	}
}

// extractMetricValue extracts the relevant value from a metric based on type
func (s *Service) extractMetricValue(metricType string, metric *models.SystemMetric) (float64, error) {
	switch metricType {
	case models.MetricTypeCPU:
		return metric.CPUUsage, nil
	case models.MetricTypeMemory:
		return metric.MemoryUsage, nil
	case models.MetricTypeDisk:
		return metric.DiskUsage, nil
	case models.MetricTypeTemp:
		return metric.CPUTemperature, nil
	case models.MetricTypeIOPS:
		return float64(metric.DiskIOPS), nil
	case models.MetricTypeNetwork:
		// Use total network throughput as a metric
		return float64(metric.NetworkRxBytesPerSec + metric.NetworkTxBytesPerSec), nil
	case models.MetricTypeHealth:
		// Get latest health score
		healthScore, err := s.metricsService.GetLatestHealthScore(context.Background())
		if err != nil {
			return 0, err
		}
		return float64(healthScore.Score), nil
	default:
		return 0, fmt.Errorf("unknown metric type: %s", metricType)
	}
}

// triggerAlert triggers an alert for a rule
func (s *Service) triggerAlert(ctx context.Context, rule *models.AlertRule, value float64, metric *models.SystemMetric) {
	now := time.Now()

	// Update rule state
	rule.LastTriggered = &now
	rule.TriggerCount++
	rule.IsActive = true
	rule.ActivatedAt = &now

	if err := s.db.Model(rule).Updates(map[string]interface{}{
		"last_triggered": now,
		"trigger_count":  rule.TriggerCount,
		"is_active":      true,
		"activated_at":   now,
	}).Error; err != nil {
		logger.Error("Failed to update rule trigger state", zap.Error(err))
	}

	// Create execution record
	message := fmt.Sprintf("Alert: %s - %s %.2f %s %.2f",
		rule.Name,
		rule.MetricType,
		value,
		conditionToSymbol(rule.Condition),
		rule.Threshold)

	execution := &models.AlertRuleExecution{
		RuleID:            rule.ID,
		MetricValue:       value,
		Threshold:         rule.Threshold,
		Triggered:         true,
		Acknowledged:      false,
		NotificationsSent: false,
		Message:           message,
	}

	if err := s.db.Create(execution).Error; err != nil {
		logger.Error("Failed to create alert execution record", zap.Error(err))
		return
	}

	// Send notifications
	s.sendNotifications(ctx, rule, execution, value, metric)

	logger.Info("Alert rule triggered",
		zap.String("rule", rule.Name),
		zap.Float64("value", value),
		zap.Float64("threshold", rule.Threshold))
}

// resolveAlert resolves an active alert
func (s *Service) resolveAlert(ctx context.Context, rule *models.AlertRule) {
	if !rule.IsActive {
		return
	}

	rule.IsActive = false

	if err := s.db.Model(rule).Update("is_active", false).Error; err != nil {
		logger.Error("Failed to resolve alert", zap.Error(err))
	}

	logger.Info("Alert rule resolved",
		zap.String("rule", rule.Name))
}

// sendNotifications sends notifications for an alert
func (s *Service) sendNotifications(ctx context.Context, rule *models.AlertRule, execution *models.AlertRuleExecution, value float64, metric *models.SystemMetric) {
	// Prepared for future webhook support
	_ = fmt.Sprintf("[%s] %s", severityToEmoji(rule.Severity), rule.Name) // subject

	_ = fmt.Sprintf(`
<html>
<body>
<h2>%s Alert: %s</h2>
<p><strong>%s</strong></p>
<ul>
<li><strong>Metric:</strong> %s</li>
<li><strong>Current Value:</strong> %.2f</li>
<li><strong>Threshold:</strong> %s %.2f</li>
<li><strong>Time:</strong> %s</li>
</ul>
%s
</body>
</html>
`, severityToString(rule.Severity), rule.Name, rule.Description,
		rule.MetricType, value,
		conditionToSymbol(rule.Condition), rule.Threshold,
		time.Now().Format("2006-01-02 15:04:05"),
		generateMetricsSummary(metric)) // htmlBody

	_ = fmt.Sprintf(`**%s Alert: %s**

%s

Metric: %s
Current Value: %.2f
Threshold: %s %.2f
Time: %s`,
		severityToString(rule.Severity), rule.Name, rule.Description,
		rule.MetricType, value,
		conditionToSymbol(rule.Condition), rule.Threshold,
		time.Now().Format("2006-01-02 15:04:05")) // textBody

	var notificationErr error

	// Send email if enabled
	if rule.NotifyEmail {
		config, err := s.alertService.GetConfig(ctx)
		if err == nil && config.Enabled {
			if err := s.alertService.SendCriticalEventAlert(ctx, "Alert Rule", rule.Name, "System", execution.Message); err != nil {
				logger.Error("Failed to send email notification", zap.Error(err))
				notificationErr = err
			}
		}
	}

	// Send webhook if enabled
	if rule.NotifyWebhook {
		config, err := s.alertService.GetConfig(ctx)
		if err == nil && config.WebhookEnabled {
			// Use the webhook sending logic from alerts service
			// This would need to be exposed or we create a custom implementation
			logger.Info("Webhook notification would be sent here")
		}
	}

	// Mark notifications as sent
	if notificationErr == nil {
		execution.NotificationsSent = true
		s.db.Model(execution).Update("notifications_sent", true)
	}
}

// Helper functions

func conditionToSymbol(condition string) string {
	switch condition {
	case models.ConditionGreaterThan:
		return ">"
	case models.ConditionLessThan:
		return "<"
	case models.ConditionEqual:
		return "="
	case models.ConditionGreaterThanOrEqual:
		return ">="
	case models.ConditionLessThanOrEqual:
		return "<="
	default:
		return condition
	}
}

func severityToString(severity string) string {
	switch severity {
	case models.SeverityInfo:
		return "Info"
	case models.SeverityWarning:
		return "Warning"
	case models.SeverityCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

func severityToEmoji(severity string) string {
	switch severity {
	case models.SeverityInfo:
		return "ℹ️"
	case models.SeverityWarning:
		return "⚠️"
	case models.SeverityCritical:
		return "🚨"
	default:
		return "•"
	}
}

func generateMetricsSummary(metric *models.SystemMetric) string {
	return fmt.Sprintf(`<h3>System Status</h3>
<ul>
<li>CPU Usage: %.2f%%</li>
<li>Memory Usage: %.2f%%</li>
<li>Disk Usage: %.2f%%</li>
<li>CPU Temperature: %.2f°C</li>
</ul>`, metric.CPUUsage, metric.MemoryUsage, metric.DiskUsage, metric.CPUTemperature)
}

// CRUD operations for alert rules

// CreateRule creates a new alert rule
func (s *Service) CreateRule(ctx context.Context, rule *models.AlertRule) error {
	return s.db.WithContext(ctx).Create(rule).Error
}

// UpdateRule updates an existing alert rule
func (s *Service) UpdateRule(ctx context.Context, rule *models.AlertRule) error {
	return s.db.WithContext(ctx).Save(rule).Error
}

// DeleteRule deletes an alert rule
func (s *Service) DeleteRule(ctx context.Context, ruleID uint) error {
	// Delete associated executions first
	if err := s.db.WithContext(ctx).Where("rule_id = ?", ruleID).Delete(&models.AlertRuleExecution{}).Error; err != nil {
		return err
	}

	// Delete rule state
	s.mu.Lock()
	delete(s.ruleStates, ruleID)
	s.mu.Unlock()

	return s.db.WithContext(ctx).Delete(&models.AlertRule{}, ruleID).Error
}

// GetRule retrieves a single alert rule by ID
func (s *Service) GetRule(ctx context.Context, ruleID uint) (*models.AlertRule, error) {
	var rule models.AlertRule
	if err := s.db.WithContext(ctx).First(&rule, ruleID).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// ListRules retrieves all alert rules
func (s *Service) ListRules(ctx context.Context) ([]models.AlertRule, error) {
	var rules []models.AlertRule
	if err := s.db.WithContext(ctx).Order("created_at DESC").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// GetExecutions retrieves executions for a rule
func (s *Service) GetExecutions(ctx context.Context, ruleID uint, limit int) ([]models.AlertRuleExecution, error) {
	var executions []models.AlertRuleExecution
	query := s.db.WithContext(ctx).Where("rule_id = ?", ruleID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&executions).Error; err != nil {
		return nil, err
	}
	return executions, nil
}

// GetRecentExecutions retrieves recent alert executions across all rules
func (s *Service) GetRecentExecutions(ctx context.Context, limit int) ([]models.AlertRuleExecution, error) {
	var executions []models.AlertRuleExecution
	query := s.db.WithContext(ctx).Preload("Rule").Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&executions).Error; err != nil {
		return nil, err
	}
	return executions, nil
}

// AcknowledgeExecution acknowledges an alert execution
func (s *Service) AcknowledgeExecution(ctx context.Context, executionID uint, acknowledgedBy, note string) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&models.AlertRuleExecution{}).
		Where("id = ?", executionID).
		Updates(map[string]interface{}{
			"acknowledged":      true,
			"acknowledged_at":   now,
			"acknowledged_by":   acknowledgedBy,
			"acknowledge_note":  note,
		}).Error
}

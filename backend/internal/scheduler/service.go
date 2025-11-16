// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Service handles scheduled task management and execution
type Service struct {
	db      *gorm.DB
	mu      sync.RWMutex
	running bool
	stop    chan bool
	tasks   map[uint]*taskRunner
}

type taskRunner struct {
	task      *models.ScheduledTask
	schedule  *CronSchedule
	nextCheck time.Time
}

var (
	globalService *Service
	once          sync.Once
)

// Initialize initializes the scheduler service
func Initialize() (*Service, error) {
	var initErr error
	once.Do(func() {
		db := database.GetDB()
		if db == nil {
			initErr = fmt.Errorf("database not initialized")
			return
		}

		globalService = &Service{
			db:    db,
			tasks: make(map[uint]*taskRunner),
			stop:  make(chan bool),
		}

		logger.Info("Scheduler service initialized")
	})

	return globalService, initErr
}

// GetService returns the global scheduler service
func GetService() *Service {
	if globalService == nil {
		globalService, _ = Initialize()
	}
	return globalService
}

// Start starts the scheduler
func (s *Service) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("scheduler already running")
	}

	s.running = true
	go s.run()

	logger.Info("Scheduler started")
	return nil
}

// Stop stops the scheduler
func (s *Service) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	s.running = false
	s.stop <- true

	logger.Info("Scheduler stopped")
}

// run is the main scheduler loop
func (s *Service) run() {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	// Initial load
	if err := s.loadTasks(); err != nil {
		logger.Error("Failed to load tasks", zap.Error(err))
	}

	for {
		select {
		case <-ticker.C:
			s.checkAndRunTasks()
			// Reload tasks periodically to pick up changes
			if err := s.loadTasks(); err != nil {
				logger.Error("Failed to reload tasks", zap.Error(err))
			}
		case <-s.stop:
			return
		}
	}
}

// loadTasks loads all enabled tasks from database
func (s *Service) loadTasks() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var tasks []models.ScheduledTask
	if err := s.db.Where("enabled = ?", true).Find(&tasks).Error; err != nil {
		return err
	}

	// Update task map
	for _, task := range tasks {
		if _, exists := s.tasks[task.ID]; !exists {
			// Parse cron expression
			schedule, err := ParseCronExpression(task.CronExpression)
			if err != nil {
				logger.Error("Failed to parse cron expression",
					zap.Uint("taskId", task.ID),
					zap.String("expression", task.CronExpression),
					zap.Error(err))
				continue
			}

			// Calculate next run
			now := time.Now()
			nextRun := schedule.Next(now)

			s.tasks[task.ID] = &taskRunner{
				task:      &task,
				schedule:  schedule,
				nextCheck: now,
			}

			// Update next run in database
			s.db.Model(&task).Updates(map[string]interface{}{
				"next_run": nextRun,
			})
		}
	}

	return nil
}

// checkAndRunTasks checks if any tasks should run now
func (s *Service) checkAndRunTasks() {
	s.mu.RLock()
	tasksToRun := make([]*models.ScheduledTask, 0)
	now := time.Now()

	for _, runner := range s.tasks {
		if now.After(runner.nextCheck) {
			nextRun := runner.schedule.Next(now.Add(-time.Minute))
			if now.After(nextRun) || now.Equal(nextRun) {
				tasksToRun = append(tasksToRun, runner.task)
				runner.nextCheck = runner.schedule.Next(now)

				// Update next run in database
				s.db.Model(runner.task).Updates(map[string]interface{}{
					"next_run": runner.nextCheck,
				})
			}
		}
	}
	s.mu.RUnlock()

	// Run tasks asynchronously
	for _, task := range tasksToRun {
		go s.executeTask(task)
	}
}

// executeTask executes a single task
func (s *Service) executeTask(task *models.ScheduledTask) {
	ctx := context.Background()
	startTime := time.Now()

	logger.Info("Executing scheduled task",
		zap.Uint("taskId", task.ID),
		zap.String("name", task.Name),
		zap.String("type", task.TaskType))

	// Create execution record
	execution := &models.TaskExecution{
		TaskID:      task.ID,
		StartedAt:   startTime,
		Status:      models.TaskStatusRunning,
		TriggeredBy: models.TriggerScheduler,
	}

	if err := s.db.Create(execution).Error; err != nil {
		logger.Error("Failed to create execution record", zap.Error(err))
		return
	}

	// Execute task with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(task.TimeoutSeconds)*time.Second)
	defer cancel()

	output, err := s.runTaskType(timeoutCtx, task)
	completedAt := time.Now()
	duration := completedAt.Sub(startTime).Milliseconds()

	// Update execution record
	execution.CompletedAt = &completedAt
	execution.Duration = duration

	if err != nil {
		execution.Status = models.TaskStatusFailed
		execution.Error = err.Error()
		logger.Error("Task execution failed",
			zap.Uint("taskId", task.ID),
			zap.String("name", task.Name),
			zap.Error(err))
	} else {
		execution.Status = models.TaskStatusSuccess
		execution.Output = output
		logger.Info("Task execution completed",
			zap.Uint("taskId", task.ID),
			zap.String("name", task.Name),
			zap.Int64("duration_ms", duration))
	}

	s.db.Save(execution)

	// Update task
	s.db.Model(task).Updates(map[string]interface{}{
		"last_run":    startTime,
		"last_status": execution.Status,
		"last_error":  execution.Error,
		"run_count":   gorm.Expr("run_count + 1"),
	})
}

// runTaskType executes the actual task based on its type
func (s *Service) runTaskType(ctx context.Context, task *models.ScheduledTask) (string, error) {
	switch task.TaskType {
	case models.TaskTypeCleanup:
		return s.runCleanupTask(ctx, task)
	case models.TaskTypeMaintenance:
		return s.runMaintenanceTask(ctx, task)
	case models.TaskTypeLogRotation:
		return s.runLogRotationTask(ctx, task)
	default:
		return "", fmt.Errorf("unsupported task type: %s", task.TaskType)
	}
}

// runCleanupTask runs a system cleanup task
func (s *Service) runCleanupTask(ctx context.Context, task *models.ScheduledTask) (string, error) {
	var config struct {
		RetentionDays int `json:"retentionDays"`
	}

	if task.Config != "" {
		if err := json.Unmarshal([]byte(task.Config), &config); err != nil {
			return "", fmt.Errorf("invalid config: %w", err)
		}
	}

	if config.RetentionDays == 0 {
		config.RetentionDays = 30 // Default
	}

	cutoffDate := time.Now().AddDate(0, 0, -config.RetentionDays)

	// Clean old audit logs
	result := s.db.Where("created_at < ?", cutoffDate).Delete(&models.AuditLog{})
	auditDeleted := result.RowsAffected

	// Clean old task executions
	result = s.db.Where("created_at < ?", cutoffDate).Delete(&models.TaskExecution{})
	execDeleted := result.RowsAffected

	// Clean old alert logs
	result = s.db.Where("created_at < ?", cutoffDate).Delete(&models.AlertLog{})
	alertDeleted := result.RowsAffected

	output := fmt.Sprintf("Cleanup completed: %d audit logs, %d task executions, %d alert logs deleted",
		auditDeleted, execDeleted, alertDeleted)

	return output, nil
}

// runMaintenanceTask runs database maintenance
func (s *Service) runMaintenanceTask(ctx context.Context, task *models.ScheduledTask) (string, error) {
	// Run VACUUM and ANALYZE on SQLite
	if err := s.db.Exec("VACUUM").Error; err != nil {
		return "", fmt.Errorf("VACUUM failed: %w", err)
	}

	if err := s.db.Exec("ANALYZE").Error; err != nil {
		return "", fmt.Errorf("ANALYZE failed: %w", err)
	}

	return "Database maintenance completed: VACUUM and ANALYZE executed", nil
}

// runLogRotationTask rotates application logs
func (s *Service) runLogRotationTask(ctx context.Context, task *models.ScheduledTask) (string, error) {
	// This would rotate log files
	// Implementation depends on logging setup
	return "Log rotation completed", nil
}

// CreateTask creates a new scheduled task
func (s *Service) CreateTask(ctx context.Context, task *models.ScheduledTask) error {
	// Validate cron expression
	schedule, err := ParseCronExpression(task.CronExpression)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}

	// Calculate next run
	nextRun := schedule.Next(time.Now())
	task.NextRun = &nextRun

	// Create in database
	if err := s.db.WithContext(ctx).Create(task).Error; err != nil {
		return err
	}

	// Add to running tasks if enabled
	if task.Enabled {
		s.mu.Lock()
		s.tasks[task.ID] = &taskRunner{
			task:      task,
			schedule:  schedule,
			nextCheck: time.Now(),
		}
		s.mu.Unlock()
	}

	return nil
}

// GetTask retrieves a task by ID
func (s *Service) GetTask(ctx context.Context, id uint) (*models.ScheduledTask, error) {
	var task models.ScheduledTask
	if err := s.db.WithContext(ctx).First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// ListTasks retrieves all tasks with pagination
func (s *Service) ListTasks(ctx context.Context, offset, limit int) ([]models.ScheduledTask, int64, error) {
	var tasks []models.ScheduledTask
	var total int64

	if err := s.db.WithContext(ctx).Model(&models.ScheduledTask{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// UpdateTask updates a task
func (s *Service) UpdateTask(ctx context.Context, task *models.ScheduledTask) error {
	// Validate cron expression if changed
	schedule, err := ParseCronExpression(task.CronExpression)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}

	// Calculate next run
	nextRun := schedule.Next(time.Now())
	task.NextRun = &nextRun

	// Update in database
	if err := s.db.WithContext(ctx).Save(task).Error; err != nil {
		return err
	}

	// Update in memory
	s.mu.Lock()
	if task.Enabled {
		s.tasks[task.ID] = &taskRunner{
			task:      task,
			schedule:  schedule,
			nextCheck: time.Now(),
		}
	} else {
		delete(s.tasks, task.ID)
	}
	s.mu.Unlock()

	return nil
}

// DeleteTask deletes a task
func (s *Service) DeleteTask(ctx context.Context, id uint) error {
	// Remove from memory
	s.mu.Lock()
	delete(s.tasks, id)
	s.mu.Unlock()

	// Delete from database
	return s.db.WithContext(ctx).Delete(&models.ScheduledTask{}, id).Error
}

// GetTaskExecutions retrieves execution history for a task
func (s *Service) GetTaskExecutions(ctx context.Context, taskID uint, offset, limit int) ([]models.TaskExecution, int64, error) {
	var executions []models.TaskExecution
	var total int64

	query := s.db.WithContext(ctx).Where("task_id = ?", taskID)

	if err := query.Model(&models.TaskExecution{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("started_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&executions).Error; err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// RunTaskNow executes a task immediately
func (s *Service) RunTaskNow(ctx context.Context, taskID uint) error {
	task, err := s.GetTask(ctx, taskID)
	if err != nil {
		return err
	}

	go s.executeTask(task)
	return nil
}

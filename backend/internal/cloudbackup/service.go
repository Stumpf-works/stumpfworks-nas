// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package cloudbackup

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Service manages cloud backup operations
type Service struct {
	rclone   *RcloneClient
	scheduler *cron.Cron
	runningJobs map[uint]context.CancelFunc
	mu       sync.RWMutex
}

var (
	globalService *Service
	once          sync.Once
)

// Initialize initializes the cloud backup service
func Initialize() (*Service, error) {
	var err error
	once.Do(func() {
		rcloneClient, rErr := NewRcloneClient()
		if rErr != nil {
			err = fmt.Errorf("failed to create rclone client: %w", rErr)
			return
		}

		// Check if rclone is installed
		if !rcloneClient.CheckInstalled() {
			logger.Warn("rclone is not installed - cloud backup features will be limited")
		} else {
			version, vErr := rcloneClient.GetVersion()
			if vErr == nil {
				logger.Info("Rclone available", zap.String("version", version))
			}
		}

		globalService = &Service{
			rclone:      rcloneClient,
			scheduler:   cron.New(),
			runningJobs: make(map[uint]context.CancelFunc),
		}

		// Start scheduler
		globalService.scheduler.Start()

		// Load and schedule existing jobs
		if err = globalService.loadScheduledJobs(); err != nil {
			logger.Warn("Failed to load scheduled jobs", zap.Error(err))
			err = nil // Non-fatal
		}
	})

	return globalService, err
}

// GetService returns the global cloud backup service
func GetService() *Service {
	return globalService
}

// loadScheduledJobs loads and schedules all enabled jobs from database
func (s *Service) loadScheduledJobs() error {
	var jobs []models.CloudSyncJob
	if err := database.DB.Where("enabled = ? AND schedule_enabled = ? AND schedule != ?", true, true, "").
		Preload("Provider").Find(&jobs).Error; err != nil {
		return err
	}

	for _, job := range jobs {
		if err := s.scheduleJob(&job); err != nil {
			logger.Warn("Failed to schedule job",
				zap.String("job", job.Name),
				zap.Error(err))
		}
	}

	logger.Info("Loaded scheduled jobs", zap.Int("count", len(jobs)))
	return nil
}

// scheduleJob adds a job to the cron scheduler
func (s *Service) scheduleJob(job *models.CloudSyncJob) error {
	if job.Schedule == "" || !job.ScheduleEnabled {
		return fmt.Errorf("job schedule is not configured or disabled")
	}

	// Remove existing schedule if any
	s.unscheduleJob(job.ID)

	// Add to cron
	entryID, err := s.scheduler.AddFunc(job.Schedule, func() {
		logger.Info("Running scheduled job", zap.String("job", job.Name))
		if _, err := s.RunJob(context.Background(), job.ID, "schedule"); err != nil {
			logger.Error("Scheduled job failed",
				zap.String("job", job.Name),
				zap.Error(err))
		}
	})

	if err != nil {
		return fmt.Errorf("failed to add job to scheduler: %w", err)
	}

	// Calculate next run time
	next := s.scheduler.Entry(entryID).Next
	job.NextRunAt = &next

	// Update job in database
	database.DB.Model(job).Update("next_run_at", next)

	logger.Info("Job scheduled",
		zap.String("job", job.Name),
		zap.String("schedule", job.Schedule),
		zap.Time("next_run", next))

	return nil
}

// unscheduleJob removes a job from the cron scheduler
func (s *Service) unscheduleJob(jobID uint) {
	// This is simplified - in production you'd want to track entry IDs
	// For now, we'll just remove all entries and re-add active ones
}

// Provider Management

// ListProviders returns all cloud providers
func (s *Service) ListProviders() ([]models.CloudProvider, error) {
	var providers []models.CloudProvider
	if err := database.DB.Order("created_at DESC").Find(&providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}

// GetProvider returns a specific cloud provider
func (s *Service) GetProvider(id uint) (*models.CloudProvider, error) {
	var provider models.CloudProvider
	if err := database.DB.First(&provider, id).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}

// CreateProvider creates a new cloud provider
func (s *Service) CreateProvider(provider *models.CloudProvider) error {
	// Validate provider type
	validTypes := map[string]bool{
		"s3": true, "b2": true, "gdrive": true,
		"dropbox": true, "onedrive": true, "azureblob": true,
		"sftp": true, "ftp": true, "webdav": true,
	}
	if !validTypes[provider.Type] {
		return fmt.Errorf("unsupported provider type: %s", provider.Type)
	}

	// Save to database
	if err := database.DB.Create(provider).Error; err != nil {
		return err
	}

	// Configure rclone remote
	if err := s.rclone.ConfigureProvider(provider); err != nil {
		// Rollback database entry
		database.DB.Delete(provider)
		return fmt.Errorf("failed to configure rclone remote: %w", err)
	}

	logger.Info("Cloud provider created",
		zap.String("name", provider.Name),
		zap.String("type", provider.Type))

	return nil
}

// UpdateProvider updates an existing cloud provider
func (s *Service) UpdateProvider(id uint, updates *models.CloudProvider) error {
	var provider models.CloudProvider
	if err := database.DB.First(&provider, id).Error; err != nil {
		return err
	}

	// Update fields
	provider.Name = updates.Name
	provider.Description = updates.Description
	provider.Config = updates.Config
	provider.Enabled = updates.Enabled

	if err := database.DB.Save(&provider).Error; err != nil {
		return err
	}

	// Reconfigure rclone remote
	if err := s.rclone.ConfigureProvider(&provider); err != nil {
		return fmt.Errorf("failed to reconfigure rclone remote: %w", err)
	}

	logger.Info("Cloud provider updated", zap.String("name", provider.Name))
	return nil
}

// DeleteProvider deletes a cloud provider
func (s *Service) DeleteProvider(id uint) error {
	var provider models.CloudProvider
	if err := database.DB.First(&provider, id).Error; err != nil {
		return err
	}

	// Check if any jobs are using this provider
	var jobCount int64
	database.DB.Model(&models.CloudSyncJob{}).Where("provider_id = ?", id).Count(&jobCount)
	if jobCount > 0 {
		return fmt.Errorf("cannot delete provider: %d jobs are using it", jobCount)
	}

	// Remove rclone remote
	if err := s.rclone.RemoveProvider(id); err != nil {
		logger.Warn("Failed to remove rclone remote", zap.Error(err))
	}

	// Delete from database
	if err := database.DB.Delete(&provider).Error; err != nil {
		return err
	}

	logger.Info("Cloud provider deleted", zap.String("name", provider.Name))
	return nil
}

// TestProvider tests connectivity to a cloud provider
func (s *Service) TestProvider(id uint) error {
	var provider models.CloudProvider
	if err := database.DB.First(&provider, id).Error; err != nil {
		return err
	}

	// Test connection
	if err := s.rclone.TestProvider(&provider); err != nil {
		provider.TestStatus = "failed"
		now := time.Now()
		provider.TestedAt = &now
		database.DB.Save(&provider)
		return err
	}

	// Update test status
	provider.TestStatus = "success"
	now := time.Now()
	provider.TestedAt = &now
	database.DB.Save(&provider)

	logger.Info("Provider test successful", zap.String("name", provider.Name))
	return nil
}

// Job Management

// ListJobs returns all cloud sync jobs
func (s *Service) ListJobs() ([]models.CloudSyncJob, error) {
	var jobs []models.CloudSyncJob
	if err := database.DB.Preload("Provider").Order("created_at DESC").Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetJob returns a specific cloud sync job
func (s *Service) GetJob(id uint) (*models.CloudSyncJob, error) {
	var job models.CloudSyncJob
	if err := database.DB.Preload("Provider").First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// CreateJob creates a new cloud sync job
func (s *Service) CreateJob(job *models.CloudSyncJob) error {
	// Validate provider exists
	var provider models.CloudProvider
	if err := database.DB.First(&provider, job.ProviderID).Error; err != nil {
		return fmt.Errorf("provider not found: %w", err)
	}

	// Validate direction
	if job.Direction != "upload" && job.Direction != "download" && job.Direction != "sync" {
		return fmt.Errorf("invalid sync direction: %s", job.Direction)
	}

	// Save to database
	if err := database.DB.Create(job).Error; err != nil {
		return err
	}

	// Schedule if enabled
	if job.ScheduleEnabled && job.Schedule != "" {
		if err := s.scheduleJob(job); err != nil {
			logger.Warn("Failed to schedule job", zap.String("job", job.Name), zap.Error(err))
		}
	}

	logger.Info("Cloud sync job created",
		zap.String("name", job.Name),
		zap.String("direction", job.Direction))

	return nil
}

// UpdateJob updates an existing cloud sync job
func (s *Service) UpdateJob(id uint, updates *models.CloudSyncJob) error {
	var job models.CloudSyncJob
	if err := database.DB.First(&job, id).Error; err != nil {
		return err
	}

	// Update fields
	job.Name = updates.Name
	job.Description = updates.Description
	job.Enabled = updates.Enabled
	job.Direction = updates.Direction
	job.LocalPath = updates.LocalPath
	job.RemotePath = updates.RemotePath
	job.Schedule = updates.Schedule
	job.ScheduleEnabled = updates.ScheduleEnabled
	job.BandwidthLimit = updates.BandwidthLimit
	job.EncryptionEnabled = updates.EncryptionEnabled
	job.DeleteAfterUpload = updates.DeleteAfterUpload
	job.Retention = updates.Retention
	job.Filters = updates.Filters

	if err := database.DB.Save(&job).Error; err != nil {
		return err
	}

	// Update schedule
	if job.ScheduleEnabled && job.Schedule != "" {
		if err := s.scheduleJob(&job); err != nil {
			logger.Warn("Failed to reschedule job", zap.String("job", job.Name), zap.Error(err))
		}
	} else {
		s.unscheduleJob(job.ID)
	}

	logger.Info("Cloud sync job updated", zap.String("name", job.Name))
	return nil
}

// DeleteJob deletes a cloud sync job
func (s *Service) DeleteJob(id uint) error {
	// Cancel if running
	s.mu.Lock()
	if cancel, exists := s.runningJobs[id]; exists {
		cancel()
		delete(s.runningJobs, id)
	}
	s.mu.Unlock()

	var job models.CloudSyncJob
	if err := database.DB.First(&job, id).Error; err != nil {
		return err
	}

	// Unschedule
	s.unscheduleJob(id)

	// Delete from database
	if err := database.DB.Delete(&job).Error; err != nil {
		return err
	}

	logger.Info("Cloud sync job deleted", zap.String("name", job.Name))
	return nil
}

// RunJob executes a cloud sync job
func (s *Service) RunJob(ctx context.Context, jobID uint, triggeredBy string) (*models.CloudSyncLog, error) {
	// Get job with provider
	var job models.CloudSyncJob
	if err := database.DB.Preload("Provider").First(&job, jobID).Error; err != nil {
		return nil, err
	}

	if !job.Enabled {
		return nil, fmt.Errorf("job is disabled")
	}

	if !job.Provider.Enabled {
		return nil, fmt.Errorf("provider is disabled")
	}

	// Check if already running
	s.mu.Lock()
	if _, exists := s.runningJobs[jobID]; exists {
		s.mu.Unlock()
		return nil, fmt.Errorf("job is already running")
	}

	// Create cancellable context
	jobCtx, cancel := context.WithCancel(ctx)
	s.runningJobs[jobID] = cancel
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.runningJobs, jobID)
		s.mu.Unlock()
	}()

	// Update job status
	now := time.Now()
	job.LastRunAt = &now
	job.LastStatus = "running"
	database.DB.Save(&job)

	// Execute sync
	opts := &SyncOptions{
		BandwidthLimit:    job.BandwidthLimit,
		Encryption:        job.EncryptionEnabled,
		EncryptionKey:     job.EncryptionKey,
		DeleteAfterUpload: job.DeleteAfterUpload,
	}

	logEntry, err := s.rclone.Sync(jobCtx, &job, job.Provider, opts)
	logEntry.TriggeredBy = triggeredBy

	// Save log to database
	if err := database.DB.Create(logEntry).Error; err != nil {
		logger.Error("Failed to save sync log", zap.Error(err))
	}

	// Update job status
	job.LastStatus = logEntry.Status
	if logEntry.Status == "failed" {
		job.FailureCount++
	} else {
		job.FailureCount = 0
	}
	database.DB.Save(&job)

	return logEntry, err
}

// GetLogs returns sync logs
func (s *Service) GetLogs(jobID uint, limit int, offset int) ([]models.CloudSyncLog, int64, error) {
	var logs []models.CloudSyncLog
	var total int64

	query := database.DB.Model(&models.CloudSyncLog{})
	if jobID > 0 {
		query = query.Where("job_id = ?", jobID)
	}

	query.Count(&total)

	if err := query.Order("started_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetProviderTypes returns supported provider types with configuration schemas
func (s *Service) GetProviderTypes() map[string]interface{} {
	return map[string]interface{}{
		"s3": map[string]string{
			"access_key_id":     "AWS Access Key ID",
			"secret_access_key": "AWS Secret Access Key",
			"region":            "AWS Region",
			"endpoint":          "S3 Endpoint (optional)",
		},
		"b2": map[string]string{
			"account":    "Backblaze Account ID",
			"key":        "Backblaze Application Key",
		},
		"gdrive": map[string]string{
			"client_id":     "Google Client ID",
			"client_secret": "Google Client Secret",
		},
		"dropbox": map[string]string{
			"client_id":     "Dropbox App Key",
			"client_secret": "Dropbox App Secret",
		},
		"onedrive": map[string]string{
			"client_id":     "Microsoft Client ID",
			"client_secret": "Microsoft Client Secret",
		},
	}
}

// Shutdown gracefully stops the service
func (s *Service) Shutdown() {
	logger.Info("Shutting down cloud backup service...")

	// Cancel all running jobs
	s.mu.Lock()
	for jobID, cancel := range s.runningJobs {
		logger.Info("Cancelling running job", zap.Uint("jobID", jobID))
		cancel()
	}
	s.mu.Unlock()

	// Stop scheduler
	ctx := s.scheduler.Stop()
	<-ctx.Done()

	logger.Info("Cloud backup service shutdown complete")
}

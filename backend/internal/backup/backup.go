package backup

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// BackupJob represents a backup job configuration
type BackupJob struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Source      string            `json:"source"`
	Destination string            `json:"destination"`
	Type        string            `json:"type"` // full, incremental, differential
	Schedule    string            `json:"schedule"` // cron expression
	Enabled     bool              `json:"enabled"`
	Retention   int               `json:"retention"` // days to keep backups
	Compression bool              `json:"compression"`
	Encryption  bool              `json:"encryption"`
	LastRun     *time.Time        `json:"lastRun,omitempty"`
	NextRun     *time.Time        `json:"nextRun,omitempty"`
	Status      string            `json:"status"` // idle, running, success, failed
	Config      map[string]string `json:"config,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

// BackupHistory represents a backup execution record
type BackupHistory struct {
	ID          string    `json:"id"`
	JobID       string    `json:"jobId"`
	JobName     string    `json:"jobName"`
	StartTime   time.Time `json:"startTime"`
	EndTime     *time.Time `json:"endTime,omitempty"`
	Status      string    `json:"status"` // running, success, failed
	BytesBackup int64     `json:"bytesBackup"`
	FilesBackup int       `json:"filesBackup"`
	Duration    int64     `json:"duration"` // seconds
	Error       string    `json:"error,omitempty"`
	BackupPath  string    `json:"backupPath"`
}

// Snapshot represents a filesystem snapshot
type Snapshot struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Filesystem  string    `json:"filesystem"`
	CreatedAt   time.Time `json:"createdAt"`
	Size        int64     `json:"size"`
	Used        int64     `json:"used"`
	Referenced  int64     `json:"referenced"`
	Type        string    `json:"type"` // zfs, btrfs, lvm
	Description string    `json:"description,omitempty"`
}

// Service handles backup operations
type Service struct {
	backupDir  string
	jobs       map[string]*BackupJob
	history    []*BackupHistory
	snapshots  []*Snapshot
	mu         sync.RWMutex
	available  bool
}

var (
	globalService *Service
	once          sync.Once
)

const (
	DefaultBackupDir = "/var/lib/stumpfworks/backups"
)

// Initialize initializes the backup service
func Initialize(backupDir string) (*Service, error) {
	var err error
	once.Do(func() {
		if backupDir == "" {
			backupDir = DefaultBackupDir
		}

		// Ensure backup directory exists
		if err = os.MkdirAll(backupDir, 0755); err != nil {
			return
		}

		globalService = &Service{
			backupDir:  backupDir,
			jobs:       make(map[string]*BackupJob),
			history:    make([]*BackupHistory, 0),
			snapshots:  make([]*Snapshot, 0),
			available:  true,
		}

		// Discover existing snapshots
		if err = globalService.discoverSnapshots(); err != nil {
			// Log but don't fail - snapshots might not be available
			err = nil
		}
	})

	return globalService, err
}

// GetService returns the global backup service
func GetService() *Service {
	return globalService
}

// discoverSnapshots discovers existing ZFS/Btrfs snapshots
func (s *Service) discoverSnapshots() error {
	// Try ZFS snapshots first
	if err := s.discoverZFSSnapshots(); err == nil {
		return nil
	}

	// Try Btrfs snapshots
	if err := s.discoverBtrfsSnapshots(); err == nil {
		return nil
	}

	return fmt.Errorf("no snapshot system available")
}

// discoverZFSSnapshots discovers ZFS snapshots
func (s *Service) discoverZFSSnapshots() error {
	cmd := exec.Command("zfs", "list", "-t", "snapshot", "-H", "-o", "name,used,refer,creation")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	// Parse ZFS snapshot output
	// Format: name used referenced creation
	// This is a simplified implementation
	_ = output

	return nil
}

// discoverBtrfsSnapshots discovers Btrfs snapshots
func (s *Service) discoverBtrfsSnapshots() error {
	// Btrfs snapshot discovery
	// This would require parsing btrfs subvolume list
	return fmt.Errorf("btrfs not implemented yet")
}

// ListJobs returns all backup jobs
func (s *Service) ListJobs(ctx context.Context) ([]*BackupJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make([]*BackupJob, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GetJob returns a specific backup job
func (s *Service) GetJob(ctx context.Context, id string) (*BackupJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, ok := s.jobs[id]
	if !ok {
		return nil, fmt.Errorf("backup job not found: %s", id)
	}

	return job, nil
}

// CreateJob creates a new backup job
func (s *Service) CreateJob(ctx context.Context, job *BackupJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate job
	if job.Name == "" {
		return fmt.Errorf("job name is required")
	}
	if job.Source == "" {
		return fmt.Errorf("source path is required")
	}
	if job.Destination == "" {
		return fmt.Errorf("destination path is required")
	}

	// Generate ID if not provided
	if job.ID == "" {
		job.ID = fmt.Sprintf("backup-%d", time.Now().Unix())
	}

	// Check if job already exists
	if _, exists := s.jobs[job.ID]; exists {
		return fmt.Errorf("backup job already exists: %s", job.ID)
	}

	// Set timestamps
	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now
	job.Status = "idle"

	s.jobs[job.ID] = job

	return nil
}

// UpdateJob updates an existing backup job
func (s *Service) UpdateJob(ctx context.Context, id string, updates *BackupJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		return fmt.Errorf("backup job not found: %s", id)
	}

	// Update fields
	if updates.Name != "" {
		job.Name = updates.Name
	}
	if updates.Description != "" {
		job.Description = updates.Description
	}
	if updates.Source != "" {
		job.Source = updates.Source
	}
	if updates.Destination != "" {
		job.Destination = updates.Destination
	}
	if updates.Type != "" {
		job.Type = updates.Type
	}
	if updates.Schedule != "" {
		job.Schedule = updates.Schedule
	}
	if updates.Retention > 0 {
		job.Retention = updates.Retention
	}

	job.Compression = updates.Compression
	job.Encryption = updates.Encryption
	job.UpdatedAt = time.Now()

	return nil
}

// DeleteJob deletes a backup job
func (s *Service) DeleteJob(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.jobs[id]; !ok {
		return fmt.Errorf("backup job not found: %s", id)
	}

	delete(s.jobs, id)
	return nil
}

// RunJob executes a backup job
func (s *Service) RunJob(ctx context.Context, id string) (*BackupHistory, error) {
	s.mu.Lock()
	job, ok := s.jobs[id]
	if !ok {
		s.mu.Unlock()
		return nil, fmt.Errorf("backup job not found: %s", id)
	}

	// Check if job is already running
	if job.Status == "running" {
		s.mu.Unlock()
		return nil, fmt.Errorf("backup job is already running")
	}

	// Update job status
	job.Status = "running"
	now := time.Now()
	job.LastRun = &now
	s.mu.Unlock()

	// Create history entry
	history := &BackupHistory{
		ID:        fmt.Sprintf("history-%d", time.Now().UnixNano()),
		JobID:     job.ID,
		JobName:   job.Name,
		StartTime: now,
		Status:    "running",
	}

	// Execute backup
	err := s.executeBackup(ctx, job, history)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Update job and history status
	endTime := time.Now()
	history.EndTime = &endTime
	history.Duration = int64(endTime.Sub(history.StartTime).Seconds())

	if err != nil {
		job.Status = "failed"
		history.Status = "failed"
		history.Error = err.Error()
	} else {
		job.Status = "success"
		history.Status = "success"
	}

	job.UpdatedAt = time.Now()
	s.history = append(s.history, history)

	return history, err
}

// executeBackup performs the actual backup operation
func (s *Service) executeBackup(ctx context.Context, job *BackupJob, history *BackupHistory) error {
	// Create backup destination directory
	timestamp := time.Now().Format("20060102-150405")
	backupPath := filepath.Join(s.backupDir, job.ID, timestamp)

	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	history.BackupPath = backupPath

	// Build rsync command for backup
	args := []string{"-av"}

	if job.Compression {
		args = append(args, "-z")
	}

	args = append(args, job.Source, backupPath+"/")

	cmd := exec.CommandContext(ctx, "rsync", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("backup failed: %w, output: %s", err, string(output))
	}

	// Get backup size
	if info, err := os.Stat(backupPath); err == nil {
		history.BytesBackup = info.Size()
	}

	return nil
}

// GetHistory returns backup history
func (s *Service) GetHistory(ctx context.Context, jobID string, limit int) ([]*BackupHistory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*BackupHistory

	for i := len(s.history) - 1; i >= 0 && len(result) < limit; i-- {
		h := s.history[i]
		if jobID == "" || h.JobID == jobID {
			result = append(result, h)
		}
	}

	return result, nil
}

// ListSnapshots returns all snapshots
func (s *Service) ListSnapshots(ctx context.Context) ([]*Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.snapshots, nil
}

// CreateSnapshot creates a new snapshot
func (s *Service) CreateSnapshot(ctx context.Context, filesystem string, name string) (*Snapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Try ZFS snapshot
	snapshotName := fmt.Sprintf("%s@%s", filesystem, name)
	cmd := exec.CommandContext(ctx, "zfs", "snapshot", snapshotName)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}

	snapshot := &Snapshot{
		ID:         snapshotName,
		Name:       name,
		Filesystem: filesystem,
		CreatedAt:  time.Now(),
		Type:       "zfs",
	}

	s.snapshots = append(s.snapshots, snapshot)

	return snapshot, nil
}

// DeleteSnapshot deletes a snapshot
func (s *Service) DeleteSnapshot(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Try ZFS snapshot deletion
	cmd := exec.CommandContext(ctx, "zfs", "destroy", id)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}

	// Remove from list
	for i, snap := range s.snapshots {
		if snap.ID == id {
			s.snapshots = append(s.snapshots[:i], s.snapshots[i+1:]...)
			break
		}
	}

	return nil
}

// RestoreSnapshot restores from a snapshot
func (s *Service) RestoreSnapshot(ctx context.Context, id string, destination string) error {
	// ZFS rollback or restore
	cmd := exec.CommandContext(ctx, "zfs", "rollback", id)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restore snapshot: %w", err)
	}

	return nil
}

// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package cloudbackup

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

const (
	rcloneConfigDir = "/etc/stumpfworks/rclone"
	rcloneConfigFile = "rclone.conf"
)

// RcloneClient wraps rclone operations
type RcloneClient struct {
	configPath string
}

// NewRcloneClient creates a new rclone client
func NewRcloneClient() (*RcloneClient, error) {
	// Ensure config directory exists
	if err := os.MkdirAll(rcloneConfigDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create rclone config directory: %w", err)
	}

	configPath := filepath.Join(rcloneConfigDir, rcloneConfigFile)

	return &RcloneClient{
		configPath: configPath,
	}, nil
}

// ConfigureProvider creates or updates an rclone remote configuration
func (r *RcloneClient) ConfigureProvider(provider *models.CloudProvider) error {
	// Parse provider config
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(provider.Config), &config); err != nil {
		return fmt.Errorf("failed to parse provider config: %w", err)
	}

	// Build rclone config command
	remoteName := fmt.Sprintf("provider_%d", provider.ID)

	// Create config using rclone config create
	args := []string{
		"config", "create",
		remoteName,
		provider.Type,
		"--config", r.configPath,
	}

	// Add provider-specific parameters
	for key, value := range config {
		args = append(args, fmt.Sprintf("%s=%v", key, value))
	}

	cmd := exec.Command("rclone", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to configure rclone remote: %w, output: %s", err, string(output))
	}

	logger.Info("Rclone remote configured",
		zap.String("provider", provider.Name),
		zap.String("type", provider.Type))

	return nil
}

// TestProvider tests connectivity to a cloud provider
func (r *RcloneClient) TestProvider(provider *models.CloudProvider) error {
	remoteName := fmt.Sprintf("provider_%d", provider.ID)

	// Use rclone lsd to test connectivity (list directories at root)
	cmd := exec.Command("rclone", "lsd", remoteName+":", "--config", r.configPath, "--max-depth", "1")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd = exec.CommandContext(ctx, "rclone", "lsd", remoteName+":", "--config", r.configPath, "--max-depth", "1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("connection test failed: %w, output: %s", err, string(output))
	}

	return nil
}

// RemoveProvider removes an rclone remote configuration
func (r *RcloneClient) RemoveProvider(providerID uint) error {
	remoteName := fmt.Sprintf("provider_%d", providerID)

	cmd := exec.Command("rclone", "config", "delete", remoteName, "--config", r.configPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove rclone remote: %w, output: %s", err, string(output))
	}

	return nil
}

// SyncJobProgress represents sync progress information
type SyncJobProgress struct {
	BytesTransferred int64
	FilesTransferred int
	TotalBytes       int64
	TotalFiles       int
	Speed            float64 // bytes per second
	ETA              int64   // seconds
	Errors           []string
}

// SyncOptions contains options for sync operations
type SyncOptions struct {
	BandwidthLimit    string // e.g., "10M", "1G"
	Encryption        bool
	EncryptionKey     string
	DeleteAfterUpload bool
	DryRun            bool
	Filters           []string // include/exclude patterns
	ProgressCallback  func(*SyncJobProgress)
}

// Sync synchronizes data between local and remote
func (r *RcloneClient) Sync(ctx context.Context, job *models.CloudSyncJob, provider *models.CloudProvider, opts *SyncOptions) (*models.CloudSyncLog, error) {
	remoteName := fmt.Sprintf("provider_%d", provider.ID)

	// Create log entry
	logEntry := &models.CloudSyncLog{
		JobID:      job.ID,
		JobName:    job.Name,
		StartedAt:  time.Now(),
		Status:     "running",
		Direction:  job.Direction,
		TriggeredBy: "manual",
	}

	// Build rclone command based on direction
	var args []string
	source := ""
	dest := ""

	switch job.Direction {
	case "upload":
		source = job.LocalPath
		dest = fmt.Sprintf("%s:%s", remoteName, job.RemotePath)
		args = append(args, "sync")
	case "download":
		source = fmt.Sprintf("%s:%s", remoteName, job.RemotePath)
		dest = job.LocalPath
		args = append(args, "sync")
	case "sync":
		source = job.LocalPath
		dest = fmt.Sprintf("%s:%s", remoteName, job.RemotePath)
		args = append(args, "sync")
	default:
		return logEntry, fmt.Errorf("unknown sync direction: %s", job.Direction)
	}

	args = append(args, source, dest)
	args = append(args, "--config", r.configPath)
	args = append(args, "--progress")
	args = append(args, "--stats", "1s")
	args = append(args, "--stats-one-line")

	// Add bandwidth limit
	if opts.BandwidthLimit != "" {
		args = append(args, "--bwlimit", opts.BandwidthLimit)
	} else if job.BandwidthLimit != "" {
		args = append(args, "--bwlimit", job.BandwidthLimit)
	}

	// Add delete flag if requested
	if opts.DeleteAfterUpload || job.DeleteAfterUpload {
		args = append(args, "--delete-after")
	}

	// Add dry run for testing
	if opts.DryRun {
		args = append(args, "--dry-run")
	}

	// Add filters
	filters := opts.Filters
	if len(filters) == 0 && job.Filters != "" {
		if err := json.Unmarshal([]byte(job.Filters), &filters); err == nil {
			for _, filter := range filters {
				if strings.HasPrefix(filter, "+") || strings.HasPrefix(filter, "-") {
					args = append(args, "--filter", filter)
				}
			}
		}
	}

	logger.Info("Starting rclone sync",
		zap.String("job", job.Name),
		zap.String("direction", job.Direction),
		zap.String("source", source),
		zap.String("dest", dest))

	// Execute rclone command
	cmd := exec.CommandContext(ctx, "rclone", args...)

	// Capture output for parsing
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logEntry.Status = "failed"
		logEntry.ErrorMessage = fmt.Sprintf("failed to create stdout pipe: %v", err)
		return logEntry, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logEntry.Status = "failed"
		logEntry.ErrorMessage = fmt.Sprintf("failed to create stderr pipe: %v", err)
		return logEntry, err
	}

	if err := cmd.Start(); err != nil {
		logEntry.Status = "failed"
		logEntry.ErrorMessage = fmt.Sprintf("failed to start rclone: %v", err)
		return logEntry, err
	}

	// Parse progress output
	progress := &SyncJobProgress{}
	outputLines := []string{}
	errors := []string{}

	// Read stdout for progress
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			outputLines = append(outputLines, line)

			// Parse progress line
			r.parseProgressLine(line, progress)

			// Call progress callback if provided
			if opts.ProgressCallback != nil {
				opts.ProgressCallback(progress)
			}
		}
	}()

	// Read stderr for errors
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			outputLines = append(outputLines, line)
			if strings.Contains(strings.ToLower(line), "error") ||
			   strings.Contains(strings.ToLower(line), "failed") {
				errors = append(errors, line)
			}
		}
	}()

	// Wait for command to complete
	err = cmd.Wait()

	// Update log entry
	completedAt := time.Now()
	logEntry.CompletedAt = &completedAt
	logEntry.Duration = int64(completedAt.Sub(logEntry.StartedAt).Seconds())
	logEntry.BytesTransferred = progress.BytesTransferred
	logEntry.FilesTransferred = progress.FilesTransferred
	logEntry.Output = strings.Join(outputLines, "\n")

	if err != nil {
		logEntry.Status = "failed"
		logEntry.ErrorMessage = fmt.Sprintf("rclone sync failed: %v", err)
		if len(errors) > 0 {
			logEntry.ErrorMessage += "\n" + strings.Join(errors, "\n")
		}
		logger.Error("Rclone sync failed",
			zap.String("job", job.Name),
			zap.Error(err))
	} else {
		logEntry.Status = "success"
		logger.Info("Rclone sync completed successfully",
			zap.String("job", job.Name),
			zap.Int64("bytes", progress.BytesTransferred),
			zap.Int("files", progress.FilesTransferred))
	}

	return logEntry, err
}

// parseProgressLine parses rclone progress output
func (r *RcloneClient) parseProgressLine(line string, progress *SyncJobProgress) {
	// Example rclone output line:
	// "Transferred:   	  100 MiB / 100 MiB, 100%, 10 MiB/s, ETA 0s"

	if strings.Contains(line, "Transferred:") {
		// Parse transferred bytes
		parts := strings.Split(line, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)

			// Parse bytes transferred
			if strings.Contains(part, "MiB") || strings.Contains(part, "GiB") || strings.Contains(part, "KiB") {
				// Simple parsing - just extract the number before the first slash
				if strings.Contains(part, "/") {
					transferred := strings.TrimSpace(strings.Split(part, "/")[0])
					// Convert to bytes (simplified, would need proper unit conversion)
					progress.BytesTransferred = parseSizeString(transferred)
				}
			}

			// Parse speed
			if strings.Contains(part, "/s") {
				speedStr := strings.TrimSpace(strings.Split(part, "/s")[0])
				progress.Speed = float64(parseSizeString(speedStr))
			}
		}
	}

	// Check for file count updates
	if strings.Contains(line, "Checks:") || strings.Contains(line, "Transferred:") {
		// Parse file counts from status lines
		// This is simplified - full implementation would parse all stats
	}
}

// parseSizeString converts size strings like "100 MiB" to bytes
func parseSizeString(sizeStr string) int64 {
	sizeStr = strings.TrimSpace(sizeStr)
	var multiplier int64 = 1

	if strings.HasSuffix(sizeStr, "GiB") {
		multiplier = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "GiB")
	} else if strings.HasSuffix(sizeStr, "MiB") {
		multiplier = 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "MiB")
	} else if strings.HasSuffix(sizeStr, "KiB") {
		multiplier = 1024
		sizeStr = strings.TrimSuffix(sizeStr, "KiB")
	} else if strings.HasSuffix(sizeStr, "GB") {
		multiplier = 1000 * 1000 * 1000
		sizeStr = strings.TrimSuffix(sizeStr, "GB")
	} else if strings.HasSuffix(sizeStr, "MB") {
		multiplier = 1000 * 1000
		sizeStr = strings.TrimSuffix(sizeStr, "MB")
	} else if strings.HasSuffix(sizeStr, "KB") {
		multiplier = 1000
		sizeStr = strings.TrimSuffix(sizeStr, "KB")
	}

	sizeStr = strings.TrimSpace(sizeStr)
	var value float64
	fmt.Sscanf(sizeStr, "%f", &value)

	return int64(value * float64(multiplier))
}

// GetVersion returns the installed rclone version
func (r *RcloneClient) GetVersion() (string, error) {
	cmd := exec.Command("rclone", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get rclone version: %w", err)
	}

	// Parse version from output
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}

	return "", fmt.Errorf("unable to parse rclone version")
}

// CheckInstalled checks if rclone is installed
func (r *RcloneClient) CheckInstalled() bool {
	_, err := exec.LookPath("rclone")
	return err == nil
}

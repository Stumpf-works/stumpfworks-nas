package system

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/stumpfworks/nas/pkg/logger"
	"go.uber.org/zap"
)

// UpdateInfo represents update information
type UpdateInfo struct {
	Available       bool      `json:"available"`
	CurrentVersion  string    `json:"currentVersion"`
	CurrentCommit   string    `json:"currentCommit"`
	LatestCommit    string    `json:"latestCommit"`
	BehindBy        int       `json:"behindBy"`
	LastChecked     time.Time `json:"lastChecked"`
	UpdateAvailable bool      `json:"updateAvailable"`
	ChangeLog       []string  `json:"changeLog"`
}

var lastUpdateCheck *UpdateInfo

// CheckForUpdates checks if updates are available from git remote
func CheckForUpdates() (*UpdateInfo, error) {
	info := &UpdateInfo{
		CurrentVersion: "0.1.0-alpha",
		LastChecked:    time.Now(),
	}

	// Get current commit
	currentCommit, err := getCurrentCommit()
	if err != nil {
		logger.Warn("Failed to get current commit", zap.Error(err))
		info.CurrentCommit = "unknown"
	} else {
		info.CurrentCommit = currentCommit
	}

	// Fetch latest from remote
	if err := fetchRemote(); err != nil {
		logger.Warn("Failed to fetch from remote", zap.Error(err))
		return info, err
	}

	// Get latest commit from remote
	latestCommit, err := getLatestRemoteCommit()
	if err != nil {
		logger.Warn("Failed to get latest remote commit", zap.Error(err))
		return info, err
	}
	info.LatestCommit = latestCommit

	// Check if we're behind
	behindBy, err := getCommitsBehind()
	if err != nil {
		logger.Warn("Failed to check commits behind", zap.Error(err))
		behindBy = 0
	}
	info.BehindBy = behindBy
	info.Available = behindBy > 0
	info.UpdateAvailable = behindBy > 0

	// Get changelog if updates available
	if info.Available {
		changelog, err := getChangeLog(currentCommit, latestCommit)
		if err == nil {
			info.ChangeLog = changelog
		}
	}

	lastUpdateCheck = info
	logger.Info("Update check completed",
		zap.Bool("available", info.Available),
		zap.Int("behindBy", info.BehindBy))

	return info, nil
}

// GetLastUpdateCheck returns the last update check result
func GetLastUpdateCheck() *UpdateInfo {
	return lastUpdateCheck
}

// PerformUpdate performs git pull to update the system
func PerformUpdate() error {
	logger.Info("Performing system update...")

	// Git pull
	cmd := exec.Command("git", "pull", "origin", getCurrentBranch())
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Git pull failed", zap.Error(err), zap.String("output", string(output)))
		return fmt.Errorf("git pull failed: %w", err)
	}

	logger.Info("System updated successfully", zap.String("output", string(output)))
	return nil
}

// Helper functions

func getCurrentCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getCurrentBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "main"
	}
	return strings.TrimSpace(string(output))
}

func fetchRemote() error {
	cmd := exec.Command("git", "fetch", "origin")
	return cmd.Run()
}

func getLatestRemoteCommit() (string, error) {
	branch := getCurrentBranch()
	cmd := exec.Command("git", "rev-parse", "--short", fmt.Sprintf("origin/%s", branch))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getCommitsBehind() (int, error) {
	branch := getCurrentBranch()
	cmd := exec.Command("git", "rev-list", "--count", fmt.Sprintf("HEAD..origin/%s", branch))
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var count int
	_, err = fmt.Sscanf(string(output), "%d", &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func getChangeLog(from, to string) ([]string, error) {
	cmd := exec.Command("git", "log", "--oneline", fmt.Sprintf("%s..%s", from, to))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	changelog := make([]string, 0, len(lines))
	for _, line := range lines {
		if line != "" {
			changelog = append(changelog, line)
		}
	}

	return changelog, nil
}

// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package updates

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

const (
	GitHubAPIURL      = "https://api.github.com/repos/%s/releases/latest"
	GitHubReleaseURL  = "https://github.com/%s/releases/tag/%s"
	DefaultRepository = "Stumpf-works/stumpfworks-nas"
	CurrentVersion    = "v1.3.0"
)

// UpdateService handles update checking and management
type UpdateService struct {
	currentVersion string
	repository     string
	client         *http.Client
	mu             sync.RWMutex
	lastCheck      time.Time
	cachedRelease  *ReleaseInfo
}

// ReleaseInfo represents GitHub release information
type ReleaseInfo struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	HTMLURL     string    `json:"html_url"`
	Assets      []Asset   `json:"assets"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
}

// Asset represents a release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// UpdateCheckResult represents the result of an update check
type UpdateCheckResult struct {
	UpdateAvailable bool         `json:"updateAvailable"`
	CurrentVersion  string       `json:"currentVersion"`
	LatestVersion   string       `json:"latestVersion"`
	ReleaseInfo     *ReleaseInfo `json:"releaseInfo,omitempty"`
	Message         string       `json:"message"`
}

var (
	globalService *UpdateService
	once          sync.Once
)

// Initialize initializes the update service
func Initialize() (*UpdateService, error) {
	once.Do(func() {
		globalService = &UpdateService{
			currentVersion: CurrentVersion,
			repository:     DefaultRepository,
			client: &http.Client{
				Timeout: 30 * time.Second,
			},
		}
	})

	return globalService, nil
}

// GetService returns the global update service
func GetService() *UpdateService {
	if globalService == nil {
		globalService, _ = Initialize()
	}
	return globalService
}

// CheckForUpdates checks GitHub for new releases
func (s *UpdateService) CheckForUpdates(ctx context.Context, forceCheck bool) (*UpdateCheckResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Use cached result if available and recent (< 1 hour)
	if !forceCheck && s.cachedRelease != nil && time.Since(s.lastCheck) < time.Hour {
		return s.buildResult(s.cachedRelease), nil
	}

	// Fetch latest release from GitHub
	apiURL := fmt.Sprintf(GitHubAPIURL, s.repository)
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add GitHub API headers
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Stumpfworks-NAS-Update-Checker")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	// Handle 404 gracefully - no releases available
	if resp.StatusCode == http.StatusNotFound {
		logger.Info("No releases found on GitHub",
			zap.String("repository", s.repository))
		return &UpdateCheckResult{
			UpdateAvailable: false,
			CurrentVersion:  s.currentVersion,
			LatestVersion:   s.currentVersion,
			Message:         "No releases available on GitHub yet",
		}, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	// Skip drafts and prereleases
	if release.Draft || release.Prerelease {
		logger.Info("Latest release is draft or prerelease, skipping",
			zap.String("version", release.TagName))
	}

	// Update cache
	s.cachedRelease = &release
	s.lastCheck = time.Now()

	logger.Info("Update check completed",
		zap.String("current", s.currentVersion),
		zap.String("latest", release.TagName))

	return s.buildResult(&release), nil
}

// buildResult builds an UpdateCheckResult from a release
func (s *UpdateService) buildResult(release *ReleaseInfo) *UpdateCheckResult {
	updateAvailable := s.isNewerVersion(release.TagName)

	message := "You are running the latest version"
	if updateAvailable {
		message = fmt.Sprintf("Update available: %s → %s", s.currentVersion, release.TagName)
	}

	return &UpdateCheckResult{
		UpdateAvailable: updateAvailable,
		CurrentVersion:  s.currentVersion,
		LatestVersion:   release.TagName,
		ReleaseInfo:     release,
		Message:         message,
	}
}

// isNewerVersion checks if the release version is newer than current
func (s *UpdateService) isNewerVersion(releaseVersion string) bool {
	// Simple version comparison (assumes semantic versioning vX.Y.Z)
	current := strings.TrimPrefix(s.currentVersion, "v")
	release := strings.TrimPrefix(releaseVersion, "v")

	// Split into parts
	currentParts := strings.Split(current, ".")
	releaseParts := strings.Split(release, ".")

	// Compare each part
	for i := 0; i < len(currentParts) && i < len(releaseParts); i++ {
		if releaseParts[i] > currentParts[i] {
			return true
		}
		if releaseParts[i] < currentParts[i] {
			return false
		}
	}

	// If all parts are equal, check length
	return len(releaseParts) > len(currentParts)
}

// GetCurrentVersion returns the current version
func (s *UpdateService) GetCurrentVersion() string {
	return s.currentVersion
}

// GetReleaseURL returns the URL to the release page
func (s *UpdateService) GetReleaseURL(version string) string {
	return fmt.Sprintf(GitHubReleaseURL, s.repository, version)
}

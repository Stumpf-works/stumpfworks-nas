package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
)

const (
	// Default plugin registry URL
	DefaultRegistryURL = "https://raw.githubusercontent.com/Stumpf-works/stumpfworks-nas-apps/main/registry.json"

	// Cache duration
	RegistryCacheDuration = 1 * time.Hour
)

// RegistryService manages the plugin registry
type RegistryService struct {
	db          *gorm.DB
	registryURL string
	lastSync    time.Time
	httpClient  *http.Client
}

// RegistryManifest represents the registry.json structure
type RegistryManifest struct {
	Version string                   `json:"version"`
	Updated time.Time                `json:"updated"`
	Plugins []RegistryPluginMetadata `json:"plugins"`
}

// RegistryPluginMetadata represents plugin metadata in the registry
type RegistryPluginMetadata struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Version       string   `json:"version"`
	Author        string   `json:"author"`
	Description   string   `json:"description"`
	Icon          string   `json:"icon"`
	Category      string   `json:"category"`
	RepositoryURL string   `json:"repository_url"`
	DownloadURL   string   `json:"download_url"`
	Homepage      string   `json:"homepage"`
	MinNasVersion string   `json:"min_nas_version"`
	RequireDocker bool     `json:"require_docker"`
	RequiredPorts []int    `json:"required_ports"`
	Screenshots   []string `json:"screenshots,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

// NewRegistryService creates a new registry service
func NewRegistryService(db *gorm.DB, registryURL string) *RegistryService {
	if registryURL == "" {
		registryURL = DefaultRegistryURL
	}

	return &RegistryService{
		db:          db,
		registryURL: registryURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Sync syncs the registry from the remote source
func (s *RegistryService) Sync() error {
	// Check cache
	if time.Since(s.lastSync) < RegistryCacheDuration {
		log.Debug().Msg("Registry cache is still valid, skipping sync")
		return nil
	}

	log.Info().Str("url", s.registryURL).Msg("Syncing plugin registry")

	// Fetch registry
	resp, err := s.httpClient.Get(s.registryURL)
	if err != nil {
		return fmt.Errorf("failed to fetch registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registry returned status %d", resp.StatusCode)
	}

	// Parse manifest
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read registry: %w", err)
	}

	var manifest RegistryManifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return fmt.Errorf("failed to parse registry: %w", err)
	}

	// Update database
	if err := s.updateDatabase(manifest); err != nil {
		return fmt.Errorf("failed to update database: %w", err)
	}

	s.lastSync = time.Now()
	log.Info().Int("plugins", len(manifest.Plugins)).Msg("Registry synced successfully")

	return nil
}

// updateDatabase updates the database with registry data
func (s *RegistryService) updateDatabase(manifest RegistryManifest) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Clear existing registry (or do upsert)
	if err := tx.Exec("DELETE FROM plugin_registry").Error; err != nil {
		tx.Rollback()
		return err
	}

	// Insert plugins
	for _, p := range manifest.Plugins {
		plugin := models.PluginRegistry{
			ID:            p.ID,
			Name:          p.Name,
			Version:       p.Version,
			Author:        p.Author,
			Description:   p.Description,
			Icon:          p.Icon,
			Category:      p.Category,
			RepositoryURL: p.RepositoryURL,
			DownloadURL:   p.DownloadURL,
			Homepage:      p.Homepage,
			MinNasVersion: p.MinNasVersion,
			RequireDocker: p.RequireDocker,
			RequiredPorts: p.RequiredPorts,
			LastUpdated:   manifest.Updated,
		}

		if err := tx.Create(&plugin).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// List returns all plugins from registry
func (s *RegistryService) List() ([]models.PluginRegistry, error) {
	// Sync if needed
	if err := s.Sync(); err != nil {
		log.Warn().Err(err).Msg("Failed to sync registry, using cached data")
	}

	var plugins []models.PluginRegistry
	if err := s.db.Find(&plugins).Error; err != nil {
		return nil, err
	}

	// Enrich with installation status
	var installedPlugins []models.InstalledPlugin
	s.db.Find(&installedPlugins)

	installedMap := make(map[string]models.InstalledPlugin)
	for _, ip := range installedPlugins {
		installedMap[ip.ID] = ip
	}

	for i := range plugins {
		if ip, ok := installedMap[plugins[i].ID]; ok {
			plugins[i].Installed = true
			plugins[i].InstalledVersion = ip.Version
		}
	}

	return plugins, nil
}

// Get returns a specific plugin from registry
func (s *RegistryService) Get(id string) (*models.PluginRegistry, error) {
	// Sync if needed
	if err := s.Sync(); err != nil {
		log.Warn().Err(err).Msg("Failed to sync registry, using cached data")
	}

	var plugin models.PluginRegistry
	if err := s.db.Where("id = ?", id).First(&plugin).Error; err != nil {
		return nil, err
	}

	// Check if installed
	var installed models.InstalledPlugin
	if err := s.db.Where("id = ?", id).First(&installed).Error; err == nil {
		plugin.Installed = true
		plugin.InstalledVersion = installed.Version
	}

	return &plugin, nil
}

// Search searches plugins by query
func (s *RegistryService) Search(query string) ([]models.PluginRegistry, error) {
	// Sync if needed
	if err := s.Sync(); err != nil {
		log.Warn().Err(err).Msg("Failed to sync registry, using cached data")
	}

	var plugins []models.PluginRegistry
	searchPattern := "%" + query + "%"

	if err := s.db.Where(
		"name LIKE ? OR description LIKE ? OR author LIKE ? OR category LIKE ?",
		searchPattern, searchPattern, searchPattern, searchPattern,
	).Find(&plugins).Error; err != nil {
		return nil, err
	}

	return plugins, nil
}

// GetByCategory returns plugins by category
func (s *RegistryService) GetByCategory(category string) ([]models.PluginRegistry, error) {
	var plugins []models.PluginRegistry
	if err := s.db.Where("category = ?", category).Find(&plugins).Error; err != nil {
		return nil, err
	}

	return plugins, nil
}

// ForceSyncNow forces an immediate sync
func (s *RegistryService) ForceSyncNow() error {
	s.lastSync = time.Time{} // Reset cache
	return s.Sync()
}

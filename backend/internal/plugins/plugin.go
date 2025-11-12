package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Plugin represents a plugin in the system
type Plugin struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Author      string                 `json:"author"`
	Description string                 `json:"description"`
	Icon        string                 `json:"icon,omitempty"`
	Enabled     bool                   `json:"enabled"`
	Installed   bool                   `json:"installed"`
	InstallPath string                 `json:"installPath,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// PluginManifest represents the plugin.json manifest file
type PluginManifest struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Author      string                 `json:"author"`
	Description string                 `json:"description"`
	Icon        string                 `json:"icon,omitempty"`
	EntryPoint  string                 `json:"entryPoint,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// Service handles plugin operations
type Service struct {
	pluginsDir string
	plugins    map[string]*Plugin
	mu         sync.RWMutex
}

var (
	globalService *Service
	once          sync.Once
)

const (
	DefaultPluginsDir = "/var/lib/stumpfworks/plugins"
)

// Initialize initializes the plugin service
func Initialize(pluginsDir string) (*Service, error) {
	var err error
	once.Do(func() {
		if pluginsDir == "" {
			pluginsDir = DefaultPluginsDir
		}

		// Ensure plugins directory exists
		if err = os.MkdirAll(pluginsDir, 0755); err != nil {
			return
		}

		globalService = &Service{
			pluginsDir: pluginsDir,
			plugins:    make(map[string]*Plugin),
		}

		// Discover installed plugins
		if err = globalService.discoverPlugins(); err != nil {
			return
		}
	})

	return globalService, err
}

// GetService returns the global plugin service
func GetService() *Service {
	return globalService
}

// discoverPlugins scans the plugins directory and loads plugin manifests
func (s *Service) discoverPlugins() error {
	entries, err := os.ReadDir(s.pluginsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No plugins directory yet
		}
		return fmt.Errorf("failed to read plugins directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginPath := filepath.Join(s.pluginsDir, entry.Name())
		manifestPath := filepath.Join(pluginPath, "plugin.json")

		// Check if plugin.json exists
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			continue
		}

		// Load manifest
		manifest, err := s.loadManifest(manifestPath)
		if err != nil {
			// Log error but continue with other plugins
			continue
		}

		// Create plugin entry
		plugin := &Plugin{
			ID:          manifest.ID,
			Name:        manifest.Name,
			Version:     manifest.Version,
			Author:      manifest.Author,
			Description: manifest.Description,
			Icon:        manifest.Icon,
			Enabled:     false, // Default to disabled
			Installed:   true,
			InstallPath: pluginPath,
			Config:      manifest.Config,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		s.plugins[plugin.ID] = plugin
	}

	return nil
}

// loadManifest loads a plugin manifest from a file
func (s *Service) loadManifest(path string) (*PluginManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest PluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

// ListPlugins returns all plugins
func (s *Service) ListPlugins(ctx context.Context) ([]*Plugin, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	plugins := make([]*Plugin, 0, len(s.plugins))
	for _, plugin := range s.plugins {
		plugins = append(plugins, plugin)
	}

	return plugins, nil
}

// GetPlugin returns a specific plugin by ID
func (s *Service) GetPlugin(ctx context.Context, id string) (*Plugin, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	plugin, ok := s.plugins[id]
	if !ok {
		return nil, fmt.Errorf("plugin not found: %s", id)
	}

	return plugin, nil
}

// InstallPlugin installs a plugin from a source path or URL
func (s *Service) InstallPlugin(ctx context.Context, sourcePath string) (*Plugin, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Load manifest from source
	manifestPath := filepath.Join(sourcePath, "plugin.json")
	manifest, err := s.loadManifest(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin manifest: %w", err)
	}

	// Check if plugin already exists
	if _, exists := s.plugins[manifest.ID]; exists {
		return nil, fmt.Errorf("plugin already installed: %s", manifest.ID)
	}

	// Create plugin installation directory
	installPath := filepath.Join(s.pluginsDir, manifest.ID)
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Copy plugin files to installation directory
	if err := copyDir(sourcePath, installPath); err != nil {
		os.RemoveAll(installPath) // Cleanup on error
		return nil, fmt.Errorf("failed to copy plugin files: %w", err)
	}

	// Create plugin entry
	plugin := &Plugin{
		ID:          manifest.ID,
		Name:        manifest.Name,
		Version:     manifest.Version,
		Author:      manifest.Author,
		Description: manifest.Description,
		Icon:        manifest.Icon,
		Enabled:     false,
		Installed:   true,
		InstallPath: installPath,
		Config:      manifest.Config,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.plugins[plugin.ID] = plugin

	return plugin, nil
}

// UninstallPlugin removes a plugin
func (s *Service) UninstallPlugin(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	plugin, ok := s.plugins[id]
	if !ok {
		return fmt.Errorf("plugin not found: %s", id)
	}

	// Disable plugin first if enabled
	if plugin.Enabled {
		plugin.Enabled = false
	}

	// Remove plugin directory
	if plugin.InstallPath != "" {
		if err := os.RemoveAll(plugin.InstallPath); err != nil {
			return fmt.Errorf("failed to remove plugin directory: %w", err)
		}
	}

	// Remove from plugins map
	delete(s.plugins, id)

	return nil
}

// EnablePlugin enables a plugin
func (s *Service) EnablePlugin(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	plugin, ok := s.plugins[id]
	if !ok {
		return fmt.Errorf("plugin not found: %s", id)
	}

	if !plugin.Installed {
		return fmt.Errorf("plugin not installed: %s", id)
	}

	plugin.Enabled = true
	plugin.UpdatedAt = time.Now()

	return nil
}

// DisablePlugin disables a plugin
func (s *Service) DisablePlugin(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	plugin, ok := s.plugins[id]
	if !ok {
		return fmt.Errorf("plugin not found: %s", id)
	}

	plugin.Enabled = false
	plugin.UpdatedAt = time.Now()

	return nil
}

// UpdatePluginConfig updates a plugin's configuration
func (s *Service) UpdatePluginConfig(ctx context.Context, id string, config map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	plugin, ok := s.plugins[id]
	if !ok {
		return fmt.Errorf("plugin not found: %s", id)
	}

	plugin.Config = config
	plugin.UpdatedAt = time.Now()

	return nil
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Determine destination path
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}

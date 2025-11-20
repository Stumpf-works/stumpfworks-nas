package handlers

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/Stumpf-works/stumpfworks-nas/internal/api/utils"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/errors"
	"github.com/Stumpf-works/stumpfworks-nas/internal/plugins"
)

const (
	PluginInstallPath = "/var/lib/stumpfworks/plugins"
)

var registryService *plugins.RegistryService

// InitPluginStore initializes the plugin store
func InitPluginStore() {
	db := database.GetDB()
	registryService = plugins.NewRegistryService(db, "")
}

// ListAvailablePlugins returns all available plugins from registry
func ListAvailablePlugins(w http.ResponseWriter, r *http.Request) {
	plugins, err := registryService.List()
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to list plugins", err))
		return
	}

	utils.RespondSuccess(w, plugins)
}

// GetPluginFromRegistry returns a specific plugin from registry
func GetPluginFromRegistry(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "id")

	plugin, err := registryService.Get(pluginID)
	if err != nil {
		utils.RespondError(w, errors.NotFound("Plugin not found", err))
		return
	}

	utils.RespondSuccess(w, plugin)
}

// SearchPlugins searches for plugins
func SearchPlugins(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")

	var plugins []models.PluginRegistry
	var err error

	if category != "" {
		plugins, err = registryService.GetByCategory(category)
	} else if query != "" {
		plugins, err = registryService.Search(query)
	} else {
		plugins, err = registryService.List()
	}

	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to search plugins", err))
		return
	}

	utils.RespondSuccess(w, plugins)
}

// InstallPlugin installs a plugin from registry
func InstallPlugin(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "id")

	log.Info().Str("plugin_id", pluginID).Msg("Installing plugin from registry")

	// Get plugin metadata
	plugin, err := registryService.Get(pluginID)
	if err != nil {
		utils.RespondError(w, errors.NotFound("Plugin not found", err))
		return
	}

	// Check if already installed
	db := database.GetDB()
	var existing models.InstalledPlugin
	if err := db.Where("id = ?", pluginID).First(&existing).Error; err == nil {
		utils.RespondError(w, errors.BadRequest("Plugin already installed", nil))
		return
	}

	// Download plugin
	downloadURL := plugin.DownloadURL
	if downloadURL == "" {
		utils.RespondError(w, errors.BadRequest("Plugin has no download URL", nil))
		return
	}

	log.Info().Str("url", downloadURL).Msg("Downloading plugin")

	resp, err := http.Get(downloadURL)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to download plugin", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.RespondError(w, errors.InternalServerError(
			fmt.Sprintf("Download failed with status %d", resp.StatusCode), nil))
		return
	}

	// Create plugin directory
	pluginPath := filepath.Join(PluginInstallPath, pluginID)
	if err := os.MkdirAll(pluginPath, 0755); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to create plugin directory", err))
		return
	}

	// Extract tar.gz
	if err := extractTarGz(resp.Body, pluginPath); err != nil {
		os.RemoveAll(pluginPath) // Cleanup on error
		utils.RespondError(w, errors.InternalServerError("Failed to extract plugin", err))
		return
	}

	// Save installation record
	installed := models.InstalledPlugin{
		ID:          pluginID,
		Version:     plugin.Version,
		InstallPath: pluginPath,
		Enabled:     false, // Start disabled
		AutoUpdate:  true,
		Status:      "installed",
	}

	if err := db.Create(&installed).Error; err != nil {
		os.RemoveAll(pluginPath) // Cleanup on error
		utils.RespondError(w, errors.InternalServerError("Failed to save installation record", err))
		return
	}

	log.Info().Str("plugin_id", pluginID).Msg("Plugin installed successfully")

	utils.RespondSuccess(w, map[string]interface{}{
		"message": "Plugin installed successfully",
		"plugin":  installed,
	})
}

// UninstallPlugin uninstalls a plugin
func UninstallPlugin(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "id")

	log.Info().Str("plugin_id", pluginID).Msg("Uninstalling plugin")

	db := database.GetDB()

	// Get installation record
	var installed models.InstalledPlugin
	if err := db.Where("id = ?", pluginID).First(&installed).Error; err != nil {
		utils.RespondError(w, errors.NotFound("Plugin not installed", err))
		return
	}

	// Stop plugin if running
	// TODO: Implement plugin stop

	// Remove plugin directory
	if err := os.RemoveAll(installed.InstallPath); err != nil {
		log.Warn().Err(err).Msg("Failed to remove plugin directory")
	}

	// Delete installation record
	if err := db.Delete(&installed).Error; err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to delete installation record", err))
		return
	}

	log.Info().Str("plugin_id", pluginID).Msg("Plugin uninstalled successfully")

	utils.RespondSuccess(w, map[string]interface{}{
		"message": "Plugin uninstalled successfully",
	})
}

// UpdatePlugin updates a plugin to the latest version
func UpdatePlugin(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "id")

	log.Info().Str("plugin_id", pluginID).Msg("Updating plugin")

	// Get current installation
	db := database.GetDB()
	var installed models.InstalledPlugin
	if err := db.Where("id = ?", pluginID).First(&installed).Error; err != nil {
		utils.RespondError(w, errors.NotFound("Plugin not installed", err))
		return
	}

	// Get latest version from registry
	plugin, err := registryService.Get(pluginID)
	if err != nil {
		utils.RespondError(w, errors.NotFound("Plugin not found in registry", err))
		return
	}

	// Check if update needed
	if installed.Version == plugin.Version {
		utils.RespondSuccess(w, map[string]interface{}{
			"message": "Plugin is already up to date",
			"version": installed.Version,
		})
		return
	}

	// Uninstall old version
	if err := os.RemoveAll(installed.InstallPath); err != nil {
		log.Warn().Err(err).Msg("Failed to remove old plugin directory")
	}

	// Download new version (same logic as Install)
	resp, err := http.Get(plugin.DownloadURL)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to download update", err))
		return
	}
	defer resp.Body.Close()

	// Extract
	if err := extractTarGz(resp.Body, installed.InstallPath); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to extract update", err))
		return
	}

	// Update record
	installed.Version = plugin.Version
	installed.Status = "updated"
	if err := db.Save(&installed).Error; err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to update record", err))
		return
	}

	log.Info().
		Str("plugin_id", pluginID).
		Str("version", plugin.Version).
		Msg("Plugin updated successfully")

	utils.RespondSuccess(w, map[string]interface{}{
		"message": "Plugin updated successfully",
		"version": plugin.Version,
	})
}

// SyncRegistry forces a registry sync
func SyncRegistry(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Forcing registry sync")

	if err := registryService.ForceSyncNow(); err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to sync registry", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"message": "Registry synced successfully",
	})
}

// ListInstalledPlugins returns all installed plugins
func ListInstalledPlugins(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()

	var installed []models.InstalledPlugin
	if err := db.Find(&installed).Error; err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to list installed plugins", err))
		return
	}

	utils.RespondSuccess(w, installed)
}

// extractTarGz extracts a tar.gz archive
func extractTarGz(src io.Reader, destPath string) error {
	gzr, err := gzip.NewReader(src)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %w", err)
		}

		// Security: Prevent path traversal
		target := filepath.Join(destPath, header.Name)
		if !strings.HasPrefix(target, filepath.Clean(destPath)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

		case tar.TypeReg:
			// Create parent directories
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			// Create file
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return fmt.Errorf("failed to extract file: %w", err)
			}

			f.Close()
		}
	}

	return nil
}

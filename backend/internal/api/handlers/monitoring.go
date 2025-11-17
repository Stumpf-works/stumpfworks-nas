// Revision: 2025-11-17 | Author: Claude | Version: 1.0.0
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// GetMonitoringConfig retrieves the monitoring configuration
func GetMonitoringConfig(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()

	var config models.MonitoringConfig
	result := db.First(&config)

	// If no config exists, create default
	if result.Error != nil {
		config = models.MonitoringConfig{
			PrometheusEnabled: true,
			GrafanaURL:        "http://localhost:3000",
			DatadogEnabled:    false,
			DatadogAPIKey:     "",
		}
		if err := db.Create(&config).Error; err != nil {
			logger.Error("Failed to create default monitoring config", zap.Error(err))
			http.Error(w, "Failed to create config", http.StatusInternalServerError)
			return
		}
	}

	// Don't expose API key in full, only indicate if it's set
	response := map[string]interface{}{
		"prometheus_enabled": config.PrometheusEnabled,
		"grafana_url":        config.GrafanaURL,
		"datadog_enabled":    config.DatadogEnabled,
		"datadog_api_key_set": config.DatadogAPIKey != "",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateMonitoringConfig updates the monitoring configuration
func UpdateMonitoringConfig(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()

	var input struct {
		PrometheusEnabled bool   `json:"prometheus_enabled"`
		GrafanaURL        string `json:"grafana_url"`
		DatadogEnabled    bool   `json:"datadog_enabled"`
		DatadogAPIKey     string `json:"datadog_api_key,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var config models.MonitoringConfig
	result := db.First(&config)

	// If no config exists, create new
	if result.Error != nil {
		config = models.MonitoringConfig{
			PrometheusEnabled: input.PrometheusEnabled,
			GrafanaURL:        input.GrafanaURL,
			DatadogEnabled:    input.DatadogEnabled,
			DatadogAPIKey:     input.DatadogAPIKey,
		}
		if err := db.Create(&config).Error; err != nil {
			logger.Error("Failed to create monitoring config", zap.Error(err))
			http.Error(w, "Failed to create config", http.StatusInternalServerError)
			return
		}
	} else {
		// Update existing config
		config.PrometheusEnabled = input.PrometheusEnabled
		config.GrafanaURL = input.GrafanaURL
		config.DatadogEnabled = input.DatadogEnabled

		// Only update API key if provided (non-empty)
		if input.DatadogAPIKey != "" {
			config.DatadogAPIKey = input.DatadogAPIKey
		}

		if err := db.Save(&config).Error; err != nil {
			logger.Error("Failed to update monitoring config", zap.Error(err))
			http.Error(w, "Failed to update config", http.StatusInternalServerError)
			return
		}
	}

	logger.Info("Monitoring configuration updated")

	response := map[string]interface{}{
		"success": true,
		"message": "Monitoring configuration updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

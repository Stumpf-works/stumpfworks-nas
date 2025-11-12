package handlers

import (
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/config"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
)

// HealthCheck returns the health status of the API
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	cfg := config.GlobalConfig

	utils.RespondSuccess(w, map[string]interface{}{
		"status":  "ok",
		"service": cfg.App.Name,
		"version": cfg.App.Version,
	})
}

// IndexHandler returns basic API information
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.GlobalConfig

	utils.RespondSuccess(w, map[string]interface{}{
		"name":        cfg.App.Name,
		"version":     cfg.App.Version,
		"environment": cfg.App.Environment,
		"api_version": "v1",
		"endpoints": map[string]string{
			"health":  "/health",
			"api":     "/api/v1",
			"ws":      "/ws",
			"docs":    "/api/v1/docs (coming soon)",
		},
	})
}

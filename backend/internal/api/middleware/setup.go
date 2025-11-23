package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
)

// SetupRequired checks if initial setup is needed (no admin user exists)
// If setup is required, blocks all requests except setup endpoints
func SetupRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip setup check for setup wizard endpoints and health check
		if r.URL.Path == "/api/v1/setup/status" ||
			r.URL.Path == "/api/v1/setup/initialize" ||
			r.URL.Path == "/health" ||
			r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		// Check if any admin user exists
		db := database.GetDB()
		if db == nil {
			// Database not initialized yet, allow request
			next.ServeHTTP(w, r)
			return
		}

		var count int64
		db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

		if count == 0 {
			// No admin exists - setup required
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error": map[string]string{
					"code":    "SETUP_REQUIRED",
					"message": "Initial setup required. Please complete the setup wizard.",
				},
				"setupRequired": true,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

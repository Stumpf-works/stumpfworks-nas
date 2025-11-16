// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package middleware

import (
	"net/http"
	"runtime"

	"github.com/Stumpf-works/stumpfworks-nas/internal/updates"
)

// RevisionMiddleware adds version and build information headers to all responses
func RevisionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add revision headers
		w.Header().Set("X-StumpfWorks-Version", updates.CurrentVersion)
		w.Header().Set("X-StumpfWorks-Go-Version", runtime.Version())
		w.Header().Set("X-StumpfWorks-API-Version", "v1")

		// Optional: Add build date if available (can be set via ldflags during build)
		// w.Header().Set("X-StumpfWorks-Build-Date", BuildDate)

		next.ServeHTTP(w, r)
	})
}

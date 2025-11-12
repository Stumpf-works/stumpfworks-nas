package middleware

import (
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/auth"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
)

// IPBlockMiddleware checks if an IP is blocked before allowing access
func IPBlockMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP
		ipAddress := getClientIP(r)

		// Check if IP is blocked
		service := auth.GetFailedLoginService()
		if service != nil {
			blocked, block, err := service.IsIPBlocked(ipAddress)
			if err != nil {
				// Log error but don't block request
				next.ServeHTTP(w, r)
				return
			}

			if blocked {
				// Return 403 Forbidden with block info
				utils.RespondError(w, errors.Forbidden(
					"Your IP address has been temporarily blocked due to too many failed login attempts. "+
					"Reason: "+block.Reason,
					nil))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

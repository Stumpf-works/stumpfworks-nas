package api

import (
	"github.com/gin-gonic/gin"

	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/config"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get allowed origins from config
		allowedOrigins := cfg.Security.AllowedOrigins
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if origin != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting
func RateLimitMiddleware(cfg *config.Config) gin.HandlerFunc {
	// TODO: Implement proper rate limiting
	// For now, this is a placeholder
	return func(c *gin.Context) {
		c.Next()
	}
}

// AuthMiddleware handles authentication
func AuthMiddleware() gin.HandlerFunc {
	// TODO: Implement JWT authentication
	// For now, this is a placeholder
	return func(c *gin.Context) {
		c.Next()
	}
}

// AuditLogMiddleware logs all API requests
func AuditLogMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement audit logging to file
		// For now, just use standard logging
		if cfg.Security.EnableAuditLogging {
			// Log request details
			// In production, write to audit log file
		}

		c.Next()
	}
}

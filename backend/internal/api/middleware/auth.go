package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/users"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware validates JWT tokens and adds user to context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		logger.Info("Auth header received", zap.String("header", authHeader), zap.Int("length", len(authHeader)))
		if authHeader == "" {
			utils.RespondError(w, errors.Unauthorized("Missing authorization header", nil))
			return
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		logger.Info("Auth header split", zap.Int("parts", len(parts)), zap.Strings("parts", parts))
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid auth header format", zap.Int("parts", len(parts)))
			utils.RespondError(w, errors.Unauthorized("Invalid authorization header format", nil))
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := users.ValidateToken(tokenString)
		if err != nil {
			utils.RespondError(w, errors.Unauthorized("Invalid or expired token", err))
			return
		}

		// Get user from database
		user, err := users.GetUserByID(claims.UserID)
		if err != nil {
			utils.RespondError(w, errors.Unauthorized("User not found", err))
			return
		}

		// Check if user is active
		if !user.IsActive {
			utils.RespondError(w, errors.Forbidden("User account is disabled", nil))
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnly middleware ensures user has admin role
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil {
			utils.RespondError(w, errors.Unauthorized("User not found in context", nil))
			return
		}

		if !user.IsAdmin() {
			utils.RespondError(w, errors.Forbidden("Admin access required", nil))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext retrieves user from request context
func GetUserFromContext(ctx context.Context) *users.User {
	user, ok := ctx.Value(UserContextKey).(*users.User)
	if !ok {
		return nil
	}
	return user
}

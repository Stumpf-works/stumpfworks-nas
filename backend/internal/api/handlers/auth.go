package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/api/middleware"
	"github.com/Stumpf-works/stumpfworks-nas/internal/auth"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/users"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string              `json:"accessToken"`
	RefreshToken string              `json:"refreshToken"`
	User         *users.UserResponse `json:"user"`
}

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Get client IP and user agent
	ipAddress := getClientIP(r)
	userAgent := r.UserAgent()

	// Authenticate user
	user, err := users.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		// Track failed login attempt
		failedLoginService := auth.GetFailedLoginService()
		if failedLoginService != nil {
			// Determine failure reason
			reason := models.FailureReasonInvalidPassword
			if strings.Contains(err.Error(), "not found") {
				reason = models.FailureReasonUserNotFound
			} else if strings.Contains(err.Error(), "disabled") {
				reason = models.FailureReasonAccountDisabled
			}

			// Record failed attempt
			if recordErr := failedLoginService.RecordFailedAttempt(
				r.Context(),
				req.Username,
				ipAddress,
				userAgent,
				reason,
			); recordErr != nil {
				logger.Error("Failed to record login attempt",
					zap.Error(recordErr),
					zap.String("username", req.Username))
			}
		}

		utils.RespondError(w, err)
		return
	}

	// Generate tokens
	accessToken, err := users.GenerateToken(user)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to generate access token", err))
		return
	}

	refreshToken, err := users.GenerateRefreshToken(user)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to generate refresh token", err))
		return
	}

	// Return response
	utils.RespondSuccess(w, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         users.ToResponse(user),
	})
}

// getClientIP extracts the real IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP if multiple are present
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// Logout handles user logout
func Logout(w http.ResponseWriter, r *http.Request) {
	// In a more complex system, we would invalidate the token here
	// For now, just return success (client will remove token)
	utils.RespondNoContent(w)
}

// RefreshToken handles token refresh
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate refresh token
	claims, err := users.ValidateToken(req.RefreshToken)
	if err != nil {
		utils.RespondError(w, errors.Unauthorized("Invalid refresh token", err))
		return
	}

	// Get user
	user, err := users.GetUserByID(claims.UserID)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	// Generate new access token
	accessToken, err := users.GenerateToken(user)
	if err != nil {
		utils.RespondError(w, errors.InternalServerError("Failed to generate access token", err))
		return
	}

	// Return new token
	utils.RespondSuccess(w, map[string]string{
		"accessToken": accessToken,
	})
}

// GetCurrentUser returns the currently authenticated user
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		utils.RespondError(w, errors.Unauthorized("User not found", nil))
		return
	}

	utils.RespondSuccess(w, users.ToResponse(user))
}

// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/api/middleware"
	"github.com/Stumpf-works/stumpfworks-nas/internal/auth"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/twofa"
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
	Requires2FA  bool                `json:"requires2FA,omitempty"`
	UserID       uint                `json:"userId,omitempty"`
	AccessToken  string              `json:"accessToken,omitempty"`
	RefreshToken string              `json:"refreshToken,omitempty"`
	User         *users.UserResponse `json:"user,omitempty"`
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

	// Check if 2FA is enabled for this user
	twofaService := twofa.GetService()
	if twofaService != nil {
		enabled, err := twofaService.IsEnabled(r.Context(), user.ID)
		if err != nil {
			logger.Error("Failed to check 2FA status", zap.Error(err))
		}

		if enabled {
			// Return a response indicating 2FA is required
			utils.RespondSuccess(w, LoginResponse{
				Requires2FA: true,
				UserID:      user.ID,
			})
			return
		}
	}

	// Generate tokens (no 2FA required or 2FA not enabled)
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

// LoginWith2FA completes the login process after 2FA verification
func LoginWith2FA(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID       uint   `json:"userId"`
		Code         string `json:"code"`
		IsBackupCode bool   `json:"isBackupCode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Code == "" {
		utils.RespondError(w, errors.BadRequest("Verification code is required", nil))
		return
	}

	// Verify 2FA code
	twofaService := twofa.GetService()
	if twofaService == nil {
		utils.RespondError(w, errors.InternalServerError("2FA service not available", nil))
		return
	}

	valid, err := twofaService.VerifyCode(r.Context(), twofa.VerifyRequest{
		UserID:       req.UserID,
		Code:         req.Code,
		IsBackupCode: req.IsBackupCode,
	})

	if err != nil {
		logger.Error("Failed to verify 2FA code", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to verify code", err))
		return
	}

	if !valid {
		utils.RespondError(w, errors.Unauthorized("Invalid verification code", nil))
		return
	}

	// Get user
	user, err := users.GetUserByID(req.UserID)
	if err != nil {
		utils.RespondError(w, errors.NotFound("User not found", err))
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

// ResetPasswordRequest represents a password reset request using a token
type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

// ResetPasswordResponse represents the response after password reset
type ResetPasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ResetPasswordWithToken handles password reset using a token
func ResetPasswordWithToken(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate input
	if req.Token == "" {
		utils.RespondError(w, errors.BadRequest("Token is required", nil))
		return
	}

	if req.NewPassword == "" {
		utils.RespondError(w, errors.BadRequest("New password is required", nil))
		return
	}

	if len(req.NewPassword) < 8 {
		utils.RespondError(w, errors.BadRequest("Password must be at least 8 characters long", nil))
		return
	}

	// Find and validate token
	resetToken, err := models.FindValidToken(auth.GetDB(), req.Token)
	if err != nil {
		logger.Warn("Invalid password reset token attempt",
			zap.String("token", req.Token[:8]+"..."), // Only log first 8 chars
			zap.String("ip", getClientIP(r)))
		utils.RespondError(w, errors.Unauthorized("Invalid or expired reset token"))
		return
	}

	// Verify token is still valid
	if !resetToken.IsValid() {
		utils.RespondError(w, errors.Unauthorized("Token has expired or already been used"))
		return
	}

	// Get user
	var user models.User
	if err := auth.GetDB().First(&user, resetToken.UserID).Error; err != nil {
		logger.Error("Failed to find user for password reset",
			zap.Error(err),
			zap.Uint("userId", resetToken.UserID))
		utils.RespondError(w, errors.InternalServerError("Failed to reset password", err))
		return
	}

	// Update password
	if err := user.SetPassword(req.NewPassword); err != nil {
		logger.Error("Failed to hash new password",
			zap.Error(err),
			zap.String("username", user.Username))
		utils.RespondError(w, errors.InternalServerError("Failed to reset password", err))
		return
	}

	if err := auth.GetDB().Save(&user).Error; err != nil {
		logger.Error("Failed to save new password",
			zap.Error(err),
			zap.String("username", user.Username))
		utils.RespondError(w, errors.InternalServerError("Failed to reset password", err))
		return
	}

	// Mark token as used
	if err := resetToken.MarkAsUsed(auth.GetDB()); err != nil {
		logger.Error("Failed to mark reset token as used",
			zap.Error(err),
			zap.Uint("tokenId", resetToken.ID))
		// Don't fail the request - password was already changed
	}

	logger.Info("Password reset successful",
		zap.String("username", user.Username),
		zap.String("ip", getClientIP(r)))

	utils.RespondSuccess(w, ResetPasswordResponse{
		Success: true,
		Message: "Password has been reset successfully. You can now log in with your new password.",
	})
}

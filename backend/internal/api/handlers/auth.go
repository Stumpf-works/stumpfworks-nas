package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stumpfworks/nas/internal/api/middleware"
	"github.com/stumpfworks/nas/internal/users"
	"github.com/stumpfworks/nas/pkg/errors"
	"github.com/stumpfworks/nas/pkg/utils"
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

	// Authenticate user
	user, err := users.AuthenticateUser(req.Username, req.Password)
	if err != nil {
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

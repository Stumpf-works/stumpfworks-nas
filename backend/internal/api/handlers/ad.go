package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/ad"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

// ADHandler handles Active Directory-related HTTP requests
type ADHandler struct {
	service *ad.Service
}

// NewADHandler creates a new AD handler
func NewADHandler() *ADHandler {
	return &ADHandler{
		service: ad.GetService(),
	}
}

// GetConfig gets the current AD configuration
func (h *ADHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	config := h.service.GetConfig()
	utils.RespondSuccess(w, config)
}

// UpdateConfig updates the AD configuration
func (h *ADHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	var config ad.ADConfig

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := h.service.UpdateConfig(&config); err != nil {
		logger.Error("Failed to update AD config", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to update AD config", err))
		return
	}

	logger.Info("AD configuration updated", zap.Bool("enabled", config.Enabled))
	utils.RespondSuccess(w, map[string]string{"message": "AD configuration updated successfully"})
}

// TestConnection tests the AD connection
func (h *ADHandler) TestConnection(w http.ResponseWriter, r *http.Request) {
	if err := h.service.TestConnection(r.Context()); err != nil {
		logger.Error("AD connection test failed", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("AD connection test failed", err))
		return
	}

	logger.Info("AD connection test successful")
	utils.RespondSuccess(w, map[string]string{"message": "Connection successful"})
}

// Authenticate authenticates a user against AD
func (h *ADHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.RespondError(w, errors.BadRequest("Username and password are required", nil))
		return
	}

	user, err := h.service.Authenticate(r.Context(), req.Username, req.Password)
	if err != nil {
		logger.Error("AD authentication failed",
			zap.Error(err),
			zap.String("username", req.Username))
		utils.RespondError(w, errors.Unauthorized("Authentication failed", err))
		return
	}

	logger.Info("User authenticated via AD",
		zap.String("username", user.Username),
		zap.String("email", user.Email))
	utils.RespondSuccess(w, user)
}

// ListUsers lists all users from AD
func (h *ADHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		logger.Error("Failed to list AD users", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to list AD users", err))
		return
	}

	utils.RespondSuccess(w, users)
}

// SyncUser synchronizes a user from AD
func (h *ADHandler) SyncUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Username == "" {
		utils.RespondError(w, errors.BadRequest("Username is required", nil))
		return
	}

	user, err := h.service.SyncUser(r.Context(), req.Username)
	if err != nil {
		logger.Error("Failed to sync AD user",
			zap.Error(err),
			zap.String("username", req.Username))
		utils.RespondError(w, errors.InternalServerError("Failed to sync AD user", err))
		return
	}

	logger.Info("User synchronized from AD",
		zap.String("username", user.Username),
		zap.String("email", user.Email))
	utils.RespondSuccess(w, user)
}

// GetStatus returns the AD service status
func (h *ADHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"available": h.service.IsAvailable(),
		"enabled":   h.service.GetConfig().Enabled,
	}

	utils.RespondSuccess(w, status)
}

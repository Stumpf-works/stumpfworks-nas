package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/twofa"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"go.uber.org/zap"
)

// TwoFAHandler handles 2FA-related HTTP requests
type TwoFAHandler struct {
	service *twofa.Service
}

// NewTwoFAHandler creates a new 2FA handler
func NewTwoFAHandler() *TwoFAHandler {
	return &TwoFAHandler{
		service: twofa.GetService(),
	}
}

// GetStatus returns the 2FA status for the authenticated user
func (h *TwoFAHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := getUserIDFromContext(r)

	enabled, err := h.service.IsEnabled(ctx, userID)
	if err != nil {
		logger.Error("Failed to get 2FA status", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get 2FA status", err))
		return
	}

	// Get backup codes count if enabled
	var backupCodesRemaining int
	if enabled {
		count, err := h.service.GetBackupCodes(ctx, userID)
		if err != nil {
			logger.Error("Failed to get backup codes count", zap.Error(err))
		} else {
			backupCodesRemaining = count
		}
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"enabled":               enabled,
		"backupCodesRemaining":  backupCodesRemaining,
	})
}

// SetupTwoFactor initiates 2FA setup for the authenticated user
func (h *TwoFAHandler) SetupTwoFactor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := getUserIDFromContext(r)

	response, err := h.service.SetupTwoFactor(ctx, twofa.SetupRequest{
		UserID: userID,
		Issuer: "Stumpf.Works NAS",
	})
	if err != nil {
		logger.Error("Failed to setup 2FA", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to setup 2FA", err))
		return
	}

	utils.RespondSuccess(w, response)
}

// EnableTwoFactor enables 2FA after verifying the initial code
func (h *TwoFAHandler) EnableTwoFactor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := getUserIDFromContext(r)

	var req struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Code == "" {
		utils.RespondError(w, errors.BadRequest("Verification code is required", nil))
		return
	}

	if err := h.service.EnableTwoFactor(ctx, userID, req.Code); err != nil {
		logger.Error("Failed to enable 2FA", zap.Error(err))
		utils.RespondError(w, errors.BadRequest("Failed to enable 2FA", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Two-factor authentication enabled successfully",
	})
}

// DisableTwoFactor disables 2FA for the authenticated user
func (h *TwoFAHandler) DisableTwoFactor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := getUserIDFromContext(r)

	var req struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Code == "" {
		utils.RespondError(w, errors.BadRequest("Verification code is required", nil))
		return
	}

	if err := h.service.DisableTwoFactor(ctx, userID, req.Code); err != nil {
		logger.Error("Failed to disable 2FA", zap.Error(err))
		utils.RespondError(w, errors.BadRequest("Failed to disable 2FA", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Two-factor authentication disabled successfully",
	})
}

// RegenerateBackupCodes generates new backup codes
func (h *TwoFAHandler) RegenerateBackupCodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := getUserIDFromContext(r)

	var req struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Code == "" {
		utils.RespondError(w, errors.BadRequest("Verification code is required", nil))
		return
	}

	backupCodes, err := h.service.RegenerateBackupCodes(ctx, userID, req.Code)
	if err != nil {
		logger.Error("Failed to regenerate backup codes", zap.Error(err))
		utils.RespondError(w, errors.BadRequest("Failed to regenerate backup codes", err))
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"backupCodes": backupCodes,
	})
}

// VerifyTwoFactor verifies a 2FA code during login
func (h *TwoFAHandler) VerifyTwoFactor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	valid, err := h.service.VerifyCode(ctx, twofa.VerifyRequest{
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

	utils.RespondSuccess(w, map[string]bool{
		"valid": true,
	})
}

// getUserIDFromContext extracts the user ID from the request context
func getUserIDFromContext(r *http.Request) uint {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		return 0
	}
	return userID
}

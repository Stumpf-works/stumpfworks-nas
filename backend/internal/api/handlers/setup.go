package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/go-playground/validator/v10"
)

var setupValidator = validator.New()

// SetupStatusResponse represents the setup status
type SetupStatusResponse struct {
	SetupRequired bool `json:"setupRequired"`
	AdminExists   bool `json:"adminExists"`
}

// InitialSetupRequest represents the initial setup request
type InitialSetupRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FullName  string `json:"fullName" validate:"required,min=2,max=100"`
}

// SetupStatus returns the current setup status
func SetupStatus(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	if db == nil {
		respondJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "DATABASE_UNAVAILABLE",
				"message": "Database connection not available",
			},
		})
		return
	}

	var count int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

	respondJSON(w, http.StatusOK, SetupStatusResponse{
		SetupRequired: count == 0,
		AdminExists:   count > 0,
	})
}

// InitializeSetup creates the initial admin user
func InitializeSetup(w http.ResponseWriter, r *http.Request) {
	var req InitialSetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "INVALID_JSON",
				"message": "Invalid request body",
				"details": err.Error(),
			},
		})
		return
	}

	// Validate request
	if err := setupValidator.Struct(req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
				"details": err.Error(),
			},
		})
		return
	}

	db := database.GetDB()
	if db == nil {
		respondJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "DATABASE_UNAVAILABLE",
				"message": "Database connection not available",
			},
		})
		return
	}

	// Check if admin already exists (prevent multiple initialization)
	var count int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count > 0 {
		respondJSON(w, http.StatusConflict, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "ALREADY_INITIALIZED",
				"message": "System has already been initialized",
			},
		})
		return
	}

	// Check if username already exists
	var existingUser models.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		respondJSON(w, http.StatusConflict, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "USERNAME_EXISTS",
				"message": "Username already exists",
			},
		})
		return
	}

	// Check if email already exists
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		respondJSON(w, http.StatusConflict, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "EMAIL_EXISTS",
				"message": "Email already exists",
			},
		})
		return
	}

	// Create admin user
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Role:     "admin",
		IsActive: true,
	}

	// Set password
	if err := user.SetPassword(req.Password); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "PASSWORD_HASH_ERROR",
				"message": "Failed to hash password",
				"details": err.Error(),
			},
		})
		return
	}

	// Save user to database
	if err := db.Create(&user).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "DATABASE_ERROR",
				"message": "Failed to create admin user",
				"details": err.Error(),
			},
		})
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"message":  "Initial setup completed successfully",
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// respondJSON is a helper function to send JSON responses
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

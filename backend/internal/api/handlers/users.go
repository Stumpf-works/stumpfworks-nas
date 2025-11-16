// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Stumpf-works/stumpfworks-nas/internal/users"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
)

// ListUsers returns all users
func ListUsers(w http.ResponseWriter, r *http.Request) {
	userList, err := users.ListUsers()
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, users.ToResponses(userList))
}

// GetUser returns a single user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid user ID", err))
		return
	}

	user, err := users.GetUserByID(uint(id))
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, users.ToResponse(user))
}

// CreateUser creates a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req users.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate required fields
	if req.Username == "" {
		utils.RespondError(w, errors.BadRequest("Username is required", nil))
		return
	}
	if len(req.Username) < 3 || len(req.Username) > 100 {
		utils.RespondError(w, errors.BadRequest("Username must be between 3 and 100 characters", nil))
		return
	}
	if req.Email == "" {
		utils.RespondError(w, errors.BadRequest("Email is required", nil))
		return
	}
	if req.Password == "" {
		utils.RespondError(w, errors.BadRequest("Password is required", nil))
		return
	}
	if len(req.Password) < 8 {
		utils.RespondError(w, errors.BadRequest("Password must be at least 8 characters", nil))
		return
	}
	if req.Role == "" {
		utils.RespondError(w, errors.BadRequest("Role is required", nil))
		return
	}
	if req.Role != "admin" && req.Role != "user" && req.Role != "guest" {
		utils.RespondError(w, errors.BadRequest("Role must be one of: admin, user, guest", nil))
		return
	}

	user, err := users.CreateUser(&req)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondCreated(w, users.ToResponse(user))
}

// UpdateUser updates an existing user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid user ID", err))
		return
	}

	var req users.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// Validate optional fields (only if provided)
	if req.Password != nil && len(*req.Password) < 8 {
		utils.RespondError(w, errors.BadRequest("Password must be at least 8 characters", nil))
		return
	}
	if req.Role != nil {
		role := *req.Role
		if role != "admin" && role != "user" && role != "guest" {
			utils.RespondError(w, errors.BadRequest("Role must be one of: admin, user, guest", nil))
			return
		}
	}

	user, err := users.UpdateUser(uint(id), &req)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, users.ToResponse(user))
}

// DeleteUser deletes a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid user ID", err))
		return
	}

	if err := users.DeleteUser(uint(id)); err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondNoContent(w)
}

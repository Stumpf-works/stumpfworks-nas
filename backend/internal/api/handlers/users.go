package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/stumpfworks/nas/internal/users"
	"github.com/stumpfworks/nas/pkg/errors"
	"github.com/stumpfworks/nas/pkg/utils"
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

	// TODO: Add validation

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

	// TODO: Add validation

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

// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Stumpf-works/stumpfworks-nas/internal/usergroups"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
)

// ListGroups returns all user groups
func ListGroups(w http.ResponseWriter, r *http.Request) {
	groupList, err := usergroups.ListGroups()
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, usergroups.ToResponses(groupList))
}

// GetGroup returns a single user group by ID
func GetGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid group ID", err))
		return
	}

	group, err := usergroups.GetGroupByID(uint(id))
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, usergroups.ToResponse(group))
}

// CreateGroup creates a new user group
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req usergroups.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	group, err := usergroups.CreateGroup(&req)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondCreated(w, usergroups.ToResponse(group))
}

// UpdateGroup updates an existing user group
func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid group ID", err))
		return
	}

	var req usergroups.UpdateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	group, err := usergroups.UpdateGroup(uint(id), &req)
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, usergroups.ToResponse(group))
}

// DeleteGroup deletes a user group
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid group ID", err))
		return
	}

	if err := usergroups.DeleteGroup(uint(id)); err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondNoContent(w)
}

// AddGroupMember adds a member to a user group
func AddGroupMember(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid group ID", err))
		return
	}

	var req usergroups.AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if err := usergroups.AddMember(uint(groupID), req.UserID); err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondNoContent(w)
}

// RemoveGroupMember removes a member from a user group
func RemoveGroupMember(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid group ID", err))
		return
	}

	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid user ID", err))
		return
	}

	if err := usergroups.RemoveMember(uint(groupID), uint(userID)); err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondNoContent(w)
}

// GetGroupMembers returns all members of a user group
func GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid group ID", err))
		return
	}

	members, err := usergroups.GetGroupMembers(uint(groupID))
	if err != nil {
		utils.RespondError(w, err)
		return
	}

	utils.RespondSuccess(w, members)
}

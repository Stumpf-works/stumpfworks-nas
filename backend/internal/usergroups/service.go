// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package usergroups

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"gorm.io/gorm"
)

// CreateGroupRequest represents a request to create a user group
type CreateGroupRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description"`
}

// UpdateGroupRequest represents a request to update a user group
type UpdateGroupRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty"`
}

// AddMemberRequest represents a request to add a member to a group
type AddMemberRequest struct {
	UserID uint `json:"userId" validate:"required"`
}

// CreateGroup creates a new user group
func CreateGroup(req *CreateGroupRequest) (*models.UserGroup, error) {
	// Check if group with this name already exists
	var existingGroup models.UserGroup
	err := database.DB.Where("name = ?", req.Name).First(&existingGroup).Error
	if err == nil {
		return nil, errors.Conflict("Group name already exists", nil)
	}

	// Create group
	group := &models.UserGroup{
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    false,
	}

	if err := database.DB.Create(group).Error; err != nil {
		return nil, errors.InternalServerError("Failed to create group", err)
	}

	// Sync to Unix system groups (non-fatal if fails)
	unixManager := GetUnixGroupManager()
	if unixManager.IsEnabled() {
		if err := unixManager.CreateUnixGroup(group); err != nil {
			// Log warning but don't fail - group is created in database
			// This allows the system to work even if Unix group creation fails
		}
	}

	return group, nil
}

// ListGroups retrieves all user groups
func ListGroups() ([]*models.UserGroup, error) {
	var groups []*models.UserGroup
	err := database.DB.Preload("Members").Order("created_at DESC").Find(&groups).Error
	if err != nil {
		return nil, errors.InternalServerError("Failed to query groups", err)
	}
	return groups, nil
}

// GetGroupByID retrieves a user group by ID
func GetGroupByID(id uint) (*models.UserGroup, error) {
	var group models.UserGroup
	err := database.DB.Preload("Members").First(&group, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Group not found", err)
		}
		return nil, errors.InternalServerError("Failed to query group", err)
	}
	return &group, nil
}

// GetGroupByName retrieves a user group by name
func GetGroupByName(name string) (*models.UserGroup, error) {
	var group models.UserGroup
	err := database.DB.Preload("Members").Where("name = ?", name).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Group not found", err)
		}
		return nil, errors.InternalServerError("Failed to query group", err)
	}
	return &group, nil
}

// UpdateGroup updates an existing user group
func UpdateGroup(id uint, req *UpdateGroupRequest) (*models.UserGroup, error) {
	group, err := GetGroupByID(id)
	if err != nil {
		return nil, err
	}

	// Prevent modifying system groups
	if group.IsSystem {
		return nil, errors.Forbidden("Cannot modify system group", nil)
	}

	// Update fields
	updates := make(map[string]interface{})

	if req.Name != nil {
		// Check if name is already in use by another group
		var existingGroup models.UserGroup
		err := database.DB.Where("name = ? AND id != ?", *req.Name, id).First(&existingGroup).Error
		if err == nil {
			return nil, errors.Conflict("Group name already exists", nil)
		}
		updates["name"] = *req.Name
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	// Perform update
	if err := database.DB.Model(group).Updates(updates).Error; err != nil {
		return nil, errors.InternalServerError("Failed to update group", err)
	}

	// Reload group
	if err := database.DB.Preload("Members").First(group, id).Error; err != nil {
		return nil, errors.InternalServerError("Failed to reload group", err)
	}

	return group, nil
}

// DeleteGroup deletes a user group (soft delete)
func DeleteGroup(id uint) error {
	group, err := GetGroupByID(id)
	if err != nil {
		return err
	}

	// Prevent deleting system groups
	if group.IsSystem {
		return errors.Forbidden("Cannot delete system group", nil)
	}

	// Remove all associations (clear members)
	if err := database.DB.Model(group).Association("Members").Clear(); err != nil {
		return errors.InternalServerError("Failed to clear group members", err)
	}

	// Delete group
	if err := database.DB.Delete(group).Error; err != nil {
		return errors.InternalServerError("Failed to delete group", err)
	}

	// Remove Unix system group (non-fatal if fails)
	unixManager := GetUnixGroupManager()
	if unixManager.IsEnabled() {
		if err := unixManager.DeleteUnixGroup(group); err != nil {
			// Log warning but don't fail - group is deleted from database
		}
	}

	return nil
}

// AddMember adds a user to a group
func AddMember(groupID uint, userID uint) error {
	group, err := GetGroupByID(groupID)
	if err != nil {
		return err
	}

	// Check if user exists
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFound("User not found", err)
		}
		return errors.InternalServerError("Failed to query user", err)
	}

	// Check if user is already a member
	for _, member := range group.Members {
		if member.ID == userID {
			return errors.Conflict("User is already a member of this group", nil)
		}
	}

	// Add member
	if err := database.DB.Model(group).Association("Members").Append(&user); err != nil {
		return errors.InternalServerError("Failed to add member to group", err)
	}

	// Sync to Unix system group (non-fatal if fails)
	unixManager := GetUnixGroupManager()
	if unixManager.IsEnabled() {
		if err := unixManager.AddUserToUnixGroup(user.Username, group); err != nil {
			// Log warning but don't fail - member is added in database
		}
	}

	return nil
}

// RemoveMember removes a user from a group
func RemoveMember(groupID uint, userID uint) error {
	group, err := GetGroupByID(groupID)
	if err != nil {
		return err
	}

	// Check if user exists
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFound("User not found", err)
		}
		return errors.InternalServerError("Failed to query user", err)
	}

	// Check if user is a member
	isMember := false
	for _, member := range group.Members {
		if member.ID == userID {
			isMember = true
			break
		}
	}

	if !isMember {
		return errors.NotFound("User is not a member of this group", nil)
	}

	// Remove member
	if err := database.DB.Model(group).Association("Members").Delete(&user); err != nil {
		return errors.InternalServerError("Failed to remove member from group", err)
	}

	// Sync to Unix system group (non-fatal if fails)
	unixManager := GetUnixGroupManager()
	if unixManager.IsEnabled() {
		if err := unixManager.RemoveUserFromUnixGroup(user.Username, group); err != nil {
			// Log warning but don't fail - member is removed from database
		}
	}

	return nil
}

// GetGroupMembers retrieves all members of a group
func GetGroupMembers(groupID uint) ([]models.User, error) {
	group, err := GetGroupByID(groupID)
	if err != nil {
		return nil, err
	}

	return group.Members, nil
}

// GroupResponse represents a user group response
type GroupResponse struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	IsSystem    bool           `json:"isSystem"`
	MemberCount int            `json:"memberCount"`
	Members     []MemberInfo   `json:"members,omitempty"`
	CreatedAt   string         `json:"createdAt"`
	UpdatedAt   string         `json:"updatedAt"`
}

// MemberInfo represents basic user info for group members
type MemberInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}

// ToResponse converts a UserGroup to GroupResponse
func ToResponse(group *models.UserGroup) *GroupResponse {
	members := make([]MemberInfo, len(group.Members))
	for i, member := range group.Members {
		members[i] = MemberInfo{
			ID:       member.ID,
			Username: member.Username,
			FullName: member.FullName,
			Email:    member.Email,
		}
	}

	return &GroupResponse{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		IsSystem:    group.IsSystem,
		MemberCount: len(group.Members),
		Members:     members,
		CreatedAt:   group.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   group.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToResponses converts a slice of UserGroups to GroupResponses
func ToResponses(groups []*models.UserGroup) []*GroupResponse {
	responses := make([]*GroupResponse, len(groups))
	for i, group := range groups {
		responses[i] = ToResponse(group)
	}
	return responses
}

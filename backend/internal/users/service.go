package users

import (
	"github.com/stumpfworks/nas/internal/database"
	"github.com/stumpfworks/nas/pkg/errors"
	"gorm.io/gorm"
)

// AuthenticateUser authenticates a user with username/email and password
func AuthenticateUser(identifier, password string) (*User, error) {
	var user User

	// Try to find user by username or email
	err := database.DB.Where("username = ? OR email = ?", identifier, identifier).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Unauthorized("Invalid credentials", nil)
		}
		return nil, errors.InternalServerError("Failed to query user", err)
	}

	// Check password
	if !user.CheckPassword(password) {
		return nil, errors.Unauthorized("Invalid credentials", nil)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.Forbidden("User account is disabled", nil)
	}

	// Update last login
	if err := user.UpdateLastLogin(database.DB); err != nil {
		// Log error but don't fail authentication
		// logger.Error("Failed to update last login", zap.Error(err))
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(id uint) (*User, error) {
	var user User
	err := database.DB.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("User not found", err)
		}
		return nil, errors.InternalServerError("Failed to query user", err)
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("User not found", err)
		}
		return nil, errors.InternalServerError("Failed to query user", err)
	}
	return &user, nil
}

// ListUsers retrieves all users
func ListUsers() ([]*User, error) {
	var users []*User
	err := database.DB.Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, errors.InternalServerError("Failed to query users", err)
	}
	return users, nil
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"fullName"`
	Role     string `json:"role" validate:"required,oneof=admin user guest"`
}

// CreateUser creates a new user
func CreateUser(req *CreateUserRequest) (*User, error) {
	// Check if username already exists
	var existingUser User
	err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error
	if err == nil {
		return nil, errors.Conflict("Username already exists", nil)
	}

	// Check if email already exists
	err = database.DB.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return nil, errors.Conflict("Email already exists", nil)
	}

	// Create user
	user := &User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Role:     req.Role,
		IsActive: true,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, errors.InternalServerError("Failed to hash password", err)
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, errors.InternalServerError("Failed to create user", err)
	}

	return user, nil
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	FullName *string `json:"fullName,omitempty"`
	Role     *string `json:"role,omitempty" validate:"omitempty,oneof=admin user guest"`
	IsActive *bool   `json:"isActive,omitempty"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=8"`
}

// UpdateUser updates an existing user
func UpdateUser(id uint, req *UpdateUserRequest) (*User, error) {
	user, err := GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	updates := make(map[string]interface{})

	if req.Email != nil {
		// Check if email is already in use by another user
		var existingUser User
		err := database.DB.Where("email = ? AND id != ?", *req.Email, id).First(&existingUser).Error
		if err == nil {
			return nil, errors.Conflict("Email already exists", nil)
		}
		updates["email"] = *req.Email
	}

	if req.FullName != nil {
		updates["full_name"] = *req.FullName
	}

	if req.Role != nil {
		updates["role"] = *req.Role
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if req.Password != nil {
		if err := user.SetPassword(*req.Password); err != nil {
			return nil, errors.InternalServerError("Failed to hash password", err)
		}
		updates["password_hash"] = user.PasswordHash
	}

	// Perform update
	if err := database.DB.Model(user).Updates(updates).Error; err != nil {
		return nil, errors.InternalServerError("Failed to update user", err)
	}

	// Reload user
	if err := database.DB.First(user, id).Error; err != nil {
		return nil, errors.InternalServerError("Failed to reload user", err)
	}

	return user, nil
}

// DeleteUser deletes a user (soft delete)
func DeleteUser(id uint) error {
	user, err := GetUserByID(id)
	if err != nil {
		return err
	}

	// Prevent deleting the last admin
	if user.IsAdmin() {
		var adminCount int64
		database.DB.Model(&User{}).Where("role = ? AND is_active = ?", "admin", true).Count(&adminCount)
		if adminCount <= 1 {
			return errors.Forbidden("Cannot delete the last admin user", nil)
		}
	}

	if err := database.DB.Delete(user).Error; err != nil {
		return errors.InternalServerError("Failed to delete user", err)
	}

	return nil
}

// ChangePassword changes a user's password
func ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !user.CheckPassword(oldPassword) {
		return errors.BadRequest("Invalid old password", nil)
	}

	// Set new password
	if err := user.SetPassword(newPassword); err != nil {
		return errors.InternalServerError("Failed to hash password", err)
	}

	// Save
	if err := database.DB.Model(user).Update("password_hash", user.PasswordHash).Error; err != nil {
		return errors.InternalServerError("Failed to update password", err)
	}

	return nil
}

// UserResponse represents a user response (without sensitive data)
type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"fullName"`
	Role      string `json:"role"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ToResponse converts a User to UserResponse
func ToResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToResponses converts a slice of Users to UserResponses
func ToResponses(users []*User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToResponse(user)
	}
	return responses
}

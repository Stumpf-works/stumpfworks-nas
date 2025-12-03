// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package testutil

import (
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
)

// CreateTestUser creates a test user
func CreateTestUser(username, role string) *models.User {
	return &models.User{
		Username:  username,
		Role:      role,
		Email:     username + "@test.local",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestAdminUser creates a test admin user
func CreateTestAdminUser() *models.User {
	return CreateTestUser("admin", "admin")
}

// CreateTestRegularUser creates a test regular user
func CreateTestRegularUser() *models.User {
	return CreateTestUser("user", "user")
}

// CreateTestShare creates a test share
func CreateTestShare(name, path string) *models.Share {
	return &models.Share{
		Name:       name,
		Path:       path,
		GuestOK:    false,
		ReadOnly:   false,
		ValidUsers: "",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// CreateTestDockerContainer creates a test Docker container model
func CreateTestDockerContainer(name, image string) *models.DockerContainer {
	return &models.DockerContainer{
		Name:      name,
		Image:     image,
		Status:    "running",
		CreatedAt: time.Now(),
	}
}

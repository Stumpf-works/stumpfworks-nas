package models

import (
	"time"

	"gorm.io/gorm"
)

// UserGroup represents a group of users
type UserGroup struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	IsSystem    bool   `gorm:"default:false" json:"isSystem"` // System groups can't be deleted

	// Many-to-many relationship with users
	Members []User `gorm:"many2many:user_group_members;" json:"members,omitempty"`
}

// TableName specifies the table name for UserGroup model
func (UserGroup) TableName() string {
	return "user_groups"
}

// UnixGroupName returns the Unix-safe group name (lowercase, no spaces)
func (g *UserGroup) UnixGroupName() string {
	// Convert to lowercase and replace spaces with underscores
	// This ensures compatibility with Unix group names
	name := g.Name
	result := ""
	for _, r := range name {
		if r == ' ' {
			result += "_"
		} else if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			result += string(r)
		} else if r >= 'A' && r <= 'Z' {
			result += string(r + 32) // Convert to lowercase
		}
	}
	return result
}

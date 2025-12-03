// Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
package testutil

import (
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return db
}

// SetupTestDBWithModels creates a test database and auto-migrates the given models
func SetupTestDBWithModels(t *testing.T, models ...interface{}) *gorm.DB {
	db := SetupTestDB(t)
	if err := db.AutoMigrate(models...); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	return db
}

// CleanupTestDB cleans up the test database
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		t.Logf("Warning: Failed to get SQL DB for cleanup: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		t.Logf("Warning: Failed to close test database: %v", err)
	}
}

// TruncateTable truncates a table in the test database
func TruncateTable(t *testing.T, db *gorm.DB, tableName string) {
	if err := db.Exec(fmt.Sprintf("DELETE FROM %s", tableName)).Error; err != nil {
		t.Fatalf("Failed to truncate table %s: %v", tableName, err)
	}
}

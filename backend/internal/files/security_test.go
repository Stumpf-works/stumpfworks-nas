// Revision: 2025-11-17 | Author: Claude | Version: 1.0.0
package files

import (
	"strings"
	"testing"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
)

// TestPathTraversal tests protection against path traversal attacks
func TestPathTraversal(t *testing.T) {
	tests := []struct {
		name        string
		inputPath   string
		basePaths   []string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Simple path traversal with .. (cleaned to outside path)",
			inputPath:   "/mnt/storage/../../../etc/passwd",
			basePaths:   []string{"/mnt/storage"},
			shouldError: true,
			errorMsg:    "outside allowed share locations",
		},
		{
			name:        "URL encoded path traversal",
			inputPath:   "/mnt/storage/%2e%2e/%2e%2e/etc/passwd",
			basePaths:   []string{"/mnt/storage"},
			shouldError: false, // Should be cleaned by filepath.Clean
		},
		{
			name:        "Valid absolute path within base",
			inputPath:   "/mnt/storage/documents/test.txt",
			basePaths:   []string{"/mnt/storage"},
			shouldError: false,
		},
		{
			name:        "Valid path at base directory",
			inputPath:   "/mnt/storage",
			basePaths:   []string{"/mnt/storage"},
			shouldError: false,
		},
		{
			name:        "Path outside allowed base",
			inputPath:   "/var/log/system.log",
			basePaths:   []string{"/mnt/storage"},
			shouldError: true,
			errorMsg:    "outside allowed share locations",
		},
		{
			name:        "Relative path with .. (detected as traversal)",
			inputPath:   "../../etc/passwd",
			basePaths:   []string{"/mnt/storage"},
			shouldError: true,
			errorMsg:    "path traversal detected",
		},
		{
			name:        "Multiple basePaths - first matches",
			inputPath:   "/mnt/pool1/data/file.txt",
			basePaths:   []string{"/mnt/pool1", "/mnt/pool2"},
			shouldError: false,
		},
		{
			name:        "Multiple basePaths - second matches",
			inputPath:   "/mnt/pool2/backup/archive.zip",
			basePaths:   []string{"/mnt/pool1", "/mnt/pool2"},
			shouldError: false,
		},
		{
			name:        "Multiple basePaths - none match",
			inputPath:   "/tmp/malicious/file.sh",
			basePaths:   []string{"/mnt/pool1", "/mnt/pool2"},
			shouldError: true,
			errorMsg:    "outside allowed share locations",
		},
		{
			name:        "Empty basePaths allows any absolute path",
			inputPath:   "/tmp/test.txt",
			basePaths:   []string{},
			shouldError: false,
		},
		{
			name:        "Path with .. cleaned but still within base (allowed)",
			inputPath:   "/mnt/storage/subdir/../file.txt",
			basePaths:   []string{"/mnt/storage"},
			shouldError: false, // Cleaned to /mnt/storage/file.txt - allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewPathValidator(tt.basePaths)
			result, err := validator.ValidateAndSanitize(tt.inputPath)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for path '%s', but got none. Result: %s", tt.inputPath, result)
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', but got: %s", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for path '%s', but got: %v", tt.inputPath, err)
				}
			}
		})
	}
}

// TestValidateFileName tests filename validation
func TestValidateFileName(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Valid filename",
			filename:    "document.txt",
			shouldError: false,
		},
		{
			name:        "Valid filename with spaces",
			filename:    "my document.txt",
			shouldError: false,
		},
		{
			name:        "Valid filename with dashes and underscores",
			filename:    "test-file_123.pdf",
			shouldError: false,
		},
		{
			name:        "Empty filename",
			filename:    "",
			shouldError: true,
			errorMsg:    "cannot be empty",
		},
		{
			name:        "Filename with forward slash",
			filename:    "path/to/file.txt",
			shouldError: true,
			errorMsg:    "invalid character",
		},
		{
			name:        "Filename with backslash",
			filename:    "file\\name.txt",
			shouldError: true,
			errorMsg:    "invalid character",
		},
		{
			name:        "Filename with null byte",
			filename:    "file\x00name.txt",
			shouldError: true,
			errorMsg:    "invalid character",
		},
		{
			name:        "Reserved Windows name - CON",
			filename:    "CON",
			shouldError: true,
			errorMsg:    "reserved name",
		},
		{
			name:        "Reserved Windows name - CON.txt",
			filename:    "CON.txt",
			shouldError: true,
			errorMsg:    "reserved name",
		},
		{
			name:        "Reserved Windows name - PRN",
			filename:    "PRN",
			shouldError: true,
			errorMsg:    "reserved name",
		},
		{
			name:        "Filename starting with dot",
			filename:    ".hidden",
			shouldError: true,
			errorMsg:    "cannot start or end with dots or spaces",
		},
		{
			name:        "Filename ending with dot",
			filename:    "file.",
			shouldError: true,
			errorMsg:    "cannot start or end with dots or spaces",
		},
		{
			name:        "Filename starting with space",
			filename:    " file.txt",
			shouldError: true,
			errorMsg:    "cannot start or end with dots or spaces",
		},
		{
			name:        "Filename ending with space",
			filename:    "file.txt ",
			shouldError: true,
			errorMsg:    "cannot start or end with dots or spaces",
		},
		{
			name:        "Filename with wildcard *",
			filename:    "file*.txt",
			shouldError: true,
			errorMsg:    "invalid character",
		},
		{
			name:        "Filename with question mark",
			filename:    "file?.txt",
			shouldError: true,
			errorMsg:    "invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileName(tt.filename)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for filename '%s', but got none", tt.filename)
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', but got: %s", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for filename '%s', but got: %v", tt.filename, err)
				}
			}
		})
	}
}

// TestPermissionChecker tests permission checking logic
func TestPermissionChecker(t *testing.T) {
	// Setup test shares
	shares := []*models.Share{
		{
			Name:       "public",
			Path:       "/mnt/storage/public",
			GuestOK:    true,
			ReadOnly:   false,
			ValidUsers: "",
		},
		{
			Name:       "private",
			Path:       "/mnt/storage/private",
			GuestOK:    false,
			ReadOnly:   false,
			ValidUsers: "alice,bob",
		},
		{
			Name:       "readonly",
			Path:       "/mnt/storage/readonly",
			GuestOK:    true,
			ReadOnly:   true,
			ValidUsers: "",
		},
	}

	checker := NewPermissionChecker(shares)

	tests := []struct {
		name        string
		user        *models.User
		isAdmin     bool
		path        string
		operation   string // "read", "write", "delete"
		shouldError bool
	}{
		{
			name:        "Admin can access any path",
			user:        &models.User{Username: "admin", Role: "admin"},
			isAdmin:     true,
			path:        "/mnt/storage/anywhere/file.txt",
			operation:   "read",
			shouldError: false,
		},
		{
			name:        "Admin can write anywhere",
			user:        &models.User{Username: "admin", Role: "admin"},
			isAdmin:     true,
			path:        "/mnt/storage/private/admin-file.txt",
			operation:   "write",
			shouldError: false,
		},
		{
			name:        "User can read allowed path",
			user:        &models.User{Username: "alice", Role: "user"},
			isAdmin:     false,
			path:        "/mnt/storage/private/alice-file.txt",
			operation:   "read",
			shouldError: false,
		},
		{
			name:        "User cannot read disallowed path",
			user:        &models.User{Username: "charlie", Role: "user"},
			isAdmin:     false,
			path:        "/mnt/storage/private/secret.txt",
			operation:   "read",
			shouldError: true,
		},
		{
			name:        "User can write to writable share",
			user:        &models.User{Username: "alice", Role: "user"},
			isAdmin:     false,
			path:        "/mnt/storage/private/new-file.txt",
			operation:   "write",
			shouldError: false,
		},
		{
			name:        "User cannot write to readonly share",
			user:        &models.User{Username: "alice", Role: "user"},
			isAdmin:     false,
			path:        "/mnt/storage/readonly/attempt.txt",
			operation:   "write",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowedPaths := GetAllowedPathsForUser(tt.user, shares)
			ctx := &SecurityContext{
				User:         tt.user,
				IsAdmin:      tt.isAdmin,
				AllowedPaths: allowedPaths,
			}

			var err error
			switch tt.operation {
			case "read":
				err = checker.CanAccess(ctx, tt.path)
			case "write":
				err = checker.CanWrite(ctx, tt.path)
			case "delete":
				err = checker.CanDelete(ctx, tt.path)
			}

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for user '%s' %s '%s', but got none",
						tt.user.Username, tt.operation, tt.path)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for user '%s' %s '%s', but got: %v",
						tt.user.Username, tt.operation, tt.path, err)
				}
			}
		})
	}
}

// BenchmarkPathTraversal benchmarks path validation performance
func BenchmarkPathTraversal(b *testing.B) {
	validator := NewPathValidator([]string{"/mnt/storage"})
	testPath := "/mnt/storage/documents/test.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = validator.ValidateAndSanitize(testPath)
	}
}

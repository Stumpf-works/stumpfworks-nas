package files

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
)

// SecurityContext holds security-related information for file operations
type SecurityContext struct {
	User        *models.User
	IsAdmin     bool
	AllowedPaths []string // Paths the user is allowed to access
}

// PathValidator validates and sanitizes file paths
type PathValidator struct {
	basePaths []string // Allowed base paths (e.g., mounted volumes)
}

// NewPathValidator creates a new path validator
func NewPathValidator(basePaths []string) *PathValidator {
	return &PathValidator{
		basePaths: basePaths,
	}
}

// ValidateAndSanitize validates and sanitizes a file path
// Returns the absolute, cleaned path or an error if the path is invalid
func (pv *PathValidator) ValidateAndSanitize(requestPath string) (string, error) {
	// Clean the path (removes .., ., etc.)
	cleanPath := filepath.Clean(requestPath)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return "", errors.BadRequest("Invalid path: path traversal detected (security violation)", nil)
	}

	// If not absolute, it's relative - we need a base path
	if !filepath.IsAbs(cleanPath) {
		return "", errors.BadRequest(
			fmt.Sprintf("Invalid path '%s': must be absolute (e.g., /mnt/storage/files)", requestPath),
			nil)
	}

	// Validate against allowed base paths
	if len(pv.basePaths) > 0 {
		allowed := false
		for _, basePath := range pv.basePaths {
			if strings.HasPrefix(cleanPath, basePath) {
				allowed = true
				break
			}
		}
		if !allowed {
			allowedPaths := strings.Join(pv.basePaths, ", ")
			return "", errors.Forbidden(
				fmt.Sprintf("Access denied: Path '%s' is outside allowed share locations. "+
					"Allowed: %s",
					cleanPath, allowedPaths),
				nil)
		}
	}

	return cleanPath, nil
}

// ValidatePaths validates multiple paths
func (pv *PathValidator) ValidatePaths(paths []string) ([]string, error) {
	validated := make([]string, 0, len(paths))
	for _, path := range paths {
		cleanPath, err := pv.ValidateAndSanitize(path)
		if err != nil {
			return nil, err
		}
		validated = append(validated, cleanPath)
	}
	return validated, nil
}

// PermissionChecker checks file permissions for users
type PermissionChecker struct {
	shares map[string]*models.Share // Cache of shares
}

// NewPermissionChecker creates a new permission checker
func NewPermissionChecker(shares []*models.Share) *PermissionChecker {
	shareMap := make(map[string]*models.Share)
	for _, share := range shares {
		shareMap[share.Path] = share
	}
	return &PermissionChecker{
		shares: shareMap,
	}
}

// CanAccess checks if a user can access a given path
func (pc *PermissionChecker) CanAccess(ctx *SecurityContext, path string) error {
	// Admins can access everything
	if ctx.IsAdmin {
		return nil
	}

	// Check if the path is within user's allowed paths
	if len(ctx.AllowedPaths) == 0 {
		return errors.Forbidden(
			fmt.Sprintf("Access denied: No shares configured for user '%s'. "+
				"Contact your administrator to grant access to shares or create new shares in the Storage app.",
				ctx.User.Username),
			nil)
	}

	for _, allowedPath := range ctx.AllowedPaths {
		if strings.HasPrefix(path, allowedPath) {
			return nil
		}
	}

	// Build helpful error message with available shares
	availableShares := strings.Join(ctx.AllowedPaths, ", ")
	return errors.Forbidden(
		fmt.Sprintf("Access denied: Path '%s' is not accessible for user '%s'. "+
			"Available shares: %s",
			path, ctx.User.Username, availableShares),
		nil)
}

// CanWrite checks if a user can write to a given path
func (pc *PermissionChecker) CanWrite(ctx *SecurityContext, path string) error {
	// Admins can write everywhere
	if ctx.IsAdmin {
		return nil
	}

	// First check if user can access the path
	if err := pc.CanAccess(ctx, path); err != nil {
		return err
	}

	// Check if the share is read-only
	for _, share := range pc.shares {
		if strings.HasPrefix(path, share.Path) {
			if share.ReadOnly {
				return errors.Forbidden("Access denied: share is read-only", nil)
			}
			return nil
		}
	}

	// If no share found but path is accessible, allow write
	return nil
}

// CanDelete checks if a user can delete from a given path
func (pc *PermissionChecker) CanDelete(ctx *SecurityContext, path string) error {
	// Same as CanWrite for now
	return pc.CanWrite(ctx, path)
}

// CanChangePermissions checks if a user can change permissions
func (pc *PermissionChecker) CanChangePermissions(ctx *SecurityContext, path string) error {
	// Only admins can change permissions
	if !ctx.IsAdmin {
		return errors.Forbidden("Access denied: only administrators can change permissions", nil)
	}
	return nil
}

// GetAllowedPathsForUser returns all paths a user can access based on shares
func GetAllowedPathsForUser(user *models.User, shares []*models.Share) []string {
	if user.Role == "admin" {
		// Admins can access all share paths
		paths := make([]string, len(shares))
		for i, share := range shares {
			paths[i] = share.Path
		}
		return paths
	}

	// Regular users: check ValidUsers field
	allowedPaths := []string{}
	for _, share := range shares {
		if share.GuestOK {
			allowedPaths = append(allowedPaths, share.Path)
			continue
		}

		// Check if user is in ValidUsers list
		if share.ValidUsers != "" {
			validUsers := strings.Split(share.ValidUsers, ",")
			for _, validUser := range validUsers {
				if strings.TrimSpace(validUser) == user.Username {
					allowedPaths = append(allowedPaths, share.Path)
					break
				}
			}
		}
	}

	return allowedPaths
}

// ValidateFileName checks if a filename is valid
func ValidateFileName(name string) error {
	if name == "" {
		return errors.BadRequest("Filename cannot be empty", nil)
	}

	// Check for invalid characters
	invalidChars := []string{"/", "\\", "\x00", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return errors.BadRequest(fmt.Sprintf("Filename contains invalid character: %s", char), nil)
		}
	}

	// Check for reserved names (Windows compatibility)
	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	upperName := strings.ToUpper(name)
	for _, reserved := range reservedNames {
		if upperName == reserved || strings.HasPrefix(upperName, reserved+".") {
			return errors.BadRequest(fmt.Sprintf("Filename uses reserved name: %s", reserved), nil)
		}
	}

	// Check for names starting/ending with dots or spaces
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, " ") || strings.HasSuffix(name, " ") || strings.HasSuffix(name, ".") {
		return errors.BadRequest("Filename cannot start or end with dots or spaces", nil)
	}

	return nil
}

// CheckDiskSpace checks if there's enough disk space for an operation
func CheckDiskSpace(path string, requiredBytes int64) error {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return errors.InternalServerError("Failed to check disk space", err)
	}

	// Available space = free blocks * block size
	availableSpace := int64(stat.Bavail) * int64(stat.Bsize)

	if availableSpace < requiredBytes {
		return errors.NewAppError(507, // HTTP 507 Insufficient Storage
			fmt.Sprintf("Insufficient disk space: required %d bytes, available %d bytes", requiredBytes, availableSpace),
			nil)
	}

	return nil
}

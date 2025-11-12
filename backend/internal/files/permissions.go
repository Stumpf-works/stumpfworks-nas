package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// ChangePermissions changes file or directory permissions
func (s *Service) ChangePermissions(ctx *SecurityContext, req *PermissionsRequest) error {
	// Validate path
	cleanPath, err := s.validator.ValidateAndSanitize(req.Path)
	if err != nil {
		return err
	}

	// Check if user can change permissions (admin only)
	if err := s.permissions.CanChangePermissions(ctx, cleanPath); err != nil {
		return err
	}

	// Check if path exists
	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.NotFound("Path not found", err)
		}
		return errors.InternalServerError("Failed to access path", err)
	}

	// Parse permissions
	permInt, err := strconv.ParseUint(req.Permissions, 8, 32)
	if err != nil {
		return errors.BadRequest("Invalid permissions format (use octal, e.g., 0644)", err)
	}
	perm := os.FileMode(permInt)

	// Change permissions
	if req.Recursive && info.IsDir() {
		err = s.chmodRecursive(cleanPath, perm)
	} else {
		err = os.Chmod(cleanPath, perm)
	}

	if err != nil {
		logger.Error("Failed to change permissions", zap.String("path", cleanPath), zap.Error(err))
		return errors.InternalServerError("Failed to change permissions", err)
	}

	// Change owner/group if specified (Unix only)
	if req.Owner != "" || req.Group != "" {
		if err := s.changeOwnership(cleanPath, req.Owner, req.Group, req.Recursive); err != nil {
			return err
		}
	}

	logger.Info("Permissions changed", zap.String("path", cleanPath), zap.String("permissions", req.Permissions), zap.String("user", ctx.User.Username))
	return nil
}

// GetPermissions returns detailed permissions for a file
func (s *Service) GetPermissions(ctx *SecurityContext, path string) (*PermissionsInfo, error) {
	// Validate path
	cleanPath, err := s.validator.ValidateAndSanitize(path)
	if err != nil {
		return nil, err
	}

	// Check access permissions
	if err := s.permissions.CanAccess(ctx, cleanPath); err != nil {
		return nil, err
	}

	// Get file info
	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.NotFound("Path not found", err)
		}
		return nil, errors.InternalServerError("Failed to access path", err)
	}

	permInfo := &PermissionsInfo{
		Path:        cleanPath,
		Permissions: fmt.Sprintf("%04o", info.Mode().Perm()),
		Mode:        info.Mode().String(),
		Owner:       "system",
		Group:       "system",
		UID:         0,
		GID:         0,
	}

	// Get Unix owner/group information
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		permInfo.UID = int(stat.Uid)
		permInfo.GID = int(stat.Gid)
		// TODO: Lookup user/group names from UID/GID
	}

	return permInfo, nil
}

// GetDiskUsage returns disk usage information for a path
func (s *Service) GetDiskUsage(ctx *SecurityContext, path string) (*DiskUsageInfo, error) {
	// Validate path
	cleanPath, err := s.validator.ValidateAndSanitize(path)
	if err != nil {
		return nil, err
	}

	// Check access permissions
	if err := s.permissions.CanAccess(ctx, cleanPath); err != nil {
		return nil, err
	}

	// Get filesystem stats
	var stat syscall.Statfs_t
	if err := syscall.Statfs(cleanPath, &stat); err != nil {
		return nil, errors.InternalServerError("Failed to get disk usage", err)
	}

	// Calculate sizes
	totalSize := int64(stat.Blocks) * int64(stat.Bsize)
	freeSize := int64(stat.Bfree) * int64(stat.Bsize)
	usedSize := totalSize - freeSize
	usagePercent := float64(usedSize) / float64(totalSize) * 100

	return &DiskUsageInfo{
		Path:         cleanPath,
		TotalSize:    totalSize,
		UsedSize:     usedSize,
		FreeSize:     freeSize,
		UsagePercent: usagePercent,
	}, nil
}

// Helper: chmodRecursive changes permissions recursively
func (s *Service) chmodRecursive(path string, perm os.FileMode) error {
	return filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Chmod(walkPath, perm)
	})
}

// Helper: changeOwnership changes file owner and group
func (s *Service) changeOwnership(path, owner, group string, recursive bool) error {
	// Parse UID/GID
	var uid, gid int = -1, -1 // -1 means no change

	if owner != "" {
		// TODO: Lookup UID from username
		// For now, expect numeric UID
		parsedUID, err := strconv.Atoi(owner)
		if err != nil {
			return errors.BadRequest("Invalid owner format (numeric UID expected)", err)
		}
		uid = parsedUID
	}

	if group != "" {
		// TODO: Lookup GID from group name
		// For now, expect numeric GID
		parsedGID, err := strconv.Atoi(group)
		if err != nil {
			return errors.BadRequest("Invalid group format (numeric GID expected)", err)
		}
		gid = parsedGID
	}

	// Change ownership
	if recursive {
		return filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			return os.Chown(walkPath, uid, gid)
		})
	}

	return os.Chown(path, uid, gid)
}

// PermissionsInfo holds detailed permissions information
type PermissionsInfo struct {
	Path        string `json:"path"`
	Permissions string `json:"permissions"` // Octal format (e.g., "0644")
	Mode        string `json:"mode"`        // String format (e.g., "-rw-r--r--")
	Owner       string `json:"owner"`
	Group       string `json:"group"`
	UID         int    `json:"uid"`
	GID         int    `json:"gid"`
}

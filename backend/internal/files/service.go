package files

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// Service provides file management operations
type Service struct {
	validator   *PathValidator
	permissions *PermissionChecker
}

// NewService creates a new file service
func NewService(allowedPaths []string, permChecker *PermissionChecker) *Service {
	return &Service{
		validator:   NewPathValidator(allowedPaths),
		permissions: permChecker,
	}
}

// CheckWritePermission validates that a user has write access to a path
// This is a helper method for operations that need to check permissions before acting
func (s *Service) CheckWritePermission(ctx *SecurityContext, path string) error {
	// Validate and sanitize path first
	cleanPath, err := s.validator.ValidateAndSanitize(path)
	if err != nil {
		return err
	}

	// Check write permissions
	return s.permissions.CanWrite(ctx, cleanPath)
}

// Browse lists files and directories at the specified path
func (s *Service) Browse(ctx *SecurityContext, req *BrowseRequest) (*BrowseResponse, error) {
	// Validate and sanitize path
	cleanPath, err := s.validator.ValidateAndSanitize(req.Path)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if err := s.permissions.CanAccess(ctx, cleanPath); err != nil {
		return nil, err
	}

	// Check if path exists
	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.NotFound("Path not found", err)
		}
		return nil, errors.InternalServerError("Failed to access path", err)
	}

	// Must be a directory
	if !info.IsDir() {
		return nil, errors.BadRequest("Path is not a directory", nil)
	}

	// Read directory contents
	entries, err := os.ReadDir(cleanPath)
	if err != nil {
		logger.Error("Failed to read directory", zap.String("path", cleanPath), zap.Error(err))
		return nil, errors.InternalServerError("Failed to read directory", err)
	}

	// Build response
	files := make([]FileInfo, 0, len(entries))
	var totalSize int64
	var totalFiles, totalDirs int

	for _, entry := range entries {
		// Skip hidden files unless requested
		if !req.ShowHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		entryPath := filepath.Join(cleanPath, entry.Name())
		fileInfo, err := s.getFileInfo(entryPath, entry)
		if err != nil {
			logger.Warn("Failed to get file info", zap.String("path", entryPath), zap.Error(err))
			continue
		}

		files = append(files, *fileInfo)
		if fileInfo.IsDir {
			totalDirs++
		} else {
			totalFiles++
			totalSize += fileInfo.Size
		}
	}

	return &BrowseResponse{
		Path:       cleanPath,
		Files:      files,
		TotalSize:  totalSize,
		TotalFiles: totalFiles,
		TotalDirs:  totalDirs,
	}, nil
}

// GetFileInfo returns information about a specific file or directory
func (s *Service) GetFileInfo(ctx *SecurityContext, path string) (*FileInfo, error) {
	// Validate and sanitize path
	cleanPath, err := s.validator.ValidateAndSanitize(path)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if err := s.permissions.CanAccess(ctx, cleanPath); err != nil {
		return nil, err
	}

	// Get file info and check if exists
	_, err = os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.NotFound("File not found", err)
		}
		return nil, errors.InternalServerError("Failed to access file", err)
	}

	return s.getFileInfo(cleanPath, nil)
}

// CreateDirectory creates a new directory
func (s *Service) CreateDirectory(ctx *SecurityContext, req *CreateDirRequest) error {
	// Validate filename
	if err := ValidateFileName(req.Name); err != nil {
		return err
	}

	// Validate and sanitize parent path
	parentPath, err := s.validator.ValidateAndSanitize(req.Path)
	if err != nil {
		return err
	}

	// Check write permissions
	if err := s.permissions.CanWrite(ctx, parentPath); err != nil {
		return err
	}

	// Build full path
	fullPath := filepath.Join(parentPath, req.Name)

	// Check if already exists
	if _, err := os.Stat(fullPath); err == nil {
		return errors.Conflict("Directory already exists", nil)
	}

	// Parse permissions (default 0755)
	perm := os.FileMode(0755)
	if req.Permissions != "" {
		var permInt uint32
		if _, err := fmt.Sscanf(req.Permissions, "%o", &permInt); err == nil {
			perm = os.FileMode(permInt)
		}
	}

	// Create directory
	if err := os.Mkdir(fullPath, perm); err != nil {
		logger.Error("Failed to create directory", zap.String("path", fullPath), zap.Error(err))
		return errors.InternalServerError("Failed to create directory", err)
	}

	logger.Info("Directory created", zap.String("path", fullPath), zap.String("user", ctx.User.Username))
	return nil
}

// Delete deletes files or directories
func (s *Service) Delete(ctx *SecurityContext, req *DeleteRequest) error {
	// Validate all paths
	cleanPaths, err := s.validator.ValidatePaths(req.Paths)
	if err != nil {
		return err
	}

	// Check delete permissions for all paths
	for _, path := range cleanPaths {
		if err := s.permissions.CanDelete(ctx, path); err != nil {
			return err
		}
	}

	// Delete each path
	for _, path := range cleanPaths {
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue // Already deleted or doesn't exist
			}
			return errors.InternalServerError("Failed to access path", err)
		}

		if info.IsDir() && !req.Recursive {
			// Check if directory is empty
			entries, err := os.ReadDir(path)
			if err != nil {
				return errors.InternalServerError("Failed to read directory", err)
			}
			if len(entries) > 0 {
				return errors.BadRequest("Directory is not empty, use recursive delete", nil)
			}
		}

		// Perform deletion
		if req.Recursive && info.IsDir() {
			err = os.RemoveAll(path)
		} else {
			err = os.Remove(path)
		}

		if err != nil {
			logger.Error("Failed to delete path", zap.String("path", path), zap.Error(err))
			return errors.InternalServerError(fmt.Sprintf("Failed to delete: %s", filepath.Base(path)), err)
		}

		logger.Info("Path deleted", zap.String("path", path), zap.String("user", ctx.User.Username))
	}

	return nil
}

// Rename renames a file or directory
func (s *Service) Rename(ctx *SecurityContext, req *RenameRequest) error {
	// Validate new name
	if err := ValidateFileName(req.NewName); err != nil {
		return err
	}

	// Validate and sanitize old path
	oldPath, err := s.validator.ValidateAndSanitize(req.OldPath)
	if err != nil {
		return err
	}

	// Check write permissions
	if err := s.permissions.CanWrite(ctx, oldPath); err != nil {
		return err
	}

	// Build new path (same directory)
	newPath := filepath.Join(filepath.Dir(oldPath), req.NewName)

	// Check if target already exists
	if _, err := os.Stat(newPath); err == nil {
		return errors.Conflict("Target file already exists", nil)
	}

	// Perform rename
	if err := os.Rename(oldPath, newPath); err != nil {
		logger.Error("Failed to rename", zap.String("old", oldPath), zap.String("new", newPath), zap.Error(err))
		return errors.InternalServerError("Failed to rename", err)
	}

	logger.Info("File renamed", zap.String("old", oldPath), zap.String("new", newPath), zap.String("user", ctx.User.Username))
	return nil
}

// Copy copies a file or directory
func (s *Service) Copy(ctx *SecurityContext, req *CopyMoveRequest) error {
	// Validate paths
	srcPath, err := s.validator.ValidateAndSanitize(req.Source)
	if err != nil {
		return err
	}

	dstPath, err := s.validator.ValidateAndSanitize(req.Destination)
	if err != nil {
		return err
	}

	// Check permissions
	if err := s.permissions.CanAccess(ctx, srcPath); err != nil {
		return err
	}
	if err := s.permissions.CanWrite(ctx, dstPath); err != nil {
		return err
	}

	// Check if source exists
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.NotFound("Source not found", err)
		}
		return errors.InternalServerError("Failed to access source", err)
	}

	// Check if destination exists
	if _, err := os.Stat(dstPath); err == nil && !req.Overwrite {
		return errors.Conflict("Destination already exists", nil)
	}

	// Perform copy
	if srcInfo.IsDir() {
		err = s.copyDirectory(srcPath, dstPath)
	} else {
		err = s.copyFile(srcPath, dstPath)
	}

	if err != nil {
		return err
	}

	logger.Info("Path copied", zap.String("src", srcPath), zap.String("dst", dstPath), zap.String("user", ctx.User.Username))
	return nil
}

// Move moves a file or directory
func (s *Service) Move(ctx *SecurityContext, req *CopyMoveRequest) error {
	// Validate paths
	srcPath, err := s.validator.ValidateAndSanitize(req.Source)
	if err != nil {
		return err
	}

	dstPath, err := s.validator.ValidateAndSanitize(req.Destination)
	if err != nil {
		return err
	}

	// Check permissions
	if err := s.permissions.CanWrite(ctx, srcPath); err != nil {
		return err
	}
	if err := s.permissions.CanWrite(ctx, dstPath); err != nil {
		return err
	}

	// Check if destination exists
	if _, err := os.Stat(dstPath); err == nil && !req.Overwrite {
		return errors.Conflict("Destination already exists", nil)
	}

	// Try direct rename first (same filesystem)
	if err := os.Rename(srcPath, dstPath); err == nil {
		logger.Info("Path moved", zap.String("src", srcPath), zap.String("dst", dstPath), zap.String("user", ctx.User.Username))
		return nil
	}

	// If rename fails, copy then delete (cross-filesystem move)
	if err := s.Copy(ctx, req); err != nil {
		return err
	}

	if err := s.Delete(ctx, &DeleteRequest{Paths: []string{srcPath}, Recursive: true}); err != nil {
		logger.Error("Failed to delete source after copy", zap.String("path", srcPath), zap.Error(err))
		return errors.InternalServerError("Move partially completed, failed to delete source", err)
	}

	logger.Info("Path moved (cross-filesystem)", zap.String("src", srcPath), zap.String("dst", dstPath), zap.String("user", ctx.User.Username))
	return nil
}

// Helper: getFileInfo extracts file information
func (s *Service) getFileInfo(path string, entry os.DirEntry) (*FileInfo, error) {
	var info os.FileInfo
	var err error

	if entry != nil {
		info, err = entry.Info()
	} else {
		info, err = os.Stat(path)
	}

	if err != nil {
		return nil, err
	}

	// Get file extension and MIME type
	ext := strings.ToLower(filepath.Ext(path))
	mimeType := mime.TypeByExtension(ext)

	// Determine if file can have thumbnail
	hasThumbnail := false
	if strings.HasPrefix(mimeType, "image/") {
		hasThumbnail = true
	}

	fileInfo := &FileInfo{
		Name:         info.Name(),
		Path:         path,
		Size:         info.Size(),
		IsDir:        info.IsDir(),
		ModTime:      info.ModTime(),
		Permissions:  info.Mode().Perm().String(),
		Extension:    ext,
		MimeType:     mimeType,
		HasThumbnail: hasThumbnail,
	}

	// Get owner/group (Unix-specific, requires syscall)
	// TODO: Implement owner/group extraction using syscall
	fileInfo.Owner = "system"
	fileInfo.Group = "system"

	return fileInfo, nil
}

// Helper: copyFile copies a single file
func (s *Service) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return errors.InternalServerError("Failed to open source file", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return errors.InternalServerError("Failed to create destination file", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return errors.InternalServerError("Failed to copy file data", err)
	}

	// Copy permissions
	srcInfo, err := os.Stat(src)
	if err == nil {
		os.Chmod(dst, srcInfo.Mode())
	}

	return nil
}

// Helper: copyDirectory recursively copies a directory
func (s *Service) copyDirectory(src, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return errors.InternalServerError("Failed to access source directory", err)
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return errors.InternalServerError("Failed to create destination directory", err)
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return errors.InternalServerError("Failed to read source directory", err)
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := s.copyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := s.copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

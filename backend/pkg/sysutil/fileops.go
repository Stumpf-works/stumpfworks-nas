// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package sysutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies a file from src to dst
// Preserves file permissions
func CopyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Sync to ensure data is written
	if err := dstFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

// MoveFile moves a file from src to dst
// Tries rename first, falls back to copy+delete if across filesystems
func MoveFile(src, dst string) error {
	// Try rename first (fast if same filesystem)
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// Rename failed (probably cross-filesystem), do copy+delete
	if err := CopyFile(src, dst); err != nil {
		return fmt.Errorf("failed to copy file during move: %w", err)
	}

	// Remove source after successful copy
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("failed to remove source file after copy: %w", err)
	}

	return nil
}

// CopyDir recursively copies a directory from src to dst
// Preserves directory structure and file permissions
func CopyDir(src, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source directory: %w", err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory: %s", src)
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Read source directory entries
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// MoveDir moves a directory from src to dst
// Tries rename first, falls back to copy+delete if across filesystems
func MoveDir(src, dst string) error {
	// Try rename first (fast if same filesystem)
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// Rename failed (probably cross-filesystem), do copy+delete
	if err := CopyDir(src, dst); err != nil {
		return fmt.Errorf("failed to copy directory during move: %w", err)
	}

	// Remove source after successful copy
	if err := os.RemoveAll(src); err != nil {
		return fmt.Errorf("failed to remove source directory after copy: %w", err)
	}

	return nil
}

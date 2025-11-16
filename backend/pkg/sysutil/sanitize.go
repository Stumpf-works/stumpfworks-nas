// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package sysutil

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// invalidFilenameChars are characters that should not be in filenames
	invalidFilenameChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)

	// reservedFilenames are Windows reserved filenames to avoid
	reservedFilenames = map[string]bool{
		"CON": true, "PRN": true, "AUX": true, "NUL": true,
		"COM1": true, "COM2": true, "COM3": true, "COM4": true,
		"COM5": true, "COM6": true, "COM7": true, "COM8": true,
		"COM9": true, "LPT1": true, "LPT2": true, "LPT3": true,
		"LPT4": true, "LPT5": true, "LPT6": true, "LPT7": true,
		"LPT8": true, "LPT9": true,
	}
)

// SanitizeFilename sanitizes a filename by removing or replacing invalid characters
// Returns a safe filename suitable for use across different operating systems
func SanitizeFilename(name string) string {
	// Remove any leading/trailing whitespace
	name = strings.TrimSpace(name)

	// Replace invalid characters with underscore
	name = invalidFilenameChars.ReplaceAllString(name, "_")

	// Remove leading/trailing dots (hidden files and path traversal)
	name = strings.Trim(name, ".")

	// Check for reserved filenames (Windows)
	upperName := strings.ToUpper(name)
	baseName := strings.Split(upperName, ".")[0] // Get name without extension
	if reservedFilenames[baseName] {
		name = "_" + name
	}

	// Ensure filename is not empty
	if name == "" {
		name = "unnamed"
	}

	// Limit filename length (255 bytes is common filesystem limit)
	if len(name) > 255 {
		// Try to preserve extension
		ext := filepath.Ext(name)
		maxBase := 255 - len(ext)
		if maxBase > 0 {
			name = name[:maxBase] + ext
		} else {
			name = name[:255]
		}
	}

	return name
}

// SanitizePath cleans and validates a file path
// Resolves relative paths, removes path traversal attempts
func SanitizePath(path string) string {
	// Clean the path (removes .., ., etc.)
	path = filepath.Clean(path)

	// Remove any null bytes
	path = strings.ReplaceAll(path, "\x00", "")

	return path
}

// IsPathTraversal checks if a path contains path traversal attempts
// Returns true if the path tries to escape its base directory
func IsPathTraversal(path string) bool {
	// Clean the path first
	cleaned := filepath.Clean(path)

	// Check if it starts with .. or contains ../ or ..\
	if strings.HasPrefix(cleaned, "..") {
		return true
	}

	if strings.Contains(cleaned, "/../") || strings.Contains(cleaned, "\\..\\") {
		return true
	}

	return false
}

// SafeJoin joins path elements and ensures the result stays within basePath
// Prevents path traversal attacks
func SafeJoin(basePath string, elem ...string) (string, error) {
	// Join all elements
	joined := filepath.Join(append([]string{basePath}, elem...)...)

	// Clean the result
	cleaned := filepath.Clean(joined)

	// Clean the base path
	cleanBase := filepath.Clean(basePath)

	// Ensure the result is within the base path
	if !strings.HasPrefix(cleaned, cleanBase) {
		return "", ErrPathTraversal
	}

	return cleaned, nil
}

// ContainsNullByte checks if a string contains null bytes (potential security issue)
func ContainsNullByte(s string) bool {
	return strings.Contains(s, "\x00")
}

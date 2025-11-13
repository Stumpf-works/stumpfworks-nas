// Package sysutil provides system-level utility functions for the Stumpf.Works NAS.
//
// This package centralizes common operations that interact with the operating system,
// such as finding system binaries, executing commands, checking file existence,
// privilege verification, user/group lookups, file operations, and network utilities.
//
// Key Features:
//
// Command Execution:
//   - Command discovery in system paths (FindCommand)
//   - Simplified command execution (RunCommand, RunCommandQuiet, RunCommandWithInput)
//
// Privilege and Security:
//   - Root privilege checking (IsRoot, RequireRoot)
//   - Path sanitization and validation (SanitizePath, SanitizeFilename, SafeJoin)
//   - Path traversal detection (IsPathTraversal)
//
// File Operations:
//   - File/directory existence checks (FileExists, DirExists, IsExecutable)
//   - File copying and moving (CopyFile, CopyDir, MoveFile, MoveDir)
//   - Sysfs file reading helpers (ReadSysFile)
//
// User and Group Management:
//   - UID/GID lookups by name (LookupUID, LookupGID)
//   - Username/Groupname lookups by ID (LookupUsername, LookupGroupname)
//   - Flexible parsing (ParseUIDOrUsername, ParseGIDOrGroupname)
//
// Network Utilities:
//   - IP address validation (ValidateIP, ValidateIPv4, ValidateIPv6)
//   - CIDR notation validation (ValidateCIDR)
//   - Private/Loopback IP detection (IsPrivateIP, IsLoopbackIP)
//   - Hostname validation (IsValidHostname)
//
// Example usage:
//
//	// Find and execute a system command
//	useraddPath := sysutil.FindCommand("useradd")
//	cmd := exec.Command(useraddPath, "-M", "-s", "/bin/false", username)
//
//	// Or use the helper
//	output, err := sysutil.RunCommand("useradd", "-M", "-s", "/bin/false", username)
//
//	// Check if running as root
//	if !sysutil.IsRoot() {
//	    return sysutil.ErrNotRoot
//	}
//
//	// Lookup user UID
//	uid, err := sysutil.LookupUID("john")
//
//	// Copy a file
//	err := sysutil.CopyFile("/src/file.txt", "/dst/file.txt")
//
//	// Sanitize a filename
//	safe := sysutil.SanitizeFilename(userInput)
//
//	// Validate an IP address
//	if sysutil.ValidateIP("192.168.1.1") {
//	    // Valid IP
//	}
package sysutil

// Package sysutil provides system-level utility functions for the Stumpf.Works NAS.
//
// This package centralizes common operations that interact with the operating system,
// such as finding system binaries, executing commands, checking file existence,
// and privilege verification.
//
// Key Features:
//   - Command discovery in system paths (FindCommand)
//   - Simplified command execution (RunCommand, RunCommandQuiet)
//   - Root privilege checking (IsRoot, RequireRoot)
//   - File/directory existence checks (FileExists, DirExists, IsExecutable)
//   - Sysfs file reading helpers (ReadSysFile)
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
package sysutil

// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
// Package executor provides common command execution interfaces and types
// used across all system management packages.
package executor

import "time"

// CommandResult represents the result of a shell command execution
type CommandResult struct {
	Command    string        `json:"command"`
	Args       []string      `json:"args"`
	Stdout     string        `json:"stdout"`
	Stderr     string        `json:"stderr"`
	ExitCode   int           `json:"exit_code"`
	Duration   time.Duration `json:"duration"`
	Success    bool          `json:"success"`
	Error      error         `json:"error,omitempty"`
	DryRun     bool          `json:"dry_run"`
}

// ShellExecutor defines the interface for executing shell commands
type ShellExecutor interface {
	// Execute runs a command and returns the result
	Execute(command string, args ...string) (*CommandResult, error)

	// ExecuteWithTimeout runs a command with a specific timeout
	ExecuteWithTimeout(timeout time.Duration, command string, args ...string) (*CommandResult, error)

	// CommandExists checks if a command exists in PATH
	CommandExists(command string) bool

	// SetDryRun enables or disables dry-run mode
	SetDryRun(enabled bool)

	// IsDryRun returns whether dry-run mode is enabled
	IsDryRun() bool
}

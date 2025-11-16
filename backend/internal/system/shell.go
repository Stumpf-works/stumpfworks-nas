// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package system

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// ShellExecutor provides safe command execution with logging, timeouts, and error handling
type ShellExecutor struct {
	defaultTimeout time.Duration
	dryRun         bool
	mu             sync.RWMutex
}

// CommandResult contains the result of a command execution
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

// CommandOptions holds options for command execution
type CommandOptions struct {
	// Timeout overrides the default timeout
	Timeout time.Duration

	// Dir sets the working directory
	Dir string

	// Env sets additional environment variables
	Env []string

	// User runs command as specific user (requires sudo)
	User string

	// IgnoreErrors doesn't return error on non-zero exit code
	IgnoreErrors bool

	// SuppressOutput prevents logging stdout/stderr
	SuppressOutput bool
}

// NewShellExecutor creates a new shell executor
func NewShellExecutor(defaultTimeout time.Duration, dryRun bool) (*ShellExecutor, error) {
	if defaultTimeout <= 0 {
		defaultTimeout = 30 * time.Second
	}

	return &ShellExecutor{
		defaultTimeout: defaultTimeout,
		dryRun:         dryRun,
	}, nil
}

// Execute executes a command with the given arguments
func (s *ShellExecutor) Execute(command string, args ...string) (*CommandResult, error) {
	return s.ExecuteWithOptions(command, nil, args...)
}

// ExecuteWithTimeout executes a command with a specific timeout
func (s *ShellExecutor) ExecuteWithTimeout(timeout time.Duration, command string, args ...string) (*CommandResult, error) {
	opts := &CommandOptions{Timeout: timeout}
	return s.ExecuteWithOptions(command, opts, args...)
}

// ExecuteWithOptions executes a command with advanced options
func (s *ShellExecutor) ExecuteWithOptions(command string, opts *CommandOptions, args ...string) (*CommandResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if opts == nil {
		opts = &CommandOptions{}
	}

	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = s.defaultTimeout
	}

	startTime := time.Now()

	result := &CommandResult{
		Command:  command,
		Args:     args,
		DryRun:   s.dryRun,
	}

	// Log command execution
	logger.Debug("Executing command",
		zap.String("command", command),
		zap.Strings("args", args),
		zap.Duration("timeout", timeout),
		zap.Bool("dry_run", s.dryRun))

	// Dry run mode
	if s.dryRun {
		result.Success = true
		result.Stdout = "[DRY RUN] Command not executed"
		result.Duration = time.Since(startTime)
		logger.Info("[DRY RUN] Command would be executed",
			zap.String("command", command),
			zap.Strings("args", args))
		return result, nil
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Prepare command
	cmd := exec.CommandContext(ctx, command, args...)

	// Set working directory if specified
	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}

	// Set environment variables if specified
	if len(opts.Env) > 0 {
		cmd.Env = append(cmd.Env, opts.Env...)
	}

	// If user is specified, wrap with sudo
	if opts.User != "" {
		originalArgs := append([]string{command}, args...)
		cmd = exec.CommandContext(ctx, "sudo", append([]string{"-u", opts.User}, originalArgs...)...)
	}

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute command
	err := cmd.Run()
	result.Duration = time.Since(startTime)
	result.Stdout = strings.TrimSpace(stdout.String())
	result.Stderr = strings.TrimSpace(stderr.String())

	// Get exit code
	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}

	// Handle errors
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = fmt.Errorf("command timed out after %v", timeout)
			logger.Error("Command timed out",
				zap.String("command", command),
				zap.Strings("args", args),
				zap.Duration("timeout", timeout))
		} else {
			result.Error = err
		}

		if !opts.IgnoreErrors {
			logger.Error("Command failed",
				zap.String("command", command),
				zap.Strings("args", args),
				zap.Int("exit_code", result.ExitCode),
				zap.String("stderr", result.Stderr),
				zap.Error(err))
			return result, fmt.Errorf("command failed: %w", err)
		} else {
			logger.Warn("Command failed (ignored)",
				zap.String("command", command),
				zap.Strings("args", args),
				zap.Int("exit_code", result.ExitCode),
				zap.String("stderr", result.Stderr))
		}
	}

	result.Success = err == nil || opts.IgnoreErrors

	// Log success
	if !opts.SuppressOutput {
		logger.Debug("Command executed successfully",
			zap.String("command", command),
			zap.Strings("args", args),
			zap.Duration("duration", result.Duration),
			zap.String("stdout", truncateString(result.Stdout, 500)),
			zap.String("stderr", truncateString(result.Stderr, 500)))
	}

	return result, nil
}

// ExecuteScript executes a shell script
func (s *ShellExecutor) ExecuteScript(script string, opts *CommandOptions) (*CommandResult, error) {
	return s.ExecuteWithOptions("bash", opts, "-c", script)
}

// ExecutePipe executes multiple commands in a pipe
func (s *ShellExecutor) ExecutePipe(commands ...string) (*CommandResult, error) {
	if len(commands) == 0 {
		return nil, fmt.Errorf("no commands provided")
	}

	pipeCommand := strings.Join(commands, " | ")
	return s.ExecuteScript(pipeCommand, nil)
}

// CommandExists checks if a command exists in PATH
func (s *ShellExecutor) CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// FindCommand searches for a command in PATH and common locations
func (s *ShellExecutor) FindCommand(command string, commonPaths ...string) (string, error) {
	// Try PATH first
	if path, err := exec.LookPath(command); err == nil {
		return path, nil
	}

	// Try common paths
	for _, path := range commonPaths {
		result, err := s.Execute("test", "-x", path)
		if err == nil && result.Success {
			return path, nil
		}
	}

	return "", fmt.Errorf("command '%s' not found", command)
}

// GetCommandVersion gets the version of a command
func (s *ShellExecutor) GetCommandVersion(command string) (string, error) {
	// Try --version first
	result, err := s.Execute(command, "--version")
	if err == nil && result.Success {
		return result.Stdout, nil
	}

	// Try -v
	result, err = s.Execute(command, "-v")
	if err == nil && result.Success {
		return result.Stdout, nil
	}

	// Try version
	result, err = s.Execute(command, "version")
	if err == nil && result.Success {
		return result.Stdout, nil
	}

	return "", fmt.Errorf("could not determine version for %s", command)
}

// RunAsRoot executes a command with sudo
func (s *ShellExecutor) RunAsRoot(command string, args ...string) (*CommandResult, error) {
	allArgs := append([]string{command}, args...)
	return s.Execute("sudo", allArgs...)
}

// IsSudoAvailable checks if sudo is available
func (s *ShellExecutor) IsSudoAvailable() bool {
	return s.CommandExists("sudo")
}

// IsRoot checks if running as root
func (s *ShellExecutor) IsRoot() bool {
	result, err := s.Execute("id", "-u")
	if err != nil {
		return false
	}
	return strings.TrimSpace(result.Stdout) == "0"
}

// SetDryRun enables or disables dry run mode
func (s *ShellExecutor) SetDryRun(dryRun bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dryRun = dryRun
}

// IsDryRun returns whether dry run mode is enabled
func (s *ShellExecutor) IsDryRun() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dryRun
}

// SetDefaultTimeout sets the default timeout for commands
func (s *ShellExecutor) SetDefaultTimeout(timeout time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.defaultTimeout = timeout
}

// GetDefaultTimeout returns the default timeout
func (s *ShellExecutor) GetDefaultTimeout() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.defaultTimeout
}

// Helper function to truncate strings for logging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "... (truncated)"
}

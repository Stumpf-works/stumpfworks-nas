// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package sysutil

import (
	"fmt"
	"os/exec"
)

// RunCommand executes a command and returns its combined output
// Automatically finds the command using FindCommand()
func RunCommand(name string, args ...string) (string, error) {
	cmdPath := FindCommand(name)
	cmd := exec.Command(cmdPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s failed: %s: %w", name, string(output), err)
	}
	return string(output), nil
}

// RunCommandQuiet executes a command and only returns error
// Use this when you don't need the output
func RunCommandQuiet(name string, args ...string) error {
	_, err := RunCommand(name, args...)
	return err
}

// RunCommandWithInput executes a command with stdin input
func RunCommandWithInput(input, name string, args ...string) (string, error) {
	cmdPath := FindCommand(name)
	cmd := exec.Command(cmdPath, args...)

	if input != "" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return "", fmt.Errorf("failed to create stdin pipe: %w", err)
		}

		go func() {
			defer stdin.Close()
			stdin.Write([]byte(input))
		}()
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s failed: %s: %w", name, string(output), err)
	}
	return string(output), nil
}

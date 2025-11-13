package sysutil

import "errors"

var (
	// ErrNotRoot is returned when an operation requires root privileges
	ErrNotRoot = errors.New("operation requires root privileges")

	// ErrCommandNotFound is returned when a required command is not found
	ErrCommandNotFound = errors.New("required command not found in system paths")
)

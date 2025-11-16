package handlers

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

type TerminalMessage struct {
	Type string `json:"type"` // "command", "interrupt", "resize"
	Data string `json:"data"`
}

type TerminalResponse struct {
	Type string `json:"type"` // "output", "error", "cwd"
	Data string `json:"data"`
}

type TerminalSession struct {
	conn         WSConn
	currentCmd   *exec.Cmd
	currentDir   string
	mu           sync.Mutex
	shellPath    string
	env          []string
}

// WSConn wraps the WebSocket connection for terminal
type WSConn interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
	Close() error
}

// gorillWSConn wraps gorilla/websocket.Conn to implement WSConn
type gorillaWSConn struct {
	conn interface {
		ReadJSON(v interface{}) error
		WriteJSON(v interface{}) error
		Close() error
	}
}

func (g *gorillaWSConn) ReadJSON(v interface{}) error  { return g.conn.ReadJSON(v) }
func (g *gorillaWSConn) WriteJSON(v interface{}) error { return g.conn.WriteJSON(v) }
func (g *gorillaWSConn) Close() error                  { return g.conn.Close() }

func NewTerminalSession(conn WSConn) *TerminalSession {
	homeDir, _ := os.UserHomeDir()
	if homeDir == "" {
		homeDir = "/root"
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	return &TerminalSession{
		conn:       conn,
		currentDir: homeDir,
		shellPath:  shell,
		env:        os.Environ(),
	}
}

func (ts *TerminalSession) sendOutput(output string) error {
	return ts.conn.WriteJSON(TerminalResponse{
		Type: "output",
		Data: output,
	})
}

func (ts *TerminalSession) sendError(errMsg string) error {
	return ts.conn.WriteJSON(TerminalResponse{
		Type: "error",
		Data: errMsg,
	})
}

func (ts *TerminalSession) sendCwd() error {
	return ts.conn.WriteJSON(TerminalResponse{
		Type: "cwd",
		Data: ts.currentDir,
	})
}

func (ts *TerminalSession) executeCommand(command string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// Handle built-in commands
	trimmed := strings.TrimSpace(command)

	if trimmed == "" {
		return nil
	}

	// Handle cd command separately as it changes the session state
	if strings.HasPrefix(trimmed, "cd ") || trimmed == "cd" {
		return ts.handleCd(trimmed)
	}

	// Execute command in current directory
	cmd := exec.Command(ts.shellPath, "-c", command)
	cmd.Dir = ts.currentDir
	cmd.Env = ts.env

	// Store current command for potential interrupt
	ts.currentCmd = cmd

	// Capture stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return ts.sendError("Failed to create stdout pipe: " + err.Error())
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return ts.sendError("Failed to create stderr pipe: " + err.Error())
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return ts.sendError("Failed to start command: " + err.Error())
	}

	// Stream output
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ts.streamOutput(stdout, false)
	}()

	go func() {
		defer wg.Done()
		ts.streamOutput(stderr, true)
	}()

	wg.Wait()

	// Wait for command to complete
	err = cmd.Wait()
	ts.currentCmd = nil

	if err != nil {
		// Check if it was killed by interrupt
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() {
					return ts.sendOutput("Process terminated by signal")
				}
			}
			return ts.sendError("Command exited with error: " + err.Error())
		}
		return ts.sendError("Command failed: " + err.Error())
	}

	return nil
}

func (ts *TerminalSession) streamOutput(reader io.Reader, isError bool) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if isError {
			ts.sendError(line)
		} else {
			ts.sendOutput(line)
		}
	}
}

func (ts *TerminalSession) handleCd(command string) error {
	parts := strings.Fields(command)
	var targetDir string

	if len(parts) == 1 {
		// Just "cd" - go to home
		homeDir, _ := os.UserHomeDir()
		if homeDir == "" {
			homeDir = "/root"
		}
		targetDir = homeDir
	} else {
		targetDir = parts[1]
	}

	// Expand ~ to home directory
	if strings.HasPrefix(targetDir, "~/") {
		homeDir, _ := os.UserHomeDir()
		if homeDir == "" {
			homeDir = "/root"
		}
		targetDir = homeDir + targetDir[1:]
	} else if targetDir == "~" {
		homeDir, _ := os.UserHomeDir()
		if homeDir == "" {
			homeDir = "/root"
		}
		targetDir = homeDir
	}

	// Handle relative paths
	if !strings.HasPrefix(targetDir, "/") {
		targetDir = ts.currentDir + "/" + targetDir
	}

	// Check if directory exists
	if info, err := os.Stat(targetDir); err != nil {
		return ts.sendError("cd: " + targetDir + ": No such file or directory")
	} else if !info.IsDir() {
		return ts.sendError("cd: " + targetDir + ": Not a directory")
	}

	ts.currentDir = targetDir
	return ts.sendCwd()
}

func (ts *TerminalSession) interrupt() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.currentCmd != nil && ts.currentCmd.Process != nil {
		// Send SIGINT to the process
		ts.currentCmd.Process.Signal(syscall.SIGINT)
	}
}

func (ts *TerminalSession) Handle() {
	defer ts.conn.Close()

	// Send initial working directory
	ts.sendCwd()

	for {
		var msg TerminalMessage
		err := ts.conn.ReadJSON(&msg)
		if err != nil {
			logger.Debug("Terminal WebSocket closed", zap.Error(err))
			return
		}

		switch msg.Type {
		case "command":
			if err := ts.executeCommand(msg.Data); err != nil {
				logger.Error("Failed to execute command", zap.Error(err))
			}

		case "interrupt":
			ts.interrupt()

		default:
			logger.Warn("Unknown terminal message type", zap.String("type", msg.Type))
		}
	}
}

// TerminalWebSocketHandler handles WebSocket connections for terminal access
func TerminalWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Failed to upgrade terminal WebSocket connection", zap.Error(err))
		return
	}

	logger.Info("Terminal WebSocket client connected", zap.String("remote_addr", r.RemoteAddr))

	// Create terminal session
	session := NewTerminalSession(&gorillaWSConn{conn: conn})
	session.Handle()

	logger.Info("Terminal WebSocket client disconnected", zap.String("remote_addr", r.RemoteAddr))
}

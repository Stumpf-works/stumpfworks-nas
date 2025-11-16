// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package plugins

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// Runtime manages plugin execution and lifecycle
type Runtime struct {
	service   *Service
	processes map[string]*PluginProcess
	mu        sync.RWMutex
}

// PluginProcess represents a running plugin process
type PluginProcess struct {
	PluginID  string
	Cmd       *exec.Cmd
	StartedAt time.Time
	Status    string // running, stopped, crashed, timeout
	LastError error
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewRuntime creates a new plugin runtime
func NewRuntime(service *Service) *Runtime {
	return &Runtime{
		service:   service,
		processes: make(map[string]*PluginProcess),
	}
}

// StartPlugin starts a plugin by executing its entry point
func (r *Runtime) StartPlugin(ctx context.Context, pluginID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if already running
	if proc, exists := r.processes[pluginID]; exists && proc.Status == "running" {
		return fmt.Errorf("plugin already running: %s", pluginID)
	}

	// Get plugin info
	plugin, err := r.service.GetPlugin(ctx, pluginID)
	if err != nil {
		return fmt.Errorf("failed to get plugin: %w", err)
	}

	if !plugin.Enabled {
		return fmt.Errorf("plugin is disabled: %s", pluginID)
	}

	// Load manifest to get entry point
	manifestPath := filepath.Join(plugin.InstallPath, "plugin.json")
	manifest, err := r.service.loadManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	if manifest.EntryPoint == "" {
		return fmt.Errorf("plugin has no entry point: %s", pluginID)
	}

	// Determine executable path
	execPath := filepath.Join(plugin.InstallPath, manifest.EntryPoint)
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		return fmt.Errorf("entry point not found: %s", execPath)
	}

	// Create context with timeout (plugins can run indefinitely unless stopped)
	procCtx, cancel := context.WithCancel(ctx)

	// Create command
	cmd := exec.CommandContext(procCtx, execPath)
	cmd.Dir = plugin.InstallPath
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PLUGIN_ID=%s", pluginID),
		fmt.Sprintf("PLUGIN_DIR=%s", plugin.InstallPath),
		fmt.Sprintf("NAS_API_URL=http://localhost:8080/api/v1"),
	)

	// Set up logging
	cmd.Stdout = logger.NewPluginLogger(pluginID, "stdout")
	cmd.Stderr = logger.NewPluginLogger(pluginID, "stderr")

	// Start process
	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start plugin: %w", err)
	}

	// Create process entry
	proc := &PluginProcess{
		PluginID:  pluginID,
		Cmd:       cmd,
		StartedAt: time.Now(),
		Status:    "running",
		ctx:       procCtx,
		cancel:    cancel,
	}

	r.processes[pluginID] = proc

	// Monitor process in background
	go r.monitorProcess(proc)

	logger.Info("Plugin started", zap.String("pluginID", pluginID))
	return nil
}

// StopPlugin stops a running plugin
func (r *Runtime) StopPlugin(ctx context.Context, pluginID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	proc, exists := r.processes[pluginID]
	if !exists {
		return fmt.Errorf("plugin not running: %s", pluginID)
	}

	if proc.Status != "running" {
		return fmt.Errorf("plugin not in running state: %s (status: %s)", pluginID, proc.Status)
	}

	// Cancel context to stop process
	proc.cancel()

	// Give it 5 seconds to gracefully shutdown
	done := make(chan error, 1)
	go func() {
		done <- proc.Cmd.Wait()
	}()

	select {
	case <-time.After(5 * time.Second):
		// Force kill if not stopped gracefully
		if proc.Cmd.Process != nil {
			logger.Warn("Plugin did not stop gracefully, killing", zap.String("pluginID", pluginID))
			proc.Cmd.Process.Kill()
		}
	case err := <-done:
		if err != nil {
			logger.Debug("Plugin exited with error", zap.String("pluginID", pluginID), zap.Error(err))
		}
	}

	proc.Status = "stopped"
	delete(r.processes, pluginID)

	logger.Info("Plugin stopped", zap.String("pluginID", pluginID))
	return nil
}

// RestartPlugin restarts a plugin
func (r *Runtime) RestartPlugin(ctx context.Context, pluginID string) error {
	// Stop if running
	if _, exists := r.processes[pluginID]; exists {
		if err := r.StopPlugin(ctx, pluginID); err != nil {
			logger.Warn("Failed to stop plugin before restart", zap.String("pluginID", pluginID), zap.Error(err))
		}
		// Wait a bit for cleanup
		time.Sleep(500 * time.Millisecond)
	}

	// Start again
	return r.StartPlugin(ctx, pluginID)
}

// GetPluginStatus returns the status of a plugin
func (r *Runtime) GetPluginStatus(pluginID string) (*PluginProcess, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	proc, exists := r.processes[pluginID]
	if !exists {
		return nil, fmt.Errorf("plugin not running: %s", pluginID)
	}

	return proc, nil
}

// ListRunningPlugins returns all running plugins
func (r *Runtime) ListRunningPlugins() []*PluginProcess {
	r.mu.RLock()
	defer r.mu.RUnlock()

	procs := make([]*PluginProcess, 0, len(r.processes))
	for _, proc := range r.processes {
		procs = append(procs, proc)
	}

	return procs
}

// StopAll stops all running plugins
func (r *Runtime) StopAll(ctx context.Context) error {
	r.mu.Lock()
	pluginIDs := make([]string, 0, len(r.processes))
	for id := range r.processes {
		pluginIDs = append(pluginIDs, id)
	}
	r.mu.Unlock()

	for _, id := range pluginIDs {
		if err := r.StopPlugin(ctx, id); err != nil {
			logger.Error("Failed to stop plugin during shutdown", zap.String("pluginID", id), zap.Error(err))
		}
	}

	return nil
}

// monitorProcess monitors a plugin process and updates its status
func (r *Runtime) monitorProcess(proc *PluginProcess) {
	err := proc.Cmd.Wait()

	r.mu.Lock()
	defer r.mu.Unlock()

	if err != nil {
		proc.LastError = err
		if proc.Cmd.ProcessState != nil && !proc.Cmd.ProcessState.Success() {
			proc.Status = "crashed"
			logger.Error("Plugin crashed",
				zap.String("pluginID", proc.PluginID),
				zap.Error(err))
		} else {
			// Context cancelled = intentional stop
			proc.Status = "stopped"
		}
	} else {
		proc.Status = "stopped"
		logger.Info("Plugin exited normally", zap.String("pluginID", proc.PluginID))
	}

	// Keep process in map for status reporting
	// It will be removed on explicit stop or restart
}

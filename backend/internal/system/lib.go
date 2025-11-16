// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package system

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// SystemLibrary provides a unified interface for all system-level operations
// This is the central point for ALL Debian/Linux system interactions
type SystemLibrary struct {
	// Shell provides safe command execution
	Shell *ShellExecutor

	// Metrics provides system metrics collection
	Metrics *MetricsCollector

	// Storage subsystems
	Storage *StorageManager

	// Network subsystems
	Network *NetworkManager

	// Sharing subsystems (SMB, NFS, iSCSI, WebDAV)
	Sharing *SharingManager

	// Users subsystems
	Users *UserManager

	// Internal state
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.RWMutex
}

// Config holds configuration for the system library
type Config struct {
	// EnableMetrics enables automatic metrics collection
	EnableMetrics bool

	// MetricsInterval is the interval for metrics collection
	MetricsInterval time.Duration

	// ShellTimeout is the default timeout for shell commands
	ShellTimeout time.Duration

	// DryRun mode doesn't execute actual system commands (for testing)
	DryRun bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		EnableMetrics:   true,
		MetricsInterval: 10 * time.Second,
		ShellTimeout:    30 * time.Second,
		DryRun:          false,
	}
}

// New creates a new SystemLibrary instance
// This should be called once at application startup
func New(cfg *Config) (*SystemLibrary, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	lib := &SystemLibrary{
		ctx:    ctx,
		cancel: cancel,
	}

	// Initialize shell executor
	shell, err := NewShellExecutor(cfg.ShellTimeout, cfg.DryRun)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize shell executor: %w", err)
	}
	lib.Shell = shell

	// Initialize metrics collector
	if cfg.EnableMetrics {
		metrics, err := NewMetricsCollector(cfg.MetricsInterval)
		if err != nil {
			logger.Warn("Failed to initialize metrics collector",
				zap.Error(err),
				zap.String("message", "System metrics will be limited"))
		} else {
			lib.Metrics = metrics
		}
	}

	// Initialize storage manager
	storage, err := NewStorageManager(shell)
	if err != nil {
		logger.Warn("Failed to initialize storage manager",
			zap.Error(err),
			zap.String("message", "Storage features may be limited"))
	} else {
		lib.Storage = storage
	}

	// Initialize network manager
	network, err := NewNetworkManager(shell)
	if err != nil {
		logger.Warn("Failed to initialize network manager",
			zap.Error(err),
			zap.String("message", "Network features may be limited"))
	} else {
		lib.Network = network
	}

	// Initialize sharing manager
	sharing, err := NewSharingManager(shell)
	if err != nil {
		logger.Warn("Failed to initialize sharing manager",
			zap.Error(err),
			zap.String("message", "Sharing features may be limited"))
	} else {
		lib.Sharing = sharing
	}

	// Initialize user manager
	users, err := NewUserManager(shell)
	if err != nil {
		logger.Warn("Failed to initialize user manager",
			zap.Error(err),
			zap.String("message", "User management features may be limited"))
	} else {
		lib.Users = users
	}

	logger.Info("System library initialized",
		zap.Bool("metrics_enabled", cfg.EnableMetrics),
		zap.Duration("metrics_interval", cfg.MetricsInterval),
		zap.Bool("dry_run", cfg.DryRun))

	return lib, nil
}

// Start starts all background tasks (metrics collection, etc.)
func (s *SystemLibrary) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Metrics != nil {
		if err := s.Metrics.Start(s.ctx); err != nil {
			return fmt.Errorf("failed to start metrics collector: %w", err)
		}
		logger.Info("Metrics collection started")
	}

	return nil
}

// Stop gracefully stops all background tasks and cleans up resources
func (s *SystemLibrary) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger.Info("Stopping system library...")

	// Cancel context to stop all background tasks
	s.cancel()

	// Stop metrics collector
	if s.Metrics != nil {
		if err := s.Metrics.Stop(); err != nil {
			logger.Error("Failed to stop metrics collector", zap.Error(err))
		}
	}

	logger.Info("System library stopped")
	return nil
}

// Health returns the health status of all subsystems
func (s *SystemLibrary) Health() (*HealthStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	health := &HealthStatus{
		Timestamp: time.Now(),
		Overall:   "healthy",
		Subsystems: make(map[string]SubsystemHealth),
	}

	// Check shell executor
	if s.Shell != nil {
		health.Subsystems["shell"] = SubsystemHealth{
			Status:  "healthy",
			Message: "Shell executor operational",
		}
	} else {
		health.Subsystems["shell"] = SubsystemHealth{
			Status:  "degraded",
			Message: "Shell executor not available",
		}
		health.Overall = "degraded"
	}

	// Check metrics collector
	if s.Metrics != nil {
		health.Subsystems["metrics"] = SubsystemHealth{
			Status:  "healthy",
			Message: "Metrics collection operational",
		}
	} else {
		health.Subsystems["metrics"] = SubsystemHealth{
			Status:  "disabled",
			Message: "Metrics collection disabled",
		}
	}

	// Check storage manager
	if s.Storage != nil {
		health.Subsystems["storage"] = SubsystemHealth{
			Status:  "healthy",
			Message: "Storage management operational",
		}
	} else {
		health.Subsystems["storage"] = SubsystemHealth{
			Status:  "degraded",
			Message: "Storage management limited",
		}
		health.Overall = "degraded"
	}

	// Check network manager
	if s.Network != nil {
		health.Subsystems["network"] = SubsystemHealth{
			Status:  "healthy",
			Message: "Network management operational",
		}
	} else {
		health.Subsystems["network"] = SubsystemHealth{
			Status:  "degraded",
			Message: "Network management limited",
		}
		health.Overall = "degraded"
	}

	// Check sharing manager
	if s.Sharing != nil {
		health.Subsystems["sharing"] = SubsystemHealth{
			Status:  "healthy",
			Message: "File sharing operational",
		}
	} else {
		health.Subsystems["sharing"] = SubsystemHealth{
			Status:  "degraded",
			Message: "File sharing limited",
		}
		health.Overall = "degraded"
	}

	// Check user manager
	if s.Users != nil {
		health.Subsystems["users"] = SubsystemHealth{
			Status:  "healthy",
			Message: "User management operational",
		}
	} else {
		health.Subsystems["users"] = SubsystemHealth{
			Status:  "degraded",
			Message: "User management limited",
		}
		health.Overall = "degraded"
	}

	return health, nil
}

// HealthStatus represents the overall health of the system library
type HealthStatus struct {
	Timestamp  time.Time                   `json:"timestamp"`
	Overall    string                      `json:"overall"` // healthy, degraded, unhealthy
	Subsystems map[string]SubsystemHealth  `json:"subsystems"`
}

// SubsystemHealth represents the health of a specific subsystem
type SubsystemHealth struct {
	Status  string `json:"status"`  // healthy, degraded, disabled, unhealthy
	Message string `json:"message"`
}

// Global instance (initialized by main.go)
var instance *SystemLibrary
var instanceMu sync.RWMutex

// Initialize initializes the global system library instance
func Initialize(cfg *Config) error {
	instanceMu.Lock()
	defer instanceMu.Unlock()

	if instance != nil {
		return fmt.Errorf("system library already initialized")
	}

	lib, err := New(cfg)
	if err != nil {
		return err
	}

	instance = lib
	return nil
}

// Get returns the global system library instance
func Get() *SystemLibrary {
	instanceMu.RLock()
	defer instanceMu.RUnlock()
	return instance
}

// MustGet returns the global instance or panics if not initialized
func MustGet() *SystemLibrary {
	lib := Get()
	if lib == nil {
		panic("system library not initialized - call Initialize() first")
	}
	return lib
}

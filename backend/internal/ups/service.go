// Revision: 2025-12-02 | Author: Claude | Version: 1.2.0
package ups

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// Service handles UPS monitoring and management
type Service struct {
	config          *models.UPSConfig
	mu              sync.RWMutex
	monitoring      bool
	lastStatus      *UPSStatus
	shutdownPending bool
	cancelMonitoring context.CancelFunc
}

// UPSStatus represents current UPS status
type UPSStatus struct {
	Online         bool    `json:"online"`
	BatteryCharge  int     `json:"battery_charge"`
	Runtime        int     `json:"runtime"`         // Estimated runtime in seconds
	LoadPercent    int     `json:"load_percent"`
	InputVoltage   float64 `json:"input_voltage"`
	OutputVoltage  float64 `json:"output_voltage"`
	Temperature    float64 `json:"temperature"`
	Status         string  `json:"status"`         // OL (online), OB (on battery), LB (low battery), etc.
	Model          string  `json:"model"`
	Manufacturer   string  `json:"manufacturer"`
	Serial         string  `json:"serial"`
	LastUpdate     time.Time `json:"last_update"`
}

var globalService *Service

// Initialize creates and starts the UPS service
func Initialize() (*Service, error) {
	logger.Info("Initializing UPS service...")

	service := &Service{
		monitoring: false,
	}

	// Load configuration from database
	if err := service.LoadConfig(); err != nil {
		logger.Warn("Failed to load UPS config, using defaults", zap.Error(err))
	}

	// Start monitoring if enabled
	if service.config != nil && service.config.Enabled {
		if err := service.StartMonitoring(); err != nil {
			logger.Error("Failed to start UPS monitoring", zap.Error(err))
			return service, err
		}
	}

	globalService = service
	return service, nil
}

// GetService returns the global UPS service instance
func GetService() *Service {
	return globalService
}

// LoadConfig loads UPS configuration from database
func (s *Service) LoadConfig() error {
	var config models.UPSConfig
	result := database.DB.First(&config)

	if result.Error != nil {
		// Create default config if not exists
		config = models.UPSConfig{
			Enabled:                false,
			UPSName:                "ups",
			UPSHost:                "localhost",
			UPSPort:                3493,
			PollInterval:           30,
			LowBatteryShutdown:     true,
			LowBatteryThreshold:    20,
			ShutdownDelay:          120,
			ShutdownCommand:        "shutdown -h now",
			NotifyOnPowerLoss:      true,
			NotifyOnBatteryLow:     true,
			NotifyOnPowerRestored:  true,
		}

		if err := database.DB.Create(&config).Error; err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
	}

	s.mu.Lock()
	s.config = &config
	s.mu.Unlock()

	return nil
}

// SaveConfig saves UPS configuration to database
func (s *Service) SaveConfig(config *models.UPSConfig) error {
	if err := database.DB.Save(config).Error; err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	s.mu.Lock()
	oldEnabled := s.config.Enabled
	s.config = config
	s.mu.Unlock()

	// Restart monitoring if enabled status changed
	if config.Enabled && !oldEnabled {
		return s.StartMonitoring()
	} else if !config.Enabled && oldEnabled {
		s.StopMonitoring()
	}

	return nil
}

// GetConfig returns the current UPS configuration
func (s *Service) GetConfig() *models.UPSConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

// StartMonitoring starts UPS monitoring
func (s *Service) StartMonitoring() error {
	s.mu.Lock()
	if s.monitoring {
		s.mu.Unlock()
		return fmt.Errorf("monitoring already started")
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelMonitoring = cancel
	s.monitoring = true
	config := s.config
	s.mu.Unlock()

	logger.Info("Starting UPS monitoring",
		zap.String("ups_name", config.UPSName),
		zap.String("host", config.UPSHost),
		zap.Int("poll_interval", config.PollInterval))

	// Start monitoring goroutine
	go s.monitorLoop(ctx, config)

	return nil
}

// StopMonitoring stops UPS monitoring
func (s *Service) StopMonitoring() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.monitoring {
		return
	}

	logger.Info("Stopping UPS monitoring")

	if s.cancelMonitoring != nil {
		s.cancelMonitoring()
	}

	s.monitoring = false
}

// monitorLoop continuously monitors UPS status
func (s *Service) monitorLoop(ctx context.Context, config *models.UPSConfig) {
	ticker := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	defer ticker.Stop()

	// Initial poll
	s.pollUPS(config)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.pollUPS(config)
		}
	}
}

// pollUPS polls the UPS for current status
func (s *Service) pollUPS(config *models.UPSConfig) {
	status, err := s.QueryUPS(config)
	if err != nil {
		logger.Error("Failed to query UPS", zap.Error(err))
		return
	}

	s.mu.Lock()
	lastStatus := s.lastStatus
	s.lastStatus = status
	s.mu.Unlock()

	// Check for events
	s.checkEvents(lastStatus, status, config)
}

// checkEvents checks for power events and triggers actions
func (s *Service) checkEvents(lastStatus, currentStatus *UPSStatus, config *models.UPSConfig) {
	// Check for power loss
	if lastStatus != nil && lastStatus.Online && !currentStatus.Online {
		s.logEvent("POWER_LOSS", "AC power lost, running on battery", currentStatus, "warning")

		if config.NotifyOnPowerLoss {
			// TODO: Send notification through alert service
			logger.Warn("UPS: AC power lost, running on battery")
		}
	}

	// Check for power restored
	if lastStatus != nil && !lastStatus.Online && currentStatus.Online {
		s.logEvent("POWER_RESTORED", "AC power restored", currentStatus, "info")

		if config.NotifyOnPowerRestored {
			// TODO: Send notification through alert service
			logger.Info("UPS: AC power restored")
		}

		// Cancel shutdown if pending
		s.mu.Lock()
		if s.shutdownPending {
			s.shutdownPending = false
			logger.Info("Canceling pending shutdown, power restored")
		}
		s.mu.Unlock()
	}

	// Check for low battery
	if !currentStatus.Online && currentStatus.BatteryCharge <= config.LowBatteryThreshold {
		if lastStatus == nil || lastStatus.BatteryCharge > config.LowBatteryThreshold {
			s.logEvent("BATTERY_LOW", fmt.Sprintf("Battery level critical: %d%%", currentStatus.BatteryCharge), currentStatus, "critical")

			if config.NotifyOnBatteryLow {
				logger.Error("UPS: Battery level critical", zap.Int("charge", currentStatus.BatteryCharge))
			}
		}

		// Trigger shutdown if enabled and not already pending
		s.mu.RLock()
		shutdownPending := s.shutdownPending
		s.mu.RUnlock()

		if config.LowBatteryShutdown && !shutdownPending {
			go s.initiateShutdown(config)
		}
	}
}

// initiateShutdown initiates system shutdown
func (s *Service) initiateShutdown(config *models.UPSConfig) {
	s.mu.Lock()
	if s.shutdownPending {
		s.mu.Unlock()
		return
	}
	s.shutdownPending = true
	status := s.lastStatus
	s.mu.Unlock()

	logger.Error("UPS: Initiating system shutdown due to low battery",
		zap.Int("delay_seconds", config.ShutdownDelay))

	s.logEvent("SHUTDOWN_INITIATED",
		fmt.Sprintf("System shutdown initiated, delay: %d seconds", config.ShutdownDelay),
		status, "critical")

	// Wait for shutdown delay
	time.Sleep(time.Duration(config.ShutdownDelay) * time.Second)

	// Check if shutdown was cancelled (power restored)
	s.mu.RLock()
	if !s.shutdownPending {
		s.mu.RUnlock()
		return
	}
	s.mu.RUnlock()

	// Execute shutdown command
	logger.Error("UPS: Executing shutdown command", zap.String("command", config.ShutdownCommand))

	cmd := exec.Command("sh", "-c", config.ShutdownCommand)
	if err := cmd.Run(); err != nil {
		logger.Error("Failed to execute shutdown command", zap.Error(err))
		s.logEvent("SHUTDOWN_FAILED", fmt.Sprintf("Shutdown command failed: %v", err), status, "critical")

		s.mu.Lock()
		s.shutdownPending = false
		s.mu.Unlock()
	}
}

// logEvent logs a UPS event to the database
func (s *Service) logEvent(eventType, description string, status *UPSStatus, severity string) {
	event := models.UPSEvent{
		EventType:   eventType,
		Description: description,
		Severity:    severity,
	}

	if status != nil {
		event.BatteryLevel = status.BatteryCharge
		event.Runtime = status.Runtime
		event.LoadPercent = status.LoadPercent
		event.Voltage = status.InputVoltage
	}

	if err := database.DB.Create(&event).Error; err != nil {
		logger.Error("Failed to log UPS event", zap.Error(err))
	}
}

// GetStatus returns the current UPS status
func (s *Service) GetStatus() (*UPSStatus, error) {
	s.mu.RLock()
	config := s.config
	lastStatus := s.lastStatus
	s.mu.RUnlock()

	if config == nil || !config.Enabled {
		return nil, fmt.Errorf("UPS monitoring is not enabled")
	}

	// Return cached status if recent (< 1 minute old)
	if lastStatus != nil && time.Since(lastStatus.LastUpdate) < time.Minute {
		return lastStatus, nil
	}

	// Query fresh status
	return s.QueryUPS(config)
}

// GetEvents returns UPS events from the database
func (s *Service) GetEvents(limit int, offset int) ([]models.UPSEvent, error) {
	var events []models.UPSEvent

	query := database.DB.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}

	return events, nil
}

// TestUPS tests the UPS connection and returns status
func (s *Service) TestUPS(config *models.UPSConfig) (*UPSStatus, error) {
	return s.QueryUPS(config)
}

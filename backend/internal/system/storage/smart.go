// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package storage

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// SMARTManager manages disk S.M.A.R.T. monitoring
type SMARTManager struct {
	shell      executor.ShellExecutor
	enabled bool
}

// SMARTInfo contains S.M.A.R.T. information for a disk
type SMARTInfo struct {
	Device           string            `json:"device"`
	Model            string            `json:"model"`
	SerialNumber     string            `json:"serial_number"`
	Firmware         string            `json:"firmware"`
	Capacity         uint64            `json:"capacity"`
	SectorSize       uint64            `json:"sector_size"`
	RotationRate     string            `json:"rotation_rate"`
	FormFactor       string            `json:"form_factor"`
	SmartSupported   bool              `json:"smart_supported"`
	SmartEnabled     bool              `json:"smart_enabled"`
	SmartStatus      string            `json:"smart_status"` // PASSED, FAILED
	Temperature      int               `json:"temperature_celsius"`
	PowerOnHours     uint64            `json:"power_on_hours"`
	PowerCycleCount  uint64            `json:"power_cycle_count"`
	ReallocatedSectors uint64          `json:"reallocated_sectors"`
	PendingSectors   uint64            `json:"pending_sectors"`
	UncorrectableErrors uint64         `json:"uncorrectable_errors"`
	Attributes       []SMARTAttribute  `json:"attributes"`
	HealthScore      int               `json:"health_score"` // 0-100
}

// SMARTAttribute represents a single SMART attribute
type SMARTAttribute struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Value      int    `json:"value"`
	Worst      int    `json:"worst"`
	Threshold  int    `json:"threshold"`
	Raw        uint64 `json:"raw"`
	Status     string `json:"status"` // OK, WARN, FAIL
}

// NewSMARTManager creates a new SMART manager
func NewSMARTManager(shell executor.ShellExecutor) (*SMARTManager, error) {
	// Check if smartctl is available
	if !shell.CommandExists("smartctl") {
		return nil, fmt.Errorf("smartmontools not installed (smartctl command not found)")
	}

	return &SMARTManager{
		shell:   shell,
		enabled: true,
	}, nil
}

// IsEnabled returns whether SMART monitoring is available
func (s *SMARTManager) IsEnabled() bool {
	return s.enabled
}

// GetInfo gets S.M.A.R.T. info for a device
func (s *SMARTManager) GetInfo(device string) (*SMARTInfo, error) {
	if !s.enabled {
		return nil, fmt.Errorf("SMART not available")
	}

	// Ensure device path
	if !strings.HasPrefix(device, "/dev/") {
		device = "/dev/" + device
	}

	// Get all SMART data
	result, err := s.shell.Execute("smartctl", "-a", device)
	if err != nil && result.ExitCode > 4 {
		// Exit codes 0-4 are acceptable (0=OK, 1-4=disk has issues but command succeeded)
		return nil, fmt.Errorf("failed to get SMART data: %w", err)
	}

	info := &SMARTInfo{
		Device: device,
		SmartStatus: "UNKNOWN",
	}

	// Parse output
	s.parseSmartOutput(result.Stdout, info)

	// Calculate health score
	info.HealthScore = s.calculateHealthScore(info)

	return info, nil
}

// parseSmartOutput parses smartctl output
func (s *SMARTManager) parseSmartOutput(output string, info *SMARTInfo) {
	lines := strings.Split(output, "\n")

	inAttributeSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Model
		if strings.HasPrefix(line, "Device Model:") || strings.HasPrefix(line, "Model Number:") {
			info.Model = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		}

		// Serial Number
		if strings.HasPrefix(line, "Serial Number:") || strings.HasPrefix(line, "Serial number:") {
			info.SerialNumber = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		}

		// Firmware
		if strings.HasPrefix(line, "Firmware Version:") {
			info.Firmware = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		}

		// Capacity
		if strings.Contains(line, "User Capacity:") {
			re := regexp.MustCompile(`\[([\d,]+) bytes\]`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				bytesStr := strings.ReplaceAll(matches[1], ",", "")
				if capacity, err := strconv.ParseUint(bytesStr, 10, 64); err == nil {
					info.Capacity = capacity
				}
			}
		}

		// Rotation Rate
		if strings.HasPrefix(line, "Rotation Rate:") {
			info.RotationRate = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		}

		// SMART Support
		if strings.Contains(line, "SMART support is: Available") {
			info.SmartSupported = true
		}
		if strings.Contains(line, "SMART support is: Enabled") {
			info.SmartEnabled = true
		}

		// SMART Status
		if strings.Contains(line, "SMART overall-health self-assessment test result:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) > 1 {
				status := strings.TrimSpace(parts[1])
				if strings.Contains(status, "PASSED") {
					info.SmartStatus = "PASSED"
				} else if strings.Contains(status, "FAILED") {
					info.SmartStatus = "FAILED"
				}
			}
		}

		// Temperature
		if strings.Contains(line, "Temperature:") || strings.Contains(line, "194 Temperature") {
			re := regexp.MustCompile(`(\d+) Celsius`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if temp, err := strconv.Atoi(matches[1]); err == nil {
					info.Temperature = temp
				}
			}
		}

		// Check if we're in the attribute section
		if strings.Contains(line, "ID# ATTRIBUTE_NAME") {
			inAttributeSection = true
			continue
		}

		// Parse attributes
		if inAttributeSection && len(line) > 0 && line[0] >= '0' && line[0] <= '9' {
			attr := s.parseAttribute(line)
			if attr != nil {
				info.Attributes = append(info.Attributes, *attr)

				// Extract specific important values
				switch attr.ID {
				case 5: // Reallocated Sectors
					info.ReallocatedSectors = attr.Raw
				case 9: // Power On Hours
					info.PowerOnHours = attr.Raw
				case 12: // Power Cycle Count
					info.PowerCycleCount = attr.Raw
				case 197: // Current Pending Sectors
					info.PendingSectors = attr.Raw
				case 198: // Uncorrectable Errors
					info.UncorrectableErrors = attr.Raw
				}
			}
		}
	}
}

// parseAttribute parses a single SMART attribute line
func (s *SMARTManager) parseAttribute(line string) *SMARTAttribute {
	fields := strings.Fields(line)
	if len(fields) < 10 {
		return nil
	}

	attr := &SMARTAttribute{}

	// ID
	if id, err := strconv.Atoi(fields[0]); err == nil {
		attr.ID = id
	}

	// Name
	attr.Name = fields[1]

	// Value
	if value, err := strconv.Atoi(fields[3]); err == nil {
		attr.Value = value
	}

	// Worst
	if worst, err := strconv.Atoi(fields[4]); err == nil {
		attr.Worst = worst
	}

	// Threshold
	if threshold, err := strconv.Atoi(fields[5]); err == nil {
		attr.Threshold = threshold
	}

	// Raw Value (last field, may contain spaces)
	rawValue := fields[len(fields)-1]
	if raw, err := strconv.ParseUint(rawValue, 10, 64); err == nil {
		attr.Raw = raw
	}

	// Determine status
	if attr.Value <= attr.Threshold {
		attr.Status = "FAIL"
	} else if attr.Value < attr.Threshold+10 {
		attr.Status = "WARN"
	} else {
		attr.Status = "OK"
	}

	return attr
}

// calculateHealthScore calculates overall health score (0-100)
func (s *SMARTManager) calculateHealthScore(info *SMARTInfo) int {
	if !info.SmartSupported || !info.SmartEnabled {
		return 0
	}

	score := 100

	// SMART status failed = 0
	if info.SmartStatus == "FAILED" {
		return 0
	}

	// Critical attributes
	if info.ReallocatedSectors > 0 {
		score -= int(info.ReallocatedSectors * 5) // -5 per reallocated sector
	}
	if info.PendingSectors > 0 {
		score -= int(info.PendingSectors * 10) // -10 per pending sector
	}
	if info.UncorrectableErrors > 0 {
		score -= int(info.UncorrectableErrors * 10) // -10 per uncorrectable error
	}

	// Temperature penalty (over 50°C)
	if info.Temperature > 50 {
		score -= (info.Temperature - 50) * 2
	}

	// Attribute health
	failCount := 0
	warnCount := 0
	for _, attr := range info.Attributes {
		if attr.Status == "FAIL" {
			failCount++
		} else if attr.Status == "WARN" {
			warnCount++
		}
	}

	score -= failCount * 20
	score -= warnCount * 5

	// Clamp to 0-100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// RunTest runs a SMART self-test
func (s *SMARTManager) RunTest(device string, testType string) error {
	if !s.enabled {
		return fmt.Errorf("SMART not available")
	}

	if !strings.HasPrefix(device, "/dev/") {
		device = "/dev/" + device
	}

	// Valid test types: short, long, conveyance
	validTypes := map[string]bool{
		"short":      true,
		"long":       true,
		"conveyance": true,
	}

	if !validTypes[testType] {
		return fmt.Errorf("invalid test type: %s (must be short, long, or conveyance)", testType)
	}

	_, err := s.shell.Execute("smartctl", "-t", testType, device)
	if err != nil {
		return fmt.Errorf("failed to start SMART test: %w", err)
	}

	return nil
}

// EnableSMART enables SMART on a device
func (s *SMARTManager) EnableSMART(device string) error {
	if !strings.HasPrefix(device, "/dev/") {
		device = "/dev/" + device
	}

	_, err := s.shell.Execute("smartctl", "-s", "on", device)
	if err != nil {
		return fmt.Errorf("failed to enable SMART: %w", err)
	}

	return nil
}

// GetTestLog gets the self-test log
func (s *SMARTManager) GetTestLog(device string) (string, error) {
	if !strings.HasPrefix(device, "/dev/") {
		device = "/dev/" + device
	}

	result, err := s.shell.Execute("smartctl", "-l", "selftest", device)
	if err != nil {
		return "", fmt.Errorf("failed to get test log: %w", err)
	}

	return result.Stdout, nil
}

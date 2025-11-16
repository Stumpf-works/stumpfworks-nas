// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package storage

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// RAIDManager manages software RAID (mdadm)
type RAIDManager struct {
	shell   ShellExecutor
	enabled bool
}

// RAIDArray represents a RAID array
type RAIDArray struct {
	Device      string       `json:"device"`
	Name        string       `json:"name"`
	Level       string       `json:"level"` // raid0, raid1, raid5, raid6, raid10
	State       string       `json:"state"` // clean, active, degraded, recovering
	Size        uint64       `json:"size"`
	UsedDevices int          `json:"used_devices"`
	TotalDevices int         `json:"total_devices"`
	ActiveDevices int        `json:"active_devices"`
	WorkingDevices int       `json:"working_devices"`
	FailedDevices int        `json:"failed_devices"`
	SpareDevices int         `json:"spare_devices"`
	UUID        string       `json:"uuid"`
	Devices     []RAIDDevice `json:"devices"`
}

// RAIDDevice represents a device in a RAID array
type RAIDDevice struct {
	Device string `json:"device"`
	Number int    `json:"number"`
	State  string `json:"state"` // active, faulty, spare
	Role   string `json:"role"`
}

// NewRAIDManager creates a new RAID manager
func NewRAIDManager(shell ShellExecutor) (*RAIDManager, error) {
	if !shell.CommandExists("mdadm") {
		return nil, fmt.Errorf("mdadm not installed")
	}

	return &RAIDManager{
		shell:   shell,
		enabled: true,
	}, nil
}

// IsEnabled returns whether RAID is available
func (r *RAIDManager) IsEnabled() bool {
	return r.enabled
}

// ListArrays lists all RAID arrays
func (r *RAIDManager) ListArrays() ([]RAIDArray, error) {
	result, err := r.shell.Execute("cat", "/proc/mdstat")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc/mdstat: %w", err)
	}

	var arrays []RAIDArray
	lines := strings.Split(result.Stdout, "\n")

	var current *RAIDArray
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// New array line (starts with md)
		if strings.HasPrefix(line, "md") {
			if current != nil {
				arrays = append(arrays, *current)
			}

			current = &RAIDArray{}
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				current.Device = "/dev/" + fields[0]
				current.State = strings.TrimSuffix(fields[2], ":")

				// Parse level
				if len(fields) >= 4 {
					current.Level = fields[3]
				}
			}
		} else if current != nil && strings.TrimSpace(line) != "" {
			// Parse additional info (second line has size, devices)
			if strings.Contains(line, "blocks") {
				re := regexp.MustCompile(`(\d+) blocks`)
				if matches := re.FindStringSubmatch(line); len(matches) > 1 {
					if size, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
						current.Size = size * 1024 // Convert KB to bytes
					}
				}

				// Parse device count [N/M] [UU_]
				re = regexp.MustCompile(`\[(\d+)/(\d+)\]`)
				if matches := re.FindStringSubmatch(line); len(matches) > 2 {
					if used, err := strconv.Atoi(matches[1]); err == nil {
						current.UsedDevices = used
					}
					if total, err := strconv.Atoi(matches[2]); err == nil {
						current.TotalDevices = total
					}
				}
			}
		}
	}

	if current != nil {
		arrays = append(arrays, *current)
	}

	// Get detailed info for each array
	for i := range arrays {
		if err := r.getArrayDetails(&arrays[i]); err != nil {
			// Log but don't fail
			continue
		}
	}

	return arrays, nil
}

// getArrayDetails gets detailed information about an array
func (r *RAIDManager) getArrayDetails(array *RAIDArray) error {
	result, err := r.shell.Execute("mdadm", "--detail", array.Device)
	if err != nil {
		return err
	}

	lines := strings.Split(result.Stdout, "\n")
	inDeviceSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "UUID :") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				array.UUID = strings.TrimSpace(parts[1])
			}
		}

		if strings.Contains(line, "Name :") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				array.Name = strings.TrimSpace(parts[1])
			}
		}

		if strings.Contains(line, "Active Devices :") {
			re := regexp.MustCompile(`Active Devices : (\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if count, err := strconv.Atoi(matches[1]); err == nil {
					array.ActiveDevices = count
				}
			}
		}

		if strings.Contains(line, "Working Devices :") {
			re := regexp.MustCompile(`Working Devices : (\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if count, err := strconv.Atoi(matches[1]); err == nil {
					array.WorkingDevices = count
				}
			}
		}

		if strings.Contains(line, "Failed Devices :") {
			re := regexp.MustCompile(`Failed Devices : (\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if count, err := strconv.Atoi(matches[1]); err == nil {
					array.FailedDevices = count
				}
			}
		}

		if strings.Contains(line, "Spare Devices :") {
			re := regexp.MustCompile(`Spare Devices : (\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if count, err := strconv.Atoi(matches[1]); err == nil {
					array.SpareDevices = count
				}
			}
		}

		// Device section starts
		if strings.Contains(line, "Number") && strings.Contains(line, "Major") {
			inDeviceSection = true
			continue
		}

		// Parse device lines
		if inDeviceSection && len(line) > 0 {
			fields := strings.Fields(line)
			if len(fields) >= 7 {
				device := RAIDDevice{}

				if num, err := strconv.Atoi(fields[0]); err == nil {
					device.Number = num
				}

				device.Device = fields[6]
				device.State = fields[4]
				device.Role = fields[5]

				array.Devices = append(array.Devices, device)
			}
		}
	}

	return nil
}

// CreateArray creates a new RAID array
func (r *RAIDManager) CreateArray(device string, level string, devices []string, spares []string) error {
	if len(devices) == 0 {
		return fmt.Errorf("no devices specified")
	}

	args := []string{"--create", device, "--level=" + level, fmt.Sprintf("--raid-devices=%d", len(devices))}

	if len(spares) > 0 {
		args = append(args, fmt.Sprintf("--spare-devices=%d", len(spares)))
	}

	args = append(args, devices...)
	args = append(args, spares...)

	_, err := r.shell.Execute("mdadm", args...)
	if err != nil {
		return fmt.Errorf("failed to create RAID array: %w", err)
	}

	return nil
}

// StopArray stops a RAID array
func (r *RAIDManager) StopArray(device string) error {
	_, err := r.shell.Execute("mdadm", "--stop", device)
	if err != nil {
		return fmt.Errorf("failed to stop array: %w", err)
	}

	return nil
}

// AssembleArray assembles a RAID array
func (r *RAIDManager) AssembleArray(device string, uuid string) error {
	args := []string{"--assemble", device}
	if uuid != "" {
		args = append(args, "--uuid="+uuid)
	}

	_, err := r.shell.Execute("mdadm", args...)
	if err != nil {
		return fmt.Errorf("failed to assemble array: %w", err)
	}

	return nil
}

// AddDevice adds a device to a RAID array
func (r *RAIDManager) AddDevice(arrayDevice string, newDevice string) error {
	_, err := r.shell.Execute("mdadm", "--add", arrayDevice, newDevice)
	if err != nil {
		return fmt.Errorf("failed to add device: %w", err)
	}

	return nil
}

// RemoveDevice removes a device from a RAID array
func (r *RAIDManager) RemoveDevice(arrayDevice string, device string) error {
	// Mark as failed first
	if _, err := r.shell.Execute("mdadm", "--fail", arrayDevice, device); err != nil {
		return fmt.Errorf("failed to mark device as failed: %w", err)
	}

	// Then remove
	_, err := r.shell.Execute("mdadm", "--remove", arrayDevice, device)
	if err != nil {
		return fmt.Errorf("failed to remove device: %w", err)
	}

	return nil
}

// GrowArray grows a RAID array (add more devices)
func (r *RAIDManager) GrowArray(device string, newDeviceCount int) error {
	_, err := r.shell.Execute("mdadm", "--grow", device, fmt.Sprintf("--raid-devices=%d", newDeviceCount))
	if err != nil {
		return fmt.Errorf("failed to grow array: %w", err)
	}

	return nil
}

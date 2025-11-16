// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package storage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetDiskIOStats retrieves I/O statistics for all disks
func GetDiskIOStats() ([]DiskIOStats, error) {
	var stats []DiskIOStats

	// Read /proc/diskstats
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	timestamp := time.Now()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 14 {
			continue
		}

		diskName := fields[2]

		// Skip partitions and special devices
		if strings.HasPrefix(diskName, "loop") ||
			strings.HasPrefix(diskName, "ram") ||
			strings.HasPrefix(diskName, "dm-") ||
			len(diskName) > 3 && diskName[len(diskName)-1] >= '0' && diskName[len(diskName)-1] <= '9' {
			continue
		}

		// Parse fields according to /proc/diskstats format
		readOps, _ := strconv.ParseUint(fields[3], 10, 64)
		readSectors, _ := strconv.ParseUint(fields[5], 10, 64)
		writeOps, _ := strconv.ParseUint(fields[7], 10, 64)
		writeSectors, _ := strconv.ParseUint(fields[9], 10, 64)

		stat := DiskIOStats{
			DiskName:   diskName,
			ReadBytes:  readSectors * 512,  // Sectors are 512 bytes
			WriteBytes: writeSectors * 512,
			ReadOps:    readOps,
			WriteOps:   writeOps,
			Timestamp:  timestamp,
		}

		stats = append(stats, stat)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

// GetDiskIOStatsForDisk retrieves I/O statistics for a specific disk
func GetDiskIOStatsForDisk(diskName string) (*DiskIOStats, error) {
	allStats, err := GetDiskIOStats()
	if err != nil {
		return nil, err
	}

	for _, stat := range allStats {
		if stat.DiskName == diskName {
			return &stat, nil
		}
	}

	return nil, fmt.Errorf("disk not found: %s", diskName)
}

// CalculateIORate calculates the I/O rate between two stat snapshots
func CalculateIORate(previous, current *DiskIOStats) *DiskIOStats {
	if previous == nil || current == nil {
		return current
	}

	duration := current.Timestamp.Sub(previous.Timestamp).Seconds()
	if duration <= 0 {
		return current
	}

	rate := &DiskIOStats{
		DiskName:  current.DiskName,
		Timestamp: current.Timestamp,
	}

	// Calculate bytes/sec
	rate.ReadBytes = uint64(float64(current.ReadBytes-previous.ReadBytes) / duration)
	rate.WriteBytes = uint64(float64(current.WriteBytes-previous.WriteBytes) / duration)

	// Calculate ops/sec
	rate.ReadOps = uint64(float64(current.ReadOps-previous.ReadOps) / duration)
	rate.WriteOps = uint64(float64(current.WriteOps-previous.WriteOps) / duration)

	return rate
}

// MonitorDiskIO monitors disk I/O continuously
type DiskIOMonitor struct {
	interval      time.Duration
	previousStats map[string]*DiskIOStats
	stopChan      chan struct{}
	statsChan     chan []DiskIOStats
}

// NewDiskIOMonitor creates a new disk I/O monitor
func NewDiskIOMonitor(interval time.Duration) *DiskIOMonitor {
	return &DiskIOMonitor{
		interval:      interval,
		previousStats: make(map[string]*DiskIOStats),
		stopChan:      make(chan struct{}),
		statsChan:     make(chan []DiskIOStats, 10),
	}
}

// Start starts the monitoring
func (m *DiskIOMonitor) Start() {
	go m.monitor()
}

// Stop stops the monitoring
func (m *DiskIOMonitor) Stop() {
	close(m.stopChan)
}

// Stats returns the stats channel
func (m *DiskIOMonitor) Stats() <-chan []DiskIOStats {
	return m.statsChan
}

// monitor continuously monitors disk I/O
func (m *DiskIOMonitor) monitor() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			close(m.statsChan)
			return
		case <-ticker.C:
			currentStats, err := GetDiskIOStats()
			if err != nil {
				continue
			}

			var rates []DiskIOStats

			for i, current := range currentStats {
				previous, exists := m.previousStats[current.DiskName]
				if exists {
					rate := CalculateIORate(previous, &current)
					rates = append(rates, *rate)
				}
				m.previousStats[current.DiskName] = &currentStats[i]
			}

			if len(rates) > 0 {
				select {
				case m.statsChan <- rates:
				default:
					// Channel full, skip this update
				}
			}
		}
	}
}

// DiskHealth represents the health assessment of a disk
type DiskHealth struct {
	DiskName    string     `json:"diskName"`
	Status      DiskStatus `json:"status"`
	Issues      []string   `json:"issues"`
	Temperature int        `json:"temperature"`
	Score       int        `json:"score"` // 0-100
}

// AssessDiskHealth assesses the health of a disk
func AssessDiskHealth(diskName string) (*DiskHealth, error) {
	disk, err := GetDiskInfo(diskName)
	if err != nil {
		return nil, err
	}

	health := &DiskHealth{
		DiskName:    diskName,
		Status:      disk.Status,
		Temperature: disk.Temperature,
		Score:       100,
		Issues:      []string{},
	}

	if disk.SMART == nil {
		health.Status = DiskStatusUnknown
		health.Score = 0
		health.Issues = append(health.Issues, "SMART data not available")
		return health, nil
	}

	smart := disk.SMART

	// Check critical issues
	if smart.ReallocatedSectors > 10 {
		health.Issues = append(health.Issues, fmt.Sprintf("High reallocated sectors: %d", smart.ReallocatedSectors))
		health.Score -= 30
	} else if smart.ReallocatedSectors > 0 {
		health.Issues = append(health.Issues, fmt.Sprintf("Reallocated sectors detected: %d", smart.ReallocatedSectors))
		health.Score -= 10
	}

	if smart.PendingSectors > 5 {
		health.Issues = append(health.Issues, fmt.Sprintf("High pending sectors: %d", smart.PendingSectors))
		health.Score -= 30
	} else if smart.PendingSectors > 0 {
		health.Issues = append(health.Issues, fmt.Sprintf("Pending sectors detected: %d", smart.PendingSectors))
		health.Score -= 10
	}

	if smart.UncorrectableErrors > 0 {
		health.Issues = append(health.Issues, fmt.Sprintf("Uncorrectable errors: %d", smart.UncorrectableErrors))
		health.Score -= 40
	}

	// Check temperature
	if smart.Temperature > 70 {
		health.Issues = append(health.Issues, fmt.Sprintf("High temperature: %d°C", smart.Temperature))
		health.Score -= 20
	} else if smart.Temperature > 60 {
		health.Issues = append(health.Issues, fmt.Sprintf("Elevated temperature: %d°C", smart.Temperature))
		health.Score -= 10
	}

	// Check SSD wear
	if disk.Type == DiskTypeSSD || disk.Type == DiskTypeNVMe {
		if smart.PercentLifeUsed > 95 {
			health.Issues = append(health.Issues, fmt.Sprintf("SSD near end of life: %d%% used", smart.PercentLifeUsed))
			health.Score -= 30
		} else if smart.PercentLifeUsed > 80 {
			health.Issues = append(health.Issues, fmt.Sprintf("SSD wear level high: %d%% used", smart.PercentLifeUsed))
			health.Score -= 15
		}
	}

	// Check CRC errors
	if smart.CRCErrors > 100 {
		health.Issues = append(health.Issues, fmt.Sprintf("High CRC errors: %d", smart.CRCErrors))
		health.Score -= 10
	}

	// Ensure score doesn't go below 0
	if health.Score < 0 {
		health.Score = 0
	}

	// Determine final status
	if health.Score >= 90 {
		health.Status = DiskStatusHealthy
	} else if health.Score >= 70 {
		health.Status = DiskStatusWarning
	} else if health.Score >= 50 {
		health.Status = DiskStatusCritical
	} else {
		health.Status = DiskStatusFailed
	}

	return health, nil
}

// GetAllDiskHealth gets health assessment for all disks
func GetAllDiskHealth() ([]DiskHealth, error) {
	disks, err := ListDisks()
	if err != nil {
		return nil, err
	}

	var healthList []DiskHealth

	for _, disk := range disks {
		health, err := AssessDiskHealth(disk.Name)
		if err != nil {
			continue
		}
		healthList = append(healthList, *health)
	}

	return healthList, nil
}

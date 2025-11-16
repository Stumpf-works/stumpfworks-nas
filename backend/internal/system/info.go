package system

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

var startTime = time.Now()

// SystemInfo represents basic system information
type SystemInfo struct {
	Hostname     string `json:"hostname"`
	Platform     string `json:"platform"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	CPUCores     int    `json:"cpuCores"`
	Uptime       uint64 `json:"uptime"` // seconds
	BootTime     uint64 `json:"bootTime"`
}

// GetSystemInfo returns basic system information
func GetSystemInfo() (*SystemInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	cpuCount, err := cpu.Counts(true)
	if err != nil {
		cpuCount = runtime.NumCPU()
	}

	return &SystemInfo{
		Hostname:     hostInfo.Hostname,
		Platform:     hostInfo.Platform,
		OS:           hostInfo.OS,
		Architecture: runtime.GOARCH,
		CPUCores:     cpuCount,
		Uptime:       hostInfo.Uptime,
		BootTime:     hostInfo.BootTime,
	}, nil
}

// RealtimeSystemMetrics represents real-time system metrics
type RealtimeSystemMetrics struct {
	CPU    RealtimeCPUMetrics    `json:"cpu"`
	Memory RealtimeMemoryMetrics `json:"memory"`
	Disk   []RealtimeDiskMetrics `json:"disk"`
	Network RealtimeNetworkMetrics `json:"network"`
	Timestamp int64 `json:"timestamp"`
}

// RealtimeCPUMetrics represents CPU usage metrics
type RealtimeCPUMetrics struct {
	UsagePercent float64   `json:"usagePercent"`
	PerCore      []float64 `json:"perCore,omitempty"`
}

// RealtimeMemoryMetrics represents memory usage metrics
type RealtimeMemoryMetrics struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// RealtimeDiskMetrics represents disk usage metrics
type RealtimeDiskMetrics struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// RealtimeNetworkMetrics represents network usage metrics
type RealtimeNetworkMetrics struct {
	BytesSent   uint64 `json:"bytesSent"`
	BytesRecv   uint64 `json:"bytesRecv"`
	PacketsSent uint64 `json:"packetsSent"`
	PacketsRecv uint64 `json:"packetsRecv"`
}

// GetRealtimeSystemMetrics returns real-time system metrics
func GetRealtimeSystemMetrics() (*RealtimeSystemMetrics, error) {
	metrics := &RealtimeSystemMetrics{
		Timestamp: time.Now().Unix(),
	}

	// CPU metrics
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil && len(cpuPercent) > 0 {
		metrics.CPU.UsagePercent = cpuPercent[0]
	}

	cpuPerCore, err := cpu.Percent(time.Second, true)
	if err == nil {
		metrics.CPU.PerCore = cpuPerCore
	}

	// Memory metrics
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		metrics.Memory = RealtimeMemoryMetrics{
			Total:       memInfo.Total,
			Available:   memInfo.Available,
			Used:        memInfo.Used,
			UsedPercent: memInfo.UsedPercent,
		}
	}

	// Disk metrics
	partitions, err := disk.Partitions(false)
	if err == nil {
		for _, partition := range partitions {
			usage, err := disk.Usage(partition.Mountpoint)
			if err == nil {
				metrics.Disk = append(metrics.Disk, RealtimeDiskMetrics{
					Device:      partition.Device,
					Mountpoint:  partition.Mountpoint,
					Fstype:      partition.Fstype,
					Total:       usage.Total,
					Free:        usage.Free,
					Used:        usage.Used,
					UsedPercent: usage.UsedPercent,
				})
			}
		}
	}

	// Network metrics
	netIO, err := net.IOCounters(false)
	if err == nil && len(netIO) > 0 {
		metrics.Network = RealtimeNetworkMetrics{
			BytesSent:   netIO[0].BytesSent,
			BytesRecv:   netIO[0].BytesRecv,
			PacketsSent: netIO[0].PacketsSent,
			PacketsRecv: netIO[0].PacketsRecv,
		}
	}

	return metrics, nil
}

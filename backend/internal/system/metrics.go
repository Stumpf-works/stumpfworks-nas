// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package system

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// MetricsCollector collects system metrics for monitoring and Prometheus export
type MetricsCollector struct {
	interval time.Duration
	mu       sync.RWMutex
	current  *SystemMetrics
	history  []*SystemMetrics
	maxHistory int
}

// SystemMetrics contains all system metrics
type SystemMetrics struct {
	Timestamp time.Time `json:"timestamp"`

	// CPU Metrics
	CPUUsagePercent    float64   `json:"cpu_usage_percent"`
	CPUCores           int       `json:"cpu_cores"`
	CPUPerCoreUsage    []float64 `json:"cpu_per_core_usage"`
	LoadAverage1       float64   `json:"load_average_1"`
	LoadAverage5       float64   `json:"load_average_5"`
	LoadAverage15      float64   `json:"load_average_15"`

	// Memory Metrics
	MemoryTotal        uint64  `json:"memory_total_bytes"`
	MemoryUsed         uint64  `json:"memory_used_bytes"`
	MemoryFree         uint64  `json:"memory_free_bytes"`
	MemoryUsagePercent float64 `json:"memory_usage_percent"`
	SwapTotal          uint64  `json:"swap_total_bytes"`
	SwapUsed           uint64  `json:"swap_used_bytes"`
	SwapFree           uint64  `json:"swap_free_bytes"`
	SwapUsagePercent   float64 `json:"swap_usage_percent"`

	// Disk Metrics (root partition)
	DiskTotal          uint64  `json:"disk_total_bytes"`
	DiskUsed           uint64  `json:"disk_used_bytes"`
	DiskFree           uint64  `json:"disk_free_bytes"`
	DiskUsagePercent   float64 `json:"disk_usage_percent"`

	// Network Metrics
	NetworkBytesSent   uint64 `json:"network_bytes_sent_total"`
	NetworkBytesRecv   uint64 `json:"network_bytes_recv_total"`
	NetworkPacketsSent uint64 `json:"network_packets_sent_total"`
	NetworkPacketsRecv uint64 `json:"network_packets_recv_total"`

	// System Info
	Uptime             uint64 `json:"uptime_seconds"`
	BootTime           uint64 `json:"boot_time"`
	Processes          uint64 `json:"processes_total"`

	// Go Runtime Metrics
	GoVersion          string `json:"go_version"`
	GoRoutines         int    `json:"goroutines"`
	GoMemAlloc         uint64 `json:"go_mem_alloc_bytes"`
	GoMemSys           uint64 `json:"go_mem_sys_bytes"`
	GoGCPauseNs        uint64 `json:"go_gc_pause_ns"`
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(interval time.Duration) (*MetricsCollector, error) {
	if interval <= 0 {
		interval = 10 * time.Second
	}

	mc := &MetricsCollector{
		interval:   interval,
		history:    make([]*SystemMetrics, 0, 100),
		maxHistory: 100, // Keep last 100 readings
	}

	// Collect initial metrics
	if err := mc.collect(); err != nil {
		return nil, fmt.Errorf("failed to collect initial metrics: %w", err)
	}

	return mc, nil
}

// Start starts the metrics collection in background
func (mc *MetricsCollector) Start(ctx context.Context) error {
	go mc.run(ctx)
	logger.Info("Metrics collector started", zap.Duration("interval", mc.interval))
	return nil
}

// Stop stops the metrics collection
func (mc *MetricsCollector) Stop() error {
	logger.Info("Metrics collector stopped")
	return nil
}

// run is the main loop for metrics collection
func (mc *MetricsCollector) run(ctx context.Context) {
	ticker := time.NewTicker(mc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := mc.collect(); err != nil {
				logger.Error("Failed to collect metrics", zap.Error(err))
			}
		}
	}
}

// collect collects all system metrics
func (mc *MetricsCollector) collect() error {
	metrics := &SystemMetrics{
		Timestamp: time.Now(),
	}

	// Collect CPU metrics
	if err := mc.collectCPU(metrics); err != nil {
		logger.Warn("Failed to collect CPU metrics", zap.Error(err))
	}

	// Collect memory metrics
	if err := mc.collectMemory(metrics); err != nil {
		logger.Warn("Failed to collect memory metrics", zap.Error(err))
	}

	// Collect disk metrics
	if err := mc.collectDisk(metrics); err != nil {
		logger.Warn("Failed to collect disk metrics", zap.Error(err))
	}

	// Collect network metrics
	if err := mc.collectNetwork(metrics); err != nil {
		logger.Warn("Failed to collect network metrics", zap.Error(err))
	}

	// Collect system info
	if err := mc.collectSystemInfo(metrics); err != nil {
		logger.Warn("Failed to collect system info", zap.Error(err))
	}

	// Collect Go runtime metrics
	mc.collectGoRuntime(metrics)

	// Store metrics
	mc.mu.Lock()
	mc.current = metrics
	mc.history = append(mc.history, metrics)
	if len(mc.history) > mc.maxHistory {
		mc.history = mc.history[1:]
	}
	mc.mu.Unlock()

	return nil
}

// collectCPU collects CPU metrics
func (mc *MetricsCollector) collectCPU(metrics *SystemMetrics) error {
	// CPU usage (average across all cores)
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}
	if len(percentages) > 0 {
		metrics.CPUUsagePercent = percentages[0]
	}

	// Per-core usage
	perCorePercentages, err := cpu.Percent(0, true)
	if err == nil {
		metrics.CPUPerCoreUsage = perCorePercentages
		metrics.CPUCores = len(perCorePercentages)
	}

	// Load average
	loadAvg, err := host.LoadAvg()
	if err == nil {
		metrics.LoadAverage1 = loadAvg.Load1
		metrics.LoadAverage5 = loadAvg.Load5
		metrics.LoadAverage15 = loadAvg.Load15
	}

	return nil
}

// collectMemory collects memory metrics
func (mc *MetricsCollector) collectMemory(metrics *SystemMetrics) error {
	// Virtual memory
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	metrics.MemoryTotal = vmem.Total
	metrics.MemoryUsed = vmem.Used
	metrics.MemoryFree = vmem.Free
	metrics.MemoryUsagePercent = vmem.UsedPercent

	// Swap memory
	swap, err := mem.SwapMemory()
	if err == nil {
		metrics.SwapTotal = swap.Total
		metrics.SwapUsed = swap.Used
		metrics.SwapFree = swap.Free
		metrics.SwapUsagePercent = swap.UsedPercent
	}

	return nil
}

// collectDisk collects disk metrics
func (mc *MetricsCollector) collectDisk(metrics *SystemMetrics) error {
	// Root partition
	usage, err := disk.Usage("/")
	if err != nil {
		return err
	}

	metrics.DiskTotal = usage.Total
	metrics.DiskUsed = usage.Used
	metrics.DiskFree = usage.Free
	metrics.DiskUsagePercent = usage.UsedPercent

	return nil
}

// collectNetwork collects network metrics
func (mc *MetricsCollector) collectNetwork(metrics *SystemMetrics) error {
	counters, err := net.IOCounters(false)
	if err != nil {
		return err
	}

	if len(counters) > 0 {
		metrics.NetworkBytesSent = counters[0].BytesSent
		metrics.NetworkBytesRecv = counters[0].BytesRecv
		metrics.NetworkPacketsSent = counters[0].PacketsSent
		metrics.NetworkPacketsRecv = counters[0].PacketsRecv
	}

	return nil
}

// collectSystemInfo collects system info
func (mc *MetricsCollector) collectSystemInfo(metrics *SystemMetrics) error {
	info, err := host.Info()
	if err != nil {
		return err
	}

	metrics.Uptime = info.Uptime
	metrics.BootTime = info.BootTime
	metrics.Processes = info.Procs

	return nil
}

// collectGoRuntime collects Go runtime metrics
func (mc *MetricsCollector) collectGoRuntime(metrics *SystemMetrics) {
	metrics.GoVersion = runtime.Version()
	metrics.GoRoutines = runtime.NumGoroutine()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	metrics.GoMemAlloc = m.Alloc
	metrics.GoMemSys = m.Sys
	metrics.GoGCPauseNs = m.PauseNs[(m.NumGC+255)%256]
}

// GetCurrent returns the current metrics
func (mc *MetricsCollector) GetCurrent() *SystemMetrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.current
}

// GetHistory returns the metrics history
func (mc *MetricsCollector) GetHistory() []*SystemMetrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	// Return a copy to prevent concurrent access issues
	history := make([]*SystemMetrics, len(mc.history))
	copy(history, mc.history)
	return history
}

// GetHistorySince returns metrics since a specific time
func (mc *MetricsCollector) GetHistorySince(since time.Time) []*SystemMetrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	var result []*SystemMetrics
	for _, m := range mc.history {
		if m.Timestamp.After(since) {
			result = append(result, m)
		}
	}
	return result
}

// ToPrometheusFormat converts metrics to Prometheus text format
func (m *SystemMetrics) ToPrometheusFormat() string {
	var output string

	// Helper function to add metric
	addMetric := func(name, metricType, help string, value interface{}) {
		output += fmt.Sprintf("# HELP %s %s\n", name, help)
		output += fmt.Sprintf("# TYPE %s %s\n", name, metricType)
		output += fmt.Sprintf("%s %v\n\n", name, value)
	}

	// CPU Metrics
	addMetric("stumpfworks_cpu_usage_percent", "gauge", "CPU usage percentage", m.CPUUsagePercent)
	addMetric("stumpfworks_cpu_cores", "gauge", "Number of CPU cores", m.CPUCores)
	addMetric("stumpfworks_load_average_1", "gauge", "Load average 1 minute", m.LoadAverage1)
	addMetric("stumpfworks_load_average_5", "gauge", "Load average 5 minutes", m.LoadAverage5)
	addMetric("stumpfworks_load_average_15", "gauge", "Load average 15 minutes", m.LoadAverage15)

	// Memory Metrics
	addMetric("stumpfworks_memory_total_bytes", "gauge", "Total memory in bytes", m.MemoryTotal)
	addMetric("stumpfworks_memory_used_bytes", "gauge", "Used memory in bytes", m.MemoryUsed)
	addMetric("stumpfworks_memory_free_bytes", "gauge", "Free memory in bytes", m.MemoryFree)
	addMetric("stumpfworks_memory_usage_percent", "gauge", "Memory usage percentage", m.MemoryUsagePercent)
	addMetric("stumpfworks_swap_total_bytes", "gauge", "Total swap in bytes", m.SwapTotal)
	addMetric("stumpfworks_swap_used_bytes", "gauge", "Used swap in bytes", m.SwapUsed)
	addMetric("stumpfworks_swap_usage_percent", "gauge", "Swap usage percentage", m.SwapUsagePercent)

	// Disk Metrics
	addMetric("stumpfworks_disk_total_bytes", "gauge", "Total disk space in bytes", m.DiskTotal)
	addMetric("stumpfworks_disk_used_bytes", "gauge", "Used disk space in bytes", m.DiskUsed)
	addMetric("stumpfworks_disk_free_bytes", "gauge", "Free disk space in bytes", m.DiskFree)
	addMetric("stumpfworks_disk_usage_percent", "gauge", "Disk usage percentage", m.DiskUsagePercent)

	// Network Metrics
	addMetric("stumpfworks_network_bytes_sent_total", "counter", "Total bytes sent", m.NetworkBytesSent)
	addMetric("stumpfworks_network_bytes_recv_total", "counter", "Total bytes received", m.NetworkBytesRecv)
	addMetric("stumpfworks_network_packets_sent_total", "counter", "Total packets sent", m.NetworkPacketsSent)
	addMetric("stumpfworks_network_packets_recv_total", "counter", "Total packets received", m.NetworkPacketsRecv)

	// System Info
	addMetric("stumpfworks_uptime_seconds", "counter", "System uptime in seconds", m.Uptime)
	addMetric("stumpfworks_processes_total", "gauge", "Total number of processes", m.Processes)

	// Go Runtime Metrics
	addMetric("stumpfworks_go_goroutines", "gauge", "Number of goroutines", m.GoRoutines)
	addMetric("stumpfworks_go_mem_alloc_bytes", "gauge", "Go memory allocated bytes", m.GoMemAlloc)
	addMetric("stumpfworks_go_mem_sys_bytes", "gauge", "Go memory system bytes", m.GoMemSys)
	addMetric("stumpfworks_go_gc_pause_ns", "gauge", "Go GC pause nanoseconds", m.GoGCPauseNs)

	return output
}

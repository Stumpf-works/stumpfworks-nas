// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package metrics

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	// CollectionInterval is how often to collect metrics
	CollectionInterval = 60 * time.Second
	// MetricsRetention is how long to keep metrics (30 days)
	MetricsRetention = 30 * 24 * time.Hour
	// HealthScoreRetention is how long to keep health scores (90 days)
	HealthScoreRetention = 90 * 24 * time.Hour
)

// Service manages metrics collection and storage
type Service struct {
	db      *gorm.DB
	mu      sync.RWMutex
	running bool
	stop    chan bool

	// Previous values for rate calculations
	prevNetStats  map[string]net.IOCountersStat
	prevDiskStats map[string]disk.IOCountersStat
	prevTime      time.Time
}

var (
	globalService *Service
	once          sync.Once
)

// Initialize initializes the metrics service
func Initialize() (*Service, error) {
	var initErr error
	once.Do(func() {
		db := database.GetDB()
		if db == nil {
			initErr = fmt.Errorf("database not initialized")
			return
		}

		globalService = &Service{
			db:            db,
			stop:          make(chan bool),
			prevNetStats:  make(map[string]net.IOCountersStat),
			prevDiskStats: make(map[string]disk.IOCountersStat),
			prevTime:      time.Now(),
		}

		logger.Info("Metrics service initialized")
	})

	return globalService, initErr
}

// GetService returns the global metrics service
func GetService() *Service {
	if globalService == nil {
		globalService, _ = Initialize()
	}
	return globalService
}

// Start starts the metrics collection
func (s *Service) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("metrics service already running")
	}

	s.running = true
	go s.run()

	logger.Info("Metrics collection started")
	return nil
}

// Stop stops the metrics collection
func (s *Service) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	s.running = false
	s.stop <- true

	logger.Info("Metrics collection stopped")
}

// run is the main metrics collection loop
func (s *Service) run() {
	ticker := time.NewTicker(CollectionInterval)
	defer ticker.Stop()

	// Collect initial metric
	s.collectMetrics()

	for {
		select {
		case <-ticker.C:
			s.collectMetrics()
		case <-s.stop:
			return
		}
	}
}

// collectMetrics collects current system metrics and stores them
func (s *Service) collectMetrics() {
	metric := &models.SystemMetric{
		Timestamp: time.Now(),
	}

	// Collect CPU metrics
	if cpuPercent, err := cpu.Percent(time.Second, false); err == nil && len(cpuPercent) > 0 {
		metric.CPUUsage = cpuPercent[0]
	}

	if loadAvg, err := load.Avg(); err == nil {
		metric.CPULoadAvg1 = loadAvg.Load1
		metric.CPULoadAvg5 = loadAvg.Load5
		metric.CPULoadAvg15 = loadAvg.Load15
	}

	// CPU temperature (may not be available on all systems)
	if temps, err := host.SensorsTemperatures(); err == nil && len(temps) > 0 {
		// Find the highest CPU temperature
		maxTemp := 0.0
		for _, temp := range temps {
			// Look for CPU-related sensors (coretemp, k10temp, etc.)
			if strings.Contains(strings.ToLower(temp.SensorKey), "core") ||
				strings.Contains(strings.ToLower(temp.SensorKey), "cpu") ||
				strings.Contains(strings.ToLower(temp.SensorKey), "package") {
				if temp.Temperature > maxTemp {
					maxTemp = temp.Temperature
				}
			}
		}
		metric.CPUTemperature = maxTemp
	}

	// Collect memory metrics
	if vmem, err := mem.VirtualMemory(); err == nil {
		metric.MemoryUsedBytes = vmem.Used
		metric.MemoryTotalBytes = vmem.Total
		metric.MemoryUsage = vmem.UsedPercent
	}

	if swap, err := mem.SwapMemory(); err == nil {
		metric.SwapUsedBytes = swap.Used
		metric.SwapTotalBytes = swap.Total
		metric.SwapUsage = swap.UsedPercent
	}

	// Collect disk metrics
	s.collectDiskMetrics(metric)

	// Collect network metrics
	s.collectNetworkMetrics(metric)

	// Collect process metrics
	if processes, err := process.Processes(); err == nil {
		metric.ProcessCount = len(processes)
		threadCount := 0
		for _, p := range processes {
			if numThreads, err := p.NumThreads(); err == nil {
				threadCount += int(numThreads)
			}
		}
		metric.ThreadCount = threadCount
	}

	// Store metric
	if err := s.db.Create(metric).Error; err != nil {
		logger.Error("Failed to store metric", zap.Error(err))
		return
	}

	// Calculate and store health score
	s.calculateHealthScore(metric)

	// Cleanup old metrics periodically (every hour)
	if time.Now().Minute() == 0 {
		s.cleanupOldMetrics()
	}
}

// collectDiskMetrics collects disk-related metrics
func (s *Service) collectDiskMetrics(metric *models.SystemMetric) {
	// Get all partitions
	partitions, err := disk.Partitions(false)
	if err != nil {
		return
	}

	var totalUsed, totalSize uint64
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}
		totalUsed += usage.Used
		totalSize += usage.Total
	}

	metric.DiskUsedBytes = totalUsed
	metric.DiskTotalBytes = totalSize
	if totalSize > 0 {
		metric.DiskUsage = float64(totalUsed) / float64(totalSize) * 100
	}

	// Get disk IO stats
	ioCounters, err := disk.IOCounters()
	if err != nil {
		return
	}

	now := time.Now()
	elapsed := now.Sub(s.prevTime).Seconds()
	if elapsed == 0 {
		elapsed = 1
	}

	var totalReadBytes, totalWriteBytes, totalIOPS uint64
	for name, counter := range ioCounters {
		if prev, ok := s.prevDiskStats[name]; ok {
			readDiff := counter.ReadBytes - prev.ReadBytes
			writeDiff := counter.WriteBytes - prev.WriteBytes
			ioDiff := (counter.ReadCount + counter.WriteCount) - (prev.ReadCount + prev.WriteCount)

			totalReadBytes += uint64(float64(readDiff) / elapsed)
			totalWriteBytes += uint64(float64(writeDiff) / elapsed)
			totalIOPS += uint64(float64(ioDiff) / elapsed)
		}
		s.prevDiskStats[name] = counter
	}

	metric.DiskReadBytesPerSec = totalReadBytes
	metric.DiskWriteBytesPerSec = totalWriteBytes
	metric.DiskIOPS = totalIOPS
	s.prevTime = now
}

// collectNetworkMetrics collects network-related metrics
func (s *Service) collectNetworkMetrics(metric *models.SystemMetric) {
	ioCounters, err := net.IOCounters(true)
	if err != nil {
		return
	}

	now := time.Now()
	elapsed := now.Sub(s.prevTime).Seconds()
	if elapsed == 0 {
		elapsed = 1
	}

	var totalRxBytes, totalTxBytes, totalRxPackets, totalTxPackets uint64
	for _, counter := range ioCounters {
		// Skip loopback
		if counter.Name == "lo" {
			continue
		}

		if prev, ok := s.prevNetStats[counter.Name]; ok {
			rxDiff := counter.BytesRecv - prev.BytesRecv
			txDiff := counter.BytesSent - prev.BytesSent
			rxPacketsDiff := counter.PacketsRecv - prev.PacketsRecv
			txPacketsDiff := counter.PacketsSent - prev.PacketsSent

			totalRxBytes += uint64(float64(rxDiff) / elapsed)
			totalTxBytes += uint64(float64(txDiff) / elapsed)
			totalRxPackets += uint64(float64(rxPacketsDiff) / elapsed)
			totalTxPackets += uint64(float64(txPacketsDiff) / elapsed)
		}
		s.prevNetStats[counter.Name] = counter
	}

	metric.NetworkRxBytesPerSec = totalRxBytes
	metric.NetworkTxBytesPerSec = totalTxBytes
	metric.NetworkRxPacketsPerSec = totalRxPackets
	metric.NetworkTxPacketsPerSec = totalTxPackets
}

// calculateHealthScore calculates and stores the system health score
func (s *Service) calculateHealthScore(metric *models.SystemMetric) {
	score := &models.HealthScore{
		Timestamp: metric.Timestamp,
	}

	// CPU score (inverse of usage, load avg consideration)
	cpuScore := 100
	if metric.CPUUsage > 90 {
		cpuScore = 20
	} else if metric.CPUUsage > 75 {
		cpuScore = 50
	} else if metric.CPUUsage > 50 {
		cpuScore = 75
	}
	score.CPUScore = cpuScore

	// Memory score (inverse of usage)
	memoryScore := 100
	if metric.MemoryUsage > 90 {
		memoryScore = 20
	} else if metric.MemoryUsage > 80 {
		memoryScore = 50
	} else if metric.MemoryUsage > 70 {
		memoryScore = 75
	}
	score.MemoryScore = memoryScore

	// Disk score (inverse of usage)
	diskScore := 100
	if metric.DiskUsage > 95 {
		diskScore = 10
	} else if metric.DiskUsage > 90 {
		diskScore = 30
	} else if metric.DiskUsage > 80 {
		diskScore = 60
	} else if metric.DiskUsage > 70 {
		diskScore = 80
	}
	score.DiskScore = diskScore

	// Network score (based on packet errors and utilization)
	networkScore := 100

	// Get current network stats to check for errors
	if netIO, err := net.IOCounters(true); err == nil {
		totalPackets := uint64(0)
		totalErrors := uint64(0)

		for _, io := range netIO {
			// Skip loopback
			if io.Name == "lo" {
				continue
			}
			totalPackets += io.PacketsSent + io.PacketsRecv
			totalErrors += io.Errin + io.Errout + io.Dropin + io.Dropout
		}

		// Calculate error rate
		if totalPackets > 0 {
			errorRate := float64(totalErrors) / float64(totalPackets) * 100

			// Penalize based on error rate
			if errorRate > 5.0 {
				networkScore = 10 // Very high error rate
			} else if errorRate > 2.0 {
				networkScore = 40 // High error rate
			} else if errorRate > 0.5 {
				networkScore = 70 // Moderate error rate
			} else if errorRate > 0.1 {
				networkScore = 90 // Low error rate
			}
			// else: networkScore = 100 (very low/no errors)
		}
	}

	score.NetworkScore = networkScore

	// Overall score (weighted average)
	score.Score = (cpuScore*30 + memoryScore*30 + diskScore*30 + score.NetworkScore*10) / 100

	// Detect issues
	issues := []string{}
	if metric.CPUUsage > 90 {
		issues = append(issues, "High CPU usage")
	}
	if metric.MemoryUsage > 90 {
		issues = append(issues, "High memory usage")
	}
	if metric.DiskUsage > 90 {
		issues = append(issues, "Low disk space")
	}
	if len(issues) > 0 {
		score.Issues = fmt.Sprintf(`["%s"]`, issues[0])
		for i := 1; i < len(issues); i++ {
			score.Issues = fmt.Sprintf(`%s,"%s"`, score.Issues[:len(score.Issues)-1], issues[i])
		}
	}

	// Store health score
	if err := s.db.Create(score).Error; err != nil {
		logger.Error("Failed to store health score", zap.Error(err))
	}
}

// cleanupOldMetrics removes metrics older than the retention period
func (s *Service) cleanupOldMetrics() {
	metricsCutoff := time.Now().Add(-MetricsRetention)
	healthScoreCutoff := time.Now().Add(-HealthScoreRetention)

	// Delete old metrics
	if err := s.db.Where("timestamp < ?", metricsCutoff).Delete(&models.SystemMetric{}).Error; err != nil {
		logger.Error("Failed to cleanup old metrics", zap.Error(err))
	}

	// Delete old health scores
	if err := s.db.Where("timestamp < ?", healthScoreCutoff).Delete(&models.HealthScore{}).Error; err != nil {
		logger.Error("Failed to cleanup old health scores", zap.Error(err))
	}
}

// GetMetrics retrieves metrics within a time range
func (s *Service) GetMetrics(ctx context.Context, start, end time.Time, limit int) ([]models.SystemMetric, error) {
	var metrics []models.SystemMetric

	query := s.db.WithContext(ctx).
		Where("timestamp >= ? AND timestamp <= ?", start, end).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&metrics).Error; err != nil {
		return nil, err
	}

	return metrics, nil
}

// GetHealthScores retrieves health scores within a time range
func (s *Service) GetHealthScores(ctx context.Context, start, end time.Time, limit int) ([]models.HealthScore, error) {
	var scores []models.HealthScore

	query := s.db.WithContext(ctx).
		Where("timestamp >= ? AND timestamp <= ?", start, end).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&scores).Error; err != nil {
		return nil, err
	}

	return scores, nil
}

// GetLatestMetric gets the most recent metric
func (s *Service) GetLatestMetric(ctx context.Context) (*models.SystemMetric, error) {
	var metric models.SystemMetric
	if err := s.db.WithContext(ctx).Order("timestamp DESC").First(&metric).Error; err != nil {
		return nil, err
	}
	return &metric, nil
}

// GetLatestHealthScore gets the most recent health score
func (s *Service) GetLatestHealthScore(ctx context.Context) (*models.HealthScore, error) {
	var score models.HealthScore
	if err := s.db.WithContext(ctx).Order("timestamp DESC").First(&score).Error; err != nil {
		return nil, err
	}
	return &score, nil
}

// GetTrends calculates trends for key metrics
func (s *Service) GetTrends(ctx context.Context, duration time.Duration) ([]models.MetricsTrend, error) {
	now := time.Now()
	recentStart := now.Add(-duration)
	previousStart := now.Add(-duration * 2)
	previousEnd := recentStart

	// Get recent metrics
	var recentMetrics []models.SystemMetric
	if err := s.db.WithContext(ctx).
		Where("timestamp >= ?", recentStart).
		Order("timestamp ASC").
		Find(&recentMetrics).Error; err != nil {
		return nil, err
	}

	// Get previous period metrics
	var previousMetrics []models.SystemMetric
	if err := s.db.WithContext(ctx).
		Where("timestamp >= ? AND timestamp < ?", previousStart, previousEnd).
		Order("timestamp ASC").
		Find(&previousMetrics).Error; err != nil {
		return nil, err
	}

	if len(recentMetrics) == 0 || len(previousMetrics) == 0 {
		return []models.MetricsTrend{}, nil
	}

	// Calculate averages
	recentAvg := calculateAverages(recentMetrics)
	previousAvg := calculateAverages(previousMetrics)

	// Create trends
	trends := []models.MetricsTrend{
		createTrend("CPU Usage", recentAvg.CPUUsage, previousAvg.CPUUsage, now),
		createTrend("Memory Usage", recentAvg.MemoryUsage, previousAvg.MemoryUsage, now),
		createTrend("Disk Usage", recentAvg.DiskUsage, previousAvg.DiskUsage, now),
	}

	return trends, nil
}

func calculateAverages(metrics []models.SystemMetric) models.SystemMetric {
	if len(metrics) == 0 {
		return models.SystemMetric{}
	}

	avg := models.SystemMetric{}
	for _, m := range metrics {
		avg.CPUUsage += m.CPUUsage
		avg.MemoryUsage += m.MemoryUsage
		avg.DiskUsage += m.DiskUsage
	}

	count := float64(len(metrics))
	avg.CPUUsage /= count
	avg.MemoryUsage /= count
	avg.DiskUsage /= count

	return avg
}

func createTrend(name string, current, previous float64, timestamp time.Time) models.MetricsTrend {
	change := current - previous
	changePercent := 0.0
	if previous != 0 {
		changePercent = (change / previous) * 100
	}

	direction := "stable"
	if change > 1 {
		direction = "up"
	} else if change < -1 {
		direction = "down"
	}

	return models.MetricsTrend{
		MetricName:    name,
		CurrentValue:  current,
		PreviousValue: previous,
		Change:        change,
		ChangePercent: changePercent,
		Direction:     direction,
		Timestamp:     timestamp,
	}
}

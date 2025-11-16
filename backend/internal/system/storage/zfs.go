// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package storage

import (
	"time"
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"fmt"
	"strconv"
	"strings"
)

// Shell executor interface (to avoid circular import)


// ZFSManager manages ZFS pools and datasets
type ZFSManager struct {
	shell      executor.ShellExecutor
	enabled bool
}

// ZFSPool represents a ZFS storage pool
type ZFSPool struct {
	Name        string  `json:"name"`
	Size        uint64  `json:"size"`
	Allocated   uint64  `json:"allocated"`
	Free        uint64  `json:"free"`
	Capacity    float64 `json:"capacity"`
	Health      string  `json:"health"`
	Dedup       float64 `json:"dedup"`
	Fragmentation float64 `json:"fragmentation"`
	ReadErrors  uint64  `json:"read_errors"`
	WriteErrors uint64  `json:"write_errors"`
	ChecksumErrors uint64 `json:"checksum_errors"`
}

// ZFSDataset represents a ZFS dataset (filesystem/volume)
type ZFSDataset struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // filesystem, volume, snapshot
	Used        uint64 `json:"used"`
	Available   uint64 `json:"available"`
	Refer       uint64 `json:"refer"`
	Mountpoint  string `json:"mountpoint"`
	Compression string `json:"compression"`
	Dedup       string `json:"dedup"`
	Quota       uint64 `json:"quota"`
	Reservation uint64 `json:"reservation"`
}

// ZFSSnapshot represents a ZFS snapshot
type ZFSSnapshot struct {
	Name     string    `json:"name"`
	Dataset  string    `json:"dataset"`
	Used     uint64    `json:"used"`
	Refer    uint64    `json:"refer"`
	Created  time.Time `json:"created"`
}

// NewZFSManager creates a new ZFS manager
func NewZFSManager(shell executor.ShellExecutor) (*ZFSManager, error) {
	// Check if ZFS is available
	if !shell.CommandExists("zpool") || !shell.CommandExists("zfs") {
		return nil, fmt.Errorf("ZFS tools not installed")
	}

	return &ZFSManager{
		shell:   shell,
		enabled: true,
	}, nil
}

// IsEnabled returns whether ZFS is available
func (z *ZFSManager) IsEnabled() bool {
	return z.enabled
}

// ListPools lists all ZFS pools
func (z *ZFSManager) ListPools() ([]ZFSPool, error) {
	if !z.enabled {
		return nil, fmt.Errorf("ZFS not available")
	}

	result, err := z.shell.Execute("zpool", "list", "-H", "-p")
	if err != nil {
		return nil, fmt.Errorf("failed to list pools: %w", err)
	}

	var pools []ZFSPool
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		pool := ZFSPool{
			Name:   fields[0],
			Health: fields[9],
		}

		// Parse sizes
		if size, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
			pool.Size = size
		}
		if alloc, err := strconv.ParseUint(fields[2], 10, 64); err == nil {
			pool.Allocated = alloc
		}
		if free, err := strconv.ParseUint(fields[3], 10, 64); err == nil {
			pool.Free = free
		}

		// Calculate capacity
		if pool.Size > 0 {
			pool.Capacity = float64(pool.Allocated) / float64(pool.Size) * 100
		}

		// Get additional pool properties
		if err := z.getPoolProperties(&pool); err == nil {
			pools = append(pools, pool)
		}
	}

	return pools, nil
}

// getPoolProperties fetches additional pool properties
func (z *ZFSManager) getPoolProperties(pool *ZFSPool) error {
	result, err := z.shell.Execute("zpool", "get", "-H", "-p", "all", pool.Name)
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		prop := fields[1]
		value := fields[2]

		switch prop {
		case "dedupratio":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				pool.Dedup = f
			}
		case "fragmentation":
			if strings.HasSuffix(value, "%") {
				valueStr := strings.TrimSuffix(value, "%")
				if f, err := strconv.ParseFloat(valueStr, 64); err == nil {
					pool.Fragmentation = f
				}
			}
		}
	}

	return nil
}

// GetPoolStatus gets detailed status of a pool
func (z *ZFSManager) GetPoolStatus(poolName string) (string, error) {
	result, err := z.shell.Execute("zpool", "status", poolName)
	if err != nil {
		return "", fmt.Errorf("failed to get pool status: %w", err)
	}
	return result.Stdout, nil
}

// CreatePool creates a new ZFS pool
func (z *ZFSManager) CreatePool(name string, raidType string, devices []string, opts map[string]string) error {
	if !z.enabled {
		return fmt.Errorf("ZFS not available")
	}

	if len(devices) == 0 {
		return fmt.Errorf("no devices specified")
	}

	args := []string{"create"}

	// Add options
	for key, value := range opts {
		args = append(args, "-o", fmt.Sprintf("%s=%s", key, value))
	}

	args = append(args, name)

	// Add RAID type if specified
	switch raidType {
	case "mirror":
		args = append(args, "mirror")
	case "raidz":
		args = append(args, "raidz")
	case "raidz2":
		args = append(args, "raidz2")
	case "raidz3":
		args = append(args, "raidz3")
	}

	// Add devices
	args = append(args, devices...)

	_, err := z.shell.Execute("zpool", args...)
	if err != nil {
		return fmt.Errorf("failed to create pool: %w", err)
	}

	return nil
}

// DestroyPool destroys a ZFS pool
func (z *ZFSManager) DestroyPool(name string, force bool) error {
	args := []string{"destroy"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, name)

	_, err := z.shell.Execute("zpool", args...)
	if err != nil {
		return fmt.Errorf("failed to destroy pool: %w", err)
	}

	return nil
}

// ScrubPool starts a scrub on a pool
func (z *ZFSManager) ScrubPool(name string) error {
	_, err := z.shell.Execute("zpool", "scrub", name)
	if err != nil {
		return fmt.Errorf("failed to start scrub: %w", err)
	}
	return nil
}

// StopScrub stops a running scrub
func (z *ZFSManager) StopScrub(name string) error {
	_, err := z.shell.Execute("zpool", "scrub", "-s", name)
	if err != nil {
		return fmt.Errorf("failed to stop scrub: %w", err)
	}
	return nil
}

// ListDatasets lists all datasets in a pool
func (z *ZFSManager) ListDatasets(poolName string) ([]ZFSDataset, error) {
	if !z.enabled {
		return nil, fmt.Errorf("ZFS not available")
	}

	args := []string{"list", "-H", "-p", "-r", "-t", "filesystem,volume"}
	if poolName != "" {
		args = append(args, poolName)
	}

	result, err := z.shell.Execute("zfs", args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list datasets: %w", err)
	}

	var datasets []ZFSDataset
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		dataset := ZFSDataset{
			Name:       fields[0],
			Mountpoint: fields[4],
		}

		if used, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
			dataset.Used = used
		}
		if avail, err := strconv.ParseUint(fields[2], 10, 64); err == nil {
			dataset.Available = avail
		}
		if refer, err := strconv.ParseUint(fields[3], 10, 64); err == nil {
			dataset.Refer = refer
		}

		datasets = append(datasets, dataset)
	}

	return datasets, nil
}

// CreateDataset creates a new dataset
func (z *ZFSManager) CreateDataset(name string, opts map[string]string) error {
	args := []string{"create"}

	for key, value := range opts {
		args = append(args, "-o", fmt.Sprintf("%s=%s", key, value))
	}

	args = append(args, name)

	_, err := z.shell.Execute("zfs", args...)
	if err != nil {
		return fmt.Errorf("failed to create dataset: %w", err)
	}

	return nil
}

// DestroyDataset destroys a dataset
func (z *ZFSManager) DestroyDataset(name string, recursive bool) error {
	args := []string{"destroy"}
	if recursive {
		args = append(args, "-r")
	}
	args = append(args, name)

	_, err := z.shell.Execute("zfs", args...)
	if err != nil {
		return fmt.Errorf("failed to destroy dataset: %w", err)
	}

	return nil
}

// CreateSnapshot creates a snapshot
func (z *ZFSManager) CreateSnapshot(datasetName string, snapshotName string) error {
	fullName := fmt.Sprintf("%s@%s", datasetName, snapshotName)
	_, err := z.shell.Execute("zfs", "snapshot", fullName)
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}
	return nil
}

// ListSnapshots lists all snapshots
func (z *ZFSManager) ListSnapshots(datasetName string) ([]ZFSSnapshot, error) {
	args := []string{"list", "-H", "-p", "-t", "snapshot"}
	if datasetName != "" {
		args = append(args, "-r", datasetName)
	}

	result, err := z.shell.Execute("zfs", args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshots: %w", err)
	}

	var snapshots []ZFSSnapshot
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		// Parse dataset@snapshot format
		parts := strings.Split(fields[0], "@")
		if len(parts) != 2 {
			continue
		}

		snapshot := ZFSSnapshot{
			Name:    fields[0],
			Dataset: parts[0],
		}

		if used, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
			snapshot.Used = used
		}
		if refer, err := strconv.ParseUint(fields[3], 10, 64); err == nil {
			snapshot.Refer = refer
		}

		snapshots = append(snapshots, snapshot)
	}

	return snapshots, nil
}

// RollbackSnapshot rolls back to a snapshot
func (z *ZFSManager) RollbackSnapshot(snapshotName string, destroyRecent bool) error {
	args := []string{"rollback"}
	if destroyRecent {
		args = append(args, "-r")
	}
	args = append(args, snapshotName)

	_, err := z.shell.Execute("zfs", args...)
	if err != nil {
		return fmt.Errorf("failed to rollback snapshot: %w", err)
	}

	return nil
}

// DestroySnapshot destroys a snapshot
func (z *ZFSManager) DestroySnapshot(snapshotName string) error {
	_, err := z.shell.Execute("zfs", "destroy", snapshotName)
	if err != nil {
		return fmt.Errorf("failed to destroy snapshot: %w", err)
	}
	return nil
}

// SetProperty sets a property on a dataset or pool
func (z *ZFSManager) SetProperty(name string, property string, value string) error {
	_, err := z.shell.Execute("zfs", "set", fmt.Sprintf("%s=%s", property, value), name)
	if err != nil {
		return fmt.Errorf("failed to set property: %w", err)
	}
	return nil
}

// GetProperty gets a property value
func (z *ZFSManager) GetProperty(name string, property string) (string, error) {
	result, err := z.shell.Execute("zfs", "get", "-H", "-o", "value", property, name)
	if err != nil {
		return "", fmt.Errorf("failed to get property: %w", err)
	}
	return strings.TrimSpace(result.Stdout), nil
}

// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package storage

import "time"

// DiskType represents the type of disk
type DiskType string

const (
	DiskTypeHDD  DiskType = "hdd"
	DiskTypeSSD  DiskType = "ssd"
	DiskTypeNVMe DiskType = "nvme"
	DiskTypeUSB  DiskType = "usb"
)

// DiskStatus represents the health status of a disk
type DiskStatus string

const (
	DiskStatusHealthy  DiskStatus = "healthy"
	DiskStatusWarning  DiskStatus = "warning"
	DiskStatusCritical DiskStatus = "critical"
	DiskStatusFailed   DiskStatus = "failed"
	DiskStatusUnknown  DiskStatus = "unknown"
)

// Disk represents a physical disk
type Disk struct {
	Name         string     `json:"name"`         // e.g., "sda", "nvme0n1"
	Path         string     `json:"path"`         // e.g., "/dev/sda"
	Label        string     `json:"label"`        // User-defined friendly name (optional)
	Model        string     `json:"model"`        // Disk model
	Serial       string     `json:"serial"`       // Serial number
	Size         uint64     `json:"size"`         // Size in bytes
	Type         DiskType   `json:"type"`         // Disk type
	Status       DiskStatus `json:"status"`       // Health status
	Temperature  int        `json:"temperature"`  // Temperature in Celsius
	IsSystem     bool       `json:"isSystem"`     // Is system disk
	IsRemovable  bool       `json:"isRemovable"`  // Is removable
	Partitions   []Partition `json:"partitions"`  // Partitions on this disk
	SMARTEnabled bool       `json:"smartEnabled"` // SMART enabled
	SMART        *SMARTData `json:"smart,omitempty"` // SMART data
}

// Partition represents a disk partition
type Partition struct {
	Name       string `json:"name"`       // e.g., "sda1"
	Path       string `json:"path"`       // e.g., "/dev/sda1"
	Size       uint64 `json:"size"`       // Size in bytes
	Used       uint64 `json:"used"`       // Used space in bytes
	Filesystem string `json:"filesystem"` // Filesystem type (ext4, xfs, etc.)
	MountPoint string `json:"mountPoint"` // Mount point
	Label      string `json:"label"`      // Partition label
	UUID       string `json:"uuid"`       // Partition UUID
	IsMounted  bool   `json:"isMounted"`  // Is currently mounted
}

// SMARTData represents SMART monitoring data
type SMARTData struct {
	Healthy           bool      `json:"healthy"`
	Temperature       int       `json:"temperature"`
	PowerOnHours      uint64    `json:"powerOnHours"`
	PowerCycleCount   uint64    `json:"powerCycleCount"`
	ReallocatedSectors uint64   `json:"reallocatedSectors"`
	PendingSectors    uint64    `json:"pendingSectors"`
	UncorrectableErrors uint64  `json:"uncorrectableErrors"`
	CRCErrors         uint64    `json:"crcErrors"`
	PercentLifeUsed   int       `json:"percentLifeUsed"` // For SSDs
	LastUpdated       time.Time `json:"lastUpdated"`
}

// VolumeType represents the type of volume
type VolumeType string

const (
	VolumeTypeSingle VolumeType = "single"
	VolumeTypeRAID0  VolumeType = "raid0"
	VolumeTypeRAID1  VolumeType = "raid1"
	VolumeTypeRAID5  VolumeType = "raid5"
	VolumeTypeRAID6  VolumeType = "raid6"
	VolumeTypeRAID10 VolumeType = "raid10"
	VolumeTypeLVM    VolumeType = "lvm"
	VolumeTypeZFS    VolumeType = "zfs"
	VolumeTypeBtrfs  VolumeType = "btrfs"
)

// VolumeStatus represents the status of a volume
type VolumeStatus string

const (
	VolumeStatusOnline    VolumeStatus = "online"
	VolumeStatusDegraded  VolumeStatus = "degraded"
	VolumeStatusOffline   VolumeStatus = "offline"
	VolumeStatusRebuilding VolumeStatus = "rebuilding"
	VolumeStatusFailed    VolumeStatus = "failed"
)

// Volume represents a storage volume/pool
type Volume struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Type        VolumeType   `json:"type"`
	Status      VolumeStatus `json:"status"`
	Size        uint64       `json:"size"`        // Total size in bytes
	Used        uint64       `json:"used"`        // Used space in bytes
	Available   uint64       `json:"available"`   // Available space in bytes
	Filesystem  string       `json:"filesystem"`  // Filesystem type
	MountPoint  string       `json:"mountPoint"`  // Mount point
	Disks       []string     `json:"disks"`       // Disk names in this volume
	RaidLevel   string       `json:"raidLevel,omitempty"` // RAID level if applicable
	Health      int          `json:"health"`      // Health percentage (0-100)
	CreatedAt   time.Time    `json:"createdAt"`
	Snapshots   []Snapshot   `json:"snapshots,omitempty"`
}

// Snapshot represents a volume snapshot
type Snapshot struct {
	ID        string    `json:"id"`
	VolumeID  string    `json:"volumeId"`
	Name      string    `json:"name"`
	Size      uint64    `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
}

// ShareType represents the type of network share
type ShareType string

const (
	ShareTypeSMB ShareType = "smb"
	ShareTypeNFS ShareType = "nfs"
	ShareTypeFTP ShareType = "ftp"
)

// Share represents a network share
type Share struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	VolumeID    string    `json:"volumeId,omitempty"` // Optional - linked volume
	Type        ShareType `json:"type"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	ReadOnly    bool      `json:"readOnly"`
	Browseable  bool      `json:"browseable"`
	GuestOK     bool      `json:"guestOk"`
	ValidUsers  []string  `json:"validUsers,omitempty"`
	ValidGroups []string  `json:"validGroups,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// StorageStats represents overall storage statistics
type StorageStats struct {
	TotalDisks      int    `json:"totalDisks"`
	TotalCapacity   uint64 `json:"totalCapacity"`
	UsedCapacity    uint64 `json:"usedCapacity"`
	AvailableCapacity uint64 `json:"availableCapacity"`
	TotalVolumes    int    `json:"totalVolumes"`
	TotalShares     int    `json:"totalShares"`
	HealthyDisks    int    `json:"healthyDisks"`
	WarningDisks    int    `json:"warningDisks"`
	CriticalDisks   int    `json:"criticalDisks"`
}

// DiskIOStats represents disk I/O statistics
type DiskIOStats struct {
	DiskName      string  `json:"diskName"`
	ReadBytes     uint64  `json:"readBytes"`
	WriteBytes    uint64  `json:"writeBytes"`
	ReadOps       uint64  `json:"readOps"`
	WriteOps      uint64  `json:"writeOps"`
	ReadLatency   float64 `json:"readLatency"`  // ms
	WriteLatency  float64 `json:"writeLatency"` // ms
	Utilization   float64 `json:"utilization"`  // Percentage
	Timestamp     time.Time `json:"timestamp"`
}

// CreateVolumeRequest represents a request to create a new volume
type CreateVolumeRequest struct {
	Name       string     `json:"name" validate:"required,min=1,max=255"`
	Type       VolumeType `json:"type" validate:"required"`
	Disks      []string   `json:"disks" validate:"required,min=1"`
	Filesystem string     `json:"filesystem" validate:"required,oneof=ext4 xfs btrfs zfs"`
	MountPoint string     `json:"mountPoint,omitempty"` // Optional - auto-generated from Name if empty
	RaidLevel  string     `json:"raidLevel,omitempty"`
}

// CreateShareRequest represents a request to create a new share
type CreateShareRequest struct {
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	VolumeID    string    `json:"volumeId,omitempty"` // Optional - select from managed volumes
	Path        string    `json:"path,omitempty"`     // Optional - manual path (used if VolumeID not provided)
	Type        ShareType `json:"type" validate:"required,oneof=smb nfs ftp"`
	Description string    `json:"description"`
	ReadOnly    bool      `json:"readOnly"`
	Browseable  bool      `json:"browseable"`
	GuestOK     bool      `json:"guestOk"`
	ValidUsers  []string  `json:"validUsers,omitempty"`
	ValidGroups []string  `json:"validGroups,omitempty"`
}

// FormatDiskRequest represents a request to format a disk/partition
type FormatDiskRequest struct {
	Disk       string `json:"disk" validate:"required"`
	Filesystem string `json:"filesystem" validate:"required,oneof=ext4 xfs btrfs"`
	Label      string `json:"label"`
	Force      bool   `json:"force"` // Force format even if mounted
}

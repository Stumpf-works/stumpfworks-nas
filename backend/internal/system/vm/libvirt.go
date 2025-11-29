// Package vm provides Virtual Machine management via libvirt/KVM
package vm

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// LibvirtManager manages KVM/QEMU virtual machines via libvirt
type LibvirtManager struct {
	shell   executor.ShellExecutor
	enabled bool
}

// VM represents a virtual machine
type VM struct {
	UUID        string   `json:"uuid"`
	Name        string   `json:"name"`
	State       string   `json:"state"`        // running, shutoff, paused
	Memory      int64    `json:"memory"`       // MB
	VCPUs       int      `json:"vcpus"`
	DiskSize    int64    `json:"disk_size"`    // GB
	Autostart   bool     `json:"autostart"`
	OSType      string   `json:"os_type"`      // linux, windows, other
	Architecture string  `json:"architecture"` // x86_64, aarch64
	Disks       []VMDisk `json:"disks"`
	Networks    []VMNetwork `json:"networks"`
}

// VMDisk represents a VM disk
type VMDisk struct {
	Path   string `json:"path"`
	Size   int64  `json:"size"` // GB
	Format string `json:"format"` // qcow2, raw
	Bus    string `json:"bus"`    // virtio, sata, scsi
}

// VMNetwork represents a VM network interface
type VMNetwork struct {
	Type    string `json:"type"`    // bridge, network
	Source  string `json:"source"`  // br0, default
	MAC     string `json:"mac"`
	Model   string `json:"model"`   // virtio, e1000
}

// VMCreateRequest represents a request to create a VM
type VMCreateRequest struct {
	Name         string   `json:"name"`
	Memory       int64    `json:"memory"`        // MB
	VCPUs        int      `json:"vcpus"`
	DiskSize     int64    `json:"disk_size"`     // GB
	DiskFormat   string   `json:"disk_format"`   // qcow2, raw
	OSType       string   `json:"os_type"`       // linux, windows
	ISOPath      string   `json:"iso_path"`      // Optional boot ISO
	Network      string   `json:"network"`       // bridge name or 'default'
	Autostart    bool     `json:"autostart"`
}

// VMSnapshot represents a VM snapshot
type VMSnapshot struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	State       string `json:"state"`
	CreatedAt   string `json:"created_at"`
}

// NewLibvirtManager creates a new libvirt manager
func NewLibvirtManager(shell executor.ShellExecutor) (*LibvirtManager, error) {
	manager := &LibvirtManager{
		shell:   shell,
		enabled: false,
	}

	// Check if virsh is available
	result, err := shell.Execute("which", "virsh")
	if err != nil || result.Stdout == "" {
		logger.Warn("virsh not found, VM features will be disabled")
		return manager, fmt.Errorf("virsh not available: install libvirt-clients package")
	}

	// Check if libvirt is running
	result, err = shell.Execute("systemctl", "is-active", "libvirtd")
	if err != nil || strings.TrimSpace(result.Stdout) != "active" {
		logger.Warn("libvirtd not running, VM features will be disabled")
		return manager, fmt.Errorf("libvirtd not running: start libvirtd service")
	}

	manager.enabled = true
	logger.Info("Libvirt manager initialized successfully")
	return manager, nil
}

// IsEnabled returns whether libvirt is available
func (lm *LibvirtManager) IsEnabled() bool {
	return lm.enabled
}

// ListVMs lists all virtual machines
func (lm *LibvirtManager) ListVMs() ([]VM, error) {
	if !lm.enabled {
		return nil, fmt.Errorf("libvirt is not enabled")
	}

	vms := []VM{}

	// List all VMs (running and shutoff)
	result, err := lm.shell.Execute("virsh", "list", "--all", "--uuid")
	if err != nil {
		return nil, fmt.Errorf("failed to list VMs: %w", err)
	}

	uuids := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	for _, uuid := range uuids {
		uuid = strings.TrimSpace(uuid)
		if uuid == "" {
			continue
		}

		vm, err := lm.GetVM(uuid)
		if err != nil {
			logger.Warn("Failed to get VM details", zap.String("uuid", uuid), zap.Error(err))
			continue
		}

		vms = append(vms, *vm)
	}

	return vms, nil
}

// GetVM gets details of a specific VM
func (lm *LibvirtManager) GetVM(nameOrUUID string) (*VM, error) {
	if !lm.enabled {
		return nil, fmt.Errorf("libvirt is not enabled")
	}

	// Get VM XML description
	result, err := lm.shell.Execute("virsh", "dumpxml", nameOrUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get VM info: %w", err)
	}

	// Parse XML
	vm, err := lm.parseVMXML(result.Stdout)
	if err != nil {
		return nil, fmt.Errorf("failed to parse VM XML: %w", err)
	}

	// Get VM state
	result, err = lm.shell.Execute("virsh", "domstate", nameOrUUID)
	if err == nil {
		vm.State = strings.TrimSpace(result.Stdout)
	}

	// Check autostart
	result, err = lm.shell.Execute("virsh", "dominfo", nameOrUUID)
	if err == nil && strings.Contains(result.Stdout, "Autostart:") {
		vm.Autostart = strings.Contains(result.Stdout, "Autostart:             enable")
	}

	return vm, nil
}

// parseVMXML parses libvirt XML to VM struct
func (lm *LibvirtManager) parseVMXML(xmlData string) (*VM, error) {
	type Domain struct {
		UUID   string `xml:"uuid"`
		Name   string `xml:"name"`
		Memory struct {
			Value int64  `xml:",chardata"`
			Unit  string `xml:"unit,attr"`
		} `xml:"memory"`
		VCPU int    `xml:"vcpu"`
		OS   struct {
			Type struct {
				Arch string `xml:"arch,attr"`
				Name string `xml:",chardata"`
			} `xml:"type"`
		} `xml:"os"`
	}

	var domain Domain
	if err := xml.Unmarshal([]byte(xmlData), &domain); err != nil {
		return nil, err
	}

	// Convert memory to MB
	memory := domain.Memory.Value
	if domain.Memory.Unit == "KiB" {
		memory = memory / 1024
	} else if domain.Memory.Unit == "GiB" {
		memory = memory * 1024
	}

	vm := &VM{
		UUID:         domain.UUID,
		Name:         domain.Name,
		Memory:       memory,
		VCPUs:        domain.VCPU,
		OSType:       domain.OS.Type.Name,
		Architecture: domain.OS.Type.Arch,
		Disks:        []VMDisk{},
		Networks:     []VMNetwork{},
	}

	return vm, nil
}

// CreateVM creates a new virtual machine
func (lm *LibvirtManager) CreateVM(req VMCreateRequest) error {
	if !lm.enabled {
		return fmt.Errorf("libvirt is not enabled")
	}

	// Validate request
	if req.Name == "" {
		return fmt.Errorf("VM name is required")
	}
	if req.Memory == 0 {
		req.Memory = 2048 // Default 2GB
	}
	if req.VCPUs == 0 {
		req.VCPUs = 2 // Default 2 vCPUs
	}
	if req.DiskSize == 0 {
		req.DiskSize = 20 // Default 20GB
	}
	if req.DiskFormat == "" {
		req.DiskFormat = "qcow2"
	}
	if req.OSType == "" {
		req.OSType = "linux"
	}
	if req.Network == "" {
		req.Network = "default"
	}

	// Create disk image
	diskPath := fmt.Sprintf("/var/lib/libvirt/images/%s.%s", req.Name, req.DiskFormat)
	result, err := lm.shell.Execute("qemu-img", "create", "-f", req.DiskFormat, diskPath, fmt.Sprintf("%dG", req.DiskSize))
	if err != nil {
		return fmt.Errorf("failed to create disk image: %s: %w", result.Stderr, err)
	}

	// Build virt-install command
	args := []string{
		"virt-install",
		"--name", req.Name,
		"--memory", strconv.FormatInt(req.Memory, 10),
		"--vcpus", strconv.Itoa(req.VCPUs),
		"--disk", fmt.Sprintf("path=%s,format=%s,bus=virtio", diskPath, req.DiskFormat),
		"--os-variant", req.OSType,
		"--graphics", "vnc,listen=0.0.0.0",
		"--noautoconsole",
	}

	// Add network
	if req.Network == "default" {
		args = append(args, "--network", "network=default,model=virtio")
	} else {
		args = append(args, "--network", fmt.Sprintf("bridge=%s,model=virtio", req.Network))
	}

	// Add ISO if specified
	if req.ISOPath != "" {
		args = append(args, "--cdrom", req.ISOPath)
	} else {
		args = append(args, "--boot", "hd")
		args = append(args, "--import")
	}

	// Create VM
	result, err = lm.shell.Execute(args[0], args[1:]...)
	if err != nil {
		return fmt.Errorf("failed to create VM: %s: %w", result.Stderr, err)
	}

	// Set autostart if requested
	if req.Autostart {
		lm.shell.Execute("virsh", "autostart", req.Name)
	}

	logger.Info("VM created", zap.String("name", req.Name))
	return nil
}

// DeleteVM deletes a virtual machine
func (lm *LibvirtManager) DeleteVM(nameOrUUID string, deleteDisks bool) error {
	if !lm.enabled {
		return fmt.Errorf("libvirt is not enabled")
	}

	// Stop VM if running
	lm.StopVM(nameOrUUID, true)

	// Undefine VM
	args := []string{"virsh", "undefine", nameOrUUID}
	if deleteDisks {
		args = append(args, "--remove-all-storage")
	}

	result, err := lm.shell.Execute(args[0], args[1:]...)
	if err != nil {
		return fmt.Errorf("failed to delete VM: %s: %w", result.Stderr, err)
	}

	logger.Info("VM deleted", zap.String("name", nameOrUUID))
	return nil
}

// StartVM starts a virtual machine
func (lm *LibvirtManager) StartVM(nameOrUUID string) error {
	if !lm.enabled {
		return fmt.Errorf("libvirt is not enabled")
	}

	result, err := lm.shell.Execute("virsh", "start", nameOrUUID)
	if err != nil {
		return fmt.Errorf("failed to start VM: %s: %w", result.Stderr, err)
	}

	logger.Info("VM started", zap.String("name", nameOrUUID))
	return nil
}

// StopVM stops a virtual machine
func (lm *LibvirtManager) StopVM(nameOrUUID string, force bool) error {
	if !lm.enabled {
		return fmt.Errorf("libvirt is not enabled")
	}

	var result *executor.ExecResult
	var err error

	if force {
		result, err = lm.shell.Execute("virsh", "destroy", nameOrUUID)
	} else {
		result, err = lm.shell.Execute("virsh", "shutdown", nameOrUUID)
	}

	if err != nil {
		return fmt.Errorf("failed to stop VM: %s: %w", result.Stderr, err)
	}

	logger.Info("VM stopped", zap.String("name", nameOrUUID), zap.Bool("force", force))
	return nil
}

// RebootVM reboots a virtual machine
func (lm *LibvirtManager) RebootVM(nameOrUUID string) error {
	if !lm.enabled {
		return fmt.Errorf("libvirt is not enabled")
	}

	result, err := lm.shell.Execute("virsh", "reboot", nameOrUUID)
	if err != nil {
		return fmt.Errorf("failed to reboot VM: %s: %w", result.Stderr, err)
	}

	logger.Info("VM rebooted", zap.String("name", nameOrUUID))
	return nil
}

// GetVNCPort gets the VNC port for a VM
func (lm *LibvirtManager) GetVNCPort(nameOrUUID string) (int, error) {
	if !lm.enabled {
		return 0, fmt.Errorf("libvirt is not enabled")
	}

	result, err := lm.shell.Execute("virsh", "vncdisplay", nameOrUUID)
	if err != nil {
		return 0, fmt.Errorf("failed to get VNC port: %w", err)
	}

	// Parse display (e.g., ":0" or "127.0.0.1:0")
	display := strings.TrimSpace(result.Stdout)
	parts := strings.Split(display, ":")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid VNC display: %s", display)
	}

	displayNum := 0
	fmt.Sscanf(parts[len(parts)-1], "%d", &displayNum)

	// VNC port = 5900 + display number
	return 5900 + displayNum, nil
}

// SetAutostart sets VM autostart
func (lm *LibvirtManager) SetAutostart(nameOrUUID string, enabled bool) error {
	if !lm.enabled {
		return fmt.Errorf("libvirt is not enabled")
	}

	var result *executor.ExecResult
	var err error

	if enabled {
		result, err = lm.shell.Execute("virsh", "autostart", nameOrUUID)
	} else {
		result, err = lm.shell.Execute("virsh", "autostart", "--disable", nameOrUUID)
	}

	if err != nil {
		return fmt.Errorf("failed to set autostart: %s: %w", result.Stderr, err)
	}

	logger.Info("VM autostart updated", zap.String("name", nameOrUUID), zap.Bool("enabled", enabled))
	return nil
}

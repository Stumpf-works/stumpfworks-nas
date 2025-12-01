// Package addons provides the addon/plugin system for StumpfWorks NAS
package addons

// Manifest describes an installable addon
type Manifest struct {
	ID          string   `json:"id"`           // Unique identifier (e.g., "vm-manager")
	Name        string   `json:"name"`         // Display name
	Description string   `json:"description"`  // Short description
	Icon        string   `json:"icon"`         // Icon (emoji or path)
	Category    string   `json:"category"`     // virtualization, storage, media, etc.
	Version     string   `json:"version"`      // Addon version
	Author      string   `json:"author"`       // Author name

	// Dependencies
	SystemPackages []string `json:"system_packages"` // apt packages to install
	Services       []string `json:"services"`        // systemd services to enable

	// Installation
	InstallScript   string `json:"install_script"`   // Optional bash script to run on install
	UninstallScript string `json:"uninstall_script"` // Optional bash script to run on uninstall

	// Frontend Integration
	AppComponent string `json:"app_component"` // React component name (if addon has UI)
	RoutePrefix  string `json:"route_prefix"`  // API route prefix (if addon has API)

	// Requirements
	MinimumMemory int64 `json:"minimum_memory"` // MB
	MinimumDisk   int64 `json:"minimum_disk"`   // GB
	Architecture  []string `json:"architecture"` // amd64, arm64

	// Service Management
	RequiresRestart bool `json:"requires_restart"` // Whether service restart is needed after installation
}

// Installation status
type InstallationStatus struct {
	AddonID       string `json:"addon_id"`
	Installed     bool   `json:"installed"`
	Version       string `json:"version"`
	InstallDate   string `json:"install_date"`
	PackagesOK    bool   `json:"packages_ok"`     // All packages installed?
	ServicesOK    bool   `json:"services_ok"`     // All services running?
	Error         string `json:"error,omitempty"` // Installation error if any
}

// Predefined addon manifests
var BuiltinAddons = []Manifest{
	{
		ID:          "vm-manager",
		Name:        "VM Manager",
		Description: "KVM/QEMU virtual machine management with live migration and HA support",
		Icon:        "🖥️",
		Category:    "virtualization",
		Version:     "1.0.0",
		Author:      "StumpfWorks",
		SystemPackages: []string{
			"qemu-kvm",
			"libvirt-daemon-system",
			"libvirt-clients",
			"bridge-utils",
			"virt-manager",
			"qemu-utils",
		},
		Services: []string{
			"libvirtd",
		},
		AppComponent: "VMManager",
		RoutePrefix:  "/api/v1/vms",
		MinimumMemory: 4096, // 4GB
		MinimumDisk:   50,   // 50GB
		Architecture:  []string{"amd64", "arm64"},
		RequiresRestart: true, // Requires restart to initialize VM manager
	},
	{
		ID:          "lxc-manager",
		Name:        "LXC Container Manager",
		Description: "Lightweight Linux container management for efficient workload isolation",
		Icon:        "📦",
		Category:    "virtualization",
		Version:     "1.0.0",
		Author:      "StumpfWorks",
		SystemPackages: []string{
			"lxc",
			"lxc-templates",
			"debootstrap",
		},
		AppComponent: "LXCManager",
		RoutePrefix:  "/api/v1/lxc",
		MinimumMemory: 1024, // 1GB
		MinimumDisk:   10,   // 10GB
		Architecture:  []string{"amd64", "arm64"},
		RequiresRestart: true, // Requires restart to initialize LXC manager
	},
	{
		ID:          "minio",
		Name:        "MinIO S3 Storage",
		Description: "High-performance S3-compatible object storage for cloud-native applications",
		Icon:        "☁️",
		Category:    "storage",
		Version:     "1.0.0",
		Author:      "StumpfWorks",
		SystemPackages: []string{
			"minio",
		},
		Services: []string{
			"minio",
		},
		AppComponent: "MinIOManager",
		RoutePrefix:  "/api/v1/minio",
		MinimumMemory: 2048, // 2GB
		MinimumDisk:   20,   // 20GB
		Architecture:  []string{"amd64", "arm64"},
	},
	{
		ID:          "iscsi-target",
		Name:        "iSCSI Target",
		Description: "Block-level storage sharing via iSCSI protocol for SAN environments",
		Icon:        "🎯",
		Category:    "storage",
		Version:     "1.0.0",
		Author:      "StumpfWorks",
		SystemPackages: []string{
			"tgt",
			"open-iscsi",
		},
		Services: []string{
			"tgt",
		},
		AppComponent: "ISCSIManager",
		RoutePrefix:  "/api/v1/iscsi",
		MinimumMemory: 512,  // 512MB
		MinimumDisk:   5,    // 5GB
		Architecture:  []string{"amd64", "arm64"},
	},
	{
		ID:          "vpn-server",
		Name:        "VPN Server",
		Description: "Multi-protocol VPN server with WireGuard, OpenVPN, PPTP, and L2TP/IPsec support",
		Icon:        "🔒",
		Category:    "networking",
		Version:     "1.0.0",
		Author:      "StumpfWorks",
		SystemPackages: []string{
			// NOTE: Protocols are installed on-demand when user enables them
			// No packages installed during addon installation
		},
		Services: []string{
			// Services are started on-demand per protocol
		},
		AppComponent: "VPNServer",
		RoutePrefix:  "/api/v1/vpn",
		MinimumMemory: 512,  // 512MB
		MinimumDisk:   2,    // 2GB
		Architecture:  []string{"amd64", "arm64"},
		RequiresRestart: false, // No restart needed - protocols are managed dynamically
	},
}

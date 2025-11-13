package dependencies

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// PackageManager represents different Linux package managers
type PackageManager string

const (
	APT  PackageManager = "apt"     // Debian/Ubuntu
	YUM  PackageManager = "yum"     // RHEL/CentOS 7
	DNF  PackageManager = "dnf"     // RHEL/CentOS 8+, Fedora
	PACMAN PackageManager = "pacman" // Arch Linux
	ZYPPER PackageManager = "zypper" // openSUSE
	UNKNOWN PackageManager = "unknown"
)

// Package represents a system package dependency
type Package struct {
	Name         string   // Package name
	Required     bool     // If true, system won't work without it
	CheckCommand string   // Command to check if installed (e.g., "samba --version")
	AptName      string   // Package name in apt (Debian/Ubuntu)
	YumName      string   // Package name in yum/dnf (RHEL/CentOS)
	PacmanName   string   // Package name in pacman (Arch)
	Description  string   // What this package is used for
	Installed    bool     // Current installation status
}

// Checker checks and manages system dependencies
type Checker struct {
	packageManager PackageManager
	packages       []*Package
}

// NewChecker creates a new dependency checker
func NewChecker() *Checker {
	checker := &Checker{
		packageManager: detectPackageManager(),
		packages:       getRequiredPackages(),
	}

	logger.Info("Dependency checker initialized",
		zap.String("packageManager", string(checker.packageManager)),
		zap.String("os", runtime.GOOS),
		zap.String("arch", runtime.GOARCH))

	return checker
}

// detectPackageManager detects which package manager is available
func detectPackageManager() PackageManager {
	// Check in order of preference
	managers := []struct {
		pm      PackageManager
		command string
	}{
		{APT, "apt"},
		{DNF, "dnf"},
		{YUM, "yum"},
		{PACMAN, "pacman"},
		{ZYPPER, "zypper"},
	}

	for _, m := range managers {
		if _, err := exec.LookPath(m.command); err == nil {
			return m.pm
		}
	}

	return UNKNOWN
}

// getRequiredPackages returns list of packages needed by the system
func getRequiredPackages() []*Package {
	return []*Package{
		{
			Name:         "samba",
			Required:     true,
			CheckCommand: "smbd",
			AptName:      "samba",
			YumName:      "samba",
			PacmanName:   "samba",
			Description:  "SMB/CIFS file server (for Windows network drives)",
		},
		{
			Name:         "smbclient",
			Required:     true,
			CheckCommand: "smbclient",
			AptName:      "smbclient",
			YumName:      "samba-client",
			PacmanName:   "smbclient",
			Description:  "Samba client tools (for user management)",
		},
		{
			Name:         "smartmontools",
			Required:     true,
			CheckCommand: "smartctl",
			AptName:      "smartmontools",
			YumName:      "smartmontools",
			PacmanName:   "smartmontools",
			Description:  "SMART disk monitoring tools (for disk health)",
		},
		{
			Name:         "nfs-kernel-server",
			Required:     false,
			CheckCommand: "exportfs",
			AptName:      "nfs-kernel-server",
			YumName:      "nfs-utils",
			PacmanName:   "nfs-utils",
			Description:  "NFS server (for Unix/Linux network shares)",
		},
		{
			Name:         "lvm2",
			Required:     false,
			CheckCommand: "lvm",
			AptName:      "lvm2",
			YumName:      "lvm2",
			PacmanName:   "lvm2",
			Description:  "Logical Volume Manager (for advanced disk management)",
		},
		{
			Name:         "mdadm",
			Required:     false,
			CheckCommand: "mdadm",
			AptName:      "mdadm",
			YumName:      "mdadm",
			PacmanName:   "mdadm",
			Description:  "Software RAID management tool",
		},
		{
			Name:         "docker",
			Required:     false,
			CheckCommand: "docker",
			AptName:      "docker.io",
			YumName:      "docker",
			PacmanName:   "docker",
			Description:  "Container runtime (for Docker management features)",
		},
	}
}

// CheckAll checks all packages and returns status
func (c *Checker) CheckAll() error {
	logger.Info("Checking system dependencies...")

	missingRequired := []string{}
	missingOptional := []string{}

	for _, pkg := range c.packages {
		installed := c.isPackageInstalled(pkg)
		pkg.Installed = installed

		if !installed {
			if pkg.Required {
				missingRequired = append(missingRequired, pkg.Name)
				logger.Warn("Required package missing",
					zap.String("package", pkg.Name),
					zap.String("description", pkg.Description))
			} else {
				missingOptional = append(missingOptional, pkg.Name)
				logger.Info("Optional package missing",
					zap.String("package", pkg.Name),
					zap.String("description", pkg.Description))
			}
		} else {
			logger.Debug("Package installed",
				zap.String("package", pkg.Name))
		}
	}

	if len(missingRequired) > 0 {
		return fmt.Errorf("missing required packages: %s", strings.Join(missingRequired, ", "))
	}

	if len(missingOptional) > 0 {
		logger.Info("Optional packages missing (system will work but some features disabled)",
			zap.Strings("packages", missingOptional))
	}

	return nil
}

// isPackageInstalled checks if a package is installed
func (c *Checker) isPackageInstalled(pkg *Package) bool {
	// First try to find the command
	if pkg.CheckCommand != "" {
		if _, err := exec.LookPath(pkg.CheckCommand); err == nil {
			return true
		}
	}

	// Fallback: check with package manager
	switch c.packageManager {
	case APT:
		return c.checkApt(pkg.AptName)
	case DNF, YUM:
		return c.checkYum(pkg.YumName)
	case PACMAN:
		return c.checkPacman(pkg.PacmanName)
	default:
		// If we can't detect package manager, assume not installed
		return false
	}
}

// checkApt checks if package is installed via apt (Debian/Ubuntu)
func (c *Checker) checkApt(packageName string) bool {
	if packageName == "" {
		return false
	}
	cmd := exec.Command("dpkg", "-l", packageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	// dpkg -l shows "ii" prefix for installed packages
	return strings.Contains(string(output), "ii  "+packageName)
}

// checkYum checks if package is installed via yum/dnf (RHEL/CentOS)
func (c *Checker) checkYum(packageName string) bool {
	if packageName == "" {
		return false
	}
	cmd := exec.Command("rpm", "-q", packageName)
	return cmd.Run() == nil
}

// checkPacman checks if package is installed via pacman (Arch)
func (c *Checker) checkPacman(packageName string) bool {
	if packageName == "" {
		return false
	}
	cmd := exec.Command("pacman", "-Q", packageName)
	return cmd.Run() == nil
}

// GetMissingPackages returns list of missing packages
func (c *Checker) GetMissingPackages() []*Package {
	missing := []*Package{}
	for _, pkg := range c.packages {
		if !pkg.Installed {
			missing = append(missing, pkg)
		}
	}
	return missing
}

// GetMissingRequired returns list of missing required packages
func (c *Checker) GetMissingRequired() []*Package {
	missing := []*Package{}
	for _, pkg := range c.packages {
		if !pkg.Installed && pkg.Required {
			missing = append(missing, pkg)
		}
	}
	return missing
}

// GetInstallCommand returns the command to install missing packages
func (c *Checker) GetInstallCommand() string {
	missing := c.GetMissingPackages()
	if len(missing) == 0 {
		return ""
	}

	var packageNames []string
	for _, pkg := range missing {
		name := c.getPackageName(pkg)
		if name != "" {
			packageNames = append(packageNames, name)
		}
	}

	if len(packageNames) == 0 {
		return ""
	}

	switch c.packageManager {
	case APT:
		return fmt.Sprintf("sudo apt update && sudo apt install -y %s", strings.Join(packageNames, " "))
	case DNF:
		return fmt.Sprintf("sudo dnf install -y %s", strings.Join(packageNames, " "))
	case YUM:
		return fmt.Sprintf("sudo yum install -y %s", strings.Join(packageNames, " "))
	case PACMAN:
		return fmt.Sprintf("sudo pacman -S --noconfirm %s", strings.Join(packageNames, " "))
	case ZYPPER:
		return fmt.Sprintf("sudo zypper install -y %s", strings.Join(packageNames, " "))
	default:
		return fmt.Sprintf("# Unknown package manager - install these packages: %s", strings.Join(packageNames, " "))
	}
}

// getPackageName returns the package name for current package manager
func (c *Checker) getPackageName(pkg *Package) string {
	switch c.packageManager {
	case APT:
		return pkg.AptName
	case DNF, YUM:
		return pkg.YumName
	case PACMAN:
		return pkg.PacmanName
	default:
		return pkg.Name
	}
}

// PrintStatus prints a human-readable status report
func (c *Checker) PrintStatus() {
	fmt.Println("\n=== System Dependencies Status ===")
	fmt.Printf("Package Manager: %s\n", c.packageManager)
	fmt.Println()

	installed := 0
	for _, pkg := range c.packages {
		status := "✗ MISSING"
		if pkg.Installed {
			status = "✓ INSTALLED"
			installed++
		}

		required := "Optional"
		if pkg.Required {
			required = "Required"
		}

		fmt.Printf("  [%s] %s (%s)\n", status, pkg.Name, required)
		fmt.Printf("      %s\n", pkg.Description)
	}

	fmt.Printf("\nStatus: %d/%d packages installed\n", installed, len(c.packages))

	if missing := c.GetMissingPackages(); len(missing) > 0 {
		fmt.Println("\nTo install missing packages, run:")
		fmt.Printf("  %s\n", c.GetInstallCommand())
	}

	fmt.Println()
}

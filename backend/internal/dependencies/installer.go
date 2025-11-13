package dependencies

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// InstallMode defines how to handle missing packages
type InstallMode string

const (
	CheckOnly   InstallMode = "check"      // Only check, don't install
	AutoInstall InstallMode = "auto"       // Automatically install missing packages
	Interactive InstallMode = "interactive" // Ask user before installing
)

// Installer handles automatic installation of missing packages
type Installer struct {
	checker *Checker
	mode    InstallMode
}

// NewInstaller creates a new package installer
func NewInstaller(mode InstallMode) *Installer {
	return &Installer{
		checker: NewChecker(),
		mode:    mode,
	}
}

// CheckAndInstall checks dependencies and installs if needed
func (i *Installer) CheckAndInstall() error {
	// Always check first
	if err := i.checker.CheckAll(); err != nil {
		logger.Warn("Dependency check found issues", zap.Error(err))
	}

	missing := i.checker.GetMissingPackages()
	if len(missing) == 0 {
		logger.Info("All dependencies satisfied")
		return nil
	}

	// Handle based on mode
	switch i.mode {
	case CheckOnly:
		return i.handleCheckOnly(missing)
	case AutoInstall:
		return i.handleAutoInstall(missing)
	case Interactive:
		return i.handleInteractive(missing)
	default:
		return fmt.Errorf("unknown install mode: %s", i.mode)
	}
}

// handleCheckOnly logs missing packages and provides install command
func (i *Installer) handleCheckOnly(missing []*Package) error {
	requiredMissing := []*Package{}
	optionalMissing := []*Package{}

	for _, pkg := range missing {
		if pkg.Required {
			requiredMissing = append(requiredMissing, pkg)
		} else {
			optionalMissing = append(optionalMissing, pkg)
		}
	}

	if len(requiredMissing) > 0 {
		logger.Error("Required dependencies missing - some features will not work",
			zap.Int("count", len(requiredMissing)))
		for _, pkg := range requiredMissing {
			logger.Error("Missing required package",
				zap.String("package", pkg.Name),
				zap.String("description", pkg.Description))
		}
	}

	if len(optionalMissing) > 0 {
		logger.Warn("Optional dependencies missing - some features will be disabled",
			zap.Int("count", len(optionalMissing)))
		for _, pkg := range optionalMissing {
			logger.Warn("Missing optional package",
				zap.String("package", pkg.Name),
				zap.String("description", pkg.Description))
		}
	}

	// Show install command
	if cmd := i.checker.GetInstallCommand(); cmd != "" {
		logger.Info("To install missing packages, run this command:",
			zap.String("command", cmd))
		fmt.Printf("\n╔════════════════════════════════════════════════════════════╗\n")
		fmt.Printf("║  MISSING DEPENDENCIES DETECTED                             ║\n")
		fmt.Printf("╠════════════════════════════════════════════════════════════╣\n")
		fmt.Printf("║  To install missing packages, run:                         ║\n")
		fmt.Printf("║                                                            ║\n")
		fmt.Printf("║  %s\n", cmd)
		fmt.Printf("╚════════════════════════════════════════════════════════════╝\n\n")
	}

	// Return error if required packages are missing
	if len(requiredMissing) > 0 {
		return fmt.Errorf("required dependencies missing: %d packages", len(requiredMissing))
	}

	return nil
}

// handleAutoInstall automatically installs missing packages
func (i *Installer) handleAutoInstall(missing []*Package) error {
	logger.Info("Auto-install mode enabled - attempting to install missing packages",
		zap.Int("count", len(missing)))

	// Check if we have root privileges
	if !isRoot() {
		logger.Error("Auto-install requires root privileges - run with sudo or as root")
		return i.handleCheckOnly(missing) // Fallback to check-only
	}

	// Get install command
	cmd := i.checker.GetInstallCommand()
	if cmd == "" {
		return fmt.Errorf("cannot generate install command")
	}

	logger.Info("Installing packages...", zap.String("command", cmd))

	// Parse and execute command
	if err := i.executeInstallCommand(cmd); err != nil {
		logger.Error("Failed to install packages", zap.Error(err))
		return err
	}

	logger.Info("Package installation completed")

	// Re-check to verify installation
	if err := i.checker.CheckAll(); err != nil {
		logger.Warn("Some packages still missing after installation", zap.Error(err))
		return err
	}

	logger.Info("All dependencies satisfied after installation")
	return nil
}

// handleInteractive asks user before installing
func (i *Installer) handleInteractive(missing []*Package) error {
	fmt.Println("\n=== Missing Dependencies Detected ===")
	fmt.Printf("Found %d missing packages\n\n", len(missing))

	for _, pkg := range missing {
		status := "Optional"
		if pkg.Required {
			status = "REQUIRED"
		}
		fmt.Printf("  - %s (%s): %s\n", pkg.Name, status, pkg.Description)
	}

	fmt.Printf("\nInstall command: %s\n\n", i.checker.GetInstallCommand())
	fmt.Print("Would you like to install missing packages now? [y/N]: ")

	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))

	if response == "y" || response == "yes" {
		// Switch to auto-install mode
		i.mode = AutoInstall
		return i.handleAutoInstall(missing)
	}

	fmt.Println("Skipping installation. You can install manually later.")
	return i.handleCheckOnly(missing)
}

// executeInstallCommand executes the package installation command
func (i *Installer) executeInstallCommand(command string) error {
	// Parse command (remove "sudo" prefix if present, we'll handle it separately)
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	// Remove "sudo" if present (we're already checking for root)
	if parts[0] == "sudo" {
		parts = parts[1:]
	}

	// Handle command with "&&" (e.g., "apt update && apt install")
	if strings.Contains(command, "&&") {
		commands := strings.Split(command, "&&")
		for _, cmd := range commands {
			cmd = strings.TrimSpace(cmd)
			if strings.HasPrefix(cmd, "sudo ") {
				cmd = strings.TrimPrefix(cmd, "sudo ")
			}
			if err := i.executeSingleCommand(cmd); err != nil {
				return err
			}
		}
		return nil
	}

	// Execute single command
	return i.executeSingleCommand(strings.Join(parts, " "))
}

// executeSingleCommand executes a single shell command
func (i *Installer) executeSingleCommand(command string) error {
	logger.Info("Executing command", zap.String("command", command))

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

// isRoot checks if the current process has root privileges
func isRoot() bool {
	return os.Geteuid() == 0
}

// PrintHelp prints usage instructions
func PrintHelp() {
	fmt.Print(`
Dependency Management:

  The system requires several packages to function properly:

  Required Packages:
    - samba: SMB/CIFS file server for Windows network drives
    - smbclient: Samba client tools for user management
    - smartmontools: Disk health monitoring

  Optional Packages:
    - nfs-kernel-server: NFS server for Unix/Linux shares
    - lvm2: Logical Volume Manager for advanced disk management
    - mdadm: Software RAID management
    - docker: Container runtime for Docker features

  Installation Modes:
    check (default): Only check and report missing packages
    auto: Automatically install missing packages (requires root)
    interactive: Ask before installing

  Configuration (config.yaml):
    dependencies:
      checkOnStartup: true
      installMode: "check"  # check | auto | interactive
      autoInstall: false    # Deprecated - use installMode instead
`)
}

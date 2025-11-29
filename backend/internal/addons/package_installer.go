package addons

import (
	"fmt"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// PackageInstaller handles system package installation via apt/dpkg
type PackageInstaller struct {
	shell executor.ShellExecutor
}

// NewPackageInstaller creates a new package installer
func NewPackageInstaller(shell executor.ShellExecutor) *PackageInstaller {
	return &PackageInstaller{
		shell: shell,
	}
}

// InstallPackages installs a list of system packages
func (pi *PackageInstaller) InstallPackages(packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	logger.Info("Installing packages", zap.Strings("packages", packages))

	// Update package lists first
	result, err := pi.shell.Execute("apt-get", "update")
	if err != nil {
		logger.Error("Failed to update package lists", zap.Error(err), zap.String("stderr", result.Stderr))
		return fmt.Errorf("failed to update package lists: %w", err)
	}

	// Install packages
	args := []string{"apt-get", "install", "-y"}
	args = append(args, packages...)

	result, err = pi.shell.Execute(args[0], args[1:]...)
	if err != nil {
		logger.Error("Failed to install packages", zap.Error(err), zap.String("stderr", result.Stderr))
		return fmt.Errorf("failed to install packages: %s: %w", result.Stderr, err)
	}

	logger.Info("Packages installed successfully", zap.Strings("packages", packages))
	return nil
}

// UninstallPackages uninstalls a list of system packages
func (pi *PackageInstaller) UninstallPackages(packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	logger.Info("Uninstalling packages", zap.Strings("packages", packages))

	args := []string{"apt-get", "remove", "-y"}
	args = append(args, packages...)

	result, err := pi.shell.Execute(args[0], args[1:]...)
	if err != nil {
		logger.Error("Failed to uninstall packages", zap.Error(err), zap.String("stderr", result.Stderr))
		return fmt.Errorf("failed to uninstall packages: %s: %w", result.Stderr, err)
	}

	logger.Info("Packages uninstalled successfully", zap.Strings("packages", packages))
	return nil
}

// IsPackageInstalled checks if a package is installed
func (pi *PackageInstaller) IsPackageInstalled(packageName string) bool {
	result, err := pi.shell.Execute("dpkg", "-s", packageName)
	if err != nil {
		return false
	}

	// Check if package is installed and not just config-files
	return strings.Contains(result.Stdout, "Status: install ok installed")
}

// AreAllPackagesInstalled checks if all packages in a list are installed
func (pi *PackageInstaller) AreAllPackagesInstalled(packages []string) bool {
	for _, pkg := range packages {
		if !pi.IsPackageInstalled(pkg) {
			logger.Debug("Package not installed", zap.String("package", pkg))
			return false
		}
	}
	return true
}

// EnableService enables and starts a systemd service
func (pi *PackageInstaller) EnableService(serviceName string) error {
	logger.Info("Enabling service", zap.String("service", serviceName))

	// Enable service
	result, err := pi.shell.Execute("systemctl", "enable", serviceName)
	if err != nil {
		logger.Warn("Failed to enable service", zap.String("service", serviceName), zap.Error(err))
		// Continue anyway - might already be enabled
	}

	// Start service
	result, err = pi.shell.Execute("systemctl", "start", serviceName)
	if err != nil {
		logger.Error("Failed to start service", zap.String("service", serviceName), zap.Error(err), zap.String("stderr", result.Stderr))
		return fmt.Errorf("failed to start service %s: %w", serviceName, err)
	}

	logger.Info("Service enabled and started", zap.String("service", serviceName))
	return nil
}

// DisableService stops and disables a systemd service
func (pi *PackageInstaller) DisableService(serviceName string) error {
	logger.Info("Disabling service", zap.String("service", serviceName))

	// Stop service
	result, err := pi.shell.Execute("systemctl", "stop", serviceName)
	if err != nil {
		logger.Warn("Failed to stop service", zap.String("service", serviceName), zap.Error(err))
		// Continue anyway
	}

	// Disable service
	result, err = pi.shell.Execute("systemctl", "disable", serviceName)
	if err != nil {
		logger.Error("Failed to disable service", zap.String("service", serviceName), zap.Error(err), zap.String("stderr", result.Stderr))
		return fmt.Errorf("failed to disable service %s: %w", serviceName, err)
	}

	logger.Info("Service disabled", zap.String("service", serviceName))
	return nil
}

// IsServiceRunning checks if a systemd service is running
func (pi *PackageInstaller) IsServiceRunning(serviceName string) bool {
	result, err := pi.shell.Execute("systemctl", "is-active", serviceName)
	if err != nil {
		return false
	}

	return strings.TrimSpace(result.Stdout) == "active"
}

// AreAllServicesRunning checks if all services in a list are running
func (pi *PackageInstaller) AreAllServicesRunning(services []string) bool {
	for _, svc := range services {
		if !pi.IsServiceRunning(svc) {
			logger.Debug("Service not running", zap.String("service", svc))
			return false
		}
	}
	return true
}

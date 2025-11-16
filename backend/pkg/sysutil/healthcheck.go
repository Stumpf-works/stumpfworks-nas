// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package sysutil

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// SystemCheck represents the result of a single system component check
type SystemCheck struct {
	Name        string    `json:"name"`
	Required    bool      `json:"required"`
	Installed   bool      `json:"installed"`
	Version     string    `json:"version,omitempty"`
	Path        string    `json:"path,omitempty"`
	Status      string    `json:"status"` // ok, warning, error, missing
	Message     string    `json:"message,omitempty"`
	CheckedAt   time.Time `json:"checkedAt"`
}

// SystemHealthReport contains all system checks
type SystemHealthReport struct {
	OverallStatus string        `json:"overallStatus"` // healthy, degraded, unhealthy
	CheckedAt     time.Time     `json:"checkedAt"`
	Hostname      string        `json:"hostname"`
	OS            string        `json:"os"`
	Checks        []SystemCheck `json:"checks"`
	Summary       HealthSummary `json:"summary"`
}

// HealthSummary provides a quick overview
type HealthSummary struct {
	TotalChecks    int `json:"totalChecks"`
	Passed         int `json:"passed"`
	Warnings       int `json:"warnings"`
	Errors         int `json:"errors"`
	Missing        int `json:"missing"`
	RequiredMissing int `json:"requiredMissing"`
}

// ComponentDefinition defines what to check for each component
type ComponentDefinition struct {
	Name        string
	Command     string
	Required    bool
	VersionFlag string
	ServiceName string // for systemd service checks
}

// Standard components to check
var standardComponents = []ComponentDefinition{
	// Samba (SMB shares)
	{Name: "Samba (smbd)", Command: "smbd", Required: false, VersionFlag: "--version", ServiceName: "smbd"},
	{Name: "Samba (nmbd)", Command: "nmbd", Required: false, ServiceName: "nmbd"},
	{Name: "smbpasswd", Command: "smbpasswd", Required: false},
	{Name: "pdbedit", Command: "pdbedit", Required: false},
	{Name: "testparm", Command: "testparm", Required: false},

	// NFS
	{Name: "NFS exportfs", Command: "exportfs", Required: false},
	{Name: "NFS rpcbind", Command: "rpcbind", Required: false, ServiceName: "rpcbind"},

	// System utilities
	{Name: "useradd", Command: "useradd", Required: true},
	{Name: "userdel", Command: "userdel", Required: true},
	{Name: "usermod", Command: "usermod", Required: true},
	{Name: "groupadd", Command: "groupadd", Required: true},
	{Name: "chown", Command: "chown", Required: true},
	{Name: "chmod", Command: "chmod", Required: true},

	// Disk management
	{Name: "lsblk", Command: "lsblk", Required: true},
	{Name: "fdisk", Command: "fdisk", Required: false},
	{Name: "parted", Command: "parted", Required: false},
	{Name: "mkfs.ext4", Command: "mkfs.ext4", Required: false},
	{Name: "mkfs.xfs", Command: "mkfs.xfs", Required: false},
	{Name: "mkfs.btrfs", Command: "mkfs.btrfs", Required: false},

	// Monitoring
	{Name: "smartctl", Command: "smartctl", Required: false},
	{Name: "iostat", Command: "iostat", Required: false},

	// Systemd
	{Name: "systemctl", Command: "systemctl", Required: false},
}

// PerformSystemHealthCheck runs all system checks
func PerformSystemHealthCheck() *SystemHealthReport {
	now := time.Now()
	report := &SystemHealthReport{
		CheckedAt: now,
		Checks:    make([]SystemCheck, 0, len(standardComponents)),
	}

	// Get hostname
	if hostname, err := RunCommand("hostname"); err == nil {
		report.Hostname = strings.TrimSpace(hostname)
	}

	// Get OS info
	if osInfo, err := RunCommand("uname", "-sr"); err == nil {
		report.OS = strings.TrimSpace(osInfo)
	}

	// Perform checks for each component
	for _, component := range standardComponents {
		check := checkComponent(component, now)
		report.Checks = append(report.Checks, check)
	}

	// Calculate summary
	report.Summary = calculateSummary(report.Checks)

	// Determine overall status
	if report.Summary.RequiredMissing > 0 {
		report.OverallStatus = "unhealthy"
	} else if report.Summary.Warnings > 0 || report.Summary.Errors > 0 {
		report.OverallStatus = "degraded"
	} else {
		report.OverallStatus = "healthy"
	}

	return report
}

// checkComponent performs a check for a single component
func checkComponent(def ComponentDefinition, now time.Time) SystemCheck {
	check := SystemCheck{
		Name:      def.Name,
		Required:  def.Required,
		CheckedAt: now,
	}

	// Check if command exists
	path := FindCommand(def.Command)
	if path == def.Command {
		// Not found in system paths
		check.Installed = false
		if def.Required {
			check.Status = "error"
			check.Message = fmt.Sprintf("Required component not found: %s", def.Command)
		} else {
			check.Status = "missing"
			check.Message = fmt.Sprintf("Optional component not installed: %s", def.Command)
		}
		return check
	}

	// Command found
	check.Installed = true
	check.Path = path

	// Try to get version if flag is specified
	if def.VersionFlag != "" {
		if version, err := getVersion(path, def.VersionFlag); err == nil {
			check.Version = version
		}
	}

	// Check service status if applicable
	if def.ServiceName != "" {
		serviceStatus := checkServiceStatus(def.ServiceName)
		check.Status = serviceStatus.Status
		check.Message = serviceStatus.Message
	} else {
		check.Status = "ok"
		check.Message = "Component installed and accessible"
	}

	return check
}

// serviceStatus represents the status of a systemd service
type serviceStatus struct {
	Status  string
	Message string
}

// checkServiceStatus checks if a systemd service is running
func checkServiceStatus(serviceName string) serviceStatus {
	// Check if systemctl is available
	if !CommandExists("systemctl") {
		return serviceStatus{
			Status:  "warning",
			Message: "systemctl not available - cannot check service status",
		}
	}

	// Check service status
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := cmd.Output()
	status := strings.TrimSpace(string(output))

	if err == nil && status == "active" {
		return serviceStatus{
			Status:  "ok",
			Message: fmt.Sprintf("Service %s is running", serviceName),
		}
	}

	// Service not running - check if it exists
	cmd = exec.Command("systemctl", "status", serviceName)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 4 {
				return serviceStatus{
					Status:  "warning",
					Message: fmt.Sprintf("Service %s not found (not installed?)", serviceName),
				}
			}
		}
	}

	return serviceStatus{
		Status:  "warning",
		Message: fmt.Sprintf("Service %s is not running (status: %s)", serviceName, status),
	}
}

// getVersion tries to get version information from a command
func getVersion(path, versionFlag string) (string, error) {
	cmd := exec.Command(path, versionFlag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// Get first line of output
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}

	return "", fmt.Errorf("no version output")
}

// calculateSummary calculates the summary statistics
func calculateSummary(checks []SystemCheck) HealthSummary {
	summary := HealthSummary{
		TotalChecks: len(checks),
	}

	for _, check := range checks {
		switch check.Status {
		case "ok":
			summary.Passed++
		case "warning":
			summary.Warnings++
		case "error":
			summary.Errors++
			if check.Required {
				summary.RequiredMissing++
			}
		case "missing":
			summary.Missing++
			if check.Required {
				summary.RequiredMissing++
			}
		}
	}

	return summary
}

// ToJSON converts the report to JSON string
func (r *SystemHealthReport) ToJSON() (string, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJSONCompact converts the report to compact JSON string
func (r *SystemHealthReport) ToJSONCompact() (string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// PrintReport prints a human-readable report
func (r *SystemHealthReport) PrintReport() {
	fmt.Println("=====================================")
	fmt.Println("  System Health Check Report")
	fmt.Println("=====================================")
	fmt.Printf("Hostname:       %s\n", r.Hostname)
	fmt.Printf("OS:             %s\n", r.OS)
	fmt.Printf("Overall Status: %s\n", r.OverallStatus)
	fmt.Printf("Checked At:     %s\n", r.CheckedAt.Format(time.RFC3339))
	fmt.Println()
	fmt.Println("Summary:")
	fmt.Printf("  Total Checks:     %d\n", r.Summary.TotalChecks)
	fmt.Printf("  ✓ Passed:         %d\n", r.Summary.Passed)
	fmt.Printf("  ⚠ Warnings:       %d\n", r.Summary.Warnings)
	fmt.Printf("  ✗ Errors:         %d\n", r.Summary.Errors)
	fmt.Printf("  - Missing:        %d\n", r.Summary.Missing)
	fmt.Printf("  ! Required Missing: %d\n", r.Summary.RequiredMissing)
	fmt.Println()
	fmt.Println("Component Details:")
	fmt.Println("-------------------------------------")

	for _, check := range r.Checks {
		statusSymbol := getStatusSymbol(check.Status)
		required := ""
		if check.Required {
			required = " [REQUIRED]"
		}

		fmt.Printf("%s %s%s\n", statusSymbol, check.Name, required)
		if check.Installed {
			fmt.Printf("    Path: %s\n", check.Path)
			if check.Version != "" {
				fmt.Printf("    Version: %s\n", check.Version)
			}
		}
		if check.Message != "" {
			fmt.Printf("    Status: %s\n", check.Message)
		}
		fmt.Println()
	}

	fmt.Println("=====================================")
}

// getStatusSymbol returns an emoji/symbol for the status
func getStatusSymbol(status string) string {
	switch status {
	case "ok":
		return "✓"
	case "warning":
		return "⚠"
	case "error":
		return "✗"
	case "missing":
		return "-"
	default:
		return "?"
	}
}

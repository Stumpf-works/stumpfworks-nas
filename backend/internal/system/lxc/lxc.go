// Package lxc provides LXC container management
package lxc

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// LXCManager manages LXC containers
type LXCManager struct {
	shell   executor.ShellExecutor
	enabled bool
}

// Container represents an LXC container
type Container struct {
	Name        string `json:"name"`
	State       string `json:"state"`       // RUNNING, STOPPED, FROZEN
	PID         int    `json:"pid"`
	Memory      int64  `json:"memory"`      // MB
	MemoryLimit int64  `json:"memory_limit"` // MB
	CPUUsage    float64 `json:"cpu_usage"`   // Percentage
	IPv4        string `json:"ipv4"`
	IPv6        string `json:"ipv6"`
	Autostart   bool   `json:"autostart"`
	Template    string `json:"template"`
}

// ContainerCreateRequest represents a request to create a container
type ContainerCreateRequest struct {
	Name        string `json:"name"`
	Template    string `json:"template"`    // ubuntu, debian, alpine, etc.
	Release     string `json:"release"`     // 22.04, bullseye, 3.18, etc.
	Architecture string `json:"architecture"` // amd64, arm64
	MemoryLimit int64  `json:"memory_limit"` // MB
	CPULimit    int    `json:"cpu_limit"`    // Number of CPUs
	Autostart   bool   `json:"autostart"`
	NetworkMode string `json:"network_mode"` // "internal" (lxcbr0) or "bridged" (custom bridge)
	Bridge      string `json:"bridge"`       // Bridge name when network_mode is "bridged" (e.g., br0, vmbr0, vmbr1)
	Password    string `json:"password"`     // Root password for SSH access
	SSHKey      string `json:"ssh_key"`      // SSH public key for passwordless authentication
}

// Template represents an LXC template
type Template struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Releases     []string `json:"releases"`
	Architectures []string `json:"architectures"`
}

// NewLXCManager creates a new LXC manager
func NewLXCManager(shell executor.ShellExecutor) (*LXCManager, error) {
	manager := &LXCManager{
		shell:   shell,
		enabled: false,
	}

	// Check if lxc-ls is available
	result, err := shell.Execute("which", "lxc-ls")
	if err != nil || result.Stdout == "" {
		logger.Warn("lxc-ls not found, LXC features will be disabled")
		return manager, fmt.Errorf("lxc not available: install lxc package")
	}

	manager.enabled = true
	logger.Info("LXC manager initialized successfully")
	return manager, nil
}

// IsEnabled returns whether LXC is available
func (lm *LXCManager) IsEnabled() bool {
	return lm.enabled
}

// ListContainers lists all LXC containers
func (lm *LXCManager) ListContainers() ([]Container, error) {
	if !lm.enabled {
		return nil, fmt.Errorf("LXC is not enabled")
	}

	containers := []Container{}

	// List all containers with details
	result, err := lm.shell.Execute("lxc-ls", "-f")
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	lines := strings.Split(result.Stdout, "\n")
	for i, line := range lines {
		// Skip header and empty lines
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		container := lm.parseContainerLine(line)
		if container != nil {
			containers = append(containers, *container)
		}
	}

	return containers, nil
}

// parseContainerLine parses a line from lxc-ls -f output
func (lm *LXCManager) parseContainerLine(line string) *Container {
	// Example line: "mycontainer  RUNNING  0  10.0.3.123  -  false"
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil
	}

	container := &Container{
		Name:  fields[0],
		State: fields[1],
	}

	if len(fields) > 2 && fields[2] != "-" {
		fmt.Sscanf(fields[2], "%d", &container.PID)
	}

	if len(fields) > 3 && fields[3] != "-" {
		container.IPv4 = fields[3]
	}

	if len(fields) > 4 && fields[4] != "-" {
		container.IPv6 = fields[4]
	}

	if len(fields) > 5 {
		container.Autostart = fields[5] == "true" || fields[5] == "YES"
	}

	return container
}

// GetContainer gets details of a specific container
func (lm *LXCManager) GetContainer(name string) (*Container, error) {
	if !lm.enabled {
		return nil, fmt.Errorf("LXC is not enabled")
	}

	// Get container info
	result, err := lm.shell.Execute("lxc-info", "-n", name)
	if err != nil {
		return nil, fmt.Errorf("failed to get container info: %w", err)
	}

	container := &Container{
		Name: name,
	}

	// Parse lxc-info output
	lines := strings.Split(result.Stdout, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "State:") {
			container.State = strings.TrimSpace(strings.TrimPrefix(line, "State:"))
		} else if strings.HasPrefix(line, "PID:") {
			pidStr := strings.TrimSpace(strings.TrimPrefix(line, "PID:"))
			fmt.Sscanf(pidStr, "%d", &container.PID)
		} else if strings.HasPrefix(line, "IP:") {
			ip := strings.TrimSpace(strings.TrimPrefix(line, "IP:"))
			if strings.Contains(ip, ":") {
				container.IPv6 = ip
			} else {
				container.IPv4 = ip
			}
		} else if strings.HasPrefix(line, "Memory use:") {
			memStr := strings.TrimSpace(strings.TrimPrefix(line, "Memory use:"))
			// Parse memory (e.g., "128.50 MiB")
			memRegex := regexp.MustCompile(`([\d.]+)\s*(\w+)`)
			if matches := memRegex.FindStringSubmatch(memStr); len(matches) > 2 {
				var mem float64
				fmt.Sscanf(matches[1], "%f", &mem)
				if matches[2] == "GiB" {
					mem = mem * 1024
				} else if matches[2] == "KiB" {
					mem = mem / 1024
				}
				container.Memory = int64(mem)
			}
		}
	}

	// Check autostart
	result, err = lm.shell.Execute("cat", fmt.Sprintf("/var/lib/lxc/%s/config", name))
	if err == nil && strings.Contains(result.Stdout, "lxc.start.auto = 1") {
		container.Autostart = true
	}

	return container, nil
}

// CreateContainer creates a new LXC container
func (lm *LXCManager) CreateContainer(req ContainerCreateRequest) error {
	if !lm.enabled {
		return fmt.Errorf("LXC is not enabled")
	}

	// Validate request
	if req.Name == "" {
		return fmt.Errorf("container name is required")
	}
	if req.Template == "" {
		req.Template = "ubuntu"
	}
	if req.Release == "" {
		req.Release = "22.04"
	}
	if req.Architecture == "" {
		req.Architecture = "amd64"
	}

	// Build lxc-create command
	// Modern LXC uses the "download" template
	args := []string{
		"lxc-create",
		"-n", req.Name,
		"-t", "download",
		"--",
		"--dist", req.Template,
		"--release", req.Release,
		"--arch", req.Architecture,
	}

	// Container creation can take several minutes (downloading packages, etc.)
	// Use a 10-minute timeout instead of the default 30 seconds
	result, err := lm.shell.ExecuteWithTimeout(10*time.Minute, args[0], args[1:]...)
	if err != nil {
		return fmt.Errorf("failed to create container: %s: %w", result.Stderr, err)
	}

	// Set resource limits if specified
	configPath := fmt.Sprintf("/var/lib/lxc/%s/config", req.Name)

	if req.MemoryLimit > 0 {
		limitStr := fmt.Sprintf("lxc.cgroup2.memory.max = %dM", req.MemoryLimit)
		lm.shell.Execute("sh", "-c", fmt.Sprintf("echo '%s' >> %s", limitStr, configPath))
	}

	if req.CPULimit > 0 {
		limitStr := fmt.Sprintf("lxc.cgroup2.cpu.max = %d00000 100000", req.CPULimit)
		lm.shell.Execute("sh", "-c", fmt.Sprintf("echo '%s' >> %s", limitStr, configPath))
	}

	// Set autostart if requested
	if req.Autostart {
		lm.shell.Execute("sh", "-c", fmt.Sprintf("echo 'lxc.start.auto = 1' >> %s", configPath))
	}

	// Configure network mode
	if req.NetworkMode == "" {
		req.NetworkMode = "internal"
	}

	// Remove default network configuration and add custom one
	lm.shell.Execute("sh", "-c", fmt.Sprintf("sed -i '/lxc.net.0/d' %s", configPath))

	if req.NetworkMode == "bridged" {
		// Use custom bridge or default to br0
		bridge := req.Bridge
		if bridge == "" {
			bridge = "br0"
		}

		// Configure bridged network for DHCP from router
		lm.shell.Execute("sh", "-c", fmt.Sprintf("echo 'lxc.net.0.type = veth' >> %s", configPath))
		lm.shell.Execute("sh", "-c", fmt.Sprintf("echo 'lxc.net.0.link = %s' >> %s", bridge, configPath))
		lm.shell.Execute("sh", "-c", fmt.Sprintf("echo 'lxc.net.0.flags = up' >> %s", configPath))
		lm.shell.Execute("sh", "-c", fmt.Sprintf("echo 'lxc.net.0.hwaddr = 00:16:3e:xx:xx:xx' >> %s", configPath))
		logger.Info("Container configured with bridged network", zap.String("name", req.Name), zap.String("bridge", bridge))
	} else {
		// Configure internal network (lxcbr0)
		lm.shell.Execute("sh", "-c", fmt.Sprintf("echo 'lxc.net.0.type = veth' >> %s", configPath))
		lm.shell.Execute("sh", "-c", fmt.Sprintf("echo 'lxc.net.0.link = lxcbr0' >> %s", configPath))
		lm.shell.Execute("sh", "-c", fmt.Sprintf("echo 'lxc.net.0.flags = up' >> %s", configPath))
		lm.shell.Execute("sh", "-c", fmt.Sprintf("echo 'lxc.net.0.hwaddr = 00:16:3e:xx:xx:xx' >> %s", configPath))
		logger.Info("Container configured with internal network", zap.String("name", req.Name))
	}

	logger.Info("Container created", zap.String("name", req.Name))

	// Start container temporarily to configure password and SSH key
	if req.Password != "" || req.SSHKey != "" {
		logger.Info("Configuring security settings for container", zap.String("name", req.Name))

		// Start container
		if err := lm.StartContainer(req.Name); err != nil {
			logger.Warn("Failed to start container for security configuration", zap.Error(err))
			return nil // Don't fail the entire creation if we can't configure security
		}

		// Wait for container to fully start
		time.Sleep(3 * time.Second)

		// Configure root password
		if req.Password != "" {
			logger.Info("Setting root password", zap.String("name", req.Name))
			// Use echo to pass password to chpasswd
			passwdCmd := fmt.Sprintf("echo 'root:%s' | chpasswd", req.Password)
			_, err := lm.shell.Execute("lxc-attach", "-n", req.Name, "--", "sh", "-c", passwdCmd)
			if err != nil {
				logger.Warn("Failed to set root password", zap.Error(err), zap.String("name", req.Name))
			} else {
				logger.Info("Root password set successfully", zap.String("name", req.Name))
			}
		}

		// Configure SSH key
		if req.SSHKey != "" {
			logger.Info("Configuring SSH key", zap.String("name", req.Name))

			// Create .ssh directory
			_, err := lm.shell.Execute("lxc-attach", "-n", req.Name, "--", "sh", "-c", "mkdir -p /root/.ssh && chmod 700 /root/.ssh")
			if err != nil {
				logger.Warn("Failed to create .ssh directory", zap.Error(err), zap.String("name", req.Name))
			} else {
				// Add SSH key to authorized_keys
				sshCmd := fmt.Sprintf("echo '%s' >> /root/.ssh/authorized_keys && chmod 600 /root/.ssh/authorized_keys", req.SSHKey)
				_, err = lm.shell.Execute("lxc-attach", "-n", req.Name, "--", "sh", "-c", sshCmd)
				if err != nil {
					logger.Warn("Failed to configure SSH key", zap.Error(err), zap.String("name", req.Name))
				} else {
					logger.Info("SSH key configured successfully", zap.String("name", req.Name))
				}
			}

			// Ensure SSH server is installed and enabled (for Ubuntu/Debian)
			if req.Template == "ubuntu" || req.Template == "debian" {
				logger.Info("Installing and enabling SSH server", zap.String("name", req.Name))
				lm.shell.Execute("lxc-attach", "-n", req.Name, "--", "sh", "-c", "apt-get update && apt-get install -y openssh-server")
				lm.shell.Execute("lxc-attach", "-n", req.Name, "--", "sh", "-c", "systemctl enable ssh && systemctl start ssh")
			}
		}

		logger.Info("Security configuration completed", zap.String("name", req.Name))
	}

	return nil
}

// DeleteContainer deletes an LXC container
func (lm *LXCManager) DeleteContainer(name string) error {
	if !lm.enabled {
		return fmt.Errorf("LXC is not enabled")
	}

	// Stop container if running
	lm.StopContainer(name, true)

	// Destroy container
	result, err := lm.shell.Execute("lxc-destroy", "-n", name, "-f")
	if err != nil {
		return fmt.Errorf("failed to delete container: %s: %w", result.Stderr, err)
	}

	logger.Info("Container deleted", zap.String("name", name))
	return nil
}

// StartContainer starts an LXC container
func (lm *LXCManager) StartContainer(name string) error {
	if !lm.enabled {
		return fmt.Errorf("LXC is not enabled")
	}

	result, err := lm.shell.Execute("lxc-start", "-n", name)
	if err != nil {
		return fmt.Errorf("failed to start container: %s: %w", result.Stderr, err)
	}

	logger.Info("Container started", zap.String("name", name))
	return nil
}

// StopContainer stops an LXC container
func (lm *LXCManager) StopContainer(name string, force bool) error {
	if !lm.enabled {
		return fmt.Errorf("LXC is not enabled")
	}

	var result *executor.CommandResult
	var err error

	if force {
		result, err = lm.shell.Execute("lxc-stop", "-n", name, "-k")
	} else {
		result, err = lm.shell.Execute("lxc-stop", "-n", name)
	}

	if err != nil {
		return fmt.Errorf("failed to stop container: %s: %w", result.Stderr, err)
	}

	logger.Info("Container stopped", zap.String("name", name), zap.Bool("force", force))
	return nil
}

// RestartContainer restarts an LXC container
func (lm *LXCManager) RestartContainer(name string) error {
	if !lm.enabled {
		return fmt.Errorf("LXC is not enabled")
	}

	if err := lm.StopContainer(name, false); err != nil {
		return err
	}

	if err := lm.StartContainer(name); err != nil {
		return err
	}

	logger.Info("Container restarted", zap.String("name", name))
	return nil
}

// FreezeContainer freezes (pauses) an LXC container
func (lm *LXCManager) FreezeContainer(name string) error {
	if !lm.enabled {
		return fmt.Errorf("LXC is not enabled")
	}

	result, err := lm.shell.Execute("lxc-freeze", "-n", name)
	if err != nil {
		return fmt.Errorf("failed to freeze container: %s: %w", result.Stderr, err)
	}

	logger.Info("Container frozen", zap.String("name", name))
	return nil
}

// UnfreezeContainer unfreezes (resumes) an LXC container
func (lm *LXCManager) UnfreezeContainer(name string) error {
	if !lm.enabled {
		return fmt.Errorf("LXC is not enabled")
	}

	result, err := lm.shell.Execute("lxc-unfreeze", "-n", name)
	if err != nil {
		return fmt.Errorf("failed to unfreeze container: %s: %w", result.Stderr, err)
	}

	logger.Info("Container unfrozen", zap.String("name", name))
	return nil
}

// GetAvailableTemplates returns available LXC templates
func (lm *LXCManager) GetAvailableTemplates() ([]Template, error) {
	if !lm.enabled {
		return nil, fmt.Errorf("LXC is not enabled")
	}

	// Common LXC templates
	templates := []Template{
		{
			Name:        "ubuntu",
			Description: "Ubuntu Linux",
			Releases:    []string{"20.04", "22.04", "24.04"},
			Architectures: []string{"amd64", "arm64"},
		},
		{
			Name:        "debian",
			Description: "Debian Linux",
			Releases:    []string{"bullseye", "bookworm"},
			Architectures: []string{"amd64", "arm64"},
		},
		{
			Name:        "alpine",
			Description: "Alpine Linux",
			Releases:    []string{"3.17", "3.18", "3.19"},
			Architectures: []string{"amd64", "arm64"},
		},
		{
			Name:        "centos",
			Description: "CentOS Linux",
			Releases:    []string{"7", "8", "9"},
			Architectures: []string{"amd64", "arm64"},
		},
	}

	return templates, nil
}

// AttachConsole attaches to a container console (returns command for execution)
func (lm *LXCManager) AttachConsole(name string) (string, error) {
	if !lm.enabled {
		return "", fmt.Errorf("LXC is not enabled")
	}

	// Return the command that can be executed via terminal
	return fmt.Sprintf("lxc-attach -n %s", name), nil
}

// ExecCommand executes a command in a container and returns the result
func (lm *LXCManager) ExecCommand(name string, command string) (*executor.CommandResult, error) {
	if !lm.enabled {
		return nil, fmt.Errorf("LXC is not enabled")
	}

	// Execute command in container
	result, err := lm.shell.Execute("lxc-attach", "-n", name, "--", "sh", "-c", command)
	if err != nil {
		return result, fmt.Errorf("failed to execute command in container: %w", err)
	}

	logger.Info("Command executed in container",
		zap.String("name", name),
		zap.String("command", command))
	return result, nil
}

// GetConsoleURL returns the console access URL/command for a container
func (lm *LXCManager) GetConsoleURL(name string) (string, error) {
	if !lm.enabled {
		return "", fmt.Errorf("LXC is not enabled")
	}

	// Check if container is running
	container, err := lm.GetContainer(name)
	if err != nil {
		return "", err
	}

	if container.State != "RUNNING" {
		return "", fmt.Errorf("container must be running to access console")
	}

	// Return terminal command for shell access
	return fmt.Sprintf("lxc-attach -n %s", name), nil
}

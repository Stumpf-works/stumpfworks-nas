// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ComposeStack represents a Docker Compose stack
type ComposeStack struct {
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Status    string            `json:"status"`
	Services  []ComposeService  `json:"services"`
	CreatedAt string            `json:"createdAt"`
	UpdatedAt string            `json:"updatedAt"`
}

// ComposeService represents a service in a Compose stack
type ComposeService struct {
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	Status     string   `json:"status"`
	Containers []string `json:"containers"`
}

// ComposeFile represents a docker-compose.yml structure
type ComposeFile struct {
	Version  string                            `yaml:"version"`
	Services map[string]ComposeServiceConfig   `yaml:"services"`
	Networks map[string]interface{}            `yaml:"networks,omitempty"`
	Volumes  map[string]interface{}            `yaml:"volumes,omitempty"`
}

// ComposeServiceConfig represents service configuration in compose file
type ComposeServiceConfig struct {
	Image       string              `yaml:"image,omitempty"`
	Build       interface{}         `yaml:"build,omitempty"`
	Ports       []string            `yaml:"ports,omitempty"`
	Environment interface{}         `yaml:"environment,omitempty"`
	Volumes     []string            `yaml:"volumes,omitempty"`
	DependsOn   interface{}         `yaml:"depends_on,omitempty"`
	Restart     string              `yaml:"restart,omitempty"`
	Networks    interface{}         `yaml:"networks,omitempty"`
}

// ListStacks lists all Docker Compose stacks
func (s *Service) ListStacks(ctx context.Context, stacksDir string) ([]ComposeStack, error) {
	if !s.available {
		return nil, fmt.Errorf("Docker is not available")
	}

	var stacks []ComposeStack

	// Check if stacks directory exists
	if _, err := os.Stat(stacksDir); os.IsNotExist(err) {
		return stacks, nil
	}

	// Read all subdirectories
	entries, err := os.ReadDir(stacksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read stacks directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		stackName := entry.Name()
		stackPath := filepath.Join(stacksDir, stackName)
		composePath := filepath.Join(stackPath, "docker-compose.yml")

		// Check if docker-compose.yml exists
		if _, err := os.Stat(composePath); os.IsNotExist(err) {
			// Try alternative names
			composePath = filepath.Join(stackPath, "docker-compose.yaml")
			if _, err := os.Stat(composePath); os.IsNotExist(err) {
				continue
			}
		}

		// Get stack info
		stack, err := s.GetStack(ctx, stackName, stackPath)
		if err != nil {
			// If we can't get full info, create basic stack info
			stack = ComposeStack{
				Name:   stackName,
				Path:   stackPath,
				Status: "unknown",
			}
		}

		stacks = append(stacks, stack)
	}

	return stacks, nil
}

// GetStack gets detailed information about a stack
func (s *Service) GetStack(ctx context.Context, name string, stackPath string) (ComposeStack, error) {
	if !s.available {
		return ComposeStack{}, fmt.Errorf("Docker is not available")
	}

	composePath := filepath.Join(stackPath, "docker-compose.yml")
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		composePath = filepath.Join(stackPath, "docker-compose.yaml")
	}

	// Parse compose file
	composeFile, err := s.parseComposeFile(composePath)
	if err != nil {
		return ComposeStack{}, fmt.Errorf("failed to parse compose file: %w", err)
	}

	// Get running containers for this stack
	containers, err := s.ListContainers(ctx, true)
	if err != nil {
		return ComposeStack{}, fmt.Errorf("failed to list containers: %w", err)
	}

	// Build services list
	var services []ComposeService

	for serviceName := range composeFile.Services {
		service := ComposeService{
			Name:       serviceName,
			Status:     "stopped",
			Containers: []string{},
		}

		// Find containers for this service
		for _, container := range containers {
			// Check if container belongs to this stack and service
			if container.Labels != nil {
				if project, ok := container.Labels["com.docker.compose.project"]; ok && project == name {
					if svc, ok := container.Labels["com.docker.compose.service"]; ok && svc == serviceName {
						service.Containers = append(service.Containers, container.ID[:12])
						if container.State == "running" {
							service.Status = "running"
						}
					}
				}
			}
		}

		services = append(services, service)
	}

	// Determine overall stack status
	status := "stopped"
	runningCount := 0
	for _, svc := range services {
		if svc.Status == "running" {
			runningCount++
		}
	}
	if runningCount > 0 {
		if runningCount == len(services) {
			status = "running"
		} else {
			status = "partial"
		}
	}

	// Get file info for timestamps
	fileInfo, _ := os.Stat(composePath)
	var createdAt, updatedAt string
	if fileInfo != nil {
		updatedAt = fileInfo.ModTime().Format("2006-01-02T15:04:05Z")
		createdAt = updatedAt // We don't track creation separately
	}

	return ComposeStack{
		Name:      name,
		Path:      stackPath,
		Status:    status,
		Services:  services,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// DeployStack deploys a Docker Compose stack
func (s *Service) DeployStack(ctx context.Context, stackPath string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	// Run docker compose up -d
	cmd := exec.CommandContext(ctx, "docker", "compose", "up", "-d")
	cmd.Dir = stackPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to deploy stack: %w, output: %s", err, string(output))
	}

	return nil
}

// StopStack stops a Docker Compose stack
func (s *Service) StopStack(ctx context.Context, stackPath string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	// Run docker compose stop
	cmd := exec.CommandContext(ctx, "docker", "compose", "stop")
	cmd.Dir = stackPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop stack: %w, output: %s", err, string(output))
	}

	return nil
}

// RemoveStack removes a Docker Compose stack
func (s *Service) RemoveStack(ctx context.Context, stackPath string, removeVolumes bool) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	// Run docker compose down
	args := []string{"compose", "down"}
	if removeVolumes {
		args = append(args, "-v")
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Dir = stackPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove stack: %w, output: %s", err, string(output))
	}

	return nil
}

// GetStackLogs gets logs from a Docker Compose stack
func (s *Service) GetStackLogs(ctx context.Context, stackPath string, tail int) (string, error) {
	if !s.available {
		return "", fmt.Errorf("Docker is not available")
	}

	args := []string{"compose", "logs"}
	if tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", tail))
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Dir = stackPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get stack logs: %w", err)
	}

	return string(output), nil
}

// RestartStack restarts a Docker Compose stack
func (s *Service) RestartStack(ctx context.Context, stackPath string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	// Run docker compose restart
	cmd := exec.CommandContext(ctx, "docker", "compose", "restart")
	cmd.Dir = stackPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to restart stack: %w, output: %s", err, string(output))
	}

	return nil
}

// CreateStack creates a new stack from compose content
func (s *Service) CreateStack(ctx context.Context, stacksDir string, name string, composeContent string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	// Validate stack name (only alphanumeric, dash, underscore)
	if !isValidStackName(name) {
		return fmt.Errorf("invalid stack name: only alphanumeric characters, dash and underscore allowed")
	}

	stackPath := filepath.Join(stacksDir, name)

	// Check if stack already exists
	if _, err := os.Stat(stackPath); !os.IsNotExist(err) {
		return fmt.Errorf("stack %s already exists", name)
	}

	// Create stack directory
	if err := os.MkdirAll(stackPath, 0755); err != nil {
		return fmt.Errorf("failed to create stack directory: %w", err)
	}

	// Write compose file
	composePath := filepath.Join(stackPath, "docker-compose.yml")
	if err := os.WriteFile(composePath, []byte(composeContent), 0644); err != nil {
		// Clean up on error
		os.RemoveAll(stackPath)
		return fmt.Errorf("failed to write compose file: %w", err)
	}

	return nil
}

// UpdateStack updates an existing stack's compose file
func (s *Service) UpdateStack(ctx context.Context, stackPath string, composeContent string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	composePath := filepath.Join(stackPath, "docker-compose.yml")

	// Write new compose file
	if err := os.WriteFile(composePath, []byte(composeContent), 0644); err != nil {
		return fmt.Errorf("failed to update compose file: %w", err)
	}

	return nil
}

// DeleteStack deletes a stack directory
func (s *Service) DeleteStack(ctx context.Context, stackPath string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	// Remove stack directory
	if err := os.RemoveAll(stackPath); err != nil {
		return fmt.Errorf("failed to delete stack directory: %w", err)
	}

	return nil
}

// GetComposeFile reads and returns the compose file content
func (s *Service) GetComposeFile(stackPath string) (string, error) {
	composePath := filepath.Join(stackPath, "docker-compose.yml")
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		composePath = filepath.Join(stackPath, "docker-compose.yaml")
	}

	content, err := os.ReadFile(composePath)
	if err != nil {
		return "", fmt.Errorf("failed to read compose file: %w", err)
	}

	return string(content), nil
}

// parseComposeFile parses a docker-compose.yml file
func (s *Service) parseComposeFile(path string) (*ComposeFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var composeFile ComposeFile
	if err := yaml.Unmarshal(data, &composeFile); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &composeFile, nil
}

// isValidStackName checks if a stack name is valid
func isValidStackName(name string) bool {
	if len(name) == 0 || len(name) > 100 {
		return false
	}

	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return false
		}
	}

	return true
}

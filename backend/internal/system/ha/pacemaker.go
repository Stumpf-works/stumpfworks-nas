// Package ha provides High Availability features for Stumpf.Works NAS
package ha

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// PacemakerManager manages Pacemaker/Corosync cluster resources
type PacemakerManager struct {
	shell   executor.ShellExecutor
	enabled bool
}

// ClusterNode represents a node in the cluster
type ClusterNode struct {
	Name   string `json:"name"`
	ID     int    `json:"id"`
	IP     string `json:"ip"`
	Online bool   `json:"online"`
}

// ClusterStatus represents the current cluster status
type ClusterStatus struct {
	Name              string        `json:"name"`
	Nodes             []ClusterNode `json:"nodes"`
	Resources         []Resource    `json:"resources"`
	Quorum            bool          `json:"quorum"`
	MaintenanceMode   bool          `json:"maintenance_mode"`
	StonithEnabled    bool          `json:"stonith_enabled"`
	SymmetricCluster  bool          `json:"symmetric_cluster"`
}

// Resource represents a Pacemaker resource
type Resource struct {
	ID        string `json:"id"`
	Type      string `json:"type"`     // ocf, systemd, service
	Agent     string `json:"agent"`    // e.g., ocf:heartbeat:IPaddr2
	Node      string `json:"node"`     // Node where resource is running
	Active    bool   `json:"active"`
	Managed   bool   `json:"managed"`
	Failed    bool   `json:"failed"`
}

// ResourceConfig represents a resource configuration
type ResourceConfig struct {
	ID     string            `json:"id"`
	Type   string            `json:"type"`   // ocf, systemd, service
	Agent  string            `json:"agent"`  // e.g., ocf:heartbeat:IPaddr2
	Params map[string]string `json:"params"` // Resource parameters
	Op     []OpConfig        `json:"op"`     // Operations (monitor, start, stop)
}

// OpConfig represents an operation configuration
type OpConfig struct {
	Name     string `json:"name"`     // monitor, start, stop
	Interval string `json:"interval"` // e.g., "10s", "30s"
	Timeout  string `json:"timeout"`  // e.g., "20s"
}

// NewPacemakerManager creates a new Pacemaker/Corosync manager
func NewPacemakerManager(shell executor.ShellExecutor) (*PacemakerManager, error) {
	manager := &PacemakerManager{
		shell:   shell,
		enabled: false,
	}

	// Check if pcs (Pacemaker Configuration System) is available
	result, err := shell.Execute("which", "pcs")
	if err != nil || result.Stdout == "" {
		logger.Warn("pcs not found, Pacemaker features will be disabled")
		return manager, fmt.Errorf("pcs not available: install pcs package")
	}

	manager.enabled = true
	logger.Info("Pacemaker manager initialized successfully")
	return manager, nil
}

// IsEnabled returns whether Pacemaker is available
func (pm *PacemakerManager) IsEnabled() bool {
	return pm.enabled
}

// GetClusterStatus gets the current cluster status
func (pm *PacemakerManager) GetClusterStatus() (*ClusterStatus, error) {
	if !pm.enabled {
		return nil, fmt.Errorf("Pacemaker is not enabled")
	}

	status := &ClusterStatus{
		Nodes:     []ClusterNode{},
		Resources: []Resource{},
	}

	// Get cluster name
	result, err := pm.shell.Execute("sudo", "pcs", "cluster", "status")
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster status: %s: %w", result.Stderr, err)
	}

	// Parse cluster name
	lines := strings.Split(result.Stdout, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Cluster name:") {
			status.Name = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	// Get quorum status
	result, err = pm.shell.Execute("sudo", "pcs", "quorum", "status")
	if err == nil && strings.Contains(result.Stdout, "Quorate: Yes") {
		status.Quorum = true
	}

	// Get node status
	result, err = pm.shell.Execute("sudo", "pcs", "status", "nodes")
	if err == nil {
		status.Nodes = pm.parseNodes(result.Stdout)
	}

	// Get resource status
	result, err = pm.shell.Execute("sudo", "pcs", "status", "resources")
	if err == nil {
		status.Resources = pm.parseResources(result.Stdout)
	}

	// Check maintenance mode
	result, err = pm.shell.Execute("sudo", "pcs", "property", "show", "maintenance-mode")
	if err == nil && strings.Contains(result.Stdout, "true") {
		status.MaintenanceMode = true
	}

	// Check STONITH
	result, err = pm.shell.Execute("sudo", "pcs", "property", "show", "stonith-enabled")
	if err == nil && strings.Contains(result.Stdout, "true") {
		status.StonithEnabled = true
	}

	return status, nil
}

// parseNodes parses node information from pcs output
func (pm *PacemakerManager) parseNodes(output string) []ClusterNode {
	nodes := []ClusterNode{}
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Pacemaker") || strings.HasPrefix(line, "Online:") || strings.HasPrefix(line, "Offline:") {
			continue
		}

		// Parse node information (format: "  node1")
		nodeRegex := regexp.MustCompile(`^\s*(\S+)`)
		if matches := nodeRegex.FindStringSubmatch(line); len(matches) > 1 {
			nodeName := matches[1]
			online := strings.Contains(output, "Online:") && strings.Contains(output, nodeName)

			node := ClusterNode{
				Name:   nodeName,
				Online: online,
			}
			nodes = append(nodes, node)
		}
	}

	return nodes
}

// parseResources parses resource information from pcs output
func (pm *PacemakerManager) parseResources(output string) []Resource {
	resources := []Resource{}
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Full List") || strings.HasPrefix(line, "*") {
			continue
		}

		// Parse resource (format: "  resource_id (type:agent): Started node1")
		resourceRegex := regexp.MustCompile(`^\s*(\S+)\s+\(([^:]+):([^)]+)\):\s+(\w+)(?:\s+(\S+))?`)
		if matches := resourceRegex.FindStringSubmatch(line); len(matches) > 4 {
			resource := Resource{
				ID:     matches[1],
				Type:   matches[2],
				Agent:  matches[3],
				Active: matches[4] == "Started",
			}
			if len(matches) > 5 {
				resource.Node = matches[5]
			}
			resources = append(resources, resource)
		}
	}

	return resources
}

// CreateResource creates a new Pacemaker resource
func (pm *PacemakerManager) CreateResource(config ResourceConfig) error {
	if !pm.enabled {
		return fmt.Errorf("Pacemaker is not enabled")
	}

	// Build pcs command
	args := []string{"sudo", "pcs", "resource", "create", config.ID}

	// Add agent
	if config.Type != "" && config.Agent != "" {
		args = append(args, fmt.Sprintf("%s:%s", config.Type, config.Agent))
	} else if config.Agent != "" {
		args = append(args, config.Agent)
	}

	// Add parameters
	for key, value := range config.Params {
		args = append(args, fmt.Sprintf("%s=%s", key, value))
	}

	// Add operations
	for _, op := range config.Op {
		args = append(args, "op", op.Name)
		if op.Interval != "" {
			args = append(args, "interval="+op.Interval)
		}
		if op.Timeout != "" {
			args = append(args, "timeout="+op.Timeout)
		}
	}

	result, err := pm.shell.Execute(args[0], args[1:]...)
	if err != nil {
		return fmt.Errorf("failed to create resource: %s: %w", result.Stderr, err)
	}

	logger.Info("Pacemaker resource created", zap.String("id", config.ID))
	return nil
}

// DeleteResource deletes a Pacemaker resource
func (pm *PacemakerManager) DeleteResource(resourceID string) error {
	if !pm.enabled {
		return fmt.Errorf("Pacemaker is not enabled")
	}

	result, err := pm.shell.Execute("sudo", "pcs", "resource", "delete", resourceID)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %s: %w", result.Stderr, err)
	}

	logger.Info("Pacemaker resource deleted", zap.String("id", resourceID))
	return nil
}

// EnableResource enables (unmanage -> manage) a resource
func (pm *PacemakerManager) EnableResource(resourceID string) error {
	if !pm.enabled {
		return fmt.Errorf("Pacemaker is not enabled")
	}

	result, err := pm.shell.Execute("sudo", "pcs", "resource", "enable", resourceID)
	if err != nil {
		return fmt.Errorf("failed to enable resource: %s: %w", result.Stderr, err)
	}

	logger.Info("Pacemaker resource enabled", zap.String("id", resourceID))
	return nil
}

// DisableResource disables (manage -> unmanage) a resource
func (pm *PacemakerManager) DisableResource(resourceID string) error {
	if !pm.enabled {
		return fmt.Errorf("Pacemaker is not enabled")
	}

	result, err := pm.shell.Execute("sudo", "pcs", "resource", "disable", resourceID)
	if err != nil {
		return fmt.Errorf("failed to disable resource: %s: %w", result.Stderr, err)
	}

	logger.Info("Pacemaker resource disabled", zap.String("id", resourceID))
	return nil
}

// MoveResource moves a resource to a specific node
func (pm *PacemakerManager) MoveResource(resourceID, targetNode string) error {
	if !pm.enabled {
		return fmt.Errorf("Pacemaker is not enabled")
	}

	result, err := pm.shell.Execute("sudo", "pcs", "resource", "move", resourceID, targetNode)
	if err != nil {
		return fmt.Errorf("failed to move resource: %s: %w", result.Stderr, err)
	}

	logger.Info("Pacemaker resource moved", zap.String("id", resourceID), zap.String("node", targetNode))
	return nil
}

// ClearResource clears the failed state of a resource
func (pm *PacemakerManager) ClearResource(resourceID string) error {
	if !pm.enabled {
		return fmt.Errorf("Pacemaker is not enabled")
	}

	result, err := pm.shell.Execute("sudo", "pcs", "resource", "cleanup", resourceID)
	if err != nil {
		return fmt.Errorf("failed to clear resource: %s: %w", result.Stderr, err)
	}

	logger.Info("Pacemaker resource cleared", zap.String("id", resourceID))
	return nil
}

// SetMaintenanceMode enables or disables maintenance mode
func (pm *PacemakerManager) SetMaintenanceMode(enabled bool) error {
	if !pm.enabled {
		return fmt.Errorf("Pacemaker is not enabled")
	}

	value := "false"
	if enabled {
		value = "true"
	}

	result, err := pm.shell.Execute("sudo", "pcs", "property", "set", "maintenance-mode="+value)
	if err != nil {
		return fmt.Errorf("failed to set maintenance mode: %s: %w", result.Stderr, err)
	}

	logger.Info("Pacemaker maintenance mode changed", zap.Bool("enabled", enabled))
	return nil
}

// StandbyNode puts a node in standby mode
func (pm *PacemakerManager) StandbyNode(nodeName string) error {
	if !pm.enabled {
		return fmt.Errorf("Pacemaker is not enabled")
	}

	result, err := pm.shell.Execute("sudo", "pcs", "node", "standby", nodeName)
	if err != nil {
		return fmt.Errorf("failed to standby node: %s: %w", result.Stderr, err)
	}

	logger.Info("Node put in standby", zap.String("node", nodeName))
	return nil
}

// UnstandbyNode removes a node from standby mode
func (pm *PacemakerManager) UnstandbyNode(nodeName string) error {
	if !pm.enabled {
		return fmt.Errorf("Pacemaker is not enabled")
	}

	result, err := pm.shell.Execute("sudo", "pcs", "node", "unstandby", nodeName)
	if err != nil {
		return fmt.Errorf("failed to unstandby node: %s: %w", result.Stderr, err)
	}

	logger.Info("Node removed from standby", zap.String("node", nodeName))
	return nil
}

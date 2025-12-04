// Revision: 2025-12-04 | Author: Claude | Version: 1.0.0
package docker

import "time"

// ImageResponse represents a Docker image for API responses
type ImageResponse struct {
	ID           string            `json:"id"`
	RepoTags     []string          `json:"repoTags"`
	RepoDigests  []string          `json:"repoDigests"`
	Created      int64             `json:"created"`      // Unix timestamp
	Size         int64             `json:"size"`         // Bytes
	VirtualSize  int64             `json:"virtualSize"`  // Bytes
	SharedSize   int64             `json:"sharedSize"`   // Bytes
	Labels       map[string]string `json:"labels"`
	Containers   int64             `json:"containers"`   // Number of containers using this image
}

// ContainerResponse represents a Docker container for API responses
type ContainerResponse struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	State   string            `json:"state"`
	Status  string            `json:"status"`
	Created string            `json:"created"`
	Ports   []PortMapping     `json:"ports,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
	Networks map[string]interface{} `json:"networks,omitempty"`
}

// PortMapping represents a port mapping
type PortMapping struct {
	PrivatePort uint16 `json:"privatePort"`
	PublicPort  uint16 `json:"publicPort,omitempty"`
	Type        string `json:"type"`
}

// NetworkResponse represents a Docker network for API responses
type NetworkResponse struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Scope      string            `json:"scope"`
	Internal   bool              `json:"internal"`
	EnableIPv6 bool              `json:"enableIPv6"`
	IPAM       *IPAMConfig       `json:"ipam,omitempty"`
	Containers map[string]interface{} `json:"containers,omitempty"`
	Options    map[string]string `json:"options,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
	CreatedAt  time.Time         `json:"createdAt"`
}

// IPAMConfig represents IPAM configuration
type IPAMConfig struct {
	Driver string        `json:"driver"`
	Config []IPAMSubnet  `json:"config,omitempty"`
}

// IPAMSubnet represents an IPAM subnet configuration
type IPAMSubnet struct {
	Subnet  string `json:"subnet,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

// VolumeResponse represents a Docker volume for API responses
type VolumeResponse struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	CreatedAt  string            `json:"createdAt"`
	Labels     map[string]string `json:"labels,omitempty"`
	Scope      string            `json:"scope"`
	Options    map[string]string `json:"options,omitempty"`
}

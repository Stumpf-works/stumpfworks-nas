// Revision: 2025-12-02 | Author: Claude | Version: 1.2.0
package docker

import (
	"archive/tar"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// Service handles Docker operations
type Service struct {
	client    *client.Client
	available bool
}

var (
	globalService *Service
)

// Initialize creates a new Docker service and checks availability
func Initialize() (*Service, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return &Service{available: false}, fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Check if Docker is available
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = cli.Ping(ctx)
	if err != nil {
		cli.Close()
		return &Service{available: false}, fmt.Errorf("Docker is not available: %w", err)
	}

	globalService = &Service{
		client:    cli,
		available: true,
	}

	return globalService, nil
}

// GetService returns the global Docker service instance
func GetService() *Service {
	return globalService
}

// IsAvailable returns whether Docker is available
func (s *Service) IsAvailable() bool {
	if s == nil {
		return false
	}
	return s.available
}

// Close closes the Docker client connection
func (s *Service) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

// Container Operations

// ListContainers lists all containers
func (s *Service) ListContainers(ctx context.Context, all bool) ([]ContainerResponse, error) {
	if !s.available {
		return nil, fmt.Errorf("Docker is not available")
	}

	containers, err := s.client.ContainerList(ctx, container.ListOptions{All: all})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	// Map to our response format
	response := make([]ContainerResponse, len(containers))
	for i, c := range containers {
		// Map ports
		ports := make([]PortMapping, 0, len(c.Ports))
		for _, p := range c.Ports {
			ports = append(ports, PortMapping{
				PrivatePort: p.PrivatePort,
				PublicPort:  p.PublicPort,
				Type:        p.Type,
			})
		}

		// Get first name (remove leading slash)
		name := ""
		if len(c.Names) > 0 {
			name = c.Names[0]
			if len(name) > 0 && name[0] == '/' {
				name = name[1:]
			}
		}

		response[i] = ContainerResponse{
			ID:       c.ID,
			Name:     name,
			Image:    c.Image,
			State:    c.State,
			Status:   c.Status,
			Created:  fmt.Sprintf("%d", c.Created),
			Ports:    ports,
			Labels:   c.Labels,
			Networks: make(map[string]interface{}),
		}

		// Map networks
		for netName, net := range c.NetworkSettings.Networks {
			response[i].Networks[netName] = net
		}
	}

	return response, nil
}

// StartContainer starts a container
func (s *Service) StartContainer(ctx context.Context, containerID string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	if err := s.client.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

// StopContainer stops a container
func (s *Service) StopContainer(ctx context.Context, containerID string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	timeout := 10 // seconds
	if err := s.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	return nil
}

// RestartContainer restarts a container
func (s *Service) RestartContainer(ctx context.Context, containerID string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	timeout := 10 // seconds
	if err := s.client.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		return fmt.Errorf("failed to restart container: %w", err)
	}

	return nil
}

// RemoveContainer removes a container
func (s *Service) RemoveContainer(ctx context.Context, containerID string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	if err := s.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	return nil
}

// GetContainerLogs gets container logs
func (s *Service) GetContainerLogs(ctx context.Context, containerID string) (string, error) {
	if !s.available {
		return "", fmt.Errorf("Docker is not available")
	}

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "500", // Last 500 lines
	}

	reader, err := s.client.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read container logs: %w", err)
	}

	return string(logs), nil
}

// Image Operations

// ListImages lists all images
func (s *Service) ListImages(ctx context.Context) ([]ImageResponse, error) {
	if !s.available {
		return nil, fmt.Errorf("Docker is not available")
	}

	images, err := s.client.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	// Map to our response format
	response := make([]ImageResponse, len(images))
	for i, img := range images {
		response[i] = ImageResponse{
			ID:          img.ID,
			RepoTags:    img.RepoTags,
			RepoDigests: img.RepoDigests,
			Created:     img.Created,
			Size:        img.Size,
			VirtualSize: img.VirtualSize,
			SharedSize:  img.SharedSize,
			Labels:      img.Labels,
			Containers:  img.Containers,
		}
	}

	return response, nil
}

// PullImage pulls an image
func (s *Service) PullImage(ctx context.Context, imageName string) (string, error) {
	if !s.available {
		return "", fmt.Errorf("Docker is not available")
	}

	reader, err := s.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()

	// Read the pull output
	output, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read pull output: %w", err)
	}

	return string(output), nil
}

// RemoveImage removes an image
func (s *Service) RemoveImage(ctx context.Context, imageID string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	_, err := s.client.ImageRemove(ctx, imageID, image.RemoveOptions{Force: true})
	if err != nil {
		return fmt.Errorf("failed to remove image: %w", err)
	}

	return nil
}

// Volume Operations

// ListVolumes lists all volumes
func (s *Service) ListVolumes(ctx context.Context) ([]VolumeResponse, error) {
	if !s.available {
		return nil, fmt.Errorf("Docker is not available")
	}

	volumeList, err := s.client.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	// Map to our response format
	response := make([]VolumeResponse, len(volumeList.Volumes))
	for i, vol := range volumeList.Volumes {
		response[i] = VolumeResponse{
			Name:       vol.Name,
			Driver:     vol.Driver,
			Mountpoint: vol.Mountpoint,
			CreatedAt:  vol.CreatedAt,
			Labels:     vol.Labels,
			Scope:      vol.Scope,
			Options:    vol.Options,
		}
	}

	return response, nil
}

// RemoveVolume removes a volume
func (s *Service) RemoveVolume(ctx context.Context, volumeName string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	if err := s.client.VolumeRemove(ctx, volumeName, true); err != nil {
		return fmt.Errorf("failed to remove volume: %w", err)
	}

	return nil
}

// Network Operations

// ListNetworks lists all networks
func (s *Service) ListNetworks(ctx context.Context) ([]NetworkResponse, error) {
	if !s.available {
		return nil, fmt.Errorf("Docker is not available")
	}

	networks, err := s.client.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	// Map to our response format
	response := make([]NetworkResponse, len(networks))
	for i, net := range networks {
		resp := NetworkResponse{
			ID:         net.ID,
			Name:       net.Name,
			Driver:     net.Driver,
			Scope:      net.Scope,
			Internal:   net.Internal,
			EnableIPv6: net.EnableIPv6,
			Options:    net.Options,
			Labels:     net.Labels,
			CreatedAt:  net.Created,
			Containers: make(map[string]interface{}),
		}

		// Map IPAM configuration
		if net.IPAM.Driver != "" {
			resp.IPAM = &IPAMConfig{
				Driver: net.IPAM.Driver,
				Config: make([]IPAMSubnet, 0, len(net.IPAM.Config)),
			}
			for _, cfg := range net.IPAM.Config {
				resp.IPAM.Config = append(resp.IPAM.Config, IPAMSubnet{
					Subnet:  cfg.Subnet,
					Gateway: cfg.Gateway,
				})
			}
		}

		// Map containers
		for containerID, containerInfo := range net.Containers {
			resp.Containers[containerID] = containerInfo
		}

		response[i] = resp
	}

	return response, nil
}

// Advanced Container Operations (Portainer-like features)

// InspectContainer gets detailed container information
func (s *Service) InspectContainer(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	if !s.available {
		return types.ContainerJSON{}, fmt.Errorf("Docker is not available")
	}

	containerInfo, err := s.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return types.ContainerJSON{}, fmt.Errorf("failed to inspect container: %w", err)
	}

	return containerInfo, nil
}

// GetContainerStats gets container resource usage statistics
func (s *Service) GetContainerStats(ctx context.Context, containerID string) (container.StatsResponse, error) {
	if !s.available {
		return container.StatsResponse{}, fmt.Errorf("Docker is not available")
	}

	stats, err := s.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return container.StatsResponse{}, fmt.Errorf("failed to get container stats: %w", err)
	}
	defer stats.Body.Close()

	var statsJSON container.StatsResponse
	data, err := io.ReadAll(stats.Body)
	if err != nil {
		return container.StatsResponse{}, fmt.Errorf("failed to read stats: %w", err)
	}

	if err := json.Unmarshal(data, &statsJSON); err != nil {
		return container.StatsResponse{}, fmt.Errorf("failed to parse stats: %w", err)
	}

	return statsJSON, nil
}

// CreateContainer creates a new container
func (s *Service) CreateContainer(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.CreateResponse, error) {
	if !s.available {
		return container.CreateResponse{}, fmt.Errorf("Docker is not available")
	}

	resp, err := s.client.ContainerCreate(ctx, config, hostConfig, networkingConfig, nil, containerName)
	if err != nil {
		return container.CreateResponse{}, fmt.Errorf("failed to create container: %w", err)
	}

	return resp, nil
}

// PauseContainer pauses a container
func (s *Service) PauseContainer(ctx context.Context, containerID string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	if err := s.client.ContainerPause(ctx, containerID); err != nil {
		return fmt.Errorf("failed to pause container: %w", err)
	}

	return nil
}

// UnpauseContainer unpauses a container
func (s *Service) UnpauseContainer(ctx context.Context, containerID string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	if err := s.client.ContainerUnpause(ctx, containerID); err != nil {
		return fmt.Errorf("failed to unpause container: %w", err)
	}

	return nil
}

// Advanced Image Operations

// InspectImage gets detailed image information
func (s *Service) InspectImage(ctx context.Context, imageID string) (types.ImageInspect, error) {
	if !s.available {
		return types.ImageInspect{}, fmt.Errorf("Docker is not available")
	}

	imageInfo, _, err := s.client.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		return types.ImageInspect{}, fmt.Errorf("failed to inspect image: %w", err)
	}

	return imageInfo, nil
}

// SearchImages searches for images on Docker Hub
func (s *Service) SearchImages(ctx context.Context, term string) ([]registry.SearchResult, error) {
	if !s.available {
		return nil, fmt.Errorf("Docker is not available")
	}

	results, err := s.client.ImageSearch(ctx, term, registry.SearchOptions{Limit: 25})
	if err != nil {
		return nil, fmt.Errorf("failed to search images: %w", err)
	}

	return results, nil
}

// Advanced Volume Operations

// InspectVolume gets detailed volume information
func (s *Service) InspectVolume(ctx context.Context, volumeName string) (volume.Volume, error) {
	if !s.available {
		return volume.Volume{}, fmt.Errorf("Docker is not available")
	}

	vol, err := s.client.VolumeInspect(ctx, volumeName)
	if err != nil {
		return volume.Volume{}, fmt.Errorf("failed to inspect volume: %w", err)
	}

	return vol, nil
}

// CreateVolume creates a new volume
func (s *Service) CreateVolume(ctx context.Context, name string, driver string, labels map[string]string) (volume.Volume, error) {
	if !s.available {
		return volume.Volume{}, fmt.Errorf("Docker is not available")
	}

	vol, err := s.client.VolumeCreate(ctx, volume.CreateOptions{
		Name:   name,
		Driver: driver,
		Labels: labels,
	})
	if err != nil {
		return volume.Volume{}, fmt.Errorf("failed to create volume: %w", err)
	}

	return vol, nil
}

// Advanced Network Operations

// InspectNetwork gets detailed network information
func (s *Service) InspectNetwork(ctx context.Context, networkID string) (network.Inspect, error) {
	if !s.available {
		return network.Inspect{}, fmt.Errorf("Docker is not available")
	}

	networkInfo, err := s.client.NetworkInspect(ctx, networkID, network.InspectOptions{})
	if err != nil {
		return network.Inspect{}, fmt.Errorf("failed to inspect network: %w", err)
	}

	return networkInfo, nil
}

// CreateNetwork creates a new network
func (s *Service) CreateNetwork(ctx context.Context, name string, driver string) (network.CreateResponse, error) {
	if !s.available {
		return network.CreateResponse{}, fmt.Errorf("Docker is not available")
	}

	resp, err := s.client.NetworkCreate(ctx, name, network.CreateOptions{
		Driver: driver,
	})
	if err != nil {
		return network.CreateResponse{}, fmt.Errorf("failed to create network: %w", err)
	}

	return resp, nil
}

// RemoveNetwork removes a network
func (s *Service) RemoveNetwork(ctx context.Context, networkID string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	if err := s.client.NetworkRemove(ctx, networkID); err != nil {
		return fmt.Errorf("failed to remove network: %w", err)
	}

	return nil
}

// ConnectContainerToNetwork connects a container to a network
func (s *Service) ConnectContainerToNetwork(ctx context.Context, networkID string, containerID string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	if err := s.client.NetworkConnect(ctx, networkID, containerID, nil); err != nil {
		return fmt.Errorf("failed to connect container to network: %w", err)
	}

	return nil
}

// DisconnectContainerFromNetwork disconnects a container from a network
func (s *Service) DisconnectContainerFromNetwork(ctx context.Context, networkID string, containerID string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	if err := s.client.NetworkDisconnect(ctx, networkID, containerID, true); err != nil {
		return fmt.Errorf("failed to disconnect container from network: %w", err)
	}

	return nil
}

// System Operations

// GetDockerInfo gets Docker system information
func (s *Service) GetDockerInfo(ctx context.Context) (system.Info, error) {
	if !s.available {
		return system.Info{}, fmt.Errorf("Docker is not available")
	}

	info, err := s.client.Info(ctx)
	if err != nil {
		return system.Info{}, fmt.Errorf("failed to get Docker info: %w", err)
	}

	return info, nil
}

// GetDockerVersion gets Docker version information
func (s *Service) GetDockerVersion(ctx context.Context) (types.Version, error) {
	if !s.available {
		return types.Version{}, fmt.Errorf("Docker is not available")
	}

	version, err := s.client.ServerVersion(ctx)
	if err != nil {
		return types.Version{}, fmt.Errorf("failed to get Docker version: %w", err)
	}

	return version, nil
}

// PruneSystem prunes unused Docker objects (containers, networks, images, volumes)
func (s *Service) PruneSystem(ctx context.Context) (types.DiskUsage, error) {
	if !s.available {
		return types.DiskUsage{}, fmt.Errorf("Docker is not available")
	}

	usage, err := s.client.DiskUsage(ctx, types.DiskUsageOptions{})
	if err != nil {
		return types.DiskUsage{}, fmt.Errorf("failed to get disk usage: %w", err)
	}

	return usage, nil
}

// UpdateContainerResources updates resource limits for a container
func (s *Service) UpdateContainerResources(ctx context.Context, containerID string, resources container.Resources) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	updateConfig := container.UpdateConfig{
		Resources: resources,
	}

	_, err := s.client.ContainerUpdate(ctx, containerID, updateConfig)
	if err != nil {
		return fmt.Errorf("failed to update container resources: %w", err)
	}

	return nil
}

// ExecContainer executes a command in a container
func (s *Service) ExecContainer(ctx context.Context, containerID string, cmd []string) (string, error) {
	if !s.available {
		return "", fmt.Errorf("Docker is not available")
	}

	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	execID, err := s.client.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec: %w", err)
	}

	resp, err := s.client.ContainerExecAttach(ctx, execID.ID, container.ExecAttachOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec: %w", err)
	}
	defer resp.Close()

	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to read exec output: %w", err)
	}

	return string(output), nil
}

// GetContainerTop gets running processes in a container
func (s *Service) GetContainerTop(ctx context.Context, containerID string) (container.ContainerTopOKBody, error) {
	if !s.available {
		return container.ContainerTopOKBody{}, fmt.Errorf("Docker is not available")
	}

	top, err := s.client.ContainerTop(ctx, containerID, []string{})
	if err != nil {
		return container.ContainerTopOKBody{}, fmt.Errorf("failed to get container processes: %w", err)
	}

	return top, nil
}

// BuildImage builds a Docker image from a Dockerfile
func (s *Service) BuildImage(ctx context.Context, dockerfile []byte, tags []string, buildArgs map[string]*string, labels map[string]string) (string, error) {
	if !s.available {
		return "", fmt.Errorf("Docker is not available")
	}

	// Create a tar archive with the Dockerfile
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		defer pipeWriter.Close()
		tw := tar.NewWriter(pipeWriter)
		defer tw.Close()

		// Add Dockerfile to tar
		header := &tar.Header{
			Name: "Dockerfile",
			Mode: 0644,
			Size: int64(len(dockerfile)),
		}
		if err := tw.WriteHeader(header); err != nil {
			return
		}
		if _, err := tw.Write(dockerfile); err != nil {
			return
		}
	}()

	// Build options
	buildOptions := types.ImageBuildOptions{
		Tags:       tags,
		Dockerfile: "Dockerfile",
		Remove:     true, // Remove intermediate containers
		BuildArgs:  buildArgs,
		Labels:     labels,
	}

	// Build the image
	response, err := s.client.ImageBuild(ctx, pipeReader, buildOptions)
	if err != nil {
		return "", fmt.Errorf("failed to build image: %w", err)
	}
	defer response.Body.Close()

	// Read the build output
	output, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read build output: %w", err)
	}

	return string(output), nil
}

// TagImage tags an image
func (s *Service) TagImage(ctx context.Context, imageID string, repo string, tag string) error {
	if !s.available {
		return fmt.Errorf("Docker is not available")
	}

	err := s.client.ImageTag(ctx, imageID, repo+":"+tag)
	if err != nil {
		return fmt.Errorf("failed to tag image: %w", err)
	}

	return nil
}

// PushImage pushes an image to a registry
func (s *Service) PushImage(ctx context.Context, imageName string, registryAuth string) (string, error) {
	if !s.available {
		return "", fmt.Errorf("Docker is not available")
	}

	options := image.PushOptions{
		RegistryAuth: registryAuth, // Base64 encoded auth config JSON
	}

	reader, err := s.client.ImagePush(ctx, imageName, options)
	if err != nil {
		return "", fmt.Errorf("failed to push image: %w", err)
	}
	defer reader.Close()

	// Read the push output
	output, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read push output: %w", err)
	}

	return string(output), nil
}

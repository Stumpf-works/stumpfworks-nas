package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// SimpleClient provides simplified Docker operations
type SimpleClient struct {
	cli *client.Client
}

// NewSimpleClient creates a new simple Docker client
func NewSimpleClient() (*SimpleClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	return &SimpleClient{cli: cli}, nil
}

// Close closes the Docker client
func (c *SimpleClient) Close() error {
	return c.cli.Close()
}

// Ping checks if Docker daemon is running
func (c *SimpleClient) Ping(ctx context.Context) error {
	_, err := c.cli.Ping(ctx)
	return err
}

// ListContainers lists all containers
func (c *SimpleClient) ListContainers(ctx context.Context, all bool) ([]types.Container, error) {
	return c.cli.ContainerList(ctx, container.ListOptions{All: all})
}

// StartContainer starts a container
func (c *SimpleClient) StartContainer(ctx context.Context, id string) error {
	return c.cli.ContainerStart(ctx, id, container.StartOptions{})
}

// StopContainer stops a container
func (c *SimpleClient) StopContainer(ctx context.Context, id string) error {
	return c.cli.ContainerStop(ctx, id, container.StopOptions{})
}

// RestartContainer restarts a container
func (c *SimpleClient) RestartContainer(ctx context.Context, id string) error {
	return c.cli.ContainerRestart(ctx, id, container.StopOptions{})
}

// RemoveContainer removes a container
func (c *SimpleClient) RemoveContainer(ctx context.Context, id string, force bool) error {
	return c.cli.ContainerRemove(ctx, id, container.RemoveOptions{Force: force})
}

// GetContainerLogs gets container logs
func (c *SimpleClient) GetContainerLogs(ctx context.Context, id string) (io.ReadCloser, error) {
	return c.cli.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "100",
	})
}

// ListImages lists all images
func (c *SimpleClient) ListImages(ctx context.Context) ([]image.Summary, error) {
	return c.cli.ImageList(ctx, image.ListOptions{})
}

// RemoveImage removes an image
func (c *SimpleClient) RemoveImage(ctx context.Context, id string, force bool) error {
	_, err := c.cli.ImageRemove(ctx, id, image.RemoveOptions{Force: force})
	return err
}

// PullImage pulls an image
func (c *SimpleClient) PullImage(ctx context.Context, ref string) (io.ReadCloser, error) {
	return c.cli.ImagePull(ctx, ref, image.PullOptions{})
}

// ListVolumes lists all volumes
func (c *SimpleClient) ListVolumes(ctx context.Context) (volume.ListResponse, error) {
	return c.cli.VolumeList(ctx, volume.ListOptions{})
}

// RemoveVolume removes a volume
func (c *SimpleClient) RemoveVolume(ctx context.Context, id string, force bool) error {
	return c.cli.VolumeRemove(ctx, id, force)
}

// ListNetworks lists all networks
func (c *SimpleClient) ListNetworks(ctx context.Context) ([]network.Summary, error) {
	return c.cli.NetworkList(ctx, network.ListOptions{})
}

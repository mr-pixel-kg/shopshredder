package services

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"io"
	"os"
	"time"
)

type DockerServiceInterface interface {
	CreateContainer(ctx context.Context, containerName string) (string, error)
	StopContainer(ctx context.Context, containerID string) error
	ListContainers(ctx context.Context) ([]ContainerInfo, error)
}

type DockerService struct {
	Client *client.Client
}

func NewDockerService() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerService{Client: cli}, nil
}

func (ds *DockerService) ListContainers(ctx context.Context) ([]ContainerInfo, error) {
	containers, err := ds.Client.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var containerInfos []ContainerInfo
	for _, container := range containers {
		created := time.Unix(container.Created, 0).Format(time.RFC3339)
		containerInfo := ContainerInfo{
			ID:        container.ID,
			Name:      container.Names[0],
			Image:     container.Image,
			CreatedAt: created,
			State:     container.State,
			Status:    container.Status,
		}
		containerInfos = append(containerInfos, containerInfo)
	}

	return containerInfos, nil
}

func (ds *DockerService) CreateContainer(ctx context.Context, imageName, instanceName, host string) (string, error) {
	// Pull image
	out, err := ds.Client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return "", err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	// Create container
	labels := map[string]string{
		"sandbox_container": "true",
		"traefik.enable":    "true",
		fmt.Sprintf("traefik.http.routers.http-%s.rule", instanceName): fmt.Sprintf("Host(`%s`)", host),
	}
	resp, err := ds.Client.ContainerCreate(ctx, &container.Config{
		Image:  imageName,
		Labels: labels,
	}, nil, nil, nil, instanceName)
	if err != nil {
		return "", err
	}

	// Start container
	if err := ds.Client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (ds *DockerService) StopContainer(ctx context.Context, containerID string) error {
	noWaitTimeout := 0 // to not wait for the container to exit gracefully
	return ds.Client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &noWaitTimeout})
}

type ContainerInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
	State     string `json:"state"`
	Status    string `json:"status"`
}

package services

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

type DockerImage struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Tag     string    `json:"tag"`
	Created time.Time `json:"created_at"`
	Size    int64     `json:"size"`
}

type DockerContainer struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	Created time.Time         `json:"created_at"`
	Status  string            `json:"status"`
	Labels  map[string]string `json:"labels"`
}

type DockerServiceInterface interface {
	ListImages(ctx context.Context) ([]DockerImage, error)
	GetImage(ctx context.Context, imageId string) (DockerImage, error)
	PullImage(ctx context.Context, imageName string) (DockerImage, error)
	RemoveImage(ctx context.Context, imageId string) error
	ListContainers(ctx context.Context) ([]DockerContainer, error)
	GetContainer(ctx context.Context, containerId string) (DockerContainer, error)
	CreateContainer(ctx context.Context, imageName, containerName string, labels map[string]string) (string, error)
	RemoveContainer(ctx context.Context, containerId string) error
}

type DockerService struct {
	dockerClient *client.Client
}

func NewDockerService() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	dockerService := &DockerService{dockerClient: cli}
	return dockerService, nil
}

func (ds *DockerService) ListImages(ctx context.Context) ([]DockerImage, error) {
	images, err := ds.dockerClient.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return nil, err
	}

	var outputImages = make([]DockerImage, 0)
	for _, img := range images {

		// Extract docker image id from format "sha256:hash"
		/*id := strings.Split(img.ID, ":")[1]

		// Validate
		if len(img.RepoTags) != 1 {
			slog.Debug("Skip Docker Image since it has not exactly one RepoTag assigned", "imageId", img.ID)
			continue
		}
		repoTag := img.RepoTags[0] // Format: "name:tag"
		if !strings.Contains(repoTag, ":") {
			slog.Debug("Skip Docker Image since RepoTag is in wrong format", "imageId", img.ID, "repoTag", repoTag)
			continue
		}

		// Extract image name and tag
		imageName := strings.Split(repoTag, ":")[0]
		imageTag := strings.Split(repoTag, ":")[1]

		outputImages = append(outputImages, DockerImage{
			ID:      id,
			Name:    imageName,
			Tag:     imageTag,
			Created: time.Unix(img.Created, 0),
			Size:    img.Size,
		})*/

		parsedImage, err := parseDockerImage(img.ID, img.RepoTags, time.Unix(img.Created, 0), img.Size)
		if err != nil {
			slog.Debug("Skipping Docker Image", "imageId", img.ID, "reason", err)
			continue
		}
		outputImages = append(outputImages, parsedImage)

	}

	return outputImages, nil
}

func (ds *DockerService) GetImage(ctx context.Context, imageId string) (DockerImage, error) {
	img, _, err := ds.dockerClient.ImageInspectWithRaw(ctx, imageId)
	if err != nil {
		return DockerImage{}, err
	}

	// Extract docker image id from format "sha256:hash"
	/*id := strings.Split(image.ID, ":")[1]

	// Validate
	if len(image.RepoTags) != 1 {
		slog.Debug("Skip Docker Image since it has not exactly one RepoTag assigned", "imageId", image.ID)
		return DockerImage{}, nil
	}
	repoTag := image.RepoTags[0] // Format: "name:tag"
	if !strings.Contains(repoTag, ":") {
		slog.Debug("Skip Docker Image since RepoTag is in wrong format", "imageId", image.ID, "repoTag", repoTag)
		return DockerImage{}, nil
	}

	// Extract image name and tag
	imageName := strings.Split(repoTag, ":")[0]
	imageTag := strings.Split(repoTag, ":")[1]

	// Parse image creation time
	createdTime, err := time.Parse(time.RFC3339Nano, image.Created)
	if err != nil {
		slog.Error("Error parsing the creation time of docker image", "image", repoTag, "error", err)
		return DockerImage{}, nil
	}

	return DockerImage{
		ID:      id,
		Name:    imageName,
		Tag:     imageTag,
		Created: createdTime,
		Size:    image.Size,
	}, nil*/

	// Parse image creation time
	createdTime, err := time.Parse(time.RFC3339Nano, img.Created)
	if err != nil {
		slog.Error("Error parsing the creation time of docker image", "image", img.ID, "error", err)
		return DockerImage{}, err
	}

	dockerImage, err := parseDockerImage(img.ID, img.RepoTags, createdTime, img.Size)
	if err != nil {
		slog.Debug("Skipping Docker Image", "imageId", img.ID, "reason", err)
		return DockerImage{}, nil
	}

	return dockerImage, nil
}

func parseDockerImage(imageID string, repoTags []string, created time.Time, size int64) (DockerImage, error) {
	// Extract Docker Image Id (Format: "sha256:hash")
	parts := strings.SplitN(imageID, ":", 2)
	if len(parts) != 2 {
		return DockerImage{}, fmt.Errorf("Invalid image ID format")
	}
	id := parts[1]

	// Validate RepoTags
	if len(repoTags) != 1 {
		return DockerImage{}, fmt.Errorf("Image has not exactly one RepoTag assigned")
	}
	repoTag := repoTags[0]

	// Extract name and tag
	name, tag, found := strings.Cut(repoTag, ":")
	if !found {
		return DockerImage{}, fmt.Errorf("RepoTag is in wrong format")
	}

	return DockerImage{
		ID:      id,
		Name:    name,
		Tag:     tag,
		Created: created,
		Size:    size,
	}, nil
}

func (ds *DockerService) PullImage(ctx context.Context, imageName string) (DockerImage, error) {
	out, err := ds.dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return DockerImage{}, err
	}

	defer out.Close()

	io.Copy(os.Stdout, out)

	return ds.GetImage(ctx, imageName)
}

func (ds *DockerService) RemoveImage(ctx context.Context, imageId string) error {
	_, err := ds.dockerClient.ImageRemove(ctx, imageId, image.RemoveOptions{Force: false})
	if err != nil {
		return err
	}
	return err
}

func (ds *DockerService) ListContainers(ctx context.Context) ([]DockerContainer, error) {
	containers, err := ds.dockerClient.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var outputContainers = make([]DockerContainer, 0)
	for _, cont := range containers {

		c := DockerContainer{
			ID:      cont.ID,
			Name:    cont.Names[0],
			Image:   cont.Image,
			Created: time.Unix(cont.Created, 0),
			Status:  cont.Status,
			Labels:  cont.Labels,
		}
		outputContainers = append(outputContainers, c)

	}

	return outputContainers, nil
}

func (ds *DockerService) GetContainer(ctx context.Context, containerId string) (DockerContainer, error) {
	cont, err := ds.dockerClient.ContainerInspect(ctx, containerId)
	if err != nil {
		return DockerContainer{}, err
	}

	// Parse image creation time (TODO test if this works)
	createdTime, err := time.Parse(time.RFC3339Nano, cont.Created)
	if err != nil {
		slog.Error("Error parsing the creation time of docker container", "container", cont.Name, "error", err)
		return DockerContainer{}, err
	}

	return DockerContainer{
		ID:      cont.ID,
		Name:    cont.Name,
		Image:   cont.Image,
		Created: createdTime,
		Status:  cont.State.Status,
		Labels:  cont.Config.Labels,
	}, nil
}

func (ds *DockerService) CreateContainer(ctx context.Context, imageName, containerName string, labels map[string]string) (string, error) {

	// Prepare container configuration
	containerConfig := &container.Config{
		Image:  imageName,
		Labels: labels,
	}

	// Create container
	cont, err := ds.dockerClient.ContainerCreate(ctx, containerConfig, nil, nil, nil, containerName)
	if err != nil {
		return "", err
	}

	// Start container
	err = ds.dockerClient.ContainerStart(ctx, cont.ID, container.StartOptions{})
	if err != nil {
		return "", err
	}

	return cont.ID, nil
}

func (ds *DockerService) RemoveContainer(ctx context.Context, containerId string) error {
	// TODO consider using a proper timeout
	noWaitTimeout := 0 // to not wait for the container to exit gracefully
	err := ds.dockerClient.ContainerStop(ctx, containerId, container.StopOptions{Timeout: &noWaitTimeout})
	if err != nil {
		return fmt.Errorf("Failed to stop container %s: %w", containerId, err)
	}

	err = ds.dockerClient.ContainerRemove(ctx, containerId, container.RemoveOptions{Force: true, RemoveVolumes: true})
	if err != nil {
		return fmt.Errorf("Failed to delete container %s: %w", containerId, err)
	}

	return nil
}

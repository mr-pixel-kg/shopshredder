package docker

import "context"

type SandboxCreateRequest struct {
	ImageName     string
	ContainerName string
	Hostname      string
}

type SandboxContainer struct {
	ID   string
	Name string
	URL  string
}

type Client interface {
	CreateContainer(ctx context.Context, request SandboxCreateRequest) (*SandboxContainer, error)
	DeleteContainer(ctx context.Context, containerID string) error
	CommitContainer(ctx context.Context, containerID, targetImage string) error
}

type NoopClient struct{}

func NewNoopClient() *NoopClient {
	return &NoopClient{}
}

func (c *NoopClient) CreateContainer(_ context.Context, request SandboxCreateRequest) (*SandboxContainer, error) {
	return &SandboxContainer{
		ID:   "stub-" + request.ContainerName,
		Name: request.ContainerName,
		URL:  "https://" + request.Hostname,
	}, nil
}

func (c *NoopClient) DeleteContainer(_ context.Context, _ string) error {
	return nil
}

func (c *NoopClient) CommitContainer(_ context.Context, _ string, _ string) error {
	return nil
}

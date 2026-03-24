package registry

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	parsed, err := time.ParseDuration(value.Value)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", value.Value, err)
	}
	d.Duration = parsed
	return nil
}

type ImageRegistry struct {
	Images []ImageEntry `yaml:"images"`
}

type ImageEntry struct {
	Match        string             `yaml:"match"`
	InternalPort *int               `yaml:"internal_port,omitempty"`
	Env          map[string]string  `yaml:"env,omitempty"`
	PostStart    []ExecCommand      `yaml:"post_start,omitempty"`
	PreStop      []ExecCommand      `yaml:"pre_stop,omitempty"`
	HealthCheck  *HealthCheckConfig `yaml:"health_check,omitempty"`
	Labels       map[string]string  `yaml:"labels,omitempty"`
	Volumes      []string           `yaml:"volumes,omitempty"`
}

type ExecCommand struct {
	Command    []string `yaml:"command"`
	Delay      Duration `yaml:"delay,omitempty"`
	Timeout    Duration `yaml:"timeout,omitempty"`
	Retries    int      `yaml:"retries,omitempty"`
	RetryDelay Duration `yaml:"retry_delay,omitempty"`
}

type HealthCheckConfig struct {
	Path     string   `yaml:"path"`
	Timeout  Duration `yaml:"timeout"`
	Interval Duration `yaml:"interval"`
}

type TemplateContext struct {
	Hostname       string
	URL            string
	Scheme         string
	Port           string
	ContainerName  string
	TrustedProxies string
	DockerMode     string
	Network        string
	InternalPort   string
	ImageName      string
	ImageRepo      string
	ImageTag       string
	SandboxID      string
	HostSuffix     string
	TTL            string
	ExpiresAt      string
	ClientIP       string
}

type ResolvedImage struct {
	InternalPort int
	Env          []string
	PostStart    []ExecCommand
	PreStop      []ExecCommand
	HealthCheck  *HealthCheckConfig
	Labels       map[string]string
	Volumes      []string
}

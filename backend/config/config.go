package config

import (
	"github.com/spf13/viper"
	"log"
)

const CONFIG_FILE = "./config.yml"

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Sandbox  SandboxConfig  `mapstructure:"sandbox"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Guard    GuardConfig    `mapstructure:"guard"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Port           int      `mapstructure:"port"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type SandboxConfig struct {
	UrlPrefix       string `mapstructure:"url_prefix"`
	UrlSuffix       string `mapstructure:"url_suffix"`
	DefaultLifetime int    `mapstructure:"default_lifetime"`
}

type AuthConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type GuardConfig struct {
	MaxTotalSandboxes      int `mapstructure:"max_total_sandboxes"`
	MaxSandboxesPerIP      int `mapstructure:"max_sandboxes_per_ip"`
	MaxSandboxLifetimeMins int `mapstructure:"max_sandbox_lifetime"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

// LoadConfig loads configuration from a file or from environment variables
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(CONFIG_FILE)

	// Init defaults and environment variables
	initConfig()

	// Support environment variables
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println(CONFIG_FILE + " not found!")
		} else {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func initConfig() {
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.SetDefault("server.port", 8080)

	viper.BindEnv("sandbox.url_prefix", "SANDBOX_URL_PREFIX")
	viper.BindEnv("sandbox.url_suffix", "SANDBOX_URL_SUFFIX")
	viper.BindEnv("sandbox.default_lifetime", "SANDBOX_DEFAULT_LIFETIME")
	viper.SetDefault("sandbox.default_lifetime", 60)

	viper.BindEnv("auth.username", "AUTH_USERNAME")
	viper.BindEnv("auth.password", "AUTH_PASSWORD")

	viper.BindEnv("guard.max_total_sandboxes", "GUARD_MAX_TOTAL_SANDBOXES")
	viper.BindEnv("guard.max_sandboxes_per_ip", "GUARD_MAX_SANDBOXES_PER_IP")
	viper.BindEnv("guard.max_sandbox_lifetime", "GUARD_MAX_SANDBOX_LIFETIME")
	viper.SetDefault("guard.max_total_sandboxes", 32)
	viper.SetDefault("guard.max_sandboxes_per_ip", 5)
	viper.SetDefault("guard.max_sandbox_lifetime", 60)
}

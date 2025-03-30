package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Hypervisor HypervisorConfig `yaml:"hypervisor"`
	Services   []ServiceConfig  `yaml:"services"`
}

type HypervisorConfig struct {
	LimaInstance string `yaml:"lima_instance"`
}

type ServiceConfig struct {
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	ComposeFile string `yaml:"compose_file"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func (c *Config) validate() error {
	if c.Hypervisor.LimaInstance == "" {
		return fmt.Errorf("hypervisor.lima_instance is required")
	}

	if len(c.Services) == 0 {
		return fmt.Errorf("at least one service is required")
	}

	for _, svc := range c.Services {
		if err := svc.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceConfig) validate() error {
	if s.Name == "" {
		return fmt.Errorf("service name is required")
	}

	if s.Path == "" {
		return fmt.Errorf("service path is required")
	}

	// Ensure path is absolute
	if !filepath.IsAbs(s.Path) {
		return fmt.Errorf("service path must be absolute: %s", s.Path)
	}

	return nil
} 
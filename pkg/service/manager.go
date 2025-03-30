package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"launchd_docker/pkg/config"
)

type Manager struct {
	services []config.ServiceConfig
	mu       sync.Mutex
	verbose  bool
}

func NewManager(services []config.ServiceConfig) *Manager {
	return &Manager{
		services: services,
	}
}

func (m *Manager) SetVerbose(verbose bool) {
	m.verbose = verbose
}

func (m *Manager) StartAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// First, validate all services have required files
	if m.verbose {
		fmt.Println("Validating service configurations...")
	}
	for _, svc := range m.services {
		if err := m.validateService(svc); err != nil {
			return fmt.Errorf("service %s validation failed: %w", svc.Name, err)
		}
	}
	if m.verbose {
		fmt.Println("Service validation complete")
	}

	// Start services in order
	if m.verbose {
		fmt.Printf("Starting %d services...\n", len(m.services))
	}
	for _, svc := range m.services {
		if err := m.startService(svc); err != nil {
			// Log the error but continue with other services
			fmt.Printf("Error starting service %s: %v\n", svc.Name, err)
		}
	}

	return nil
}

func (m *Manager) StopAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.verbose {
		fmt.Printf("Stopping %d services...\n", len(m.services))
	}

	// Stop services in reverse order
	for i := len(m.services) - 1; i >= 0; i-- {
		svc := m.services[i]
		if err := m.stopService(svc); err != nil {
			// Log the error but continue with other services
			fmt.Printf("Error stopping service %s: %v\n", svc.Name, err)
		}
	}

	if m.verbose {
		fmt.Println("All services stopped")
	}

	return nil
}

func (m *Manager) validateService(svc config.ServiceConfig) error {
	if m.verbose {
		fmt.Printf("Validating service: %s\n", svc.Name)
	}

	// Check if service directory exists
	if _, err := os.Stat(svc.Path); err != nil {
		return fmt.Errorf("service directory does not exist: %w", err)
	}

	// Check for compose file
	composeFile := "docker-compose.yaml"
	if svc.ComposeFile != "" {
		composeFile = svc.ComposeFile
	}

	composePath := filepath.Join(svc.Path, composeFile)
	if _, err := os.Stat(composePath); err != nil {
		return fmt.Errorf("compose file %s not found: %w", composeFile, err)
	}

	if m.verbose {
		fmt.Printf("Service %s validation complete\n", svc.Name)
	}

	return nil
}

func (m *Manager) startService(svc config.ServiceConfig) error {
	if m.verbose {
		fmt.Printf("Starting service: %s\n", svc.Name)
	}

	// Change to service directory
	if err := os.Chdir(svc.Path); err != nil {
		return fmt.Errorf("failed to change to service directory: %w", err)
	}

	// Build compose command
	args := []string{"compose"}
	if svc.ComposeFile != "" {
		args = append(args, "-f", svc.ComposeFile)
	}
	args = append(args, "up", "-d")

	if m.verbose {
		fmt.Printf("Running command: docker %v\n", args)
	}

	// Run docker compose up -d
	cmd := exec.Command("docker", args...)
	cmd.Dir = svc.Path // Ensure we're in the correct directory
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start service: %s: %w", string(output), err)
	}

	if m.verbose {
		fmt.Printf("Service %s started successfully\n", svc.Name)
	}

	return nil
}

func (m *Manager) stopService(svc config.ServiceConfig) error {
	if m.verbose {
		fmt.Printf("Stopping service: %s\n", svc.Name)
	}

	// Change to service directory
	if err := os.Chdir(svc.Path); err != nil {
		return fmt.Errorf("failed to change to service directory: %w", err)
	}

	// Build compose command
	args := []string{"compose"}
	if svc.ComposeFile != "" {
		args = append(args, "-f", svc.ComposeFile)
	}
	args = append(args, "down")

	if m.verbose {
		fmt.Printf("Running command: docker %v\n", args)
	}

	// Run docker compose down
	cmd := exec.Command("docker", args...)
	cmd.Dir = svc.Path // Ensure we're in the correct directory
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop service: %s: %w", string(output), err)
	}

	if m.verbose {
		fmt.Printf("Service %s stopped successfully\n", svc.Name)
	}

	return nil
} 
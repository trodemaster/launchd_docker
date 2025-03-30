package vm

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Manager struct {
	instanceName string
}

func NewManager(instanceName string) *Manager {
	return &Manager{
		instanceName: instanceName,
	}
}

func (m *Manager) EnsureRunning() error {
	// Check if VM is running
	isRunning, err := m.isRunning()
	if err != nil {
		return fmt.Errorf("failed to check VM status: %w", err)
	}

	if isRunning {
		return nil
	}

	// Start the VM
	cmd := exec.Command("limactl", "start", m.instanceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start VM: %s: %w", string(output), err)
	}

	return nil
}

func (m *Manager) WaitForHealthy() error {
	// Wait for VM to be fully started and healthy
	// This is a simple implementation - you might want to add more sophisticated health checks
	maxAttempts := 30
	attempt := 0
	for attempt < maxAttempts {
		cmd := exec.Command("limactl", "list", "--format", "{{.Name}}\t{{.Status}}")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to check VM status: %w", err)
		}

		// Parse output to find our instance
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			fields := strings.Split(line, "\t")
			if len(fields) == 2 && fields[0] == m.instanceName {
				if fields[1] == "Running" {
					return nil
				}
			}
		}

		time.Sleep(2 * time.Second)
		attempt++
	}

	return fmt.Errorf("timeout waiting for VM to be healthy")
}

func (m *Manager) isRunning() (bool, error) {
	cmd := exec.Command("limactl", "list", "--format", "{{.Name}}\t{{.Status}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to check VM status: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		if len(fields) == 2 && fields[0] == m.instanceName {
			return fields[1] == "Running", nil
		}
	}

	return false, nil
} 
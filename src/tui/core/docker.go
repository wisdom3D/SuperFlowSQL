package core

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// ServiceStatus represents the status of a Docker service
type ServiceStatus struct {
	Name   string
	Status string // "running", "stopped", "error", "not found"
	Port   string
	Health string
}

// StackAction represents an action to perform on the stack
type StackAction int

const (
	StackStart StackAction = iota
	StackStop
	StackRestart
)

// DockerComposeCmd runs a docker compose command in the project directory
func DockerComposeCmd(projectDir string, args ...string) (string, error) {
	// Try "docker compose" first (Docker Compose V2), fall back to "docker-compose"
	cmdArgs := append([]string{"compose"}, args...)
	cmd := exec.Command("docker", cmdArgs...)
	cmd.Dir = projectDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Try legacy docker-compose
		cmd2 := exec.Command("docker-compose", args...)
		cmd2.Dir = projectDir
		var stdout2, stderr2 bytes.Buffer
		cmd2.Stdout = &stdout2
		cmd2.Stderr = &stderr2

		err2 := cmd2.Run()
		if err2 != nil {
			return "", fmt.Errorf("docker compose failed: %s\n%s", stderr.String(), stderr2.String())
		}
		return stdout2.String(), nil
	}
	return stdout.String(), nil
}

// StartStack starts all services using docker compose
func StartStack(projectDir string) (string, error) {
	return DockerComposeCmd(projectDir, "up", "-d", "--build")
}

// StopStack stops all services
func StopStack(projectDir string) (string, error) {
	return DockerComposeCmd(projectDir, "down")
}

// GetStackStatus returns the status of all services
func GetStackStatus(projectDir string) ([]ServiceStatus, error) {
	output, err := DockerComposeCmd(projectDir, "ps", "--format", "{{.Name}}|{{.Status}}|{{.Ports}}")
	if err != nil {
		return nil, err
	}

	var services []ServiceStatus
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 3)
		s := ServiceStatus{
			Name:   "unknown",
			Status: "unknown",
			Port:   "",
		}
		if len(parts) >= 1 {
			s.Name = parts[0]
		}
		if len(parts) >= 2 {
			status := strings.ToLower(parts[1])
			if strings.Contains(status, "up") {
				s.Status = "running"
				if strings.Contains(status, "healthy") {
					s.Health = "healthy"
				} else if strings.Contains(status, "unhealthy") {
					s.Health = "unhealthy"
				} else {
					s.Health = "starting"
				}
			} else {
				s.Status = "stopped"
				s.Health = "-"
			}
		}
		if len(parts) >= 3 {
			s.Port = parts[2]
		}

		services = append(services, s)
	}

	// If no services found, return default expected services as "not found"
	if len(services) == 0 {
		services = []ServiceStatus{
			{Name: "postgres", Status: "not found", Health: "-"},
			{Name: "pgadmin", Status: "not found", Health: "-"},
			{Name: "airflow-webserver", Status: "not found", Health: "-"},
			{Name: "airflow-scheduler", Status: "not found", Health: "-"},
			{Name: "superset", Status: "not found", Health: "-"},
		}
	}

	return services, nil
}

// IsDockerRunning checks if Docker daemon is running
func IsDockerRunning() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

// OpenBrowser opens a URL in the default browser
func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Run()
}

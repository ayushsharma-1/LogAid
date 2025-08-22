package plugin

import (
	"context"
	"strings"
)

// DockerPlugin handles Docker command errors
type DockerPlugin struct {
	*BasePlugin
}

// NewDockerPlugin creates a new Docker plugin
func NewDockerPlugin() *DockerPlugin {
	return &DockerPlugin{
		BasePlugin: NewBasePlugin(
			"docker",
			"Handles Docker and Docker Compose errors",
			85,
		),
	}
}

// Match determines if this plugin should handle the command
func (p *DockerPlugin) Match(cmd string, output string, exitCode int) bool {
	if exitCode == 0 {
		return false
	}

	cmdLower := strings.ToLower(strings.TrimSpace(cmd))
	dockerCommands := []string{"docker ", "docker-compose ", "docker compose "}
	
	for _, dockerCmd := range dockerCommands {
		if strings.HasPrefix(cmdLower, dockerCmd) {
			return true
		}
	}

	return false
}

// Suggest generates a suggestion for fixing the Docker error
func (p *DockerPlugin) Suggest(ctx context.Context, cmd string, output string, exitCode int) (string, string, error) {
	outputLower := strings.ToLower(output)
	
	// Check for common Docker errors
	
	// Command typos
	if correctedCmd, explanation := p.checkForTypos(cmd); correctedCmd != "" {
		return correctedCmd, explanation, nil
	}
	
	// Docker daemon not running
	if strings.Contains(outputLower, "cannot connect to the docker daemon") ||
	   strings.Contains(outputLower, "docker daemon socket") {
		return "sudo systemctl start docker", "Start the Docker daemon", nil
	}
	
	// Permission denied
	if strings.Contains(outputLower, "permission denied") && 
	   strings.Contains(outputLower, "docker.sock") {
		return "sudo usermod -aG docker $USER && newgrp docker", 
			   "Add user to docker group for permission", nil
	}
	
	// Image not found
	if strings.Contains(outputLower, "pull access denied") ||
	   strings.Contains(outputLower, "repository does not exist") ||
	   strings.Contains(outputLower, "image not found") {
		imageName := extractImageName(cmd)
		if imageName != "" {
			return "docker search " + imageName, "Search for the correct image name", nil
		}
	}
	
	// Port already in use
	if strings.Contains(outputLower, "port is already allocated") ||
	   strings.Contains(outputLower, "bind: address already in use") {
		return "docker ps", "Check running containers using the port", nil
	}
	
	// Container name already exists
	if strings.Contains(outputLower, "conflict") && 
	   strings.Contains(outputLower, "name") && 
	   strings.Contains(outputLower, "already in use") {
		return "docker ps -a", "Check existing containers with the same name", nil
	}
	
	// Dockerfile not found
	if strings.Contains(outputLower, "no such file") && 
	   strings.Contains(outputLower, "dockerfile") {
		return "ls -la", "Check if Dockerfile exists in current directory", nil
	}
	
	// No space left on device
	if strings.Contains(outputLower, "no space left on device") {
		return "docker system prune", "Clean up unused Docker resources", nil
	}
	
	// Network errors
	if strings.Contains(outputLower, "network") && 
	   strings.Contains(outputLower, "not found") {
		return "docker network ls", "List available Docker networks", nil
	}
	
	// Volume errors
	if strings.Contains(outputLower, "volume") && 
	   strings.Contains(outputLower, "not found") {
		return "docker volume ls", "List available Docker volumes", nil
	}
	
	// Build context errors
	if strings.Contains(outputLower, "unable to prepare context") {
		return "ls -la", "Check files in build context", nil
	}
	
	return "", "No specific Docker suggestion available", nil
}

// checkForTypos checks for common typos in Docker commands
func (p *DockerPlugin) checkForTypos(cmd string) (string, string) {
	corrections := map[string]string{
		"run":         "run",
		"rnu":         "run",
		"runn":        "run",
		"build":       "build",
		"buil":        "build",
		"buidl":       "build",
		"biuld":       "build",
		"pull":        "pull",
		"pul":         "pull",
		"pll":         "pull",
		"push":        "push",
		"pus":         "push",
		"psh":         "push",
		"images":      "images",
		"image":       "images",
		"imgs":        "images",
		"ps":          "ps",
		"p":           "ps",
		"start":       "start",
		"stat":        "start",
		"satrt":       "start",
		"stop":        "stop",
		"stp":         "stop",
		"sotp":        "stop",
		"restart":     "restart",
		"restat":      "restart",
		"resart":      "restart",
		"remove":      "remove",
		"remov":       "remove",
		"rm":          "rm",
		"logs":        "logs",
		"log":         "logs",
		"lgs":         "logs",
		"exec":        "exec",
		"exe":         "exec",
		"exce":        "exec",
		"inspect":     "inspect",
		"inspec":      "inspect",
		"insect":      "inspect",
		"network":     "network",
		"netwrk":      "network",
		"netowrk":     "network",
		"volume":      "volume",
		"volum":       "volume",
		"volme":       "volume",
		"compose":     "compose",
		"compos":      "compose",
		"compse":      "compose",
		"system":      "system",
		"sytem":       "system",
		"systm":       "system",
	}

	words := strings.Fields(cmd)
	changed := false
	
	for i, word := range words {
		cleanWord := strings.ToLower(word)
		if correction, exists := corrections[cleanWord]; exists {
			words[i] = correction
			changed = true
		}
	}

	if changed {
		newCmd := strings.Join(words, " ")
		return newCmd, "Fixed typos in Docker command"
	}

	return "", ""
}

// extractImageName extracts Docker image name from command
func extractImageName(cmd string) string {
	words := strings.Fields(cmd)
	
	for i, word := range words {
		if (word == "run" || word == "pull" || word == "push" || word == "build") && i+1 < len(words) {
			// Skip flags and return the first non-flag argument
			for j := i + 1; j < len(words); j++ {
				if !strings.HasPrefix(words[j], "-") {
					return words[j]
				}
			}
		}
	}
	
	return ""
}

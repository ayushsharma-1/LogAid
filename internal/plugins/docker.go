package plugins

import (
	"context"
	"fmt"
	"strings"

	"github.com/ayushsharma-1/LogAid/internal/ai"
)

// DockerPlugin handles Docker command errors with AI-powered suggestions
type DockerPlugin struct{}

func (p *DockerPlugin) Name() string {
	return "docker"
}

// Match checks if this plugin should handle the command/output
func (p *DockerPlugin) Match(cmd string, output string) bool {
	// Check if command uses docker
	if !strings.Contains(strings.ToLower(cmd), "docker") {
		return false
	}

	// Check for common docker errors
	dockerErrors := []string{
		"unable to find image",
		"is not a docker command",
		"permission denied while trying to connect to the docker daemon",
		"cannot connect to the docker daemon",
		"docker daemon not running",
		"no such container",
		"no such image",
		"error response from daemon",
		"pull access denied",
		"repository does not exist",
		"unauthorized",
		"manifest unknown",
		"tag does not exist",
	}

	return containsAny(output, dockerErrors)
}

// Suggest generates an AI-powered suggestion for the error
func (p *DockerPlugin) Suggest(cmd string, output string) string {
	// First try manual corrections for speed
	if quickFix := p.getQuickFix(cmd, output); quickFix != "" {
		return quickFix
	}

	// Use AI for complex suggestions
	return p.getAISuggestion(cmd, output)
}

// getQuickFix provides immediate fixes for common issues
func (p *DockerPlugin) getQuickFix(cmd string, output string) string {
	outputLower := strings.ToLower(output)

	// Handle permission errors
	if strings.Contains(outputLower, "permission denied") && !strings.Contains(cmd, "sudo") {
		return "sudo " + cmd
	}

	// Handle daemon not running
	if strings.Contains(outputLower, "cannot connect to the docker daemon") {
		return "sudo systemctl start docker && " + cmd
	}

	// Handle command typos
	if strings.Contains(outputLower, "is not a docker command") {
		return p.correctDockerCommand(cmd)
	}

	// Handle image name typos
	if strings.Contains(outputLower, "unable to find image") {
		return p.correctImageName(cmd, output)
	}

	return ""
}

// correctDockerCommand fixes common Docker command typos
func (p *DockerPlugin) correctDockerCommand(cmd string) string {
	corrections := map[string]string{
		"ru":    "run",
		"rn":    "run",
		"buil":  "build",
		"buid":  "build",
		"pul":   "pull",
		"pll":   "pull",
		"pus":   "push",
		"psh":   "push",
		"exe":   "exec",
		"exec":  "exec",
		"p":     "ps",
		"log":   "logs",
		"stp":   "stop",
		"stop":  "stop",
		"stat":  "start",
		"strt":  "start",
		"rm":    "rm",
		"rmi":   "rmi",
		"img":   "images",
		"image": "images",
		"net":   "network",
		"vol":   "volume",
		"cp":    "cp",
		"insp":  "inspect",
		"inspt": "inspect",
	}

	parts := strings.Fields(cmd)
	if len(parts) >= 2 {
		command := parts[1]
		if correction, exists := corrections[command]; exists {
			parts[1] = correction
			return strings.Join(parts, " ")
		}
	}

	return cmd
}

// correctImageName fixes common Docker image name typos
func (p *DockerPlugin) correctImageName(cmd string, output string) string {
	imageCorrections := map[string]string{
		"ubntu":      "ubuntu",
		"ubunt":      "ubuntu",
		"ubunut":     "ubuntu",
		"ngnix":      "nginx",
		"ngin":       "nginx",
		"nginc":      "nginx",
		"alpin":      "alpine",
		"alpne":      "alpine",
		"redi":       "redis",
		"redis":      "redis",
		"rediss":     "redis",
		"postgre":    "postgres",
		"postgrs":    "postgres",
		"postgresql": "postgres",
		"mysq":       "mysql",
		"mysl":       "mysql",
		"mysql":      "mysql",
		"mong":       "mongo",
		"mongo":      "mongo",
		"mongod":     "mongo",
		"node":       "node",
		"nodejs":     "node",
		"pythn":      "python",
		"pythno":     "python",
		"pyhton":     "python",
		"centso":     "centos",
		"cenot":      "centos",
		"deban":      "debian",
		"debain":     "debian",
		"fedra":      "fedora",
		"fedro":      "fedora",
		"archlinx":   "archlinux",
		"arch":       "archlinux",
	}

	// Extract image name from output
	for typo, correct := range imageCorrections {
		if strings.Contains(strings.ToLower(output), typo) {
			return strings.Replace(cmd, typo, correct, 1)
		}
	}

	return cmd
}

// getAISuggestion uses AI to generate intelligent suggestions
func (p *DockerPlugin) getAISuggestion(cmd string, output string) string {
	prompt := p.buildAIPrompt(cmd, output)

	ctx := context.Background()
	suggestion, err := ai.GetSuggestion(ctx, prompt)
	if err != nil {
		// Fallback to generic suggestion
		return "docker --help # Check the correct Docker command syntax"
	}

	return suggestion
}

// buildAIPrompt creates a detailed prompt for the AI
func (p *DockerPlugin) buildAIPrompt(cmd string, output string) string {
	return fmt.Sprintf(`
You are an expert Docker administrator and DevOps engineer.

CONTEXT:
- User executed command: %s
- Command output/error: %s
- System: Linux with Docker installed
- Goal: Provide the EXACT corrected command to fix the issue

TASK:
Analyze the command and error, then provide a single, executable command that will resolve the issue.

RULES:
1. Return ONLY the corrected command, no explanations
2. Use proper Docker syntax and image names
3. Include sudo if needed for permissions
4. Handle common issues: typos, missing images, daemon not running, permission errors
5. If image doesn't exist, suggest the closest alternative
6. For service issues, suggest the complete fix including service management
7. Always prioritize safety and standard practices

COMMON DOCKER PATTERNS TO CONSIDER:
- Image name typos (ubntu → ubuntu, ngnix → nginx)
- Command typos (ru → run, pul → pull, pus → push)
- Missing sudo for daemon access
- Docker daemon not running (need to start service)
- Image tag issues
- Port binding problems
- Volume mount issues

EXAMPLES:
- Input: "docker ru ubuntu" + "docker: 'ru' is not a docker command"
- Output: "docker run ubuntu"

- Input: "docker run ubuntu" + "permission denied while trying to connect"
- Output: "sudo docker run ubuntu"

- Input: "docker run ubntu" + "Unable to find image 'ubntu:latest'"
- Output: "docker run ubuntu"

Provide the corrected command:`, cmd, output)
}

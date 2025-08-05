package plugins

import (
	"context"
	"fmt"
	"strings"

	"github.com/ayushsharma-1/LogAid/internal/ai"
)

// PipPlugin handles Python pip command errors with AI-powered suggestions
type PipPlugin struct{}

func (p *PipPlugin) Name() string {
	return "pip"
}

// Match checks if this plugin should handle the command/output
func (p *PipPlugin) Match(cmd string, output string) bool {
	// Check if command uses pip
	if !strings.Contains(strings.ToLower(cmd), "pip") {
		return false
	}

	// Check for common pip errors
	pipErrors := []string{
		"no such option:",
		"unknown command",
		"could not find a version",
		"no matching distribution found",
		"permission denied",
		"externally-managed-environment",
		"pip: command not found",
		"error: could not install packages",
		"certificate verify failed",
		"connection error",
		"timeout",
		"requirement already satisfied",
		"syntax error in requirements",
		"invalid requirement",
		"pip is being invoked by an old script",
	}

	return containsAny(output, pipErrors)
}

// Suggest generates an AI-powered suggestion for the error
func (p *PipPlugin) Suggest(cmd string, output string) string {
	// First try manual corrections for speed
	if quickFix := p.getQuickFix(cmd, output); quickFix != "" {
		return quickFix
	}

	// Use AI for complex suggestions
	return p.getAISuggestion(cmd, output)
}

// getQuickFix provides immediate fixes for common issues
func (p *PipPlugin) getQuickFix(cmd string, output string) string {
	outputLower := strings.ToLower(output)

	// Handle pip not found
	if strings.Contains(outputLower, "pip: command not found") {
		return "sudo apt install python3-pip && " + strings.Replace(cmd, "pip", "pip3", 1)
	}

	// Handle permission errors
	if strings.Contains(outputLower, "permission denied") {
		if !strings.Contains(cmd, "--user") && !strings.Contains(cmd, "sudo") {
			return cmd + " --user"
		}
	}

	// Handle externally managed environment (Ubuntu 23.04+)
	if strings.Contains(outputLower, "externally-managed-environment") {
		return cmd + " --break-system-packages # or use: python3 -m venv myenv && source myenv/bin/activate && " + cmd
	}

	// Handle package name corrections
	if strings.Contains(outputLower, "could not find a version") || strings.Contains(outputLower, "no matching distribution") {
		return p.correctPackageName(cmd)
	}

	// Handle pip vs pip3
	if strings.Contains(cmd, "pip ") && !strings.Contains(cmd, "pip3") {
		return strings.Replace(cmd, "pip ", "pip3 ", 1)
	}

	return ""
}

// correctPackageName fixes common Python package name typos
func (p *PipPlugin) correctPackageName(cmd string) string {
	packageCorrections := map[string]string{
		// Popular Python packages with common typos
		"beautifulsoup":   "beautifulsoup4",
		"bs4":             "beautifulsoup4",
		"beautiful-soup":  "beautifulsoup4",
		"request":         "requests",
		"requets":         "requests",
		"reqeusts":        "requests",
		"numpy":           "numpy",
		"numpi":           "numpy",
		"numpyy":          "numpy",
		"pandas":          "pandas",
		"panda":           "pandas",
		"pandass":         "pandas",
		"matplotlib":      "matplotlib",
		"matplot":         "matplotlib",
		"matplotlb":       "matplotlib",
		"scipy":           "scipy",
		"scipi":           "scipy",
		"scypy":           "scipy",
		"scikit-learn":    "scikit-learn",
		"sklearn":         "scikit-learn",
		"scikit":          "scikit-learn",
		"tensorflow":      "tensorflow",
		"tensorflw":       "tensorflow",
		"tensoflow":       "tensorflow",
		"torch":           "torch",
		"pytorch":         "torch",
		"pyyaml":          "pyyaml",
		"yaml":            "pyyaml",
		"yml":             "pyyaml",
		"flask":           "flask",
		"flsk":            "flask",
		"flaskk":          "flask",
		"django":          "django",
		"djnago":          "django",
		"djangoo":         "django",
		"fastapi":         "fastapi",
		"fastap":          "fastapi",
		"fast-api":        "fastapi",
		"sqlalchemy":      "sqlalchemy",
		"sqlalchmy":       "sqlalchemy",
		"sql-alchemy":     "sqlalchemy",
		"pillow":          "pillow",
		"pil":             "pillow",
		"pillw":           "pillow",
		"opencv-python":   "opencv-python",
		"opencv":          "opencv-python",
		"cv2":             "opencv-python",
		"jupyter":         "jupyter",
		"jupytr":          "jupyter",
		"jupyterr":        "jupyter",
		"ipython":         "ipython",
		"ipythoon":        "ipython",
		"py-python":       "ipython",
		"pytz":            "pytz",
		"pyttz":           "pytz",
		"timezone":        "pytz",
		"dateutil":        "python-dateutil",
		"python-dateutil": "python-dateutil",
		"date-util":       "python-dateutil",
		"click":           "click",
		"clik":            "click",
		"clickk":          "click",
		"setuptools":      "setuptools",
		"setup-tools":     "setuptools",
		"setuptool":       "setuptools",
		"wheel":           "wheel",
		"whel":            "wheel",
		"wheell":          "wheel",
		"virtualenv":      "virtualenv",
		"virtual-env":     "virtualenv",
		"venv":            "virtualenv",
		"pipenv":          "pipenv",
		"pip-env":         "pipenv",
		"pipenev":         "pipenv",
	}

	// Try to extract package name and correct it
	parts := strings.Fields(cmd)
	for i, part := range parts {
		if part == "install" {
			if i+1 < len(parts) {
				packageName := parts[i+1]
				// Remove flags and get clean package name
				cleanPackage := strings.Split(packageName, "==")[0]
				cleanPackage = strings.Split(cleanPackage, ">=")[0]
				cleanPackage = strings.Split(cleanPackage, "<=")[0]
				cleanPackage = strings.Split(cleanPackage, ">")[0]
				cleanPackage = strings.Split(cleanPackage, "<")[0]
				cleanPackage = strings.Split(cleanPackage, "!=")[0]

				if correction, exists := packageCorrections[cleanPackage]; exists {
					parts[i+1] = strings.Replace(packageName, cleanPackage, correction, 1)
					return strings.Join(parts, " ")
				}
			}
		}
	}

	return cmd
}

// getAISuggestion uses AI to generate intelligent suggestions
func (p *PipPlugin) getAISuggestion(cmd string, output string) string {
	prompt := p.buildAIPrompt(cmd, output)

	ctx := context.Background()
	suggestion, err := ai.GetSuggestion(ctx, prompt)
	if err != nil {
		// Fallback to generic suggestion
		return "pip3 --help # Check the correct pip command syntax"
	}

	return suggestion
}

// buildAIPrompt creates a detailed prompt for the AI
func (p *PipPlugin) buildAIPrompt(cmd string, output string) string {
	return fmt.Sprintf(`
You are an expert Python developer and package management specialist.

CONTEXT:
- User executed command: %s
- Command output/error: %s
- System: Linux with Python and pip package manager
- Goal: Provide the EXACT corrected command to fix the issue

TASK:
Analyze the command and error, then provide a single, executable command that will resolve the issue.

RULES:
1. Return ONLY the corrected command, no explanations
2. Use proper pip/pip3 syntax and package names
3. Handle common issues: typos, missing packages, permission errors, version conflicts
4. If package doesn't exist, suggest the closest alternative
5. For permission issues, suggest --user flag or virtual environment
6. For system management issues, provide appropriate workarounds
7. Always prioritize safety and best practices

COMMON PIP PATTERNS TO CONSIDER:
- pip vs pip3 usage (prefer pip3 for Python 3)
- Package name typos (request → requests, beautifulsoup → beautifulsoup4)
- Permission issues (use --user flag)
- Externally managed environments (Ubuntu 23.04+)
- Virtual environment recommendations
- Version specification syntax
- Requirements file issues

EXAMPLES:
- Input: "pip install request" + "Could not find a version that satisfies the requirement request"
- Output: "pip3 install requests"

- Input: "pip install numpy" + "Permission denied"
- Output: "pip3 install numpy --user"

- Input: "pip install flask" + "externally-managed-environment"
- Output: "python3 -m venv myenv && source myenv/bin/activate && pip install flask"

Provide the corrected command:`, cmd, output)
}

// SystemctlPlugin handles systemctl service management errors
type SystemctlPlugin struct{}

func (p *SystemctlPlugin) Name() string {
	return "systemctl"
}

// Match checks if this plugin should handle the command/output
func (p *SystemctlPlugin) Match(cmd string, output string) bool {
	// Check if command uses systemctl
	if !strings.Contains(strings.ToLower(cmd), "systemctl") {
		return false
	}

	// Check for common systemctl errors
	systemctlErrors := []string{
		"unit not found",
		"failed to start",
		"failed to stop",
		"failed to restart",
		"failed to reload",
		"permission denied",
		"authentication required",
		"could not find",
		"unknown operation",
		"invalid option",
		"unit file not found",
		"masked unit",
		"inactive unit",
		"job failed",
	}

	return containsAny(output, systemctlErrors)
}

// Suggest generates an AI-powered suggestion for the error
func (p *SystemctlPlugin) Suggest(cmd string, output string) string {
	// First try manual corrections for speed
	if quickFix := p.getQuickFix(cmd, output); quickFix != "" {
		return quickFix
	}

	// Use AI for complex suggestions
	return p.getAISuggestion(cmd, output)
}

// getQuickFix provides immediate fixes for common issues
func (p *SystemctlPlugin) getQuickFix(cmd string, output string) string {
	outputLower := strings.ToLower(output)

	// Handle permission errors
	if strings.Contains(outputLower, "permission denied") || strings.Contains(outputLower, "authentication required") {
		if !strings.Contains(cmd, "sudo") {
			return "sudo " + cmd
		}
	}

	// Handle service name corrections
	if strings.Contains(outputLower, "unit not found") || strings.Contains(outputLower, "could not find") {
		return p.correctServiceName(cmd)
	}

	// Handle masked units
	if strings.Contains(outputLower, "masked unit") {
		parts := strings.Fields(cmd)
		if len(parts) >= 3 {
			serviceName := parts[2]
			return fmt.Sprintf("sudo systemctl unmask %s && %s", serviceName, cmd)
		}
	}

	return ""
}

// correctServiceName fixes common service name typos
func (p *SystemctlPlugin) correctServiceName(cmd string) string {
	serviceCorrections := map[string]string{
		"apache":     "apache2",
		"httpd":      "apache2",
		"nginx":      "nginx",
		"ngnix":      "nginx",
		"docker":     "docker",
		"dockerd":    "docker",
		"mysql":      "mysql",
		"mariadb":    "mariadb",
		"postgresql": "postgresql",
		"postgres":   "postgresql",
		"redis":      "redis-server",
		"redis-srv":  "redis-server",
		"ssh":        "ssh",
		"sshd":       "ssh",
		"openssh":    "ssh",
		"network":    "networking",
		"net":        "networking",
		"firewall":   "ufw",
		"iptables":   "iptables",
		"cron":       "cron",
		"crond":      "cron",
		"systemd":    "systemd",
		"dbus":       "dbus",
		"avahi":      "avahi-daemon",
		"bluetooth":  "bluetooth",
		"cups":       "cups",
		"printer":    "cups",
	}

	parts := strings.Fields(cmd)
	if len(parts) >= 3 {
		serviceName := parts[2]
		// Remove .service suffix if present
		cleanService := strings.TrimSuffix(serviceName, ".service")

		if correction, exists := serviceCorrections[cleanService]; exists {
			parts[2] = correction + ".service"
			return strings.Join(parts, " ")
		}

		// If no exact match, try without .service suffix
		if !strings.HasSuffix(serviceName, ".service") {
			parts[2] = cleanService + ".service"
			return strings.Join(parts, " ")
		}
	}

	return cmd
}

// getAISuggestion uses AI to generate intelligent suggestions
func (p *SystemctlPlugin) getAISuggestion(cmd string, output string) string {
	prompt := p.buildAIPrompt(cmd, output)

	ctx := context.Background()
	suggestion, err := ai.GetSuggestion(ctx, prompt)
	if err != nil {
		// Fallback to generic suggestion
		return "systemctl --help # Check the correct systemctl command syntax"
	}

	return suggestion
}

// buildAIPrompt creates a detailed prompt for the AI
func (p *SystemctlPlugin) buildAIPrompt(cmd string, output string) string {
	return fmt.Sprintf(`
You are an expert Linux system administrator specializing in systemd service management.

CONTEXT:
- User executed command: %s
- Command output/error: %s
- System: Linux with systemd service manager
- Goal: Provide the EXACT corrected command to fix the issue

TASK:
Analyze the command and error, then provide a single, executable command that will resolve the issue.

RULES:
1. Return ONLY the corrected command, no explanations
2. Use proper systemctl syntax and service names
3. Include sudo if needed for permissions
4. Handle common issues: typos, missing services, permission errors, masked units
5. If service doesn't exist, suggest the closest alternative
6. For masked units, provide unmask command first
7. Always prioritize safety and standard practices

COMMON SYSTEMCTL PATTERNS TO CONSIDER:
- Service name corrections (apache → apache2, mysql → mysql)
- Missing .service suffix
- Permission issues requiring sudo
- Masked units needing to be unmasked
- Service not found vs service not enabled
- Start/stop/restart/reload commands
- Enable/disable for boot time behavior

EXAMPLES:
- Input: "systemctl start apache" + "Unit apache.service not found"
- Output: "sudo systemctl start apache2.service"

- Input: "systemctl restart nginx" + "Permission denied"
- Output: "sudo systemctl restart nginx"

- Input: "systemctl start docker" + "Unit is masked"
- Output: "sudo systemctl unmask docker && sudo systemctl start docker"

Provide the corrected command:`, cmd, output)
}

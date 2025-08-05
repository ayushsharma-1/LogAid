# Contributing to LogAid

Thank you for your interest in contributing to LogAid! This document provides guidelines and information for contributors.

## 📋 Table of Contents
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Guidelines](#contributing-guidelines)
- [Plugin Development](#plugin-development)
- [Testing](#testing)
- [Code Style](#code-style)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## 🚀 Getting Started

### Prerequisites

- **Go**: Version 1.21 or higher
- **Git**: For version control
- **Make**: For build automation (optional)
- **Docker**: For testing Docker-related features (optional)

### Quick Setup

```bash
# Clone the repository
git clone https://github.com/ayushsharma-1/LogAid.git
cd LogAid

# Install dependencies
go mod download

# Build the project
go build -o logaid .

# Run tests
go test ./...
```

## 🛠️ Development Setup

### Environment Configuration

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Configure your API keys (optional for core development):
```bash
# For AI features (optional)
GEMINI_API_KEY=your_api_key_here
# or
OPENAI_API_KEY=your_api_key_here
```

3. Set development-friendly options:
```bash
DEBUG_MODE=true
LOG_LEVEL=debug
TEST_MODE=true
```

### IDE Setup

#### VS Code
Recommended extensions:
- Go (golang.go)
- gopls language server
- Test Explorer for Go
- GitLens

#### GoLand/IntelliJ
- Enable Go modules support
- Configure code style (see below)
- Enable test coverage display

## 📖 Contributing Guidelines

### Types of Contributions

We welcome several types of contributions:

1. **🐛 Bug Reports**: Issues with existing functionality
2. **✨ Feature Requests**: New features or enhancements
3. **🔌 Plugin Development**: New tool/command support
4. **📚 Documentation**: Improvements to docs
5. **🧪 Tests**: Additional test coverage
6. **🔧 Infrastructure**: Build, CI/CD improvements

### Before You Start

1. **Check existing issues**: Avoid duplicate work
2. **Discuss large changes**: Open an issue for discussion
3. **Follow conventions**: Read this guide thoroughly
4. **Start small**: Begin with minor fixes or improvements

## 🔌 Plugin Development

### Creating a New Plugin

Plugins are the core of LogAid's extensibility. Here's how to create one:

#### 1. Plugin Structure

```go
package plugins

import (
    "context"
    "strings"
    
    "github.com/ayushsharma-1/LogAid/internal/ai"
)

// YourToolPlugin handles errors for YourTool commands
type YourToolPlugin struct{}

func (p *YourToolPlugin) Name() string {
    return "yourtool"
}

func (p *YourToolPlugin) Match(cmd string, output string) bool {
    // Check if this plugin should handle this error
    if !strings.Contains(cmd, "yourtool") {
        return false
    }
    
    // Check for error patterns
    errorPatterns := []string{
        "command not found",
        "invalid option",
        // Add your patterns
    }
    
    return containsAny(output, errorPatterns)
}

func (p *YourToolPlugin) Suggest(cmd string, output string) string {
    // Try quick fixes first
    if quickFix := p.getQuickFix(cmd, output); quickFix != "" {
        return quickFix
    }
    
    // Fall back to AI for complex cases
    return p.getAISuggestion(cmd, output)
}
```

#### 2. Quick Fix Examples

```go
func (p *YourToolPlugin) getQuickFix(cmd string, output string) string {
    // Common typo corrections
    corrections := map[string]string{
        "yourtool instal": "yourtool install",
        "yourtool updat":  "yourtool update",
        // Add more corrections
    }
    
    for typo, fix := range corrections {
        if strings.Contains(cmd, typo) {
            return strings.Replace(cmd, typo, fix, 1)
        }
    }
    
    // Permission fixes
    if strings.Contains(output, "permission denied") {
        return "sudo " + cmd
    }
    
    return ""
}
```

#### 3. AI Integration

```go
func (p *YourToolPlugin) getAISuggestion(cmd string, output string) string {
    prompt := fmt.Sprintf(`You are a YourTool expert. Fix this command:

Command: %s
Error: %s

Rules:
- Only return the corrected command
- No explanations or markdown
- Focus on common YourTool issues
- Consider package names, syntax, permissions

Corrected command:`, cmd, output)

    ctx := context.Background()
    suggestion, err := ai.GetSuggestion(ctx, prompt)
    if err != nil {
        return ""
    }
    
    return suggestion
}
```

#### 4. Register Your Plugin

Add your plugin to `internal/plugins/plugins.go`:

```go
func LoadAllPlugins() []Plugin {
    return []Plugin{
        &AptPlugin{},
        &GitPlugin{},
        &DockerPlugin{},
        &YourToolPlugin{}, // Add here
        // ... other plugins
    }
}
```

### Plugin Testing

Create comprehensive tests for your plugin:

```go
func TestYourToolPlugin(t *testing.T) {
    plugin := &YourToolPlugin{}
    
    tests := []struct {
        name           string
        cmd            string
        output         string
        expectMatch    bool
        expectSuggestion string
    }{
        {
            name:           "typo in command",
            cmd:            "yourtool instal package",
            output:         "unknown command: instal",
            expectMatch:    true,
            expectSuggestion: "yourtool install package",
        },
        // Add more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            match := plugin.Match(tt.cmd, tt.output)
            assert.Equal(t, tt.expectMatch, match)
            
            if match {
                suggestion := plugin.Suggest(tt.cmd, tt.output)
                assert.Equal(t, tt.expectSuggestion, suggestion)
            }
        })
    }
}
```

## 🧪 Testing

### Test Categories

1. **Unit Tests**: Test individual components
2. **Integration Tests**: Test component interactions
3. **Plugin Tests**: Test specific plugin behavior
4. **Performance Tests**: Ensure speed requirements

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestAptPlugin ./tests/

# Run benchmarks
go test -bench=. ./tests/

# Run with race detection
go test -race ./...
```

### Test Requirements

- **Coverage**: Maintain >80% code coverage
- **Performance**: Plugin tests must run <1ms each
- **Isolation**: Tests should not depend on external services
- **Documentation**: Complex tests need explanatory comments

### Mock Requirements

For AI-dependent tests, use test mode:

```go
func TestWithAI(t *testing.T) {
    os.Setenv("TEST_MODE", "true")
    defer os.Unsetenv("TEST_MODE")
    
    // Your test code here
}
```

## 🎨 Code Style

### Go Style Guidelines

We follow standard Go conventions plus these specifics:

#### 1. Formatting
```bash
# Format code
go fmt ./...

# Run linters
golangci-lint run
```

#### 2. Naming Conventions
- **Packages**: lowercase, single word when possible
- **Files**: lowercase with underscores
- **Types**: PascalCase
- **Functions**: PascalCase (exported), camelCase (private)
- **Constants**: SCREAMING_SNAKE_CASE for public, camelCase for private

#### 3. Documentation
```go
// Package plugins provides error detection and correction for various CLI tools.
package plugins

// Plugin represents a tool-specific error handler that can detect and suggest
// fixes for command-line errors.
type Plugin interface {
    // Match determines if this plugin should handle the given command error.
    // It returns true if the plugin recognizes the error pattern.
    Match(cmd string, output string) bool
    
    // Suggest generates a corrected command based on the error.
    // It should return an empty string if no suggestion can be made.
    Suggest(cmd string, output string) string
    
    // Name returns the plugin's identifier for logging and debugging.
    Name() string
}
```

#### 4. Error Handling
```go
// Preferred: Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to parse command %q: %w", cmd, err)
}

// Avoid: Silent failures
if err != nil {
    return ""
}
```

#### 5. Logging
```go
// Use structured logging
logger.Info("Plugin loaded", 
    "name", plugin.Name(),
    "commands", len(plugin.Commands()))

// Avoid string concatenation
logger.Info("Plugin " + name + " loaded") // Don't do this
```

## 📤 Submitting Changes

### Pull Request Process

1. **Fork & Branch**
```bash
git checkout -b feature/your-feature-name
```

2. **Make Changes**
- Follow code style guidelines
- Add tests for new functionality
- Update documentation if needed

3. **Test Thoroughly**
```bash
go test ./...
go test -race ./...
go test -cover ./...
```

4. **Commit with Clear Messages**
```bash
git commit -m "feat: add support for new-tool commands

- Implement NewToolPlugin with pattern matching
- Add tests for common error scenarios
- Update documentation

Closes #123"
```

#### Commit Message Format
```
type(scope): short description

Longer description if needed.

- Bullet points for changes
- Reference issues with #123
- Break lines at 72 characters

Closes #issue-number
```

**Types**: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `chore`

5. **Submit Pull Request**
- Clear title and description
- Reference related issues
- Include screenshots if UI changes
- Ensure CI passes

### PR Review Process

1. **Automated Checks**: CI must pass
2. **Code Review**: At least one maintainer approval
3. **Testing**: Manual testing for significant changes
4. **Documentation**: Updates must be included

## 🚀 Release Process

### Version Strategy

We use [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

### Release Checklist

For maintainers releasing new versions:

1. **Pre-release**
```bash
# Update version in relevant files
# Update CHANGELOG.md
# Run full test suite
go test ./...
```

2. **Create Release**
```bash
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

3. **GitHub Actions** will automatically:
- Build binaries for multiple platforms
- Create GitHub release
- Publish to package registries

### Binary Distribution

Release binaries are provided for:
- Linux (amd64, arm64)
- macOS (amd64, arm64) 
- Windows (amd64)

## 🤝 Community

### Getting Help

- **Issues**: GitHub issues for bugs and features
- **Discussions**: GitHub discussions for questions
- **Discord**: [Coming soon] Real-time chat

### Code of Conduct

We follow the [Contributor Covenant](https://www.contributor-covenant.org/). Be respectful, inclusive, and collaborative.

### Recognition

Contributors are recognized in:
- CHANGELOG.md for each release
- GitHub contributors page
- Special recognition for significant contributions

---

Thank you for contributing to LogAid! Your help makes this tool better for everyone. 🙏

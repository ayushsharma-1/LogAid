# Changelog

All notable changes to LogAid will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.2] - 2025-09-08

### Added
- 🔒 **Enhanced Security Sanitization System**
  - Comprehensive sensitive data detection and redaction
  - Advanced pattern matching for API keys, tokens, passwords, and SSH keys
  - Risk level assessment (Low, Medium, High, Critical)
  - User consent flows for sensitive data handling
  - Sanitization integration across all plugins and PTY wrapper
  - Support for environment variables, database URLs, and cloud credentials
- 🔌 **Plugin Security Enhancements**
  - `SuggestSecure` methods for all plugins (apt, docker, git, npm, kubernetes, generic)
  - Individual plugin sanitization with context-aware patterns
  - Security-first approach with automatic sensitive data detection
- 🛡️ **AI Provider Security**
  - Pre-flight security checks before sending data to AI services
  - Configurable security consent prompts
  - Sanitized data transmission to protect user privacy

### Security
- 15+ sensitive data patterns including:
  - API keys and access tokens
  - SSH private keys and certificates
  - Database connection strings
  - Cloud provider credentials (AWS, GCP, Azure)
  - Authentication headers and bearer tokens
  - Environment variables with sensitive names
  - Password patterns and hashes
- User consent requirements for high-risk data
- Automatic data sanitization with placeholder replacement
- Security logging and audit trails

## [1.0.0] - 2025-08-23

### Added
- ✨ Initial release of LogAid
- 🤖 Dual AI provider support (Google Gemini + OpenAI GPT)
- 🔄 Automatic fallback between AI providers
- 🔌 Plugin system for common command errors
  - apt (Ubuntu/Debian package manager)
  - git (Version control)
  - npm (Node.js package manager)
  - docker (Container platform)
  - kubernetes (Container orchestration)
  - Generic command typo corrections
- 🖥️ PTY wrapper for real-time shell monitoring
- 📊 Comprehensive logging and analytics
- 🎨 Colored terminal output with user-friendly interface
- ⚙️ Environment-based configuration system
- 📋 CLI commands: start, config, logs, help, version
- 🧪 Comprehensive test suite with integration tests
- 🚀 Multi-platform build support (Linux, macOS, Windows)
- 📦 Automated CI/CD with GitHub Actions
- 🔒 Security scanning and code quality checks
- 📖 Complete documentation and development guide

### Technical Features
- Built with Go 1.21+ for performance and reliability
- Modular architecture with clear separation of concerns
- JSON-based configuration and logging
- HTTP client with timeout and retry logic
- Signal handling for graceful shutdown
- Cross-platform compatibility
- Memory-efficient command monitoring
- Real-time error detection and suggestion

### Security
- No hardcoded API keys in source code
- Environment variable-based configuration
- Input validation and sanitization
- Secure HTTP communications with AI providers
- Minimal permission requirements

### Documentation
- README with quick start guide
- Development guide with contributing instructions
- .env.example with all configuration options
- API documentation for plugin development
- GitHub Actions workflow documentation

[1.0.0]: https://github.com/ayushsharma-1/LogAid/releases/tag/v1.0.0

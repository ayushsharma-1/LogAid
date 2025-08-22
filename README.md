# LogAid 🛡️

> Your intelligent CLI companion that learns from your mistakes and suggests fixes in real-time.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI/CD](https://github.com/ayushsharma-1/LogAid/actions/workflows/ci.yml/badge.svg)](https://github.com/ayushsharma-1/LogAid/actions)
[![Release](https://img.shields.io/github/v/release/ayushsharma-1/LogAid)](https://github.com/ayushsharma-1/LogAid/releases)

LogAid is an AI-powered command-line assistant that monitors your shell commands, detects errors, and provides intelligent suggestions to fix them. With support for multiple AI providers and a comprehensive plugin system, LogAid learns from common mistakes and helps you become more productive.

## ✨ Features

### 🤖 **Dual AI Provider Support**
- **Google Gemini API** integration with intelligent prompting
- **OpenAI GPT** support with automatic fallback
- Smart provider switching when one fails
- Optimized prompts for command-line error correction

### 🔌 **Comprehensive Plugin System**
- **apt** - Ubuntu/Debian package management errors
- **git** - Version control command suggestions
- **npm** - Node.js package manager assistance
- **docker** - Container platform error resolution
- **kubernetes** - K8s command corrections
- **Generic** - Common command typos and fixes

### 🖥️ **Real-time Shell Monitoring**
- PTY wrapper for seamless command interception
- Real-time error detection and analysis
- Non-intrusive monitoring that doesn't affect workflow
- Support for bash, zsh, and other popular shells

### 🎨 **Beautiful User Interface**
- Colored terminal output with clear suggestions
- Interactive prompts with timeout support
- Progress indicators and status updates
- ASCII art logo and branded experience

### 📊 **Analytics & Logging**
- Comprehensive command history tracking
- Success/failure analytics
- Performance metrics and timing
- Configurable log levels and output formats

## 🚀 Quick Start

### Prerequisites
- Go 1.20+ installed
- At least one AI API key (Google Gemini or OpenAI)
- Linux, macOS, or Windows

### Installation

1. **Clone the repository:**
```bash
git clone https://github.com/ayushsharma-1/LogAid.git
cd LogAid
```

2. **Set up your environment:**
```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your API keys
nano .env
```

3. **Build LogAid:**
```bash
# Using Make (recommended)
make build

# Or with Go directly
go build -o logaid ./cmd
```

4. **Configure API Keys:**
```bash
# Option 1: Edit .env file
GEMINI_API_KEY=your_gemini_api_key_here
OPENAI_API_KEY=your_openai_api_key_here

# Option 2: Export environment variables
export GEMINI_API_KEY="your_gemini_api_key_here"
export OPENAI_API_KEY="your_openai_api_key_here"
```

5. **Run LogAid:**
```bash
# Using the runner script (loads .env automatically)
./run_logaid.sh version

# Or run directly
./logaid version
```

### Getting API Keys

#### Google Gemini API
1. Visit [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Create a new API key
3. Copy the key to your `.env` file

#### OpenAI API
1. Visit [OpenAI API Keys](https://platform.openai.com/api-keys)
2. Create a new secret key
3. Copy the key to your `.env` file

## 📖 Usage

### Basic Commands

```bash
# Start monitoring your shell
./run_logaid.sh start

# Show current configuration
./run_logaid.sh config

# Display help information
./run_logaid.sh help

# Show version and active providers
./run_logaid.sh version

# View recent activity logs
./run_logaid.sh logs
```

### Interactive Mode

When you start LogAid with `./run_logaid.sh start`, it begins monitoring your commands:

```bash
$ ./run_logaid.sh start
LogAid is monitoring your commands...
Type 'exit' to quit LogAid

[LogAid] $ ls -la /nonexistent
ls: cannot access '/nonexistent': No such file or directory

╭─ LogAid Suggestion
│ ls -la ~/
│ Explanation: The directory /nonexistent doesn't exist. Try listing your home directory instead.
│ Plugin: generic
╰─
Execute suggestion? [y/N]: y
Executing: ls -la ~/
# ... output appears ...
```

## 🧪 Testing

LogAid includes a comprehensive test suite:

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run only integration tests
cd tests && go test -v .

# Run linting and formatting
make lint
```

## 🏗️ Development

### Building from Source

```bash
# Install dependencies
make deps

# Development build (fast)
make dev

# Production build with optimizations
make build

# Multi-platform builds
make build-all
```

See [DEVELOPMENT.md](DEVELOPMENT.md) for detailed development guidelines.

## 📦 Installation Methods

### Method 1: Build from Source (Recommended)
```bash
git clone https://github.com/ayushsharma-1/LogAid.git
cd LogAid
make build
```

### Method 2: Download Binary Release
```bash
# Download from GitHub Releases (when available)
curl -L https://github.com/ayushsharma-1/LogAid/releases/latest/download/logaid-linux-amd64.tar.gz | tar xz
```

## 🤝 Support

- **Issues**: [GitHub Issues](https://github.com/ayushsharma-1/LogAid/issues)
- **Documentation**: [DEVELOPMENT.md](DEVELOPMENT.md)
- **Changelog**: [CHANGELOG.md](CHANGELOG.md)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**LogAid - Your CLI Guardian** 🛡️

[GitHub](https://github.com/ayushsharma-1/LogAid) • [Documentation](DEVELOPMENT.md) • [Changelog](CHANGELOG.md)

</div>

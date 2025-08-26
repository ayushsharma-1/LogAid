# LogAid - Flow-State Agent

An intelligent CLI agent that detects command errors and provides real-time AI-generated solutions to maintain your development flow.

## 🚀 Features

- **Pseudo-terminal wrapper** for transparent operation
- **AI-powered error analysis** using Gemini or OpenAI
- **Non-intrusive, calm UX** design 
- **Privacy-first** with local fallback options
- **Cross-distribution** packaging
- **Real-time suggestions** without breaking flow state

## 📦 Installation

### Quick Start (AppImage)
```bash
# Download the latest release
wget https://github.com/ayushsharma-1/LogAid/releases/latest/download/logaid-linux-x86_64.AppImage

# Make executable
chmod +x logaid-linux-x86_64.AppImage

# Run
./logaid-linux-x86_64.AppImage
```

### Build from Source
```bash
# Clone the repository
git clone https://github.com/ayushsharma-1/LogAid.git
cd LogAid

# Build
go build -o logaid .

# Install (optional)
sudo mv logaid /usr/local/bin/
```

## ⚙️ Configuration

LogAid reads configuration from environment variables or a `.env` file:

```bash
# AI Configuration
LOGAID_AI_PROVIDER=gemini              # or "openai"
GEMINI_API_KEY=your-gemini-api-key     # Required for Gemini
OPENAI_API_KEY=your-openai-api-key     # Required for OpenAI
LOGAID_AI_MODEL=gemini-1.5-flash       # AI model to use
LOGAID_MAX_TOKENS=1000                 # Max response tokens
LOGAID_TEMPERATURE=0.3                 # AI creativity (0.0-1.0)

# Feature Flags
LOGAID_ENABLE_COLORS=true              # Enable colorized output
LOGAID_ENABLE_LOGGING=true             # Enable logging
LOGAID_ENABLE_LOCAL_FALLBACK=false     # Use local AI models

# Terminal Configuration
LOGAID_SHELL=/bin/bash                 # Shell to wrap
LOGAID_PROMPT_TIMEOUT=30               # AI request timeout (seconds)
```

## 🏃 Usage

### Start the Agent
```bash
# Start with default shell
logaid run

# Start with specific shell
logaid run --shell /bin/zsh
```

### Test AI Integration
```bash
# Verify your configuration and API keys
logaid test
```

### View Configuration
```bash
# Display current configuration
logaid config
```

### Help
```bash
# Show help
logaid --help
logaid run --help
```

## 🎯 How It Works

1. **Command Detection**: LogAid wraps your terminal using a pseudo-terminal (PTY)
2. **Error Monitoring**: Monitors command exit codes and stderr output
3. **AI Analysis**: Sends error context to AI for analysis when failures occur
4. **Smart Suggestions**: Displays non-intrusive, actionable suggestions
5. **Flow Preservation**: Maintains your development flow with calm UX

## 🧠 AI Providers

### Gemini (Google)
- Fast and cost-effective
- Large context windows
- Excellent for CLI error analysis

### OpenAI
- GPT-3.5 and GPT-4 support
- High-quality reasoning
- Good for complex error scenarios

## 🔒 Privacy & Security

- **No auto-execution**: All suggestions require explicit user approval
- **Minimal data collection**: Only command and error output when failures occur
- **Local fallback**: Option to use on-device models for complete privacy
- **No training data**: Your data is never used for model training (when configured)

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User Input    │────│  LogAid Agent   │────│   Shell/PTY     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │   AI Provider   │
                    │ (Gemini/OpenAI) │
                    └─────────────────┘
```

## 🛠️ Development

### Prerequisites
- Go 1.21+
- API key for Gemini or OpenAI

### Building
```bash
# Install dependencies
go mod tidy

# Build
go build -o logaid .

# Run tests
go test ./...
```

### Contributing
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📚 Commands Reference

| Command | Description |
|---------|-------------|
| `logaid run` | Start the terminal wrapper agent |
| `logaid test` | Test AI integration and configuration |
| `logaid config` | Display current configuration |
| `logaid version` | Show version information |
| `logaid --help` | Show help information |

## 🔧 Troubleshooting

### API Key Issues
```bash
# Test your configuration
logaid test

# Check configuration
logaid config
```

### Permission Issues
```bash
# Ensure executable permissions
chmod +x logaid

# Check shell access
echo $SHELL
```

### Shell Compatibility
LogAid supports:
- bash
- zsh
- fish
- sh

## 📄 License

MIT License - see LICENSE file for details.

## 🙏 Acknowledgments

- [creack/pty](https://github.com/creack/pty) - PTY interface for Go
- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [fatih/color](https://github.com/fatih/color) - Terminal colors

## 🌟 Support

- 📧 Email: support@logaid.dev
- 🐛 Issues: [GitHub Issues](https://github.com/ayushsharma-1/LogAid/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/ayushsharma-1/LogAid/discussions)

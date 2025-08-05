# LogAid 🚀

**AI-Powered Linux CLI Assistant**

LogAid is a CLI-first AI assistant that intercepts shell commands and error logs in real time, identifies mistakes (typos, wrong package names, syntax errors, etc.), and suggests or auto-applies corrections with user confirmation.

## Features

- 🔍 **Real-time Command Monitoring** - Intercepts every command and its output
- 🧠 **AI-Powered Error Detection** - Uses Gemini 2.5 Pro/Flash for intelligent suggestions
- 🔌 **Plugin Architecture** - Extensible with built-in plugins for apt, npm, git, docker, pip, systemctl
- 🎨 **Beautiful CLI UX** - Color-coded output with ASCII art
- 📝 **Command History** - Logs all commands, suggestions, and outcomes

## Quick Start

### Installation

```bash
# Via Go (recommended)
go install github.com/ayushsharma-1/LogAid@latest

# Or download binary from releases
curl -L https://github.com/ayushsharma-1/LogAid/releases/latest/download/logaid-linux-amd64 -o logaid
chmod +x logaid
sudo mv logaid /usr/local/bin/
```

### Usage

```bash
# Start LogAid shell
logaid

# Or wrap existing commands
logaid exec "sudo apt install rediscli"
```

### Configuration

Create `~/.logaid/.env`:

```env
AI_PROVIDER=gemini
GEMINI_API_KEY=your_key_here
LOG_LEVEL=info
```

## Plugin Development

LogAid uses a plugin architecture. Each plugin implements:

```go
type Plugin interface {
    Match(cmd string, output string) bool       // When to trigger
    Suggest(cmd string, output string) string   // AI-powered fix
    Name() string                              // Plugin identifier
}
```

See `plugins/` directory for examples.

## Testing

```bash
# Run all tests
go test ./...

# Run integration tests
go test -tags=integration ./tests/

# Run in Docker
docker build -t logaid-test .
docker run --rm logaid-test
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) file for details.

# LogAid Architecture & Technology Stack

## 📋 Table of Contents
- [Architecture Overview](#architecture-overview)
- [Technology Stack](#technology-stack)
- [Library Choices & Rationale](#library-choices--rationale)
- [Alternative Libraries Considered](#alternative-libraries-considered)
- [Plugin System](#plugin-system)
- [AI Integration](#ai-integration)
- [Performance Considerations](#performance-considerations)

## 🏗️ Architecture Overview

LogAid follows a modular, plugin-based architecture designed for extensibility and performance:

```
┌─────────────────────────────────────────────────────────────┐
│                        CLI Interface                        │
│                      (Cobra Commands)                       │
├─────────────────────────────────────────────────────────────┤
│                     Command Execution                      │
│                    (Engine + Monitoring)                   │
├─────────────────────┬───────────────────┬───────────────────┤
│    Plugin System    │   Error Detection │   AI Integration  │
│   (Rule-based)      │     (Patterns)    │  (Fallback/Enh.)  │
├─────────────────────┼───────────────────┼───────────────────┤
│  APT │ Git │ Docker │   NPM │ Pip │ ... │  Gemini │ OpenAI  │
└─────────────────────┴───────────────────┴───────────────────┘
                               │
                    ┌──────────┴──────────┐
                    │   Configuration     │
                    │   (.env + Viper)    │
                    └─────────────────────┘
```

### Core Components

1. **CLI Layer**: Built with Cobra for command parsing and routing
2. **Engine**: Core execution and monitoring logic
3. **Plugin System**: Modular error detection and correction
4. **AI Integration**: Fallback for complex scenarios
5. **Configuration**: Centralized settings management

## 🛠️ Technology Stack

### Core Dependencies

| Library | Version | Purpose | License |
|---------|---------|---------|---------|
| **Go** | 1.21+ | Core language | BSD-3-Clause |
| **Cobra** | v1.8.0 | CLI framework | Apache-2.0 |
| **Viper** | v1.18.2 | Configuration management | MIT |
| **Color** | v1.16.0 | Terminal colors | MIT |
| **godotenv** | v1.5.1 | Environment variables | MIT |

### Development Dependencies

| Library | Purpose | Alternative Considered |
|---------|---------|----------------------|
| Standard `testing` | Unit tests | Testify (unnecessary complexity) |
| Standard `net/http` | HTTP client | Resty (overkill for simple requests) |
| Standard `encoding/json` | JSON parsing | GJSON (performance not critical) |

## 📚 Library Choices & Rationale

### 1. CLI Framework: Cobra vs Alternatives

**Chosen: `spf13/cobra`**
- ✅ Industry standard for Go CLI applications
- ✅ Excellent subcommand support
- ✅ Auto-generated help and completions
- ✅ Used by kubectl, Hugo, GitHub CLI

**Alternatives Considered:**
- `urfave/cli`: Simpler but less powerful
- `kingpin`: Less maintained, complex API
- `flag` (stdlib): Too basic for our needs

### 2. Configuration: Viper vs Alternatives

**Chosen: `spf13/viper`**
- ✅ Multiple format support (JSON, YAML, ENV)
- ✅ Live reloading capabilities
- ✅ Nested configuration handling
- ✅ Pairs perfectly with Cobra

**Alternatives Considered:**
- `kelseyhightower/envconfig`: ENV-only, limited
- `caarlos0/env`: Simple but lacks features
- Custom implementation: Reinventing the wheel

### 3. Colors: fatih/color vs Alternatives

**Chosen: `fatih/color`**
- ✅ Simple, intuitive API
- ✅ Cross-platform compatibility
- ✅ Minimal dependencies
- ✅ Wide adoption in Go ecosystem

**Alternatives Considered:**
- `logrus` with color: Overkill for simple coloring
- `aurora`: More features but unnecessary complexity
- ANSI codes directly: Platform compatibility issues

### 4. HTTP Client: Standard Library vs Alternatives

**Chosen: Standard `net/http`**
- ✅ No external dependencies
- ✅ Sufficient for our API calls
- ✅ Better control over requests
- ✅ Reduces binary size

**Alternatives Considered:**
- `go-resty/resty`: Great API but adds dependency
- `hashicorp/go-retryablehttp`: Useful but complex
- `parnurzeal/gorequest`: Deprecated

### 5. AI Integration: Direct HTTP vs SDKs

**Chosen: Direct HTTP calls**
- ✅ Full control over requests/responses
- ✅ No vendor lock-in with SDKs
- ✅ Smaller binary size
- ✅ Easier to support multiple providers

**Alternatives Considered:**
- Google AI SDK: Provider lock-in, large dependency
- OpenAI SDK: Same issues as Google
- LangChain Go: Overkill for simple API calls

## 🔌 Plugin System

### Design Principles

1. **Interface-based**: All plugins implement the same interface
2. **Performance-first**: Fast pattern matching before AI calls
3. **Extensible**: Easy to add new tools/commands
4. **Isolated**: Plugins don't interfere with each other

### Plugin Interface

```go
type Plugin interface {
    Match(cmd string, output string) bool    // Pattern matching
    Suggest(cmd string, output string) string // Generate suggestion
    Name() string                           // Plugin identifier
}
```

### Current Plugins

| Plugin | Commands Covered | Pattern Types |
|--------|------------------|---------------|
| **APT** | apt, apt-get | Package names, permissions, locks |
| **Git** | git * | Command typos, branch issues |
| **Docker** | docker * | Image names, daemon issues |
| **NPM** | npm * | Package names, command typos |
| **Pip** | pip * | Package names, permissions |
| **Systemctl** | systemctl * | Service names, permissions |

### Performance Characteristics

- **Pattern Matching**: 25-85ns per plugin
- **Memory Usage**: <1MB per plugin
- **Concurrent Safe**: All plugins are stateless

## 🧠 AI Integration

### Multi-Provider Support

LogAid supports multiple AI providers for maximum flexibility:

1. **Google Gemini**: Primary choice
   - Fast response times
   - Cost-effective
   - Good command understanding

2. **OpenAI GPT**: Alternative option
   - High accuracy
   - More expensive
   - Well-documented API

### AI Usage Strategy

```
Command Error
     ↓
Plugin Matching (Fast)
     ↓
Found? → Return Rule-based Suggestion
     ↓
Not Found? → AI Fallback
     ↓
Generate Context-aware Suggestion
```

### Prompt Engineering

LogAid uses carefully crafted prompts for each tool:

- **Context**: Command + error output
- **Instructions**: Tool-specific guidance
- **Format**: Request specific output format
- **Temperature**: Low (0.1) for consistent suggestions

## ⚡ Performance Considerations

### Design Decisions for Speed

1. **Plugin First**: Rule-based matching before AI
2. **Concurrent Plugins**: Parallel pattern matching
3. **Caching**: Optional suggestion caching
4. **Minimal Dependencies**: Faster startup
5. **Compiled Binary**: No runtime interpretation

### Benchmarks

- **Plugin Matching**: 25-85ns average
- **Binary Size**: ~8MB (optimized build)
- **Startup Time**: <100ms
- **Memory Usage**: <50MB typical

### Optimization Techniques

1. **String Interning**: Reduced memory for common patterns
2. **Lazy Loading**: Plugins loaded on demand
3. **Batch Processing**: Multiple error patterns in one check
4. **Response Streaming**: AI responses processed incrementally

## 🔮 Future Architecture Considerations

### Planned Improvements

1. **Plugin SDK**: External plugin development
2. **Distributed Caching**: Shared suggestion cache
3. **Machine Learning**: Local model integration
4. **Telemetry**: Usage analytics and improvement
5. **Shell Integration**: Deeper shell hooks

### Scalability

- **Horizontal**: Multiple LogAid instances
- **Vertical**: Increased plugin capacity
- **Cloud**: SaaS suggestion service
- **Offline**: Local AI model support

---

This architecture ensures LogAid remains fast, extensible, and maintainable while providing intelligent command suggestions across multiple tools and platforms.

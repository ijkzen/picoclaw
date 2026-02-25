# PicoClaw Knowledge Base

**Generated:** 2026-02-26
**Commit:** 094d6591
**Branch:** main

## OVERVIEW

PicoClaw is an ultra-lightweight personal AI Assistant in Go. Runs on $10 hardware with <10MB RAM. Cobra-based CLI with modular architecture supporting multiple LLM providers and chat channels.

## STRUCTURE

```
.
├── cmd/picoclaw/          # CLI entry point and commands
│   ├── main.go           # Cobra root command setup
│   └── internal/         # Subcommand implementations
├── pkg/                   # Core library packages
│   ├── agent/            # Agent orchestration logic
│   ├── auth/             # Authentication utilities
│   ├── channels/         # Chat platform integrations (20+ channels)
│   ├── config/           # Configuration management
│   ├── cron/             # Scheduled task execution
│   ├── heartbeat/        # Periodic task processing
│   ├── migrate/          # Config/workspace migrations
│   ├── providers/        # LLM provider implementations
│   ├── routing/          # Request routing & session management
│   ├── skills/           # Built-in skill system
│   ├── tools/            # Tool framework (i2c, spi, file ops)
│   └── utils/            # Shared utilities
├── config/               # Default configuration templates
├── docs/                 # Documentation
├── skills/               # Built-in skill definitions
└── workspace/            # Default workspace templates
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add new CLI command | `cmd/picoclaw/internal/` | Follow Cobra pattern, add to main.go |
| Add LLM provider | `pkg/providers/` | Implement Provider interface |
| Add chat channel | `pkg/channels/` | Implement Channel interface |
| Add tool capability | `pkg/tools/` | Register in tool registry |
| Configuration schema | `pkg/config/` | JSON tags define config format |
| Agent core logic | `pkg/agent/` | Tool loop, message handling |
| Auth flows | `cmd/picoclaw/internal/auth/` | OAuth, API key management |

## CODE MAP

### Key Interfaces

| Interface | Package | Purpose |
|-----------|---------|---------|
| `Provider` | `pkg/providers/` | LLM provider abstraction |
| `Channel` | `pkg/channels/` | Chat platform abstraction |
| `Tool` | `pkg/tools/` | Agent tool capability |
| `MessageHandler` | `pkg/agent/` | Message processing |

### Entry Points

| Entry | File | Role |
|-------|------|------|
| CLI root | `cmd/picoclaw/main.go` | Cobra command setup |
| Agent cmd | `cmd/picoclaw/internal/agent/command.go` | Interactive/chat mode |
| Gateway | `cmd/picoclaw/internal/gateway/command.go` | HTTP server mode |

## CONVENTIONS

### Go Code Style

- **Linting**: golangci-lint with 50+ linters (see `.golangci.yaml`)
- **Line length**: 120 characters max
- **Imports**: gci formatter with standard/default/localmodule sections
- **Complexity**: gocyclo max 20, gocognit max 25
- **Function length**: max 120 lines, 40 statements

### Project-Specific

- **CGO disabled**: `CGO_ENABLED=0` for all builds
- **Version injection**: Use ldflags with `internal.version`, `internal.gitCommit`
- **Workspace path**: `~/.picoclaw/workspace` default
- **Config path**: `~/.picoclaw/config.json`

### File Organization

```
cmd/picoclaw/internal/{command}/
├── command.go      # Cobra command implementation
├── helpers.go      # Command-specific helpers
└── *_test.go       # Tests

pkg/{domain}/
├── {type}.go       # Core types and interfaces
├── {type}_impl.go  # Implementation
└── *_test.go       # Tests
```

## ANTI-PATTERNS

- **Don't** use `interface{}` - use `any` (gofmt rewrite rule enforces)
- **Don't** use naked returns (max 3 func lines)
- **Don't** ignore errors - always handle or explicitly ignore with comment
- **Don't** use magic numbers without context (mnd linter)
- **Don't** create package-level variables (gochecknoglobals disabled but discouraged)

## COMMANDS

```bash
# Development
make build              # Build for current platform
make build-all          # Build all platforms (linux/darwin/windows)
make install            # Install to ~/.local/bin
make test               # Run tests
make lint               # Run golangci-lint
make fmt                # Format code

# Running
picoclaw onboard        # Initialize config & workspace
picoclaw agent -m "..." # One-shot query
picoclaw agent          # Interactive mode
picoclaw gateway        # Start HTTP gateway
picoclaw cron list      # List scheduled jobs

# Docker
docker compose --profile gateway up -d
docker compose run --rm picoclaw-agent -m "..."
```

## NOTES

- **Multi-arch**: Supports x86_64, ARM64, RISC-V, LoongArch
- **Sandbox**: Tools respect `restrict_to_workspace` config
- **Heartbeat**: HEARTBEAT.md in workspace defines periodic tasks
- **Skills**: Extensible skill system in `skills/` and workspace
- **Model routing**: Uses `vendor/model` format (e.g., `openai/gpt-4`)

## EXTERNAL DEPENDENCIES

### Key Libraries

| Library | Purpose |
|---------|---------|
| cobra | CLI framework |
| anthropic-sdk-go | Claude API |
| openai-go | OpenAI API |
| discordgo | Discord bot |
| telego | Telegram bot |
| gorilla/websocket | WebSocket support |

### LLM Providers Supported

- OpenAI, Anthropic, Gemini, Zhipu, DeepSeek, Groq, Cerebras
- OpenRouter (multi-provider), Ollama (local), vLLM

### Chat Channels

Telegram, Discord, Slack, LINE, DingTalk, QQ, WeCom, Feishu, WhatsApp, OneBot

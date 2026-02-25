# cmd/picoclaw/internal Knowledge Base

**Path:** `cmd/picoclaw/internal/`
**Files:** ~30
**Subcommands:** 13

## OVERVIEW

CLI subcommand implementations using Cobra. Each subdirectory is a self-contained command with its own flags, args, and logic.

## STRUCTURE

```
cmd/picoclaw/internal/
├── helpers.go           # Shared command utilities
├── agent/               # Interactive/chat mode
│   ├── command.go       # Cobra command
│   └── helpers.go       # Agent-specific helpers
├── auth/                # Authentication management
│   ├── command.go       # Auth subcommands
│   ├── login.go         # Provider login
│   ├── logout.go        # Logout/clear
│   └── status.go        # Auth status
├── gateway/             # HTTP server mode
│   └── command.go
├── cron/                # Scheduled tasks
│   └── command.go
├── onboard/             # Initial setup
│   └── command.go
├── status/              # System status
│   └── command.go
├── skills/              # Skill management
│   └── command.go
├── migrate/             # Config migrations
│   └── command.go
└── version/             # Version info
    └── command.go
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add CLI command | Create `{cmd}/command.go` | Follow Cobra pattern |
| Command helpers | `{cmd}/helpers.go` | Command-specific utils |
| Shared utilities | `helpers.go` | Cross-command helpers |
| Agent mode | `agent/` | Interactive chat |
| Gateway mode | `gateway/` | HTTP server |
| Auth flows | `auth/` | OAuth, API keys |
| Scheduled jobs | `cron/` | Cron management |

## CONVENTIONS

### Command Pattern

```go
func NewCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "name",
        Short: "Description",
        RunE:  runCommand,
    }
    cmd.Flags().StringP("flag", "f", "", "Usage")
    return cmd
}
```

### Registration

All commands registered in `../main.go`:

```go
cmd.AddCommand(
    onboard.NewOnboardCommand(),
    agent.NewAgentCommand(),
    // ...
)
```

### Naming

- Directory: lowercase command name
- Main file: `command.go` with `New{X}Command()`
- Helpers: `helpers.go`
- Tests: `*_test.go`

## ANTI-PATTERNS

- **Don't** put business logic in commands - delegate to pkg/
- **Don't** ignore flag validation errors
- **Don't** print directly - use configured output writer
- **Don't** hardcode paths - respect env vars/config

## NOTES

- All commands use RunE (return errors, don't exit)
- Global flags defined in root command
- Configuration loaded early in command chain
- Tests use command helpers pattern

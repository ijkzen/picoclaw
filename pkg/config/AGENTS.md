# pkg/config Knowledge Base

**Path:** `pkg/config/`
**Files:** 6
**Purpose:** Configuration management and migrations

## OVERVIEW

Central configuration management with JSON config file, environment variable support, and migration system for config updates.

## STRUCTURE

```
pkg/config/
├── config.go          # Core Config struct and loading
├── defaults.go        # Default configuration values
├── migration.go       # Config version migrations
└── model_config.go    # Model list configuration types
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Config schema | `config.go` | Config struct with JSON tags |
| Default values | `defaults.go` | Agent defaults, tool defaults |
| Config migration | `migration.go` | Version upgrade logic |
| Model config | `model_config.go` | model_list types |
| Load config | `config.go` | Load() function |

## CONVENTIONS

### Config Struct Pattern

```go
type Config struct {
    Agents    AgentsConfig    `json:"agents"`
    ModelList []ModelConfig   `json:"model_list"`
    Channels  ChannelsConfig  `json:"channels"`
    Tools     ToolsConfig     `json:"tools"`
    // ...
}
```

### Flexible Types

```go
// FlexibleStringSlice accepts JSON strings or numbers
type FlexibleStringSlice []string

// Usage in config:
// "allow_from": ["123", 456]  // Both work
```

### Environment Override

Environment variables override config file values:
- `PICOCLAW_AGENTS_DEFAULTS_MODEL`
- `PICOCLAW_HEARTBEAT_INTERVAL`

## ANTI-PATTERNS

- **Don't** break backward compatibility without migration
- **Don't** use `interface{}` for config fields (use specific types)
- **Don't** ignore UnmarshalJSON errors
- **Don't** mutate config after load (use copy)

## NOTES

- Config path: `~/.picoclaw/config.json`
- Supports hot-reload for some fields
- Migrations run automatically on config load
- Round-robin load balancing via atomic counter

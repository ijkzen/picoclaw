# pkg/tools Knowledge Base

**Path:** `pkg/tools/`
**Files:** 35
**Subdirectories:** 37 (one per tool category)

## OVERVIEW

Tool framework providing agent capabilities. Each subdirectory implements a specific tool category (file ops, hardware interfaces, system commands).

## STRUCTURE

```
pkg/tools/
├── registry.go           # Tool registration and discovery
├── execute.go            # Command execution utilities
├── i2c.go               # I2C hardware interface
├── spi.go               # SPI hardware interface
├── gpio/                # GPIO pin control
├── http/                # HTTP requests
├── file/                # File operations
├── exec/                # Shell execution
├── python/              # Python script execution
├── search/              # Web search (Brave, DuckDuckGo, Tavily)
├── browser/             # Browser automation (Playwright)
├── image/               # Image processing
├── vision/              # Vision/ML models
├── tmux/                # Terminal multiplexer control
├── git/                 # Git operations
└── ... (additional tools)
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add new tool | Create subdirectory | Implement Tool interface |
| Tool registration | `registry.go` | Auto-registration via init() |
| Command execution | `execute.go` | Sandbox-aware exec |
| Hardware tools | `i2c.go`, `spi.go`, `gpio/` | Linux-specific |
| Web search | `search/` | Multiple provider support |
| Browser automation | `browser/` | Playwright integration |

## CONVENTIONS

### Tool Interface

```go
type Tool interface {
    Name() string
    Description() string
    Schema() *ToolSchema
    Execute(ctx context.Context, input json.RawMessage) (any, error)
}
```

### Registration Pattern

```go
func init() {
    Register(NewMyTool())
}
```

### Naming

- Tool files: `{tool}.go` for single file, `{tool}/` for complex tools
- Test files: `{tool}_test.go` or `{tool}/{tool}_test.go`

## ANTI-PATTERNS

- **Don't** bypass workspace sandbox in file operations
- **Don't** execute commands without timeout context
- **Don't** register tools without schema validation
- **Don't** use blocking operations without context cancellation

## NOTES

- Tools auto-discovered via registry
- Sandbox restricts file/exec access to workspace when enabled
- Hardware tools (i2c/spi) require appropriate permissions
- Browser tool requires Playwright installation

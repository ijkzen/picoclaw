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

- **CGO disabled**: `CGO_ENABLED=0` for all builds
- **Version injection**: Use ldflags with `internal.version`, `internal.gitCommit`
- **Workspace path**: `~/.picoclaw/workspace` default
- **Config path**: `~/.picoclaw/config.json`
- **Gateway port**: `18790` (default), Web UI runs on port `18791` (gateway port + 1)


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


## FRONTEND DEVELOPMENT WORKFLOW

This section describes the complete workflow for developing and verifying the PicoClaw Web UI using Playwright MCP.

### Quick Reference

```bash
# 1. Make changes to frontend files in web/src/
# 2. Build frontend
cd /Users/ijkzen/Projects/GO-Project/picoclaw/web && pnpm run build

# 3. Copy build artifacts
rm -rf /Users/ijkzen/Projects/GO-Project/picoclaw/pkg/web/dist/*
cp -r /Users/ijkzen/Projects/GO-Project/picoclaw/web/dist/web/* /Users/ijkzen/Projects/GO-Project/picoclaw/pkg/web/dist/

# 4. Build and install backend
cd /Users/ijkzen/Projects/GO-Project/picoclaw && make install

# 5. Restart gateway
picoclaw gateway stop && sleep 1 && picoclaw gateway start

# 6. Verify with Playwright MCP (see examples below)
```

### Frontend Project Structure

```
web/
├── src/
│   ├── app/
│   │   ├── components/     # Shared components
│   │   │   └── layout/     # Sidebar layout component
│   │   ├── pages/          # Page components
│   │   │   ├── chat/       # Chat interface
│   │   │   └── settings/   # Settings interface
│   │   ├── services/       # API services
│   │   └── models/         # TypeScript interfaces
│   ├── styles.scss         # Global styles + Angular Material theme
│   └── index.html
├── angular.json
└── package.json
```

### Component Template & Style Policy

- For all Angular components except the root component (`web/src/app/app.ts`):
  - Do not use inline `template` or inline `styles` in `@Component`.
  - Use `templateUrl` with a separate `.component.html` file.
  - Do not create component-level `.scss`/`.css` style files (`styleUrl`/`styleUrls` should not be used).
  - Keep component styling in template markup using Tailwind utility classes.
  - If host-level layout styling is required, use `host: { class: '...' }` in the component decorator.

### Component Complexity & Decomposition Policy

- Frontend UI components must stay small and single-responsibility.
- If a component becomes too complex (for example: multiple large sections/tabs, mixed responsibilities, or hard-to-maintain template/logic), split it into smaller standalone child components.
- Prefer decomposition by feature area (for example: header, tab panel, list item editor, overlay dialog).
- Parent components should focus on orchestration/state flow; child components should focus on presentation and local interactions.

### Common Modifications

#### Global Flat Card Style

Edit `web/src/app/app.config.ts`:

```typescript
import { MAT_CARD_CONFIG } from '@angular/material/card';

export const appConfig: ApplicationConfig = {
  providers: [
    // ... other providers
    {
      provide: MAT_CARD_CONFIG,
      useValue: {
        appearance: 'outlined'  // Flat style for all cards
      }
    }
  ]
};
```

#### Fixed Input at Bottom of Chat

Edit `web/src/app/pages/chat/chat.component.ts`:

```typescript
// Template structure
<div class="chat-wrapper">
  <div class="messages-area">     <!-- flex: 1, scrollable -->
    <!-- Messages list -->
  </div>
  <div class="input-area">        <!-- Fixed at bottom -->
    <!-- Input field and send button -->
  </div>
</div>

// Key CSS styles
.chat-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.messages-area {
  flex: 1;
  overflow-y: auto;
}

.input-area {
  flex-shrink: 0;
  background: var(--mat-sys-surface);
}
```

#### Angular Material Theme Setup

Edit `web/src/styles.scss`:

```scss
// Angular Material Theming - @use must come first
@use '@angular/material' as mat;

// Import Tailwind CSS (optional, can coexist)
@import "tailwindcss";

// Define theme
$web-theme: (
  color: (
    theme-type: light,
    primary: mat.$azure-palette,
    tertiary: mat.$blue-palette,
  ),
  typography: Roboto,
  density: 0,
);

@include mat.core();
@include mat.theme($web-theme);

// Dark theme support
.dark {
  @include mat.theme((
    color: (
      theme-type: dark,
      primary: mat.$azure-palette,
      tertiary: mat.$blue-palette,
    ),
    typography: Roboto,
    density: 0,
  ));
}
```

### Playwright MCP Verification

#### Basic Navigation and Screenshot

```typescript
// Navigate to web UI
skill_mcp({
  mcp_name: "playwright",
  tool_name: "browser_navigate",
  arguments: { url: "http://127.0.0.1:18791/" }
});

// Set viewport size for consistent screenshots
skill_mcp({
  mcp_name: "playwright",
  tool_name: "browser_resize",
  arguments: { width: 1280, height: 800 }
});

// Take screenshot
skill_mcp({
  mcp_name: "playwright",
  tool_name: "browser_take_screenshot",
  arguments: {
    filename: "verification.png",
    type: "png"
  }
});

// Cleanup
skill_mcp({
  mcp_name: "playwright",
  tool_name: "browser_close"
});
```

#### Interactive Testing Example

```typescript
// Navigate to settings page
skill_mcp({
  mcp_name: "playwright",
  tool_name: "browser_navigate",
  arguments: { url: "http://127.0.0.1:18791/settings" }
});

// Click on Models tab
skill_mcp({
  mcp_name: "playwright",
  tool_name: "browser_click",
  arguments: {
    element: "Models tab",
    ref: "e89"  // Reference from page snapshot
  }
});

// Fill form field
skill_mcp({
  mcp_name: "playwright",
  tool_name: "browser_type",
  arguments: {
    ref: "e115",
    text: "gpt-4"
  }
});

// Take full page screenshot
skill_mcp({
  mcp_name: "playwright",
  tool_name: "browser_take_screenshot",
  arguments: {
    filename: "settings_full.png",
    type: "png",
    fullPage: true
  }
});
```

### Verification Checklist

Before submitting changes:

- [ ] Frontend builds without errors (`pnpm run build`)
- [ ] Backend builds successfully (`make install`)
- [ ] Gateway starts and web UI is accessible
- [ ] Playwright MCP can navigate to page
- [ ] Screenshot shows expected layout
- [ ] Input components are properly positioned
- [ ] No console errors in browser
- [ ] Responsive design works at 1280x800 and mobile sizes

### Troubleshooting

#### Build Errors

**`@use rules must be written before any other rules`**
- Ensure `@use '@angular/material' as mat;` comes before `@import` in styles.scss

**`NG5002: @else block must be last`**
- Check template for duplicate `@else` blocks in @if statements

**Module not found errors**
- Run `pnpm install` to ensure all dependencies are installed

#### Layout Issues

**Input not sticking to bottom:**
- Verify parent has `display: flex; flex-direction: column; height: 100%;`
- Messages container needs `flex: 1; overflow-y: auto;`
- Input container needs `flex-shrink: 0;`

**Content not scrolling:**
- Ensure `overflow-y: auto` is set on scrollable containers
- Check that parent containers have defined heights

#### Playwright MCP Issues

**Cannot connect to gateway:**
- Verify gateway status: `picoclaw gateway status`
- Check if port 18791 is accessible
- Ensure no firewall blocking local connections

**Page elements not found:**
- Use `browser_snapshot` to get current page structure
- Element references change between page loads
- Always get fresh references before interactions

### Complete Example: Layout Modification

Here's a complete workflow example for modifying the chat page layout:

```bash
# Step 1: Edit chat.component.ts
# - Remove welcome card
# - Fix input at bottom
# - Update styles

# Step 2: Build
$ cd /Users/ijkzen/Projects/GO-Project/picoclaw/web
$ pnpm run build
✔ Building...
Output location: /Users/ijkzen/Projects/GO-Project/picoclaw/web/dist/web

# Step 3: Copy to pkg
$ rm -rf /Users/ijkzen/Projects/GO-Project/picoclaw/pkg/web/dist/*
$ cp -r /Users/ijkzen/Projects/GO-Project/picoclaw/web/dist/web/* \
    /Users/ijkzen/Projects/GO-Project/picoclaw/pkg/web/dist/

# Step 4: Build backend
$ cd /Users/ijkzen/Projects/GO-Project/picoclaw
$ make install
Building picoclaw for darwin/arm64...
Build complete: build/picoclaw-darwin-arm64
Installation complete!

# Step 5: Restart gateway
$ picoclaw gateway stop && sleep 1 && picoclaw gateway start
✓ Gateway stopped
✓ Gateway started in background (PID: xxxxx)

# Step 6: Verify via Playwright MCP
# (Execute MCP commands as shown in examples above)
```

### References

- [Angular Material Components](https://material.angular.dev/components/categories)
- [Angular Material Card API](https://material.angular.dev/components/card/api)
- [Playwright MCP Documentation](https://github.com/microsoft/playwright-mcp)
- [Angular Flex Layout Guide](https://material.angular.dev/guides)

# pkg/agent Knowledge Base

**Path:** `pkg/agent/`
**Files:** 11

## OVERVIEW

Core agent orchestration. Manages the tool loop, message processing, and LLM interaction. The "brain" of PicoClaw.

## STRUCTURE

```
pkg/agent/
├── agent.go             # Core agent struct and lifecycle
├── loop.go              # Main tool iteration loop
├── executor.go          # Tool execution orchestration
├── context.go           # Agent context management
├── message.go           # Message handling utilities
├── history.go           # Conversation history management
├── system_prompt.go     # System prompt construction
├── config.go            # Agent configuration types
└── ... (additional files)
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Agent lifecycle | `agent.go` | Start, stop, reset |
| Tool loop | `loop.go` | Main iteration logic |
| Tool execution | `executor.go` | Calls tool registry |
| Context management | `context.go` | Session, state |
| System prompts | `system_prompt.go` | Prompt construction |
| Message history | `history.go` | Conversation state |

## CONVENTIONS

### Agent Loop Flow

```
1. Receive message
2. Build system prompt + history
3. Call LLM
4. If tool calls requested:
   a. Execute tools
   b. Append results
   c. Go to 3 (continue loop)
5. Return final response
```

### Configuration

```go
type Config struct {
    Model            string
    MaxTokens        int
    Temperature      float64
    MaxIterations    int
    Workspace        string
    RestrictToWorkspace bool
}
```

## ANTI-PATTERNS

- **Don't** exceed max_iterations - always bound the loop
- **Don't** lose error context when tools fail
- **Don't** expose system internals in user-facing errors
- **Don't** allow infinite recursion in tool calls

## NOTES

- Supports both streaming and non-streaming responses
- Tool results feed back into LLM context
- Heartbeat creates separate agent instances
- Subagents spawned for async tasks
- Context cancellation stops execution cleanly

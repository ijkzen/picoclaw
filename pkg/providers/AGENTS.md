# pkg/providers Knowledge Base

**Path:** `pkg/providers/`
**Files:** 32
**Subdirectories:** 37 (one per provider)

## OVERVIEW

LLM provider implementations supporting multiple protocols: OpenAI-compatible, Anthropic, and custom adapters. Zero-code provider addition via `vendor/model` format.

## STRUCTURE

```
pkg/providers/
├── provider.go          # Provider interface definition
├── registry.go          # Provider registration
├── common/              # Shared provider utilities
├── openai_compat/       # OpenAI-compatible base implementation
│   └── client.go
├── anthropic/           # Claude API implementation
├── openai/              # OpenAI direct implementation
├── gemini/              # Google Gemini implementation
├── zhipu/               # 智谱 AI (GLM) implementation
├── deepseek/            # DeepSeek implementation
├── groq/                # Groq (fast inference)
├── cerebras/            # Cerebras implementation
├── qwen/                # 通义千问 implementation
├── moonshot/            # Moonshot implementation
├── openrouter/          # OpenRouter (multi-provider)
├── ollama/              # Local Ollama support
├── volcengine/          # 火山引擎 implementation
└── ... (additional providers)
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add new provider | Create subdirectory | Implement Provider interface |
| Provider interface | `provider.go` | Core abstraction |
| OpenAI-compatible | `openai_compat/` | Base for most providers |
| Registry | `registry.go` | Provider discovery |
| Common utilities | `common/` | Shared HTTP, auth, retry |

## CONVENTIONS

### Provider Interface

```go
type Provider interface {
    Name() string
    CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error)
    CreateChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (ChatCompletionStream, error)
}
```

### Configuration

Providers configured via `model_list` in config.json:

```json
{
  "model_name": "gpt-4",
  "model": "openai/gpt-4",
  "api_key": "sk-...",
  "api_base": "https://api.openai.com/v1"
}
```

### Naming

- Directory: lowercase vendor name
- Files: `{provider}.go`, `{provider}_test.go`

## ANTI-PATTERNS

- **Don't** hardcode API endpoints - allow `api_base` override
- **Don't** ignore rate limits - implement exponential backoff
- **Don't** leak API keys in errors
- **Don't** block streaming responses - use proper context handling

## NOTES

- Protocol families: OpenAI-compatible, Anthropic, custom
- Most providers use OpenAI-compatible protocol
- Model routing: `vendor/model` format (e.g., `zhipu/glm-4`)
- Load balancing: Multiple endpoints per model_name supported

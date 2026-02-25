# pkg/channels Knowledge Base

**Path:** `pkg/channels/`
**Files:** 20
**Complexity:** High (20+ channel implementations)

## OVERVIEW

Chat platform integrations. Each channel implements the Channel interface for bidirectional messaging with various platforms.

## STRUCTURE

```
pkg/channels/
├── base.go              # Base channel interface and utilities
├── manager.go           # Channel lifecycle management
├── telegram.go          # Telegram bot (telego)
├── telegram_commands.go # Telegram command handlers
├── discord.go           # Discord bot (discordgo)
├── slack.go             # Slack integration
├── line.go              # LINE messaging
├── dingtalk.go          # DingTalk
├── wecom.go             # WeCom bot
├── wecom_app.go         # WeCom app
├── qq.go                # QQ bot
├── feishu_32.go         # Feishu 32-bit
├── feishu_64.go         # Feishu 64-bit
├── whatsapp.go          # WhatsApp
├── onebot.go            # OneBot protocol
├── maixcam.go           # MaixCAM hardware
└── ... (additional channels)
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add new channel | Create `{channel}.go` | Implement Channel interface |
| Channel interface | `base.go` | Core abstraction |
| Channel manager | `manager.go` | Start/stop/monitor |
| Telegram | `telegram.go`, `telegram_commands.go` | Most popular channel |
| Discord | `discord.go` | Rich embeds supported |
| Webhook handlers | `*_app.go` files | HTTP webhook receivers |

## CONVENTIONS

### Channel Interface

```go
type Channel interface {
    Name() string
    Start(ctx context.Context) error
    Stop() error
    SendMessage(ctx context.Context, msg Message) error
}
```

### Configuration

Channels configured in config.json:

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "...",
      "allow_from": ["user_id"]
    }
  }
}
```

### Naming

- Single file per channel: `{platform}.go`
- Complex channels: `{platform}_commands.go` for command handling
- App variants: `{platform}_app.go` for webhook-based

## ANTI-PATTERNS

- **Don't** expose tokens in logs
- **Don't** ignore message delivery errors
- **Don't** block on channel send - use buffered channels
- **Don't** forget context cancellation on Stop()

## NOTES

- Most channels use polling (Telegram, Discord)
- Some use webhooks (LINE, WeCom)
- `allow_from` restricts to specific users
- Gateway mode starts all enabled channels
- Heartbeat tasks can send proactive messages

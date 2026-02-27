# pkg/routing Knowledge Base

**Path:** `pkg/routing/`
**Files:** 6
**Purpose:** Agent routing and session key resolution

## OVERVIEW

Routes incoming messages to the appropriate agent based on a 7-level priority cascade. Manages session key generation for conversation continuity.

## STRUCTURE

```
pkg/routing/
├── route.go           # Route resolution logic
├── agent_id.go        # Agent ID generation
└── session_key.go     # Session key construction
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Route resolution | `route.go` | ResolveRoute() - 7-level cascade |
| Agent ID logic | `agent_id.go` | Generate/resolve agent IDs |
| Session keys | `session_key.go` | Build session identifiers |
| Routing config | `route.go` | RouteResolver struct |

## CONVENTIONS

### 7-Level Priority Cascade

```
1. binding.peer           (specific user peer)
2. binding.peer.parent    (parent/reply context)
3. binding.guild          (Discord guild/server)
4. binding.team           (team/org scope)
5. binding.account        (account-level)
6. binding.channel.*      (wildcard channel)
7. default                (fallback)
```

### Route Resolution

```go
resolver := NewRouteResolver(cfg)
result := resolver.ResolveRoute(RouteInput{
    Channel:   "discord",
    AccountID: "123456",
    Peer:      &RoutePeer{ID: "...", Username: "..."},
})
// result.AgentID, result.SessionKey, result.MatchedBy
```

### DM Scope

```go
type DMScope string
const (
    DMScopeMain     DMScope = "main"      // One session per user
    DMScopeChannel  DMScope = "channel"   // Session per channel
)
```

## ANTI-PATTERNS

- **Don't** change priority order without careful consideration
- **Don't** use non-normalized account IDs (always NormalizeAccountID)
- **Don't** ignore MatchedBy field (useful for debugging)
- **Don't** create circular routing rules

## NOTES

- Session keys determine conversation continuity
- Route resolution is synchronous and cached
- Peer bindings override guild bindings
- DMScope affects how direct messages are routed

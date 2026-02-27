# pkg/auth Knowledge Base

**Path:** `pkg/auth/`
**Files:** 7
**Purpose:** OAuth authentication and token management

## OVERVIEW

Authentication utilities for OAuth flows including PKCE, token storage, and provider-specific configurations (OpenAI, Google/Antigravity).

## STRUCTURE

```
pkg/auth/
├── oauth.go          # OAuth flow implementation
├── pkce.go           # PKCE (Proof Key for Code Exchange)
├── store.go          # Token storage (keyring/file fallback)
└── token.go          # Token types and utilities
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add OAuth provider | `oauth.go` | Define OAuthProviderConfig |
| PKCE flow | `pkce.go` | Code verifier/challenge generation |
| Token persistence | `store.go` | Keychain → file fallback |
| Provider config | `oauth.go` | OpenAIOAuthConfig(), GoogleAntigravityOAuthConfig() |

## CONVENTIONS

### OAuth Provider Config

```go
type OAuthProviderConfig struct {
    Issuer       string
    ClientID     string
    ClientSecret string  // For confidential clients (Google)
    TokenURL     string  // Override token endpoint
    Scopes       string
    Originator   string  // Client identifier
    Port         int     // Local callback port
}
```

### Token Storage Flow

1. Try system keyring (macOS Keychain, Linux secret-service, Windows Credential Manager)
2. Fall back to file-based storage in ~/.picoclaw/

## ANTI-PATTERNS

- **Don't** store tokens in plain text without encryption
- **Don't** hardcode client secrets (except for public clients)
- **Don't** ignore token refresh errors
- **Don't** use predictable PKCE verifiers

## NOTES

- OpenAI OAuth uses public client (no secret required)
- Google Antigravity uses confidential client (requires secret)
- PKCE required for all OAuth flows (security best practice)
- Token storage prefers OS keyring, falls back to file

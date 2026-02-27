# pkg/utils Knowledge Base

**Path:** `pkg/utils/`
**Files:** 8
**Purpose:** Shared utility functions

## OVERVIEW

Common utility functions used across the codebase: string manipulation, media handling, downloads, and compression.

## STRUCTURE

```
pkg/utils/
├── string.go          # String utilities (Truncate, DerefStr)
├── message.go         # Message formatting utilities
├── media.go           # Media type detection
├── download.go        # HTTP download utilities
├── zip.go             # ZIP archive handling
└── skills.go          # Skill-related utilities
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| String truncation | `string.go` | Truncate() with Unicode support |
| Pointer dereference | `string.go` | DerefStr(ptr, fallback) |
| Media type detection | `media.go` | MIME type helpers |
| HTTP download | `download.go` | Download with progress |
| ZIP extraction | `zip.go` | Extract ZIP archives |
| Message formatting | `message.go` | Channel message builders |

## CONVENTIONS

### String Truncation

```go
// Unicode-aware truncation with "..." suffix
short := utils.Truncate(longString, 100)
```

### Pointer Dereference

```go
// Safe dereference with fallback
value := utils.DerefStr(maybeNil, "default")
```

### Pure Functions

Utilities are stateless pure functions:
- No package-level variables
- No side effects
- Thread-safe

## ANTI-PATTERNS

- **Don't** add state to utils (keep functions pure)
- **Don't** duplicate utility functions elsewhere
- **Don't** use `interface{}` when specific type works
- **Don't** ignore Unicode (use []rune for string ops)

## NOTES

- Functions are intentionally simple and focused
- No external dependencies beyond standard library
- Test coverage required for all utilities
- Consider upstreaming reusable utilities

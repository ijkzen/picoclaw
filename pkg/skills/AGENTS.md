# pkg/skills Knowledge Base

**Path:** `pkg/skills/`
**Files:** 9

## OVERVIEW

Skill system for extending agent capabilities. Skills are reusable, composable units of functionality that can be dynamically loaded.

## STRUCTURE

```
pkg/skills/
├── skill.go             # Skill interface and registry
├── loader.go            # Skill loading from filesystem
├── builtin.go           # Built-in skill definitions
├── executor.go          # Skill execution runtime
├── parser.go            # Skill manifest parsing
└── ... (additional files)
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Skill interface | `skill.go` | Core abstraction |
| Load skills | `loader.go` | From workspace/skills/ |
| Built-in skills | `builtin.go` | Embedded defaults |
| Execute skill | `executor.go` | Runtime execution |
| Parse manifest | `parser.go` | YAML/JSON parsing |

## CONVENTIONS

### Skill Structure

```
workspace/skills/{skill-name}/
├── skill.yaml          # Manifest
├── handler.go          # Implementation
└── README.md           # Documentation
```

### Manifest Format

```yaml
name: skill-name
description: What it does
tools:
  - tool1
  - tool2
entry: handler.go
```

## ANTI-PATTERNS

- **Don't** load skills outside workspace
- **Don't** execute skills without validation
- **Don't** ignore manifest parse errors

## NOTES

- Skills in workspace override built-ins
- Hot-reload on file change (configurable)
- Skills can depend on other skills

# Data Model: Basic Subcommands

**Feature**: 004-basic-commands  
**Date**: 2024-12-14

## Overview

This feature adds CLI subcommands and does not introduce new persistent data models. The only data artifact is the configuration file template used by the `init` command.

## Entities

### ConfigTemplate

The `init` command generates a `.bumpkin.yaml` configuration file. This uses the existing `Config` struct from `internal/config/config.go`.

**Existing Entity** (no changes needed):
```go
// internal/config/config.go
type Config struct {
    Prefix string `yaml:"prefix"`
    Remote string `yaml:"remote"`
    Hooks  Hooks  `yaml:"hooks"`
}

type Hooks struct {
    PreTag   []string `yaml:"pre-tag"`
    PostTag  []string `yaml:"post-tag"`
    PostPush []string `yaml:"post-push"`
}
```

**Default Template** (for init command):
```yaml
# Bumpkin configuration
# See https://github.com/benny123tw/bumpkin for documentation

# Tag prefix (default: "v")
prefix: v

# Git remote (default: "origin")
remote: origin

# Hooks
hooks:
  # Commands to run before creating the tag
  # pre-tag:
  #   - go test ./...
  
  # Commands to run after creating the tag
  # post-tag:
  #   - echo "Tagged ${BUMPKIN_NEW_VERSION}"
  
  # Commands to run after pushing
  # post-push:
  #   - goreleaser release
```

## State Transitions

N/A - This feature does not introduce stateful entities.

## Validation Rules

| Field | Rule | Error Message |
|-------|------|---------------|
| Shell (completion) | Must be one of: bash, zsh, fish, powershell | "Unsupported shell: {shell}" |
| Config file (init) | Must not exist | ".bumpkin.yaml already exists" |
| Repository (current) | Must be a git repository | "not a git repository" |

## Relationships

```
init command ──creates──> .bumpkin.yaml (uses Config struct)
current command ──reads──> git repository (uses existing git.Repository)
completion command ──generates──> shell scripts (uses Cobra built-in)
```

## No API Contracts

This feature consists of CLI commands only. No REST/GraphQL API contracts are applicable.

The `/contracts/` directory is not needed for this feature.

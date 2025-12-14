# Research: Version Tagger CLI

**Feature**: 001-version-tagger
**Date**: 2025-12-14

## Technology Decisions

### 1. TUI Framework: Bubbletea

**Decision**: Use `github.com/charmbracelet/bubbletea` with `bubbles` components

**Rationale**:
- Elm Architecture provides clean separation of state, updates, and views
- Active development and strong community (Charm.sh ecosystem)
- `bubbles` provides ready-made components: List, TextInput, Spinner, Viewport
- `lipgloss` enables consistent styling across terminal themes
- Battle-tested in production tools (gum, soft-serve, etc.)

**Alternatives Considered**:
- `tview`: More widget-oriented but less compositional
- `termui`: Focused on dashboards, less suited for wizard-style UI
- Raw terminal codes: Too low-level, unnecessary complexity

### 2. CLI Framework: Cobra

**Decision**: Use `github.com/spf13/cobra` for command structure

**Rationale**:
- Industry standard for Go CLIs (Kubernetes, Docker, Hugo, GitHub CLI)
- Excellent flag handling with both persistent and command-specific flags
- Built-in help generation and shell completion
- Integrates well with Viper for configuration (optional)

**Alternatives Considered**:
- `urfave/cli`: Good but less widely adopted
- Standard `flag` package: Too basic for complex flag scenarios
- `kong`: Struct-based, different paradigm

### 3. Git Operations: go-git

**Decision**: Use `github.com/go-git/go-git/v5` for git operations

**Rationale**:
- Pure Go implementation - no external git binary dependency
- Full git protocol support including push to remotes
- Can read tags, commits, and repository state
- Cross-platform without needing git installed

**Alternatives Considered**:
- Shell out to `git` command: Adds external dependency, parsing complexity
- `libgit2` bindings: CGO dependency, harder cross-compilation

**Limitations to Note**:
- Some operations slower than native git
- Complex merge scenarios may need fallback to git command
- Tag iteration returns tag object hash, not commit hash (need resolution)

### 4. Semantic Versioning: Masterminds/semver

**Decision**: Use `github.com/Masterminds/semver/v3`

**Rationale**:
- Most mature semver library in Go ecosystem
- Handles version parsing with/without "v" prefix
- Provides comparison operators
- Supports prerelease and build metadata

**Alternatives Considered**:
- `blang/semver`: Good but less actively maintained
- `adamwasila/go-semver`: Has built-in bumping but less tested
- Custom implementation: Unnecessary when good libraries exist

**Implementation Note**: Need to implement bumping functions manually:
```go
func BumpPatch(v *semver.Version) *semver.Version
func BumpMinor(v *semver.Version) *semver.Version
func BumpMajor(v *semver.Version) *semver.Version
```

### 5. Conventional Commits: leodido/go-conventionalcommits

**Decision**: Use `github.com/leodido/go-conventionalcommits` for parsing

**Rationale**:
- Parses conventional commit messages to structured data
- Detects breaking changes (BREAKING CHANGE footer, `!` suffix)
- Extracts type, scope, and description
- Well-maintained and follows spec

**Alternatives Considered**:
- Custom regex parsing: Error-prone, edge cases
- Simple prefix matching: Misses breaking change detection

### 6. Configuration Format

**Decision**: Use YAML configuration file (`.bumpkin.yaml` or `.bumpkin.yml`)

**Rationale**:
- Human-readable and editable
- Standard for CLI tool configuration
- Supports complex structures for hooks
- Can use Viper for loading if needed

**Configuration Structure**:
```yaml
# .bumpkin.yaml
prefix: "v"              # Tag prefix (default: "v")
remote: "origin"         # Git remote (default: "origin")
hooks:
  pre-tag:
    - "./scripts/update-version.sh"
  post-tag:
    - "./scripts/notify.sh"
```

### 7. Hook System Design

**Decision**: Simple shell command execution with environment variables

**Rationale**:
- Maximum flexibility - users write any script
- Language-agnostic (bash, python, node, etc.)
- Familiar pattern from git hooks

**Environment Variables Passed to Hooks**:
- `BUMPKIN_VERSION`: New version being created
- `BUMPKIN_PREVIOUS_VERSION`: Previous version
- `BUMPKIN_TAG`: Full tag name (e.g., "v1.2.3")
- `BUMPKIN_BUMP_TYPE`: Type of bump (major/minor/patch/custom)

**Behavior**:
- Pre-tag hooks run before tag creation
- If pre-tag hook exits non-zero, abort tagging
- Post-tag hooks run after successful push
- Hook failures logged but don't rollback

## Architecture Patterns

### Hybrid CLI/TUI Mode

```
User runs `bumpkin`
        │
        ▼
┌─────────────────────┐
│ Parse CLI flags     │
│ (Cobra root cmd)    │
└─────────────────────┘
        │
        ▼
┌─────────────────────┐     Yes     ┌─────────────────────┐
│ Has --patch/minor/  │────────────▶│ Non-Interactive     │
│ major/version flag? │             │ Execute directly    │
└─────────────────────┘             └─────────────────────┘
        │ No
        ▼
┌─────────────────────┐
│ Launch Bubbletea    │
│ Interactive TUI     │
└─────────────────────┘
        │
        ▼
┌─────────────────────┐
│ Both modes call     │
│ executor.Bump()     │
└─────────────────────┘
```

### TUI State Machine

```
┌──────────────┐
│ Loading      │ ← Initial state, fetch git info
└──────────────┘
        │
        ▼
┌──────────────┐
│ CommitList   │ ← Show commits since last tag
└──────────────┘
        │ Enter
        ▼
┌──────────────┐
│ VersionSelect│ ← Major/Minor/Patch/Custom/Prerelease
└──────────────┘
        │ Enter
        ▼
┌──────────────┐
│ Prerelease   │ ← Only if prerelease selected
│ (optional)   │
└──────────────┘
        │ Enter
        ▼
┌──────────────┐
│ Confirm      │ ← Show summary, confirm/cancel
└──────────────┘
        │ Yes
        ▼
┌──────────────┐
│ Executing    │ ← Run hooks, create tag, push
└──────────────┘
        │
        ▼
┌──────────────┐
│ Done/Error   │ ← Show result
└──────────────┘
```

### Shared Execution Logic

Both TUI and CLI modes delegate to the same executor:

```go
// internal/executor/bump.go
type BumpRequest struct {
    Repository   *git.Repository
    BumpType     BumpType  // Major, Minor, Patch, Custom, Prerelease
    CustomVersion string   // Only for Custom type
    PreType      string    // alpha, beta, rc
    DryRun       bool
    NoPush       bool
    Config       *config.Config
}

type BumpResult struct {
    PreviousVersion string
    NewVersion      string
    TagName         string
    CommitHash      string
    Pushed          bool
}

func Execute(ctx context.Context, req BumpRequest) (*BumpResult, error)
```

## Dependency Versions

```go
require (
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/bubbles v0.18.0
    github.com/charmbracelet/lipgloss v0.10.0
    github.com/spf13/cobra v1.8.0
    github.com/go-git/go-git/v5 v5.11.0
    github.com/Masterminds/semver/v3 v3.2.1
    github.com/leodido/go-conventionalcommits v0.11.0
    gopkg.in/yaml.v3 v3.0.1
)
```

## Testing Strategy

### Unit Tests
- `internal/version/`: Test all bump operations, edge cases
- `internal/conventional/`: Test commit parsing with various formats
- `internal/config/`: Test config loading and validation

### Integration Tests
- Create temporary git repositories with tags/commits
- Test full bump workflow end-to-end
- Test hook execution with mock scripts

### TUI Testing
- Use Bubbletea's test utilities for model testing
- Test state transitions without actual rendering

## Open Questions Resolved

1. **Config file location**: `.bumpkin.yaml` in repository root (standard pattern)
2. **Tag format**: Annotated tags with version as message (matches goreleaser)
3. **No existing tags**: Start from v0.0.0 (documented in edge cases)
4. **Commit listing**: Reverse chronological, first line only for display

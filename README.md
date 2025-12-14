# Bumpkin

A semantic version tagger CLI for git repositories. Inspired by [antfu/bumpp](https://github.com/antfu/bumpp).

## Features

- **Interactive TUI** - Select version bumps with keyboard navigation
- **Non-interactive CLI** - Automate versioning in CI/CD pipelines
- **Conventional Commits** - Auto-detect version bump from commit history
- **Prerelease Support** - Full alpha/beta/rc workflow
- **Hook System** - Run scripts before/after tagging
- **Configurable** - Via CLI flags or `.bumpkin.yml`

## Installation

```bash
go install github.com/benny123tw/bumpkin/cmd/bumpkin@latest
```

Or build from source:

```bash
git clone https://github.com/benny123tw/bumpkin.git
cd bumpkin
just build
```

## Usage

### Interactive Mode

Run without flags to launch the interactive TUI:

```bash
bumpkin
```

Navigate with arrow keys, select a version bump, and confirm.

### Non-Interactive Mode

Use flags for automation:

```bash
# Bump patch version (1.2.3 -> 1.2.4)
bumpkin --patch --yes

# Bump minor version (1.2.3 -> 1.3.0)
bumpkin --minor --yes

# Bump major version (1.2.3 -> 2.0.0)
bumpkin --major --yes

# Set specific version
bumpkin --set-version 2.0.0 --yes

# Auto-detect from conventional commits
bumpkin --conventional --yes
```

### Prerelease Versions

```bash
# Create alpha (1.2.3 -> 1.2.4-alpha.0)
bumpkin --alpha --yes

# Create beta (1.2.4-alpha.0 -> 1.2.4-beta.0)
bumpkin --beta --yes

# Create release candidate (1.2.4-beta.0 -> 1.2.4-rc.0)
bumpkin --rc --yes

# Promote to release (1.2.4-rc.0 -> 1.2.4)
bumpkin --release --yes
```

### Additional Options

```bash
# Preview changes without executing
bumpkin --patch --dry-run

# Create tag but don't push
bumpkin --patch --yes --no-push

# Skip hook execution
bumpkin --patch --yes --no-hooks

# JSON output for scripting
bumpkin --patch --yes --json

# Custom tag prefix (default: v)
bumpkin --patch --yes --prefix "ver"

# Custom remote (default: origin)
bumpkin --patch --yes --remote upstream
```

## Configuration

Create `.bumpkin.yml` in your repository root:

```yaml
# Tag prefix (default: "v")
prefix: "v"

# Git remote (default: "origin")
remote: "origin"

# Hooks
hooks:
  # Run before creating tag (aborts on failure)
  pre-tag:
    - "npm version $BUMPKIN_VERSION --no-git-tag-version"
    - "./scripts/update-changelog.sh"
  
  # Run after creating tag (aborts on failure)
  post-tag:
    - "echo Tagged $BUMPKIN_TAG"
  
  # Run after pushing tag (fail-open: continues on failure, reports warnings)
  post-push:
    - "curl -X POST $SLACK_WEBHOOK -d '{\"text\": \"Released $BUMPKIN_TAG\"}'"
    - "./scripts/notify-team.sh"
```

### Hook Phases

Hooks execute in this order:

```
pre-tag → create tag → post-tag → push → post-push
```

| Phase | Behavior on Failure |
|-------|---------------------|
| `pre-tag` | Aborts - tag not created |
| `post-tag` | Warning - tag already created |
| `post-push` | Warning - tag already pushed (fail-open) |

**Note:** `post-push` hooks use fail-open behavior: if a hook fails, subsequent hooks still execute and warnings are reported. This is ideal for notifications where you don't want one failing webhook to block others.

### Hook Environment Variables

Hooks receive these environment variables:

| Variable | Description |
|----------|-------------|
| `BUMPKIN_VERSION` | New version (without prefix) |
| `BUMPKIN_PREVIOUS_VERSION` | Previous version |
| `BUMPKIN_TAG` | Full tag name (with prefix) |
| `BUMPKIN_PREFIX` | Tag prefix |
| `BUMPKIN_REMOTE` | Remote name |
| `BUMPKIN_COMMIT` | Commit hash being tagged |
| `BUMPKIN_DRY_RUN` | "true" if dry run mode |

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Not a git repository |
| 5 | User cancelled |
| 6 | Hook execution failed |

## Conventional Commits

When using `--conventional`, bumpkin analyzes commits since the last tag:

| Commit Type | Version Bump |
|-------------|--------------|
| `feat!:` or `BREAKING CHANGE:` | Major |
| `feat:` | Minor |
| `fix:`, `docs:`, `chore:`, etc. | Patch |

Example:
```bash
git commit -m "feat: add user authentication"  # -> minor bump
git commit -m "fix: resolve login bug"         # -> patch bump
git commit -m "feat!: redesign API"            # -> major bump
```

## Development

Requires:
- Go 1.24+
- [just](https://github.com/casey/just) (optional, for task running)
- [golangci-lint](https://golangci-lint.run/) (for linting)

```bash
# Run tests
just test

# Run linter
just lint

# Run both
just check

# Build binary
just build

# Run interactive mode
just run
```

## License

MIT

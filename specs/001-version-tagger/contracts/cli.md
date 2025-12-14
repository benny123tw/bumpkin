# CLI Contract: bumpkin

**Version**: 1.0.0
**Date**: 2025-12-14

## Command Structure

```
bumpkin [flags]
```

No subcommands - single command with optional flags for non-interactive mode.

## Flags

### Version Bump Flags (mutually exclusive)

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--patch` | `-p` | bool | Bump patch version (x.y.Z) |
| `--minor` | `-m` | bool | Bump minor version (x.Y.0) |
| `--major` | `-M` | bool | Bump major version (X.0.0) |
| `--version` | `-v` | string | Set specific version |
| `--prerelease` | | string | Create prerelease (alpha/beta/rc) |
| `--release` | `-r` | bool | Remove prerelease suffix |

### Behavior Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--dry-run` | `-d` | bool | false | Show what would happen without executing |
| `--no-push` | | bool | false | Create tag locally but don't push |
| `--no-hooks` | | bool | false | Skip pre/post tag hooks |
| `--yes` | `-y` | bool | false | Skip confirmation in non-interactive mode |

### Configuration Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--remote` | | string | origin | Git remote to push to |
| `--prefix` | | string | v | Tag prefix |
| `--config` | `-c` | string | .bumpkin.yml | Config file path |

### Output Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--quiet` | `-q` | bool | false | Suppress output except errors |
| `--json` | | bool | false | Output result as JSON (non-interactive only) |

## Usage Examples

### Interactive Mode (default)

```bash
# Launch interactive TUI
bumpkin

# Interactive with custom remote
bumpkin --remote upstream
```

### Non-Interactive Mode

```bash
# Bump patch version
bumpkin --patch

# Bump minor version with confirmation skip
bumpkin --minor --yes

# Set specific version
bumpkin --version 2.0.0

# Create alpha prerelease
bumpkin --prerelease alpha

# Dry run to see what would happen
bumpkin --patch --dry-run

# Create tag locally only (don't push)
bumpkin --minor --no-push

# JSON output for scripting
bumpkin --patch --yes --json
```

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error (git operation failed, hook failed, etc.) |
| 2 | Invalid arguments |
| 3 | Not a git repository |
| 4 | No commits since last tag (with warning) |
| 5 | User cancelled operation |
| 6 | Hook execution failed |

## Output Formats

### Standard Output (default)

```
Current version: v1.2.3
Commits since last tag:
  abc1234 feat: add new feature
  def5678 fix: resolve bug
  
Creating tag: v1.3.0
Pushing to origin...
Done! Tagged v1.3.0
```

### JSON Output (`--json` flag)

```json
{
  "previous_version": "v1.2.3",
  "new_version": "v1.3.0",
  "tag_name": "v1.3.0",
  "commit_hash": "abc123def456...",
  "pushed": true,
  "commits": [
    {
      "hash": "abc1234",
      "message": "feat: add new feature",
      "type": "feat",
      "is_breaking": false
    }
  ],
  "hooks": {
    "pre_tag": {"success": true, "exit_code": 0},
    "post_tag": {"success": true, "exit_code": 0}
  }
}
```

### Dry Run Output

```
[DRY RUN] Would create tag: v1.3.0
[DRY RUN] Would push to remote: origin
[DRY RUN] No changes made
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `BUMPKIN_CONFIG` | Override config file path |
| `BUMPKIN_PREFIX` | Override tag prefix |
| `BUMPKIN_REMOTE` | Override git remote |
| `NO_COLOR` | Disable colored output |

## Hook Environment Variables

When hooks are executed, these variables are set:

| Variable | Description |
|----------|-------------|
| `BUMPKIN_VERSION` | New version (without prefix) |
| `BUMPKIN_PREVIOUS_VERSION` | Previous version (without prefix) |
| `BUMPKIN_TAG` | Full tag name (with prefix) |
| `BUMPKIN_BUMP_TYPE` | Bump type (patch/minor/major/custom/prerelease) |
| `BUMPKIN_COMMIT` | Commit hash being tagged |
| `BUMPKIN_REMOTE` | Remote name |
| `BUMPKIN_DRY_RUN` | "true" if dry run mode |

## Configuration File

Location: `.bumpkin.yml` (or `.bumpkin.yaml`) in repository root

```yaml
# Tag prefix (default: "v")
prefix: "v"

# Git remote (default: "origin")
remote: "origin"

# Tag message template (supports Go templates)
tag_message: "Release {{.Version}}"

# Hooks
hooks:
  # Commands to run before creating tag
  pre_tag:
    - "./scripts/update-changelog.sh"
    - "npm version {{.Version}} --no-git-tag-version"
  
  # Commands to run after pushing tag
  post_tag:
    - "./scripts/notify-release.sh"
```

### Template Variables in Config

| Variable | Description |
|----------|-------------|
| `{{.Version}}` | New version (without prefix) |
| `{{.PreviousVersion}}` | Previous version |
| `{{.Tag}}` | Full tag name |
| `{{.BumpType}}` | Bump type |

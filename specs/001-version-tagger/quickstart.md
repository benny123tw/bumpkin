# Quickstart: bumpkin

**Version**: 1.0.0
**Date**: 2025-12-14

## Installation

```bash
# Using go install
go install github.com/benny123tw/bumpkin/cmd/bumpkin@latest

# Or build from source
git clone https://github.com/benny123tw/bumpkin.git
cd bumpkin
go build -o bumpkin ./cmd/bumpkin
```

## Basic Usage

### Interactive Mode

Run `bumpkin` without arguments to launch the interactive TUI:

```bash
cd your-project
bumpkin
```

The TUI will:
1. Show commits since the last tag
2. Display the current version
3. Present version bump options
4. Create and push the new tag

### Non-Interactive Mode

For CI/CD or quick bumps:

```bash
# Bump patch version (1.2.3 → 1.2.4)
bumpkin --patch

# Bump minor version (1.2.3 → 1.3.0)
bumpkin --minor

# Bump major version (1.2.3 → 2.0.0)
bumpkin --major

# Set specific version
bumpkin --version 2.0.0

# Skip confirmation prompt
bumpkin --patch --yes
```

## First Run

If your repository has no tags, bumpkin starts from `v0.0.0`:

```bash
# Create initial release
bumpkin --minor  # Creates v0.1.0
```

## Prerelease Versions

```bash
# Create alpha release (1.2.3 → 1.2.4-alpha.0)
bumpkin --prerelease alpha

# Continue alpha series (1.2.4-alpha.0 → 1.2.4-alpha.1)
bumpkin --prerelease alpha

# Move to beta (1.2.4-alpha.1 → 1.2.4-beta.0)
bumpkin --prerelease beta

# Final release (1.2.4-beta.0 → 1.2.4)
bumpkin --release
```

## Dry Run

Preview what would happen without making changes:

```bash
bumpkin --patch --dry-run
```

Output:
```
[DRY RUN] Current version: v1.2.3
[DRY RUN] Would create tag: v1.2.4
[DRY RUN] Would push to remote: origin
[DRY RUN] No changes made
```

## Configuration

Create `.bumpkin.yaml` in your repository root:

```yaml
prefix: "v"
remote: "origin"

hooks:
  pre_tag:
    - "npm version {{.Version}} --no-git-tag-version"
    - "git add package.json"
    - "git commit -m 'chore: bump version to {{.Version}}'"
  post_tag:
    - "echo 'Released {{.Tag}}!'"
```

## Common Workflows

### Simple Project (no version files)

```bash
# Just tag and push
bumpkin --patch
```

### Node.js Project

`.bumpkin.yaml`:
```yaml
hooks:
  pre_tag:
    - "npm version {{.Version}} --no-git-tag-version"
    - "git add package.json package-lock.json"
    - "git commit -m 'chore: bump version to {{.Version}}'"
```

### Go Project

`.bumpkin.yaml`:
```yaml
hooks:
  pre_tag:
    - "go generate ./..."
    - "git add ."
    - "git diff --cached --quiet || git commit -m 'chore: update generated files for {{.Version}}'"
```

### CI/CD Integration

```yaml
# GitHub Actions example
- name: Bump version
  run: bumpkin --patch --yes --json > version.json
  
- name: Read new version
  id: version
  run: echo "version=$(jq -r .new_version version.json)" >> $GITHUB_OUTPUT
```

## Troubleshooting

### "Not a git repository"

Make sure you're in a directory with a `.git` folder:

```bash
cd your-project
git status  # Should show repository status
bumpkin
```

### "No remote configured"

Add a remote or use `--no-push`:

```bash
git remote add origin https://github.com/you/repo.git
bumpkin --patch
```

Or skip pushing:

```bash
bumpkin --patch --no-push
```

### "Hook failed"

Check your hook scripts:
1. Ensure scripts are executable: `chmod +x scripts/your-hook.sh`
2. Test hooks manually before running bumpkin
3. Use `--no-hooks` to skip hooks temporarily

### Push authentication failed

Make sure you have push access:
- For HTTPS: Check credentials or use a token
- For SSH: Ensure your SSH key is configured

## Getting Help

```bash
bumpkin --help
```

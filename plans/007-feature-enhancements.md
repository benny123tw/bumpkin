# Feature Enhancements Plan

## Overview

This plan identifies feature enhancements and new feature requests based on codebase analysis and common use cases for semantic versioning tools.

---

## High Priority

### 1. Changelog Generation

**Description**: Automatically generate or update CHANGELOG.md from conventional commits when creating a new version tag.

**Use Case**: Teams want automated release notes without manual changelog maintenance.

**Implementation**:
- Add `--changelog` flag to enable changelog generation
- Parse commits since last tag using existing conventional commit analyzer
- Group commits by type (features, fixes, breaking changes, etc.)
- Support keep-a-changelog format
- Optionally create as post-push hook or integrated command

**Configuration**:
```yaml
changelog:
  enabled: true
  path: CHANGELOG.md
  template: keepachangelog  # or custom template path
```

**Effort**: Medium

---

### 2. Dry-Run Visual Diff

**Description**: Show exact commands and changes that will be executed in dry-run mode with clear visual formatting.

**Use Case**: Users want to preview exactly what will happen before committing to a version bump.

**Implementation**:
- Enhance `--dry-run` output to show:
  - Tag that will be created
  - Git commands that will run
  - Hooks that will execute
  - Remote that will receive the push
- Use color-coded diff-style output

**Effort**: Small

---

### 3. Tag Annotation Templates

**Description**: Customizable tag message formatting with template support.

**Use Case**: Teams have specific tag message requirements (include commit summary, issue links, etc.).

**Implementation**:
- Add template support for tag messages
- Provide template variables: `{{.Version}}`, `{{.PreviousVersion}}`, `{{.Commits}}`, `{{.Date}}`
- Allow template in config file

**Configuration**:
```yaml
tag:
  template: |
    Release {{.Version}}

    Changes since {{.PreviousVersion}}:
    {{range .Commits}}- {{.Type}}: {{.Description}}
    {{end}}
```

**Effort**: Medium

---

### 4. Version Constraints

**Description**: Allow specifying minimum/maximum version boundaries to prevent accidental downgrades or major jumps.

**Use Case**: Prevent mistakes in CI/CD pipelines or protect against version conflicts.

**Implementation**:
- Add `--min-version` and `--max-version` flags
- Add validation before tag creation
- Support in config file

**Configuration**:
```yaml
constraints:
  min: "1.0.0"
  max: "2.0.0"  # optional, for pre-2.0 projects
```

**Effort**: Small

---

## Medium Priority

### 5. Interactive Config Builder

**Description**: TUI-based `.bumpkin.yaml` generator with guided setup.

**Use Case**: New users find it easier to configure through interactive prompts than reading documentation.

**Implementation**:
- Enhance `bumpkin init` to launch interactive TUI
- Step through common options: prefix, hooks, conventional commits
- Show preview of generated config
- Validate configuration before saving

**Effort**: Medium

---

### 6. Hook Output Streaming

**Description**: Display hook execution output in real-time instead of blocking until completion.

**Use Case**: Long-running hooks (e.g., build scripts) should show progress.

**Implementation**:
- Stream stdout/stderr from hook commands
- Show in TUI as scrollable output pane
- Maintain separation between different hooks
- Support for non-interactive mode with line buffering

**Effort**: Medium

---

### 7. Multi-Remote Support

**Description**: Push tags to multiple remotes with configurable behavior.

**Use Case**: Projects mirrored to multiple Git hosts (GitHub + GitLab, or internal + public).

**Implementation**:
- Add `remotes` config option accepting array
- Add `--remote` flag accepting comma-separated list
- Push to all specified remotes with individual error handling
- Report which remotes succeeded/failed

**Configuration**:
```yaml
remotes:
  - origin
  - mirror
```

**Effort**: Small

---

### 8. Commit Message Filtering

**Description**: Skip certain commit types from analysis or include only specific types.

**Use Case**: Teams want to exclude chore/docs commits from bump calculations or focus only on feat/fix.

**Implementation**:
- Add filter options in config
- Support include/exclude patterns
- Apply during conventional commit analysis

**Configuration**:
```yaml
conventional:
  include: [feat, fix, perf]
  exclude: [chore, docs, style]
```

**Effort**: Small

---

### 9. GPG Tag Signing

**Description**: Support for cryptographically signed tags.

**Use Case**: Security-conscious projects requiring verifiable releases.

**Implementation**:
- Add `--sign` flag
- Use configured GPG key or prompt for selection
- Pass signing option to git tag command
- Verify GPG is available before attempting

**Effort**: Small

---

### 10. Monorepo Support

**Description**: Tag multiple packages independently within one repository.

**Use Case**: Monorepos with multiple versioned packages (e.g., `@scope/pkg-a@1.0.0`, `@scope/pkg-b@2.0.0`).

**Implementation**:
- Add `--package` or `--scope` flag
- Support package-specific prefixes
- Read package paths from config
- Analyze commits per package based on changed files

**Configuration**:
```yaml
packages:
  - name: core
    path: packages/core
    prefix: core/v
  - name: cli
    path: packages/cli
    prefix: cli/v
```

**Effort**: Large

---

## Low Priority

### 11. Export Formats

**Description**: Output version bump results in multiple formats (JSON, YAML, CSV).

**Use Case**: CI/CD integration where structured output is parsed by other tools.

**Implementation**:
- Extend existing `--json` flag to `--output-format`
- Support json, yaml, csv, text formats
- Include all relevant data (version, tag, commits, etc.)

**Effort**: Small

---

### 12. Metrics Collection

**Description**: Report timing statistics per phase of the bump operation.

**Use Case**: Debugging slow operations and optimizing CI/CD pipelines.

**Implementation**:
- Add `--metrics` flag
- Measure time for each phase: load, analyze, hooks, tag, push
- Output as part of JSON/structured output
- Optionally send to external metrics service

**Effort**: Small

---

### 13. Pre-Release Branch Protection

**Description**: Prevent creating stable releases from non-main branches.

**Use Case**: Enforce release workflow policies.

**Implementation**:
- Add `release-branches` config option
- Warn or block when releasing from non-allowed branch
- Override with `--force` flag

**Configuration**:
```yaml
release-branches:
  - main
  - master
  - release/*
```

**Effort**: Small

---

### 14. Version History Command

**Description**: Add `bumpkin history` command to show version release history.

**Use Case**: Quick reference for release timeline without opening git log.

**Implementation**:
- New subcommand `bumpkin history`
- List all version tags with dates
- Show conventional commit summary per version
- Support `--limit` flag

**Effort**: Small

---

### 15. Undo Last Tag

**Description**: Add `bumpkin undo` command to remove the last created tag.

**Use Case**: Quick recovery from accidental version bumps.

**Implementation**:
- Delete local tag
- Optionally delete remote tag (with confirmation)
- Only allow undoing tags created by bumpkin (track in local state)
- Add `--force` for remote deletion

**Effort**: Small

---

## Implementation Roadmap

### Phase 1: Core Improvements
- Dry-run visual diff (#2)
- Version constraints (#4)
- GPG tag signing (#9)

### Phase 2: Workflow Enhancements
- Changelog generation (#1)
- Tag annotation templates (#3)
- Commit message filtering (#8)

### Phase 3: User Experience
- Interactive config builder (#5)
- Hook output streaming (#6)
- Export formats (#11)

### Phase 4: Advanced Features
- Multi-remote support (#7)
- Monorepo support (#10)
- Version history command (#14)
- Undo last tag (#15)

### Phase 5: Observability
- Metrics collection (#12)
- Pre-release branch protection (#13)

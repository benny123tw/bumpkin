# Data Model: Version Tagger CLI

**Feature**: 001-version-tagger
**Date**: 2025-12-14

## Core Entities

### Version

Represents a semantic version following the semver 2.0 specification.

| Field | Type | Description |
|-------|------|-------------|
| Major | uint64 | Major version number |
| Minor | uint64 | Minor version number |
| Patch | uint64 | Patch version number |
| Prerelease | string | Prerelease identifier (e.g., "alpha.0", "beta.1", "rc.2") |
| Metadata | string | Build metadata (optional, not used in comparison) |

**Validation Rules**:
- Major, Minor, Patch must be non-negative integers
- Prerelease follows pattern: `<type>.<number>` where type is alpha/beta/rc
- Version string format: `v<major>.<minor>.<patch>[-<prerelease>][+<metadata>]`

**State Transitions**:
```
v1.0.0 --[patch]--> v1.0.1
v1.0.0 --[minor]--> v1.1.0
v1.0.0 --[major]--> v2.0.0
v1.0.0 --[prerelease:alpha]--> v1.0.1-alpha.0
v1.0.1-alpha.0 --[prerelease:alpha]--> v1.0.1-alpha.1
v1.0.1-alpha.1 --[prerelease:beta]--> v1.0.1-beta.0
v1.0.1-beta.0 --[release]--> v1.0.1
```

---

### Commit

Represents a git commit with parsed conventional commit information.

| Field | Type | Description |
|-------|------|-------------|
| Hash | string | Full SHA-1 commit hash (40 chars) |
| ShortHash | string | Abbreviated hash (7 chars) |
| Message | string | Full commit message |
| Subject | string | First line of commit message |
| Author | string | Author name |
| AuthorEmail | string | Author email |
| Timestamp | time.Time | Commit timestamp |
| Type | string | Conventional commit type (feat, fix, etc.) |
| Scope | string | Conventional commit scope (optional) |
| IsBreaking | bool | Whether commit contains breaking change |

**Validation Rules**:
- Hash must be valid hexadecimal string
- Subject extracted as first line of Message (before first newline)
- Type parsed from conventional commit prefix (empty if non-conventional)

---

### Tag

Represents a git tag pointing to a specific commit.

| Field | Type | Description |
|-------|------|-------------|
| Name | string | Tag name (e.g., "v1.2.3") |
| Version | *Version | Parsed semantic version (nil if not semver) |
| CommitHash | string | SHA of the tagged commit |
| Tagger | string | Name of person who created tag |
| TaggerEmail | string | Email of tagger |
| Message | string | Tag annotation message |
| Timestamp | time.Time | Tag creation timestamp |
| IsAnnotated | bool | Whether tag is annotated (vs lightweight) |

**Validation Rules**:
- Name must match configured prefix pattern (default: `v<semver>`)
- Version extracted by stripping prefix and parsing semver

---

### Config

User configuration for bumpkin behavior.

| Field | Type | Description |
|-------|------|-------------|
| Prefix | string | Tag prefix (default: "v") |
| Remote | string | Git remote name (default: "origin") |
| PreTagHooks | []string | Commands to run before tagging |
| PostTagHooks | []string | Commands to run after tagging |
| TagMessage | string | Template for tag annotation (default: "Release {{.Version}}") |

**Validation Rules**:
- Prefix must be a valid tag name prefix (no special characters)
- Remote must exist in repository configuration
- Hook commands must be valid shell commands

**Default Configuration**:
```yaml
prefix: "v"
remote: "origin"
tag_message: "Release {{.Version}}"
hooks:
  pre_tag: []
  post_tag: []
```

---

### BumpType

Enumeration of version bump operations.

| Value | Description | Example |
|-------|-------------|---------|
| Patch | Increment patch version | v1.0.0 → v1.0.1 |
| Minor | Increment minor, reset patch | v1.0.0 → v1.1.0 |
| Major | Increment major, reset minor and patch | v1.0.0 → v2.0.0 |
| Custom | User-specified version | v1.0.0 → v3.0.0 |
| PrereleaseAlpha | Create/increment alpha prerelease | v1.0.0 → v1.0.1-alpha.0 |
| PrereleaseBeta | Create/increment beta prerelease | v1.0.0 → v1.0.1-beta.0 |
| PrereleaseRC | Create/increment RC prerelease | v1.0.0 → v1.0.1-rc.0 |
| Release | Remove prerelease suffix | v1.0.1-rc.0 → v1.0.1 |

---

### BumpRequest

Input for the bump execution operation.

| Field | Type | Description |
|-------|------|-------------|
| BumpType | BumpType | Type of version bump |
| CustomVersion | string | Custom version string (only for Custom type) |
| DryRun | bool | If true, show actions without executing |
| NoPush | bool | If true, create tag but don't push |
| NoHooks | bool | If true, skip hook execution |

---

### BumpResult

Output from the bump execution operation.

| Field | Type | Description |
|-------|------|-------------|
| PreviousVersion | string | Version before bump |
| NewVersion | string | New version created |
| TagName | string | Full tag name (with prefix) |
| CommitHash | string | Commit the tag points to |
| TagCreated | bool | Whether tag was created |
| Pushed | bool | Whether tag was pushed to remote |
| HooksExecuted | []HookResult | Results of hook executions |

---

### HookResult

Result of a single hook execution.

| Field | Type | Description |
|-------|------|-------------|
| Command | string | The command that was executed |
| ExitCode | int | Exit code (0 = success) |
| Stdout | string | Standard output |
| Stderr | string | Standard error |
| Duration | time.Duration | Execution time |

---

### RepositoryState

Current state of the git repository.

| Field | Type | Description |
|-------|------|-------------|
| CurrentBranch | string | Name of current branch |
| LatestTag | *Tag | Most recent semver tag (nil if none) |
| CommitsSinceTag | []Commit | Commits between latest tag and HEAD |
| IsDirty | bool | Whether working tree has uncommitted changes |
| RemoteURL | string | URL of configured remote |
| HasRemote | bool | Whether remote is configured |

---

## Relationships

```
RepositoryState
    │
    ├── 1:1 ──► LatestTag (Tag)
    │               │
    │               └── 1:1 ──► Version
    │
    └── 1:N ──► CommitsSinceTag (Commit)

Config
    │
    └── 1:N ──► Hooks (pre-tag, post-tag)

BumpRequest + RepositoryState ──► BumpResult
                                      │
                                      └── 1:N ──► HookResult
```

## State Machine: TUI Application

```
                    ┌─────────────────┐
                    │     Init        │
                    └────────┬────────┘
                             │ Load repo state
                             ▼
                    ┌─────────────────┐
         Error ◄────│    Loading      │
                    └────────┬────────┘
                             │ Success
                             ▼
                    ┌─────────────────┐
         Quit ◄─────│   CommitList    │───────► Show commits
                    └────────┬────────┘         Current version
                             │ Enter/Select
                             ▼
                    ┌─────────────────┐
         Back ◄─────│ VersionSelect   │───────► Bump options
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
              ▼              │              ▼
     ┌────────────┐          │     ┌────────────┐
     │  Custom    │          │     │ Prerelease │
     │  Input     │          │     │  Select    │
     └─────┬──────┘          │     └─────┬──────┘
           │                 │           │
           └─────────────────┼───────────┘
                             │
                             ▼
                    ┌─────────────────┐
         Back ◄─────│    Confirm      │───────► Summary
                    └────────┬────────┘
                             │ Yes
                             ▼
                    ┌─────────────────┐
                    │   Executing     │───────► Hooks, Tag, Push
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              ▼                             ▼
     ┌────────────────┐            ┌────────────────┐
     │    Success     │            │     Error      │
     └────────────────┘            └────────────────┘
              │                             │
              └──────────────┬──────────────┘
                             │ Any key
                             ▼
                          [Exit]
```

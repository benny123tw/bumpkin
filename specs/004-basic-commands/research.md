# Research: Basic Subcommands

**Feature**: 004-basic-commands  
**Date**: 2024-12-14

## Executive Summary

This feature adds 5 subcommands (`version`, `help`, `init`, `current`, `completion`) to bumpkin. Research confirms Cobra provides built-in support for most of these patterns, requiring minimal custom implementation.

## Research Tasks

### 1. Cobra Subcommand Patterns

**Decision**: Use Cobra's standard subcommand registration pattern with `AddCommand()`.

**Rationale**: 
- Cobra is already the CLI framework in use (v1.10.2)
- Native support for `help` subcommand (automatic)
- Native support for `completion` subcommand via `GenBashCompletion`, `GenZshCompletion`, `GenFishCompletion`, `GenPowerShellCompletion`
- `version` subcommand is a simple custom command

**Alternatives Considered**:
- urfave/cli: Would require migration, not worth it
- Custom flag handling: Less idiomatic, more work

### 2. Help Subcommand Implementation

**Decision**: Cobra provides automatic `help` subcommand - no implementation needed.

**Rationale**:
- Running `bumpkin help` already works with Cobra
- Running `bumpkin help <subcommand>` also works automatically
- The `--help` flag works identically to `help` subcommand

**Alternatives Considered**:
- Custom help command: Unnecessary duplication of Cobra's built-in

### 3. Version Subcommand Implementation

**Decision**: Create a simple `version` subcommand that reuses existing version output logic.

**Rationale**:
- Mirrors the `--version` flag behavior
- Follows Go CLI conventions (e.g., `go version`, `docker version`)
- Single file implementation: `internal/cli/version.go`

**Implementation Pattern**:
```go
var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print version information",
    Run: func(cmd *cobra.Command, args []string) {
        // Reuse existing version output logic
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}
```

### 4. Init Subcommand Implementation

**Decision**: Create `init` command that generates `.bumpkin.yaml` with commented defaults.

**Rationale**:
- Follows goreleaser's behavior (fail if config exists)
- Provides a quick start for new users
- Single file implementation: `internal/cli/init.go`

**Implementation Pattern**:
```go
var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Create a .bumpkin.yaml configuration file",
    RunE: func(cmd *cobra.Command, args []string) error {
        if _, err := os.Stat(".bumpkin.yaml"); err == nil {
            return errors.New(".bumpkin.yaml already exists")
        }
        // Write default config with comments
    },
}
```

**Config Template** (with comments):
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

### 5. Current Subcommand Implementation

**Decision**: Create `current` command that displays the latest version tag.

**Rationale**:
- Useful for scripting and CI/CD pipelines
- Leverages existing `repo.LatestTag()` functionality
- Single file implementation: `internal/cli/current.go`

**Implementation Pattern**:
```go
var currentCmd = &cobra.Command{
    Use:   "current",
    Short: "Show the current version (latest tag)",
    RunE: func(cmd *cobra.Command, args []string) error {
        repo, err := git.OpenFromCurrent()
        if err != nil {
            return fmt.Errorf("not a git repository")
        }
        tag, err := repo.LatestTag(prefix)
        if err != nil || tag == nil {
            fmt.Println("No version tags found")
            return nil
        }
        fmt.Println(tag.Name)
        return nil
    },
}
```

**Flags**:
- `--prefix, -p`: Tag prefix filter (default: "v", inherited from root)

### 6. Completion Subcommand Implementation

**Decision**: Use Cobra's built-in completion generation with a custom wrapper command.

**Rationale**:
- Cobra provides `GenBashCompletion`, `GenZshCompletion`, `GenFishCompletion`, `GenPowerShellCompletion`
- Common pattern in Go CLI tools (kubectl, hugo, helm)
- Single file implementation: `internal/cli/completion.go`

**Implementation Pattern**:
```go
var completionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish|powershell]",
    Short: "Generate shell completion scripts",
    ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        switch args[0] {
        case "bash":
            return rootCmd.GenBashCompletion(os.Stdout)
        case "zsh":
            return rootCmd.GenZshCompletion(os.Stdout)
        case "fish":
            return rootCmd.GenFishCompletion(os.Stdout, true)
        case "powershell":
            return rootCmd.GenPowerShellCompletion(os.Stdout)
        default:
            return fmt.Errorf("unsupported shell: %s", args[0])
        }
    },
}
```

**Usage Instructions** (to include in help text):
- Bash: `source <(bumpkin completion bash)`
- Zsh: `bumpkin completion zsh > "${fpath[1]}/_bumpkin"`
- Fish: `bumpkin completion fish | source`
- PowerShell: `bumpkin completion powershell | Out-String | Invoke-Expression`

## Dependencies

No new dependencies required. All functionality is provided by:
- `github.com/spf13/cobra` v1.10.2 (existing)
- Standard library (`os`, `fmt`, `errors`)

## File Structure

| File | Purpose |
|------|---------|
| `internal/cli/version.go` | version subcommand |
| `internal/cli/init.go` | init subcommand |
| `internal/cli/current.go` | current subcommand |
| `internal/cli/completion.go` | completion subcommand |

Note: `help` subcommand is automatic via Cobra - no implementation file needed.

## Testing Strategy

Each subcommand will have corresponding tests:
- `internal/cli/version_test.go`
- `internal/cli/init_test.go`
- `internal/cli/current_test.go`
- `internal/cli/completion_test.go`

Tests will follow existing patterns in `internal/cli/root_test.go`.

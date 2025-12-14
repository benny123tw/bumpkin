# Quickstart: Basic Subcommands

**Feature**: 004-basic-commands  
**Date**: 2024-12-14

## Prerequisites

- Go 1.24+
- Existing bumpkin codebase checked out

## Implementation Steps

### Step 1: Create version subcommand

Create `internal/cli/version.go`:

```go
package cli

import (
    "fmt"
    "github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print version information",
    Long:  "Print the version, commit hash, and build date of bumpkin.",
    Run: func(cmd *cobra.Command, args []string) {
        commit := GitCommit
        if len(commit) > 7 {
            commit = commit[:7]
        }
        fmt.Fprintf(cmd.OutOrStdout(), "bumpkin %s (%s, built %s)\n", AppVersion, commit, BuildDate)
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}
```

### Step 2: Create init subcommand

Create `internal/cli/init.go`:

```go
package cli

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

const configTemplate = `# Bumpkin configuration
prefix: v
remote: origin

hooks:
  # pre-tag:
  #   - go test ./...
  # post-tag:
  #   - echo "Tagged ${BUMPKIN_NEW_VERSION}"
  # post-push:
  #   - goreleaser release
`

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Create a .bumpkin.yaml configuration file",
    RunE: func(cmd *cobra.Command, args []string) error {
        if _, err := os.Stat(".bumpkin.yaml"); err == nil {
            return fmt.Errorf(".bumpkin.yaml already exists")
        }
        return os.WriteFile(".bumpkin.yaml", []byte(configTemplate), 0644)
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
}
```

### Step 3: Create current subcommand

Create `internal/cli/current.go`:

```go
package cli

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/benny123tw/bumpkin/internal/git"
)

var currentCmd = &cobra.Command{
    Use:   "current",
    Short: "Show the current version (latest tag)",
    RunE: func(cmd *cobra.Command, args []string) error {
        repo, err := git.OpenFromCurrent()
        if err != nil {
            return fmt.Errorf("not a git repository")
        }
        
        prefix, _ := cmd.Flags().GetString("prefix")
        tag, err := repo.LatestTag(prefix)
        if err != nil || tag == nil {
            fmt.Fprintln(cmd.OutOrStdout(), "No version tags found")
            return nil
        }
        fmt.Fprintln(cmd.OutOrStdout(), tag.Name)
        return nil
    },
}

func init() {
    currentCmd.Flags().StringP("prefix", "p", "v", "Tag prefix")
    rootCmd.AddCommand(currentCmd)
}
```

### Step 4: Create completion subcommand

Create `internal/cli/completion.go`:

```go
package cli

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish|powershell]",
    Short: "Generate shell completion scripts",
    Long: `Generate shell completion scripts for bumpkin.

To load completions:

Bash:
  source <(bumpkin completion bash)

Zsh:
  bumpkin completion zsh > "${fpath[1]}/_bumpkin"

Fish:
  bumpkin completion fish | source

PowerShell:
  bumpkin completion powershell | Out-String | Invoke-Expression
`,
    ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
    Args:      cobra.ExactArgs(1),
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

func init() {
    rootCmd.AddCommand(completionCmd)
}
```

### Step 5: Verify help subcommand

Cobra provides `help` automatically. Verify it works:

```bash
go run ./cmd/bumpkin help
go run ./cmd/bumpkin help version
go run ./cmd/bumpkin help init
```

## Testing

Run the test suite:

```bash
go test ./internal/cli/...
```

Run linting:

```bash
golangci-lint run
```

## Verification Checklist

- [ ] `bumpkin version` outputs version info
- [ ] `bumpkin --version` still works (backward compatibility)
- [ ] `bumpkin help` shows all commands
- [ ] `bumpkin help <cmd>` shows specific command help
- [ ] `bumpkin init` creates `.bumpkin.yaml`
- [ ] `bumpkin init` fails if config exists
- [ ] `bumpkin current` shows latest tag
- [ ] `bumpkin current` handles no-tags case
- [ ] `bumpkin completion bash` outputs valid script
- [ ] `bumpkin completion zsh` outputs valid script
- [ ] `bumpkin completion fish` outputs valid script
- [ ] `bumpkin completion powershell` outputs valid script
- [ ] All tests pass
- [ ] golangci-lint passes

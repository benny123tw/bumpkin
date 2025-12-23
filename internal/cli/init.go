package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const configTemplate = `# Bumpkin configuration
# See https://github.com/benny123tw/bumpkin for documentation

# Tag prefix (default: "v")
prefix: v

# Git remote (default: "origin")
remote: origin

# Hooks - commands to run at different stages
hooks:
  # Commands to run before creating the tag
  # pre-tag:
  #   - go test ./...
  #   - golangci-lint run

  # Commands to run after creating the tag
  # post-tag:
  #   - echo "Tagged ${BUMPKIN_NEW_VERSION}"

  # Commands to run after pushing
  # post-push:
  #   - goreleaser release
`

type initCommand struct {
	cmd *cobra.Command
}

// newInitCommand constructs an initCommand with a Cobra command configured to create
// a starter .bumpkin.yaml configuration file in the current directory.
// The returned initCommand contains the prepared Cobra command accessible via its
// cmd field.
func newInitCommand() *initCommand {
	c := &initCommand{}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Create a .bumpkin.yaml configuration file",
		Long: `Create a .bumpkin.yaml configuration file with default settings.

This command creates a starter configuration file in the current directory
with sensible defaults and commented examples for hooks.`,
		RunE: c.execute,
	}

	c.cmd = initCmd
	return c
}

func (c *initCommand) execute(cmd *cobra.Command, _ []string) error {
	const defaultConfigFile = ".bumpkin.yaml"
	configFiles := []string{defaultConfigFile, ".bumpkin.yml"}

	// Check if config already exists
	for _, f := range configFiles {
		if _, err := os.Stat(f); err == nil {
			return fmt.Errorf("%s already exists", f)
		}
	}

	// Write config file
	//nolint:gosec // Config file needs to be readable by user
	if err := os.WriteFile(defaultConfigFile, []byte(configTemplate), 0o644); err != nil {
		return fmt.Errorf("failed to create %s: %w", defaultConfigFile, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", defaultConfigFile)
	return nil
}

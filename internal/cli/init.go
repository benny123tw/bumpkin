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

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a .bumpkin.yaml configuration file",
	Long: `Create a .bumpkin.yaml configuration file with default settings.

This command creates a starter configuration file in the current directory
with sensible defaults and commented examples for hooks.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		configFile := ".bumpkin.yaml"

		// Check if config already exists (both .yaml and .yml)
		if _, err := os.Stat(".bumpkin.yaml"); err == nil {
			return fmt.Errorf(".bumpkin.yaml already exists")
		}
		if _, err := os.Stat(".bumpkin.yml"); err == nil {
			return fmt.Errorf(".bumpkin.yml already exists")
		}

		// Write config file
		//nolint:gosec // Config file needs to be readable by user
		if err := os.WriteFile(configFile, []byte(configTemplate), 0o644); err != nil {
			return fmt.Errorf("failed to create %s: %w", configFile, err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", configFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

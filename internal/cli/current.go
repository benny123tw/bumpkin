package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/benny123tw/bumpkin/internal/git"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the current version (latest tag)",
	Long: `Show the current version by displaying the latest semver tag.

This command is useful for scripting and CI/CD pipelines where you need
to quickly check the current version without launching the interactive UI.`,
	RunE: runCurrent,
}

func init() {
	currentCmd.Flags().StringP("prefix", "p", "v", "Tag prefix to filter versions")
	rootCmd.AddCommand(currentCmd)
}

func runCurrent(cmd *cobra.Command, _ []string) error {
	prefix, _ := cmd.Flags().GetString("prefix")

	repo, err := git.OpenFromCurrent()
	if err != nil {
		return fmt.Errorf("not a git repository")
	}

	tag, err := repo.LatestTag(prefix)
	if err != nil {
		return fmt.Errorf("failed to get latest tag: %w", err)
	}

	if tag == nil {
		fmt.Fprintln(cmd.OutOrStdout(), "No version tags found")
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), tag.Name)
	return nil
}

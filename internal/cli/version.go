package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print the version, commit hash, and build date of bumpkin.",
	Run: func(cmd *cobra.Command, _ []string) {
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

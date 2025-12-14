package cli

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print the version, commit hash, and build date of bumpkin.",
	Run: func(cmd *cobra.Command, _ []string) {
		PrintVersionInfo(cmd)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

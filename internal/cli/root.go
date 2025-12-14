package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/benny123tw/bumpkin/internal/git"
	"github.com/benny123tw/bumpkin/internal/tui"
)

var (
	// Flags
	prefix  string
	remote  string
	dryRun  bool
	noPush  bool
	version bool
)

// Version information (set at build time)
var (
	Version   = "dev"
	BuildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "bumpkin",
	Short: "Semantic version tagger for git repositories",
	Long: `Bumpkin is a CLI tool that helps you tag commits by analyzing
conventional commit history and providing version options to select or customize.

Run without flags for interactive mode, or use flags for automation.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if version {
			fmt.Printf("bumpkin %s (built %s)\n", Version, BuildDate)
			return nil
		}

		// Open the repository from current directory
		repo, err := git.OpenFromCurrent()
		if err != nil {
			return fmt.Errorf("failed to open repository: %w", err)
		}

		// Launch interactive TUI
		cfg := tui.Config{
			Repository: repo,
			Prefix:     prefix,
			Remote:     remote,
			DryRun:     dryRun,
			NoPush:     noPush,
		}

		return tui.Run(cfg)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&prefix, "prefix", "p", "v", "Tag prefix")
	rootCmd.Flags().StringVarP(&remote, "remote", "r", "origin", "Git remote name")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Preview without making changes")
	rootCmd.Flags().BoolVar(&noPush, "no-push", false, "Create tag but don't push")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Show version information")
}

// Execute runs the root command
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return nil
}

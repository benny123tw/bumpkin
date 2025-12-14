package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/benny123tw/bumpkin/internal/config"
	"github.com/benny123tw/bumpkin/internal/conventional"
	"github.com/benny123tw/bumpkin/internal/executor"
	"github.com/benny123tw/bumpkin/internal/git"
	"github.com/benny123tw/bumpkin/internal/tui"
	"github.com/benny123tw/bumpkin/internal/version"
)

// Version information (set at build time)
var (
	AppVersion = "dev"
	BuildDate  = "unknown"
)

// Flag variables
var (
	// Bump type flags (mutually exclusive)
	flagPatch        bool
	flagMinor        bool
	flagMajor        bool
	flagSetVersion   string
	flagConventional bool
	flagAlpha        bool
	flagBeta         bool
	flagRC           bool
	flagRelease      bool

	// Behavior flags
	flagPrefix      string
	flagRemote      string
	flagConfig      string
	flagDryRun      bool
	flagNoPush      bool
	flagNoHooks     bool
	flagYes         bool
	flagJSON        bool
	flagShowVersion bool
)

// JSONOutput represents the JSON output format for non-interactive mode
type JSONOutput struct {
	Success         bool   `json:"success"`
	PreviousVersion string `json:"previous_version"`
	NewVersion      string `json:"new_version"`
	TagName         string `json:"tag_name"`
	CommitHash      string `json:"commit_hash"`
	TagCreated      bool   `json:"tag_created"`
	Pushed          bool   `json:"pushed"`
	DryRun          bool   `json:"dry_run"`
	Error           string `json:"error,omitempty"`
}

var rootCmd = &cobra.Command{
	Use:   "bumpkin",
	Short: "Semantic version tagger for git repositories",
	Long: `Bumpkin is a CLI tool that helps you tag commits by analyzing
conventional commit history and providing version options to select or customize.

Run without flags for interactive mode, or use flags for automation.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runRoot,
}

// NewRootCmd creates a new root command instance (for testing)
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bumpkin",
		Short: "Semantic version tagger for git repositories",
		Long: `Bumpkin is a CLI tool that helps you tag commits by analyzing
conventional commit history and providing version options to select or customize.

Run without flags for interactive mode, or use flags for automation.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runRoot,
	}

	addFlags(cmd)
	return cmd
}

func init() {
	addFlags(rootCmd)
}

func addFlags(cmd *cobra.Command) {
	// Bump type flags
	cmd.Flags().BoolVar(&flagPatch, "patch", false, "Bump patch version (x.y.Z)")
	cmd.Flags().BoolVar(&flagMinor, "minor", false, "Bump minor version (x.Y.0)")
	cmd.Flags().BoolVar(&flagMajor, "major", false, "Bump major version (X.0.0)")
	cmd.Flags().StringVar(&flagSetVersion, "set-version", "", "Set specific version")
	cmd.Flags().BoolVarP(
		&flagConventional,
		"conventional",
		"c",
		false,
		"Auto-detect bump type from conventional commits",
	)

	// Prerelease flags
	cmd.Flags().BoolVar(&flagAlpha, "alpha", false, "Bump to alpha prerelease")
	cmd.Flags().BoolVar(&flagBeta, "beta", false, "Bump to beta prerelease")
	cmd.Flags().BoolVar(&flagRC, "rc", false, "Bump to release candidate")
	cmd.Flags().BoolVar(&flagRelease, "release", false, "Promote prerelease to release")

	// Behavior flags
	cmd.Flags().StringVarP(&flagPrefix, "prefix", "p", "v", "Tag prefix")
	cmd.Flags().StringVarP(&flagRemote, "remote", "r", "origin", "Git remote name")
	cmd.Flags().StringVarP(&flagConfig, "config", "C", ".bumpkin.yml", "Config file path")
	cmd.Flags().BoolVarP(&flagDryRun, "dry-run", "d", false, "Preview without making changes")
	cmd.Flags().BoolVar(&flagNoPush, "no-push", false, "Create tag but don't push")
	cmd.Flags().BoolVar(&flagNoHooks, "no-hooks", false, "Skip hook execution")
	cmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "Skip confirmation in non-interactive mode")
	cmd.Flags().BoolVar(&flagJSON, "json", false, "Output result as JSON")
	cmd.Flags().BoolVar(&flagShowVersion, "show-version", false, "Show version information")

	// Keep -v as alias for show-version for backwards compatibility
	cmd.Flags().BoolP("version", "v", false, "Show version information")
	//nolint:errcheck // Best effort to mark hidden
	cmd.Flags().MarkHidden("version")
}

func runRoot(cmd *cobra.Command, _ []string) error {
	// Handle version flag
	showVer, _ := cmd.Flags().GetBool("version")
	if flagShowVersion || showVer {
		fmt.Fprintf(cmd.OutOrStdout(), "bumpkin %s (built %s)\n", AppVersion, BuildDate)
		return nil
	}

	// Determine if we're in non-interactive mode
	isNonInteractive := flagPatch || flagMinor || flagMajor || flagSetVersion != "" ||
		flagConventional || flagAlpha || flagBeta || flagRC || flagRelease

	// Open the repository from current directory
	repo, err := git.OpenFromCurrent()
	if err != nil {
		return handleErrorWithCode(cmd, ExitNotGitRepo, "not a git repository", err)
	}

	if isNonInteractive {
		return runNonInteractive(cmd, repo)
	}

	return runInteractive(repo)
}

func runNonInteractive(cmd *cobra.Command, repo *git.Repository) error {
	// Validate mutually exclusive flags
	bumpCount := 0
	if flagPatch {
		bumpCount++
	}
	if flagMinor {
		bumpCount++
	}
	if flagMajor {
		bumpCount++
	}
	if flagSetVersion != "" {
		bumpCount++
	}
	if flagConventional {
		bumpCount++
	}
	if flagAlpha {
		bumpCount++
	}
	if flagBeta {
		bumpCount++
	}
	if flagRC {
		bumpCount++
	}
	if flagRelease {
		bumpCount++
	}

	if bumpCount > 1 {
		return handleErrorWithCode(
			cmd,
			ExitInvalidArgs,
			"only one bump type flag can be specified",
			nil,
		)
	}

	// Determine bump type
	var bumpType version.BumpType
	var customVersion string

	switch {
	case flagAlpha:
		bumpType = version.BumpPrereleaseAlpha
	case flagBeta:
		bumpType = version.BumpPrereleaseBeta
	case flagRC:
		bumpType = version.BumpPrereleaseRC
	case flagRelease:
		bumpType = version.BumpRelease
	case flagPatch:
		bumpType = version.BumpPatch
	case flagMinor:
		bumpType = version.BumpMinor
	case flagMajor:
		bumpType = version.BumpMajor
	case flagSetVersion != "":
		bumpType = version.BumpCustom
		customVersion = flagSetVersion
	case flagConventional:
		// Analyze commits to determine bump type
		bumpType = analyzeConventionalCommits(repo)
	}

	// If not --yes, require confirmation (unless dry-run)
	if !flagYes && !flagDryRun {
		// Get current version for display
		latestTag, err := repo.LatestTag(flagPrefix)
		if err != nil {
			return handleError(cmd, err, "failed to get latest tag")
		}

		var prevVersion version.Version
		if latestTag == nil || latestTag.Version == nil {
			prevVersion = version.Zero()
		} else {
			prevVersion = *latestTag.Version
		}

		var newVersion version.Version
		if bumpType == version.BumpCustom {
			newVersion, err = version.Parse(customVersion)
			if err != nil {
				return handleError(cmd, err, "invalid version")
			}
		} else {
			newVersion = version.Bump(prevVersion, bumpType)
		}

		fmt.Fprintf(
			cmd.OutOrStdout(),
			"Will bump version: %s → %s\n",
			prevVersion.String(),
			newVersion.String(),
		)
		fmt.Fprintf(
			cmd.OutOrStdout(),
			"Use --yes to skip this confirmation, or run in interactive mode.\n",
		)
		return fmt.Errorf("confirmation required: use --yes flag to proceed")
	}

	// Load configuration
	cfg, _ := config.LoadFile(flagConfig)

	// Execute the bump
	req := executor.Request{
		Repository:    repo,
		BumpType:      bumpType,
		CustomVersion: customVersion,
		Prefix:        flagPrefix,
		Remote:        flagRemote,
		DryRun:        flagDryRun,
		NoPush:        flagNoPush,
		NoHooks:       flagNoHooks,
		PreTagHooks:   cfg.Hooks.PreTag,
		PostTagHooks:  cfg.Hooks.PostTag,
	}

	result, err := executor.Execute(context.Background(), req)
	if err != nil {
		return handleError(cmd, err, "bump failed")
	}

	// Output result
	if flagJSON {
		return outputJSON(cmd, result, nil)
	}

	return outputText(cmd, result)
}

func runInteractive(repo *git.Repository) error {
	cfg := tui.Config{
		Repository: repo,
		Prefix:     flagPrefix,
		Remote:     flagRemote,
		DryRun:     flagDryRun,
		NoPush:     flagNoPush,
	}

	return tui.Run(cfg)
}

func handleError(cmd *cobra.Command, err error, context string) error {
	return handleErrorWithCode(cmd, ExitGeneralError, context, err)
}

func handleErrorWithCode(cmd *cobra.Command, code int, message string, err error) error {
	exitErr := NewExitError(code, message, err)
	if flagJSON {
		//nolint:errcheck // Best effort output
		outputJSON(cmd, nil, exitErr)
	}
	return exitErr
}

func outputJSON(cmd *cobra.Command, result *executor.Result, err error) error {
	output := JSONOutput{
		Success: err == nil,
		DryRun:  flagDryRun,
	}

	if err != nil {
		output.Error = err.Error()
	}

	if result != nil {
		output.PreviousVersion = result.PreviousVersion
		output.NewVersion = result.NewVersion
		output.TagName = result.TagName
		output.CommitHash = result.CommitHash
		output.TagCreated = result.TagCreated
		output.Pushed = result.Pushed
	}

	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

func outputText(cmd *cobra.Command, result *executor.Result) error {
	out := cmd.OutOrStdout()

	if flagDryRun {
		fmt.Fprintln(out, "[DRY RUN]")
	}

	fmt.Fprintf(out, "Version: %s → %s\n", result.PreviousVersion, result.NewVersion)
	fmt.Fprintf(out, "Tag: %s\n", result.TagName)
	fmt.Fprintf(out, "Commit: %s\n", result.CommitHash[:7])

	if result.TagCreated {
		fmt.Fprintln(out, "Tag created: yes")
	} else {
		fmt.Fprintln(out, "Tag created: no (dry run)")
	}

	switch {
	case result.Pushed:
		fmt.Fprintln(out, "Pushed: yes")
	case flagNoPush:
		fmt.Fprintln(out, "Pushed: no (--no-push)")
	case flagDryRun:
		fmt.Fprintln(out, "Pushed: no (dry run)")
	default:
		fmt.Fprintln(out, "Pushed: no")
	}

	return nil
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(GetExitCode(err))
	}
}

// analyzeConventionalCommits analyzes commits and returns recommended bump type
func analyzeConventionalCommits(repo *git.Repository) version.BumpType {
	// Get latest tag
	latestTag, err := repo.LatestTag(flagPrefix)
	if err != nil {
		return version.BumpPatch // Default on error
	}

	var commits []*git.Commit
	if latestTag != nil {
		commits, err = repo.GetCommitsSinceTag(latestTag.Name)
	} else {
		commits, err = repo.GetAllCommits()
	}

	if err != nil || len(commits) == 0 {
		return version.BumpPatch // Default
	}

	// Extract commit messages
	var messages []string
	for _, c := range commits {
		messages = append(messages, c.Message)
	}

	// Analyze and return recommended bump
	analysis := conventional.AnalyzeCommits(messages)
	return analysis.RecommendedBump
}

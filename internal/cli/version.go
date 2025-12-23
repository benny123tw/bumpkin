package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// BuildInfo contains version information about the build.
// This is populated in main.go and passed to Execute().
type BuildInfo struct {
	GoVersion string `json:"goVersion"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
}

// String returns a formatted version string.
func (b BuildInfo) String() string {
	return fmt.Sprintf("bumpkin %s built with %s from %s on %s",
		b.Version, b.GoVersion, b.Commit, b.Date)
}

type versionCommand struct {
	cmd  *cobra.Command
	info BuildInfo
}

func newVersionCommand(info BuildInfo) *versionCommand {
	c := &versionCommand{info: info}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print the version, commit hash, and build date of bumpkin.",
		Args:  cobra.NoArgs,
		RunE:  c.execute,
	}

	c.cmd = versionCmd
	return c
}

func (c *versionCommand) execute(cmd *cobra.Command, _ []string) error {
	return printVersion(cmd.OutOrStdout(), c.info)
}

func printVersion(w io.Writer, info BuildInfo) error {
	_, err := fmt.Fprintln(w, info.String())
	return err
}

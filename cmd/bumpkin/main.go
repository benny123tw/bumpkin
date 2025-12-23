package main

import (
	"cmp"
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/benny123tw/bumpkin/internal/cli"
)

var (
	goVersion = "unknown"

	// Populated by goreleaser during build
	version = "unknown"
	commit  = "?"
	date    = ""
)

// main is the program entry point. It constructs runtime build metadata via createBuildInfo and invokes cli.Execute with that metadata.
func main() {
	info := createBuildInfo()
	cli.Execute(info)
}

// createBuildInfo returns version info from goreleaser variables or debug.ReadBuildInfo.
func createBuildInfo() cli.BuildInfo {
	info := cli.BuildInfo{
		Commit:    commit,
		Version:   version,
		GoVersion: goVersion,
		Date:      date,
	}

	buildInfo, available := debug.ReadBuildInfo()
	if !available {
		return info
	}

	info.GoVersion = buildInfo.GoVersion

	// If date is set (goreleaser build), use goreleaser values
	if date != "" {
		return info
	}

	// For go install / go build, extract version from build info
	info.Version = buildInfo.Main.Version

	// Strip "v" prefix only from proper semver versions (vX.Y.Z)
	if matched, _ := regexp.MatchString(`v\d+\.\d+\.\d+`, buildInfo.Main.Version); matched {
		info.Version = strings.TrimPrefix(buildInfo.Main.Version, "v")
	}

	// VCS info is only available with "go build", not "go install"
	var revision, modified string
	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.time":
			info.Date = setting.Value
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value
		}
	}

	info.Date = cmp.Or(info.Date, "(unknown)")

	// Format commit info with fallback to module sum for go install
	info.Commit = fmt.Sprintf("(%s, modified: %s, mod sum: %q)",
		cmp.Or(revision, "unknown"), cmp.Or(modified, "?"), buildInfo.Main.Sum)

	return info
}

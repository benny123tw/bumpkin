package conventional

import (
	"github.com/benny123tw/bumpkin/internal/version"
)

// AnalysisResult contains the result of analyzing commits
type AnalysisResult struct {
	RecommendedBump version.BumpType
	TypeCounts      map[string]int
	BreakingCount   int
	TotalCommits    int
}

// Commit types that trigger minor version bump
var minorTypes = map[string]bool{
	"feat": true,
	"perf": true, // Performance improvements add value
}

// AnalyzeCommits analyzes a list of commit messages and recommends a version bump
func AnalyzeCommits(messages []string) *AnalysisResult {
	result := &AnalysisResult{
		RecommendedBump: version.BumpPatch, // Default
		TypeCounts:      make(map[string]int),
		TotalCommits:    len(messages),
	}

	hasMinor := false
	hasMajor := false

	for _, msg := range messages {
		cc, err := ParseCommit(msg)
		if err != nil {
			continue
		}

		// Count by type (only for conventional commits)
		if cc.Type != "other" {
			result.TypeCounts[cc.Type]++
		}

		// Check for breaking changes
		if cc.IsBreaking {
			result.BreakingCount++
			hasMajor = true
		}

		// Check for minor bump types
		if minorTypes[cc.Type] {
			hasMinor = true
		}
	}

	// Determine recommendation based on priority: major > minor > patch
	if hasMajor {
		result.RecommendedBump = version.BumpMajor
	} else if hasMinor {
		result.RecommendedBump = version.BumpMinor
	}

	return result
}

// AnalyzeCommitMessages is a convenience function that takes git.Commit objects
// This allows integration with the git package
func AnalyzeCommitMessages(subjects []string) *AnalysisResult {
	return AnalyzeCommits(subjects)
}

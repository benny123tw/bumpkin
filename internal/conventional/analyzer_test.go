package conventional

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/benny123tw/bumpkin/internal/version"
)

// T077: Test feat commits recommend minor
func TestAnalyzeCommits_FeatRecommendsMinor(t *testing.T) {
	commits := []string{
		"feat: add new feature",
		"feat(api): add endpoint",
	}

	result := AnalyzeCommits(commits)
	assert.Equal(t, version.BumpMinor, result.RecommendedBump)
}

// T078: Test fix only commits recommend patch
func TestAnalyzeCommits_FixRecommendsPatch(t *testing.T) {
	commits := []string{
		"fix: resolve bug",
		"fix(ui): correct alignment",
	}

	result := AnalyzeCommits(commits)
	assert.Equal(t, version.BumpPatch, result.RecommendedBump)
}

// T079: Test breaking change recommends major
func TestAnalyzeCommits_BreakingRecommendsMajor(t *testing.T) {
	commits := []string{
		"feat!: remove deprecated API",
	}

	result := AnalyzeCommits(commits)
	assert.Equal(t, version.BumpMajor, result.RecommendedBump)
}

// T081: Test mixed commits use highest priority
func TestAnalyzeCommits_MixedUsesHighestPriority(t *testing.T) {
	tests := []struct {
		name     string
		commits  []string
		expected version.BumpType
	}{
		{
			name: "feat overrides fix",
			commits: []string{
				"fix: bug fix",
				"feat: new feature",
			},
			expected: version.BumpMinor,
		},
		{
			name: "breaking overrides feat",
			commits: []string{
				"feat: new feature",
				"fix!: breaking fix",
			},
			expected: version.BumpMajor,
		},
		{
			name: "breaking footer overrides all",
			commits: []string{
				"feat: feature\n\nBREAKING CHANGE: API changed",
				"fix: bug",
			},
			expected: version.BumpMajor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AnalyzeCommits(tt.commits)
			assert.Equal(t, tt.expected, result.RecommendedBump)
		})
	}
}

func TestAnalyzeCommits_EmptyCommits(t *testing.T) {
	result := AnalyzeCommits([]string{})
	assert.Equal(t, version.BumpPatch, result.RecommendedBump)
}

func TestAnalyzeCommits_NonConventionalOnly(t *testing.T) {
	commits := []string{
		"Update readme",
		"Merge branch 'main'",
		"WIP",
	}

	result := AnalyzeCommits(commits)
	// Default to patch when no conventional commits found
	assert.Equal(t, version.BumpPatch, result.RecommendedBump)
}

func TestAnalyzeCommits_CountsByType(t *testing.T) {
	commits := []string{
		"feat: feature 1",
		"feat: feature 2",
		"fix: bug 1",
		"docs: update docs",
		"chore: update deps",
	}

	result := AnalyzeCommits(commits)

	assert.Equal(t, 2, result.TypeCounts["feat"])
	assert.Equal(t, 1, result.TypeCounts["fix"])
	assert.Equal(t, 1, result.TypeCounts["docs"])
	assert.Equal(t, 1, result.TypeCounts["chore"])
}

func TestAnalyzeCommits_BreakingCount(t *testing.T) {
	commits := []string{
		"feat!: breaking feature",
		"fix: normal fix",
		"refactor!: breaking refactor",
	}

	result := AnalyzeCommits(commits)

	assert.Equal(t, 2, result.BreakingCount)
	assert.Equal(t, version.BumpMajor, result.RecommendedBump)
}

func TestAnalyzeCommits_PerfRecommendsMinor(t *testing.T) {
	// Performance improvements are like features - they add value
	commits := []string{
		"perf: optimize database queries",
	}

	result := AnalyzeCommits(commits)
	assert.Equal(t, version.BumpMinor, result.RecommendedBump)
}

func TestAnalyzeCommits_DocsRecommendsPatch(t *testing.T) {
	commits := []string{
		"docs: update API documentation",
	}

	result := AnalyzeCommits(commits)
	assert.Equal(t, version.BumpPatch, result.RecommendedBump)
}

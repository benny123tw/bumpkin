package cli

// testBuildInfo returns a BuildInfo for testing purposes
func testBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   "test",
		Commit:    "abc1234",
		Date:      "2024-01-01",
		GoVersion: "go1.21",
	}
}

package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/benny123tw/bumpkin/internal/executor"
)

func TestOutputJSON_IncludesPostPushWarnings(t *testing.T) {
	cmd := &cobra.Command{}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	result := &executor.Result{
		PreviousVersion:  "1.0.0",
		NewVersion:       "1.1.0",
		TagName:          "v1.1.0",
		CommitHash:       "abc1234",
		TagCreated:       true,
		Pushed:           true,
		PostPushWarnings: []string{"hook 'notify': exit 1", "hook 'changelog': timeout"},
	}

	require.NoError(t, outputJSON(cmd, result, nil))

	var out JSONOutput
	require.NoError(t, json.Unmarshal(buf.Bytes(), &out))
	assert.True(t, out.Success)
	assert.Equal(t, "v1.1.0", out.TagName)
	assert.Equal(t,
		[]string{"hook 'notify': exit 1", "hook 'changelog': timeout"},
		out.PostPushWarnings,
	)
}

func TestOutputJSON_OmitsEmptyPostPushWarnings(t *testing.T) {
	cmd := &cobra.Command{}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	result := &executor.Result{
		PreviousVersion: "1.0.0",
		NewVersion:      "1.1.0",
		TagName:         "v1.1.0",
		CommitHash:      "abc1234",
		TagCreated:      true,
		Pushed:          true,
	}

	require.NoError(t, outputJSON(cmd, result, nil))
	assert.NotContains(t, buf.String(), "post_push_warnings")
}

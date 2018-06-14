package github

import (
	"context"
	"testing"

	updater "github.com/ktr0731/go-updater"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ updater.Means = (*GitHubClient)(nil)
)

func TestGitHubReleaseMeans_LatestTag(t *testing.T) {
	// TODO: don't use real http request
	c, err := GitHubReleaseMeans("ktr0731", "evans", TarDecompresser)()
	require.NoError(t, err)
	_, err = c.LatestTag(context.Background())
	assert.NoError(t, err)
}

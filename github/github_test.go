package github

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
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
	c, err := GitHubReleaseMeans("ktr0731", "evans")()
	require.NoError(t, err)
	_, err = c.LatestTag(context.Background())
	assert.NoError(t, err)
}

func TestUpdateBinaryWithBackup(t *testing.T) {
	const defaultStr = "violet snow"
	setup := func(t *testing.T) *os.File {
		old, err := ioutil.TempFile("", "")
		require.NoError(t, err)

		_, err = io.WriteString(old, defaultStr)
		require.NoError(t, err)
		return old
	}

	t.Run("normal", func(t *testing.T) {
		f := setup(t)
		defer f.Close()

		in := strings.NewReader("violet-evergarden")
		err := updateBinaryWithBackup(f.Name(), in)
		assert.NoError(t, err)
	})

	t.Run("has error", func(t *testing.T) {
		f := setup(t)
		defer f.Close()

		ioCopy = func(_ io.Writer, _ io.Reader) (int64, error) {
			return 0, errors.New("an error")
		}

		in := strings.NewReader("violet-evergarden")
		err := updateBinaryWithBackup(f.Name(), in)
		assert.Error(t, err)

		b, err := ioutil.ReadFile(f.Name())
		assert.NoError(t, err)
		assert.Equal(t, defaultStr, string(b))
	})

	t.Run("has panic", func(t *testing.T) {
		f := setup(t)
		defer f.Close()

		ioCopy = func(_ io.Writer, _ io.Reader) (int64, error) {
			panic("a panic")
			return 0, nil
		}
		defer func(t *testing.T) {
			pErr := recover()
			assert.NotNil(t, pErr)
			b, err := ioutil.ReadFile(f.Name())
			assert.NoError(t, err)
			assert.Equal(t, defaultStr, string(b))
		}(t)
		in := strings.NewReader("violet-evergarden")
		err := updateBinaryWithBackup(f.Name(), in)
		assert.Error(t, err)
	})
}

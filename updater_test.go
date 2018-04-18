package updater

import (
	"testing"

	semver "github.com/ktr0731/go-semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mock MeansType = "mock"

var (
	v = semver.MustParse("0.1.0")
)

type mockMeans struct {
	t MeansType
	Means
}

func (m *mockMeans) Type() MeansType {
	// for Updatable(), returns GitHubRelease type
	if m.t == Empty {
		return GitHubRelease
	}
	return mock
}

func TestUpdater(t *testing.T) {
	t.Run("initialized updater has GitHub updater", func(t *testing.T) {
		u := New("ktr0731", "evans", v)
		_, ok := u.m[GitHubRelease]
		assert.True(t, ok)
	})

	t.Run("can register means", func(t *testing.T) {
		u := newUpdater("ktr0731", "evans", v)
		err := u.RegisterMeans(&mockMeans{})
		assert.NoError(t, err)

		// duplicated means
		err = u.RegisterMeans(&mockMeans{})
		assert.Error(t, err)
	})

	newMockUpdater := func(t *testing.T) *Updater {
		u := newUpdater("ktr0731", "evans", v)
		err := u.RegisterMeans(&mockMeans{})
		require.NoError(t, err)
		return u
	}

	t.Run("can register means with condition", func(t *testing.T) {
		u := newMockUpdater(t)
		u.UpdateIf = FoundPatchUpdate
	})
}

package updater

import (
	"context"
	"fmt"
	"testing"

	semver "github.com/ktr0731/go-semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	v = semver.MustParse("0.1.0")
)

type mockMeans struct {
	t      MeansType
	latest *semver.Version

	Means
}

func (m *mockMeans) LatestTag(_ context.Context) (*semver.Version, error) {
	return m.latest, nil
}

func (m *mockMeans) Type() MeansType {
	// for Updatable(), returns GitHubRelease type
	if m.t == Empty {
		return GitHubRelease
	}
	return m.t
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

}

func TestUpdater_Updatable(t *testing.T) {
	newMockUpdater := func(t *testing.T) (*Updater, *mockMeans) {
		u := newUpdater("ktr0731", "evans", v)
		m := &mockMeans{}
		err := u.RegisterMeans(m)
		require.NoError(t, err)
		return u, m
	}

	cases := []struct {
		cond            UpdateCondition
		current, latest string
		updatable       bool
	}{
		{FoundPatchUpdate, "0.1.5", "0.1.6", true},
		{FoundPatchUpdate, "0.1.5", "0.2.0", true},
		{FoundPatchUpdate, "0.1.5", "1.0.0", true},
		{FoundPatchUpdate, "0.1.5", "0.1.0", false},
		{FoundMinorUpdate, "0.1.5", "0.2.0", true},
		{FoundMinorUpdate, "0.1.5", "1.0.0", true},
		{FoundMinorUpdate, "0.1.5", "0.1.6", false},
		{FoundMajorUpdate, "0.1.5", "1.2.0", true},
		{FoundMajorUpdate, "0.1.5", "0.2.0", false},
	}

	for _, c := range cases {
		name := fmt.Sprintf("current: %s, latest: %s", c.current, c.latest)
		t.Run(name, func(t *testing.T) {
			u, m := newMockUpdater(t)
			u.UpdateIf = c.cond
			u.current = semver.MustParse(c.current)

			m.latest = semver.MustParse(c.latest)

			updatable, err := u.Updatable()
			require.NoError(t, err)
			assert.Equal(t, c.updatable, updatable)
		})
	}
}

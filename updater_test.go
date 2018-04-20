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
	latest       *semver.Version
	updateCalled bool

	Means
}

func (m *mockMeans) LatestTag(_ context.Context) (*semver.Version, error) {
	return m.latest, nil
}

func (m *mockMeans) Update(_ context.Context) (*semver.Version, error) {
	m.updateCalled = true
	return m.latest, nil
}

func TestUpdater_Update(t *testing.T) {
	newMockUpdater := func(t *testing.T) (*Updater, *mockMeans) {
		m := &mockMeans{}
		return New(v, m), m
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

			updatable, _, err := u.Updatable()
			require.NoError(t, err)
			assert.Equal(t, c.updatable, updatable)

			err = u.Update()
			require.NoError(t, err)

			if updatable {
				assert.True(t, m.updateCalled)
			} else {
				assert.False(t, m.updateCalled)
			}
		})
	}
}

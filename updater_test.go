package updater

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	v = version.Must(version.NewSemver("0.1.0"))
)

type mockMeans struct {
	latest       *version.Version
	updateCalled bool

	Means
}

func (m *mockMeans) LatestTag(_ context.Context) (*version.Version, error) {
	return m.latest, nil
}

func (m *mockMeans) Update(_ context.Context, _ *version.Version) error {
	m.updateCalled = true
	return nil
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
			u.current = version.Must(version.NewSemver(c.current))

			m.latest = version.Must(version.NewSemver(c.latest))

			updatable, _, err := u.Updatable(context.Background())
			require.NoError(t, err)
			assert.Equal(t, c.updatable, updatable)

			err = u.Update(context.Background())
			require.NoError(t, err)

			if updatable {
				assert.True(t, m.updateCalled)
			} else {
				assert.False(t, m.updateCalled)
			}
		})
	}
}

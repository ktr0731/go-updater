package updater

import (
	"testing"

	semver "github.com/ktr0731/go-semver"
	"github.com/stretchr/testify/assert"
)

func TestUpdater(t *testing.T) {
	u := New()

	t.Run("can register means with func literal condition", func(t *testing.T) {
		m := NewGitHubReleaseMeans("ktr0731", "evans")
		m.UpdateIf = func(current, latest *semver.Version) bool {
			return current.LessThan(latest)
		}
		err := u.RegisterMeans(m)
		assert.NoError(t, err)
	})
}

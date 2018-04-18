package updater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdater(t *testing.T) {
	u := New()

	t.Run("can register means with func literal condition", func(t *testing.T) {
		m := NewGitHubReleaseMeans("ktr0731", "evans")
		err := u.RegisterMeans(m)
		assert.NoError(t, err)
	})
}

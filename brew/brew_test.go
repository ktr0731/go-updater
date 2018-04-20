package brew

import (
	"context"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/require"
)

func TestHomeBrewFormula(t *testing.T) {
	m := NewHomeBrewMeans("ktr0731/evans", "evans")
	v, err := m.LatestTag(context.TODO())
	require.NoError(t, err)
	pp.Println(v)
}

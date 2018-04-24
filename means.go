package updater

import (
	"context"

	semver "github.com/ktr0731/go-semver"
)

type MeansType string

// Means manages methods related to specified update means
// for example, fetches the latest tag, update binary, or
// check whether the software is installed by this.
type Means interface {
	LatestTag(context.Context) (*semver.Version, error)
	Update(context.Context, *semver.Version) error

	Installed() bool

	CommandText(*semver.Version) string

	Type() MeansType
}

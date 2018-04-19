package updater

import (
	"context"

	semver "github.com/ktr0731/go-semver"
)

type MeansType string

const (
	Empty                   = ""
	GitHubRelease MeansType = "github-release"
	GoGet                   = "go-get"
	HomeBrew                = "homebrew"
)

type Means interface {
	LatestTag(context.Context) (*semver.Version, error)
	Update(context.Context) (*semver.Version, error)

	Type() MeansType

	Installed() bool
}

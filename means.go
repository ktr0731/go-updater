package updater

import (
	"context"
)

type MeansType string

const (
	GitHubRelease MeansType = "github-release"
	GoGet                   = "go-get"
	HomeBrew                = "homebrew"
)

type Means interface {
	LatestTag(context.Context) (string, error)
	Update(context.Context) (string, error)

	Type() MeansType

	Installed() bool
}

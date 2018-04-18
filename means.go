package updater

import (
	"context"

	semver "github.com/ktr0731/go-semver"
)

type MeansType int

const (
	GitHubRelease = "gh-release"
	GoGet         = "go-get"
	HomeBrew      = "brew"
)

type UpdateCondition func(*semver.Version, *semver.Version) bool

type Means struct {
	typ      MeansType
	UpdateIf UpdateCondition
	client   Client
}

func NewMeans(typ MeansType, client Client) *Means {
	return &Means{
		typ:      typ,
		UpdateIf: FoundMinorUpdate,
		client:   client,
	}
}

func (m *Means) LatestTag(ctx context.Context) (string, error) {
	return m.client.FetchLatestTag(ctx)
}

func (m *Means) Update() error {
	return m.client.Update()
}

var (
	FoundMajorUpdate = func(current, latest *semver.Version) bool {
		return current.LessThan(latest) && current.Major < latest.Major
	}

	FoundMinorUpdate = func(current, latest *semver.Version) bool {
		return current.LessThan(latest) && current.Minor < latest.Minor
	}

	FoundPatchUpdate = func(current, latest *semver.Version) bool {
		return current.LessThan(latest) && current.Patch < latest.Patch
	}
)

func NewGitHubReleaseMeans(owner, repo string) *Means {
	return NewMeans(GitHubRelease, newGitHubClient(owner, repo))
}

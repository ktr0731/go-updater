package updater

import (
	"context"
	"sync"

	semver "github.com/ktr0731/go-semver"
	"github.com/pkg/errors"
)

type Updater struct {
	UpdateIf UpdateCondition

	current *semver.Version

	mu sync.Mutex
	m  map[MeansType]Means
}

// New receives repository info for fetch release tags
// TODO: other hosting services, like BitBucket
func New(owner, repo string, current *semver.Version) *Updater {
	u := newUpdater(owner, repo, current)
	if err := u.RegisterMeans(newGitHubReleaseMeans(owner, repo)); err != nil {
		panic(err)
	}
	return u
}

func newUpdater(owner, repo string, current *semver.Version) *Updater {
	return &Updater{
		UpdateIf: FoundMinorUpdate,
		current:  current,
		m:        map[MeansType]Means{},
	}
}

func (u *Updater) RegisterMeans(m Means) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if _, ok := u.m[m.Type()]; ok {
		return errors.Errorf("duplicated means: %s", m.Type())
	}
	u.m[m.Type()] = m
	return nil
}

func (u *Updater) UpdateBy(typ MeansType) error {
	m, ok := u.m[typ]
	if !ok {
		return errors.Errorf("no such means: %s", typ)
	}
	_, err := m.Update(context.TODO())
	return err
}

func (u *Updater) Updatable() bool {
	m := u.m[GitHubRelease]
	m.LatestTag(context.Background())
	return true
}

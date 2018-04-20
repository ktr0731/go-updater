package updater

import (
	"context"
	"io"

	semver "github.com/ktr0731/go-semver"
)

type Updater struct {
	UpdateIf UpdateCondition

	current *semver.Version

	m Means
}

// New receives repository info for fetch release tags
// TODO: other hosting services, like BitBucket
func New(current *semver.Version, m Means) *Updater {
	return &Updater{
		UpdateIf: FoundMinorUpdate,
		current:  current,
		m:        m,
	}
}

func (u *Updater) Update() error {
	_, err := u.m.Update(context.TODO())
	return err
}

func (u *Updater) Updatable() (bool, error) {
	latest, err := u.m.LatestTag(context.Background())
	if err != nil {
		return false, err
	}
	return u.UpdateIf(u.current, latest), nil
}

func (u *Updater) PrintInstruction(typ MeansType, w io.Writer, v *semver.Version) error {
	_, err := io.WriteString(w, u.m.CommandText(v))
	return err
}

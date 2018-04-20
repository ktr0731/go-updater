package updater

import (
	"context"
	"io"

	semver "github.com/ktr0731/go-semver"
)

type Updater struct {
	UpdateIf UpdateCondition

	current      *semver.Version
	cachedLatest *semver.Version // to be set when call Updatable()

	m Means
}

func New(current *semver.Version, m Means) *Updater {
	return &Updater{
		UpdateIf: FoundMinorUpdate,
		current:  current,
		m:        m,
	}
}

func (u *Updater) Update() error {
	updatable, _, err := u.Updatable()
	if err != nil {
		return err
	}
	if updatable {
		_, err = u.m.Update(context.TODO())
	}
	return err
}

func (u *Updater) Updatable() (bool, *semver.Version, error) {
	var err error
	if u.cachedLatest == nil {
		u.cachedLatest, err = u.m.LatestTag(context.Background())
		if err != nil {
			return false, nil, err
		}
	}
	return u.UpdateIf(u.current, u.cachedLatest), u.cachedLatest, nil
}

func (u *Updater) PrintInstruction(w io.Writer, v *semver.Version) error {
	_, err := io.WriteString(w, u.m.CommandText(v))
	return err
}

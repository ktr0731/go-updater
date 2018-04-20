package updater

import (
	"context"
	"io"

	semver "github.com/ktr0731/go-semver"
	"github.com/pkg/errors"
)

// Updater updates the binary using with UpdateCondition and Means.
// updating can be executed by m when UpdateIf is true.
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

// Update updates the binary if Updatable()
func (u *Updater) Update() error {
	updatable, _, err := u.Updatable()
	if err != nil {
		return errors.Wrap(err, "failed to check the latest release")
	}
	if updatable {
		err = u.m.Update(context.TODO(), u.cachedLatest)
	}
	if err != nil {
		return errors.Wrap(err, "failed to update the binary")
	}
	return nil
}

func (u *Updater) Updatable() (bool, *semver.Version, error) {
	var err error
	if u.cachedLatest == nil {
		u.cachedLatest, err = u.m.LatestTag(context.Background())
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to cache the latest release version")
		}
	}
	return u.UpdateIf(u.current, u.cachedLatest), u.cachedLatest, nil
}

func (u *Updater) PrintInstruction(w io.Writer, v *semver.Version) error {
	if _, err := io.WriteString(w, u.m.CommandText(v)); err != nil {
		return errors.Wrap(err, "failed to write command text")
	}
	return nil
}

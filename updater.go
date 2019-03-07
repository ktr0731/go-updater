package updater

import (
	"context"
	"io"

	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

// Updater updates the binary using with UpdateCondition and Means.
// updating can be executed by m when UpdateIf is true.
type Updater struct {
	UpdateIf UpdateCondition

	current      *version.Version
	cachedLatest *version.Version // to be set when call Updatable()

	m Means
}

func New(current *version.Version, m Means) *Updater {
	return &Updater{
		UpdateIf: FoundMinorUpdate,
		current:  current,
		m:        m,
	}
}

// Update updates the binary if Updatable()
func (u *Updater) Update(ctx context.Context) error {
	updatable, _, err := u.Updatable(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to check the latest release")
	}
	if updatable {
		err = u.m.Update(ctx, u.cachedLatest)
	}
	if err != nil {
		return errors.Wrap(err, "failed to update the binary")
	}
	return nil
}

func (u *Updater) Updatable(ctx context.Context) (bool, *version.Version, error) {
	var err error
	if u.cachedLatest == nil {
		u.cachedLatest, err = u.m.LatestTag(ctx)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to cache the latest release version")
		}
	}
	return u.UpdateIf(u.current, u.cachedLatest), u.cachedLatest, nil
}

func (u *Updater) PrintInstruction(w io.Writer, v *version.Version) error {
	if _, err := io.WriteString(w, u.m.CommandText(v)); err != nil {
		return errors.Wrap(err, "failed to write command text")
	}
	return nil
}

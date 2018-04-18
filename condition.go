package updater

import semver "github.com/ktr0731/go-semver"

type UpdateCondition func(*semver.Version, *semver.Version) bool

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

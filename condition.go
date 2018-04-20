package updater

import semver "github.com/ktr0731/go-semver"

// UpdateCondition is the condition for update the binary.
// when the binary will be updated, always go-updater checks the condition specified by this.
type UpdateCondition func(*semver.Version, *semver.Version) bool

var (
	FoundMajorUpdate UpdateCondition = func(current, latest *semver.Version) bool {
		return current.LessThan(latest) && current.Major < latest.Major
	}

	FoundMinorUpdate UpdateCondition = func(current, latest *semver.Version) bool {
		return current.LessThan(latest) && (current.Major < latest.Major || current.Minor < latest.Minor)
	}

	FoundPatchUpdate UpdateCondition = func(current, latest *semver.Version) bool {
		return current.LessThan(latest)
	}
)

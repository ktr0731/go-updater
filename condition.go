package updater

import "github.com/hashicorp/go-version"

// UpdateCondition is the condition for update the binary.
// when the binary will be updated, always go-updater checks the condition specified by this.
type UpdateCondition func(*version.Version, *version.Version) bool

var (
	FoundMajorUpdate UpdateCondition = func(current, latest *version.Version) bool {
		return current.LessThan(latest) && current.Segments()[0] < latest.Segments()[0]
	}

	FoundMinorUpdate UpdateCondition = func(current, latest *version.Version) bool {
		return current.LessThan(latest) && (current.Segments()[0] < latest.Segments()[0] || current.Segments()[1] < latest.Segments()[1])
	}

	FoundPatchUpdate UpdateCondition = func(current, latest *version.Version) bool {
		return current.LessThan(latest)
	}
)

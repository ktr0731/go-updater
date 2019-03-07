package updater

import "errors"

var (
	// ErrUnavailable is used when the means unavailable.
	// For example, Homebrew cannot be used from Linux.
	//
	// The condition is checked at the constructor like github.New(...).
	// Clients which use go-updater must check this error.
	ErrUnavailable = errors.New("unavailable means")
)

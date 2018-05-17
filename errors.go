package updater

import "errors"

var (
	// ErrUnavailable is used when the means unavailable
	// for example, Homebrew cannot be used from Linux
	//
	// the condition is checked at the constructor like github.New(...)
	// client which using go-updater must check this error
	ErrUnavailable = errors.New("unavailable means")
)

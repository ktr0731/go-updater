package updater

import (
	"context"
)

const (
	owner = "ktr0731"
	repo  = "evans"
)

type Client interface {
	FetchLatestTag(context.Context) (string, error)
	Installed() bool
}

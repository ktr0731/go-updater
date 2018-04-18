package updater

import (
	"context"

	semver "github.com/ktr0731/go-semver"
)

type UpdateChecker struct {
	client  Client
	version semver.Version
}

func NewUpdateChecker(version semver.Version, client Client) *UpdateChecker {
	return &UpdateChecker{
		client:  client,
		version: version,
	}
}

type ReleaseTag struct {
	LatestVersion   string
	CurrentIsLatest bool
}

func (u *UpdateChecker) Check(ctx context.Context) (*ReleaseTag, error) {
	tag, err := u.client.FetchLatestTag(ctx)
	if err != nil {
		return nil, err
	}
	return &ReleaseTag{
		LatestVersion:   tag,
		CurrentIsLatest: tag == u.version,
	}, nil
}

package updater

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	semver "github.com/ktr0731/go-semver"
)

type gitHubClient struct {
	client      *github.Client
	owner, repo string
}

func newGitHubReleaseMeans(owner, repo string) Means {
	return &gitHubClient{
		client: github.NewClient(nil),
	}
}

func (c *gitHubClient) LatestTag(ctx context.Context) (*semver.Version, error) {
	r, _, err := c.client.Repositories.GetLatestRelease(ctx, c.owner, c.repo)
	if err != nil {
		return nil, err
	}
	return semver.MustParse(r.GetTagName()), nil
}

func (c *gitHubClient) Update(ctx context.Context) (*semver.Version, error) {
	panic("not implemented yet")
	return nil, nil
}

func (c *gitHubClient) Installed() bool {
	return false
}

func (c *gitHubClient) CommandText(v *semver.Version) string {
	return fmt.Sprintf(
		"curl -sL https://github.com/%s/%s/releases/download/%s/%s_{OS}_{ARCH}.tar.gz | tar xf -",
		c.owner,
		c.repo,
		v,
		c.repo,
	)
}

func (c *gitHubClient) Type() MeansType {
	return GitHubRelease
}

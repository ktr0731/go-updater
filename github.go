package updater

import (
	"context"

	"github.com/google/go-github/github"
)

type gitHubClient struct {
	client      *github.Client
	owner, repo string
}

func NewGitHubReleaseMeans(owner, repo string) Means {
	return &gitHubClient{
		client: github.NewClient(nil),
	}
}

func (c *gitHubClient) LatestTag(ctx context.Context) (string, error) {
	r, _, err := c.client.Repositories.GetLatestRelease(ctx, c.owner, c.repo)
	if err != nil {
		return "", err
	}
	return r.GetTagName(), nil
}

func (c *gitHubClient) Update(ctx context.Context) (string, error) {
	return "", nil
}

func (c *gitHubClient) Installed() bool {
	return false
}

func (c *gitHubClient) Type() MeansType {
	return GitHubRelease
}

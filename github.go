package updater

import (
	"context"

	"github.com/google/go-github/github"
)

type gitHubClient struct {
	client      *github.Client
	owner, repo string
}

func newGitHubClient(owner, repo string) Client {
	return &gitHubClient{
		client: github.NewClient(nil),
	}
}

func (c *gitHubClient) FetchLatestTag(ctx context.Context) (string, error) {
	r, _, err := c.client.Repositories.GetLatestRelease(ctx, c.owner, c.repo)
	if err != nil {
		return "", err
	}
	return r.GetTagName(), nil
}

package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/google/go-github/github"
	"github.com/hashicorp/go-version"
	update "github.com/inconshreveable/go-update"
	updater "github.com/ktr0731/go-updater"
	"github.com/pkg/errors"
)

var (
	// IsGitHubReleasedBinary means the app is installed by GitHub Releases.
	// The app must change this value to true to use GitHub means.
	IsGitHubReleasedBinary bool

	releaseURLFormat = fmt.Sprintf(
		"https://github.com/%%s/%%s/releases/download/%%s/%%s_%s_%s.tar.gz",
		runtime.GOOS,
		runtime.GOARCH,
	)
)

const MeansTypeGitHubRelease updater.MeansType = "github-release"

type GitHubClient struct {
	client       *github.Client
	owner, repo  string
	decompresser Decompresser
}

// GitHubReleaseMeans returns updater.MeansBuilder for GitHubReleases.
// if dec is nil, DefaultDecompresser is used for extract binary from the compressed file.
func GitHubReleaseMeans(owner, repo string, dec Decompresser) updater.MeansBuilder {
	c := &GitHubClient{
		client:       github.NewClient(nil),
		owner:        owner,
		repo:         repo,
		decompresser: dec,
	}
	// if didn't set Decompresser, use default compresser (tar.gz)
	if dec == nil {
		c.decompresser = DefaultDecompresser
	}
	return func() (updater.Means, error) {
		return c, nil
	}
}

func (c *GitHubClient) LatestTag(ctx context.Context) (*version.Version, error) {
	r, _, err := c.client.Repositories.GetLatestRelease(ctx, c.owner, c.repo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the latest release from GitHub")
	}
	return version.Must(version.NewSemver(r.GetTagName())), nil
}

func (c *GitHubClient) Update(ctx context.Context, latest *version.Version) error {
	p, err := exec.LookPath(os.Args[0])
	if err != nil {
		return errors.Wrap(err, "failed to lookup the command, are you installed?")
	}

	req, err := http.NewRequest(http.MethodGet, c.releaseURL(latest), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create new http request to get latest GitHub release")
	}
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to get release binary")
	}
	defer res.Body.Close()

	dec, err := c.decompresser(res.Body)
	if err != nil && err != io.EOF {
		return errors.Wrap(err, "failed to decompress downloaded release file")
	}

	return update.Apply(dec, update.Options{
		TargetPath: p,
	})
}

func (c *GitHubClient) Installed(_ context.Context) bool {
	return IsGitHubReleasedBinary
}

func (c *GitHubClient) CommandText(v *version.Version) string {
	return fmt.Sprintf("curl -sL %s | tar xf -\n", c.releaseURL(v))
}

func (c *GitHubClient) Type() updater.MeansType {
	return MeansTypeGitHubRelease
}

func (c *GitHubClient) releaseURL(v *version.Version) string {
	return fmt.Sprintf(releaseURLFormat, c.owner, c.repo, v, c.repo)
}

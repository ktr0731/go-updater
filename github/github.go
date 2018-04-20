package github

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/google/go-github/github"
	"github.com/k0kubun/pp"
	semver "github.com/ktr0731/go-semver"
	updater "github.com/ktr0731/go-updater"
)

var (
	isGitHubReleasedBinary string

	releaseURLFormat = fmt.Sprintf("https://github.com/%%s/%%s/releases/download/%%s/%%s_%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)
)

type GitHubClient struct {
	client              *github.Client
	owner, repo         string
	DecompresserBuilder func(io.Reader) (io.ReadCloser, error)
}

func NewGitHubReleaseMeans(owner, repo string) updater.Means {
	c := &GitHubClient{
		client: github.NewClient(nil),
		owner:  owner,
		repo:   repo,
	}
	// if didn't set Decompresser, use default compresser (tar.gz)
	if c.DecompresserBuilder == nil {
		c.DecompresserBuilder = func(w io.Reader) (io.ReadCloser, error) {
			// return gzip.NewReader(tar.NewReader(w))
			// return ioutil.NopCloser(tar.NewReader(w)), nil
			return ioutil.NopCloser(w), nil
		}
	}

	return c
}

func (c *GitHubClient) LatestTag(ctx context.Context) (*semver.Version, error) {
	r, _, err := c.client.Repositories.GetLatestRelease(ctx, c.owner, c.repo)
	if err != nil {
		return nil, err
	}
	return semver.MustParse(r.GetTagName()), nil
}

func (c *GitHubClient) Update(ctx context.Context) (*semver.Version, error) {
	p, err := exec.LookPath(os.Args[0])
	if err != nil {
		return nil, err
	}

	latest, err := c.LatestTag(ctx)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(c.releaseURL(latest))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// TODO: rollback
	pp.Println(p)
	f, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec, err := c.DecompresserBuilder(res.Body)
	if err != nil && err != io.EOF {
		return nil, err
	}
	_, err = io.Copy(f, dec)

	return latest, err
}

func (c *GitHubClient) Installed() bool {
	// return isGitHubReleasedBinary != ""
	return true
}

func (c *GitHubClient) CommandText(v *semver.Version) string {
	return fmt.Sprintf("curl -sL %s | tar xf -", c.releaseURL(v))
}

func (c *GitHubClient) releaseURL(v *semver.Version) string {
	return fmt.Sprintf(releaseURLFormat, c.owner, c.repo, v, c.repo)
}

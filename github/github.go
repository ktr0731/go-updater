package github

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/google/go-github/github"
	semver "github.com/ktr0731/go-semver"
	updater "github.com/ktr0731/go-updater"
)

var (
	isGitHubReleasedBinary string

	releaseURLFormat = fmt.Sprintf(
		"https://github.com/%%s/%%s/releases/download/%%s/%%s_%s_%s.tar.gz",
		runtime.GOOS,
		runtime.GOARCH,
	)
)

type GitHubClient struct {
	client       *github.Client
	owner, repo  string
	Decompresser updater.Decompresser
}

func NewGitHubReleaseMeans(owner, repo string) *GitHubClient {
	c := &GitHubClient{
		client: github.NewClient(nil),
		owner:  owner,
		repo:   repo,
	}
	// if didn't set Decompresser, use default compresser (tar.gz)
	if c.Decompresser == nil {
		c.Decompresser = updater.DefaultDecompresser
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

	dec, err := c.Decompresser(res.Body)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return latest, updateBinaryWithBackup(p, dec)
}

func (c *GitHubClient) Installed() bool {
	return isGitHubReleasedBinary != ""
}

func (c *GitHubClient) CommandText(v *semver.Version) string {
	return fmt.Sprintf("curl -sL %s | tar xf -\n", c.releaseURL(v))
}

func (c *GitHubClient) releaseURL(v *semver.Version) string {
	return fmt.Sprintf(releaseURLFormat, c.owner, c.repo, v, c.repo)
}

// for testing
var ioCopy = io.Copy

func updateBinaryWithBackup(p string, in io.Reader) error {
	tmp := &bytes.Buffer{}

	f, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	// backup current binary
	if _, err := io.Copy(tmp, f); err != nil {
		return err
	}

	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	// rollback
	defer func() {
		if err := recover(); err != nil {
			io.Copy(f, tmp)
			panic(err)
		}
		if err != nil {
			io.Copy(f, tmp)
		}
	}()

	_, err = ioCopy(f, in)
	return err
}

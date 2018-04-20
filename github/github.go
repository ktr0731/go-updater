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
	"github.com/pkg/errors"
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
	Decompresser Decompresser
}

func NewGitHubReleaseMeans(owner, repo string) *GitHubClient {
	c := &GitHubClient{
		client: github.NewClient(nil),
		owner:  owner,
		repo:   repo,
	}
	// if didn't set Decompresser, use default compresser (tar.gz)
	if c.Decompresser == nil {
		c.Decompresser = DefaultDecompresser
	}
	return c
}

func (c *GitHubClient) LatestTag(ctx context.Context) (*semver.Version, error) {
	r, _, err := c.client.Repositories.GetLatestRelease(ctx, c.owner, c.repo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the latest release from GitHub")
	}
	return semver.MustParse(r.GetTagName()), nil
}

func (c *GitHubClient) Update(ctx context.Context) (*semver.Version, error) {
	p, err := exec.LookPath(os.Args[0])
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup the command, are you installed?")
	}

	latest, err := c.LatestTag(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest tag")
	}

	res, err := http.Get(c.releaseURL(latest))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get release binary")
	}
	defer res.Body.Close()

	dec, err := c.Decompresser(res.Body)
	if err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "failed to decompress downloaded release file")
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
		return errors.Wrap(err, "failed to create the executable")
	}
	defer f.Close()

	// backup current binary
	if _, err := io.Copy(tmp, f); err != nil {
		return errors.Wrap(err, "failed to create backup")
	}

	if err := f.Truncate(0); err != nil {
		return errors.Wrap(err, "failed to truncate old executable content")
	}
	if _, err := f.Seek(0, 0); err != nil {
		return errors.Wrap(err, "failed to seek to head")
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

	if _, err = ioCopy(f, in); err != nil {
		return errors.Wrap(err, "failed to write new binary to file")
	}
	return nil
}

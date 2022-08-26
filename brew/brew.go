package brew

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/hashicorp/go-version"
	updater "github.com/ktr0731/go-updater"
	pipeline "github.com/mattn/go-pipeline"
	"github.com/pkg/errors"
)

const MeansTypeHomebrew updater.MeansType = "homebrew"

type HomebrewClient struct {
	formula, name string
	cmdPath       string
}

func HomebrewMeans(formula, name string) updater.MeansBuilder {
	return func() (updater.Means, error) {
		if runtime.GOOS != "darwin" {
			return nil, updater.ErrUnavailable
		}
		p, err := exec.LookPath("brew")
		if err != nil {
			return nil, updater.ErrUnavailable
		}
		return &HomebrewClient{
			formula: formula,
			name:    name,
			cmdPath: p,
		}, nil
	}
}

// update instruction
//   1. update formula by "brew tap <formula>" if formula is not empty
//   2. get latest version by "brew info <formula>"
func (c *HomebrewClient) LatestTag(ctx context.Context) (*version.Version, error) {
	// update formula
	if c.formula != "" {
		err := exec.CommandContext(ctx, c.cmdPath, "tap", c.formula).Run()
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed to update Homebrew formula: %s", c.formula)
		}
	}

	// get latest version
	out, err := pipeline.Output(
		[]string{c.cmdPath, "info", c.getFullName()},
		[]string{"head", "-1"},
		[]string{"grep", "-o", "-E", "([0-9]+\.){1}[0-9]+(\.[0-9]+)?"},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get the latest info: %s", c.getFullName())
	}

	return version.NewSemver(strings.TrimSpace(string(out)))
}

func (c *HomebrewClient) Update(ctx context.Context, _ *version.Version) error {
	cmd := exec.CommandContext(ctx, c.cmdPath, "upgrade", c.getFullName())
	eo := new(bytes.Buffer)
	cmd.Stderr = eo
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to upgrade the binary by Homebrew: %s", eo.String())
	}
	return nil
}

func (c *HomebrewClient) Installed(ctx context.Context) bool {
	out, err := exec.CommandContext(ctx, c.cmdPath, "list", c.getFullName()).Output()
	if err != nil {
		return false
	}
	return len(out) != 0
}

func (c *HomebrewClient) CommandText(v *version.Version) string {
	return fmt.Sprintf("brew upgrade %s\n", c.getFullName())
}

func (c *HomebrewClient) Type() updater.MeansType {
	return MeansTypeHomebrew
}

func (c *HomebrewClient) getFullName() string {
	if c.formula != "" {
		return fmt.Sprintf("%s/%s", c.formula, c.name)
	}
	return c.name
}

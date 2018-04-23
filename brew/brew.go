package brew

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	semver "github.com/ktr0731/go-semver"
	updater "github.com/ktr0731/go-updater"
	pipeline "github.com/mattn/go-pipeline"
	"github.com/pkg/errors"
)

type HomeBrewClient struct {
	formula, name string
	cmdPath       string
}

func HomeBrewMeans(formula, name string) updater.MeansBuilder {
	return func() (updater.Means, error) {
		if runtime.GOOS != "darwin" {
			return nil, updater.ErrUnavailable
		}
		p, err := exec.LookPath("brew")
		if err != nil {
			return nil, updater.ErrUnavailable
		}
		return &HomeBrewClient{
			formula: formula,
			name:    name,
			cmdPath: p,
		}, nil
	}
}

// update instruction
//   1. update formula by "brew tap <formula>" if formula is not empty
//   2. get latest version by "brew info <formula>"
func (c *HomeBrewClient) LatestTag(ctx context.Context) (*semver.Version, error) {
	// update formula
	if c.formula != "" {
		err := exec.Command(c.cmdPath, "tap", c.formula).Run()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to update formula: %s", c.formula)
		}
	}

	// get latest version
	out, err := pipeline.Output(
		[]string{c.cmdPath, "info", c.getFullName()},
		[]string{"head", "-1"},
		[]string{"awk", "{ print $3 }"},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get the latest info: %s", c.getFullName())
	}

	return semver.MustParse(strings.TrimSpace(string(out))), nil
}

func (c *HomeBrewClient) Update(ctx context.Context, _ *semver.Version) error {
	if err := exec.Command(c.cmdPath, "upgrade", c.getFullName()).Run(); err != nil {
		return errors.Wrap(err, "failed to upgrade the binary")
	}
	return nil
}

func (c *HomeBrewClient) Installed() bool {
	out, err := exec.Command(c.cmdPath, "list", c.getFullName()).Output()
	if err != nil {
		return false
	}
	return len(out) != 0
}

func (c *HomeBrewClient) CommandText(v *semver.Version) string {
	return fmt.Sprintf("brew upgrade %s\n", c.getFullName())
}

func (c *HomeBrewClient) getFullName() string {
	if c.formula != "" {
		return fmt.Sprintf("%s/%s", c.formula, c.name)
	}
	return c.name
}

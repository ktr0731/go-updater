package mod

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime/debug"

	"github.com/hashicorp/go-version"
	"github.com/ktr0731/go-updater"
	"github.com/pkg/errors"
)

const MeansTypeGoModules updater.MeansType = "go-modules"

type GoModulesClient struct {
	path string
}

// GoModulesMeans returns updater.MeansBuilder for Go modules.
// if dec is nil, DefaultDecompresser is used for extract binary from the compressed file.
func GoModulesMeans(path string) updater.MeansBuilder {
	c := &GoModulesClient{path: path}

	return func() (updater.Means, error) {
		return c, nil
	}
}

func (c *GoModulesClient) LatestTag(ctx context.Context) (*version.Version, error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, errors.New("not versioned by Go modules")
	}
	return version.Must(version.NewSemver(info.Main.Version)), nil
}

func (c *GoModulesClient) Update(ctx context.Context, latest *version.Version) error {
	var buf bytes.Buffer
	cmd := exec.CommandContext(ctx, "go", "install", fmt.Sprintf("%s@latest", c.path))
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to run command: %s", buf.String())
	}

	return nil
}

func (c *GoModulesClient) Installed(ctx context.Context) bool {
	_, err := c.LatestTag(ctx)
	return err == nil
}

func (c *GoModulesClient) CommandText(v *version.Version) string {
	return fmt.Sprintf("go install %s@latest\n", c.path)
}

func (c *GoModulesClient) Type() updater.MeansType {
	return MeansTypeGoModules
}

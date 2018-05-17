# go-updater
[![GoDoc](https://godoc.org/github.com/ktr0731/go-updater?status.svg)](https://godoc.org/github.com/ktr0731/go-updater)
[![CircleCI](https://circleci.com/gh/ktr0731/go-updater.svg?style=svg)](https://circleci.com/gh/ktr0731/go-updater)  

## Usage
a software which using go-updater
``` go
package main

import (
	"fmt"

	semver "github.com/ktr0731/go-semver"
	updater "github.com/ktr0731/go-updater"
)

var version = semver.MustParse("0.1.0")

func main() {
	// determine what use means for update this software
	// in this example, use GitHub release
	u := updater.New(version, github.NewGitHubReleaseMeans("ktr0731", "evans"))

	// in default, update if minor update found
	u.UpdateIf = updater.FoundPatchUpdate

	updatable, latest, _ := u.Updatable()
	if updatable {
		_ = u.Update()
		fmt.Println("update from %s to %s", version, latest)
	} else {
		fmt.Println("already latest")
	}
}
```

more advanced step.
``` go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	semver "github.com/ktr0731/go-semver"
	updater "github.com/ktr0731/go-updater"
	"github.com/ktr0731/go-updater/brew"
	"github.com/ktr0731/go-updater/github"
)

var version = semver.MustParse("0.1.0")

type Config struct {
	UpdateBy string `json:"updateBy"`
}

func main() {
	f, _ := os.Open("config.json")
	defer f.Close()

	cfg := Config{}
	json.NewDecoder(f).Decode(&cfg)

	var m updater.Means
	switch cfg.UpdateBy {
	case "brew":
		m = brew.NewHomebrewMeans("ktr0731/evans", "evans")
	case "gh-release":
		m = github.NewGitHubReleaseMeans("ktr0731", "evans")
	default:
		panic("unknown means")
	}

	u := updater.New(version, m)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var latest *semver.Version
	go func() {
		defer cancel()
		var updatable bool
		updatable, latest, _ = u.Updatable()
		if updatable {
			_ = u.Update()
		}
	}()

	// do something

	// always call cancel.
	// because if updating is WIP, need to cancel for stop application immediately.
	// if updating is finished, cancel() do nothing.
	cancel()
	<-ctx.Done()
	if latest != nil {
		fmt.Println("updated from %s to %s", version, latest)
	}
}
```

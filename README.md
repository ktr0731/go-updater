# go-updater

## workflow
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
	u := updater.New(version, NewGitHubReleaseMeans("ktr0731", "evans"))

	// in default, update if minor update found
	u.UpdateIf = FoundPatchUpdate

	if u.Updatable() {
		latest, err := u.Update()
		if err != nil {
			panic(err)
		}
		fmt.Println("update from %s to %s", version, latest)
	} else {
		fmt.Println("already latest")
	}
}
```

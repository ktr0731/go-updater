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
	u := updater.New("ktr0731", "evans", version)

	// in default, update if minor update found
	u.UpdateIf = FoundPatchUpdate

	if u.Updatable() {
		// determine what use means for update this software
		// in this example, use HomeBrew
		latest, err := u.UpdateBy(updater.HomeBrew)
		if err != nil {
			panic(err)
		}
		fmt.Println("update from %s to %s", version, latest)
	} else {
		fmt.Println("already latest")
	}
}
```

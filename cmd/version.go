package cmd

import (
	"fmt"
	"runtime"
)

var (
	Commit    = "unknown-commit"
	BuildTime = "unknown-buildtime"

	Version = "0.2.0"
)

// VersionFunc holds the textual version string.
func VersionFunc() string {
	return fmt.Sprintf(": %s\ncommit: %s\nbuild time: %s\ngolang version: %s\n",
		Version, Commit, BuildTime, runtime.Version()+" "+runtime.GOOS+"/"+runtime.GOARCH)
}
